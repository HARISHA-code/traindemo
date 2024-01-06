// server/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"

	pb "github.com/HARISHA-code/ticket/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedTicketServiceServer
	mu              sync.Mutex
	seatAllocations map[string]pb.SeatAllocation
}

// PurchaseTicket implements the PurchaseTicket gRPC method
func (s *server) PurchaseTicket(ctx context.Context, req *pb.TicketRequest) (*pb.TicketResponse, error) {
	price := 20.0
	receipt := generateReceipt(req.From, req.To, req.UserFirstName, req.UserLastName, req.UserEmail, price, pb.SeatAllocation{})

	section := randomSection()
	seatNumber := randomSeatNumber()

	seatAllocation := pb.SeatAllocation{
		UserEmail:  req.UserEmail,
		Section:    section,
		SeatNumber: seatNumber,
	}

	s.mu.Lock()
	s.seatAllocations[req.UserEmail] = seatAllocation
	s.mu.Unlock()

	return &pb.TicketResponse{Receipt: receipt}, nil
}

// GetReceipt implements the GetReceipt gRPC method
func (s *server) GetReceipt(ctx context.Context, req *pb.ReceiptRequest) (*pb.ReceiptResponse, error) {
	s.mu.Lock()
	allocation, exists := s.seatAllocations[req.UserEmail]
	s.mu.Unlock()

	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	price := 20.0
	receipt := generateReceipt("London", "France", allocation.UserEmail, "", "", price, allocation)

	return &pb.ReceiptResponse{Receipt: receipt}, nil
}

// GetSeatAllocation implements the GetSeatAllocation gRPC method
func (s *server) GetSeatAllocation(ctx context.Context, req *pb.SeatRequest) (*pb.SeatResponse, error) {
	var seatAllocations []*pb.SeatAllocation

	s.mu.Lock()
	for _, allocation := range s.seatAllocations {
		if allocation.Section == req.Section {
			seatAllocations = append(seatAllocations, &allocation)
		}
	}
	s.mu.Unlock()

	return &pb.SeatResponse{SeatAllocation: seatAllocations}, nil
}

// RemoveUser implements the RemoveUser gRPC method
func (s *server) RemoveUser(ctx context.Context, req *pb.RemoveUserRequest) (*pb.RemoveUserResponse, error) {
	s.mu.Lock()
	_, exists := s.seatAllocations[req.UserEmail]
	if !exists {
		s.mu.Unlock()
		return nil, fmt.Errorf("user not found")
	}

	delete(s.seatAllocations, req.UserEmail)
	s.mu.Unlock()

	return &pb.RemoveUserResponse{Message: "User removed successfully"}, nil
}

// ModifySeat implements the ModifySeat gRPC method
func (s *server) ModifySeat(ctx context.Context, req *pb.ModifySeatRequest) (*pb.ModifySeatResponse, error) {
	s.mu.Lock()
	allocation, exists := s.seatAllocations[req.UserEmail]
	if !exists {
		s.mu.Unlock()
		return nil, fmt.Errorf("user not found")
	}

	allocation.Section = req.NewSection
	allocation.SeatNumber = req.NewSeatNumber

	s.seatAllocations[req.UserEmail] = allocation
	s.mu.Unlock()

	return &pb.ModifySeatResponse{Message: "Seat modified successfully"}, nil
}

// Additional helper functions

func generateReceipt(from, to, user, firstName, lastName string, price float64, allocation pb.SeatAllocation) string {
	return fmt.Sprintf("Receipt:\nFrom: %s\nTo: %s\nUser: %s %s %s\nPrice Paid: $%.2f\nSeat Allocation: Section %s, Seat %s", from, to, user, firstName, lastName, price, allocation.Section, allocation.SeatNumber)
}

func randomSection() string {
	sections := []string{"A", "B"}
	return sections[rand.Intn(len(sections))]
}

func randomSeatNumber() string {
	return strconv.Itoa(rand.Intn(10) + 1)
}

// Function to generate random string
func randSeq(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	// Initialize gRPC server
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTicketServiceServer(s, &server{seatAllocations: make(map[string]pb.SeatAllocation)})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Printf("Server listening on :50051")
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

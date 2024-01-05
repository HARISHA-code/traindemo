package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"traindemo/pb"

	//pb "github.com/HARISHA-code/traindemo/pb"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
)

var (
	users   = make(map[string]*pb.User)
	seatMap = make(map[string]*pb.Seat)
)

type trainServer struct{}

func (s *trainServer) PurchaseTicket(ctx context.Context, req *pb.PurchaseRequest) (*pb.Receipt, error) {
	user := req.GetUser()
	price := 20.0
	seat := allocateSeat(user.GetFirstName(), user.GetLastName(), req.GetSection())
	receipt := &pb.Receipt{
		From:      req.GetFrom(),
		To:        req.GetTo(),
		User:      user,
		PricePaid: price,
		Seat:      seat,
		Time:      ptypes.TimestampNow(),
	}
	users[user.GetEmail()] = user
	return receipt, nil
}

func (s *trainServer) GetReceipt(ctx context.Context, req *pb.User) (*pb.Receipt, error) {
	user, exists := users[req.GetEmail()]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	seat, exists := seatMap[user.GetEmail()]
	if !exists {
		return nil, fmt.Errorf("seat not found")
	}
	return &pb.Receipt{
		User: user,
		Seat: seat,
	}, nil
}

func (s *trainServer) GetSectionDetails(ctx context.Context, req *pb.SectionRequest) (*pb.SectionDetails, error) {
	section := req.GetSection()
	var sectionUsers []*pb.User
	for _, user := range users {
		seat, exists := seatMap[user.GetEmail()]
		if exists && seat.GetSection() == section {
			sectionUsers = append(sectionUsers, user)
		}
	}
	return &pb.SectionDetails{
		Section: section,
		Users:   sectionUsers,
	}, nil
}

func (s *trainServer) RemoveUser(ctx context.Context, req *pb.User) (*pb.RemoveUserResponse, error) {
	email := req.GetEmail()
	delete(users, email)
	delete(seatMap, email)
	return &pb.RemoveUserResponse{
		Success: true,
	}, nil
}

func (s *trainServer) ModifySeat(ctx context.Context, req *pb.ModifySeatRequest) (*pb.ModifySeatResponse, error) {
	email := req.GetUser().GetEmail()
	section := req.GetSection()
	newSeat := allocateSeat(req.GetUser().GetFirstName(), req.GetUser().GetLastName(), section)
	seatMap[email] = newSeat
	return &pb.ModifySeatResponse{
		Success: true,
		NewSeat: newSeat,
	}, nil
}

func allocateSeat(firstName, lastName, section string) *pb.Seat {
	seat := &pb.Seat{
		Section: section,
		Number:  len(seatMap) + 1,
	}
	seatMap[fmt.Sprintf("%s_%s", firstName, lastName)] = seat
	return seat
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTrainServiceServer(s, &trainServer{})

	log.Println("Server is listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

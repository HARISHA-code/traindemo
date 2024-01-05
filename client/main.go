package main

import (
	"context"
	"fmt"
	"log"

	//"traindemo/pb"

	pb "github.com/HARISHA-code/traindemo/pb"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTrainServiceClient(conn)

	// Purchase Ticket
	purchaseReq := &pb.PurchaseRequest{
		From:    "London",
		To:      "France",
		Section: "A",
		User: &pb.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		},
	}
	receipt, err := client.PurchaseTicket(context.Background(), purchaseReq)
	if err != nil {
		log.Fatalf("Error purchasing ticket: %v", err)
	}
	fmt.Printf("Purchase Receipt: %+v\n", receipt)

	// Get Receipt
	getReceiptReq := &pb.User{
		Email: "john.doe@example.com",
	}
	receipt, err = client.GetReceipt(context.Background(), getReceiptReq)
	if err != nil {
		log.Fatalf("Error getting receipt: %v", err)
	}
	fmt.Printf("Get Receipt: %+v\n", receipt)

	// Get Section Details
	sectionDetailsReq := &pb.SectionRequest{
		Section: "A",
	}
	sectionDetails, err := client.GetSectionDetails(context.Background(), sectionDetailsReq)
	if err != nil {
		log.Fatalf("Error getting section details: %v", err)
	}
	fmt.Printf("Section Details: %+v\n", sectionDetails)

	// Remove User
	removeUserReq := &pb.User{
		Email: "john.doe@example.com",
	}
	removeUserRes, err := client.RemoveUser(context.Background(), removeUserReq)
	if err != nil {
		log.Fatalf("Error removing user: %v", err)
	}
	fmt.Printf("Remove User Response: %+v\n", removeUserRes)

	// Modify Seat
	modifySeatReq := &pb.ModifySeatRequest{
		User: &pb.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		},
		Section: "B",
	}
	modifySeatRes, err := client.ModifySeat(context.Background(), modifySeatReq)
	if err != nil {
		log.Fatalf("Error modifying seat: %v", err)
	}
	fmt.Printf("Modify Seat Response: %+v\n", modifySeatRes)
}

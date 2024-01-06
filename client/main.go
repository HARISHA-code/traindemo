// client/main.go
package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/HARISHA-code/ticket/proto"
	"google.golang.org/grpc"
)

func main() {
	// Set up gRPC connection
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTicketServiceClient(conn)

	// Get user input for first name, last name, and email
	firstName := getUserInput("Enter your first name: ")
	lastName := getUserInput("Enter your last name: ")
	email := getUserInput("Enter your email address: ")

	// Example 1: Purchase a ticket
	receipt, err := purchaseTicket(c, "London", "France", firstName, lastName, email)
	if err != nil {
		log.Fatalf("Error purchasing ticket: %v", err)
	}
	fmt.Println("Receipt Details:")
	fmt.Println(receipt)

	// Example 2: Get receipt details
	receiptDetails, err := getReceipt(c, email)
	if err != nil {
		log.Fatalf("Error getting receipt details: %v", err)
	}
	fmt.Println("Receipt Details:")
	fmt.Println(receiptDetails)

	// Example 3: Get seat allocation details
	seatAllocations, err := getSeatAllocation(c, "A")
	if err != nil {
		log.Fatalf("Error getting seat allocation details: %v", err)
	}
	fmt.Println("Seat Allocation Details:")
	for _, allocation := range seatAllocations {
		fmt.Printf("User: %s, Section: %s, Seat: %s\n", allocation.UserEmail, allocation.Section, allocation.SeatNumber)
	}

	// Example 4: Remove a user from the train
	removeUserResponse, err := removeUser(c, email)
	if err != nil {
		log.Fatalf("Error removing user: %v", err)
	}
	fmt.Println("Remove User Response:")
	fmt.Println(removeUserResponse.Message)

	// Example 5: Modify a user's seat
	modifySeatResponse, err := modifySeat(c, email, "B", "5")
	if err != nil {
		log.Fatalf("Error modifying user's seat: %v", err)
	}
	fmt.Println("Modify Seat Response:")
	fmt.Println(modifySeatResponse.Message)

	fmt.Println("Client completed.")
}

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)
	return input
}

func purchaseTicket(client pb.TicketServiceClient, from, to, firstName, lastName, email string) (string, error) {
	response, err := client.PurchaseTicket(context.Background(), &pb.TicketRequest{
		From:          from,
		To:            to,
		UserFirstName: firstName,
		UserLastName:  lastName,
		UserEmail:     email,
	})
	if err != nil {
		return "", err
	}
	return response.Receipt, nil
}

func getReceipt(client pb.TicketServiceClient, email string) (string, error) {
	response, err := client.GetReceipt(context.Background(), &pb.ReceiptRequest{
		UserEmail: email,
	})
	if err != nil {
		return "", err
	}
	return response.Receipt, nil
}

func getSeatAllocation(client pb.TicketServiceClient, section string) ([]*pb.SeatAllocation, error) {
	response, err := client.GetSeatAllocation(context.Background(), &pb.SeatRequest{
		Section: section,
	})
	if err != nil {
		return nil, err
	}
	return response.SeatAllocation, nil
}

func removeUser(client pb.TicketServiceClient, email string) (*pb.RemoveUserResponse, error) {
	response, err := client.RemoveUser(context.Background(), &pb.RemoveUserRequest{
		UserEmail: email,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func modifySeat(client pb.TicketServiceClient, email, newSection, newSeatNumber string) (*pb.ModifySeatResponse, error) {
	response, err := client.ModifySeat(context.Background(), &pb.ModifySeatRequest{
		UserEmail:     email,
		NewSection:    newSection,
		NewSeatNumber: newSeatNumber,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

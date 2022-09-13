package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/ruancaetano/go-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect to gRPC Server %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	// AddUser(client)
	// AddUserVerbose(client)
	// AddUsers(client)
	AddUsersVerbose(client)
}

func AddUser(client pb.UserServiceClient) {
	request := &pb.User{
		Name:  "Ruan Caetano",
		Email: "ruansouza_caetano@hotmail.com",
	}

	response, err := client.AddUser(context.Background(), request)

	if err != nil {
		log.Fatalf("Could not make gRPC request %v", err)
	}

	fmt.Println(response)
}

func AddUserVerbose(client pb.UserServiceClient) {
	request := &pb.User{
		Name:  "Ruan Caetano",
		Email: "ruansouza_caetano@hotmail.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), request)
	if err != nil {
		log.Fatalf("Could not make gRPC request %v", err)
	}

	for {
		stream, err := responseStream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Could not receive the message %v", err)
		}

		fmt.Println("Status: ", stream.Status)
	}

}

func AddUsers(client pb.UserServiceClient) {

	usersToAdd := []*pb.User{
		{

			Email: "ruan1@caetano.com",
			Name:  "Ruan 1",
		},
		{
			Id:    "2",
			Email: "ruan2@caetano.com",
			Name:  "Ruan 2",
		},
		{
			Id:    "3",
			Email: "ruan3@caetano.com",
			Name:  "Ruan 3",
		},
		{
			Id:    "4",
			Email: "ruan4@caetano.com",
			Name:  "Ruan 4",
		},
		{
			Id:    "5",
			Email: "ruan5@caetano.com",
			Name:  "Ruan 5",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating requests %v", err)
	}

	for _, userToAdd := range usersToAdd {
		stream.Send(userToAdd)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error recieving response %v", err)
	}

	fmt.Println(res)

}

func AddUsersVerbose(client pb.UserServiceClient) {

	usersToAdd := []*pb.User{
		{

			Email: "ruan1@caetano.com",
			Name:  "Ruan 1",
		},
		{
			Email: "ruan2@caetano.com",
			Name:  "Ruan 2",
		},
		{
			Email: "ruan3@caetano.com",
			Name:  "Ruan 3",
		},
		{
			Email: "ruan4@caetano.com",
			Name:  "Ruan 4",
		},
		{
			Email: "ruan5@caetano.com",
			Name:  "Ruan 5",
		},
	}

	stream, err := client.AddUsersVerbose(context.Background())

	if err != nil {
		log.Fatalf("Could not make gRPC request %v", err)
	}

	wait := make(chan int)

	go func() {
		for _, userToAdd := range usersToAdd {
			fmt.Println("Sending user: ", userToAdd.Name)
			stream.Send(userToAdd)
			time.Sleep(time.Second * 3)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Could not receive the message %v", err)
			}

			fmt.Printf("Recebendo user %v com Status %v\n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}

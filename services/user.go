package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/ruancaetano/go-grpc/pb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (*UserService) AddUser(ctx context.Context, req *pb.User) (*pb.User, error) {

	newId, _ := uuid.NewUUID()

	return &pb.User{
		Id:    newId.String(),
		Name:  req.GetName(),
		Email: req.GetEmail(),
	}, nil
}

func (*UserService) AddUserVerbose(req *pb.User, stream pb.UserService_AddUserVerboseServer) error {

	stream.Send(&pb.UserResultStream{
		Status: "Init",
		User:   &pb.User{},
	})

	time.Sleep(time.Second * 3)

	newId, _ := uuid.NewUUID()

	stream.Send(&pb.UserResultStream{
		Status: "Id Generated",
		User: &pb.User{
			Id: newId.String(),
		},
	})

	time.Sleep(time.Second * 3)

	stream.Send(&pb.UserResultStream{
		Status: "User has been inserted",
		User: &pb.User{
			Id:    newId.String(),
			Name:  req.GetName(),
			Email: req.GetEmail(),
		},
	})

	time.Sleep(time.Second * 3)

	stream.Send(&pb.UserResultStream{
		Status: "Completed",
		User: &pb.User{
			Id:    newId.String(),
			Name:  req.GetName(),
			Email: req.GetEmail(),
		},
	})

	return nil
}

func (*UserService) AddUsers(stream pb.UserService_AddUsersServer) error {

	users := []*pb.User{}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Users{
				User: users,
			})
		}

		if err != nil {
			log.Fatalf("Error receivign stream: %v", err)
		}

		users = append(users, &pb.User{
			Id:    uuid.New().String(),
			Email: req.GetEmail(),
			Name:  req.GetName(),
		})
		fmt.Println("Adding ", req.GetName())

	}
}

func (*UserService) AddUsersVerbose(stream pb.UserService_AddUsersVerboseServer) error {

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error receiving stream: %v", err)
		}

		fmt.Println("Recebendo: ", req.GetName())
		err = stream.Send(&pb.UserResultStream{
			Status: "Added",
			User: &pb.User{
				Id:    uuid.New().String(),
				Name:  req.GetName(),
				Email: req.GetEmail(),
			},
		})

		if err != nil {
			log.Fatalf("Error sending stream: %v", err)
		}

	}
}

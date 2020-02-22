package v1

import (
	"context"

	"fmt"

	v1 "github.com/sajanjswl/file-service/pkg/api/v1"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

type fileServiceServer struct {
	db *mongo.Client
}

func NewFileServiceServer(db *mongo.Client) v1.FileServiceServer {
	return &fileServiceServer{db: db}
}

func (s *fileServiceServer) RegisterUser(ctx context.Context, req *v1.RegisterUserRequest) (*v1.RegisterUserResponse, error) {

	fmt.Println(req.GetUserDetails().GetName())

	return &v1.RegisterUserResponse{
		Message: "successfully tested",
	}, nil
}

package v1

import (
	"context"
	"log"
	"os"

	v1 "github.com/sajanjswl/file-service/pkg/api/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion     = "v1"
	fileDatabase   = os.Getenv("FILE_SERVICE_DATABASE")
	fileCollection = os.Getenv("FILE_SERVICE_COLLECTION")
)

type fileServiceServer struct {
	db *mongo.Client
}

// checkAPI checks if the API version requested by client is supported by server
func (f *fileServiceServer) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}

func NewFileServiceServer(db *mongo.Client) v1.FileServiceServer {
	return &fileServiceServer{db: db}
}

func (f *fileServiceServer) RegisterUser(ctx context.Context, req *v1.RegisterUserRequest) (*v1.RegisterUserResponse, error) {

	if err := s.checkAPI(req.GetApi()); err != nil {
		return nil, err
	}

	//bycrpting the plaint text password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.GetUser().GetPassword()), bcrypt.MinCost)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, config.InternalError)
	}
, err := bcrypt.GenerateFromPassword([]byte(req.GetUserDetails().GetPassword(), bcrypt.MinCost)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, config.InternalError)
	}

		password:=string(passwordHash)

	dbCol := f.db.Database(fileDatabase).Collection(fileCollection)
	_, err = dbCol.InsertOne(mongoCtx, bson.D{
		{Key: "name", Value: req.GetUserDetails().GetName()},
		{Key: "email", Value: req.GetUserDetails().GetEmail()},
		{Key: "password", Value: password },
	})

	if err != nil {
		log.Println(err)
		return false
	}

	return &v1.RegisterUserResponse{
		Message: "successfully tested",
	}, nil
}

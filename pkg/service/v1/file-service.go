package v1

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	v1 "github.com/sajanjswl/file-service/pkg/api/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
	// fileDatabase   = os.Getenv("FILE_SERVICE_DATABASE")
	// fileCollection = os.Getenv("FILE_SERVICE_COLLECTION")
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

	if err := f.checkAPI(req.GetApi()); err != nil {
		log.Error(err)
		return nil, err
	}

	//bycrpting the plaint text password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.GetUserDetails().GetPassword()), bcrypt.MinCost)
	log.Println("passwor", passwordHash)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	password := string(passwordHash)

	dbCol := f.db.Database(os.Getenv("FILE_SERVICE_DATABASE")).Collection(os.Getenv("FILE_SERVICE_COLLECTION"))
	_, err = dbCol.InsertOne(ctx, bson.D{
		{Key: "name", Value: req.GetUserDetails().GetName()},
		{Key: "email", Value: req.GetUserDetails().GetEmail()},
		{Key: "password", Value: password},
	})

	if err != nil {
		log.Error(err)
		return nil, err
	} else {
		return &v1.RegisterUserResponse{
			Message: "successfully tested",
		}, nil

	}

}

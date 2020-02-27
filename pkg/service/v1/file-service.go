package v1

import (
	"context"
	"errors"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	v1 "github.com/sajanjswl/file-service/pkg/api/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
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

		return nil

	}
	return errors.New("unsupported API version: Api version cannot be nil")
}

type User struct {
	Username         string `json: "username"`
	Password         string `json: "password"`
	FileUploadStatus bool   `json: "fileuploadstatus"`
}

func NewFileServiceServer(db *mongo.Client) v1.FileServiceServer {
	return &fileServiceServer{db: db}
}

func (f *fileServiceServer) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {

	if err := f.checkAPI(req.GetApi()); err != nil {
		log.Error(err)
		return nil, err
	}

	dbCol := f.db.Database(os.Getenv("FILE_SERVICE_DATABASE")).Collection(os.Getenv("FILE_SERVICE_COLLECTION"))

	//checking if useer already exists
	var result User
	err := dbCol.FindOne(ctx, bson.D{{"username", req.GetUsername()}}).Decode(&result)
	if err != nil {
		log.Error(err)
		if err.Error() == "mongo: no documents in result" {
			goto CREATE
		}

		return nil, errors.New("failed to register into database")
	}

	log.Printf("Found document without description: %+v\n", result)
	return nil, errors.New("user already exists")

CREATE:
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.MinCost)

	if err != nil {
		log.Error(err)
		return nil, errors.New("status code:503 Service Unavailable")
	}

	password := string(passwordHash)
	user := &User{
		Username:         req.GetUsername(),
		Password:         password,
		FileUploadStatus: false,
	}

	id, err := dbCol.InsertOne(ctx, user)

	if err != nil {
		log.Error(err)
		return nil, errors.New("failed to register into database")
	}
	log.Println("user created with ", *id)

	return &v1.CreateUserResponse{
		Message: "successfully registered..",
	}, nil

}

func (f *fileServiceServer) UploadFile(stream v1.FileService_UploadFileServer) error {


	// dbCol := f.db.Database(os.Getenv("FILE_SERVICE_DATABASE")).Collection(os.Getenv("FILE_SERVICE_COLLECTION"))

	// update := bson.D{
	// 	{"$set", bson.D{
	// 		{"fileuploadstatus", true},
	// 	}},
	// }
	// updateResult, err := dbCol.UpdateOne(stream.Context(), bson.D{{"id", "1WEB"}}, update)
	// if err != nil {
	// 	log.Error(err)
	// }

	// log.Printf("Updated %v document(s).\n", updateResult.ModifiedCount)

	// bucket, err := gridfs.NewBucket(
	// 	f.db.Database(os.Getenv("FILE_SERVICE_DATABASE")),
	// )
	// if err != nil {

	// 	log.Println(err)
	// 	return errors.New("status code:503 Service Unavailable")
	// }

	// uploadStream, err := bucket.OpenUploadStream(
	// 	"something",
	// )
	// if err != nil {
	// 	log.Println(err)
	// 	return errors.New("status code:503 Service Unavailable")
	// }
	// defer uploadStream.Close()

	for {
		_, err := stream.Recv()

		if err != nil {
			if err == io.EOF {
				goto END
			}
			return err

		}

		// fileSize, err := uploadStream.Write(b.GetContent())
		// if err != nil {
		// 	log.Println(err)
		// 	return errors.New("status code:503 Service Unavailable")
		// }
		// log.Printf("Write file to DB was successful. File size: %d M\n", fileSize)

	}

END:

	return stream.SendAndClose(&v1.UploadStatus{
		Message: "Upload received with success",
		Code:    v1.UploadStatusCode_Ok,
	})

}

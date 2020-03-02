package v1

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	v1 "github.com/sajanjswl/file-service/pkg/api/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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

func NewFileServiceServer(db *mongo.Client) v1.FileServiceServer {
	return &fileServiceServer{db: db}
}

func getConn(f *fileServiceServer, connType string) (*mongo.Database, *mongo.Collection) {

	if connType == "db" {

		return f.db.Database(os.Getenv("FILE_SERVICE_DATABASE")), nil
	}
	return nil, f.db.Database(os.Getenv("FILE_SERVICE_DATABASE")).Collection(os.Getenv("FILE_SERVICE_COLLECTION"))

}

func (f *fileServiceServer) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {

	if err := f.checkAPI(req.GetApi()); err != nil {
		log.Error(err)
		return nil, err
	}

	_, collection := getConn(f, "collection")

	//checking if useer already exists
	var results bson.M
	err := collection.FindOne(ctx, bson.D{{"username", req.GetUserDetails().GetUsername()}}).Decode(&results)

	log.Println(results)
	if err != nil {
		log.Error(err)
		if err.Error() == "mongo: no documents in result" {
			goto CREATE
		}

		return nil, errors.New("failed to register into database")
	}

	return nil, errors.New("user already exists")

CREATE:
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.GetUserDetails().GetPassword()), bcrypt.MinCost)

	if err != nil {
		log.Error(err)
		return nil, errors.New("status code:503 Service Unavailable")
	}

	password := string(passwordHash)

	id, err := collection.InsertOne(ctx, bson.D{
		{"username", req.GetUserDetails().GetUsername()}, {"password", password}, {"fileid", primitive.NilObjectID},
	})

	if err != nil {
		log.Error(err)
		return nil, errors.New("failed to register into database")
	}
	log.Println("user created with ", *id)

	return &v1.CreateUserResponse{
		Message: "successfully registered..",
	}, nil

}

func (f *fileServiceServer) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*v1.DeleteUserResponse, error) {
	if err := f.checkAPI(req.GetApi()); err != nil {
		log.Error(err)
		return nil, err
	}

	_, collection := getConn(f, "collection")
	//verifying user
	_, fileID, err := getLoginStatusAndFileID(ctx, collection, req.GetUserDetails().GetUsername(), req.GetUserDetails().GetPassword())

	if err != nil && err.Error() != "file doesnt exists" {
		return nil, err
	}

	res, err := collection.DeleteOne(ctx, bson.M{"username": req.GetUserDetails().GetUsername()})

	if err != nil {
		log.Error(err)
		return nil, errors.New("failed to delete user  " + req.GetUserDetails().GetUsername())
	}
	log.Println("deleted result", res)

	if fileID.IsZero() {
		return &v1.DeleteUserResponse{
			Message: "user" + req.GetUserDetails().GetUsername() + "deleted successfully",
		}, nil

	}

	db, _ := getConn(f, "db")
	bucket, err := gridfs.NewBucket(
		db,
	)

	if err := bucket.Delete(fileID); err != nil {
		log.Error(err)
		return nil, errors.New("failed to delete user " + req.GetUserDetails().GetUsername() + "files")
	}

	return &v1.DeleteUserResponse{
		Message: "user  " + req.GetUserDetails().GetUsername() + "deleted successfully",
	}, nil

}

func (f *fileServiceServer) UploadFile(stream v1.FileService_UploadFileServer) error {

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return stream.SendAndClose(&v1.UploadStatus{
			Message: "upload failed , failed to read credential",
			Code:    v1.UploadStatusCode_Failed,
		})
	}

	usernames := md["username"]
	passwords := md["password"]
	apVersions := md["apiversion"]

	if err := f.checkAPI(apVersions[0]); err != nil {
		log.Error(err)
		return stream.SendAndClose(&v1.UploadStatus{
			Message: "upload failed, " + err.Error(),
			Code:    v1.UploadStatusCode_Failed,
		})

	}
	_, collection := getConn(f, "collection")

	//verifying user
	_, fileIDReturned, err := getLoginStatusAndFileID(stream.Context(), collection, usernames[0], passwords[0])
	fileID := primitive.NewObjectID()

	if err != nil && err.Error() != "file doesnt exists" {

		return stream.SendAndClose(&v1.UploadStatus{
			Message: "upload failed, " + err.Error(),
			Code:    v1.UploadStatusCode_Failed,
		})

	}

	if !fileIDReturned.IsZero() {

		return stream.SendAndClose(&v1.UploadStatus{
			Message: "upload failed, File upload limit exceded ",
			Code:    v1.UploadStatusCode_Failed,
		})
	}

	db, _ := getConn(f, "db")
	bucket, err := gridfs.NewBucket(
		db,
	)
	if err != nil {

		log.Println(err)
		return errors.New("status code:503 Service Unavailable")
	}

	uploadStream, err := bucket.OpenUploadStreamWithID(
		fileID,
		usernames[0],
	)
	if err != nil {
		log.Println(err)
		return errors.New("status code:503 Service Unavailable")
	}
	defer uploadStream.Close()

	for {

		b, err := stream.Recv()

		if err != nil {
			if err == io.EOF {
				goto END
			}
			return err

		}

		fileSize, err := uploadStream.Write(b.GetContent())
		if err != nil {
			log.Println(err)
			return errors.New("status code:503 Service Unavailable")
		}
		log.Printf("Write file to DB was successful. File size: %d M\n", fileSize)

	}

END:

	update := bson.D{
		{"$set", bson.D{
			{"fileid", fileID},
		}},
	}

	updateResult, err := collection.UpdateOne(stream.Context(), bson.D{{"username", usernames[0]}}, update)
	if err != nil {
		log.Error(err)
	}

	log.Printf("Updated %v document(s).\n", updateResult.ModifiedCount)

	return stream.SendAndClose(&v1.UploadStatus{
		Message: "Upload received with success",
		Code:    v1.UploadStatusCode_Ok,
	})

}

func (f *fileServiceServer) DownloadFile(req *v1.DownloadFileRequest, stream v1.FileService_DownloadFileServer) error {

	if err := f.checkAPI(req.GetApi()); err != nil {
		log.Error(err)
		return err

	}

	_, collection := getConn(f, "collection")

	//verifying user
	_, fileID, err := getLoginStatusAndFileID(stream.Context(), collection, req.GetUserDetails().GetUsername(), req.GetUserDetails().GetPassword())

	if err != nil {
		return err
	}
	db, _ := getConn(f, "db")
	bucket, err := gridfs.NewBucket(
		db,
	)

	if err != nil {

		log.Println(err)
		return errors.New("status code:503 Service Unavailable")
	}

	var buf bytes.Buffer
	//	var buf []byte
	dStream, err := bucket.DownloadToStream(fileID, &buf)

	if err != nil {
		log.Println(err)
		return errors.New("status code:503 Service Unavailable")
	}
	log.Printf("File size to download: %v \n", dStream)

	err = stream.Send(&v1.Chunk{

		Content: buf.Bytes(),
	})

	if err != nil {

		log.Error(stream, err)
		return err
	}

	log.Println("file sent successfully")

	return nil

}

func getLoginStatusAndFileID(ctx context.Context, collection *mongo.Collection, username string, password string) (bool, primitive.ObjectID, error) {

	var results bson.M

	err := collection.FindOne(ctx, bson.D{{"username", username}}).Decode(&results)
	if err != nil {
		log.Error(err)
		if err.Error() == "mongo: no documents in result" {
			return false, primitive.NilObjectID, errors.New("user doesnt exists")
		}

		return false, primitive.NilObjectID, errors.New("sevice unavailable")

	}

	var fileID primitive.ObjectID
	var userPassword string

	for i, v := range results {

		if i == "fileid" {

			fileID = v.(primitive.ObjectID)

		}
		if i == "password" {
			userPassword = v.(string)
		}

	}

	err1 := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password))

	if err1 != nil {
		log.Error(err1)
		return false, primitive.NilObjectID, errors.New("invalid credentials")
	}

	if fileID.IsZero() {
		return true, fileID, errors.New("file doesnt exists")
	}

	return true, fileID, nil
}

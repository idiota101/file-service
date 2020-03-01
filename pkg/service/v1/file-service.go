package v1

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"

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

type User struct {
	//	ID       primitive.ObjectID `bson:"_id"`
	Username string             `json: "username"`
	Password string             `json: "password"`
	fileID   primitive.ObjectID `json: "fileid"`
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
	var result User
	err := collection.FindOne(ctx, bson.D{{"username", req.GetUserDetails().GetUsername()}}).Decode(&result)
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
	some := primitive.NewObjectID()
	log.Println("some   ", some)

	// update := bson.D{
	// 	{"$set", bson.D{
	// 		{"fileid", fileId},
	// 	}},
	// }

	password := string(passwordHash)

	//bson.D{{"username", req.GetUserDetails().GetUsername()}, {"password", password}, {"fileid", some},}

	user := &User{
		Username: req.GetUserDetails().GetUsername(),
		Password: password,
		fileID:   some,
	}
	log.Println("new object id", user.fileID)

	//id, err := collection.InsertOne(ctx, user)

	id, err := collection.InsertOne(ctx, bson.D{
		{"username", req.GetUserDetails().GetUsername()}, {"password", password}, {"fileid", some},
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

	dbCol := f.db.Database(os.Getenv("FILE_SERVICE_DATABASE")).Collection(os.Getenv("FILE_SERVICE_COLLECTION"))

	//verifying user
	loginstatus, fileExists, err := loginAndFileExistsStatus(ctx, dbCol, req.GetUsername(), req.GetPassword())

	if !loginstatus {
		return nil, err
	}

	// res, err := dbCol.DeleteOne(ctx, bson.M{"username": req.GetUsername})
	// if err != nil {
	// 	log.Error(err)
	// 	return nil, err
	// }
	// log.Println(res)

	if !fileExists {
		goto END
	} else {

	}

END:

	db := f.db.Database(os.Getenv("FILE_SERVICE_DATABASE"))

	filefind(ctx, db)

	return &v1.DeleteUserResponse{
		Message: "user" + req.GetUsername() + "deleted successfully",
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
	// apVersions := md["apiVersion"]

	// log.Println(apVersions)
	// if err := f.checkAPI(apVersions[0]); err != nil {
	// 	log.Error(err)
	// 	return err
	// }
	_, collection := getConn(f, "collection")

	//verifying user
	loginstatus, fileExists, err := loginAndFileExistsStatus(stream.Context(), collection, usernames[0], passwords[0])

	log.Println(loginstatus, fileExists)
	//if !loginstatus || fileExists {
	if true {
		return stream.SendAndClose(&v1.UploadStatus{
			//	Message: "upload failed, " + err.Error(),
			Code: v1.UploadStatusCode_Failed,
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
	fileId := primitive.NewObjectID()
	uploadStream, err := bucket.OpenUploadStreamWithID(
		fileId,
		usernames[0],
	)
	if err != nil {
		log.Println(err)
		return errors.New("status code:503 Service Unavailable")
	}
	defer uploadStream.Close()

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

	update := bson.D{
		{"$set", bson.D{
			{"fileid", fileId},
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

	dbCol := f.db.Database(os.Getenv("FILE_SERVICE_DATABASE")).Collection(os.Getenv("FILE_SERVICE_COLLECTION"))

	//verifying user
	loginstatus, fileExists, err := loginAndFileExistsStatus(stream.Context(), dbCol, req.GetUsername(), req.GetPassword())

	if !loginstatus || !fileExists {
		return err
	}

	bucket, err := gridfs.NewBucket(
		f.db.Database(os.Getenv("FILE_SERVICE_DATABASE")),
	)

	if err != nil {

		log.Println(err)
		return errors.New("status code:503 Service Unavailable")
	}

	var buf bytes.Buffer
	//	var buf []byte
	dStream, err := bucket.DownloadToStreamByName(req.GetUsername(), &buf)
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

func loginAndFileExistsStatus(ctx context.Context, collection *mongo.Collection, username string, password string) (bool, bool, error) {

	var results bson.M
	// err := fsFiles.FindOne(mongoCtx, bson.M{}).Decode(&results)
	// if err != nil {
	// 	log.Println(err)
	// 	return false
	// }
	//var user User
	err := collection.FindOne(ctx, bson.D{{"username", username}}).Decode(&results)
	if err != nil {
		log.Error(err)
		if err.Error() == "mongo: no documents in result" {
			return false, false, errors.New("user doesnt exists")
		}

	}

	//type a interface {}
	var c primitive.ObjectID
	// func(a  interface{}){
	// 	c =reflect.ValueOf(a)
	// }

	//var a bson.PrimitiveCodecs
	//results.

	for i, v := range results {

		if i == "fileid" {

			a := reflect.ValueOf(v)
			log.Printf("inde"+i+"value %T", v)
			log.Println(a, "pritin value of a")
			s := v.(primitive.ObjectID)
			log.Println(v.(primitive.ObjectID), "pritin value of a")

			log.Println(s, "Printing s")
			log.Println(s.IsZero(), "Printing zero")
			log.Println(c.IsZero(), "c  Printing zero")
		}

	}

	log.Println(c)

	// err1 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	// if err1 != nil {
	// 	log.Error(err1)
	// 	return false, false, errors.New("invalid credentials")
	// }

	// log.Println("printing file size", user.fileID.IsZero())

	// log.Println(user.fileID.MarshalJSON())

	// log.Println(user)

	// if !user.fileID.IsZero() {

	// 	return true, true, errors.New("file upload limit exceded")
	// }

	return true, false, nil
}

func filefind(ctx context.Context, db *mongo.Database) {

	bucket, err := gridfs.NewBucket(
		db,
	)
	if err != nil {

		log.Println(err)

	}

	log.Println("i am in find")

	filter := bson.D{
		{"length", bson.D{{"$gt", 100}}},
	}
	cursor, err := bucket.Find(filter)

	defer func() {
		if err := cursor.Close(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	type gridfsFile struct {
		Name   string `bson:"filename"`
		Length int64  `bson:"length"`

		ID primitive.ObjectID `bson:"_id"`
	}
	var foundFiles []gridfsFile
	if err = cursor.All(context.TODO(), &foundFiles); err != nil {
		log.Fatal(err)
	}

	log.Println("printing files", foundFiles)
	for _, file := range foundFiles {
		fmt.Printf("filename: %s, length: %d\n", file.Name, file.Length, file.ID)

		if err := bucket.Delete(file.ID); err != nil {
			log.Fatal(err)
		}

		log.Println("successfully deleted files")

	}

}

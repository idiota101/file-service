package cmd

import (
	"context"
	"flag"
	"fmt"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sajanjswl/file-service/pkg/protocol/grpc"

	v1 "github.com/sajanjswl/file-service/pkg/service/v1"

	"os"
	//"fmt"
)

type Config struct {
	GRPCPort string
}

func RunServer() error {

	mongoUri := fmt.Sprintf("mongodb://%s:%s@%s:%s", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))

	ctx := context.Background()

	// Initialize MongoDb client
	fmt.Println("Connecting to MongoDB...")

	db, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))

	err = db.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v\n", err)
	} else {
		fmt.Println("Connected to Mongodb!!")
	}

	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.Parse()
	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	v1API := v1.NewFileServiceServer(db)
	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}

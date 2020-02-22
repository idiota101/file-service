package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc/testdata"

	v1 "github.com/sajanjswl/file-service/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func RunServer(ctx context.Context, v1API v1.FileServiceServer, port string) error {
	listen, err := net.Listen(os.Getenv("GRPC_NETWORK_TYPE"), ":"+port)
	if err != nil {
		return err
	}
	// Create tls based credential.
	creds, err := credentials.NewServerTLSFromFile(testdata.Path("server1.pem"), testdata.Path("server1.key"))
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}
	// register service
	server := grpc.NewServer(grpc.Creds(creds))
	v1.RegisterFileServiceServer(server, v1API)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Println("starting gRPC server...")
	return server.Serve(listen)
}

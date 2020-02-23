package main

import (
	"context"
	"flag"
	"time"

	v1 "github.com/sajanjswl/file-service/pkg/api/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

func main() {

	address := flag.String("server", "", "gRPC server in format host:port")
	flag.Parse()

	creds, err := credentials.NewClientTLSFromFile(testdata.Path("ca.pem"), "x.test.youtube.com")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	conn, err := grpc.Dial(*address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := v1.NewFileServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//register
	// req := registeration()
	// resp, err := c.Register(ctx, req)

	//Login
	//req := registeration()
	resp, err := c.RegisterUser(ctx, &v1.RegisterUserRequest{
		Api: apiVersion,
		UserDetails: &v1.UserDetails{
			Name:     "sajan",
			Email:    "sjnjaiswalfhgfjfg",
			Password: "sjfddnf",
		},
	})

	if err != nil {
		log.Println(err)
	} else {

		log.Println(resp)
	}

}

func registeration() *v1.RegisterUserRequest {
	req := &v1.RegisterUserRequest{
		UserDetails: &v1.UserDetails{
			Name:     "sajan",
			Email:    "sjnjaiswaldfd",
			Password: "sjfddnf",
		},
	}

	return req
}

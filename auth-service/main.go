package main

import (
	"log"
	"net"
	"net/http"
	"sync"

	"egnite.app/microservices/authentication/services/login"
	"egnite.app/microservices/authentication/services/register"
	"egnite.app/microservices/authentication/services/token"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func startGRPC() {

	lis, err := net.Listen("tcp", "localhost:9001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	loginServer := login.Server{}
	login.RegisterLoginServiceServer(grpcServer, &loginServer)

	registerServer := register.Server{}
	register.RegisterRegisterServiceServer(grpcServer, &registerServer)

	tokenServer := token.Server{}
	token.RegisterTokenServiceServer(grpcServer, &tokenServer)
	log.Println("gRPC server ready...")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func startHTTP() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Connect to the GRPC server
	conn, err := grpc.Dial("localhost:9001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	// Register grpc-gateway
	rmux := runtime.NewServeMux()

	loginClient := login.NewLoginServiceClient(conn)
	err = login.RegisterLoginServiceHandlerClient(ctx, rmux, loginClient)
	if err != nil {
		log.Fatal(err)
	}

	registerClient := register.NewRegisterServiceClient(conn)
	err = register.RegisterRegisterServiceHandlerClient(ctx, rmux, registerClient)
	if err != nil {
		log.Fatal(err)
	}

	tokenClient := token.NewTokenServiceClient(conn)
	err = token.RegisterTokenServiceHandlerClient(ctx, rmux, tokenClient)
	if err != nil {
		log.Fatal(err)
	}
	handler := cors.Default().Handler(rmux)
	log.Println("REST server ready...")
	err = http.ListenAndServe("localhost:8001", handler)

	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go startHTTP()
	go startGRPC()

	wg.Wait()

}

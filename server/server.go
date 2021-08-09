package main

import (
	"context"
	"flag"
	"fmt"
	"gRPCDemo/gen/productpb"
	"gRPCDemo/server/catalogservice"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	initServer()

	// server graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	log.Print("Server Stopped")
}

func initServer() {
	fmt.Println("Starting The Server")
	grpcServerAddress := flag.String("grpc-address", "", "grpc server address")
	restServerAddress := flag.String("rest-address", "", "rest server address")
	flag.Parse()
	//grpcEndPoint := "0.0.0.0:8080"
	//fmt.Printf("grpc address %s\n", *grpcServerAddress)
	//restStreamingEndpoint := "0.0.0.0:8081"
	gRPCListener, err := net.Listen("tcp", *grpcServerAddress)
	if err != nil {
		log.Fatalf("failed To create grpc listener: %+v", err)
	}

	restListener, err := net.Listen("tcp", *restServerAddress)
	if err != nil {
		log.Fatalf("failed To create rest listener: %+v", err)
	}

	productService := &catalogservice.ProductService{}

	go initGRPCServer(gRPCListener, productService)
	go initRestServer(gRPCListener, productService)
	go initRestStreamingServer(restListener, *grpcServerAddress)
}

func initGRPCServer(
	listener net.Listener,
	productService *catalogservice.ProductService,
) {
	fmt.Printf("Starting The gRPC Server at %s\n", listener.Addr().String())

	server := grpc.NewServer()

	productpb.RegisterProductServiceServer(server, productService)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %+v", err)
	}
}

func initRestServer(
	listener net.Listener,
	productService *catalogservice.ProductService,
) {

	fmt.Printf("Starting The Rest Server at %s\n", listener.Addr().String())

	mux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := productpb.RegisterProductServiceHandlerServer(ctx, mux, productService)
	if err != nil {
		log.Fatalf("unable to register rest server")
		return
	}

	if err := http.Serve(listener, mux); err != nil {
		log.Fatalf("failed to serve: %+v", err)
	}
}

func initRestStreamingServer(
	listener net.Listener,
	grpcEndPoint string,
) {

	fmt.Printf("Starting The Rest Streaming Server at %s\n", listener.Addr().String())

	mux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dialOptions := []grpc.DialOption{grpc.WithInsecure()}

	err := productpb.RegisterProductServiceHandlerFromEndpoint(ctx, mux, grpcEndPoint, dialOptions)
	if err != nil {
		log.Fatalf("unable to register rest server")
		return
	}

	if err := http.Serve(listener, mux); err != nil {
		log.Fatalf("failed to serve: %+v", err)
	}
}

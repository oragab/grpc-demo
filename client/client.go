package main

import (
	"context"
	"flag"
	"fmt"
	"gRPCDemo/gen/productpb"
	petname "github.com/dustinkirkland/golang-petname"
	"google.golang.org/grpc"
	"io"
	"log"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	targetServer := flag.String("target-server", "", "server")
	flag.Parse()

	fmt.Println("Starting the Catalog Client")
	conn, err := grpc.Dial(*targetServer,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(Unary()),
		grpc.WithStreamInterceptor(StreamInterceptor()),
	)
	if err != nil {
		log.Fatalf("error creating connection:%+v", err)
	}
	defer conn.Close()

	client := productpb.NewProductServiceClient(conn)
	fmt.Println("CatalogClient Created")

	// unary operation
	createProduct(client)

	// client streaming
	createProductsClientStreaming(client)

	// server streaming
	createProductsServerStreaming(client)

	// bidi
	createProductsBiDi(client)
}

func Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		log.Printf("--> unary interceptor: %s", method)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func StreamInterceptor() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		log.Printf("--> stream interceptor: %s", method)

		return streamer(ctx, desc, cc, method, opts...)
	}
}

func createProduct(client productpb.ProductServiceClient) {
	fmt.Println("Create Product Request")
	request := &productpb.CreateProductRequest{
		Product: getRandomProduct(),
	}

	response, err := client.CreateProduct(context.Background(), request)
	if err != nil {
		log.Fatalf("error sending request to the server: %+v\n", err)
		return
	}

	fmt.Printf("Created Product Successfully productID: %s\n", response.GetProductId())
}

func createProductsClientStreaming(client productpb.ProductServiceClient) {
	fmt.Println("Started createProductsClientStreaming RPC...")
	stream, err := client.CreateProductsClientStream(context.Background())
	if err != nil {
		log.Fatalf("error creating stream to the server: %+v\n", err)
		return
	}

	createProductRequests := []*productpb.CreateProductRequest{
		{Product: getRandomProduct()},
		{Product: getRandomProduct()},
		{Product: getRandomProduct()},
		{Product: getRandomProduct()},
	}

	for _, request := range createProductRequests {
		fmt.Printf("Sending Product: %+v\n", request.GetProduct())
		err = stream.Send(request)
		if err != nil {
			log.Fatalf("error sending message to the server: %+v\n", err)
			return
		}
		time.Sleep(1 * time.Second)
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error receving response from the server: %+v\n", err)
		return
	}

	fmt.Printf("Recevied ProductIds: %+v\n", response.GetProductIds())
}

func createProductsServerStreaming(client productpb.ProductServiceClient) {
	fmt.Println("Started createProductsServerStreaming RPC...")

	request := &productpb.CreateProductsRequest{
		Products: []*productpb.Product{
			getRandomProduct(),
			getRandomProduct(),
			getRandomProduct(),
			getRandomProduct(),
			getRandomProduct(),
		},
	}

	stream, err := client.CreateProductsServerStream(context.Background(), request)
	if err != nil {
		log.Fatalf("error sending request to the server: %+v\n", err)
		return
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("Stream Results Ended")
			return
		}

		if err != nil {
			log.Fatalf("error receiving response from the server stream: %+v\n", err)
			return
		}

		fmt.Printf("Recevied Response from the server created Product With ID: %+v\n", response.GetProductId())
	}
}

func createProductsBiDi(client productpb.ProductServiceClient) {
	fmt.Println("Started createProductsBiDi RPC...")
	stream, err := client.CreateProductsBiDi(context.Background())
	if err != nil {
		log.Fatalf("error creating stream with the server: %+v\n", err)
	}

	requests := []*productpb.CreateProductRequest{
		{Product: getRandomProduct()},
		{Product: getRandomProduct()},
		{Product: getRandomProduct()},
		{Product: getRandomProduct()},
		{Product: getRandomProduct()},
	}

	// sending stream
	go func() {
		for _, request := range requests {
			err := stream.Send(request)
			if err != nil {
				log.Fatalf("error sending request to the server: %+v\n", err)
				break
			}
			time.Sleep(600 * time.Millisecond)
		}
		err := stream.CloseSend()
		if err != nil {
			log.Fatalf("error closing send stream to the server: %+v\n", err)
		}
	}()

	waitChannel := make(chan int)
	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("error receiving response from the server stream: %+v\n", err)
				break
			}

			fmt.Printf("recevied response from the server stream product created ProductID: %s\n", response.GetProductId())
		}

		fmt.Println("Response stream ended from the server")
		waitChannel <- 1
	}()

	<-waitChannel
	fmt.Println("createProductsBiDi ended")
}

func getRandomProduct() *productpb.Product {
	min, max := 10, 250
	price := float32(min) + rand.Float32()*float32(max-min)
	maxOrder := int32(rand.Intn(250-10+1) + 10)
	return &productpb.Product{
		Name:     petname.Generate(2, "-"),
		Active:   false,
		Price:    price,
		MaxOrder: maxOrder,
	}
}

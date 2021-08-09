package catalogservice

import (
	"context"
	"fmt"
	"gRPCDemo/gen/productpb"
	"github.com/google/uuid"
	"io"
	"log"
	"time"
)

type ProductService struct{}

func (s *ProductService) CreateProduct(
	_ context.Context,
	request *productpb.CreateProductRequest,
) (*productpb.CreateProductResponse, error) {
	fmt.Println("Received CreateProduct gRPC")

	product := request.GetProduct()

	fmt.Printf("Received CreateProduct gRPC with Product: %+v\n", product)

	response := &productpb.CreateProductResponse{
		ProductId: uuid.NewString(),
	}

	return response, nil
}

func (s *ProductService) CreateProductsClientStream(
	stream productpb.ProductService_CreateProductsClientStreamServer,
) error {
	fmt.Println("Received CreateProductsClientStream gRPC")
	var productIds []string
	for {
		request, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("error receving message from client stream")
			return err
		}

		fmt.Printf("CreateProductsClientStream Product:%+v\n", request.GetProduct())

		productIds = append(productIds, uuid.NewString())
	}

	return stream.SendAndClose(&productpb.CreateProductsResponse{
		ProductIds: productIds,
	})
}

func (s *ProductService) CreateProductsServerStream(
	request *productpb.CreateProductsRequest,
	stream productpb.ProductService_CreateProductsServerStreamServer,
) error {
	fmt.Println("Received CreateProductsServerStream gRPC")
	for _, product := range request.GetProducts() {
		fmt.Printf("CreateProductsServerStream product:%+v\n", product)

		time.Sleep(1 * time.Second)

		err := stream.Send(&productpb.CreateProductResponse{ProductId: uuid.NewString()})
		if err != nil {
			fmt.Println("error sending message to the client")
			return err
		}
	}
	return nil
}

func (s *ProductService) CreateProductsBiDi(
	stream productpb.ProductService_CreateProductsBiDiServer,
) error {
	fmt.Println("Received CreateProductsServerStream gRPC")
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("error receving message from client: %+v", err)
			return err
		}

		product := request.GetProduct()
		fmt.Printf("CreateProductsServerStream Product:%+v\n", product)

		time.Sleep(600 * time.Millisecond)

		err = stream.Send(&productpb.CreateProductResponse{ProductId: uuid.NewString()})
		if err != nil {
			log.Fatalf("error sending message to the client:%+v\n", err)
			return err
		}
	}
}

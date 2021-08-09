
start-server-1:
	go run server/server.go -grpc-address 0.0.0.0:50051

start-server-2:
	go run server/server.go -grpc-address 0.0.0.0:50052

start-server-3:
	go run server/server.go -grpc-address 0.0.0.0:50053

start-server-4:
	go run server/server.go -grpc-address 0.0.0.0:50054

start-client:
	go run client/client.go -target-server localhost:9000

re-gen:
	protoc -I . --proto_path=proto proto/*.proto --go_out=plugins=grpc:gen --grpc-gateway_out gen --openapiv2_out openapi


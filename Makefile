
start-server-1:
	go run server/server.go -grpc-address 0.0.0.0:50051 -rest-address 0.0.0.0:8083

start-server-2:
	go run server/server.go -grpc-address 0.0.0.0:50052 -rest-address 0.0.0.0:8084

start-client:
	go run client/client.go -target-server localhost:50051

re-gen:
	protoc -I . --proto_path=proto proto/*.proto --go_out=plugins=grpc:gen --grpc-gateway_out gen --openapiv2_out openapi



start-server-1:
	go run server/server.go -grpc-address 0.0.0.0:8080 -rest-address 0.0.0.0:8081

start-server-2:
	go run server/server.go -grpc-address 0.0.0.0:8083 -rest-address 0.0.0.0:8081

start-client:
	go run client/client.go

re-gen:
	protoc -I . --proto_path=proto proto/*.proto --go_out=plugins=grpc:gen --grpc-gateway_out gen --openapiv2_out openapi


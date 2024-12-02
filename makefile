 
clean:
	rm proto/*.pb.go

gen:
	protoc --go_out=./proto --go_opt=paths=source_relative \
	--go-grpc_out=./proto --go-grpc_opt=paths=source_relative train.proto

grpc-server:
	go run ./server/main.go

grpc-client:
	go run ./client/main.go

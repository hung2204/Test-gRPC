gen :
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative calculator/calculator.proto

ser :
	go run .\server\server.go

cli :
	go run .\client\client.go
run:
	@go run cmd/app/main.go

build:
	@swag init -g cmd/app/main.go
	@go build -o bin/app cmd/app/main.go
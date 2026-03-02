DB_DSN := postgres://postgres:yourpassword@localhost:5432/postgres?sslmode=disable
MIGRATE := migrate -path ./migrations -database $(DB_DSN)

proto:
	protoc -I api --go_out=grpc/pb --go_opt=paths=source_relative --go-grpc_out=grpc/pb --go-grpc_opt=paths=source_relative --grpc-gateway_out=grpc/pb --grpc-gateway_opt=paths=source_relative api/UserService.proto

proto-down:
	rm grpc/pb/*

migrate-new-users:
	migrate create -ext sql -dir ./migrations users

migrate:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) disown

lint:
	golangci-lint run --color=always
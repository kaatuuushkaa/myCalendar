DB_DSN := postgres://postgres:yourpassword@localhost:5433/postgres?sslmode=disable
DB_TEST_DSN := postgres://postgres:yourpassword@localhost:5434/postgres_test?sslmode=disable
MIGRATE := migrate -path ./migrations -database "$(DB_DSN)"
MIGRATE_TEST := migrate -path ./migrations -database "$(DB_TEST_DSN)"

proto-users:
	protoc -I api --go_out=grpc/pb --go_opt=paths=source_relative --go-grpc_out=grpc/pb --go-grpc_opt=paths=source_relative --grpc-gateway_out=grpc/pb --grpc-gateway_opt=paths=source_relative api/UserService.proto

proto-events:
	protoc -I api --go_out=grpc/pb --go_opt=paths=source_relative --go-grpc_out=grpc/pb --go-grpc_opt=paths=source_relative --grpc-gateway_out=grpc/pb --grpc-gateway_opt=paths=source_relative api/EventService.proto

proto:
	protoc -I api --go_out=grpc/pb --go_opt=paths=source_relative --go-grpc_out=grpc/pb --go-grpc_opt=paths=source_relative --grpc-gateway_out=grpc/pb --grpc-gateway_opt=paths=source_relative api/UserService.proto api/EventService.proto

proto-down:
	rm grpc/pb/*

migrate-new-users:
	migrate create -ext sql -dir ./migrations users

migrate-new-events:
	migrate create -ext sql -dir ./migrations events

migrate:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) disown

migrate-test:
	$(MIGRATE_TEST) up

migrate-test-down:
	$(MIGRATE_TEST) disown


lint:
	golangci-lint run --color=always

up:
	docker compose up --build -d

# остановка
down:
	docker compose down

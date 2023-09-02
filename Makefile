.PHONY: start test test-coverage migrate-up migrate-down migrate-force

start:
	go run cmd/main.go

test:
	go test -v ./... 

test-coverage:
	go test -covermode=count -coverprofile coverage.out ./...
	go tool cover -func=coverage

migrate-up:
	migrate -database "mysql://root:root@tcp(localhost:3306)/otp_go" -path migrations up

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-down:
	migrate -database "mysql://root:root@tcp(localhost:3306)/otp_go" -path migrations down

migrate-force:
	migrate -database "mysql://root:root@tcp(localhost:3306)/otp_go" -path migrations force $(version)
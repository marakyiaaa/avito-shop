.PHONY: run down go test test_k6 cover check

run:
	@docker compose up -d

down:
	@docker compose down -v

info_docker:
	docker ps -a
	docker images
	docker volume ls

go:
	@go run cmd/main.go

test:
	@go test -coverprofile=coverage.out ./...

test_k6:
	 @k6 run loadtest.js

cover:
	@go tool cover -func=coverage.out

check:
	@go vet ./...
	@golint ./...
	@errcheck ./...

#golangci-lint run




# psql -U postgres -d shop
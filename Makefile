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
	@go test -coverprofile=coverage.out ./cmd ./internal/...

test_k6:
	 @k6 run loadtest.js

cover:
	@go tool cover -func=coverage.out

docker_test_e2e:
	@docker compose -f docker-compose.test.yml up --build -d

e2e:
	@go test -v ./test/e2e/...

check:
	@golangci-lint run
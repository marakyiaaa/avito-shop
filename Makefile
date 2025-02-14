include .env.local

build:
	@docker compose build

run:
	@docker compose up -d

del:
	@docker compose down -v  # Удаляет контейнеры и тома (вся БД удалится)

inf:
	docker ps -a
	docker images
	docker volume ls


#topic:
#	@docker exec kafka kafka-topics --bootstrap-server kafka:9092 --create --topic orders

go:
	@go run cmd/main.go

test:
	@go test -v internal/service/service_test/auth_service_test.go \
	internal/service/service_test/store_service_test.go
#
#cover:
#	@go test -cover internal/service/service_test.go internal/service/service.go
#
#brew_wrk:
#	@brew install wrk
#
#wrk:
#	@wrk -t4 -c200 -d30s http://localhost:8080/api/orders



check:
	@go vet ./...
	@golint ./...
	@errcheck ./...


#goose create create_orders_table sql

#Откат миграций - goose down


#goose up
#Goose сам будет отслеживать номера версий миграций и применять их в правильной последовательности.



.PHONY: build run reboot topic go test cover brew_wrk wrk check

# psql -U postgres -d shop
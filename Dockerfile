FROM golang:1.22

WORKDIR ${GOPATH}/avito-shop/
COPY . ${GOPATH}/avito-shop/

RUN go build -o /build ./cmd
#    && go clean -cache -modcache

EXPOSE 8080

CMD ["/build"]

## Установка образа для сборки
#FROM golang:1.22.3 AS builder
#
#WORKDIR /app
#COPY . .
#
## Собираем приложение
#RUN go mod tidy
#RUN go build -o avito-shop-service cmd/main.go
#
## Второй этап - контейнер с результатом сборки
#FROM alpine:latest
#
#WORKDIR /app
#
## Установка зависимостей
#RUN apk --no-cache add ca-certificates
#
## Копируем исполнимый файл из стадии сборки
#COPY --from=builder /app/avito-shop-service /avito-shop-service
#
## Даем права на исполнение
#RUN chmod +x /avito-shop-service
#
## Запуск сервиса
#ENTRYPOINT ["/avito-shop-service"]
#
#

#WORKDIR /avito-shop
#COPY . .
#
#RUN go mod tidy
#RUN go build -o ./bin/main cmd/main.go
#RUN chmod +x ./bin/main
#
#EXPOSE 8080
#
##CMD ["/avito-shop/bin/main"]
#CMD ["./bin/main"]
#
##    && go clean -cache -modcache

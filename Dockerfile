FROM golang:1.22.3 AS builder

WORKDIR /go/src/avito-shop
COPY . /go/src/avito-shop

RUN go mod tidy
RUN go build -o ./bin/main cmd/main.go
#    && go clean -cache -modcache

EXPOSE 8080

CMD ["/go/src/avito-shop/bin/main"]
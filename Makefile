.PHONY: lint
lint:
	golangci-lint run -v --color=always --timeout 4m ./...

.PHONY: build
build:
	go build -o bin/main main.go

.PHONY: run
run:
	go run main.go

.PHONY: clean
clean:
	rm -rf bin
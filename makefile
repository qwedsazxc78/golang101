.PHONY: build run test clean d-build d-run d-size

build:
	go build -o ./bin/app ./src/main.go

run:
	go run ./src/main.go

test:
	go test -v ./src

clean:
	go clean

dc-build:
	docker-compose build

dc-run:
	docker-compose up

d-size:
	docker images --format "{{if eq .Repository \"awoo.urlproxy\"}}{{.Size}}{{end}}" awoo.urlproxy:latest

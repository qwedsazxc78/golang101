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

db:
	docker build --tag urlproxy:0.1.0 .

dp-dev:
	docker tag urlproxy:0.1.0 asia-east1-docker.pkg.dev/asia-awoo-com-tw/develop/urlproxy:0.1.0 && docker push asia-east1-docker.pkg.dev/asia-awoo-com-tw/develop/urlproxy:0.1.0

dp-stg:
	docker tag urlproxy:0.1.0 asia-east1-docker.pkg.dev/asia-awoo-com-tw/staging/urlproxy:0.1.0 && docker push asia-east1-docker.pkg.dev/asia-awoo-com-tw/staging/urlproxy:0.1.0

dp-prd:
	docker tag urlproxy:0.1.0 asia-east1-docker.pkg.dev/asia-awoo-com-tw/production/urlproxy:0.1.0 && docker push asia-east1-docker.pkg.dev/asia-awoo-com-tw/production/urlproxy:0.1.0

dp-all: db dp-dev dp-stg dp-prd

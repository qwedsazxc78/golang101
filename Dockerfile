# Build Stage
FROM --platform=linux/amd64 golang:1.20 AS build

WORKDIR /app

COPY src/ .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

# Production Stage
FROM alpine:3.14 AS production

WORKDIR /app

COPY --from=build /app/app .

EXPOSE 8080

CMD ["./app"]

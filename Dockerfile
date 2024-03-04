FROM golang:1.21.5-alpine

WORKDIR /app

COPY src /app

RUN go mod tidy

RUN go mod vendor

CMD ["go", "run", "cmd/auth_app/main.go"]

EXPOSE 8000

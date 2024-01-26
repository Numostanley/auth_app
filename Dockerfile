FROM golang:1.21.5-alpine

WORKDIR /app

COPY src /app

RUN go mod tidy

RUN go mod vendor

CMD ["go", "run", "main.go"]

EXPOSE 8000

FROM golang:1.25-alpine

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o main cmd/server/main.go

EXPOSE 8080

CMD ["./main"]
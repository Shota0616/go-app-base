FROM golang:1.22.4

WORKDIR /usr/src/app

RUN go build -o main cmd/api/main.go

# EXPOSE 8080

CMD ["./main"]

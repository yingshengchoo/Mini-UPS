# dev mode
# FROM golang:1.24.2

# WORKDIR /app

# RUN go install github.com/air-verse/air@latest

# COPY . .

# EXPOSE 8080

# CMD ["air"]


# product mode

FROM golang:1.24.2 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . /app/

RUN go build -o main .

# FROM alpine:latest

# WORKDIR /root/

# COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
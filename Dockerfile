FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
RUN go build -o qabz .

FROM alpine

WORKDIR /app

COPY --from=builder /app/qabz ./

EXPOSE 8080

CMD ["./qabz"]

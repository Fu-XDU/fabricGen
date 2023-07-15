FROM golang:1.18 AS builder
WORKDIR /app
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o fabric-gen-linux-amd64 .
FROM alpine:3.13.6
COPY --from=builder /app/fabric-gen-linux-amd64 /fabric-gen-linux-amd64
RUN chmod +x fabric-gen-linux-amd64
ENTRYPOINT ["/fabric-gen-linux-amd64"]
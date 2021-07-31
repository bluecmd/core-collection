FROM golang:1.16 AS builder
WORKDIR /go/src/github.com/bluecmd/core-collection
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o collect ./cmd/collect
RUN strip collect

FROM alpine:latest
WORKDIR /
COPY --from=builder /go/src/github.com/bluecmd/core-collection/collect ./
CMD ["cp", "/collect", "/corec/collect"]

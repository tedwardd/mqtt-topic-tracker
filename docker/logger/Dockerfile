FROM golang:latest as builder

WORKDIR /app
ENV GO111MODULE on
COPY go.mod go.sum ./
COPY cmd ./cmd
COPY internal ./internal
COPY vendor ./vendor
RUN go mod download
RUN GOOS=linux GOARCH=arm64 go build -o /mqtt-topic-tracker /app/cmd/mqtt-topic-tracker


FROM debian:latest

COPY --from=builder /mqtt-topic-tracker /bin/mqtt-topic-tracker
RUN chmod +x /bin/mqtt-topic-tracker

CMD ["/bin/mqtt-topic-tracker"]


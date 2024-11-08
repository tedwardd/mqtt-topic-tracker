FROM golang:latest as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /mqtt-topic-tracker


FROM alpine:latest

COPY --from=builder /mqtt-topic-tracker /bin/mqtt-topic-tracker
RUN chmod +x /bin/mqtt-topic-tracker

CMD ["/bin/mqtt-topic-tracker"]


FROM golang:latest as builder

WORKDIR /app
COPY go.mod go.sum ./
COPY cmd ./cmd
COPY internal ./internal
COPY vendor ./vendor
RUN go mod download RUN GOOS=linux GOARCH=arm64 go build -o /mqtt-topic-frontend /app/cmd/frontend FROM debian:latest WORKDIR /run COPY --from=builder /mqtt-topic-frontend /bin/mqtt-topic-frontend COPY web ./web
RUN chmod +x /bin/mqtt-topic-frontend -t /run/web/templates


EXPOSE 8080

ENTRYPOINT ["/bin/mqtt-topic-frontend"]


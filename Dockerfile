FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o irc-server ./main.go

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/irc-server .

EXPOSE 8080
ENTRYPOINT [ "./irc-server" ]

FROM golang:1.23-alpine AS builder

RUN apk --no-cache add ca-certificates gcc g++ libc-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ./leaderboard ./cmd/main

FROM alpine:latest

RUN apk --no-cache add ca-certificates 

COPY --from=builder /app/leaderboard /leaderboard

CMD [ "/leaderboard" ]
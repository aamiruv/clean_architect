FROM golang:1.21.5-alpine3.19 as builder
LABEL authors="clean architecture webserver by amir mirzayi"

WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN go build -o web_server ./cmd/web

FROM alpine

WORKDIR /app
COPY --from=builder /app/web_server .
COPY --from=builder /app/config.json .

EXPOSE 8070 8071

ENTRYPOINT ["./web_server"]
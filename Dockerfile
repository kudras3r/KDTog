FROM golang:1.24.2-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/kdtog

FROM alpine:latest

RUN adduser -D appuser

WORKDIR /home/appuser
COPY --from=build /app/app .
COPY static/ ./static/

RUN chown -R appuser:appuser /home/appuser

USER appuser

CMD ["./app"]

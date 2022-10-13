FROM golang:1.19-alpine AS builder

WORKDIR /src
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

FROM alpine:latest

COPY --from=builder src/main .
ENV LISTEN_URL=0.0.0.0:80
EXPOSE 80
CMD ["/main"]

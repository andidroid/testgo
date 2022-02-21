FROM golang:1.17-alpine AS builder

WORKDIR /src
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

FROM scratch

COPY --from=builder src/main .
ENV LISTEN_URL=0.0.0.0:8090
EXPOSE 8090
CMD ["/main"]
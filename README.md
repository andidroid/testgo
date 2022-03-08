# testgo

/d/Programmierung/Git/testgo (master)
$ go build cmd/main.go

/d/Programmierung/Git/testgo (master)
$ ./main.exe
Run Test Go
start server

## Docker build and run

docker build . -t testgo
docker run -p 8090:8090 testgo

## Go Modules

openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout certs/localhost.key -out certs/localhost.crt
go get github.com/gin-contrib/sessions
go get github.com/gin-contrib/sessions/redis
get github.com/gin-contrib/cors

### TODO

session := sessions.Default(c)
session.Set("username", user.Username)
session.Set("token", sessionToken)
session.Save()

### Open telemetry

go.opentelemetry.io/contrib/instrumentation/{IMPORT_PATH}/otel{PACKAGE_NAME}
go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux
go get go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin

go.opentelemetry.io/otel/exporters/jaeger

go get go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo
https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo

https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation/github.com/gin-gonic/gin/otelgin

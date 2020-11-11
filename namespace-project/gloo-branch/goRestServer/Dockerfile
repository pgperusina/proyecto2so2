# Start by building the application.
FROM golang:1.13-buster as build

WORKDIR /go/src/app
ADD . /go/src/app

RUN go get -d -v ./...

# RUN go mod init github.com/my/repo
# RUN go get github.com/gorilla/mux
# RUN go get go.mongodb.org/mongo_driver
# RUN go get github.com/go_redis/redis/v8

RUN go build -o /go/bin/app

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/app /
CMD ["/app"]
FROM golang:latest
MAINTAINER khoale88 "lenguyenkhoa1988@gmail.com"

RUN go get github.com/gorilla/mux
RUN go get gopkg.in/mgo.v2

COPY  . /$GOPATH
RUN go build -o bin/app src/app/server.go

#CMD ./Khoa.RestBucks
# CMD go run src/app/server.go
CMD ./bin/app
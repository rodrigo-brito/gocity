FROM golang:1.10
RUN mkdir -p /go/src/github.com/rodrigo-brito/gocity
WORKDIR /go/src/github.com/rodrigo-brito/gocity
ADD . .
RUN go build
CMD ["/go/src/github.com/rodrigo-brito/gocity/gocity"]

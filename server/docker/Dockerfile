FROM golang:1.12
RUN mkdir -p /gocity
WORKDIR /gocity
ADD . .
RUN go build -mod vendor
CMD ["/gocity/gocity", "server"]

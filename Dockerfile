FROM golang:1.17 as build
WORKDIR /app
ADD . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s"

FROM alpine
COPY --from=build /app/gocity /bin/gocity

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

CMD ["/bin/gocity", "server"]

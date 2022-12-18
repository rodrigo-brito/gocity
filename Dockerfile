FROM alpine
COPY ./gocity /bin/gocity

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

CMD ["/bin/gocity", "server"]

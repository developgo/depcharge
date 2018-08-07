FROM golang:1.10-alpine as go

RUN apk update \
	&& apk add git \
	&& rm -rf /var/cache/apk/*

RUN go get \
    github.com/ghodss/yaml \
    github.com/cbroglie/mustache/... \
    github.com/stretchr/testify

COPY . /go/src/depcharge
WORKDIR /go/src/depcharge

RUN go test .
RUN go build .

FROM alpine:latest
COPY --from=go /go/src/depcharge/depcharge /bin/depcharge
RUN mkdir /mount
WORKDIR /mount
ENTRYPOINT ["depcharge"]
CMD ["--help"]
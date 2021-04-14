FROM golang:1.16 as build-env

WORKDIR /go/src/github.com/fluffle/sp0rkle
ADD . /go/src/github.com/fluffle/sp0rkle

RUN mkdir -p /srv/vol

RUN go get -d -v ./...

RUN go build -o /srv/sp0rkle

FROM gcr.io/distroless/base

COPY --from=build-env /srv /srv
WORKDIR /srv
EXPOSE 6666/tcp
STOPSIGNAL SIGINT
ENTRYPOINT [\
	"/srv/sp0rkle",\
	"--url_cache_dir=/srv/vol/cache",\
    "--boltdb=/srv/vol/sp0rkle.boltdb",\
	"--backup_dir=/srv/vol/backup"]

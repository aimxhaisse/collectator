FROM golang:1.3-cross
MAINTAINER s. rannou <mxs@sbrk.org>

ENV CGO_ENABLED 0
ADD . /app/
RUN cd /app && GOOS=linux GOARCH=arm GOARM=7 go build -o /app/collectator

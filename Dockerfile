# Build Gfro in a stock Go builder container
FROM golang:1.16-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

ADD . /go-frogeum
RUN cd /go-frogeum && make gfro

# Pull Gfro into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go-frogeum/build/bin/gfro /usr/local/bin/

EXPOSE 9506 9507 60606 60606/udp
ENTRYPOINT ["gfro"]

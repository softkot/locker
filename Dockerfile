FROM golang:1.17-alpine as builder
RUN mkdir -p /build
COPY ./ /build/
WORKDIR /build/
RUN env CGO_ENABLED=0 go build -o locker -ldflags "-s -w"

FROM scratch
MAINTAINER Alexei Volkov <softkot@gmail.com>

EXPOSE 8443
COPY --from=builder /build/locker /locker
COPY --from=builder /etc/ssl/cert.pem /etc/ssl/cert.pem
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV ZONEINFO=/zoneinfo.zip
ENTRYPOINT ["/locker"]
CMD []

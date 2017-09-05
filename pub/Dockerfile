FROM golang:1.9-alpine3.6

RUN apk add --no-cache curl git

RUN curl -Lo - https://github.com/cortesi/modd/releases/download/v0.4/modd-0.4-linux64.tgz | tar xz -C ./ modd-0.4-linux64/modd

COPY *.go ./

ENV CGO_ENABLED=0
RUN go-wrapper download
RUN go build -ldflags="-s -w"


FROM alpine:3.6
LABEL maintainer="KenjiTakahashi <kenji.sx>"

COPY --from=0 /go/go /pub
COPY --from=0 /go/modd-0.4-linux64/modd /modd

WORKDIR /data/

VOLUME /data/

ENTRYPOINT ["/modd", "-c", "-f"]
CMD ["/data/modd.conf"]

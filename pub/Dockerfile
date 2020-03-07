FROM golang:1.14-alpine3.11

ARG moddv=0.8

RUN apk add --no-cache curl git

RUN curl -Lo - https://github.com/cortesi/modd/releases/download/v${moddv}/modd-${moddv}-linux64.tgz | tar xz -C ./ modd-${moddv}-linux64/modd

COPY *.go ./

ENV CGO_ENABLED=0
RUN go get -d -v
RUN go build -v
# RUN go install -v .


FROM alpine:3.11
LABEL maintainer="KenjiTakahashi <kenji.sx>"

ARG moddv=0.8

COPY --from=0 /go/go /pub
COPY --from=0 /go/modd-${moddv}-linux64/modd /modd
COPY modd.conf /modd.conf

WORKDIR /data/

VOLUME /data/

ENTRYPOINT ["/modd", "-c", "-f"]
CMD ["/modd.conf"]

FROM golang:1.15-alpine3.12

RUN apk add --no-cache git

ENV CGO_ENABLED=0
WORKDIR /go/src/blog

COPY go.* ./
COPY *.go ./
COPY db/ ./db/

RUN go build


FROM scratch
LABEL maintainer="KenjiTakahashi <kenji.sx>"

COPY --from=0 /go/src/blog/blog /home/blog

EXPOSE 9100

ENTRYPOINT ["/home/blog"]

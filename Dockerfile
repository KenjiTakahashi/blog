FROM golang:1.9-alpine3.6

RUN apk add --no-cache git

COPY *.go ./
COPY db/ ./db/

ENV CGO_ENABLED=0
RUN go-wrapper download
RUN go build -ldflags="-s -w"


FROM scratch
LABEL maintainer="KenjiTakahashi <kenji.sx>"

COPY --from=0 /go/go /home/blog

EXPOSE 9100

ENTRYPOINT ["/home/blog"]

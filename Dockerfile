FROM golang:1.12.7 AS builder

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY Gopkg.lock /go/src/github.com/Vualto/redis-migrate/Gopkg.lock
COPY Gopkg.toml /go/src/github.com/Vualto/redis-migrate/Gopkg.toml

WORKDIR /go/src/github.com/Vualto/redis-migrate/
RUN dep ensure --vendor-only

COPY *.go        /go/src/github.com/Vualto/redis-migrate/
RUN dep ensure

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o redis-migrate .

FROM scratch
LABEL maintainer="roger.pales@vualto.com"

COPY --from=builder /go/src/github.com/Vualto/redis-migrate/redis-migrate .

CMD ["./redis-migrate"]
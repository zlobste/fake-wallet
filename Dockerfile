FROM golang:1.16

WORKDIR /go/src/github.com/zlobste/fake-wallet
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/fake-wallet github.com/zlobste/fake-wallet

ENTRYPOINT ["fake-wallet"]
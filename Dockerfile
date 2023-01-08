FROM golang:1.19-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/jdecool/myip/
COPY . .

RUN go get -d -v . \
    && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /usr/local/bin/myip .

FROM scratch

COPY --from=builder /usr/local/bin/myip /usr/local/bin/myip

ENTRYPOINT [ "/usr/local/bin/myip", "-host", "0.0.0.0" ]

FROM golang:1

WORKDIR /build
ADD . .
RUN go build .
ENTRYPOINT ["/build/gogles"]

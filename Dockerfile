FROM golang:1.8

WORKDIR /go/src/siaworker
COPY . .

RUN go get -v ./...
RUN go build

CMD ["siaworker"]
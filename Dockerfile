FROM golang:1.14

ENV GOPATH /tmp

WORKDIR .
COPY . .
RUN go version
RUN go build .
CMD ["./poll_bot"]
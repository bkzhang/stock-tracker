FROM golang:1.11.1-alpine

WORKDIR $GOPATH/src/github.com/bkzhang/server
COPY . $GOPATH/src/github.com/bkzhang/server

RUN apk update && apk add git && apk add gcc 
RUN go get -u golang.org/x/vgo
#RUN export GO111MODULE=on
RUN CGO_ENABLED=0 vgo install . 

FROM alpine:latest
COPY --from=0 go/bin/server .

ENV PORT=8080
EXPOSE 8080
EXPOSE 443
ENTRYPOINT ["./server"]

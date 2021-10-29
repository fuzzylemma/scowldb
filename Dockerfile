FROM golang:1.17-alpine

RUN apk update && apk upgrade && \
    apk add gcc musl-dev git 

RUN mkdir -p /opt/scowldb && \
    git clone --depth 1 https://github.com/fuzzylemma/scowldb /opt/scowldb && \
    cd /opt/scowldb && \
    go mod tidy && \
    go build

EXPOSE 8888
WORKDIR /opt/scowldb
CMD ["./scowldb"]


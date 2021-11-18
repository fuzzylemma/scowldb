FROM golang:1.17-alpine

RUN apk update && apk upgrade && \
    apk add gcc musl-dev git make 

RUN mkdir -p /opt/scowldb && \
    git clone --depth 1 https://github.com/fuzzylemma/scowldb /opt/scowldb && \
    cd /opt/scowldb && \
    go mod tidy && \
    go build && \
    make getscowl && make unzip

EXPOSE 8888
WORKDIR /opt/scowldb
CMD ["./scowldb"]


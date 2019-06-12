FROM golang:latest
RUN mkdir -p /go/src/eos-cassandra-historyapi
ADD dep /usr/bin/dep
RUN chmod +x /usr/bin/dep
ADD . /go/src/eos-cassandra-historyapi/
WORKDIR /go/src/eos-cassandra-historyapi
RUN dep ensure && go build -o middleware .
EXPOSE 3000:3000
CMD ["/go/src/eos-cassandra-historyapi/middleware"]


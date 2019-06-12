FROM golang:latest 
RUN mkdir -p /go/src/app
ADD dep /usr/bin/dep
RUN chmod +x /usr/bin/dep 
ADD . /go/src/app/ 
WORKDIR /go/src/app 
RUN dep ensure && go build -o middleware . 
EXPOSE 3000:3000
CMD ["/go/src/app/middleware"]

FROM golang:1.13-alpine
WORKDIR /dcard-simple-demo
ADD . /dcard-simple-demo
RUN cd /dcard-simple-demo && go build
ENTRYPOINT ["./dcard-simple-demo"]
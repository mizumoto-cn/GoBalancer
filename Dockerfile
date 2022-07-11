FROM golang:1.18.3-alpine
LABEL github.com.mizumoto-cn.authors="mizumoto-cn<mizumotokunn@gmail.com>"

ENV GO111MODULE=on \
    CGO_ENABLE=0

# copy the source code
ADD . /gobalancer

# default entrypoint
WORKDIR /gobalancer

# install dependencies
RUN go build -v

# run the application
CMD ["/gobalancer/gobalancer"]

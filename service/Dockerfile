FROM golang:1.13-stretch

# install gin
RUN go get github.com/codegangsta/gin
# create a working directory
WORKDIR /go/src
# add Gopkg.toml and Gopkg.lock
ADD app app
CMD ["gin", "--path", "app", "--port", "4000", "run", "*.go"]

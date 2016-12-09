FROM golang:1.6
ADD . /go/src/qiniu.com/ufop/example/demo
RUN go install qiniu.com/ufop/example/demo
ENTRYPOINT /go/bin/demo

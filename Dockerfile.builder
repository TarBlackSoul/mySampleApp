FROM golang:1.15.8
RUN go env
ADD . /go/src/app/
RUN cat /go/src/app/go.mod
ENV GO111MODULE=on
RUN cd /go/src/app/ && go build -o myapp /go/src/app/main.go
#FROM alpine:3.13
#COPY --from=builder /go/src/app .
WORKDIR /go/src/app
EXPOSE 9000
CMD ["/go/src/app/myapp"]

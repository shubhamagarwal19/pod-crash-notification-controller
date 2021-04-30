FROM golang:1.16 AS build
WORKDIR /go/src
COPY . /go/src
RUN CGO_ENABLED=0 GOOS=linux go build -a -mod=readonly ./

FROM alpine:latest
WORKDIR /root/
COPY --from=build /go/src/pod-crash-notification-controller .
ENTRYPOINT ["./pod-crash-notification-controller"]

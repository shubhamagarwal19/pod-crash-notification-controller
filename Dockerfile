FROM openshift/origin-release:golang-1.15 AS build
WORKDIR /go/src
COPY . /go/src
RUN CGO_ENABLED=0 GOOS=linux go build -a -mod=readonly ./

FROM registry.access.redhat.com/ubi8/ubi-minimal
WORKDIR /root/
COPY --from=build /go/src/pod-crash-notification-controller .
ENTRYPOINT ["./pod-crash-notification-controller"]

# Phase 1: Compile httpServ
ARG GO_IMAGE=golang:1.19.2-alpine3.16
ARG RUNTIME_IMAGE=alpine:3.16
FROM $GO_IMAGE AS builder

RUN go env -w GO111MODULE=on GOPROXY=https://proxy.golang.org,direct

COPY . /go/src/httpServ/
WORKDIR /go/src/httpServ/cmd
# By default it matches to alpine-3.16 as the container name
ARG LD_FLAGS
RUN go build -ldflags "$LD_FLAGS"

# Phase 2: Setup and run runtime container from alpine:3.16 (Light weight)
FROM $RUNTIME_IMAGE

COPY --from=builder /go/src/httpServ/cmd/cmd /bin/httpServ

EXPOSE 8081

ENTRYPOINT [ "/bin/httpServ" ]

# How to run this container image:
# docker run -p 8081:8081 httpserv

# Once container is up, we can use nsenter to get into container network ns to check IP configs
# ps -ef | grep httpserv | grep -v grep
# ls -la /proc/12345/ns
# nsenter -t 12345 -n ip addr
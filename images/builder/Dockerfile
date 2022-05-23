# Build the manager binary
# FROM registry.access.redhat.com/ubi8/go-toolset:1.17.7 as builder
FROM golang:1.18 as builder

USER root
WORKDIR /workspace
COPY . .
RUN make build

FROM registry.access.redhat.com/ubi8/ubi-minimal

LABEL MAINTAINER "Adrian Riobo" "<ariobolo@redhat.com>"

COPY --from=builder /workspace/out/qe-eventmanager /workspace/images/builder/entrypoint.sh /usr/local/bin/

ENTRYPOINT entrypoint.sh 

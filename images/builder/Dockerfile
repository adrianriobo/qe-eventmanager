# Build the manager binary
FROM registry.access.redhat.com/ubi8/go-toolset:1.14.12 as builder

WORKDIR /workspace
COPY . .
RUN make build

FROM registry.access.redhat.com/ubi8/ubi-minimal

LABEL MAINTAINER "Adrian Riobo" "<ariobolo@redhat.com>"

WORKDIR /
COPY --from=builder /workspace/out/qe-eventmanager .
#USER 65532:65532

ENTRYPOINT ["/qe-eventmanager"]
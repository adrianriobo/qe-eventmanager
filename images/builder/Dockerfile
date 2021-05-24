# Build the manager binary
FROM registry.access.redhat.com/ubi8/go-toolset:1.14.12 as builder

WORKDIR /tmp
COPY . .
RUN make build

FROM registry.access.redhat.com/ubi8/ubi-minimal

LABEL MAINTAINER "Adrian Riobo" "<ariobolo@redhat.com>"

COPY --from=builder /tmp/out/qe-eventmanager /tmp/images/builder/entrypoint.sh /usr/local/bin/

ENTRYPOINT entrypoint.sh 
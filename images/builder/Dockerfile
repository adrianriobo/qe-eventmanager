# Build the manager binary
FROM registry.access.redhat.com/ubi8/go-toolset:1.14.12 as builder

WORKDIR /workspace
COPY . .
RUN make build

FROM registry.access.redhat.com/ubi8/ubi-minimal

LABEL MAINTAINER "Adrian Riobo" "<ariobolo@redhat.com>"

COPY --from=builder /workspace/out/qe-eventmanager /workspace/images/builder/entrypoint.sh /usr/local/bin/

ENTRYPOINT entrypoint.sh 



#  "start", 
#                 "-b", "messaging-devops-broker02.web.stage.ext.phx2.redhat.com:61612", 
#                 "--ca-certs", "/etc/pki/ca-trust/source/anchors/2015-RH-IT-Root-CA.pem", 
#                 "--certificate-file", "/home/ariobolo/02OFFTOPIC/umb-certificates/psi-crcqe-openstack.crt",
#                 "--private-key-file", "/home/ariobolo/02OFFTOPIC/umb-certificates/psi-crcqe-openstack.key",
#                 "-k", "/home/ariobolo/02OFFTOPIC/umb-certificates/config"]
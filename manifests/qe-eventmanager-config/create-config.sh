#!/bin/bash

# Usage create-config.sh $CA_FILE_PATH $CERTIFICATE_FILE_PATH $KEY_FILE_PATH $BROKERS $CONSUMER_ID $DRIVER

ca=$(cat $1 | base64 -w0)
certificate=$(cat $2 | base64 -w0)
key=$(cat $3 | base64 -w0)
brokers=$(echo -n $4 | base64 -w0)
consumer_id=$(echo -n $5 | base64 -w0)
driver=$(echo -n $6 | base64 -w0)

if [ "${DEBUG:-}" = "true" ]; then
  set -xuo 
fi

set -e pipefail

# Create file
cat <<EOF > config.yaml
---
apiVersion: v1
kind: Secret
metadata:
  name: qe-eventmanager-config
type: Opaque
data:
  ca: ${ca}
  certificate: ${certificate}
  key: ${key}
  brokers: ${brokers}
  consumerid: ${consumer_id}
  driver: ${driver}
EOF
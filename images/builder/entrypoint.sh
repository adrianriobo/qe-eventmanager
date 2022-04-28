#!/bin/bash

VALID_CONFIG=true
# Check required ENVs
if [ -z "${BROKERS}" ]; then 
  echo "BROKERS ENV is required"
  VALID_CONFIG=false
fi

if [ -z "${CA}" ]; then 
  echo "CA ENV is required"
  VALID_CONFIG=false  
fi

if [ -z "${CERTIFICATE}" ]; then 
  echo "CERTIFICATE ENV is required"
  VALID_CONFIG=false
fi

if [ -z "${KEY}" ]; then 
  echo "KEY ENV is required"
  VALID_CONFIG=false
fi

if [ -z "${CONSUMER_ID}" ]; then 
  echo "CONSUMER_ID ENV is required"
  VALID_CONFIG=false
fi

if [ -z "${DRIVER}" ]; then 
  DRIVER=stomp
fi


if [ "${VALID_CONFIG}" = false ]; then
  echo "Add the required ENVs"
  exit 1
fi

# Run qe-eventmanager
if [ -z "${KUBECONFIG}" ]; then
  exec qe-eventmanager start \
    --brokers "${BROKERS}" \
    --ca-certs "${CA}" \
 		--certificate-file "${CERTIFICATE}" \
		--private-key-file "${KEY}" \
    --consumerid "${CONSUMER_ID}" \
    --driver "${DRIVER}" 
else 
  exec qe-eventmanager start \
    --brokers "${BROKERS}" \
    --ca-certs "${CA}" \
    --certificate-file "${CERTIFICATE}" \
    --private-key-file "${KEY}" \
    --kubeconfig "${KUBECONFIG}" \
    --consumerid "${CONSUMER_ID}" \
    --driver "${DRIVER}" 
fi

#!/bin/bash

VALID_CONFIG=true
# Check required ENVs
if [ -z "${PROVIDERS_FILE_PATH}" ]; then 
  echo "PROVIDERS_FILE_PATH is required"
  VALID_CONFIG=false
fi

if [ -z "${FLOWS_FILE_PATH}" ]; then 
  echo "FLOWS_FILE_PATH ENV is required"
  VALID_CONFIG=false  
fi

if [ "${VALID_CONFIG}" = false ]; then
  echo "Add the required ENVs"
  exit 1
fi

# Run eventmanager
exec eventmanager start \
    --providers-filepath "${PROVIDERS_FILE_PATH}" \
    --flows-filepath "${FLOWS_FILE_PATH}"

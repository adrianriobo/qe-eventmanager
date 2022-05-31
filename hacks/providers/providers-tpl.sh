#!/bin/bash

configname=${1}
ca=$(cat $2 | base64 -w0)
certificate=$(cat $3 | base64 -w0)
key=$(cat $4 | base64 -w0)
brokers=${5}
consumer_id=${6}
driver=${7}
namespace=${8}
workspace_1=${9}
pvc_1="${10}"
consoleurl="${11}"
# kubeconfig=$(cat ${12} | base64 -w0)
# github_token="${13}"
github_app_id="${12}"
github_app_installation_id="${13}"
github_app_key=$(cat ${14} | base64 -w0)


if [ "${DEBUG:-}" = "true" ]; then
  set -xuo 
fi

set -e pipefail

# Create file
cat <<-EOF > providers.yaml
umb:
  consumerID: ${consumer_id}
  driver: ${driver}
  brokers: ${brokers}
  userCertificate: ${certificate}
  userKey: ${key}
  certificateAuthority: ${ca}
tekton:
  namespace: ${namespace}
  workspaces:
  - name: ${workspace_1}
    pvc: ${pvc_1}
  kubeconfig: ${kubeconfig}
  consoleURL: ${consoleurl}
github:
  appID: ${github_app_id} 
  appInstallationID: ${github_app_installation_id}
  appKey: ${github_app_key}
EOF

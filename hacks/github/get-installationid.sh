#!/bin/bash
#$1 app id
#$2 file path for app key  
jwt_token=$(./githubapp-jwt.sh ${1} ${2})
curl -H "Authorization: Bearer ${jwt_token}" \
     -H "Accept: application/vnd.github.v3+json" \
     https://api.github.com/app/installations
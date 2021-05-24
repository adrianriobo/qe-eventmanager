#!/bin/bash

oc create secret generic qe-eventmanager-config \
    --from-file=ca=ca.crt \
    --from-file=certificate=user.crt \
    --from-file=key=user.key
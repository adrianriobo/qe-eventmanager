# qe-eventmanager
Sample app for handling qe events

[![Container Repository on Quay](https://quay.io/repository/ariobolo/qe-eventmanager/status "Container Repository on Quay")](https://quay.io/repository/ariobolo/qe-eventmanager)

# overview

UMB integration with qe platform

![Overview](docs/diagrams/overview.jpg?raw=true)

# Deploy

```bash 
# Create config
manifest/create-config.sh $CA_FILE_PATH $CERTIFICATE_FILE_PATH $KEY_FILE_PATH $BROKERS

# Deploy resources
oc apply -f manifest/config.yaml
oc apply -f manifest/rbac.yaml
oc apply -f manifest/deployment.yaml
```

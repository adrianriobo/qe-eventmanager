# Overview  

This task allows to send messages to umb

## Configuration

The task will require the UMB provider credentials with rights to read / write on the topic

providers.yaml

```yaml
umb:
  consumerID: ${consumer_id}
  driver: ${driver}
  brokers: ${brokers}
  userCertificate: ${certificate}
  userKey: ${key}
  certificateAuthority: ${ca}
```

which then would be included as a secret:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: umb-provider
type: Opaque
data:
  providers.yaml: 
  ${content}
```

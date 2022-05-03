# qe-eventmanager

Sample app for handling qe events

[![Container Repository on Quay](https://quay.io/repository/ariobolo/qe-eventmanager/status "Container Repository on Quay")](https://quay.io/repository/ariobolo/qe-eventmanager)

## overview

UMB integration with qe platform

![Overview](docs/diagrams/overview.jpg?raw=true)

## Configuration

The eventmanager requires a set of information around the `providers` on which it can act upon  
and a set or `rules` defining the integrations and the actions to be executed.  

### Providers

Here is a sample for a providers configuration:  

```yaml
umb:
  consumerID: foo
  driver: amqp
  brokers: broker1:5556,broker2:5556
  userCertificate-data: XXXXXXX
  userKey-data: XXXXX
  certificateAuthority-data: XXXXXX
tekton:
  namespace: myNamespace
  workspaces:
  - name: workspace1
    pvc: pvc1
  - name: workspace2
    pvc: pvc2
  kubeconfig-data: XXXXXX # This value is optional is used to connect to remote cluster
                          # Otherwise eventmanager can rely on RBAC when running inside the cluster
```

### Rules  

Here is a sample for a rule defintion:  

```yaml
name:  sample-umb-rule
input:
  umb:
    topic: topic-to-consume
    filters:
      - $.estrcuture.list[?(@.field1=='value1')].field1
      - $.estrcuture.list[?(@.field2=='value2')].field1
action:
  tektonPipeline:
    name: XXX
    params:
    - name: foo
      value: $.estrcuture.list[(@.field=='foo')].id # $. jsonpath expression function
    - name: bar
      value: constant #fixed value 
    results:
    - name: baz # accessible with actionResults function
      pipelineResultName: baz 
    - name: error # accessible with actionResults function
      pipelineResultName: error 
success:
  umb:
    topic: topic-to-produce
    eventSchema: message-schema-to-send
    eventFields:
     - name: baz
       value: actionResult[baz] #actionResults allow to access results 
     - name: fixed
       value: constant #fixed value
error:
  umb:
    topic: topic-to-produce
    eventSchema: message-schema-to-send
    eventFields:
    - name: error
      value: actionResult[error] #actionResults allow to access results 
```

## Build

```bash
podman build -t quay.io/ariobolo/qe-eventmanager:$VERSION -f images/builder/Dockerfile .
```

## Deploy

```bash
# Create config
manifest/create-config.sh $CA_FILE_PATH $CERTIFICATE_FILE_PATH $KEY_FILE_PATH $BROKERS

# Deploy resources
oc apply -f manifest/config.yaml
oc apply -f manifest/rbac.yaml
oc apply -f manifest/deployment.yaml
```

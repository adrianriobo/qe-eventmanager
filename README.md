# qe-eventmanager

![avatar](docs/diagrams/eventmanager.svg)

Sample app for handling qe events

[![Container Repository on Quay](https://quay.io/repository/ariobolo/qe-eventmanager/status "Container Repository on Quay")](https://quay.io/repository/ariobolo/qe-eventmanager)

## overview

UMB integration with qe platform

## Roles

The manager can act as a manager handling integration based on flow definitions, or it can act as a tool to interact within the providers configured.

### Actioner

As an actioner the cli allows to run single actions:

* Send an UMB menssage

```bash
./qe-eventmanager umb send -p providers.yaml \
                           -m message.json \
                           -d VirtualTopic.sample
```

### Manager

As a manager integrate providers based on flow definitnions:

```bash
./qe-eventmanager start -p providers.yaml \
                        -f flow1.yaml,flow2.yaml
```

A simple overview on an umb-tekton integration

![Overview](docs/diagrams/overview.jpg?raw=true)

## Configuration

The eventmanager requires a set of information around the `providers` on which it can act upon  
and a set or `flows` defining the integrations and the actions to be executed.  

### Providers

```yaml
umb:
  consumerID: foo
  driver: amqp
  brokers: broker1:5556,broker2:5556
  userCertificate: XXXXXXX # encoded as base64
  userKey: XXXXX # encoded as base64
  certificateAuthority: XXXXXX # encoded as base64
tekton:
  namespace: myNamespace
  workspaces:
  - name: workspace1
    pvc: pvc1
  - name: workspace2
    pvc: pvc2
  kubeconfig: XXXXXX # encoded as base64. This value is optional is used to connect to remote cluster
                     # Otherwise eventmanager can rely on RBAC when running inside the cluster
github:
  token: github_pat_token
```

### Flows  

```yaml
name:  sample-flow
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
      value: bar # constant string 
  success:
    umb:
      topic: topic-to-produce
      eventSchema: message-schema-to-send
      eventFields:
      - name: foo
        value: $(pipeline.results.result1) # Pick value from pipeline results result1 
      - name: bar
        value: bar # constant string
  error:
    umb:
      topic: topic-to-produce
      eventSchema: message-schema-to-send
      eventFields:
      - name: baz
        value: baz
```

## Build

### Cli

```bash
make clean
make build
```

### Container

```bash
make container-build
```
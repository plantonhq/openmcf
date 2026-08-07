# KubernetesKafkaConnector Pulumi Module

Deploys one `kafka.strimzi.io/v1` `KafkaConnector` custom resource —
one declared data pipe on an existing KubernetesKafkaConnect cluster.
The typed spec renders into an untyped CR body in `module/main.go` —
the exact twin of the Terraform module's `local.connector_manifest`
(same keys rendered and omitted, numbers as ints).

## What the Module Creates

1. **The KafkaConnector CR** (named `metadata.name`) — in the Connect
   cluster's OWN namespace, carrying the `strimzi.io/cluster` label
   rendered from `connect_cluster` on top of the identity labels
   (without it the cluster operator never picks the connector up)

That is the whole footprint, deliberately: no namespace resource (the
namespace belongs to the KubernetesKafkaConnect resource's
lifecycle — a connector that creates its own namespace could violate
the placement contract), no satellites.

## Why an Untyped CustomResource

The Strimzi CRDs type `config` with
`x-kubernetes-preserve-unknown-fields`, which crd2pulumi cannot
carry — no generated package is shipped for the Kafka family (the
same ruling as the KubernetesKafka module). Shape errors still fail
loudly at the operator's schema validation.

## Rendering Notes

- **`class` is the one always-rendered spec key**; tasksMax renders
  as an int only when set, version/state only when non-empty, config
  only when non-empty.
- **autoRestart renders enabled plus maxRestarts** (int, only when
  set).
- **The offset ConfigMap targets are declarations only** —
  listOffsets/alterOffsets render `{toConfigMap|fromConfigMap:
  {name}}`; the list/alter ACTIONS run when the resource carries the
  `strimzi.io/connector-offsets` annotation, an operational verb
  outside this module's scope.
- **No await machinery** — reconciliation belongs to the cluster
  operator (which drives the connector through the Connect REST API
  internally); a "class not found" condition surfaces on the resource
  status, never as an apply failure.

## Usage

```shell
planton pulumi up --manifest hack/manifest.yaml --module-dir <path-to-this-module>
```

## Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the KafkaConnector lives in (the Connect cluster's namespace) |
| `connector_name` | Connector name (`metadata.name`) — what the Connect REST API and `connect-<name>` consumer groups key off |

## Module Structure

- `main.go`: entrypoint that calls the module
- `module/main.go`: labels (identity + the Strimzi binding label),
  the CR body, output exports
- `module/outputs.go`: output name constants

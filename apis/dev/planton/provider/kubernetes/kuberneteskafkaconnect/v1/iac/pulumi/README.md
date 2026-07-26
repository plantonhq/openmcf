# KubernetesKafkaConnect Pulumi Module

Deploys one Strimzi-operator-managed Kafka Connect cluster: the
optional namespace, the optional JMX metrics rules ConfigMap, and the
`kafka.strimzi.io/v1` `KafkaConnect` CR itself. The typed spec
renders into an untyped CR body in `module/connect.go` — the exact
twin of the Terraform module's `locals.tf` (same keys rendered and
omitted, numbers as ints, booleans as booleans).

## What the Module Creates

1. **Namespace** (optional) — created with the standard governance
   labels when `create_namespace` is true; otherwise the namespace
   must already exist
2. **JMX metrics ConfigMap** (optional, `metrics.enabled`) — the
   canonical Strimzi connect Prometheus rules as
   `<name>-connect-metrics` (key `metrics-config.yml`); the CR's
   `metricsConfig` only points at it
3. **The KafkaConnect CR** (named `metadata.name`) — always annotated
   `strimzi.io/use-connector-resources: "true"`, so connectors on
   this cluster are managed declaratively through
   KubernetesKafkaConnector resources (the operator reverts
   REST-API-made changes)

## Why Untyped CustomResources

The Strimzi CRDs type `spec.config` (and the custom authentication's
config block) with `x-kubernetes-preserve-unknown-fields`, which
crd2pulumi flattens into shapes that cannot hold the free-typed
bodies — so no generated package is shipped for the Kafka family at
all. Shape errors still fail loudly: the operator validates the
applied spec against its schema.

## Rendering Notes

- **The group identity quartet always renders** — group.id and the
  three storage topics, resolved in `locals.go` from spec overrides
  with metadata.name-derived fallbacks (`<name>`,
  `<name>-connect-{configs,status,offsets}`); explicit rendering
  keeps the uniqueness contract visible in the applied object.
- **`image` renders only on the prebuilt arm** — image and build are
  mutually exclusive at validation, so the build arm never carries a
  competing image.
- **Per-arm authentication** — each type renders only its own
  credential fields (tls: certificateAndKey with
  `user.crt`/`user.key` fallbacks; SCRAM/plain: username +
  passwordSecret with the `password` key fallback; custom: sasl +
  config) — the KubernetesKafkaUser credential-Secret layout.
- **Per-type build artifacts** — a maven artifact renders
  group/artifact/version (repository only when set); url-family
  artifacts render url/sha512sum/insecure (+ fileName for `other`) —
  only the keys the operator's per-type sub-schema allows.
- **OCI plugins own their artifact type** — the literal `image` is
  the only artifact type the plugins arm supports, so the module
  renders it rather than asking the author to repeat it.
- **`node_selector` translates to node affinity** — the Strimzi pod
  template carries affinity and tolerations but NO nodeSelector; one
  matchExpressions entry per label, sorted for determinism.
- **An unset optional is never inserted into the map** (the Go twin
  of Terraform's null-prune), so the apiserver applies the CRD's own
  defaults.

## No Await Machinery, Deliberately

Worker readiness depends on the operator (image pulls or an
operator-driven image BUILD, Connect group formation) that is not
part of applying the resources — the never-block-on-a-controller
posture of every operator-CR kind in the catalog.

## Usage

```shell
planton pulumi up --manifest hack/manifest.yaml --module-dir <path-to-this-module>
```

## Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the Connect cluster deploys into |
| `connect_name` | Connect cluster name (`metadata.name`) — the `strimzi.io/cluster` binding value for KubernetesKafkaConnector resources |
| `rest_api_service_name` | Connect REST API Service (`<name>-connect-api`) |
| `rest_api_endpoint` | In-cluster REST endpoint (`http://<name>-connect-api.<namespace>.svc.cluster.local:8083`) — read-only inspection |

## Module Structure

- `main.go`: entrypoint that calls the module
- `module/main.go`: namespace → metrics ConfigMap → KafkaConnect CR →
  output exports
- `module/connect.go`: the CR body (Kafka connection, group identity,
  worker config, the image/plugins/build arms, sizing, rack, metrics
  wiring, the pod template translation)
- `module/metrics_rules.go`: the module-owned JMX Prometheus rules
- `module/locals.go`: the group-identity resolution and naming
  contracts (`<name>-connect-api`, `<name>-connect-metrics`) — kept
  in lockstep with the Terraform module's `locals.tf`
- `module/vars.go`: the `kafka.strimzi.io/v1` apiVersion, the
  use-connector-resources annotation key, the metrics ConfigMap key

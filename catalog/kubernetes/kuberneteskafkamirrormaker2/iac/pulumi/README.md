# KubernetesKafkaMirrorMaker2 Pulumi Module

Deploys one Strimzi-operator-managed MirrorMaker 2 replication
engine: the optional namespace, the optional JMX metrics rules
ConfigMap, and the `kafka.strimzi.io/v1` `KafkaMirrorMaker2` CR
itself. The typed spec renders into an untyped CR body in
`module/mirrormaker2.go` — the exact twin of the Terraform module's
`locals.tf` (same keys rendered and omitted, numbers as ints,
booleans as booleans).

## What the Module Creates

1. **Namespace** (optional) — created with the standard governance
   labels when `create_namespace` is true
2. **JMX metrics ConfigMap** (optional, `metrics.enabled`) — the
   canonical Strimzi connect Prometheus rules as
   `<name>-mm2-metrics` (key `metrics-config.yml`); the CR's
   `metricsConfig` only points at it
3. **The KafkaMirrorMaker2 CR** (named `metadata.name`) — one target
   body, one mirror body per `mirrors` entry, worker sizing, rack,
   pod template

## Why an Untyped CustomResource

The Strimzi CRDs type the cluster and connector `config` blocks with
`x-kubernetes-preserve-unknown-fields`, which crd2pulumi cannot
carry — no generated package is shipped for the Kafka family (the
same ruling as the KubernetesKafka module).

## Rendering Notes

- **The target body always carries the group identity** — alias
  (default `target`), group.id and the three storage topics, resolved
  in `locals.go` from spec overrides with metadata.name-derived
  fallbacks (`<name>-mirrormaker2-{configs,status,offsets}`).
- **Alias semantics are validation-backed** — source aliases unique,
  target alias distinct from every source alias (CEL at the spec);
  the module renders aliases verbatim. Under the default replication
  policy mirrored topics arrive prefixed with the source alias;
  IdentityReplicationPolicy on a mirror's source AND checkpoint
  connectors keeps original names.
- **Shared client shapes** — target and source TLS/authentication use
  the same rendering as KubernetesKafkaConnect: each authentication
  type renders only its own credential fields, with the
  KubernetesKafkaUser credential-Secret fallbacks
  (`user.crt`/`user.key`, password key `password`).
- **Connector bodies render only what is set** — tasksMax as an int,
  config, autoRestart (enabled + maxRestarts).
- **`node_selector` translates to node affinity** — the Strimzi pod
  template carries no nodeSelector; one matchExpressions entry per
  label, sorted for determinism.
- **An unset optional is never inserted into the map** (the Go twin
  of Terraform's null-prune).

## No Await Machinery, Deliberately

Engine readiness depends on the operator (image pulls, worker group
formation, connector startup) that is not part of applying the
resources — the never-block-on-a-controller posture of every
operator-CR kind.

## Usage

```shell
planton pulumi up --manifest e2e/manifest.yaml --module-dir <path-to-this-module>
```

## Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the deployment runs in |
| `mirrormaker_name` | The deployment's name (`metadata.name`) |
| `rest_api_endpoint` | In-cluster Connect REST endpoint of the engine (`http://<name>-mirrormaker2-api.<namespace>.svc.cluster.local:8083`) — read-only inspection |

## Module Structure

- `main.go`: entrypoint that calls the module
- `module/main.go`: namespace → metrics ConfigMap → KafkaMirrorMaker2
  CR → output exports
- `module/mirrormaker2.go`: the CR body (target, mirrors, shared
  client TLS/authentication rendering, connector bodies, sizing,
  rack, metrics wiring, the pod template translation)
- `module/metrics_rules.go`: the module-owned JMX Prometheus rules
- `module/locals.go`: group-identity resolution and naming contracts
  (`<name>-mirrormaker2-api`, `<name>-mm2-metrics`) — kept in
  lockstep with the Terraform module's `locals.tf`
- `module/vars.go`: the `kafka.strimzi.io/v1` apiVersion, the metrics
  ConfigMap key

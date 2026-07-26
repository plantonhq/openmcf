# KubernetesKafka Pulumi Module

Deploys one Strimzi-operator-managed KRaft Kafka cluster: the optional
namespace, the optional JMX metrics rules ConfigMap, one
`kafka.strimzi.io/v1` `KafkaNodePool` per `node_pools` entry, and the
`kafka.strimzi.io/v1` `Kafka` CR itself. The typed spec renders into
untyped CR bodies in `module/kafka.go` / `module/nodepools.go` — the
exact twin of the Terraform module's `locals.tf` (same keys rendered
and omitted, numbers as ints, booleans as booleans).

## What the Module Creates

1. **Namespace** (optional) — created with the standard governance
   labels when `create_namespace` is true; otherwise the namespace
   must already exist
2. **JMX metrics ConfigMap** (optional, `metrics.enabled`) — the
   canonical Strimzi Prometheus rules as `<name>-kafka-metrics`
   (key `kafka-metrics-config.yml`); the Kafka CR's `metricsConfig`
   only points at it
3. **KafkaNodePool CRs** — one per pool, named after the pool, bound
   to the cluster by the `strimzi.io/cluster` label riding on top of
   the identity labels; created BEFORE the Kafka CR (Strimzi tolerates
   either order, but a Kafka CR with no matching pools reports a
   transient warning state)
4. **The Kafka CR** (named `metadata.name`) — listeners, config,
   authorization, entity operators, Cruise Control, exporter, CAs,
   rack, JVM, maintenance windows

## Why Untyped CustomResources

The typed crd2pulumi tree is structurally unable to carry this CR: the
CRD types `spec.kafka.config` (and the listener and Cruise Control
configuration blocks) with `x-kubernetes-preserve-unknown-fields`,
which the generated types flatten into shapes that cannot hold the
free-typed bodies — so no generated package is shipped for the Kafka
family at all. Shape errors still fail loudly: the operator validates
the applied spec against its schema.

## Rendering Notes

- **Every listener carries the CRD-required quartet**
  (name/port/type/tls, type coalesced to `internal`); authentication
  and configuration render only when declared. Custom-authentication
  knobs (`sasl`, `listenerConfig`) render only on the custom arm, and
  custom-authorizer knobs (`authorizerClass`, `supportsAdminApi`) only
  on the custom authorization arm.
- **Storage renders per arm** — ephemeral (type only),
  persistent-claim (size/class/deleteClaim), jbod (per-volume bodies,
  at most one carrying `kraftMetadata: shared`).
- **`node_selector` translates to node affinity** — the Strimzi pod
  template carries affinity and tolerations but NO nodeSelector; the
  module renders a required node affinity with one matchExpressions
  entry per label, sorted for determinism (the Terraform module
  renders the same translation).
- **The entity operator omits cleanly** — each sub-operator renders
  when enabled (spec defaults both true); with BOTH disabled the block
  is omitted entirely and KafkaTopic/KafkaUser declarations for this
  cluster become inert.
- **cert-manager key fallbacks** — `broker_cert_chain_and_key` renders
  `certificate`/`key` defaulted to `tls.crt`/`tls.key`, the names
  cert-manager writes.
- **An unset optional is never inserted into the map** (the Go twin of
  Terraform's null-prune), so the apiserver applies the CRD's own
  defaults.

## No Await Machinery, Deliberately

Cluster readiness depends on the operator (image pulls, KRaft quorum
formation, listener provisioning) that is not part of applying the
resources — the never-block-on-a-controller posture of every
operator-CR kind in the catalog.

## Usage

```shell
planton pulumi up --manifest hack/manifest.yaml --module-dir <path-to-this-module>
```

## Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the cluster deploys into |
| `cluster_name` | Kafka cluster name (`metadata.name`) — the `strimzi.io/cluster` binding value |
| `bootstrap_service_name` | Internal bootstrap Service (`<cluster>-kafka-bootstrap`) |
| `internal_bootstrap_endpoint` | In-cluster bootstrap address for the first internal listener (empty when none is declared) |
| `cluster_ca_cert_secret_name` | Cluster CA Secret (`<cluster>-cluster-ca-cert`, key `ca.crt`) |

## Module Structure

- `main.go`: entrypoint that calls the module
- `module/main.go`: namespace → metrics ConfigMap → node pools →
  Kafka CR → output exports
- `module/kafka.go`: the Kafka CR body (listeners, config,
  authorization, entity operators, Cruise Control, exporter, CAs,
  rack, JVM, maintenance windows)
- `module/nodepools.go`: the KafkaNodePool bodies (roles, storage
  arms, resources, the node-affinity translation)
- `module/metrics.go` / `module/metrics_rules.go`: the module-owned
  JMX Prometheus rules ConfigMap
- `module/locals.go`: naming contracts (`<name>-kafka-bootstrap`,
  `<name>-cluster-ca-cert`, `<name>-kafka-metrics`) and the
  first-internal-listener bootstrap endpoint — kept in lockstep with
  the Terraform module's `locals.tf`
- `module/vars.go`: the `kafka.strimzi.io/v1` apiVersion, the
  `strimzi.io/cluster` label key, the metrics ConfigMap key

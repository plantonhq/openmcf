# KubernetesKafka: Research and Design

## Introduction

KubernetesKafka declares one KRaft-mode Apache Kafka cluster
reconciled by the Strimzi cluster operator
(KubernetesStrimziKafkaOperator must be on the cluster and watching
the namespace; the operator line is pinned at 1.1.0). The spec renders
a `kafka.strimzi.io/v1` `Kafka` custom resource plus one
`KafkaNodePool` per pool entry — node pools, listeners,
broker configuration, authorization, entity operators, Cruise Control,
metrics, CAs, rack awareness, and maintenance windows in one resource.

## The Deployment Landscape

Kafka on Kubernetes without an operator is the classic stateful
anti-pattern: brokers carry node IDs, partition placements, and quorum
membership that no StatefulSet or Helm chart understands. Rolling two
brokers at once can drop a partition below its minimum in-sync
replicas; a Kafka version upgrade has a mandated order; TLS
certificates need rotation that touches every node. Strimzi encodes
that Day-2 expertise into a reconciler — which is why this catalog
splits the concern in two: the operator kind installs the engine, this
kind declares each cluster.

Two upstream shifts define the current architecture:

- **KRaft only.** Strimzi removed ZooKeeper support (and
  ZooKeeper-to-KRaft migration) in 0.46. The metadata quorum runs
  inside Kafka on nodes carrying the `controller` role — there is no
  ZooKeeper anywhere in this spec, and no migration arm to model.
- **The `v1` API and node pools.** Strimzi 1.0 moved the CRDs to the
  served-and-stored `v1` API. Topology is declared exclusively as
  `KafkaNodePool` resources bound to their cluster by the
  `strimzi.io/cluster` label; the spec folds them in as the repeated
  `node_pools` block, and the module renders one CR per entry.

## Upstream Architecture

The operator reconciles the declared resources into:

- **Broker/controller pods** — grouped by pool, each pool with its own
  storage shape, sizing, and scheduling. At least one pool must carry
  the `controller` role and one the `broker` role (spec-validated); a
  single dual-role pool is the dev shape, separate pools the
  production norm (brokers scale without disturbing the quorum;
  controller counts stay odd).
- **The bootstrap Service** `<name>-kafka-bootstrap` — the internal
  client entry point; the module derives the exported
  `internal_bootstrap_endpoint` from the FIRST internal-type listener
  (empty when only external listeners are declared — an honest signal
  that in-cluster clients have no plain path).
- **Certificate Secrets** — the operator generates and renews a
  cluster CA (`<name>-cluster-ca-cert`, key `ca.crt` — exported for
  client truststores) and a clients CA (signs KafkaUser mTLS
  certificates). Validity/renewal windows are the `cluster_ca` /
  `clients_ca` knobs; renewals can be fenced into
  `maintenance_time_windows`.
- **Per-cluster entity operators** — the topic operator and user
  operator (both default-enabled in this spec). They reconcile
  `KafkaTopic` / `KafkaUser` resources — declared through
  KubernetesKafkaTopic / KubernetesKafkaUser — into real topics and
  authenticated principals with credential Secrets.

### The placement contract for topics and users

KafkaTopic and KafkaUser resources must live in the SAME namespace as
their Kafka cluster and carry the `strimzi.io/cluster: <cluster-name>`
label — the entity operators watch only that namespace. The
KubernetesKafkaTopic / KubernetesKafkaUser kinds render both from the
cluster reference; nothing here needs manual labeling. Disabling an
entity operator makes the corresponding declarations silently inert —
the reason both default to enabled.

## Listeners: the Client Surface

Seven exposure types, straight from the upstream listener schema:
`internal` (ClusterIP, the default), `cluster-ip` (per-broker
ClusterIP services for custom proxies), `nodeport`, `loadbalancer`
(one cloud LoadBalancer per broker plus bootstrap), `ingress`
(NGINX-Ingress TLS passthrough), `route` (OpenShift), and `tlsroute`
(Gateway API, new in Strimzi 1.1).

Rules the spec enforces because the operator would otherwise reject
the resource at reconcile time:

- `ingress`, `route`, and `tlsroute` listeners are TLS-passthrough by
  construction — `tls: true` is required.
- `ingress` listeners REQUIRE `configuration.bootstrap.host` AND a
  host per broker (each broker must be individually addressable
  through SNI), plus — outside the manifest — an ingress controller
  with SSL passthrough enabled. Kafka is a binary protocol over TCP;
  passthrough with SNI is how it rides an HTTP-shaped ingress.
  Upstream has deprecated the `ingress` type (the NGINX ingress
  controller it targets was archived upstream in March 2026); it
  remains functional in the pinned release.
- Mutual-TLS authentication requires `tls: true` on the listener.

Authentication per listener: `tls` (client certificates issued via
KubernetesKafkaUser), `scram-sha-512` (username/password Secrets
generated by KubernetesKafkaUser), or `custom` (bring-your-own SASL —
OAuth lost its first-class Strimzi type in 1.x and routes through
custom's `sasl` + `listener_config`).

### Environment injection: cloud annotations ride the listeners

This component calls no cloud APIs, but external listeners are where
managed-Kubernetes integration rides — annotations on
`configuration.bootstrap.annotations` and
`configuration.brokers[].annotations` are the component's
environment-injection surface, playing the role that
credential/identity blocks play in components that do call cloud APIs.

| Cloud / posture | Where | Annotations |
|---|---|---|
| AWS NLB (via the AWS Load Balancer Controller) | bootstrap + each broker | `service.beta.kubernetes.io/aws-load-balancer-type: external` |
| DNS automation (any cloud, external-dns) | bootstrap + each broker | `external-dns.alpha.kubernetes.io/hostname: <name>` — one hostname for bootstrap, one per broker |
| Client IP preservation | `configuration.external_traffic_policy` | `Local` (no extra hop) vs `Cluster` (the Kubernetes default) |
| Source fencing | `configuration.load_balancer_source_ranges` | CIDR list rendered as `loadBalancerSourceRanges` |

Per-broker `advertised_host` / `advertised_port` overrides complete
the picture for NAT/proxy topologies: the LoadBalancer annotations
provision the path in, the advertised addresses are what brokers tell
clients to reconnect to.

## Broker Configuration: Strings, and Who Owns What

`config` entries are Kafka configuration strings — numbers and
booleans are written as strings ("3", "false"); the operator
serializes every value into Java properties form. The operator OWNS
listener, node identity, security, and quorum configuration: entries
with the prefixes `listeners`, `advertised.`, `broker.`, `listener.`,
`host.name`, `port`, `inter.broker.listener.name`, `sasl.`, `ssl.`,
`security.`, `password.`, `log.dir`, `authorizer.`, `super.user`,
`node.id`, `process.roles`, `controller.` (and the Cruise Control
metrics reporter keys) are IGNORED with an operator log warning, not
applied — those concerns have typed fields.

### The durability story

Topic durability is governed by per-topic replication and two cluster
entries this spec recommends: `default.replication.factor: "3"` and
`min.insync.replicas: "2"` (plus the matching internal-topic entries).
A 3-broker cluster with those settings survives one broker loss
without losing acknowledged writes — producers using acks=all get
their acknowledgment only after two replicas hold the write, so the
failure of any single broker leaves at least one acknowledged copy.
RF 1 (the dev shape) makes every broker restart a data-availability
event; the presets are explicit about which side of that line they
sit on.

## Authorization

Omitted = no authorizer: every authenticated (or anonymous, on a
no-auth listener) client can do everything. `simple` enables Kafka's
built-in ACL authorizer and enforces the ACLs declared on
KubernetesKafkaUser resources — from that moment, clients WITHOUT
matching ACLs are denied, so `super_users` must list real operational
principals (Kafka principal form: `User:CN=platform-admin` for TLS
users, `User:admin` for SCRAM users) BEFORE authorization turns on.
`custom` takes a fully-qualified authorizer class (on the broker
classpath via a custom image): Keycloak authorization and the OPA
authorizer lost their first-class Strimzi types in 1.x — the OPA
plugin is no longer bundled in the images — and both route through
`custom` now.

## Design Decisions

- **Untyped CustomResources on both engines.** The CRD types
  `spec.kafka.config` (and the listener and Cruise Control
  configuration blocks) with `x-kubernetes-preserve-unknown-fields`,
  which generated typed SDKs flatten into shapes that cannot hold the
  free-typed bodies — so the Pulumi module renders untyped
  `apiextensions.CustomResource` bodies and the Terraform module
  renders `kubectl_manifest` resources, kept as exact twins (same keys
  rendered and omitted, numbers as numbers, booleans as booleans).
- **`kubectl_manifest` over `kubernetes_manifest`.** The alekc/kubectl
  provider needs no cluster connection at plan time — an infra chart
  can plan the operator and its Kafka clusters in one run, before the
  CRDs exist.
- **Pools apply before the Kafka CR.** Strimzi tolerates either order,
  but a Kafka CR with no matching pools reports a transient warning
  state; the modules sequence pools first. Pools are keyed by their
  OWN NAME in Terraform state so list reorderings never churn
  addresses.
- **No await machinery, deliberately.** Cluster readiness depends on
  the operator (image pulls, KRaft quorum formation, listener
  provisioning) that is not part of applying the resources — the
  never-block-on-a-controller posture of every operator-CR kind in
  the catalog.
- **The module owns the metrics ConfigMap.** `metrics.enabled` renders
  the canonical Strimzi JMX Prometheus rules as
  `<name>-kafka-metrics` and wires it as `metricsConfig` (port 9404 in
  the pods); the CR only points at it. Pair with `kafka_exporter` for
  consumer-lag metrics the JMX exporter cannot see.
- **`node_selector` translates to node affinity.** The Strimzi pod
  template carries affinity and tolerations but no nodeSelector — the
  modules render a required node affinity with one matchExpressions
  entry per label (semantically identical for exact-match selection),
  sorted for determinism.
- **KafkaRebalance is an operational verb, not a declared resource.**
  Cruise Control deploys from this spec; rebalance operations are
  applied with kubectl or automation when needed.
- **CAs are operator-managed.** Bringing your own cluster/clients CA
  is deliberately not modeled — the operator renews self-generated CAs
  automatically, and the per-listener `broker_cert_chain_and_key` (the
  cert-manager seam) covers the actual client-trust need. The
  certificate's SANs must cover the listener's advertised names — the
  operator does NOT validate this.

## Deliberately Unmodeled

Kept off the typed spec on purpose; every item remains reachable by
declaring the raw `Kafka` CR through KubernetesManifest:

- **Tiered storage** — upstream carries a `tieredStorage` block
  requiring a custom RemoteStorageManager plugin in a custom image; a
  bring-your-own-plugin surface with no portable story.
- **The quotas plugin** — upstream's `quotas` block configures the
  Strimzi quotas plugin; workload-specific tuning that composes badly
  with a typed cross-environment spec.
- **Per-listener `networkPolicyPeers`** — network fencing is a
  cluster-level policy concern (and the operator already generates
  baseline NetworkPolicies, toggleable at the operator kind); a
  per-listener peer list embeds policy in the workload.
- **`parentRefs` on tlsroute listeners** — the Gateway API attachment
  point for the new tlsroute type; the Gateway composition story
  belongs to the gateway kinds, not inline here.
- **Custom pod templates beyond node_selector/tolerations** — the
  upstream `template` trees (pod metadata, security contexts, init
  containers, affinity in full generality) are an unbounded surface;
  the spec models the two scheduling knobs that cover the real
  placement need.
- **Remote JMX options** — `jmxOptions` opens an authenticated remote
  JMX port on brokers; the modeled JMX Prometheus metrics path covers
  observability without exposing a management protocol.

## Version Pins and Naming Contracts

| What | Value | Notes |
|---|---|---|
| CR API | `kafka.strimzi.io/v1` | The only served API from Strimzi 1.0 onward |
| Cluster binding label | `strimzi.io/cluster` | Stamped on every KafkaNodePool (and by the topic/user kinds) |
| Bootstrap Service | `<name>-kafka-bootstrap` | Exported as `bootstrap_service_name` |
| Cluster CA Secret | `<name>-cluster-ca-cert` (key `ca.crt`) | Exported as `cluster_ca_cert_secret_name` |
| Metrics ConfigMap | `<name>-kafka-metrics` (key `kafka-metrics-config.yml`) | Module-owned, rendered when `metrics.enabled` |
| `kafka_version` | empty = operator default | Strimzi 1.1 supports Kafka 4.3.0 and 4.2.1 |
| `metadata_version` | empty = pinned to the Kafka version | One-way door during upgrades — hold at the old format until every node runs new binaries |

## IaC Twins

Pulumi (untyped CustomResources, `module/kafka.go` +
`module/nodepools.go`) and Terraform (`kubectl_manifest` + null-prune
locals, `locals.tf`) render identical CR bodies and the same
module-owned metrics ConfigMap. Keep `locals.go`/`kafka.go`/
`nodepools.go` and `locals.tf` in lockstep.

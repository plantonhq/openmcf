# Kubernetes Kafka UI

Deploys the kafbat UI console for Kafka from the served `kafka-ui`
Helm chart — one installation observing and managing many clusters:
browse topics and live messages, inspect consumer groups and lag,
view and register schemas through a connected registry, and monitor
Connect pipes. Cluster, schema-registry and Connect wiring compose
directly from KubernetesKafka, KubernetesKarapace and
KubernetesKafkaConnect siblings; credentials ride Secret-backed
environment variables, never rendered configuration.

## What Gets Created

- **Namespace** (optional) — created and owned when `create_namespace`
  is set
- **Console Secret** (`<name>-secrets`, only when `auth` is declared)
  — the module-materialized home of the console login password,
  wired into the app through the chart's `envs.secretMappings`
- **Helm release** (`kafka-ui` at https://ui.charts.kafbat.io, pinned
  1.6.4) — release name and chart fullname pinned to `metadata.name`,
  the typed spec rendered into chart values (the app's cluster
  config, service, sizing, TLS volume mounts), `helm_values` merged
  last

The chart renders the console Deployment, Service, and the app-config
ConfigMap; the deploy waits for the console to become Ready (wait +
atomic).

## Prerequisites

- Reachable Kafka clusters (`clusters[].bootstrap_servers` — sibling
  references or literal addresses); the console deploys fine before
  its clusters exist, it just reports them offline
- For TLS clusters: the CA Secret in this namespace (a Strimzi
  cluster-CA Secret works directly — PEM truststore, no conversion)
- For SASL clusters: the credential Secret (a KubernetesKafkaUser's
  operator-generated Secret by reference)

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesKafkaUi
metadata:
  name: kafka-console
spec:
  namespace:
    value: kafka-console
  create_namespace: true
  clusters:
    - name: dev
      bootstrap_servers:
        value: dev-kafka-kafka-bootstrap.dev-kafka.svc.cluster.local:9092
```

Reach it without any exposure via the exported port-forward command
(`kubectl port-forward svc/kafka-console -n kafka-console 8080:80`).

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the console runs in |
| `service_name` | Console Service (`<name>`, pinned via fullnameOverride) |
| `endpoint` | In-cluster endpoint (`http://<name>.<namespace>.svc.cluster.local:<port>`) |
| `port_forward_command` | Workstation access without any exposure |

## Next Steps

Enable `auth` (the single login-form account — multi-user, OAuth2 and
LDAP ride `helm_values`) before exposing the console, and mark
production clusters `read_only`. Wire the full stack: a
KubernetesKarapace endpoint as `schema_registry.url` for schema
browsing, and each KubernetesKafkaConnect's REST endpoint under
`kafka_connect` for pipe monitoring. Expose by composing a
KubernetesIngress or Gateway route against the `service_name`
output.

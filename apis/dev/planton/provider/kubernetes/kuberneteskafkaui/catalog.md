# Kafka UI

Deploys kafbat UI — the Apache-2.0 web console for Kafka — from the served `kafka-ui` Helm chart. One installation observes and manages **many** clusters: browse topics and live messages, inspect consumer groups and lag, view schemas through a connected registry, and monitor Connect pipes. Each `clusters` entry wires one Kafka cluster (plus its optional schema registry and Connect clusters) into the console; the foreign keys compose directly with KubernetesKafka, KubernetesKarapace, and KubernetesKafkaConnect siblings.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Helm Release** (`<metadata.name>`) — the kafbat UI chart with the Service pinned to the resource name (`fullnameOverride`), so outputs stay deterministic and several consoles can coexist in one cluster
- **Console Service** — ClusterIP by default; compose KubernetesIngress or Gateway routes against the exported `service_name` / `endpoint` handles for shared access
- **Module-owned Secrets** — the console login password (when declared) materializes into `<name>-secrets`; referenced cluster SASL and basic-auth passwords mount from their source Secrets — nothing sensitive lands in rendered chart values
- **Namespace** (optional) — created with standard governance labels when `create_namespace` is true

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** — an active connection in the Connect module with credentials for the target cluster.

### Kafka Family Side

- At least one **Kubernetes Kafka** cluster to observe (bootstrap endpoint wired by reference or literal)
- Optional **KubernetesKarapace** for schema browsing and **KubernetesKafkaConnect** for connector monitoring — wired per cluster entry
- Optional **KubernetesKafkaUser** for SASL credentials — reference its operator-generated Secret; never type passwords into the console manifest

## Deploy

### Console

Open the deployment store, find **Kafka UI**, and click **Deploy**. The creation wizard walks you through the console's namespace, the served-chart pin, the cluster composition step (where the whole Kafka family's outputs meet), console login, Service exposure, sizing, scheduling, and the Helm-values escape hatch. Start from the **Single cluster readonly** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesKafkaUi
metadata:
  name: kafka-console
  org: acme-corp
  env: prod
spec:
  namespace:
    value: kafka-console
  create_namespace: true
  clusters:
    - name: production
      bootstrap_servers:
        value: events-kafka-bootstrap.kafka.svc.cluster.local:9092
      read_only: true
```

```shell
planton apply -f kafka-console.yaml
```

Reach the console through the exported port-forward command until you compose exposure — and add `login_form` before anything shared can reach the Service.

### InfraChart

Wire sibling outputs so the family cannot drift:

```yaml
spec:
  namespace:
    value: kafka-console
  clusters:
    - name: events
      bootstrap_servers:
        valueFrom:
          kind: KubernetesKafka
          name: event-bus
          fieldPath: status.outputs.internal_bootstrap_endpoint
      tls:
        ca_secret_name:
          valueFrom:
            kind: KubernetesKafka
            name: event-bus
            fieldPath: status.outputs.cluster_ca_cert_secret_name
      sasl:
        mechanism: SCRAM-SHA-512
        username: kafka-ui
        password_secret:
          secret_name:
            valueFrom:
              kind: KubernetesKafkaUser
              name: kafka-ui-user
              fieldPath: status.outputs.secret_name
      schema_registry:
        url:
          valueFrom:
            kind: KubernetesKarapace
            name: schema-registry
            fieldPath: status.outputs.endpoint
      kafka_connect:
        - name: cdc
          address:
            valueFrom:
              kind: KubernetesKafkaConnect
              name: cdc-connect
              fieldPath: status.outputs.rest_api_endpoint
```

## Key Configuration

These are the most important decisions when configuring the console. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Open vs login-gated console** — omitting `auth` means NO authentication: anyone who can reach the Service has full console access on every wired cluster. Acceptable only for cluster-internal evaluation paired with ClusterIP and no composed route. The typed login is exactly ONE `login_form` account; OAuth2/OIDC/LDAP compose through `helm_values`.

**Per-cluster `read_only`** — an app-side switch that hides every mutating action (topic create/delete, message produce, config edits) for that cluster. The right posture for production clusters on a shared console — set it where the risk lives, not globally.

**Exposure composes, never embeds** — prefer ClusterIP + a composed Ingress/Gateway against the exported handles. NodePort/LoadBalancer exist as Service knobs, not a hostname story.

**Credentials never land in rendered configuration** — SASL and basic-auth passwords are Secret references; the console login password is the one literal the module materializes into a Secret-backed environment variable.

**`helm_values` is the escape hatch** — YAML merged LAST over everything the typed fields render: probes, security contexts, OAuth2/LDAP login. Never for secrets.

## Outputs and Dependencies

### What This Component Consumes

| Field | References | Purpose |
|-------|-----------|---------|
| `spec.namespace` | KubernetesNamespace (`spec.name`) | The console's own home — distinct from observed clusters' namespaces |
| `spec.clusters[].bootstrap_servers` | KubernetesKafka (`status.outputs.internal_bootstrap_endpoint`) | The cluster the console observes |
| `spec.clusters[].tls.ca_secret_name` | KubernetesKafka (`status.outputs.cluster_ca_cert_secret_name`) | TLS trust for the cluster connection |
| `spec.clusters[].sasl.password_secret.secret_name` | KubernetesKafkaUser (`status.outputs.secret_name`) | SASL credentials for secured listeners |
| `spec.clusters[].schema_registry.url` | KubernetesKarapace (`status.outputs.endpoint`) | Schema browsing for the cluster |
| `spec.clusters[].kafka_connect[].address` | KubernetesKafkaConnect (`status.outputs.rest_api_endpoint`) | Connect pipe monitoring |

### What This Component Provides

After provisioning, `status.outputs` contains values downstream Cloud Resources can consume:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `namespace` | Namespace the console runs in | Composed route placement context |
| `service_name` | Name of the console Service | KubernetesIngress / Gateway backend target |
| `endpoint` | In-cluster HTTP endpoint | Internal tooling links |
| `port_forward_command` | Copyable kubectl port-forward | Quick look with no exposure |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Single cluster readonly** — one cluster wired with `read_only` on, no login, no exposure: the safe first console for observing production. Start from the **Single cluster readonly** preset.

**Full stack console** — TLS + SCRAM cluster connection, schema registry, Connect monitoring, and a login gate. Start from the **Full stack console** preset.

**Multi-cluster org console** — staging with full powers and production locked to observe-only on one shared pane. Start from the **Multi cluster** preset.

## Works With

- **Kubernetes Kafka** — bootstrap endpoint and cluster CA for each wired cluster
- **Kubernetes Kafka User** — SASL credential Secrets for secured listeners
- **Kubernetes Karapace** — schema registry endpoint for schema-aware browsing
- **Kubernetes Kafka Connect** — REST endpoints for connector monitoring
- **Kubernetes Ingress / Gateway API** — shared exposure composed against the exported Service handles

# KubernetesKafkaConnect Terraform Module

Deploys one Strimzi-operator-managed Kafka Connect cluster: the
optional namespace, the optional JMX metrics rules ConfigMap, and the
`kafka.strimzi.io/v1` `KafkaConnect` CR itself. The typed spec
renders into the CR body in `locals.tf` — the exact twin of the
Pulumi module's `connect.go` (same keys rendered and omitted, numbers
as numbers, booleans as booleans).

## Module Behavior

- **The Connect cluster name is `metadata.name`** — the KafkaConnect
  CR name and the value KubernetesKafkaConnector resources bind to
  (rendered as their `strimzi.io/cluster` label).
- **The group identity quartet always renders** — group.id and the
  three storage topics resolve from spec overrides with
  metadata.name-derived fallbacks (`<name>`,
  `<name>-connect-{configs,status,offsets}`), keeping the uniqueness
  contract visible in the applied object.
- **The use-connector-resources annotation is module-owned and
  unconditional** — `strimzi.io/use-connector-resources: "true"`
  switches connector management to declarative mode; the operator
  reverts REST-API-made changes.
- **CRs apply through `kubectl_manifest` (alekc/kubectl)** — unlike
  the hashicorp provider's `kubernetes_manifest` it needs no cluster
  connection at plan time, so the Connect cluster can be PLANNED
  before the Strimzi operator's CRDs exist.
- **No wait_for block, deliberately** — worker readiness depends on
  the operator (image pulls or an operator-driven image build, group
  formation) that is not part of applying the resources.
- **The module (not the operator) owns the metrics ConfigMap** —
  `metrics.enabled` renders the canonical Strimzi connect JMX rules
  as `<name>-connect-metrics` (`kubernetes_config_map_v1`).
- **The module owns namespace creation** — `create_namespace` drives
  a `kubernetes_namespace_v1` resource with the standard governance
  labels.

## Rendering Quirks

- **`image` renders only on the prebuilt arm** — image and build are
  mutually exclusive at validation.
- **Per-arm authentication** — each type renders only its own
  credential fields, with the KubernetesKafkaUser credential-Secret
  fallbacks (`user.crt`/`user.key`, password key `password`) applied
  via coalesce.
- **Per-type build artifacts** — maven artifacts render coordinates
  only, url-family artifacts render url/sha512sum/insecure (+
  fileName for `other`) — only the keys the operator's per-type
  sub-schema allows.
- **OCI plugin artifacts always carry `type = "image"`** — the only
  artifact type the plugins arm supports; module-owned.
- **`node_selector` translates to node affinity** — the Strimzi pod
  template carries no nodeSelector; one matchExpressions entry per
  label, sorted by key (the Pulumi module renders the same
  translation).
- **Null-prune idiom throughout** — conditional entries are written
  as `key = cond ? value : null` inside one object literal and
  pruned, so numbers and booleans keep their types in the rendered
  CR.

## Resources

| Resource | Condition |
|---|---|
| `kubernetes_namespace_v1.namespace` | `spec.create_namespace` |
| `kubernetes_config_map_v1.connect_metrics` | `spec.metrics.enabled` |
| `kubectl_manifest.connect` | always |

## Usage

```bash
planton tofu apply --manifest connect.yaml
```

## Local Development

```bash
terraform init
terraform validate
terraform plan -var-file=terraform.tfvars.json
terraform apply -var-file=terraform.tfvars.json
```

## Inputs

See `variables.tf` for the full variable specification (generated
from the spec proto). The spec arrives from the proto→tfvars
converter in snake_case with the `namespace`, `bootstrap_servers`,
trusted-certificate and credential-Secret references resolved to
literal strings before Terraform runs.

## Outputs

| Output | Description |
|--------|-------------|
| `namespace` | Namespace the Connect cluster deploys into |
| `connect_name` | Connect cluster name (`metadata.name`) — the `strimzi.io/cluster` binding value |
| `rest_api_service_name` | Connect REST API Service (`<name>-connect-api`) |
| `rest_api_endpoint` | In-cluster REST endpoint (port 8083) — read-only inspection |

## Parity

Kept in lockstep with the Pulumi module (`../pulumi/module/`): same
`kafka.strimzi.io/v1` apiVersion, same CR body (group identity
resolution, per-arm authentication, per-type build artifacts, the
affinity translation), same module-owned metrics ConfigMap and
unconditional use-connector-resources annotation, same outputs.

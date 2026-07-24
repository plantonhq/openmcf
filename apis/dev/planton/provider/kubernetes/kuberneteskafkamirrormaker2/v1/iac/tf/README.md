# KubernetesKafkaMirrorMaker2 Terraform Module

Deploys one Strimzi-operator-managed MirrorMaker 2 replication
engine: the optional namespace, the optional JMX metrics rules
ConfigMap, and the `kafka.strimzi.io/v1` `KafkaMirrorMaker2` CR
itself. The typed spec renders into the CR body in `locals.tf` — the
exact twin of the Pulumi module's `mirrormaker2.go` (same keys
rendered and omitted, numbers as numbers, booleans as booleans).

## Module Behavior

- **The deployment name is `metadata.name`** — the CR name, and the
  root of the engine's group identity: group.id and the three storage
  topics resolve from spec overrides with metadata.name-derived
  fallbacks (`<name>-mirrormaker2-{configs,status,offsets}`), unique
  per Connect-protocol workload sharing the target cluster.
- **The target body always renders alias + group identity** — target
  alias coalesced to `target`; the spec's CEL rules guarantee it
  differs from every source alias.
- **CRs apply through `kubectl_manifest` (alekc/kubectl)** — no
  cluster connection at plan time, so the engine can be planned
  before the Strimzi CRDs exist.
- **No wait_for block, deliberately** — engine readiness depends on
  the operator (image pulls, worker group formation, connector
  startup).
- **The module owns the metrics ConfigMap** — `metrics.enabled`
  renders the canonical Strimzi connect JMX rules as
  `<name>-mm2-metrics` (`kubernetes_config_map_v1`).
- **The module owns namespace creation** — `create_namespace` drives
  a `kubernetes_namespace_v1` resource with the standard governance
  labels.

## Rendering Quirks

- **Shared client shapes** — target and source TLS/authentication
  render identically to the KubernetesKafkaConnect module: per-type
  credential fields only, with the KubernetesKafkaUser
  credential-Secret fallbacks (`user.crt`/`user.key`, password key
  `password`) applied via coalesce.
- **Mirror bodies render only what is set** — scope patterns when
  non-empty, connector bodies with tasksMax as a number, config, and
  autoRestart.
- **`node_selector` translates to node affinity** — the Strimzi pod
  template carries no nodeSelector; one matchExpressions entry per
  label, sorted by key (the Pulumi module renders the same
  translation).
- **Null-prune idiom throughout** — conditional entries are
  `key = cond ? value : null` inside one object literal and pruned,
  so numbers and booleans keep their types in the rendered CR.

## Resources

| Resource | Condition |
|---|---|
| `kubernetes_namespace_v1.namespace` | `spec.create_namespace` |
| `kubernetes_config_map_v1.mirrormaker2_metrics` | `spec.metrics.enabled` |
| `kubectl_manifest.mirrormaker2` | always |

## Usage

```bash
planton tofu apply --manifest mirrormaker2.yaml
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
converter in snake_case with the `namespace`, bootstrap,
trusted-certificate and credential-Secret references resolved to
literal strings before Terraform runs.

## Outputs

| Output | Description |
|--------|-------------|
| `namespace` | Namespace the deployment runs in |
| `mirrormaker_name` | The deployment's name (`metadata.name`) |
| `rest_api_endpoint` | In-cluster Connect REST endpoint of the engine (port 8083) — read-only inspection |

## Parity

Kept in lockstep with the Pulumi module (`../pulumi/module/`): same
`kafka.strimzi.io/v1` apiVersion, same CR body (target group-identity
resolution, mirror/source/connector bodies, shared client shapes, the
affinity translation), same module-owned metrics ConfigMap, same
outputs.

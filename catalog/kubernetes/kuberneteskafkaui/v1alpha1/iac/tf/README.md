# KubernetesKafkaUi Terraform Module

Deploys the kafbat UI console from the served `kafka-ui` Helm chart
(https://ui.charts.kafbat.io, default pin 1.6.4): the optional
namespace, the optional module-materialized console Secret, and a
`helm_release` with the typed spec rendered into chart values and the
`helm_values` escape hatch merged LAST (native provider semantics:
`values = [yamlencode(typed), helm_values]`). Every rendering has an
exact twin in the Pulumi module (`../pulumi/module/`).

## Module Behavior

- **Release name and chart fullname are `metadata.name`**
  (`fullnameOverride`) — the Service IS the resource name; several
  consoles coexist per cluster and outputs stay deterministic.
- **No credential ever lands in rendered values** — the app config
  ships in a ConfigMap, so every password position carries a
  `${ENV_VAR}` placeholder and `envs.secretMappings` wires each env
  var to its Secret key (secretKeyRef). Referenced credentials map at
  their source Secrets; the console login password lives in the
  module-materialized `<name>-secrets` Secret
  (`kubernetes_secret_v1.console`, key `console-user-password`).
- **LOGIN_FORM is a single account** — Spring Boot's default security
  user (the app registers no user store); auth renders `DISABLED`
  explicitly when the block is absent.
- **PEM truststores** — one secret volume per TLS cluster mounted at
  `/etc/kafkaui/cluster-<i>-ca`; client properties set
  `ssl.truststore.type=PEM` with the module-owned
  `security.protocol` and JAAS line.
- **Wait-for-ready, atomic** — `wait = true`, atomic, 600s timeout: a
  console that never starts fails the apply.
- **The module owns namespace creation** — `create_namespace` drives
  `kubernetes_namespace_v1` with the standard governance labels (the
  release itself never creates namespaces).

## Resources

| Resource | Condition |
|---|---|
| `kubernetes_namespace_v1.kafka_ui` | `spec.create_namespace` |
| `kubernetes_secret_v1.console` | `spec.auth` declared |
| `helm_release.kafka_ui` | always |

## Usage

```bash
planton tofu apply --manifest kafka-ui.yaml
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
CA/password-Secret, `schema_registry.url` and `kafka_connect[].address`
references resolved to literal strings before Terraform runs.

## Outputs

| Output | Description |
|--------|-------------|
| `namespace` | Namespace the console runs in |
| `service_name` | Console Service (`<name>`, pinned via fullnameOverride) |
| `endpoint` | In-cluster endpoint (`http://<name>.<namespace>.svc.cluster.local:<port>`) |
| `port_forward_command` | Workstation access without any exposure |

## Parity

Kept in lockstep with the Pulumi module: same chart identity and
1.6.4 default pin (verified against the served repository index),
same rendered values documents (app config with deterministic
`KAFKA_CLUSTER_<i>_...` placeholders, secretMappings, TLS
volumes/mounts, resolved service/replica defaults), same
helm_values-merges-last semantics, same wait/atomic posture, same
outputs.

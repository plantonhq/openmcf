# KubernetesKarapace Terraform Module

Deploys a Karapace schema registry as module-owned native manifests —
upstream ships no Helm chart or operator, so the module IS the
deployment: `kubernetes_deployment_v1` / `kubernetes_service_v1`
resources per role, `KARAPACE_*` environment-variable configuration
on upstream's own image, and Secret-mounted file material. Every
rendering has an exact twin in the Pulumi module
(`../pulumi/module/`).

## Module Behavior

- **Names derive from `metadata.name`** — the registry Deployment and
  Service are `<name>`, the REST-proxy pair `<name>-rest`, the
  module-materialized SASL Secret `<name>-sasl`.
- **Image pin**: `ghcr.io/aiven-open/karapace:6.2.1` when
  `spec.image` is empty (the `default_image` local; bump in lockstep
  with the Pulumi module's `vars.KarapaceImage`).
- **Per-pod advertised identity**: `KARAPACE_ADVERTISED_HOSTNAME`
  injects each pod's IP via the downward API (`status.podIP`) —
  followers forward writes to the leader's advertised address, so a
  shared Service name or an unresolvable pod name would break
  forwarding.
- **The SASL password always arrives via secretKeyRef** — a
  referenced `password_secret` wires directly; a literal `password`
  drives `kubernetes_secret_v1.sasl_password` first. Never a
  plaintext env value.
- **server_tls flips the coupled trio** — advertised protocol https,
  cert/key file paths under `/etc/karapace/server-tls`, and the probe
  scheme. The REST proxy always serves plain HTTP.
- **Probes hit `/_health`** — unauthenticated by engine design
  (skip-auth list), with startup headroom for schemas-topic replay.
- **Per-role selector identity** — the role-specific `app` label
  keeps each Service from selecting the other role's pods.
- **Scheduling is registry-scoped** — node_selector/tolerations apply
  to the registry Deployment only, per the spec contract.

## Resources

| Resource | Condition |
|---|---|
| `kubernetes_namespace_v1.this` | `spec.create_namespace` |
| `kubernetes_secret_v1.sasl_password` | literal `kafka.sasl.password` declared |
| `kubernetes_deployment_v1.registry` | always |
| `kubernetes_service_v1.registry` | always |
| `kubernetes_deployment_v1.rest_proxy` | `spec.rest_proxy.enabled` |
| `kubernetes_service_v1.rest_proxy` | `spec.rest_proxy.enabled` |

## Usage

```bash
planton tofu apply --manifest karapace.yaml
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
CA/client-certificate/password Secret and `server_tls` references
resolved to literal strings before Terraform runs.

## Outputs

| Output | Description |
|--------|-------------|
| `namespace` | Namespace the registry runs in |
| `service_name` | Registry Service (`<name>`) |
| `endpoint` | In-cluster endpoint — the `schema.registry.url` value |
| `rest_proxy_endpoint` | REST-proxy endpoint; empty when the role is off |
| `schemas_topic` | The Kafka topic storing the schemas |

## Parity

Kept in lockstep with the Pulumi module: same image pin and
entrypoints (`python3 -m karapace` /
`python3 -m karapace.kafka_rest_apis` with the role flags), same env
var sets and mount paths (byte-identical: `/etc/karapace/kafka-ca`,
`/etc/karapace/kafka-cert`, `/etc/karapace/server-tls`,
`/etc/karapace/auth`), same probes, naming contracts, and outputs.

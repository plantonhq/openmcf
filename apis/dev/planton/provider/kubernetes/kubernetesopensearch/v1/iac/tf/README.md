# KubernetesOpenSearch Terraform Module

Deploys one operator-managed OpenSearch cluster: the optional namespace
and the `opensearch.opster.io/v1` OpenSearchCluster CR. The custom
resource applies through `kubectl_manifest` (alekc/kubectl provider,
server-side apply), which needs no cluster connection at plan time —
the cluster can be planned before the OpenSearch operator's CRDs exist.

Prerequisites at apply time: the OpenSearch Kubernetes Operator
(`KubernetesOpenSearchOperator`) on the cluster; the Prometheus
Operator's ServiceMonitor CRD when `monitoring` is enabled.

## Module Behavior

- **One CR, everything else operator-created** — node StatefulSets
  (`<name>-<pool>`), the main Service `<name>` (the module pins
  `general.serviceName` to `metadata.name` — the `service_name` output
  depends on this), the discovery Service `<name>-discovery`, TLS
  Secrets, the admin bootstrap Secret `<name>-admin-password`, and the
  optional Dashboards deployment + Service `<name>-dashboards`. There
  are NO ingress resources by design — exposure composes from
  first-class kinds referencing the exported handles.
- **`http_endpoint` is always https** — the operator's HTTP layer
  serves TLS in EVERY posture: with `security.tls.http` it uses
  operator-generated (or provided) certificates; with `spec.security`
  entirely absent the TLS reconciler generates nothing (upstream
  `pkg/reconcilers/tls.go`) and the opensearchproject image's demo
  security configuration serves TLS instead — the operator itself
  always connects over https (`pkg/builders/cluster.go` URLForCluster,
  and the node readiness probe curls `https://localhost`).
- **`admin_credentials_secret_name`** is the operator-generated
  `<name>-admin-password` Secret (fields `username`/`password`) —
  exported EMPTY when `security.config` is set: a custom security
  config replaces the operator bootstrap and the user's
  `admin_credentials_secret` is authoritative.
- **Unset optionals are omitted** — the rendered CR body is null-pruned
  so the apiserver applies the CRD's own defaults. Presence-sensitive
  booleans whose spec default (true) diverges from the CRD default
  (false) — `setVMMaxMapCount`, TLS `generate`/`perNode` — render
  explicitly when true and are omitted when false.
- **PDB bounds follow intstr semantics** — an integer string ("2")
  renders as a YAML number, a percentage ("25%") as a string
  (`try(tonumber(v), v)` — a conditional would unify the number/string
  branches to string); the Pulumi twin applies strconv.Atoi for the
  identical result.
- **The CRD's ImageSpec takes ONE image string** — the shared
  ContainerImage folds into `repo:tag`; `imagePullSecrets` come from
  the spec's own list.
- **The module (not the operator) owns namespace creation** —
  `create_namespace` drives a `kubernetes_namespace_v1` resource
  carrying the standard governance labels.

## Resources

| Resource | Condition |
|---|---|
| `kubernetes_namespace_v1.namespace` | `spec.create_namespace` |
| `kubectl_manifest.opensearch_cluster` | always |

## Usage

```bash
planton tofu apply --manifest kubernetes-open-search.yaml
```

## Local Development

```bash
tofu init -backend=false
tofu validate
tofu plan -var-file=terraform.tfvars.json
tofu apply -var-file=terraform.tfvars.json
```

The full-surface hack manifest for offline proofs lives in `../hack/`.

## Inputs

See `variables.tf` for the full variable specification (generated from
the spec proto). The spec arrives from the proto→tfvars converter in
snake_case with every `StringValueOrRef` foreign key — `namespace`
(KubernetesNamespace), pool `persistence.pvc.storage_class`
(KubernetesStorageClass), the TLS secrets
(KubernetesCertificate/KubernetesSecret), the security-config,
keystore, dashboards and monitoring secrets (KubernetesSecret) —
resolved to a literal string before Terraform runs.

## State Import

Existing deployments can be adopted into state. `kubectl_manifest` uses
the composed import ID `apiVersion//kind//name//namespace`; the CR name
is deterministic (`metadata.name`), so the address derives blind.

## Outputs

| Output | Description |
|--------|-------------|
| `namespace` | Namespace the cluster runs in |
| `cluster_name` | Name of the OpenSearchCluster resource (equals `metadata.name`) |
| `service_name` | The cluster's main Service (= `metadata.name`) |
| `http_endpoint` | In-cluster HTTP API endpoint (`https://<name>.<namespace>.svc.cluster.local:<http_port>`) |
| `admin_credentials_secret_name` | `<name>-admin-password` — empty when a custom security config is declared |
| `dashboards_service_name` | `<name>-dashboards` — empty when dashboards are not enabled |
| `dashboards_endpoint` | In-cluster Dashboards endpoint on port 5601 — empty when dashboards are not enabled |
| `port_forward_command` | Port-forward command for workstation access |

## Parity

This module is the behavioral twin of the Pulumi module
(`../pulumi/module/`): same rendered CR body (null-pruned,
presence-sensitive booleans on divergence, intstr PDB bounds), same
module-owned constants (serviceName = metadata.name, vendor
`opensearch`, Dashboards version defaulting to the cluster version),
same outputs — keep them in lockstep.

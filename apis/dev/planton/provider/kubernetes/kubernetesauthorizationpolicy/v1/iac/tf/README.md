# KubernetesAuthorizationPolicy Terraform Module

Creates a namespaced Istio `AuthorizationPolicy` via the `kubectl_manifest`
resource. The Istio CRDs must already be installed on the target cluster (see
`KubernetesIstioBaseCrds`), a running istiod is required to enforce the policy in the
data plane (see `KubernetesIstio`), and the target namespace must exist (see
`KubernetesNamespace`).

## Usage

```bash
planton tofu apply --manifest authorizationpolicy.yaml
```

## Local Development

```bash
terraform init
terraform validate
terraform plan -var-file=terraform.tfvars.json
terraform apply -var-file=terraform.tfvars.json
```

## Inputs

See `variables.tf` for the variable specification. `variable "spec"` is typed `any`
and passed through verbatim: the platform resolves the `StringValueOrRef` foreign keys
(`namespace`, `target_refs[].name`) to literal strings before Terraform runs, and
emits the manifest-shaped (camelCase, null-pruned) spec, so unset fields are omitted
from the manifest and upstream defaults flow through (e.g. an absent `action` becomes
ALLOW). `selector` and `target_refs` are mutually exclusive.

## Outputs

| Output | Description |
|--------|-------------|
| `authorization_policy_name` | Name of the created AuthorizationPolicy (equals `metadata.name`) |
| `namespace` | Namespace the AuthorizationPolicy was created in |

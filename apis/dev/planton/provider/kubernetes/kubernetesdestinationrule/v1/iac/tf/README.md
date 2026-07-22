# KubernetesDestinationRule Terraform Module

Creates a namespaced Istio `DestinationRule` via the `kubectl_manifest` resource. The
Istio CRDs must already be installed on the target cluster (see `KubernetesIstioBaseCrds`),
a running istiod is required to apply the policy (see `KubernetesIstio`), and the target
namespace must exist (see `KubernetesNamespace`).

## Usage

```bash
planton tofu apply --manifest destinationrule.yaml
```

## Local Development

```bash
terraform init
terraform validate
terraform plan -var-file=terraform.tfvars.json
terraform apply -var-file=terraform.tfvars.json
```

## Inputs

See `variables.tf` for the variable specification. `variable "spec"` is typed `any` and
passed through verbatim: the platform's proto-to-tfvars converter emits the
manifest-shaped (camelCase, null-pruned) spec with `StringValueOrRef` foreign keys
(`namespace`, `host`) resolved to literal strings before Terraform runs, so the module
performs no snake-to-camel mapping, null-pruning, or `oneOf` logic. `locals.tf` only
strips the Planton `namespace` key (which maps to `metadata.namespace`, not into the CR
spec) and renders the identity labels.

## Outputs

| Output | Description |
|--------|-------------|
| `destination_rule_name` | Name of the created DestinationRule (equals `metadata.name`) |
| `namespace` | Namespace the DestinationRule was created in |

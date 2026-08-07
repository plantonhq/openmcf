# KubernetesServiceEntry Terraform Module

Creates a namespaced Istio `ServiceEntry` via the `kubectl_manifest` resource. The
Istio CRDs must already be installed on the target cluster (see
`KubernetesIstioBaseCrds`), a running istiod is required to program the registry (see
`KubernetesIstio`), and the target namespace must exist (see `KubernetesNamespace`).

## Usage

```bash
planton tofu apply --manifest serviceentry.yaml
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
manifest-shaped (camelCase, null-pruned) spec with the `namespace` `StringValueOrRef`
foreign key resolved to a literal before Terraform runs, so unset fields are omitted
from the manifest and upstream defaults flow through. `locals.tf` only strips the
Planton `namespace` key (which maps to `metadata.namespace`) and renders the identity
labels.

## Outputs

| Output | Description |
|--------|-------------|
| `service_entry_name` | Name of the created ServiceEntry (equals `metadata.name`) |
| `namespace` | Namespace the ServiceEntry was created in |

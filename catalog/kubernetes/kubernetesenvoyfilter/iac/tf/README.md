# KubernetesEnvoyFilter Terraform Module

Creates a namespaced Istio `EnvoyFilter` via the `kubectl_manifest` resource (emitting
`networking.istio.io/v1alpha3` -- EnvoyFilter has not graduated to v1). The Istio CRDs must
already be installed on the target cluster (see `KubernetesIstioBaseCrds`), a running istiod is
required to translate the patches (see `KubernetesIstio`), and the target namespace must exist
(see `KubernetesNamespace`).

## Usage

```bash
planton tofu apply --manifest envoyfilter.yaml
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
passed through verbatim: the platform resolves the `StringValueOrRef` foreign keys
(`namespace`, `target_refs[].name`) to literal strings before Terraform runs, and emits the
manifest-shaped (camelCase, null-pruned) spec, so unset fields are omitted from the manifest
and upstream defaults flow through. `workload_selector` and `target_refs` are mutually
exclusive. The free-form `config_patches[].patch.value` passes through unmodified (the
upstream CRD marks it preserveUnknownFields).

## Outputs

| Output | Description |
|--------|-------------|
| `envoy_filter_name` | Name of the created EnvoyFilter (equals `metadata.name`) |
| `namespace` | Namespace the EnvoyFilter was created in |

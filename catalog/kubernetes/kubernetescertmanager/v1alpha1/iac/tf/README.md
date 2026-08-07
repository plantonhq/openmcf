# KubernetesCertManager Terraform Module

Installs cert-manager from the official Helm chart (`cert-manager` at
`https://charts.jetstack.io`) as a real Helm release. The typed spec renders
into chart values in `locals.tf`; the `helm_values` escape hatch is passed
as a second values document the provider merges last (Helm `-f` semantics) —
the exact semantic twin of the Pulumi module.

## Usage

```bash
planton tofu apply --manifest cert-manager.yaml
```

## Local Development

```bash
terraform init
terraform validate
terraform plan -var-file=terraform.tfvars.json
terraform apply -var-file=terraform.tfvars.json
```

## Inputs

See `variables.tf` for the full variable specification (generated from the
spec proto).

## Outputs

| Output | Description |
|--------|-------------|
| `namespace` | Kubernetes namespace cert-manager was installed into |
| `release_name` | Helm release name (always `cert-manager`) |
| `service_account_name` | Controller ServiceAccount — bind cloud-side for keyless DNS-01 |
| `cluster_resource_namespace` | Where ClusterIssuer credential Secrets live |

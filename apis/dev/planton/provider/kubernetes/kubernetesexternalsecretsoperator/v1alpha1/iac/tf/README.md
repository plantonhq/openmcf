# KubernetesExternalSecretsOperator Terraform Module

Installs the External Secrets Operator from the official Helm chart
(`external-secrets` at `https://charts.external-secrets.io`) as a real Helm
release. The typed spec renders into chart values in `locals.tf`; the
`helm_values` escape hatch is passed as a second values document the
provider merges last (Helm `-f` semantics) — the exact semantic twin of the
Pulumi module.

The release is always named `external-secrets` (one installation per
cluster is an upstream architectural constraint), installs the CRDs by
default, and keeps them on uninstall via the `helm.sh/resource-policy: keep`
annotation unless `crds.keep_on_uninstall` is explicitly false. The install
waits for all three Deployments (controller, webhook, cert-controller) to
become Available.

## Usage

```bash
planton tofu apply --manifest external-secrets-operator.yaml
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
| `namespace` | Kubernetes namespace the operator was installed into |
| `release_name` | Helm release name (always `external-secrets`) |
| `controller_service_account` | Controller ServiceAccount — bind cloud-side for ambient keyless store access |

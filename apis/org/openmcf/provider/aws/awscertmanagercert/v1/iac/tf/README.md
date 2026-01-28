# Terraform Module to Deploy AwsCertManagerCert

This module provisions an AWS Certificate Manager (ACM) certificate with DNS validation,
creates Route53 DNS records for validation, and completes certificate validation.

Generated `variables.tf` reflects the proto schema for `AwsCertManagerCert`.

## Usage

Use the OpenMCF CLI (tofu) with the default local backend:

```shell
openmcf tofu init --manifest hack/manifest.yaml
openmcf tofu plan --manifest hack/manifest.yaml
openmcf tofu apply --manifest hack/manifest.yaml --auto-approve
openmcf tofu destroy --manifest hack/manifest.yaml --auto-approve
```

**Note**: Credentials are provided via stack input (CLI), not in the manifest `spec`.

For more examples, see [`examples.md`](./examples.md) and [`hack/manifest.yaml`](../hack/manifest.yaml).



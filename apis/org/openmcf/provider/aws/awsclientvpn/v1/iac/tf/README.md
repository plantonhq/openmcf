# Terraform Module to Deploy AwsClientVpn

This module provisions an AWS Client VPN endpoint for secure remote access into a VPC using OpenVPN.
It sets up the endpoint, target subnet associations, authorization rules, and optional connection logging.

Generated `variables.tf` reflects the proto schema for `AwsClientVpn`.

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



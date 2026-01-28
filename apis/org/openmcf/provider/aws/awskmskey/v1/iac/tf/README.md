# Terraform Module to Deploy AwsKmsKey

This module provisions an AWS KMS (Key Management Service) key with support for symmetric and asymmetric encryption, key rotation, and optional aliases.
It includes configurable key types, deletion windows, and comprehensive encryption management capabilities.

Generated `variables.tf` reflects the proto schema for `AwsKmsKey`.

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


# Terraform Module to Deploy AwsIamRole

This module provisions an AWS IAM role with support for trust policies, managed policies, and inline policies.
It includes configurable role paths, descriptions, and comprehensive policy management capabilities.

Generated `variables.tf` reflects the proto schema for `AwsIamRole`.

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


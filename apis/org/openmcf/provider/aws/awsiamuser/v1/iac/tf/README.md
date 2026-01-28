# Terraform Module to Deploy AwsIamUser

This module provisions an AWS IAM user with support for managed policies, inline policies, and optional access key creation.
It includes configurable user names, policy attachments, and secure credential management for CI/CD and application use cases.

Generated `variables.tf` reflects the proto schema for `AwsIamUser`.

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


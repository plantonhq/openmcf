# Terraform Module to Deploy AwsS3Bucket

This module provisions an AWS S3 (Simple Storage Service) bucket with support for public and private access configurations.
It includes configurable access controls, ownership settings, and comprehensive object storage management capabilities.

Generated `variables.tf` reflects the proto schema for `AwsS3Bucket`.

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

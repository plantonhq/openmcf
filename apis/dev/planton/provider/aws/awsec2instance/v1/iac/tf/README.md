# Terraform Module to Deploy AwsEc2Instance

This module provisions exactly one AWS EC2 virtual machine: launch source
(inline AMI/type or a referenced launch template with per-instance
overrides), networking, storage, IMDSv2 posture, purchase options, and
lifecycle protections. Everything the instance composes with -- subnet,
security groups, IAM instance profile (by NAME), launch template, KMS
keys -- attaches by reference; the module creates no side resources.

Generated `variables.tf` reflects the proto schema for `AwsEc2Instance`
(generator-owned; never hand-edit).

## Usage

Use the Planton CLI (tofu) with the default local backend:

```shell
planton tofu init --manifest hack/manifest.yaml
planton tofu plan --manifest hack/manifest.yaml
planton tofu apply --manifest hack/manifest.yaml --auto-approve
planton tofu destroy --manifest hack/manifest.yaml --auto-approve
```

**Note**: Credentials are provided via stack input (CLI), not in the manifest `spec`.

For more examples, see [`examples.md`](./examples.md) and [`hack/manifest.yaml`](../hack/manifest.yaml).



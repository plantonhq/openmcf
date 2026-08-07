# Terraform Module to Deploy AliCloudSecurityGroup

This module provisions an Alibaba Cloud Security Group with bundled security
rules. Each entry in `rules` creates an `alicloud_security_group_rule` resource
linked to the security group via `for_each`.

Generated `variables.tf` reflects the proto schema for `AliCloudSecurityGroup`.

## Usage

Use the Planton CLI (tofu) with the default local backend:

```shell
planton tofu init --manifest e2e/manifest.yaml
planton tofu plan --manifest e2e/manifest.yaml
planton tofu apply --manifest e2e/manifest.yaml --auto-approve
planton tofu destroy --manifest e2e/manifest.yaml --auto-approve
```

**Note**: Credentials are provided via stack input (CLI), not in the manifest `spec`.

For more examples, see `../../e2e/manifest.yaml` and [`e2e/manifest.yaml`](../../e2e/manifest.yaml).

# Terraform Module to Deploy AliCloudKubernetesCluster

This module provisions an Alibaba Cloud ACK Managed Kubernetes cluster with configurable networking (Flannel or Terway CNI), addons, control plane logging, maintenance windows, and automatic version upgrades.

Generated `variables.tf` reflects the proto schema for `AliCloudKubernetesCluster`.

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

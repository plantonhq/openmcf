# Terraform Module to Deploy AwsNeptuneCluster

This module provisions an Amazon Neptune graph database cluster and its instances aligned with the Planton API.

## CLI (local backend)

```shell
planton tofu init --manifest ../hack/manifest.yaml
planton tofu plan --manifest ../hack/manifest.yaml
planton tofu apply --manifest ../hack/manifest.yaml --auto-approve
planton tofu destroy --manifest ../hack/manifest.yaml --auto-approve
```

Credentials are passed via the stack input through the CLI, not in `spec`.

## Files
- `variables.tf` (generated; do not edit)
- `provider.tf` — provider setup
- `locals.tf` — computed locals and flags
- `subnet_group.tf` — Neptune subnet group when subnet IDs provided
- `cluster_param_group.tf` — managed cluster parameter group for inline parameters
- `neptune_cluster.tf` — main cluster resource
- `cluster_instances.tf` — per-name folded instances
- `outputs.tf` — outputs matching `AwsNeptuneClusterStackOutputs`

## Presets
See `../../presets/` for ready-to-adapt manifests.

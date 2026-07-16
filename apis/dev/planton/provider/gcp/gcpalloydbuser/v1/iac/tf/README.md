# GcpAlloydbUser - Terraform Module

Terraform implementation of Planton `GcpAlloydbUser`. Enables `alloydb.googleapis.com` then creates `google_alloydb_user`.

## Usage with Planton CLI

```shell
planton tofu init --manifest ../hack/manifest.yaml
planton tofu apply --manifest ../hack/manifest.yaml --auto-approve
```

## Outputs

| Name | Description |
|------|-------------|
| `name` | Fully qualified user resource path |
| `user_id` | User ID |
| `cluster_id` | Cluster resource path |

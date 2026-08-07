# GcpAlloydbInstance - Terraform Module

Terraform implementation of Planton `GcpAlloydbInstance` with feature parity to the Pulumi module. Enables `alloydb.googleapis.com` then creates `google_alloydb_instance`.

## Usage with Planton CLI

```shell
planton tofu init --manifest ../../e2e/manifest.yaml
planton tofu plan --manifest ../../e2e/manifest.yaml
planton tofu apply --manifest ../../e2e/manifest.yaml --auto-approve
```

Manifest: `../../e2e/manifest.yaml`.

## Outputs

| Name | Description |
|------|-------------|
| `instance_name` | Fully qualified instance resource path |
| `ip_address` | Private IP |
| `state` | Instance state |

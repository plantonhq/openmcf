# GcpAlloydbInstance - Terraform Module

Terraform implementation of Planton `GcpAlloydbInstance` with feature parity to the Pulumi module. Enables `alloydb.googleapis.com` then creates `google_alloydb_instance`.

## Usage with Planton CLI

```shell
planton tofu init --manifest ../hack/manifest.yaml
planton tofu plan --manifest ../hack/manifest.yaml
planton tofu apply --manifest ../hack/manifest.yaml --auto-approve
```

Manifest: `../hack/manifest.yaml`.

## Outputs

| Name | Description |
|------|-------------|
| `instance_name` | Fully qualified instance resource path |
| `ip_address` | Private IP |
| `state` | Instance state |

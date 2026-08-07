# GcpTargetHttpProxy - Terraform Module

This Terraform module provisions a GCP Compute Engine global target HTTP proxy. It is the Terraform-side implementation of the Planton `GcpTargetHttpProxy` resource kind and has feature parity with the Pulumi module.

## Overview

The module creates `google_compute_target_http_proxy` — the plaintext-HTTP frontend adapter that binds a global forwarding rule (the VIP) to a URL map (the routing brain). The standard production role is serving the http→https redirect on port 80.

`url_map` is the only mutable field (GCP swaps it in place via `setUrlMap`); everything else is ForceNew.

## Usage with Planton CLI

```shell
planton tofu init --manifest ../../e2e/manifest.yaml
planton tofu plan --manifest ../../e2e/manifest.yaml
planton tofu apply --manifest ../../e2e/manifest.yaml --auto-approve
planton tofu destroy --manifest ../../e2e/manifest.yaml --auto-approve
```

Credentials are provided via stack input (by the CLI), not in the manifest `spec`. Manifest file: `../../e2e/manifest.yaml`.

## Direct Terraform Usage

```bash
cd catalog/gcp/gcptargethttpproxy/iac/tf
terraform init
terraform plan -var-file=terraform.tfvars.json
terraform apply -var-file=terraform.tfvars.json
```

## Variables

| Name | Description | Default |
|------|-------------|---------|
| `metadata` | Resource metadata (name, labels, etc.) | — |
| `spec` | GcpTargetHttpProxy spec | — |

The `spec` object includes: the required `url_map` (plain string after ref resolution), optional `proxy_name` / `description`, the `EXTERNAL_MANAGED`-only `http_keep_alive_timeout_sec`, and the Traffic Director `proxy_bind`.

## Outputs

| Name | Description |
|------|-------------|
| `self_link` | Self-link URI — the value a global forwarding rule references |
| `proxy_name` | Name of the proxy in GCP |
| `proxy_id` | Server-assigned numeric ID |
| `fingerprint` | Server-computed fingerprint for concurrency control |

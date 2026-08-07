# GcpTargetHttpsProxy - Terraform Module

This Terraform module provisions a GCP Compute Engine global target HTTPS proxy. It is the Terraform-side implementation of the Planton `GcpTargetHttpsProxy` resource kind and has feature parity with the Pulumi module.

## Overview

The module creates `google_compute_target_https_proxy` — the TLS-termination node binding a global forwarding rule (the VIP) to a URL map (the routing brain). Certificates attach through exactly one of three mechanisms (classic list, Certificate Manager list, or SNI-scale certificate map); SSL policy, QUIC, and TLS early data live here too.

`url_map`, the certificate wiring, `ssl_policy`, `server_tls_policy`, and `quic_override` update in place; name, description, keep-alive, `tls_early_data`, and `proxy_bind` are ForceNew.

## Usage with Planton CLI

```shell
planton tofu init --manifest ../hack/manifest.yaml
planton tofu plan --manifest ../hack/manifest.yaml
planton tofu apply --manifest ../hack/manifest.yaml --auto-approve
planton tofu destroy --manifest ../hack/manifest.yaml --auto-approve
```

Credentials are provided via stack input (by the CLI), not in the manifest `spec`. Manifest file: `../hack/manifest.yaml`.

## Direct Terraform Usage

```bash
cd catalog/gcp/gcptargethttpsproxy/v1alpha1/iac/tf
terraform init
terraform plan -var-file=terraform.tfvars.json
terraform apply -var-file=terraform.tfvars.json
```

## Variables

| Name | Description | Default |
|------|-------------|---------|
| `metadata` | Resource metadata (name, labels, etc.) | — |
| `spec` | GcpTargetHttpsProxy spec | — |

The `spec` object includes: the required `url_map`, one certificate mechanism (`ssl_certificates` / `certificate_manager_certificates` / `certificate_map`), TLS behavior (`ssl_policy`, `server_tls_policy`, `quic_override`, `tls_early_data`), and the frontend dials (`http_keep_alive_timeout_sec`, `proxy_bind`). Ref fields arrive as plain strings after CLI-side resolution.

## Outputs

| Name | Description |
|------|-------------|
| `self_link` | Self-link URI — the value a global forwarding rule references |
| `proxy_name` | Name of the proxy in GCP |
| `proxy_id` | Server-assigned numeric ID |
| `fingerprint` | Server-computed fingerprint for concurrency control |

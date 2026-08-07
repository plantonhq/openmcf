# GcpGlobalForwardingRule - Terraform Module

This Terraform module provisions a GCP Compute Engine global forwarding rule. It is the Terraform-side implementation of the Planton `GcpGlobalForwardingRule` resource kind and has feature parity with the Pulumi module.

## Overview

The module creates `google_compute_global_forwarding_rule` — the VIP node binding an IP address and port to a target proxy, and the Private Service Connect entry point (scheme `NONE`, sent to the API as an empty scheme).

`target` and `labels` update in place (`setTarget` is the zero-downtime frontend swap); every other field is ForceNew — which is why production frontends bind a reserved `GcpGlobalAddress` rather than an ephemeral IP.

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
cd catalog/gcp/gcpglobalforwardingrule/v1alpha1/iac/tf
terraform init
terraform plan -var-file=terraform.tfvars.json
terraform apply -var-file=terraform.tfvars.json
```

## Variables

| Name | Description | Default |
|------|-------------|---------|
| `metadata` | Resource metadata (name, labels, etc.) | — |
| `spec` | GcpGlobalForwardingRule spec | — |

The `spec` object includes: the required `target` (plain string after ref resolution), the VIP (`ip_address`, `ip_protocol`, `ip_version`, `port_range`), the `load_balancing_scheme` (with the `NONE` PSC sentinel), network wiring for internal/PSC schemes, Traffic Director `metadata_filters`, PSC extras (`service_directory_registration`, `no_automate_dns_zone`), `labels`, and the backend-bucket migration canary dials.

## Outputs

| Name | Description |
|------|-------------|
| `ip_address` | The VIP — the value DNS records point at |
| `self_link` | Self-link URI of the rule |
| `forwarding_rule_name` | Name of the rule in GCP |
| `forwarding_rule_id` | Server-assigned numeric ID |
| `psc_connection_id` | PSC connection id (PSC frontends only) |
| `psc_connection_status` | PSC connection status (PSC frontends only) |

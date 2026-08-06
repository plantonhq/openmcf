# GcpWorkloadIdentityPool - Terraform Module

This Terraform module provisions a GCP Workload Identity Pool (`google_iam_workload_identity_pool`). It is the Terraform-side implementation of the Planton `GcpWorkloadIdentityPool` resource kind and has feature parity with the Pulumi module.

## Overview

The module creates the trust boundary for keyless authentication: external identities (GitHub Actions, AWS workloads, enterprise IdPs) federate through the pool instead of holding service-account keys. `workload_identity_pool_id`, `project`, and `mode` are immutable (ForceNew; the API rejects mode updates even when a plan shows one); `display_name`, `description`, `disabled`, and the inline certificate/trust configs update in place. GCP soft-deletes pools for ~30 days, during which the pool ID cannot be reused — and a create against a soft-deleted ID fails (no undelete-on-create), so prefer `disabled = true` for temporary shutoffs.

The pool resource runs on the `google-beta` provider: on the 6.x line, `mode` and the inline certificate-issuance/trust-config blocks exist only there. The beta provider is a strict superset, so GA-only configurations behave identically.

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
cd apis/dev/planton/provider/gcp/gcpworkloadidentitypool/v1alpha1/iac/tf
terraform init
terraform plan -var-file=terraform.tfvars.json
terraform apply -var-file=terraform.tfvars.json
```

## Variables

| Name | Description | Default |
|------|-------------|---------|
| `metadata` | Resource metadata (name, labels, etc.) | — |
| `spec` | GcpWorkloadIdentityPool spec | — |

The `spec` object includes: `workload_identity_pool_id` (4-32 chars; lowercase, digits, hyphens; `gcp-` reserved), `project_id` (plain string — the CLI resolves references before the module runs; empty falls back to the provider default project), optional `display_name`, `description`, `disabled` (default `false`), `mode` (default `FEDERATION_ONLY`), `inline_certificate_issuance_config`, and `inline_trust_config`.

## Outputs

| Name | Description |
|------|-------------|
| `name` | Full pool resource name — the handle IAM principals and providers are built from |
| `workload_identity_pool_id` | The bare pool ID (providers reference this) |
| `state` | `ACTIVE`, or `DELETED` while soft-deleted |

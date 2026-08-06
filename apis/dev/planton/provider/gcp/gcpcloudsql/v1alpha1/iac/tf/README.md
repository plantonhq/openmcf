# GcpCloudSql - Terraform Module

This Terraform module provisions a Cloud SQL instance (`google_sql_database_instance`). It is the Terraform-side implementation of the Planton `GcpCloudSql` resource kind and has feature parity with the Pulumi module.

## Overview

The module creates a managed MySQL, PostgreSQL, or SQL Server instance — a primary, or (with `master_instance_name`) a read replica — with edition/HA, disk tuning, explicit connectivity (public IPv4 / private VPC IP / Private Service Connect), backups with point-in-time recovery, maintenance scheduling, Query Insights, password policies, managed connection pooling, CMEK, and both delete guards. It also enables the Cloud SQL Admin API (`sqladmin.googleapis.com`) so a fresh project works on the first deploy.

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
cd apis/dev/planton/provider/gcp/gcpcloudsql/v1alpha1/iac/tf
terraform init
terraform plan -var-file=terraform.tfvars.json
terraform apply -var-file=terraform.tfvars.json
```

## Variables

| Name | Description | Default |
|------|-------------|---------|
| `spec` | GcpCloudSql spec (instance_name, region, engine/version/tier, disk, network, backup, …) | — |
| `metadata` | Resource metadata (name, org, env, labels) | — |

The `spec` object mirrors the proto spec; `StringValueOrRef` fields (`project_id`, `network.private_network`, `encryption_key_name`, `master_instance_name`) arrive as plain strings after the CLI's ref resolution.

## Outputs

| Name | Description |
|------|-------------|
| `instance_name` | The composition key databases, users, and replicas reference |
| `connection_name` | `project:region:instance` for the Auth Proxy and connectors |
| `private_ip` | Private IP address (empty unless private_network is configured) |
| `public_ip` | Public IPv4 address (empty unless ipv4_enabled) |
| `self_link` | Self-link URL of the Cloud SQL instance |
| `service_account_email` | The instance's Google-managed service account |
| `dns_name` | DNS name (PSC-enabled instances) |
| `psc_service_attachment_link` | PSC service attachment link (PSC-enabled instances) |

# AwsMskServerlessCluster — Terraform IaC Module

Terraform/OpenTofu module for provisioning an Amazon MSK Serverless cluster using the Planton `AwsMskServerlessClusterSpec`.

## Overview

This module creates a single `aws_msk_serverless_cluster`: network interfaces in the referenced subnets, the referenced security groups attached, and SASL/IAM client authentication enabled unconditionally (the only scheme serverless MSK supports — AWS requires it, so it is not a variable).

The resource is effectively immutable: everything except tags is create-time (ForceNew).

## Inputs

| Variable | Type | Required | Description |
|----------|------|----------|-------------|
| `metadata` | object | yes | Resource ID, name, org, env |
| `spec` | object | yes | `AwsMskServerlessClusterSpec` — see `variables.tf` (generator-owned contract) |

### Spec Fields

| Field | Type | Default | Description |
|---|---|---|---|
| `region` | string | **required** | AWS region |
| `subnet_ids` | list(string) | **required** (≥1) | Subnets for the cluster network interfaces. ForceNew. |
| `security_group_ids` | list(string) | `[]` (VPC default group) | Security groups attached to the network interfaces (max 5). ForceNew. |

## Outputs

| Output | Description |
|--------|-------------|
| `cluster_arn` | Cluster ARN (also the resource identifier) |
| `cluster_name` | Cluster name |
| `cluster_uuid` | UUID extracted from the ARN |
| `bootstrap_brokers_sasl_iam` | SASL/IAM broker endpoints (port 9098) |

## File Structure

| File | Purpose |
|------|---------|
| `provider.tf` | AWS provider configuration (hashicorp/aws >= 6.0.0) |
| `backend.tf` | Local state backend (overridden by the runtime) |
| `variables.tf` | Input variable definitions (generator-owned contract) |
| `locals.tf` | Cluster name basis and resource-identity tags |
| `cluster.tf` | The MSK Serverless cluster resource |
| `outputs.tf` | Output definitions |

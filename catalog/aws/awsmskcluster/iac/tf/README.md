# AwsMskCluster — Terraform IaC Module

Terraform module for provisioning AWS MSK (Managed Streaming for Apache Kafka) clusters using the Planton `AwsMskClusterSpec`.

## Overview

This module creates:
- An MSK Cluster (`aws_msk_cluster`) with configurable brokers, encryption, authentication, connectivity, logging, and monitoring. The referenced `security_group_ids` (required, ≥1) attach directly to the broker network interfaces; ingress rules live on those first-class security-group nodes, never on a module-managed shadow group.
- An inline MSK Configuration (`aws_msk_configuration`) from `server_properties` — conditional on the map being non-empty.
- SCRAM secret associations (`aws_msk_single_scram_secret_association`, one per ARN in `scram_secret_arns`), a cluster policy (`aws_msk_cluster_policy` from `cluster_policy`, serialized to JSON), and declared Kafka topics (`aws_msk_topic` per `topics` entry, keyed by name, exported as the `topic_arns` map) — folded satellites in `satellites.tf`.

## Usage

```hcl
module "msk" {
  source = "./path/to/this/module"

  provider_config = {
    region = "us-east-1"
  }

  metadata = {
    id   = "prod-events"
    name = "prod-events"
    org  = "myorg"
    env  = "production"
  }

  spec = {
    kafka_version          = "3.6.0"
    number_of_broker_nodes = 3
    instance_type          = "kafka.m5.large"
    subnet_ids             = ["subnet-aaa", "subnet-bbb", "subnet-ccc"]
    security_group_ids     = ["sg-0abc1234"]

    authentication = {
      sasl_iam_enabled = true
    }

    server_properties = {
      "auto.create.topics.enable"  = "false"
      "default.replication.factor" = "3"
      "min.insync.replicas"        = "2"
    }
  }
}
```

## Inputs

| Variable | Type | Required | Description |
|----------|------|----------|-------------|
| `provider_config` | object | yes | AWS region and optional credentials |
| `metadata` | object | yes | Resource ID, name, org, env |
| `spec` | object | yes | `AwsMskClusterSpec` — see `variables.tf` for full type |

See `variables.tf` for the complete type definition of `spec`, including all optional fields and their defaults.

## Outputs

| Output | Description |
|--------|-------------|
| `cluster_arn` | ARN of the MSK cluster |
| `cluster_name` | Cluster name |
| `cluster_uuid` | UUID extracted from ARN |
| `current_version` | Cluster version (for updates) |
| `bootstrap_brokers` | Plaintext broker endpoints (port 9092) |
| `bootstrap_brokers_tls` | TLS broker endpoints (port 9094) |
| `bootstrap_brokers_sasl_iam` | SASL/IAM broker endpoints (port 9098) |
| `bootstrap_brokers_sasl_scram` | SASL/SCRAM broker endpoints (port 9096) |
| `bootstrap_brokers_public_tls` | Public TLS broker endpoints |
| `bootstrap_brokers_public_sasl_iam` | Public SASL/IAM broker endpoints |
| `bootstrap_brokers_public_sasl_scram` | Public SASL/SCRAM broker endpoints |
| `bootstrap_brokers_vpc_connectivity_tls` | PrivateLink mTLS broker endpoints |
| `bootstrap_brokers_vpc_connectivity_sasl_iam` | PrivateLink SASL/IAM broker endpoints |
| `bootstrap_brokers_vpc_connectivity_sasl_scram` | PrivateLink SASL/SCRAM broker endpoints |
| `zookeeper_connect_string` | ZooKeeper plaintext endpoints (empty on KRaft-mode clusters) |
| `zookeeper_connect_string_tls` | ZooKeeper TLS endpoints (empty on KRaft-mode clusters) |
| `configuration_arn` | Inline config ARN (empty if not created) |

## File Structure

| File | Purpose |
|------|---------|
| `provider.tf` | AWS provider configuration (hashicorp/aws >= 6.41) |
| `variables.tf` | Generator-owned input variable definitions (do not edit by hand) |
| `locals.tf` | Tags, configuration condition, server_properties serialization |
| `cluster.tf` | MSK Configuration + MSK Cluster resources |
| `satellites.tf` | SCRAM secret associations + cluster policy + declared topics |
| `outputs.tf` | 17 output definitions |

## Prerequisites

- Terraform 1.5+ / OpenTofu
- AWS provider >= 6.41
- AWS credentials (via provider config or ambient)

## Related

- [Spec reference](../../README.md)

# Basic IAM-Authenticated Serverless Kafka

A minimal MSK Serverless cluster: network interfaces in two private subnets, a referenced security group controlling access, and SASL/IAM authentication (always on — the only scheme serverless MSK supports).

## What this preset gives you

- Fully managed Kafka with automatic capacity scaling and pay-per-throughput billing — no brokers, storage, or version to manage.
- IAM-native client authentication on port 9098: no passwords or certificates to rotate.
- Network access controlled by the referenced security group, where the port-9098 ingress rule lives.

## Before you deploy

- Replace the subnet placeholders with private subnets in two different Availability Zones.
- Replace the security group placeholder with a group whose ingress opens port 9098 to your client workloads (for example, from the application tier's security group).

## Remix ideas

- Reference Planton-managed resources with `valueFrom` instead of literal IDs so the cluster composes into the resource graph.
- Omit `securityGroupIds` for a quick sandbox — AWS attaches the VPC's default security group.
- Grant client workloads `kafka-cluster:*` IAM permissions scoped to the exported `cluster_arn`.

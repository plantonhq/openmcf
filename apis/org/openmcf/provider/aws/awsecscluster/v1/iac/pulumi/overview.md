# Overview

The **AwsEcsCluster** API resource simplifies the process of creating and operating Amazon ECS clusters on AWS. It
leverages OpenMCF’s unified API model to provide a consistent, declarative way of defining cluster
configurations—from basic setups using Fargate to more advanced deployments with multiple capacity providers and ECS
Exec enabled.

By encapsulating all necessary settings into a single YAML specification, **AwsEcsCluster** ensures that developers and
platform teams can consistently deploy production-grade ECS clusters without juggling multiple AWS APIs or tooling
differences. As part of the broader OpenMCF framework, it supports both Pulumi and Terraform under the hood,
offering flexibility in your Infrastructure as Code approach.

## Key Features

- **Fargate & EC2 Capacity**  
  Provision ECS clusters with pure Fargate capacity or combine Fargate and EC2 providers (including Spot options) for
  cost optimization and operational flexibility.

- **Built-in Observability**  
  Easily enable CloudWatch Container Insights to monitor resource utilization and application health, streamlining
  performance troubleshooting.

- **ECS Exec Support**  
  Optionally allow secure shell access (`exec`) into running containers for real-time debugging and maintenance tasks.

- **Consistent Resource Model**  
  Adheres to the Kubernetes-like OpenMCF resource structure (`apiVersion`, `kind`, `metadata`, `spec`, `status`),
  providing a familiar, validated schema.

- **Pulumi & Terraform Integration**  
  Seamlessly switch between Pulumi and Terraform provisioning workflows using the same manifest, maintaining consistency
  across your multi-cloud deployments.

## Next Steps

- See the [README.md](./README.md) for detailed usage instructions and configuration settings.
- Refer to the [examples.md](./examples.md) to explore different manifest samples and common deployment scenarios.
- Consult the official OpenMCF documentation for broader insights into multi-cloud deployment patterns and CLI
  workflows.

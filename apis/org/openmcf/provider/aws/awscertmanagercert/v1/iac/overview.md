# Overview

The **AWS Cert Manager Cert** component in OpenMCF extends our multi-cloud framework to offer a straightforward way of
deploying containerized workloads on Amazon Elastic Container Service (ECS). It manages essential configurations—like
cluster, container images, networking, and IAM roles—within a single resource spec, simplifying both Fargate and
EC2-based deployments under the same approach used across all OpenMCF components.

## Purpose and Functionality

- **Unified ECS Deployment**: Create and manage ECS services using a standard OpenMCF manifest. This streamlines
  container service deployments without juggling multiple AWS configurations or manual setup steps.
- **Fargate or EC2**: Choose the best ECS launch type for your workload, whether you need a serverless model with
  Fargate or the flexibility of EC2 capacity.
- **Essential Resource Settings**: Configure CPU, memory, subnets, security groups, and environment variables in one
  consistent schema, validated against Protobuf definitions.
- **Consistent Multi-Cloud Model**: Stay within the OpenMCF YAML-based workflow and CLI you already use for other
  resources, eliminating the need to learn specialized ECS tools.

## Key Benefits

- **One Workflow for All**: Leverage the same CLI-driven process—validated through Protobuf—for AWS ECS services as you
  do for other cloud deployments. No additional per-provider complexities or scripts.
- **Streamlined Provisioning**: Automatically manage tasks like attaching security groups and subnets. Whether you use
  Pulumi or Terraform under the hood, OpenMCF handles resource provisioning seamlessly.
- **Scalability & Observability**: Scale services by adjusting the desired task count, or enable advanced ECS features
  without rewriting your entire deployment process.
- **Faster Iterations**: Rapidly test changes to images, environment variables, or service parameters in a single place,
  reducing the operational overhead of updating ECS configurations manually.

Use the **AwsCertManagerCert** component to reduce friction when deploying container-based workloads on AWS. By integrating with
the standard OpenMCF CLI and resource model, teams can accelerate delivery of ECS applications while preserving a
unified, multi-cloud operational approach.

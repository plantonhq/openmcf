---
title: "Presets"
description: "Ready-to-deploy configuration presets for MSK Serverless Cluster"
type: "preset-list"
componentSlug: "msk-serverless-cluster"
componentTitle: "MSK Serverless Cluster"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-basic-iam"
    rank: "01"
    title: "Basic IAM-Authenticated Serverless Kafka"
    excerpt: "A minimal MSK Serverless cluster: network interfaces in two private subnets, a referenced security group controlling access, and SASL/IAM authentication (always on — the only scheme serverless MSK..."
  - slug: "02-composed-references"
    rank: "02"
    title: "Composed Serverless Kafka (Full References)"
    excerpt: "An MSK Serverless cluster where every cross-resource input resolves from Planton-managed resources: subnets from `AwsSubnet` nodes and the attached security group from an `AwsSecurityGroup` node. The..."
---

# MSK Serverless Cluster Presets

Ready-to-deploy configuration presets for MSK Serverless Cluster. Each preset is a complete manifest you can copy, customize, and deploy.

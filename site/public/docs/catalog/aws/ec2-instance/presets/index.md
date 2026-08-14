---
title: "Presets"
description: "Ready-to-deploy configuration presets for EC2 Instance"
type: "preset-list"
componentSlug: "ec2-instance"
componentTitle: "EC2 Instance"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-ssm-managed"
    rank: "01"
    title: "SSM-Managed Hardened Instance"
    excerpt: "This preset creates an EC2 instance accessible via AWS Systems Manager Session Manager, hardened to the modern baseline: IMDSv2 enforced, an encrypted gp3 root volume, and termination protection. SSM..."
  - slug: "02-launch-template"
    rank: "02"
    title: "Launch-Template Instance"
    excerpt: "This preset launches an instance from an `AwsLaunchTemplate` -- the org's golden baseline of AMI, hardening (IMDSv2, encrypted volumes), and storage defaults -- and overrides only what makes this..."
  - slug: "03-spot-worker"
    rank: "03"
    title: "Spot Worker Instance"
    excerpt: "This preset runs a standalone instance on Spot capacity -- typically 60-90% cheaper than On-Demand -- configured as a persistent request that stops (rather than terminates) on interruption and..."
  - slug: "04-static-eni-identity"
    rank: "04"
    title: "Static Network Identity (Pre-Provisioned ENI)"
    excerpt: "This preset attaches a pre-provisioned network interface (ENI) as the instance's primary interface instead of creating one. The ENI -- not the instance -- owns the network identity: the subnet, the..."
  - slug: "05-capacity-block-ml"
    rank: "05"
    title: "Capacity Block ML Training Node"
    excerpt: "This preset launches a GPU instance into a pre-purchased EC2 Capacity Block -- reserved GPU capacity for a defined time window at a committed price, the way AWS sells scarce accelerators (P5/P4,..."
---

# EC2 Instance Presets

Ready-to-deploy configuration presets for EC2 Instance. Each preset is a complete manifest you can copy, customize, and deploy.

---
title: "Presets"
description: "Ready-to-deploy configuration presets for NLB"
type: "preset-list"
componentSlug: "nlb"
componentTitle: "NLB"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-internal"
    rank: "01"
    title: "Internal NLB"
    excerpt: "This preset creates the simplest useful Network Load Balancer: an internal NLB with one node in one private subnet. It carries no listeners by design — attach `AwsLbListener` resources..."
  - slug: "02-static-ip-internet-facing"
    rank: "02"
    title: "Static-IP Internet-Facing NLB"
    excerpt: "This preset creates the headline NLB use case: an internet-facing load balancer whose public IPs never change. Each subnet mapping pins one NLB node to an Elastic IP, referenced from an..."
  - slug: "03-private-link-hardened"
    rank: "03"
    title: "PrivateLink-Hardened NLB"
    excerpt: "This preset creates the NLB posture for exposing a service over AWS PrivateLink: an internal load balancer with security groups attached, security-group rules enforced on traffic arriving through..."
---

# NLB Presets

Ready-to-deploy configuration presets for NLB. Each preset is a complete manifest you can copy, customize, and deploy.

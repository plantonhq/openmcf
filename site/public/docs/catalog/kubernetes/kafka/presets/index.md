---
title: "Presets"
description: "Ready-to-deploy configuration presets for Kafka"
type: "preset-list"
componentSlug: "kafka"
componentTitle: "Kafka"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-dev-single-node"
    rank: "01"
    title: "Dev single node preset"
    excerpt: "The smallest declarable Kafka that actually serves: one dual-role node (KRaft controller + broker in a single pod), one plaintext internal listener, and single-node replication settings (RF 1..."
  - slug: "02-production-three-broker"
    rank: "02"
    title: "Production three broker preset"
    excerpt: "The standard production posture: a 3-node controller pool (the KRaft quorum tolerates one loss) plus a 3-node broker pool on JBOD storage, RF-3 / min-ISR-2 replication (one broker can die without..."
  - slug: "03-aws-external-loadbalancer"
    rank: "03"
    title: "AWS external loadbalancer preset"
    excerpt: "The production three-broker shape plus a `loadbalancer` listener for clients OUTSIDE the Kubernetes cluster: the AWS Load Balancer Controller provisions one NLB per broker plus one for bootstrap..."
  - slug: "04-scram-app-cluster"
    rank: "04"
    title: "SCRAM app cluster preset"
    excerpt: "A quorum-safe application cluster with exactly one credential story: a single SCRAM-SHA-512-over-TLS listener plus simple ACL authorization. Every client — including the admin — is a..."
---

# Kafka Presets

Ready-to-deploy configuration presets for Kafka. Each preset is a complete manifest you can copy, customize, and deploy.

---
title: "Presets"
description: "Ready-to-deploy configuration presets for Kafka User"
type: "preset-list"
componentSlug: "kafka-user"
componentTitle: "Kafka User"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-producer"
    rank: "01"
    title: "Producer (SCRAM)"
    excerpt: "This preset declares a producer identity: SCRAM-SHA-512 authentication with an operator-generated password, and a prefix ACL granting write access to every topic starting with `orders-`. The user..."
  - slug: "02-consumer"
    rank: "02"
    title: "Consumer (SCRAM)"
    excerpt: "This preset declares a consumer identity: SCRAM-SHA-512 authentication with an operator-generated password, and the two-part ACL grant every consumer needs — Read/Describe on the topic it consumes..."
  - slug: "03-mtls-service"
    rank: "03"
    title: "mTLS Service (Producer + Consumer, Quotas)"
    excerpt: "This preset declares a full service identity on mutual TLS: the user operator issues a client certificate from the cluster's clients CA into the user's Secret, the ACLs cover both directions (produce..."
---

# Kafka User Presets

Ready-to-deploy configuration presets for Kafka User. Each preset is a complete manifest you can copy, customize, and deploy.

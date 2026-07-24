---
title: "Presets"
description: "Ready-to-deploy configuration presets for RabbitMQ"
type: "preset-list"
componentSlug: "rabbitmq"
componentTitle: "RabbitMQ"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-dev-single-node"
    rank: "01"
    title: "Dev single node preset"
    excerpt: "The fastest declarable RabbitMQ: one node on ephemeral storage (emptyDir — no PVC, no StorageClass involved), 1Gi of memory with requests equal to limits, and a 60-second termination grace period in..."
  - slug: "02-production-quorum"
    rank: "02"
    title: "Production quorum preset"
    excerpt: "A production RabbitMQ: 3 nodes (the quorum posture — quorum queues and the Raft-based metadata store survive one node loss), a 50Gi data volume per node, 4Gi of memory with requests equal to limits..."
  - slug: "03-mqtt-iot-broker"
    rank: "03"
    title: "MQTT IoT broker preset"
    excerpt: "A RabbitMQ cluster serving an IoT / device fleet over MQTT 5.0 on the shared broker core: 3 nodes (device fleets reconnect in herds, so availability matters from day one), the rabbitmq_mqtt and..."
---

# RabbitMQ Presets

Ready-to-deploy configuration presets for RabbitMQ. Each preset is a complete manifest you can copy, customize, and deploy.

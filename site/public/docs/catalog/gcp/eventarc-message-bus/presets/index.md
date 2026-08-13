---
title: "Presets"
description: "Ready-to-deploy configuration presets for Eventarc Message Bus"
type: "preset-list"
componentSlug: "eventarc-message-bus"
componentTitle: "Eventarc Message Bus"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-bus-with-topic-pipeline"
    rank: "01"
    title: "Bus with Topic Pipeline"
    excerpt: "The smallest useful hub: everything published to the bus lands on one Pub/Sub topic, where downstream consumers subscribe — Eventarc Advanced as a managed firehose."
  - slug: "02-audit-fan-out"
    rank: "02"
    title: "Audit Fan-Out"
    excerpt: "One hub, two consumers: a Google API source feeds the bus, and CEL enrollments split the stream — storage events to an analytics topic, everything else to an ops workflow."
---

# Eventarc Message Bus Presets

Ready-to-deploy configuration presets for Eventarc Message Bus. Each preset is a complete manifest you can copy, customize, and deploy.

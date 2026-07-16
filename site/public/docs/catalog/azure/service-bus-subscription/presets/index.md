---
title: "Presets"
description: "Ready-to-deploy configuration presets for Service Bus Subscription"
type: "preset-list"
componentSlug: "service-bus-subscription"
componentTitle: "Service Bus Subscription"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-catch-all-consumer"
    rank: "01"
    title: "Catch-All Consumer"
    excerpt: "This preset creates an unfiltered subscription: the consumer receives every message published to the topic (Azure's auto-created `$Default` catch-all rule stays in place). The right starting point..."
  - slug: "02-filtered-consumer"
    rank: "02"
    title: "Filtered Consumer"
    excerpt: "This preset creates a subscription with a SQL filter rule. Declared rules are ADDITIVE alongside Azure's auto-created `$Default` catch-all -- for restrictive delivery (ONLY matches), remove the..."
  - slug: "03-fanout-to-work-queue"
    rank: "03"
    title: "Fan-Out to Work Queue"
    excerpt: "This preset creates a routing subscription: it filters the topic's stream and auto-forwards every match into a work queue that a processing fleet drains. Filtering and processing decouple -- the..."
---

# Service Bus Subscription Presets

Ready-to-deploy configuration presets for Service Bus Subscription. Each preset is a complete manifest you can copy, customize, and deploy.

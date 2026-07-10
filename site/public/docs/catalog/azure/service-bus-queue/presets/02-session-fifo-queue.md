---
title: "Session FIFO Queue"
description: "This preset creates a session-aware queue with duplicate detection: strict per-session ordering with exclusive consumption, and idempotent producers that can retry sends safely. The right shape for..."
type: "preset"
rank: "02"
presetSlug: "02-session-fifo-queue"
componentSlug: "service-bus-queue"
componentTitle: "Service Bus Queue"
provider: "azure"
icon: "package"
order: 2
---

# Session FIFO Queue

This preset creates a session-aware queue with duplicate detection:
strict per-session ordering with exclusive consumption, and idempotent
producers that can retry sends safely. The right shape for
entity-ordered processing (per-customer, per-device, per-conversation).

## When to Use

- Events for one entity must process in order (account transactions,
  device telemetry commands)
- Producers retry on ambiguous failures and duplicates are unacceptable

## Key Configuration Choices

- **`requiresSession: true`** (ForceNew) -- consumers must be
  session-aware; each session is delivered to one consumer at a time
- **`requiresDuplicateDetection: true`** (ForceNew) + **`PT10M` window**
  -- MessageId-based dedup; size the window to the producer's retry
  horizon
- **`lockDuration: PT2M`** -- stateful session processing usually needs
  more than the 1-minute default

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-app-bus` | The AzureServiceBusNamespace's Planton resource name | Your messaging composition |
| `account-transactions` | Unique within the namespace | Your naming convention |

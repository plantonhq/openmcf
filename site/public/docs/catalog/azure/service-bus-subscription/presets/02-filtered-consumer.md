---
title: "Filtered Consumer"
description: "This preset creates a subscription with a SQL filter rule. Declared rules are ADDITIVE alongside Azure's auto-created `$Default` catch-all -- for restrictive delivery (ONLY matches), remove the..."
type: "preset"
rank: "02"
presetSlug: "02-filtered-consumer"
componentSlug: "service-bus-subscription"
componentTitle: "Service Bus Subscription"
provider: "azure"
icon: "package"
order: 2
---

# Filtered Consumer

This preset creates a subscription with a SQL filter rule. Declared
rules are ADDITIVE alongside Azure's auto-created `$Default` catch-all
-- for restrictive delivery (ONLY matches), remove the catch-all once
after creation (the service-created rule cannot be declared or adopted
by the management plane):

```shell
az servicebus topic subscription rule delete --name '$Default' \
  --namespace-name <ns> --topic-name <topic> \
  --subscription-name emea-consumer --resource-group <rg>
```

## When to Use

- A consumer that cares about a slice of the stream (one region, one
  event type, one priority band)
- Cost/noise control -- with the catch-all removed, unmatched messages
  are never delivered, so the consumer never pays to discard them

## Key Configuration Choices

- **The one-time catch-all removal** is the load-bearing step: until
  it runs, the auto-created catch-all still delivers EVERYTHING
  alongside your rule
- **`filterType: SQL_FILTER`** -- SQL-92-like expressions over system
  (`sys.Label`) and user properties; switch to CORRELATION_FILTER for
  cheap equality matching at high throughput

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-event-topic` | The AzureServiceBusTopic's Planton resource name | Your messaging composition |
| `emea-consumer` | ≤50 chars | Your naming convention |
| `region = 'emea' AND priority > 3` | e.g. `region = 'emea' AND priority > 3` | Your message property contract |

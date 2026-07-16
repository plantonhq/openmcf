# Standard Streaming Namespace

This preset creates a STANDARD-tier Event Hubs namespace with elastic
throughput -- the full-featured multi-tenant tier (Kafka endpoint, 20
consumer groups per hub, 7-day retention) that fits most production
streaming workloads.

## When to Use

- The default starting point for any event-streaming domain
- Kafka migrations -- the namespace IS a Kafka-compatible endpoint, no
  cluster to run
- Workloads that don't need reserved capacity or VNet-scale isolation
  (those are PREMIUM/dedicated features)

## Key Configuration Choices

- **`sku: STANDARD`** -- multi-tenant with the full feature set;
  BASIC↔STANDARD changes update in place, but moving to PREMIUM later
  replaces the namespace
- **Auto-inflate with a ceiling of 10 TUs** -- Azure grows throughput
  under load instead of throttling producers; it never scales back
  down, so trim `capacity` manually after a traffic spike
- **`tags`** -- ARM tags are Azure's governance surface; user tags merge
  over the platform's identity tags

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `eastus` | The namespace's region | Your region strategy |
| `my-streaming-rg` | The AzureResourceGroup's Planton resource name | Your foundation composition |
| `myorg-telemetry-hubs` | 6-50 chars, letters/numbers/hyphens, globally unique | Your naming convention |
| `telemetry-ingestion` | What this namespace carries | Your data taxonomy |

## Downstream Wiring

Hubs and the other family kinds reference the namespace's ARM id:

```yaml
# On an AzureEventHub
namespaceId:
  valueFrom:
    kind: AzureEventHubNamespace
    name: my-telemetry-hubs
    fieldPath: status.outputs.namespace_id
```

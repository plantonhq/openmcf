---
title: "Hub-Scoped Producer Credential"
description: "This preset mints a send-only SAS credential on ONE event hub -- the least-privilege connection string a producer service should hold."
type: "preset"
rank: "01"
presetSlug: "01-hub-producer"
componentSlug: "event-hub-authorization-rule"
componentTitle: "Event Hub Authorization Rule"
provider: "azure"
icon: "package"
order: 1
---

# Hub-Scoped Producer Credential

This preset mints a send-only SAS credential on ONE event hub -- the
least-privilege connection string a producer service should hold.

## When to Use

- Each producer service gets its own send-only rule on exactly the
  stream it writes -- a leaked credential cannot read data or touch
  other streams
- SAS-based estates; for a keyless posture skip rules entirely (disable
  the namespace's `localAuthenticationEnabled` and grant Entra
  identities data-plane roles)

## Key Configuration Choices

- **`eventHubId` scope** -- rights end at this hub's boundary
- **`send: true` only** -- producers don't need listen; the spec rejects
  a rule with no rights at all
- **Key rotation** -- the secondary key/connection-string outputs are
  the rotation partner: move clients to the secondary, regenerate the
  primary, move back

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-telemetry-stream` | The AzureEventHub's Planton resource name | Your streaming composition |
| `telemetry-producer` | The rule name (the SharedAccessKeyName clients present) | Your credential taxonomy |

## Downstream Wiring

The connection string surfaces as a sensitive output:

```yaml
# Consumed by an application's configuration
connectionString:
  valueFrom:
    kind: AzureEventHubAuthorizationRule
    name: my-producer-credential
    fieldPath: status.outputs.primary_connection_string
```

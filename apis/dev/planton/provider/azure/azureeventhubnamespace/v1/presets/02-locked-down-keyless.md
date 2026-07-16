# Locked-Down Keyless Namespace

This preset creates a STANDARD namespace in the production security
posture: SAS authentication disabled (Entra-only data plane) and a DENY
firewall admitting only named sources.

## When to Use

- Estates standardizing on keyless (Entra ID) authentication -- no
  connection strings to rotate, leak, or vault
- Compliance regimes requiring network allow-lists on messaging planes

## Key Configuration Choices

- **`localAuthenticationEnabled: false`** -- every SAS rule's keys stop
  working namespace-wide; pair with `AzureRoleAssignment` grants of
  "Azure Event Hubs Data Owner/Sender/Receiver" to workload identities
- **`defaultAction: DENY` + admitted sources** -- Azure rejects a DENY
  rule set with nothing admitted, so the preset admits an IP range and
  a service-endpoint subnet
- **`trustedServiceAccessEnabled: true`** -- platform services (Azure
  Monitor diagnostic streaming, Event Grid delivery) bypass the
  firewall; without it a locked-down namespace silently stops receiving
  platform telemetry

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `myorg-secure-hubs` | Globally unique namespace name | Your naming convention |
| `203.0.113.0/24` | The admitted public range | Your egress IPs |
| `my-app-subnet` | The AzureSubnet carrying the Microsoft.EventHub service endpoint | Your network composition |

## Downstream Wiring

Grant a workload identity the data plane instead of handing it keys:

```yaml
# An AzureRoleAssignment granting send on the namespace
roleDefinitionName: Azure Event Hubs Data Sender
scope:
  valueFrom:
    kind: AzureEventHubNamespace
    name: my-secure-hubs
    fieldPath: status.outputs.namespace_id
```

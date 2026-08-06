# Azure Event Hub Authorization Rule

Creates a SAS credential for Azure Event Hubs -- listen/send/manage rights scoped to exactly one of a namespace or a single event hub. The least-privilege alternative to the namespace's root key.

## What Gets Created

When you deploy an AzureEventHubAuthorizationRule resource, Planton provisions ONE of:

- **Namespace Authorization Rule** -- rights over every hub in the namespace
- **Event Hub Authorization Rule** -- rights over one hub only

The scope you reference picks the resource; the credential surface is identical across both.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **The scope entity** -- an AzureEventHubNamespace or an AzureEventHub

## Quick Start

Create a file `auth-rule.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureEventHubAuthorizationRule
metadata:
  name: telemetry-producer
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureEventHubAuthorizationRule.telemetry-producer
spec:
  ruleName: telemetry-producer
  eventHubId:
    valueFrom:
      kind: AzureEventHub
      name: telemetry-stream
      fieldPath: status.outputs.event_hub_id
  send: true
```

Deploy:

```shell
planton apply -f auth-rule.yaml
```

Azure's rights contract: at least one of listen/send/manage must be true, and manage requires BOTH listen and send. `RootManageSharedAccessKey` is reserved -- the root rule's keys already surface as AzureEventHubNamespace outputs.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `primary_connection_string` | The ready-to-use credential applications hold (hub-scoped rules append `EntityPath={hub}`); rotate via the secondary |
| `primary_key` / `secondary_key` | For SDKs that take key and key name separately or mint their own SAS tokens |
| `authorization_rule_id` | The rule's ARM identity |
| `*_connection_string_alias` | Populated only when the namespace carries a geo-DR pairing -- failover-stable credentials |

## Related Resources

- [Azure Event Hub Namespace](/docs/catalog/azure/azureeventhubnamespace) / [Azure Event Hub](/docs/catalog/azure/azureeventhub) -- the two scopes
- [Azure Role Assignment](/docs/catalog/azure/azureroleassignment) -- the keyless alternative

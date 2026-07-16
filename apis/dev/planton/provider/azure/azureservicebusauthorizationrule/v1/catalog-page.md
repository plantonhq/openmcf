# Azure Service Bus Authorization Rule

Creates a SAS credential for Azure Service Bus -- listen/send/manage rights scoped to exactly one of a namespace, queue, or topic. The least-privilege alternative to the namespace's root key.

## What Gets Created

When you deploy an AzureServiceBusAuthorizationRule resource, Planton provisions ONE of:

- **Namespace Authorization Rule** -- rights over every entity in the namespace
- **Queue Authorization Rule** -- rights over one queue only
- **Topic Authorization Rule** -- rights over one topic only

The scope you reference picks the resource; the credential surface is identical across all three.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **The scope entity** -- an AzureServiceBusNamespace, AzureServiceBusQueue, or AzureServiceBusTopic

## Quick Start

Create a file `auth-rule.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureServiceBusAuthorizationRule
metadata:
  name: orders-sender
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureServiceBusAuthorizationRule.orders-sender
spec:
  ruleName: orders-sender
  queueId:
    valueFrom:
      kind: AzureServiceBusQueue
      name: orders-queue
      fieldPath: status.outputs.queue_id
  send: true
```

Deploy:

```shell
planton apply -f auth-rule.yaml
```

Azure's rights contract: at least one of listen/send/manage must be true, and manage requires BOTH listen and send.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `primary_connection_string` | The ready-to-use credential applications hold (rotate via the secondary) |
| `primary_key` / `secondary_key` | For SDKs that take key and key name separately |
| `authorization_rule_id` | Consumed by the geo-DR pairing's `aliasAuthorizationRuleId` with zero translation |
| `*_connection_string_alias` | Populated only when the namespace carries a geo-DR pairing -- failover-stable credentials |

## Related Resources

- [Azure Service Bus Namespace](/docs/catalog/azure/azureservicebusnamespace) / [Queue](/docs/catalog/azure/azureservicebusqueue) / [Topic](/docs/catalog/azure/azureservicebustopic) -- the three scopes
- [Azure Role Assignment](/docs/catalog/azure/azureroleassignment) -- the keyless alternative

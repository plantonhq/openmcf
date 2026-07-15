# Azure Service Bus Disaster Recovery Config

Creates a geo-disaster-recovery pairing between two PREMIUM Service Bus namespaces -- continuous metadata replication plus a failover-stable alias DNS name clients connect through.

## What Gets Created

When you deploy an AzureServiceBusDisasterRecoveryConfig resource, Planton provisions:

- **Disaster Recovery Config (alias)** -- an `azurerm_servicebus_namespace_disaster_recovery_config` pairing the primary and partner namespaces under the alias name

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **Two PREMIUM AzureServiceBusNamespaces** in DIFFERENT regions; the partner must be EMPTY (no entities) at pairing time

## Quick Start

Create a file `geo-dr.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureServiceBusDisasterRecoveryConfig
metadata:
  name: orders-bus-dr
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureServiceBusDisasterRecoveryConfig.orders-bus-dr
spec:
  aliasName: myorg-orders-bus
  primaryNamespaceId:
    valueFrom:
      kind: AzureServiceBusNamespace
      name: orders-bus-eastus
      fieldPath: status.outputs.namespace_id
  partnerNamespaceId:
    valueFrom:
      kind: AzureServiceBusNamespace
      name: orders-bus-westus
      fieldPath: status.outputs.namespace_id
```

Deploy:

```shell
planton apply -f geo-dr.yaml
```

Point clients at the ALIAS connection strings, not either namespace -- that is what makes failover transparent. Failover itself is triggered on the secondary during an incident (an operational action, not a config change).

## Key Outputs

| Output | Purpose |
|--------|---------|
| `primary_connection_string_alias` | The failover-stable credential DR-aware clients hold |
| `alias_name` | The stable DNS identity (`{alias}.servicebus.windows.net`) |
| `default_primary_key` | The paired rule's key (rotation partner: `default_secondary_key`) |

## Related Resources

- [Azure Service Bus Namespace](/docs/catalog/azure/azureservicebusnamespace) -- the paired namespaces (PREMIUM)
- [Azure Service Bus Authorization Rule](/docs/catalog/azure/azureservicebusauthorizationrule) -- least-privilege alias credentials via `aliasAuthorizationRuleId`

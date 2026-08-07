# Azure Event Hub Disaster Recovery Config

Creates a geo-disaster-recovery pairing between two Event Hubs namespaces -- continuous metadata replication plus a failover-stable alias DNS name clients connect through.

## What Gets Created

When you deploy an AzureEventHubDisasterRecoveryConfig resource, Planton provisions:

- **Disaster Recovery Config (alias)** -- an `azurerm_eventhub_namespace_disaster_recovery_config` pairing the primary and partner namespaces under the alias name

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **Two AzureEventHubNamespaces** at STANDARD tier or above, in DIFFERENT regions, on the SAME tier; the partner must be EMPTY (no hubs) at pairing time

## Quick Start

Create a file `geo-dr.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureEventHubDisasterRecoveryConfig
metadata:
  name: telemetry-dr
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureEventHubDisasterRecoveryConfig.telemetry-dr
spec:
  aliasName: myorg-telemetry-alias
  primaryNamespaceId:
    valueFrom:
      kind: AzureEventHubNamespace
      name: telemetry-hubs-eastus
      fieldPath: status.outputs.namespace_id
  partnerNamespaceId:
    valueFrom:
      kind: AzureEventHubNamespace
      name: telemetry-hubs-westus
      fieldPath: status.outputs.namespace_id
```

Deploy:

```shell
planton apply -f geo-dr.yaml
```

Point clients at the ALIAS connection strings, not either namespace -- that is what makes failover transparent. Those credentials surface as the `*_connection_string_alias` outputs on the namespace and authorization-rule kinds once the pairing exists. Failover itself is triggered on the secondary during an incident (an operational action, not a config change). Metadata replicates; event data does not. Destroys take minutes by design: Azure breaks the pairing, deletes the config, and holds the alias name briefly before releasing it.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `disaster_recovery_config_id` | The Azure Resource Manager ID of the pairing (under the primary namespace) |
| `alias_name` | The stable DNS identity (`{alias}.servicebus.windows.net`) |

## Related Resources

- [Azure Event Hub Namespace](/docs/catalog/azure/azureeventhubnamespace) -- the paired namespaces; its `*_connection_string_alias` outputs carry the alias credentials
- [Azure Event Hub Authorization Rule](/docs/catalog/azure/azureeventhubauthorizationrule) -- least-privilege alias credentials via scoped rules

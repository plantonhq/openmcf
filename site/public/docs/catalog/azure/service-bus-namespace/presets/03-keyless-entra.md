---
title: "Keyless Namespace (Entra-Only)"
description: "This preset creates a namespace with SAS authentication disabled: no connection strings, no shared keys -- clients authenticate with Microsoft Entra identities holding Service Bus data-plane roles...."
type: "preset"
rank: "03"
presetSlug: "03-keyless-entra"
componentSlug: "service-bus-namespace"
componentTitle: "Service Bus Namespace"
provider: "azure"
icon: "package"
order: 3
---

# Keyless Namespace (Entra-Only)

This preset creates a namespace with SAS authentication disabled: no
connection strings, no shared keys -- clients authenticate with
Microsoft Entra identities holding Service Bus data-plane roles. The
zero-secret posture security teams standardize on.

## When to Use

- Estates standardizing on managed identities and workload identity
  federation (no secrets to rotate or leak)
- Anywhere a leaked connection string is an unacceptable risk

## Key Configuration Choices

- **`localAuthEnabled: false`** -- namespace-wide; the root rule's
  credential outputs still exist but stop being usable, and
  AzureServiceBusAuthorizationRule credentials will not work either

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `eastus` | The namespace's region | Your region strategy |
| `my-messaging-rg` | The AzureResourceGroup's Planton resource name | Your foundation composition |
| `myorg-keyless-bus` | 6-50 chars, globally unique | Your naming convention |
| `order-processing` | What this namespace carries | Your data taxonomy |

## Downstream Wiring

Grant an application's identity data-plane access:

```yaml
# On an AzureRoleAssignment
scope:
  valueFrom:
    kind: AzureServiceBusNamespace
    name: my-keyless-bus
    fieldPath: status.outputs.namespace_id
roleDefinitionName: Azure Service Bus Data Owner
principalId:
  valueFrom:
    kind: AzureUserAssignedIdentity
    name: my-app-identity
    fieldPath: status.outputs.principal_id
```

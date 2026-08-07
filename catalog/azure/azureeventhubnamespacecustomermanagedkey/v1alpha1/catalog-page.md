# Azure Event Hub Namespace Customer Managed Key

Applies customer-managed-key (BYOK) encryption onto an existing Event Hubs namespace -- event data at rest encrypted with your Key Vault keys instead of Microsoft-managed keys.

## What Gets Created

When you deploy an AzureEventHubNamespaceCustomerManagedKey resource, Planton provisions:

- **Customer Managed Key configuration** -- an `azurerm_eventhub_namespace_customer_managed_key` patching the referenced namespace's encryption to your Key Vault keys

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A single-tenant AzureEventHubNamespace** -- placed on a dedicated cluster or on the PREMIUM tier; Azure rejects the encryption patch on multi-tenant BASIC/STANDARD
- **An AzureKeyVaultKey** in a vault with purge protection enabled
- **Vault access for the unwrapping identity** -- the namespace's identity (system-assigned, or a user-assigned identity already in its identity block) needs wrap/unwrap access, e.g. an AzureRoleAssignment granting "Key Vault Crypto Service Encryption User"

## Quick Start

Create a file `cmk.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureEventHubNamespaceCustomerManagedKey
metadata:
  name: streaming-cmk
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureEventHubNamespaceCustomerManagedKey.streaming-cmk
spec:
  eventhubNamespaceId:
    valueFrom:
      kind: AzureEventHubNamespace
      name: premium-hubs
      fieldPath: status.outputs.namespace_id
  keyVaultKeyIds:
    - valueFrom:
        kind: AzureKeyVaultKey
        name: streaming-key
        fieldPath: status.outputs.versionless_id
```

Deploy:

```shell
planton apply -f cmk.yaml
```

CMK is a second step by Azure's design: create the namespace with its identity, grant that identity vault access, then apply this resource. Versionless key references make vault-side rotation propagate automatically. The lifecycle is ADD-ONLY -- once enabled, CMK can never be removed (deleting this resource changes nothing on the namespace; returning to Microsoft-managed keys means replacing the namespace).

## Key Outputs

| Output | Purpose |
|--------|---------|
| `customer_managed_key_id` | The provider's identity for the configuration -- the namespace's ARM ID (CMK is a namespace property, not an ARM object of its own) |

## Related Resources

- [Azure Event Hub Namespace](/docs/catalog/azure/azureeventhubnamespace) -- the namespace to encrypt (dedicated cluster or PREMIUM)
- [Azure Event Hub Cluster](/docs/catalog/azure/azureeventhubcluster) -- dedicated placement that unlocks CMK eligibility
- [Azure Key Vault Key](/docs/catalog/azure/azurekeyvaultkey) -- the encryption keys (versionless for auto-rotation)
- [Azure Role Assignment](/docs/catalog/azure/azureroleassignment) -- the wrap/unwrap grant on the vault
- [Azure User Assigned Identity](/docs/catalog/azure/azureuserassignedidentity) -- an optional dedicated unwrapping identity

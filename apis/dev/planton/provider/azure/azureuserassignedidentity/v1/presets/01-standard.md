# Standard Managed Identity

This preset creates a plain user-assigned managed identity -- the anchor of
every keyless-auth story on Azure. The identity is deliberately just the
identity: what it may do and who may act as it are separate, composable
resources you add alongside it.

Grant it permissions with an `AzureRoleAssignment` referencing the identity's
`principal_id` output:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRoleAssignment
metadata:
  name: my-identity-kv-reader
spec:
  scope:
    value: "<target-resource-id>"
  roleDefinitionName: Key Vault Secrets User
  principalId:
    valueFrom:
      name: my-identity
```

## When to Use

- Any workload (AKS pods, Container Apps, Function Apps, VMs) that needs
  credential-free authentication to Azure services
- As the shared anchor an application's grants and trust rules attach to
- Whenever you would otherwise create a service principal with a client
  secret

## Key Configuration Choices

- **Identity only** -- permissions live in `AzureRoleAssignment` resources
  and external trust in `AzureFederatedIdentityCredential` resources, each
  individually reviewable and removable
- **Durable name** -- renaming replaces the identity and mints a new
  principal, invalidating existing grants; name it after the workload or
  duty ("ci-deployer", "payments-api")

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-identity-name>` | Name for the managed identity (3-128 chars) | Your naming convention |

## Related Presets

- **02-ci-deployer** -- the complete keyless-CI composition
- **03-governance-tagged** -- identity carrying org governance tags and regional isolation

# Premium Isolated Namespace

This preset creates a PREMIUM namespace with the enterprise posture:
dedicated messaging units, customer-managed-key encryption, and a
deny-by-default VNet firewall admitting only the application subnets.

## When to Use

- Compliance workloads that mandate BYOK encryption at rest
- Namespaces that must be unreachable from the public internet
- Predictable-latency workloads or messages beyond 256 KB

## Key Configuration Choices

- **`capacity: 1`** -- one messaging unit; scale to 2/4/8/16 in place as
  throughput grows
- **`premiumMessagingPartitions: 1`** -- not partitioned; the layout is
  ForceNew, so choose 2/4 up front only when one store's throughput
  ceiling is a real constraint
- **`customerManagedKey`** -- rides the versionless key ID so rotations
  propagate automatically; grant the identity "Key Vault Crypto Service
  Encryption User" on the vault; CMK is irreversible once set
- **`networkRuleSet` with DENY** -- trusted Microsoft services stay
  admitted (`trustedServicesAllowed`) so platform integrations keep
  working

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `eastus` | The namespace's region | Your region strategy |
| `my-messaging-rg` | The AzureResourceGroup's Planton resource name | Your foundation composition |
| `myorg-enterprise-bus` | 6-50 chars, globally unique | Your naming convention |
| `my-messaging-identity` | The AzureUserAssignedIdentity holding vault access | Your identity composition |
| `my-messaging-cmk` | The AzureKeyVaultKey used for BYOK | Your Key Vault composition |
| `my-app-subnet` | The AzureSubnet admitted through the firewall | Your network composition |
| `payments` | What this namespace carries | Your data taxonomy |

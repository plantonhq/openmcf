# Identity-Secured Worker

This preset deploys a Python processing worker with a fully credential-free posture: the app authenticates to its deployment storage as its own managed identity, and the classic username/password publishing path is closed -- no storage key or deploy password exists anywhere in the configuration.

## When to Use

- Queue, event, and blob processing where idle cost must be zero
- Environments with a no-static-credentials policy
- Workloads that want the 4 GB instance size for memory-heavy processing

## Key Configuration Choices

- **`SYSTEM_ASSIGNED_IDENTITY` storage auth** -- no key travels anywhere; the grant is day-2 by construction: create the app, then give its `identity_principal_id` output "Storage Blob Data Contributor" on the storage account (package deployments fail until the grant lands)
- **`webdeployPublishBasicAuthenticationEnabled: false`** -- closes basic-auth publishing; pair with identity-based CI/CD (the `site_credential_password` output stops being usable)
- **`instanceMemoryInMb: 4096`** -- the largest flex instance size, for memory-heavy batches
- **No `alwaysReady` pool** -- a worker tolerates cold starts; the app costs nothing idle

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The Planton name of your `AzureResourceGroup` resource | Planton console |
| `<your-flex-service-plan>` | The Planton name of your FC1-SKU `AzureServicePlan` | Planton console |
| `<your-storage-account>` | The account holding the deployment container (its name forms the endpoint) | Planton console |

## Related Presets

- **Node HTTP API** -- the warm-path HTTP shape
- **Entra-Protected API** -- the platform-authenticated shape

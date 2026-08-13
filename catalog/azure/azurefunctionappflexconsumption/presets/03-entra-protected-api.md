# Entra-Protected API

This preset deploys a .NET isolated function API behind App Service's built-in authentication: Azure validates Entra ID tokens at the platform layer -- before any request reaches function code -- with Application Insights telemetry wired in.

## When to Use

- Internal APIs that every caller must authenticate to with an organizational identity
- Function backends consumed by SPAs or services that already hold Entra tokens
- Teams that want auth enforced by the platform rather than re-implemented per function

## Key Configuration Choices

- **`unauthenticatedAction: RETURN_401`** -- APIs reject anonymous requests outright; use the login-redirect default only for browser-facing apps
- **The secret lives in `appSettings`, by NAME** -- `clientSecretSettingName` points at `AAD_CLIENT_SECRET`, whose value is a Key Vault reference; the auth block itself never carries a secret
- **`login.tokenStoreEnabled: true`** -- durably stores tokens so `/.auth/me` and token refresh work
- **`applicationInsightsConnectionString`** -- references the `AzureApplicationInsights` resource's output; requests, dependencies, and exceptions land in APM automatically

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The Planton name of your `AzureResourceGroup` resource | Planton console |
| `<your-flex-service-plan>` | The Planton name of your FC1-SKU `AzureServicePlan` | Planton console |
| `<your-storage-account>` | The account holding the deployment container | Planton console |
| `<your-application-insights>` | The Planton name of your `AzureApplicationInsights` resource | Planton console |
| `<your-entra-app-client-id>` | The Entra app registration's application (client) ID | Entra admin center -> App registrations |
| `<your-tenant-id>` | Your Entra tenant GUID | Entra admin center -> Overview |
| `<your-secret-uri>` | The Key Vault secret URI holding the client secret | The vault's secret page |

## Related Presets

- **Node HTTP API** -- the warm-path HTTP shape
- **Identity-Secured Worker** -- the credential-free storage-auth shape

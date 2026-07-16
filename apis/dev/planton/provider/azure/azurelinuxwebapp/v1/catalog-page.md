# Azure Linux Web App

Deploys an Azure Linux Web App -- a managed web hosting platform for running long-lived web applications, APIs, and containerized services on Azure App Service. Supports .NET, Node.js, Python, PHP, Java (with Tomcat, JBoss EAP, or embedded SE), and Docker containers with configurable managed identity, VNet integration, built-in authentication (Easy Auth v2), auto-heal, scheduled backups, Application Insights telemetry, logging to file system or blob storage, IP restrictions, CORS, and connection strings.

## What Gets Created

When you deploy an AzureLinuxWebApp resource, Planton provisions:

- **Linux Web App** -- an `azurerm_linux_web_app` / `appservice.LinuxWebApp` resource in the specified region and resource group, configured with the chosen application stack, operational settings, logging, authentication, and security configuration
- **Managed Identity** -- created only when `identity` is configured, provides credential-free authentication to Azure services
- **VNet Integration** -- created only when `virtualNetworkSubnetId` is set, routes outbound traffic through a VNet subnet
- **Azure Tags** -- platform metadata tags merged with your own `tags` (your tags win on key collision)

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An Azure Resource Group** where the web app will be created (can reference an AzureResourceGroup resource)
- **An Azure Service Plan** providing compute resources -- Basic (`BASIC_B1`-`B3`) for dedicated compute, Standard (`STANDARD_S1`-`S3`) for autoscale and deployment slots, or Premium (`PREMIUM_P1V3`-`P3V3`) for enhanced performance and zone redundancy
- **A globally unique app name** -- the name becomes the hostname `{webAppName}.azurewebsites.net`

## Quick Start

Create a file `webapp.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureLinuxWebApp
metadata:
  name: my-web
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureLinuxWebApp.my-web
spec:
  region: eastus
  resourceGroup:
    value: my-rg
  webAppName: my-web-app
  servicePlanId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Web/serverFarms/my-plan
  siteConfig:
    applicationStack:
      nodeVersion: "20-lts"
```

Deploy:

```shell
planton apply -f webapp.yaml
```

This creates a Node.js 20 LTS Web App with HTTPS-only access, TLS 1.2, and 64-bit worker processes.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region for the web app. **ForceNew**. | Required, minimum length 1 |
| `resourceGroup` | `StringValueOrRef` | Azure Resource Group name. Can reference an AzureResourceGroup resource via `valueFrom`. **ForceNew**. | Required |
| `webAppName` | `string` | Globally unique app name. Becomes `{webAppName}.azurewebsites.net`. **ForceNew**. | Required, 2-60 characters, pattern `^[a-zA-Z0-9][a-zA-Z0-9-]{0,58}[a-zA-Z0-9]$` |
| `servicePlanId` | `StringValueOrRef` | Service Plan providing compute resources. Can reference an AzureServicePlan resource via `valueFrom`. | Required |
| `siteConfig` | `object` | Site configuration containing the application stack, security dials, and auto-heal rules. | Required |
| `siteConfig.applicationStack` | `object` | Runtime selection. Exactly one runtime: `dotnetVersion`, `nodeVersion`, `pythonVersion`, `phpVersion`, `javaVersion` (with `javaServer` + `javaServerVersion`), or `docker`. | Required |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `httpsOnly` | `bool` | `true` | Redirect all HTTP to HTTPS. |
| `publicNetworkAccessEnabled` | `bool` | `true` | Allow public internet access. |
| `enabled` | `bool` | `true` | Enable/disable the web app without deleting it. |
| `clientAffinityEnabled` | `bool` | `false` | Enable ARR session affinity cookies. Use for stateful apps only. |
| `clientCertificateEnabled` | `bool` | `false` | Require/request client certificates (mutual TLS). |
| `clientCertificateMode` | `enum` | `OPTIONAL` | `REQUIRED`, `OPTIONAL`, or `OPTIONAL_INTERACTIVE_USER`. |
| `applicationInsightsConnectionString` | `StringValueOrRef` | -- | App Insights connection string. Can reference an AzureApplicationInsights resource via `valueFrom`. |
| `virtualNetworkSubnetId` | `StringValueOrRef` | -- | Subnet ID for VNet integration. Can reference an AzureSubnet resource via `valueFrom`. |
| `vnetImagePullEnabled` | `bool` | `false` | Pull container images over the VNet (private registries). |
| `virtualNetworkBackupRestoreEnabled` | `bool` | `false` | Route backup/restore traffic over the VNet. |
| `identity.type` | `enum` | -- | Managed identity: `SYSTEM_ASSIGNED`, `USER_ASSIGNED`, or `SYSTEM_AND_USER_ASSIGNED`. |
| `keyVaultReferenceIdentityId` | `StringValueOrRef` | -- | Identity used for Key Vault references in app settings. |
| `webdeployPublishBasicAuthenticationEnabled` | `bool` | `true` | Basic-auth publishing over Web Deploy; disable to harden. |
| `ftpPublishBasicAuthenticationEnabled` | `bool` | `true` | Basic-auth publishing over FTP/FTPS; disable to harden. |
| `zipDeployFile` | `string` | -- | Local ZIP package deployed on create/update. |
| `appSettings` | `map<string, string>` | `{}` | Application environment variables. |
| `connectionStrings` | `list` | `[]` | Named connection strings with `name`, enum `type` (e.g. `SQL_AZURE`, `POSTGRESQL`, `CUSTOM`), and `value`. |
| `stickySettings` | `object` | -- | App setting / connection string names pinned to the production slot during swaps. |
| `storageMounts` | `list` | `[]` | Azure Files (`AZURE_FILES`) or Blob (`AZURE_BLOB`) mounts. |
| `backup` | `object` | -- | Scheduled backups to a storage container SAS URL (Standard tier+): cadence (`DAY`/`HOUR`), retention, keep-at-least-one. |
| `authSettingsV2` | `object` | -- | Built-in authentication: platform settings, `login` behavior, and provider blocks (`activeDirectoryV2`, `githubV2`, `googleV2`, `microsoftV2`, `facebookV2`, `appleV2`, `twitterV2`, repeated `customOidcV2`). Secrets are referenced by app setting name. |
| `siteConfig.alwaysOn` | `bool` | -- | Keep app loaded in memory. Critical for Standard/Premium plans; rejected on Free/Shared. |
| `siteConfig.healthCheckPath` | `string` | -- | Health check endpoint (e.g., `/health`). |
| `siteConfig.healthCheckEvictionTimeInMin` | `int` | -- | Minutes before unhealthy instance eviction (2-10). Requires the path. |
| `siteConfig.minimumTlsVersion` | `enum` | `TLS_1_2` | TLS floor: `TLS_1_0` - `TLS_1_3` (SCM twin: `scmMinimumTlsVersion`). |
| `siteConfig.minimumTlsCipherSuite` | `string` | -- | Minimum accepted cipher suite (Azure's `TLS_*` identifiers). |
| `siteConfig.ftpsState` | `enum` | `DISABLED` | FTP endpoint posture: `ALL_ALLOWED`, `FTPS_ONLY`, `DISABLED`. |
| `siteConfig.loadBalancingMode` | `enum` | `LEAST_REQUESTS` | Request distribution across instances. |
| `siteConfig.managedPipelineMode` | `enum` | `INTEGRATED` | Request pipeline mode (`CLASSIC` for legacy compatibility). |
| `siteConfig.autoHealSetting` | `object` | -- | Recycle on `requests`, `statusCodes`, `slowRequest`, or `slowRequestWithPath` triggers, guarded by `minimumProcessExecutionTime`. |
| `siteConfig.cors.allowedOrigins` | `string[]` | -- | CORS allowed origins (wildcard cannot combine with credentials). |
| `siteConfig.ipRestrictions` | `list` | `[]` | IP/service-tag/subnet access rules with `ALLOW`/`DENY` actions and Front Door header filters. |
| `siteConfig.ipRestrictionDefaultAction` | `enum` | `ALLOW` | Action for unmatched traffic (`DENY` = allow-list semantics). |
| `logs.applicationLogs.fileSystemLevel` | `enum` | `ERROR` | `OFF`, `ERROR`, `WARNING`, `INFORMATION`, `VERBOSE`; blob-storage destination available. |
| `logs.httpLogs` | `object` | -- | Exactly one destination: `fileSystem` (25-100 MB rotation) or `azureBlobStorage` (SAS URL). |
| `logs.failedRequestTracing` | `bool` | `false` | Capture detailed traces for failed requests. |
| `logs.detailedErrorMessages` | `bool` | `false` | Return detailed error pages. Disable in production. |
| `tags` | `map` | -- | Free-form Azure resource tags merged over metadata-derived tags. |

## Examples

### Node.js Web API

A Node.js 20 LTS Web App with Application Insights and health checks:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureLinuxWebApp
metadata:
  name: node-api
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureLinuxWebApp.node-api
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  webAppName: node-api-app
  servicePlanId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/prod-rg/providers/Microsoft.Web/serverFarms/prod-plan
  applicationInsightsConnectionString:
    value: "InstrumentationKey=00000000-0000-0000-0000-000000000000;IngestionEndpoint=https://eastus-0.in.applicationinsights.azure.com/"
  siteConfig:
    applicationStack:
      nodeVersion: "20-lts"
    alwaysOn: true
    healthCheckPath: /health
    http2Enabled: true
  appSettings:
    NODE_ENV: production
    DATABASE_URL: "postgresql://..."
  logs:
    applicationLogs:
      fileSystemLevel: INFORMATION
    httpLogs:
      fileSystem:
        retentionInMb: 50
        retentionInDays: 7
```

### Docker Container Web App

A containerized Web App with VNet integration and managed identity:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureLinuxWebApp
metadata:
  name: docker-web
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureLinuxWebApp.docker-web
spec:
  region: westeurope
  resourceGroup:
    value: prod-rg
  webAppName: docker-web-app
  servicePlanId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/prod-rg/providers/Microsoft.Web/serverFarms/premium-plan
  virtualNetworkSubnetId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/webapp
  identity:
    type: SYSTEM_ASSIGNED
  siteConfig:
    applicationStack:
      docker:
        registryUrl: https://myregistry.azurecr.io
        imageName: myorg/my-web-app
        imageTag: v2.0.0
    containerRegistryUseManagedIdentity: true
    alwaysOn: true
    healthCheckPath: /healthz
    vnetRouteAllEnabled: true
```

### Web App with Built-in Authentication (Easy Auth)

A Web App that authenticates every request against Microsoft Entra ID before it reaches application code. The client secret lives in an app setting; the auth block references it by name:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureLinuxWebApp
metadata:
  name: internal-portal
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureLinuxWebApp.internal-portal
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  webAppName: internal-portal-app
  servicePlanId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/prod-rg/providers/Microsoft.Web/serverFarms/prod-plan
  appSettings:
    AAD_CLIENT_SECRET: "@Microsoft.KeyVault(SecretUri=https://prod-kv.vault.azure.net/secrets/portal-aad-secret)"
  authSettingsV2:
    authEnabled: true
    requireAuthentication: true
    unauthenticatedAction: REDIRECT_TO_LOGIN_PAGE
    excludedPaths:
      - /health
    login:
      tokenStoreEnabled: true
    activeDirectoryV2:
      clientId: 11111111-2222-3333-4444-555555555555
      tenantAuthEndpoint: https://login.microsoftonline.com/v2.0/99999999-8888-7777-6666-555555555555/
      clientSecretSettingName: AAD_CLIENT_SECRET
  siteConfig:
    applicationStack:
      dotnetVersion: "8.0"
    alwaysOn: true
    healthCheckPath: /health
```

### Enterprise Private Web App

A Premium-tier Web App with private-only access, client certificate authentication, auto-heal, hardened publishing, and comprehensive logging:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureLinuxWebApp
metadata:
  name: private-web
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureLinuxWebApp.private-web
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  webAppName: private-web-app
  servicePlanId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/prod-rg/providers/Microsoft.Web/serverFarms/premium-plan
  publicNetworkAccessEnabled: false
  clientCertificateEnabled: true
  clientCertificateMode: REQUIRED
  webdeployPublishBasicAuthenticationEnabled: false
  ftpPublishBasicAuthenticationEnabled: false
  identity:
    type: SYSTEM_ASSIGNED
  siteConfig:
    applicationStack:
      dotnetVersion: "8.0"
    alwaysOn: true
    healthCheckPath: /api/health
    ipRestrictionDefaultAction: DENY
    autoHealSetting:
      trigger:
        statusCodes:
          - statusCodeRange: 500-599
            count: 50
            interval: "00:05:00"
      minimumProcessExecutionTime: "00:01:00"
  logs:
    applicationLogs:
      fileSystemLevel: WARNING
    httpLogs:
      fileSystem:
        retentionInMb: 100
        retentionInDays: 30
    failedRequestTracing: true
```

### Using Foreign Key References

Reference Planton-managed resources:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureLinuxWebApp
metadata:
  name: ref-web
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureLinuxWebApp.ref-web
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: my-rg
      fieldPath: status.outputs.resource_group_name
  webAppName: ref-web-app
  servicePlanId:
    valueFrom:
      kind: AzureServicePlan
      name: my-plan
      fieldPath: status.outputs.service_plan_id
  applicationInsightsConnectionString:
    valueFrom:
      kind: AzureApplicationInsights
      name: my-insights
      fieldPath: status.outputs.connection_string
  siteConfig:
    applicationStack:
      pythonVersion: "3.12"
    alwaysOn: true
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `web_app_id` | `string` | Azure Resource Manager ID of the Web App |
| `default_hostname` | `string` | Default hostname (`{webAppName}.azurewebsites.net`) |
| `outbound_ip_addresses` | `string[]` | Currently active outbound IP addresses |
| `possible_outbound_ip_addresses` | `string[]` | Every outbound IP the platform could ever use -- for durable firewall allowlists |
| `identity_principal_id` | `string` | System-assigned identity principal ID (when identity is configured) |
| `identity_tenant_id` | `string` | System-assigned identity tenant ID |
| `custom_domain_verification_id` | `string` | TXT record value for custom domain verification |
| `kind` | `string` | Resource kind (e.g., `app,linux`) |
| `hosting_environment_id` | `string` | ARM ID of the App Service Environment (empty outside ASE) |
| `site_credential_name` | `string` | Publishing credential username (Kudu/SCM basic auth) |
| `site_credential_password` | `string` | Publishing credential password -- treat like an admin password |

## Related Components

- [AzureServicePlan](/docs/catalog/azure/azureserviceplan) -- provides the compute tier for the Web App
- [AzureApplicationInsights](/docs/catalog/azure/azureapplicationinsights) -- provides APM telemetry collection
- [AzureResourceGroup](/docs/catalog/azure/azureresourcegroup) -- provides the resource group for app placement
- [AzureSubnet](/docs/catalog/azure/azuresubnet) -- provides VNet integration for outbound connectivity
- [AzureFrontDoorProfile](/docs/catalog/azure/azurefrontdoorprofile) -- global CDN and load balancing for the Web App

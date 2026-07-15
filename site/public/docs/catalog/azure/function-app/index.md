---
title: "Function App"
description: "Function App deployment documentation"
icon: "package"
order: 100
componentName: "azurefunctionapp"
---

# Azure Function App

Deploys an Azure Linux Function App -- a serverless compute platform for event-driven workloads supporting HTTP triggers, queue triggers, timer schedules, and more. The component provides full configuration of the application runtime stack, all three storage bindings, serverless scaling dials, built-in authentication (Easy Auth v2), scheduled backups, managed identity, VNet integration, Application Insights telemetry, IP restrictions, CORS, storage mounts, and connection strings.

## What Gets Created

When you deploy an AzureFunctionApp resource, Planton provisions:

- **Linux Function App** -- an `azurerm_linux_function_app` / `appservice.LinuxFunctionApp` resource in the specified region and resource group, configured with the chosen runtime stack, storage binding, Application Insights connection, authentication, and operational settings
- **Managed Identity** -- created only when `identity` is configured, provides credential-free authentication to Azure services (Key Vault, Storage, ACR)
- **VNet Integration** -- created only when `virtualNetworkSubnetId` is set, routes outbound traffic through a VNet subnet for private connectivity
- **Azure Tags** -- platform metadata tags merged with your own `tags` (your tags win on key collision)

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An Azure Resource Group** where the function app will be created (can reference an AzureResourceGroup resource)
- **An Azure Service Plan** providing compute resources -- Consumption (`CONSUMPTION_Y1`) for pay-per-execution, Elastic Premium (`ELASTIC_PREMIUM_EP1`-`EP3`) for pre-warmed instances, or Dedicated tiers for reserved capacity
- **An Azure Storage Account** for Function App runtime state (trigger management, logs, coordination) -- or a Key Vault secret holding its connection string
- **A globally unique app name** -- the name becomes the hostname `{functionAppName}.azurewebsites.net`

## Quick Start

Create a file `functionapp.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFunctionApp
metadata:
  name: my-func
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureFunctionApp.my-func
spec:
  region: eastus
  resourceGroup:
    value: my-rg
  functionAppName: my-func-app
  servicePlanId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Web/serverFarms/my-plan
  storageAccountName:
    value: mystorageacct
  storageAccountAccessKey:
    value: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=="
  siteConfig:
    applicationStack:
      pythonVersion: "3.12"
```

Deploy:

```shell
planton apply -f functionapp.yaml
```

This creates a Python 3.12 Function App on the specified Service Plan with HTTPS-only access, TLS 1.2, and Functions runtime v4.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region for the function app. **ForceNew**. | Required, minimum length 1 |
| `resourceGroup` | `StringValueOrRef` | Azure Resource Group name. Can reference an AzureResourceGroup resource via `valueFrom`. **ForceNew**. | Required |
| `functionAppName` | `string` | Globally unique app name. Becomes `{functionAppName}.azurewebsites.net`. **ForceNew**. | Required, 2-60 characters, alphanumeric + hyphens |
| `servicePlanId` | `StringValueOrRef` | Service Plan providing compute resources. Can reference an AzureServicePlan resource via `valueFrom`. | Required |
| `storageAccountName` XOR `storageKeyVaultSecretId` | - | The runtime-state storage binding: an account name (reference-able) or a Key Vault secret URL holding the connection string. | Exactly one |
| `siteConfig` | `object` | Site configuration containing the application stack. | Required |
| `siteConfig.applicationStack` | `object` | Runtime selection. Exactly one runtime: `dotnetVersion` (+ `useDotnetIsolatedRuntime`), `nodeVersion`, `pythonVersion`, `javaVersion`, `powershellCoreVersion`, `docker`, or `useCustomRuntime`. | Required |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `storageAccountAccessKey` | `StringValueOrRef` | -- | Storage access key (conflicts with `storageUsesManagedIdentity`). |
| `storageUsesManagedIdentity` | `bool` | `false` | Credential-free storage access via the app's identity (needs Storage Blob Data Owner + Storage Queue Data Contributor). |
| `functionsExtensionVersion` | `string` | `"~4"` | Functions host version. |
| `dailyMemoryTimeQuota` | `int` | `0` | Consumption cost circuit breaker (GB-seconds per day; 0 = unlimited). |
| `httpsOnly` | `bool` | `true` | Redirect all HTTP to HTTPS. |
| `publicNetworkAccessEnabled` | `bool` | `true` | Allow public internet access. |
| `enabled` | `bool` | `true` | Enable/disable the app without deleting it. |
| `builtinLoggingEnabled` | `bool` | `true` | Legacy AzureWebJobsDashboard logging; disable when App Insights is configured. |
| `contentShareForceDisabled` | `bool` | `false` | Skip the auto-created Azure Files content share. |
| `clientCertificateEnabled` / `clientCertificateMode` | - | `false` / `OPTIONAL` | Mutual TLS posture (`REQUIRED`, `OPTIONAL`, `OPTIONAL_INTERACTIVE_USER`). |
| `virtualNetworkSubnetId` | `StringValueOrRef` | -- | Subnet ID for VNet integration (not supported on Consumption). |
| `vnetImagePullEnabled` | `bool` | `false` | Pull container images over the VNet. |
| `virtualNetworkBackupRestoreEnabled` | `bool` | `false` | Route backup/restore traffic over the VNet. |
| `identity.type` | `enum` | -- | `SYSTEM_ASSIGNED`, `USER_ASSIGNED`, or `SYSTEM_AND_USER_ASSIGNED`. |
| `keyVaultReferenceIdentityId` | `StringValueOrRef` | -- | Identity used for Key Vault references in app settings. |
| `webdeployPublishBasicAuthenticationEnabled` | `bool` | `true` | Basic-auth publishing over Web Deploy; disable to harden. |
| `ftpPublishBasicAuthenticationEnabled` | `bool` | `true` | Basic-auth publishing over FTP/FTPS; disable to harden. |
| `zipDeployFile` | `string` | -- | Local ZIP package deployed on create/update. |
| `appSettings` | `map<string, string>` | `{}` | Application environment variables. |
| `connectionStrings` | `list` | `[]` | Named connection strings with `name`, enum `type`, and `value`. |
| `stickySettings` | `object` | -- | Setting names pinned to the production slot during swaps. |
| `storageMounts` | `list` | `[]` | Azure Files (`AZURE_FILES`) or Blob (`AZURE_BLOB`) mounts. |
| `backup` | `object` | -- | Scheduled backups to a storage container SAS URL (Standard tier+). |
| `authSettingsV2` | `object` | -- | Built-in authentication: platform settings, `login` behavior, provider blocks (`activeDirectoryV2`, `githubV2`, `googleV2`, `microsoftV2`, `facebookV2`, `appleV2`, `twitterV2`, repeated `customOidcV2`). Secrets referenced by app setting name. |
| `siteConfig.appScaleLimit` | `int` | -- | Scale-out cap for Consumption/Elastic Premium (the cost lever). |
| `siteConfig.elasticInstanceMinimum` | `int` | -- | Always-ready instances on Elastic Premium. |
| `siteConfig.preWarmedInstanceCount` | `int` | -- | Warm instances beyond the minimum (Elastic Premium). |
| `siteConfig.runtimeScaleMonitoringEnabled` | `bool` | -- | KEDA-based trigger scale monitoring (EP*/Dedicated, Functions v4+). |
| `siteConfig.alwaysOn` | `bool` | -- | Keep app loaded (Dedicated plans; auto-managed on serverless tiers). |
| `siteConfig.healthCheckPath` (+ `healthCheckEvictionTimeInMin`) | - | -- | Health probe + eviction window (2-10 min). |
| `siteConfig.minimumTlsVersion` | `enum` | `TLS_1_2` | TLS floor (SCM twin: `scmMinimumTlsVersion`); `minimumTlsCipherSuite` for the cipher floor. |
| `siteConfig.ftpsState` | `enum` | `DISABLED` | FTP endpoint posture. |
| `siteConfig.ipRestrictions` (+ default actions) | `list` | `[]` | IP/service-tag/subnet access rules with `ALLOW`/`DENY`. |
| `siteConfig.appServiceLogs` | `object` | -- | Disk quota (25-100 MB) + retention days. |
| `tags` | `map` | -- | Free-form Azure resource tags merged over metadata-derived tags. |

## Examples

### Consumption HTTP API with Managed-Identity Storage

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFunctionApp
metadata:
  name: events-fn
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureFunctionApp.events-fn
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  functionAppName: events-fn-app
  servicePlanId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/prod-rg/providers/Microsoft.Web/serverFarms/consumption-plan
  storageAccountName:
    value: prodfnstorage
  storageUsesManagedIdentity: true
  identity:
    type: SYSTEM_ASSIGNED
  dailyMemoryTimeQuota: 1000000
  siteConfig:
    applicationStack:
      pythonVersion: "3.12"
    appScaleLimit: 100
```

### Elastic Premium with Pre-Warmed Instances and VNet

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFunctionApp
metadata:
  name: private-fn
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureFunctionApp.private-fn
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  functionAppName: private-fn-app
  servicePlanId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/prod-rg/providers/Microsoft.Web/serverFarms/ep1-plan
  storageAccountName:
    value: prodfnstorage
  storageAccountAccessKey:
    value: "<storage-key>"
  virtualNetworkSubnetId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/functions
  siteConfig:
    applicationStack:
      nodeVersion: "20"
    elasticInstanceMinimum: 2
    preWarmedInstanceCount: 3
    runtimeScaleMonitoringEnabled: true
    vnetRouteAllEnabled: true
```

### Using Foreign Key References

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFunctionApp
metadata:
  name: ref-fn
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureFunctionApp.ref-fn
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: my-rg
      fieldPath: status.outputs.resource_group_name
  functionAppName: ref-fn-app
  servicePlanId:
    valueFrom:
      kind: AzureServicePlan
      name: fn-plan
      fieldPath: status.outputs.service_plan_id
  storageAccountName:
    valueFrom:
      kind: AzureStorageAccount
      name: fn-storage
      fieldPath: status.outputs.storage_account_name
  storageAccountAccessKey:
    valueFrom:
      kind: AzureStorageAccount
      name: fn-storage
      fieldPath: status.outputs.primary_access_key
  applicationInsightsConnectionString:
    valueFrom:
      kind: AzureApplicationInsights
      name: fn-insights
      fieldPath: status.outputs.connection_string
  siteConfig:
    applicationStack:
      pythonVersion: "3.12"
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `function_app_id` | `string` | Azure Resource Manager ID of the Function App |
| `default_hostname` | `string` | Default hostname (`{functionAppName}.azurewebsites.net`) |
| `outbound_ip_addresses` | `string[]` | Currently active outbound IP addresses |
| `possible_outbound_ip_addresses` | `string[]` | Every outbound IP the platform could ever use -- for durable firewall allowlists |
| `identity_principal_id` | `string` | System-assigned identity principal ID (when identity is configured) |
| `identity_tenant_id` | `string` | System-assigned identity tenant ID |
| `custom_domain_verification_id` | `string` | TXT record value for custom domain verification |
| `kind` | `string` | Resource kind (e.g., `functionapp,linux`) |
| `hosting_environment_id` | `string` | ARM ID of the App Service Environment (empty outside ASE) |
| `site_credential_name` | `string` | Publishing credential username (Kudu/SCM basic auth) |
| `site_credential_password` | `string` | Publishing credential password -- treat like an admin password |

## Related Components

- [AzureServicePlan](/docs/catalog/azure/service-plan) -- provides the compute tier for the Function App
- [AzureStorageAccount](/docs/catalog/azure/storage-account) -- provides the runtime-state storage binding
- [AzureApplicationInsights](/docs/catalog/azure/application-insights) -- provides APM telemetry collection
- [AzureResourceGroup](/docs/catalog/azure/resource-group) -- provides the resource group for app placement
- [AzureSubnet](/docs/catalog/azure/subnet) -- provides VNet integration for outbound connectivity

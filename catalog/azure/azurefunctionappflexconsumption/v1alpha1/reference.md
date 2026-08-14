# AzureFunctionAppFlexConsumption

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Maximal offline-validation manifest for
# AzureFunctionAppFlexConsumption: every configurable arm populated
# with literal values -- the user-assigned-identity deployment-storage
# mode (the two live lanes prove the connection-string and
# system-assigned modes), all scale dials with multiple always-ready
# pools, the full site_config surface, Easy Auth v2 with EVERY identity
# provider plus a custom OIDC entry, connection strings, sticky
# settings, and the combined identity model. Drives validate-manifest
# and the offline tofu plan's rendered-value inspection; never deployed
# live.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureFunctionAppFlexConsumption
metadata:
  name: flex-consumption-maximal
  org: planton-oss
  env: e2e
spec:
  region: centralus
  resourceGroup:
    value: app-rg
  functionAppName: acme-orders-flex
  servicePlanId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Web/serverFarms/flex-plan
  storageContainerEndpoint: https://acmeordersflexsa.blob.core.windows.net/deployments
  storageAuthenticationType: USER_ASSIGNED_IDENTITY
  storageUserAssignedIdentityId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/orders-mi
  runtimeName: DOTNET_ISOLATED
  runtimeVersion: "8.0"
  instanceMemoryInMb: 4096
  maximumInstanceCount: 250
  httpConcurrency: 32
  alwaysReady:
    - name: http
      instanceCount: 2
    - name: durable
      instanceCount: 1
    - name: function:ProcessOrders
      instanceCount: 1
  siteConfig:
    apiManagementApiId: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.ApiManagement/service/acme-apim/apis/orders
    apiDefinitionUrl: https://acme.example.com/openapi.json
    appCommandLine: ""
    applicationInsightsKey: 00000000-0000-0000-0000-000000000000
    appServiceLogs:
      diskQuotaMb: 50
      retentionPeriodDays: 7
    containerRegistryUseManagedIdentity: true
    defaultDocuments:
      - index.html
    elasticInstanceMinimum: 1
    http2Enabled: true
    ipRestrictions:
      - name: front-door-only
        priority: 100
        action: ALLOW
        serviceTag: AzureFrontDoor.Backend
        description: Only Front Door reaches the origin
        headers:
          xForwardedFor:
            - 203.0.113.0/24
          xForwardedHost:
            - orders.acme.example.com
          xAzureFdid:
            - value: 11111111-2222-3333-4444-555555555555
          xFdHealthProbe:
            - "1"
      - name: office
        priority: 200
        action: DENY
        ipAddress: 198.51.100.0/24
      - name: hub-subnet
        priority: 300
        action: ALLOW
        virtualNetworkSubnetId:
          value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/net-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet/subnets/agents
    ipRestrictionDefaultAction: DENY
    scmUseMainIpRestriction: false
    scmIpRestrictions:
      - name: ci-runners
        priority: 100
        action: ALLOW
        ipAddress: 192.0.2.0/24
    scmIpRestrictionDefaultAction: DENY
    loadBalancingMode: WEIGHTED_ROUND_ROBIN
    managedPipelineMode: INTEGRATED
    remoteDebuggingEnabled: true
    remoteDebuggingVersion: VS2022
    runtimeScaleMonitoringEnabled: true
    websocketsEnabled: true
    healthCheckPath: /api/health
    healthCheckEvictionTimeInMin: 5
    workerCount: 4
    minimumTlsVersion: TLS_1_3
    scmMinimumTlsVersion: TLS_1_2
    cors:
      allowedOrigins:
        - https://orders.acme.example.com
        - https://admin.acme.example.com
      supportCredentials: true
    vnetRouteAllEnabled: true
  appSettings:
    ENVIRONMENT: production
    AAD_CLIENT_SECRET: set-me-via-key-vault-reference
    GITHUB_CLIENT_SECRET: set-me-via-key-vault-reference
  connectionStrings:
    - name: orders-db
      type: SQL_AZURE
      value:
        value: Server=tcp:orders.database.windows.net;Database=orders
    - name: cache
      type: REDIS_CACHE
      value:
        value: acme-cache.redis.cache.windows.net:6380,ssl=true
  stickySettings:
    appSettingNames:
      - ENVIRONMENT
    connectionStringNames:
      - orders-db
  applicationInsightsConnectionString:
    value: InstrumentationKey=00000000-0000-0000-0000-000000000000;IngestionEndpoint=https://centralus-2.in.applicationinsights.azure.com/
  httpsOnly: true
  publicNetworkAccessEnabled: true
  enabled: true
  clientCertificateEnabled: true
  clientCertificateMode: OPTIONAL_INTERACTIVE_USER
  clientCertificateExclusionPaths: /api/health;/api/status
  virtualNetworkSubnetId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/net-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet/subnets/flex-integration
  identity:
    type: SYSTEM_AND_USER_ASSIGNED
    identityIds:
      - value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/orders-mi
  webdeployPublishBasicAuthenticationEnabled: false
  authSettingsV2:
    authEnabled: true
    requireAuthentication: true
    unauthenticatedAction: RETURN_401
    defaultProvider: azureactivedirectory
    excludedPaths:
      - /api/health
    requireHttps: true
    login:
      tokenStoreEnabled: true
      tokenRefreshExtensionTime: 96
      cookieExpirationConvention: FIXED_TIME
      cookieExpirationTime: "08:00:00"
      validateNonce: true
      nonceExpirationTime: "00:05:00"
      allowedExternalRedirectUrls:
        - https://orders.acme.example.com/signed-in
    appleV2:
      clientId: com.acme.orders
      clientSecretSettingName: APPLE_CLIENT_SECRET
    activeDirectoryV2:
      clientId: 22222222-2222-2222-2222-222222222222
      tenantAuthEndpoint: https://login.microsoftonline.com/v2.0/33333333-3333-3333-3333-333333333333/
      clientSecretSettingName: AAD_CLIENT_SECRET
      allowedAudiences:
        - api://acme-orders
      allowedGroups:
        - 44444444-4444-4444-4444-444444444444
    azureStaticWebAppV2:
      clientId: 55555555-5555-5555-5555-555555555555
    customOidcV2:
      - name: corp-idp
        clientId: corp-orders-client
        openidConfigurationEndpoint: https://idp.acme.example.com/.well-known/openid-configuration
        nameClaimType: name
        scopes:
          - openid
          - profile
    facebookV2:
      appId: "666666666666666"
      appSecretSettingName: FACEBOOK_APP_SECRET
      graphApiVersion: v17.0
    githubV2:
      clientId: Iv1.acmeorders
      clientSecretSettingName: GITHUB_CLIENT_SECRET
    googleV2:
      clientId: acme-orders.apps.googleusercontent.com
      clientSecretSettingName: GOOGLE_CLIENT_SECRET
      loginScopes:
        - openid
    microsoftV2:
      clientId: 77777777-7777-7777-7777-777777777777
      clientSecretSettingName: MICROSOFT_CLIENT_SECRET
    twitterV2:
      consumerKey: acmeorderskey
      consumerSecretSettingName: TWITTER_CONSUMER_SECRET
  tags:
    cost-center: orders
    owner: platform-team
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.functionAppName` | `string` | yes |  |  |
| `spec.servicePlanId` | `string \| valueFrom` | yes |  | AzureServicePlan (`status.outputs.service_plan_id`) |
| `spec.storageContainerEndpoint` | `string` | yes |  |  |
| `spec.storageAuthenticationType` | `enum` | yes |  |  |
| `spec.storageAccessKey` | `string \| valueFrom` (sensitive) |  |  | AzureStorageAccount (`status.outputs.primary_access_key`) |
| `spec.storageUserAssignedIdentityId` | `string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |
| `spec.runtimeName` | `enum` | yes |  |  |
| `spec.runtimeVersion` | `string` | yes |  |  |
| `spec.instanceMemoryInMb` | `int32` |  | `2048` |  |
| `spec.maximumInstanceCount` | `int32` |  | `100` |  |
| `spec.httpConcurrency` | `int32` |  |  |  |
| `spec.alwaysReady` | `[]AzureFunctionAppFlexConsumptionAlwaysReady` |  |  |  |
| `spec.alwaysReady[].name` | `string` | yes |  |  |
| `spec.alwaysReady[].instanceCount` | `int32` |  |  |  |
| `spec.siteConfig` | `AzureFunctionAppFlexConsumptionSiteConfig` | yes |  |  |
| `spec.siteConfig.apiManagementApiId` | `string` |  |  |  |
| `spec.siteConfig.apiDefinitionUrl` | `string` |  |  |  |
| `spec.siteConfig.appCommandLine` | `string` |  |  |  |
| `spec.siteConfig.applicationInsightsKey` | `string` (sensitive) |  |  |  |
| `spec.siteConfig.appServiceLogs` | `AzureFunctionAppFlexConsumptionAppServiceLogs` |  |  |  |
| `spec.siteConfig.appServiceLogs.diskQuotaMb` | `int32` |  | `35` |  |
| `spec.siteConfig.appServiceLogs.retentionPeriodDays` | `int32` |  |  |  |
| `spec.siteConfig.containerRegistryUseManagedIdentity` | `bool` |  | `false` |  |
| `spec.siteConfig.defaultDocuments` | `[]string` |  |  |  |
| `spec.siteConfig.elasticInstanceMinimum` | `int32` |  |  |  |
| `spec.siteConfig.http2Enabled` | `bool` |  | `false` |  |
| `spec.siteConfig.ipRestrictions` | `[]AzureFunctionAppFlexConsumptionIpRestriction` |  |  |  |
| `spec.siteConfig.ipRestrictions[].name` | `string` |  |  |  |
| `spec.siteConfig.ipRestrictions[].priority` | `int32` |  |  |  |
| `spec.siteConfig.ipRestrictions[].action` | `enum` |  |  |  |
| `spec.siteConfig.ipRestrictions[].ipAddress` | `string` |  |  |  |
| `spec.siteConfig.ipRestrictions[].serviceTag` | `string` |  |  |  |
| `spec.siteConfig.ipRestrictions[].virtualNetworkSubnetId` | `string \| valueFrom` |  |  | AzureSubnet (`status.outputs.subnet_id`) |
| `spec.siteConfig.ipRestrictions[].description` | `string` |  |  |  |
| `spec.siteConfig.ipRestrictions[].headers` | `AzureFunctionAppFlexConsumptionIpRestrictionHeaders` |  |  |  |
| `spec.siteConfig.ipRestrictions[].headers.xForwardedFor` | `[]string` |  |  |  |
| `spec.siteConfig.ipRestrictions[].headers.xForwardedHost` | `[]string` |  |  |  |
| `spec.siteConfig.ipRestrictions[].headers.xAzureFdid` | `[]string \| valueFrom` |  |  | AzureFrontDoorProfile (`status.outputs.resource_guid`) |
| `spec.siteConfig.ipRestrictions[].headers.xFdHealthProbe` | `[]string` |  |  |  |
| `spec.siteConfig.ipRestrictionDefaultAction` | `enum` |  |  |  |
| `spec.siteConfig.scmUseMainIpRestriction` | `bool` |  | `false` |  |
| `spec.siteConfig.scmIpRestrictions` | `[]AzureFunctionAppFlexConsumptionIpRestriction` |  |  |  |
| `spec.siteConfig.scmIpRestrictions[].name` | `string` |  |  |  |
| `spec.siteConfig.scmIpRestrictions[].priority` | `int32` |  |  |  |
| `spec.siteConfig.scmIpRestrictions[].action` | `enum` |  |  |  |
| `spec.siteConfig.scmIpRestrictions[].ipAddress` | `string` |  |  |  |
| `spec.siteConfig.scmIpRestrictions[].serviceTag` | `string` |  |  |  |
| `spec.siteConfig.scmIpRestrictions[].virtualNetworkSubnetId` | `string \| valueFrom` |  |  | AzureSubnet (`status.outputs.subnet_id`) |
| `spec.siteConfig.scmIpRestrictions[].description` | `string` |  |  |  |
| `spec.siteConfig.scmIpRestrictions[].headers` | `AzureFunctionAppFlexConsumptionIpRestrictionHeaders` |  |  |  |
| `spec.siteConfig.scmIpRestrictions[].headers.xForwardedFor` | `[]string` |  |  |  |
| `spec.siteConfig.scmIpRestrictions[].headers.xForwardedHost` | `[]string` |  |  |  |
| `spec.siteConfig.scmIpRestrictions[].headers.xAzureFdid` | `[]string \| valueFrom` |  |  | AzureFrontDoorProfile (`status.outputs.resource_guid`) |
| `spec.siteConfig.scmIpRestrictions[].headers.xFdHealthProbe` | `[]string` |  |  |  |
| `spec.siteConfig.scmIpRestrictionDefaultAction` | `enum` |  |  |  |
| `spec.siteConfig.loadBalancingMode` | `enum` |  |  |  |
| `spec.siteConfig.managedPipelineMode` | `enum` |  |  |  |
| `spec.siteConfig.remoteDebuggingEnabled` | `bool` |  | `false` |  |
| `spec.siteConfig.remoteDebuggingVersion` | `string` |  |  |  |
| `spec.siteConfig.runtimeScaleMonitoringEnabled` | `bool` |  |  |  |
| `spec.siteConfig.websocketsEnabled` | `bool` |  | `false` |  |
| `spec.siteConfig.healthCheckPath` | `string` |  |  |  |
| `spec.siteConfig.healthCheckEvictionTimeInMin` | `int32` |  |  |  |
| `spec.siteConfig.workerCount` | `int32` |  |  |  |
| `spec.siteConfig.minimumTlsVersion` | `enum` |  |  |  |
| `spec.siteConfig.scmMinimumTlsVersion` | `enum` |  |  |  |
| `spec.siteConfig.cors` | `AzureFunctionAppFlexConsumptionCorsSettings` |  |  |  |
| `spec.siteConfig.cors.allowedOrigins` | `[]string` | yes |  |  |
| `spec.siteConfig.cors.supportCredentials` | `bool` |  | `false` |  |
| `spec.siteConfig.vnetRouteAllEnabled` | `bool` |  | `false` |  |
| `spec.appSettings` | `map<string, string>` |  |  |  |
| `spec.connectionStrings` | `[]AzureFunctionAppFlexConsumptionConnectionString` |  |  |  |
| `spec.connectionStrings[].name` | `string` | yes |  |  |
| `spec.connectionStrings[].type` | `enum` | yes |  |  |
| `spec.connectionStrings[].value` | `string \| valueFrom` (sensitive) | yes |  |  |
| `spec.stickySettings` | `AzureFunctionAppFlexConsumptionStickySettings` |  |  |  |
| `spec.stickySettings.appSettingNames` | `[]string` |  |  |  |
| `spec.stickySettings.connectionStringNames` | `[]string` |  |  |  |
| `spec.applicationInsightsConnectionString` | `string \| valueFrom` |  |  | AzureApplicationInsights (`status.outputs.connection_string`) |
| `spec.httpsOnly` | `bool` |  | `true` |  |
| `spec.publicNetworkAccessEnabled` | `bool` |  | `true` |  |
| `spec.enabled` | `bool` |  | `true` |  |
| `spec.clientCertificateEnabled` | `bool` |  | `false` |  |
| `spec.clientCertificateMode` | `enum` |  |  |  |
| `spec.clientCertificateExclusionPaths` | `string` |  |  |  |
| `spec.virtualNetworkSubnetId` | `string \| valueFrom` |  |  | AzureSubnet (`status.outputs.subnet_id`) |
| `spec.identity` | `AzureFunctionAppFlexConsumptionIdentity` |  |  |  |
| `spec.identity.type` | `enum` | yes |  |  |
| `spec.identity.identityIds` | `[]string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |
| `spec.webdeployPublishBasicAuthenticationEnabled` | `bool` |  | `true` |  |
| `spec.zipDeployFile` | `string` |  |  |  |
| `spec.authSettingsV2` | `AzureFunctionAppFlexConsumptionAuthSettingsV2` |  |  |  |
| `spec.authSettingsV2.authEnabled` | `bool` |  | `false` |  |
| `spec.authSettingsV2.runtimeVersion` | `string` |  | `~1` |  |
| `spec.authSettingsV2.configFilePath` | `string` |  |  |  |
| `spec.authSettingsV2.requireAuthentication` | `bool` |  | `false` |  |
| `spec.authSettingsV2.unauthenticatedAction` | `enum` |  |  |  |
| `spec.authSettingsV2.defaultProvider` | `string` |  |  |  |
| `spec.authSettingsV2.excludedPaths` | `[]string` |  |  |  |
| `spec.authSettingsV2.requireHttps` | `bool` |  | `true` |  |
| `spec.authSettingsV2.httpRouteApiPrefix` | `string` |  | `/.auth` |  |
| `spec.authSettingsV2.forwardProxyConvention` | `enum` |  |  |  |
| `spec.authSettingsV2.forwardProxyCustomHostHeaderName` | `string` |  |  |  |
| `spec.authSettingsV2.forwardProxyCustomSchemeHeaderName` | `string` |  |  |  |
| `spec.authSettingsV2.login` | `AzureFunctionAppFlexConsumptionAuthV2Login` | yes |  |  |
| `spec.authSettingsV2.login.logoutEndpoint` | `string` |  |  |  |
| `spec.authSettingsV2.login.tokenStoreEnabled` | `bool` |  | `false` |  |
| `spec.authSettingsV2.login.tokenRefreshExtensionTime` | `double` |  | `72` |  |
| `spec.authSettingsV2.login.tokenStorePath` | `string` |  |  |  |
| `spec.authSettingsV2.login.tokenStoreSasSettingName` | `string` |  |  |  |
| `spec.authSettingsV2.login.preserveUrlFragmentsForLogins` | `bool` |  | `false` |  |
| `spec.authSettingsV2.login.allowedExternalRedirectUrls` | `[]string` |  |  |  |
| `spec.authSettingsV2.login.cookieExpirationConvention` | `enum` |  |  |  |
| `spec.authSettingsV2.login.cookieExpirationTime` | `string` |  | `08:00:00` |  |
| `spec.authSettingsV2.login.validateNonce` | `bool` |  | `true` |  |
| `spec.authSettingsV2.login.nonceExpirationTime` | `string` |  | `00:05:00` |  |
| `spec.authSettingsV2.appleV2` | `AzureFunctionAppFlexConsumptionAuthV2Apple` |  |  |  |
| `spec.authSettingsV2.appleV2.clientId` | `string` | yes |  |  |
| `spec.authSettingsV2.appleV2.clientSecretSettingName` | `string` | yes |  |  |
| `spec.authSettingsV2.activeDirectoryV2` | `AzureFunctionAppFlexConsumptionAuthV2ActiveDirectory` |  |  |  |
| `spec.authSettingsV2.activeDirectoryV2.clientId` | `string` | yes |  |  |
| `spec.authSettingsV2.activeDirectoryV2.tenantAuthEndpoint` | `string` | yes |  |  |
| `spec.authSettingsV2.activeDirectoryV2.clientSecretSettingName` | `string` |  |  |  |
| `spec.authSettingsV2.activeDirectoryV2.clientSecretCertificateThumbprint` | `string` |  |  |  |
| `spec.authSettingsV2.activeDirectoryV2.loginParameters` | `map<string, string>` |  |  |  |
| `spec.authSettingsV2.activeDirectoryV2.wwwAuthenticationDisabled` | `bool` |  | `false` |  |
| `spec.authSettingsV2.activeDirectoryV2.jwtAllowedGroups` | `[]string` |  |  |  |
| `spec.authSettingsV2.activeDirectoryV2.jwtAllowedClientApplications` | `[]string` |  |  |  |
| `spec.authSettingsV2.activeDirectoryV2.allowedGroups` | `[]string` |  |  |  |
| `spec.authSettingsV2.activeDirectoryV2.allowedIdentities` | `[]string` |  |  |  |
| `spec.authSettingsV2.activeDirectoryV2.allowedApplications` | `[]string` |  |  |  |
| `spec.authSettingsV2.activeDirectoryV2.allowedAudiences` | `[]string` |  |  |  |
| `spec.authSettingsV2.azureStaticWebAppV2` | `AzureFunctionAppFlexConsumptionAuthV2StaticWebApp` |  |  |  |
| `spec.authSettingsV2.azureStaticWebAppV2.clientId` | `string` | yes |  |  |
| `spec.authSettingsV2.customOidcV2` | `[]AzureFunctionAppFlexConsumptionAuthV2CustomOidc` |  |  |  |
| `spec.authSettingsV2.customOidcV2[].name` | `string` | yes |  |  |
| `spec.authSettingsV2.customOidcV2[].clientId` | `string` | yes |  |  |
| `spec.authSettingsV2.customOidcV2[].openidConfigurationEndpoint` | `string` | yes |  |  |
| `spec.authSettingsV2.customOidcV2[].nameClaimType` | `string` |  |  |  |
| `spec.authSettingsV2.customOidcV2[].scopes` | `[]string` |  |  |  |
| `spec.authSettingsV2.facebookV2` | `AzureFunctionAppFlexConsumptionAuthV2Facebook` |  |  |  |
| `spec.authSettingsV2.facebookV2.appId` | `string` | yes |  |  |
| `spec.authSettingsV2.facebookV2.appSecretSettingName` | `string` | yes |  |  |
| `spec.authSettingsV2.facebookV2.graphApiVersion` | `string` |  |  |  |
| `spec.authSettingsV2.facebookV2.loginScopes` | `[]string` |  |  |  |
| `spec.authSettingsV2.githubV2` | `AzureFunctionAppFlexConsumptionAuthV2Github` |  |  |  |
| `spec.authSettingsV2.githubV2.clientId` | `string` | yes |  |  |
| `spec.authSettingsV2.githubV2.clientSecretSettingName` | `string` | yes |  |  |
| `spec.authSettingsV2.githubV2.loginScopes` | `[]string` |  |  |  |
| `spec.authSettingsV2.googleV2` | `AzureFunctionAppFlexConsumptionAuthV2Google` |  |  |  |
| `spec.authSettingsV2.googleV2.clientId` | `string` | yes |  |  |
| `spec.authSettingsV2.googleV2.clientSecretSettingName` | `string` | yes |  |  |
| `spec.authSettingsV2.googleV2.allowedAudiences` | `[]string` |  |  |  |
| `spec.authSettingsV2.googleV2.loginScopes` | `[]string` |  |  |  |
| `spec.authSettingsV2.microsoftV2` | `AzureFunctionAppFlexConsumptionAuthV2Microsoft` |  |  |  |
| `spec.authSettingsV2.microsoftV2.clientId` | `string` | yes |  |  |
| `spec.authSettingsV2.microsoftV2.clientSecretSettingName` | `string` | yes |  |  |
| `spec.authSettingsV2.microsoftV2.allowedAudiences` | `[]string` |  |  |  |
| `spec.authSettingsV2.microsoftV2.loginScopes` | `[]string` |  |  |  |
| `spec.authSettingsV2.twitterV2` | `AzureFunctionAppFlexConsumptionAuthV2Twitter` |  |  |  |
| `spec.authSettingsV2.twitterV2.consumerKey` | `string` | yes |  |  |
| `spec.authSettingsV2.twitterV2.consumerSecretSettingName` | `string` | yes |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.resourceGroup

`string | valueFrom` · required

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.functionAppName

`string` · required

- rule: function_app_name may only contain alphanumeric characters and dashes and must be 1 to 60 characters long
- rule: {"required":true}

### spec.servicePlanId

`string | valueFrom` · required

- references: AzureServicePlan (`status.outputs.service_plan_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureServicePlan, name: <that resource's name>, fieldPath: status.outputs.service_plan_id}} -- a bare string does not parse

### spec.storageContainerEndpoint

`string` · required

- rule: storage_container_endpoint must be an https blob container URL (https://{account}.blob.core.windows.net/{container})
- rule: {"required":true}

### spec.storageAuthenticationType

`enum` · required

- rule: {"required":true,"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_function_app_flex_consumption_storage_authentication_type_unspecified`
- `STORAGE_ACCOUNT_CONNECTION_STRING`
- `SYSTEM_ASSIGNED_IDENTITY`
- `USER_ASSIGNED_IDENTITY`

### spec.storageAccessKey

`string | valueFrom` · sensitive

- references: AzureStorageAccount (`status.outputs.primary_access_key`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.primary_access_key}} -- a bare string does not parse

### spec.storageUserAssignedIdentityId

`string | valueFrom`

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.runtimeName

`enum` · required

- rule: {"required":true,"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_function_app_flex_consumption_runtime_name_unspecified`
- `NODE`
- `DOTNET_ISOLATED`
- `JAVA`
- `POWERSHELL`
- `PYTHON`
- `CUSTOM_HANDLER`

### spec.runtimeVersion

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.instanceMemoryInMb

`int32` · optional (explicit presence)

- default: `2048`

### spec.maximumInstanceCount

`int32` · optional (explicit presence)

- default: `100`
- rule: {"int32":{"lte":1000,"gte":1}}

### spec.httpConcurrency

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":1000,"gte":1}}

### spec.alwaysReady

`[]AzureFunctionAppFlexConsumptionAlwaysReady`

### spec.alwaysReady[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.alwaysReady[].instanceCount

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":1000,"gte":0}}

### spec.siteConfig

`AzureFunctionAppFlexConsumptionSiteConfig` · required

- rule: {"required":true}
- rule: health_check_path and health_check_eviction_time_in_min require each other (Azure pairs them both ways)

### spec.siteConfig.apiManagementApiId

`string`

### spec.siteConfig.apiDefinitionUrl

`string`

### spec.siteConfig.appCommandLine

`string`

### spec.siteConfig.applicationInsightsKey

`string` · sensitive

### spec.siteConfig.appServiceLogs

`AzureFunctionAppFlexConsumptionAppServiceLogs`

### spec.siteConfig.appServiceLogs.diskQuotaMb

`int32` · optional (explicit presence)

- default: `35`
- rule: {"int32":{"lte":100,"gte":25}}

### spec.siteConfig.appServiceLogs.retentionPeriodDays

`int32` · optional (explicit presence)

- rule: {"int32":{"gte":0}}

### spec.siteConfig.containerRegistryUseManagedIdentity

`bool` · optional (explicit presence)

- default: `false`

### spec.siteConfig.defaultDocuments

`[]string`

### spec.siteConfig.elasticInstanceMinimum

`int32` · optional (explicit presence)

- rule: {"int32":{"gte":0}}

### spec.siteConfig.http2Enabled

`bool` · optional (explicit presence)

- default: `false`

### spec.siteConfig.ipRestrictions

`[]AzureFunctionAppFlexConsumptionIpRestriction`

### spec.siteConfig.ipRestrictions[].name

`string`

### spec.siteConfig.ipRestrictions[].priority

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":65000,"gte":1}}

### spec.siteConfig.ipRestrictions[].action

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_function_app_flex_consumption_ip_restriction_action_unspecified`
- `ALLOW`
- `DENY`

### spec.siteConfig.ipRestrictions[].ipAddress

`string`

### spec.siteConfig.ipRestrictions[].serviceTag

`string`

### spec.siteConfig.ipRestrictions[].virtualNetworkSubnetId

`string | valueFrom`

- references: AzureSubnet (`status.outputs.subnet_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.siteConfig.ipRestrictions[].description

`string`

### spec.siteConfig.ipRestrictions[].headers

`AzureFunctionAppFlexConsumptionIpRestrictionHeaders`

### spec.siteConfig.ipRestrictions[].headers.xForwardedFor

`[]string`

- rule: {"repeated":{"maxItems":"8"}}

### spec.siteConfig.ipRestrictions[].headers.xForwardedHost

`[]string`

- rule: {"repeated":{"maxItems":"8"}}

### spec.siteConfig.ipRestrictions[].headers.xAzureFdid

`[]string | valueFrom`

- references: AzureFrontDoorProfile (`status.outputs.resource_guid`)
- rule: {"repeated":{"maxItems":"8"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureFrontDoorProfile, name: <that resource's name>, fieldPath: status.outputs.resource_guid}} -- a bare string does not parse

### spec.siteConfig.ipRestrictions[].headers.xFdHealthProbe

`[]string`

- rule: {"repeated":{"maxItems":"1"}}

### spec.siteConfig.ipRestrictionDefaultAction

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_function_app_flex_consumption_ip_restriction_action_unspecified`
- `ALLOW`
- `DENY`

### spec.siteConfig.scmUseMainIpRestriction

`bool` · optional (explicit presence)

- default: `false`

### spec.siteConfig.scmIpRestrictions

`[]AzureFunctionAppFlexConsumptionIpRestriction`

### spec.siteConfig.scmIpRestrictions[].name

`string`

### spec.siteConfig.scmIpRestrictions[].priority

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":65000,"gte":1}}

### spec.siteConfig.scmIpRestrictions[].action

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_function_app_flex_consumption_ip_restriction_action_unspecified`
- `ALLOW`
- `DENY`

### spec.siteConfig.scmIpRestrictions[].ipAddress

`string`

### spec.siteConfig.scmIpRestrictions[].serviceTag

`string`

### spec.siteConfig.scmIpRestrictions[].virtualNetworkSubnetId

`string | valueFrom`

- references: AzureSubnet (`status.outputs.subnet_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.siteConfig.scmIpRestrictions[].description

`string`

### spec.siteConfig.scmIpRestrictions[].headers

`AzureFunctionAppFlexConsumptionIpRestrictionHeaders`

### spec.siteConfig.scmIpRestrictions[].headers.xForwardedFor

`[]string`

- rule: {"repeated":{"maxItems":"8"}}

### spec.siteConfig.scmIpRestrictions[].headers.xForwardedHost

`[]string`

- rule: {"repeated":{"maxItems":"8"}}

### spec.siteConfig.scmIpRestrictions[].headers.xAzureFdid

`[]string | valueFrom`

- references: AzureFrontDoorProfile (`status.outputs.resource_guid`)
- rule: {"repeated":{"maxItems":"8"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureFrontDoorProfile, name: <that resource's name>, fieldPath: status.outputs.resource_guid}} -- a bare string does not parse

### spec.siteConfig.scmIpRestrictions[].headers.xFdHealthProbe

`[]string`

- rule: {"repeated":{"maxItems":"1"}}

### spec.siteConfig.scmIpRestrictionDefaultAction

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_function_app_flex_consumption_ip_restriction_action_unspecified`
- `ALLOW`
- `DENY`

### spec.siteConfig.loadBalancingMode

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_function_app_flex_consumption_load_balancing_mode_unspecified`
- `LEAST_REQUESTS`
- `WEIGHTED_ROUND_ROBIN`
- `LEAST_RESPONSE_TIME`
- `WEIGHTED_TOTAL_TRAFFIC`
- `REQUEST_HASH`
- `PER_SITE_ROUND_ROBIN`

### spec.siteConfig.managedPipelineMode

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_function_app_flex_consumption_managed_pipeline_mode_unspecified`
- `INTEGRATED`
- `CLASSIC`

### spec.siteConfig.remoteDebuggingEnabled

`bool` · optional (explicit presence)

- default: `false`

### spec.siteConfig.remoteDebuggingVersion

`string`

- rule: remote_debugging_version must be one of: VS2017, VS2019, VS2022

### spec.siteConfig.runtimeScaleMonitoringEnabled

`bool` · optional (explicit presence)

### spec.siteConfig.websocketsEnabled

`bool` · optional (explicit presence)

- default: `false`

### spec.siteConfig.healthCheckPath

`string`

### spec.siteConfig.healthCheckEvictionTimeInMin

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":10,"gte":2}}

### spec.siteConfig.workerCount

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":100,"gte":1}}

### spec.siteConfig.minimumTlsVersion

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_function_app_flex_consumption_tls_version_unspecified`
- `TLS_1_0`
- `TLS_1_1`
- `TLS_1_2`
- `TLS_1_3`

### spec.siteConfig.scmMinimumTlsVersion

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_function_app_flex_consumption_tls_version_unspecified`
- `TLS_1_0`
- `TLS_1_1`
- `TLS_1_2`
- `TLS_1_3`

### spec.siteConfig.cors

`AzureFunctionAppFlexConsumptionCorsSettings`

- rule: support_credentials cannot be used with a wildcard '*' origin

### spec.siteConfig.cors.allowedOrigins

`[]string` · required

- rule: {"repeated":{"minItems":"1"}}

### spec.siteConfig.cors.supportCredentials

`bool` · optional (explicit presence)

- default: `false`

### spec.siteConfig.vnetRouteAllEnabled

`bool` · optional (explicit presence)

- default: `false`

### spec.appSettings

`map<string, string>`

### spec.connectionStrings

`[]AzureFunctionAppFlexConsumptionConnectionString`

### spec.connectionStrings[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.connectionStrings[].type

`enum` · required

- rule: {"required":true,"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_function_app_flex_consumption_connection_string_type_unspecified`
- `MYSQL`
- `SQL_SERVER`
- `SQL_AZURE`
- `CUSTOM`
- `NOTIFICATION_HUB`
- `SERVICE_BUS`
- `EVENT_HUB`
- `API_HUB`
- `DOC_DB`
- `REDIS_CACHE`
- `POSTGRESQL`

### spec.connectionStrings[].value

`string | valueFrom` · required · sensitive

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.stickySettings

`AzureFunctionAppFlexConsumptionStickySettings`

- rule: sticky_settings requires at least one app_setting_names or connection_string_names entry

### spec.stickySettings.appSettingNames

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.stickySettings.connectionStringNames

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.applicationInsightsConnectionString

`string | valueFrom`

- references: AzureApplicationInsights (`status.outputs.connection_string`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureApplicationInsights, name: <that resource's name>, fieldPath: status.outputs.connection_string}} -- a bare string does not parse

### spec.httpsOnly

`bool` · optional (explicit presence)

- default: `true`

### spec.publicNetworkAccessEnabled

`bool` · optional (explicit presence)

- default: `true`

### spec.enabled

`bool` · optional (explicit presence)

- default: `true`

### spec.clientCertificateEnabled

`bool` · optional (explicit presence)

- default: `false`

### spec.clientCertificateMode

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_function_app_flex_consumption_client_certificate_mode_unspecified`
- `REQUIRED`
- `OPTIONAL`
- `OPTIONAL_INTERACTIVE_USER`

### spec.clientCertificateExclusionPaths

`string`

### spec.virtualNetworkSubnetId

`string | valueFrom`

- references: AzureSubnet (`status.outputs.subnet_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.identity

`AzureFunctionAppFlexConsumptionIdentity`

- rule: identity_ids is required when type includes USER_ASSIGNED, and must be empty for SYSTEM_ASSIGNED

### spec.identity.type

`enum` · required

- rule: {"required":true,"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_function_app_flex_consumption_identity_type_unspecified`
- `SYSTEM_ASSIGNED`
- `USER_ASSIGNED`
- `SYSTEM_AND_USER_ASSIGNED`

### spec.identity.identityIds

`[]string | valueFrom`

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.webdeployPublishBasicAuthenticationEnabled

`bool` · optional (explicit presence)

- default: `true`

### spec.zipDeployFile

`string`

### spec.authSettingsV2

`AzureFunctionAppFlexConsumptionAuthSettingsV2`

- rule: forward_proxy_custom_host_header_name and forward_proxy_custom_scheme_header_name require forward_proxy_convention FORWARD_PROXY_CUSTOM

### spec.authSettingsV2.authEnabled

`bool` · optional (explicit presence)

- default: `false`

### spec.authSettingsV2.runtimeVersion

`string` · optional (explicit presence)

- default: `~1`

### spec.authSettingsV2.configFilePath

`string`

### spec.authSettingsV2.requireAuthentication

`bool` · optional (explicit presence)

- default: `false`

### spec.authSettingsV2.unauthenticatedAction

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_function_app_flex_consumption_unauthenticated_action_unspecified`
- `REDIRECT_TO_LOGIN_PAGE`
- `ALLOW_ANONYMOUS`
- `RETURN_401`
- `RETURN_403`

### spec.authSettingsV2.defaultProvider

`string`

### spec.authSettingsV2.excludedPaths

`[]string`

### spec.authSettingsV2.requireHttps

`bool` · optional (explicit presence)

- default: `true`

### spec.authSettingsV2.httpRouteApiPrefix

`string` · optional (explicit presence)

- default: `/.auth`

### spec.authSettingsV2.forwardProxyConvention

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_function_app_flex_consumption_forward_proxy_convention_unspecified`
- `FORWARD_PROXY_NO_PROXY`
- `FORWARD_PROXY_STANDARD`
- `FORWARD_PROXY_CUSTOM`

### spec.authSettingsV2.forwardProxyCustomHostHeaderName

`string`

### spec.authSettingsV2.forwardProxyCustomSchemeHeaderName

`string`

### spec.authSettingsV2.login

`AzureFunctionAppFlexConsumptionAuthV2Login` · required

- rule: {"required":true}
- rule: token_store_path and token_store_sas_setting_name are mutually exclusive

### spec.authSettingsV2.login.logoutEndpoint

`string`

### spec.authSettingsV2.login.tokenStoreEnabled

`bool` · optional (explicit presence)

- default: `false`

### spec.authSettingsV2.login.tokenRefreshExtensionTime

`double` · optional (explicit presence)

- default: `72`
- rule: {"double":{"gte":0}}

### spec.authSettingsV2.login.tokenStorePath

`string`

### spec.authSettingsV2.login.tokenStoreSasSettingName

`string`

### spec.authSettingsV2.login.preserveUrlFragmentsForLogins

`bool` · optional (explicit presence)

- default: `false`

### spec.authSettingsV2.login.allowedExternalRedirectUrls

`[]string`

### spec.authSettingsV2.login.cookieExpirationConvention

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_function_app_flex_consumption_cookie_expiration_convention_unspecified`
- `FIXED_TIME`
- `IDENTITY_PROVIDER_DERIVED`

### spec.authSettingsV2.login.cookieExpirationTime

`string` · optional (explicit presence)

- default: `08:00:00`
- rule: cookie_expiration_time must be hh:mm:ss (e.g. 08:00:00)

### spec.authSettingsV2.login.validateNonce

`bool` · optional (explicit presence)

- default: `true`

### spec.authSettingsV2.login.nonceExpirationTime

`string` · optional (explicit presence)

- default: `00:05:00`
- rule: nonce_expiration_time must be hh:mm:ss (e.g. 00:05:00)

### spec.authSettingsV2.appleV2

`AzureFunctionAppFlexConsumptionAuthV2Apple`

### spec.authSettingsV2.appleV2.clientId

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.authSettingsV2.appleV2.clientSecretSettingName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.authSettingsV2.activeDirectoryV2

`AzureFunctionAppFlexConsumptionAuthV2ActiveDirectory`

- rule: client_secret_setting_name and client_secret_certificate_thumbprint are mutually exclusive

### spec.authSettingsV2.activeDirectoryV2.clientId

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.authSettingsV2.activeDirectoryV2.tenantAuthEndpoint

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.authSettingsV2.activeDirectoryV2.clientSecretSettingName

`string`

### spec.authSettingsV2.activeDirectoryV2.clientSecretCertificateThumbprint

`string`

### spec.authSettingsV2.activeDirectoryV2.loginParameters

`map<string, string>`

### spec.authSettingsV2.activeDirectoryV2.wwwAuthenticationDisabled

`bool` · optional (explicit presence)

- default: `false`

### spec.authSettingsV2.activeDirectoryV2.jwtAllowedGroups

`[]string`

### spec.authSettingsV2.activeDirectoryV2.jwtAllowedClientApplications

`[]string`

### spec.authSettingsV2.activeDirectoryV2.allowedGroups

`[]string`

### spec.authSettingsV2.activeDirectoryV2.allowedIdentities

`[]string`

### spec.authSettingsV2.activeDirectoryV2.allowedApplications

`[]string`

### spec.authSettingsV2.activeDirectoryV2.allowedAudiences

`[]string`

### spec.authSettingsV2.azureStaticWebAppV2

`AzureFunctionAppFlexConsumptionAuthV2StaticWebApp`

### spec.authSettingsV2.azureStaticWebAppV2.clientId

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.authSettingsV2.customOidcV2

`[]AzureFunctionAppFlexConsumptionAuthV2CustomOidc`

### spec.authSettingsV2.customOidcV2[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.authSettingsV2.customOidcV2[].clientId

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.authSettingsV2.customOidcV2[].openidConfigurationEndpoint

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.authSettingsV2.customOidcV2[].nameClaimType

`string`

### spec.authSettingsV2.customOidcV2[].scopes

`[]string`

### spec.authSettingsV2.facebookV2

`AzureFunctionAppFlexConsumptionAuthV2Facebook`

### spec.authSettingsV2.facebookV2.appId

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.authSettingsV2.facebookV2.appSecretSettingName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.authSettingsV2.facebookV2.graphApiVersion

`string`

### spec.authSettingsV2.facebookV2.loginScopes

`[]string`

### spec.authSettingsV2.githubV2

`AzureFunctionAppFlexConsumptionAuthV2Github`

### spec.authSettingsV2.githubV2.clientId

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.authSettingsV2.githubV2.clientSecretSettingName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.authSettingsV2.githubV2.loginScopes

`[]string`

### spec.authSettingsV2.googleV2

`AzureFunctionAppFlexConsumptionAuthV2Google`

### spec.authSettingsV2.googleV2.clientId

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.authSettingsV2.googleV2.clientSecretSettingName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.authSettingsV2.googleV2.allowedAudiences

`[]string`

### spec.authSettingsV2.googleV2.loginScopes

`[]string`

### spec.authSettingsV2.microsoftV2

`AzureFunctionAppFlexConsumptionAuthV2Microsoft`

### spec.authSettingsV2.microsoftV2.clientId

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.authSettingsV2.microsoftV2.clientSecretSettingName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.authSettingsV2.microsoftV2.allowedAudiences

`[]string`

### spec.authSettingsV2.microsoftV2.loginScopes

`[]string`

### spec.authSettingsV2.twitterV2

`AzureFunctionAppFlexConsumptionAuthV2Twitter`

### spec.authSettingsV2.twitterV2.consumerKey

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.authSettingsV2.twitterV2.consumerSecretSettingName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.tags

`map<string, string>`

## Validation Rules

- `flex_storage_key_required_for_connection_string`: storage_access_key is required when storage_authentication_type is STORAGE_ACCOUNT_CONNECTION_STRING
- `flex_storage_uami_required_for_user_assigned`: storage_user_assigned_identity_id is required when storage_authentication_type is USER_ASSIGNED_IDENTITY
- `flex_vnet_route_all_requires_subnet`: site_config.vnet_route_all_enabled requires virtual_network_subnet_id (there is no VNet to route through otherwise)

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureFunctionAppFlexConsumption, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.function_app_id` | `string` |  |
| `status.outputs.default_hostname` | `string` |  |
| `status.outputs.outbound_ip_addresses` | `[]string` |  |
| `status.outputs.identity_principal_id` | `string` |  |
| `status.outputs.identity_tenant_id` | `string` |  |
| `status.outputs.custom_domain_verification_id` | `string` |  |
| `status.outputs.kind` | `string` |  |
| `status.outputs.possible_outbound_ip_addresses` | `[]string` |  |
| `status.outputs.site_credential_name` | `string` |  |
| `status.outputs.site_credential_password` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.servicePlanId` | AzureServicePlan | `status.outputs.service_plan_id` |
| `spec.storageAccessKey` | AzureStorageAccount | `status.outputs.primary_access_key` |
| `spec.storageUserAssignedIdentityId` | AzureUserAssignedIdentity | `status.outputs.identity_id` |
| `spec.siteConfig.ipRestrictions[].virtualNetworkSubnetId` | AzureSubnet | `status.outputs.subnet_id` |
| `spec.siteConfig.ipRestrictions[].headers.xAzureFdid` | AzureFrontDoorProfile | `status.outputs.resource_guid` |
| `spec.siteConfig.scmIpRestrictions[].virtualNetworkSubnetId` | AzureSubnet | `status.outputs.subnet_id` |
| `spec.siteConfig.scmIpRestrictions[].headers.xAzureFdid` | AzureFrontDoorProfile | `status.outputs.resource_guid` |
| `spec.applicationInsightsConnectionString` | AzureApplicationInsights | `status.outputs.connection_string` |
| `spec.virtualNetworkSubnetId` | AzureSubnet | `status.outputs.subnet_id` |
| `spec.identity.identityIds` | AzureUserAssignedIdentity | `status.outputs.identity_id` |

## See Also

- [Overview](../README.md)

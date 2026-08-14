# AzureContainerInstance

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Deep-shape example for docs and offline validation: a maximal PUBLIC
# container group exercising every arm -- two containers (one with all
# four volume forms across the pair, probes in both forms, limits,
# secure environment variables), an init container seeding a shared
# scratch volume, both registry-credential forms, combined identity,
# Log Analytics diagnostics with typed metadata, custom DNS, zones,
# narrowed exposed ports, and customer-managed-key encryption.
# References are literal values so the manifest validates standalone.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureContainerInstance
metadata:
  name: test-container-instance
  id: test-container-instance
  org: test-org
  env: test
spec:
  resourceGroup:
    value: test-rg
  name: app-aci
  region: eastus
  osType: Linux
  restartPolicy: Always
  sku: Standard
  ipAddressType: Public
  dnsNameLabel: app-aci-acme
  dnsNameLabelReusePolicy: TenantReuse
  zones:
    - "1"
  # Narrow the group's surface to the web port -- the sidecar's metrics
  # port stays group-internal.
  exposedPorts:
    - port: 80
  initContainers:
    - name: seed
      image: mcr.microsoft.com/cbl-mariner/busybox:2.0
      commands: ["sh", "-c", "echo ready > /work/ready"]
      environmentVariables:
        SEED_MODE: static
      volumes:
        - name: work
          mountPath: /work
          emptyDir: true
  containers:
    - name: web
      image: mcr.microsoft.com/azuredocs/aci-helloworld:latest
      cpu: 0.5
      memory: 1
      cpuLimit: 1
      memoryLimit: 1.5
      ports:
        - port: 80
      environmentVariables:
        APP_MODE: production
      secureEnvironmentVariables:
        API_TOKEN: test-token-value
      volumes:
        - name: work
          mountPath: /work
          emptyDir: true
        - name: content
          mountPath: /content
          readOnly: true
          azureFile:
            shareName:
              value: app-content
            storageAccountName:
              value: appcontentsa
            storageAccountKey:
              value: dGVzdC1zdG9yYWdlLWtleQ==
      livenessProbe:
        httpGet:
          path: /
          port: 80
          scheme: http
        periodSeconds: 30
        failureThreshold: 3
      readinessProbe:
        exec: ["ls", "/work/ready"]
        initialDelaySeconds: 5
        successThreshold: 1
        timeoutSeconds: 2
    - name: sidecar
      image: registry.example.com/acme/metrics-sidecar:2.1.0
      cpu: 0.25
      memory: 0.5
      ports:
        - port: 9090
      commands: ["/bin/sidecar", "--listen", ":9090"]
      security:
        privilegeEnabled: false
      volumes:
        - name: rules
          mountPath: /etc/rules
          gitRepo:
            url: https://github.com/acme/metric-rules
            directory: rules
            revision: v1.4.0
        - name: creds
          mountPath: /etc/creds
          secret:
            client.pem: dGVzdC1jbGllbnQtcGVt
  imageRegistryCredentials:
    - server:
        value: acme.azurecr.io
      userAssignedIdentityId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/app-mi
    - server:
        value: registry.example.com
      username: puller
      password: test-registry-password
  identity:
    type: SYSTEM_AND_USER_ASSIGNED
    identityIds:
      - value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/app-mi
  diagnosticsLogAnalytics:
    workspaceId:
      value: 00000000-0000-0000-0000-000000000000
    workspaceKey:
      value: dGVzdC13b3Jrc3BhY2Uta2V5
    logType: ContainerInsights
    metadata:
      team: platform
  dnsConfig:
    nameservers:
      - 10.0.0.10
    searchDomains:
      - internal.acme.com
    options:
      - ndots:2
  keyVaultKeyId:
    value: https://app-vault.vault.azure.net/keys/aci-cmk/0123456789abcdef0123456789abcdef
  keyVaultUserAssignedIdentityId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/app-mi
  tags:
    workload: web
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.region` | `string` | yes |  |  |
| `spec.osType` | `string` | yes |  |  |
| `spec.restartPolicy` | `string` |  |  |  |
| `spec.sku` | `string` |  |  |  |
| `spec.priority` | `string` |  |  |  |
| `spec.ipAddressType` | `string` |  |  |  |
| `spec.dnsNameLabel` | `string` |  |  |  |
| `spec.dnsNameLabelReusePolicy` | `string` |  |  |  |
| `spec.subnetId` | `string \| valueFrom` |  |  | AzureSubnet (`status.outputs.subnet_id`) |
| `spec.zones` | `[]string` |  |  |  |
| `spec.exposedPorts` | `[]AzureContainerInstancePort` |  |  |  |
| `spec.exposedPorts[].port` | `int32` |  |  |  |
| `spec.exposedPorts[].protocol` | `string` |  |  |  |
| `spec.containers` | `[]AzureContainerInstanceContainer` | yes |  |  |
| `spec.containers[].name` | `string` | yes |  |  |
| `spec.containers[].image` | `string` | yes |  |  |
| `spec.containers[].cpu` | `double` | yes |  |  |
| `spec.containers[].memory` | `double` | yes |  |  |
| `spec.containers[].cpuLimit` | `double` |  |  |  |
| `spec.containers[].memoryLimit` | `double` |  |  |  |
| `spec.containers[].ports` | `[]AzureContainerInstancePort` |  |  |  |
| `spec.containers[].ports[].port` | `int32` |  |  |  |
| `spec.containers[].ports[].protocol` | `string` |  |  |  |
| `spec.containers[].environmentVariables` | `map<string, string>` |  |  |  |
| `spec.containers[].secureEnvironmentVariables` | `map<string, string>` (sensitive) |  |  |  |
| `spec.containers[].commands` | `[]string` |  |  |  |
| `spec.containers[].volumes` | `[]AzureContainerInstanceVolume` |  |  |  |
| `spec.containers[].volumes[].name` | `string` | yes |  |  |
| `spec.containers[].volumes[].mountPath` | `string` | yes |  |  |
| `spec.containers[].volumes[].readOnly` | `bool` |  |  |  |
| `spec.containers[].volumes[].azureFile` | `AzureContainerInstanceVolumeAzureFile` |  |  |  |
| `spec.containers[].volumes[].azureFile.shareName` | `string \| valueFrom` | yes |  | AzureStorageShare (`status.outputs.share_name`) |
| `spec.containers[].volumes[].azureFile.storageAccountName` | `string \| valueFrom` | yes |  | AzureStorageAccount (`status.outputs.storage_account_name`) |
| `spec.containers[].volumes[].azureFile.storageAccountKey` | `string \| valueFrom` (sensitive) | yes |  | AzureStorageAccount (`status.outputs.primary_access_key`) |
| `spec.containers[].volumes[].emptyDir` | `bool` |  |  |  |
| `spec.containers[].volumes[].gitRepo` | `AzureContainerInstanceVolumeGitRepo` |  |  |  |
| `spec.containers[].volumes[].gitRepo.url` | `string` | yes |  |  |
| `spec.containers[].volumes[].gitRepo.directory` | `string` |  |  |  |
| `spec.containers[].volumes[].gitRepo.revision` | `string` |  |  |  |
| `spec.containers[].volumes[].secret` | `map<string, string>` (sensitive) |  |  |  |
| `spec.containers[].security` | `AzureContainerInstanceSecurity` |  |  |  |
| `spec.containers[].security.privilegeEnabled` | `bool` |  |  |  |
| `spec.containers[].livenessProbe` | `AzureContainerInstanceProbe` |  |  |  |
| `spec.containers[].livenessProbe.exec` | `[]string` |  |  |  |
| `spec.containers[].livenessProbe.httpGet` | `AzureContainerInstanceProbeHttpGet` |  |  |  |
| `spec.containers[].livenessProbe.httpGet.path` | `string` |  |  |  |
| `spec.containers[].livenessProbe.httpGet.port` | `int32` |  |  |  |
| `spec.containers[].livenessProbe.httpGet.scheme` | `string` |  |  |  |
| `spec.containers[].livenessProbe.httpGet.httpHeaders` | `map<string, string>` |  |  |  |
| `spec.containers[].livenessProbe.initialDelaySeconds` | `int32` |  |  |  |
| `spec.containers[].livenessProbe.periodSeconds` | `int32` |  |  |  |
| `spec.containers[].livenessProbe.failureThreshold` | `int32` |  |  |  |
| `spec.containers[].livenessProbe.successThreshold` | `int32` |  |  |  |
| `spec.containers[].livenessProbe.timeoutSeconds` | `int32` |  |  |  |
| `spec.containers[].readinessProbe` | `AzureContainerInstanceProbe` |  |  |  |
| `spec.containers[].readinessProbe.exec` | `[]string` |  |  |  |
| `spec.containers[].readinessProbe.httpGet` | `AzureContainerInstanceProbeHttpGet` |  |  |  |
| `spec.containers[].readinessProbe.httpGet.path` | `string` |  |  |  |
| `spec.containers[].readinessProbe.httpGet.port` | `int32` |  |  |  |
| `spec.containers[].readinessProbe.httpGet.scheme` | `string` |  |  |  |
| `spec.containers[].readinessProbe.httpGet.httpHeaders` | `map<string, string>` |  |  |  |
| `spec.containers[].readinessProbe.initialDelaySeconds` | `int32` |  |  |  |
| `spec.containers[].readinessProbe.periodSeconds` | `int32` |  |  |  |
| `spec.containers[].readinessProbe.failureThreshold` | `int32` |  |  |  |
| `spec.containers[].readinessProbe.successThreshold` | `int32` |  |  |  |
| `spec.containers[].readinessProbe.timeoutSeconds` | `int32` |  |  |  |
| `spec.initContainers` | `[]AzureContainerInstanceInitContainer` |  |  |  |
| `spec.initContainers[].name` | `string` | yes |  |  |
| `spec.initContainers[].image` | `string` | yes |  |  |
| `spec.initContainers[].environmentVariables` | `map<string, string>` |  |  |  |
| `spec.initContainers[].secureEnvironmentVariables` | `map<string, string>` (sensitive) |  |  |  |
| `spec.initContainers[].commands` | `[]string` |  |  |  |
| `spec.initContainers[].volumes` | `[]AzureContainerInstanceVolume` |  |  |  |
| `spec.initContainers[].volumes[].name` | `string` | yes |  |  |
| `spec.initContainers[].volumes[].mountPath` | `string` | yes |  |  |
| `spec.initContainers[].volumes[].readOnly` | `bool` |  |  |  |
| `spec.initContainers[].volumes[].azureFile` | `AzureContainerInstanceVolumeAzureFile` |  |  |  |
| `spec.initContainers[].volumes[].azureFile.shareName` | `string \| valueFrom` | yes |  | AzureStorageShare (`status.outputs.share_name`) |
| `spec.initContainers[].volumes[].azureFile.storageAccountName` | `string \| valueFrom` | yes |  | AzureStorageAccount (`status.outputs.storage_account_name`) |
| `spec.initContainers[].volumes[].azureFile.storageAccountKey` | `string \| valueFrom` (sensitive) | yes |  | AzureStorageAccount (`status.outputs.primary_access_key`) |
| `spec.initContainers[].volumes[].emptyDir` | `bool` |  |  |  |
| `spec.initContainers[].volumes[].gitRepo` | `AzureContainerInstanceVolumeGitRepo` |  |  |  |
| `spec.initContainers[].volumes[].gitRepo.url` | `string` | yes |  |  |
| `spec.initContainers[].volumes[].gitRepo.directory` | `string` |  |  |  |
| `spec.initContainers[].volumes[].gitRepo.revision` | `string` |  |  |  |
| `spec.initContainers[].volumes[].secret` | `map<string, string>` (sensitive) |  |  |  |
| `spec.initContainers[].security` | `AzureContainerInstanceSecurity` |  |  |  |
| `spec.initContainers[].security.privilegeEnabled` | `bool` |  |  |  |
| `spec.imageRegistryCredentials` | `[]AzureContainerInstanceImageRegistryCredential` |  |  |  |
| `spec.imageRegistryCredentials[].server` | `string \| valueFrom` | yes |  | AzureContainerRegistry (`status.outputs.login_server`) |
| `spec.imageRegistryCredentials[].username` | `string` |  |  |  |
| `spec.imageRegistryCredentials[].password` | `string` (sensitive) |  |  |  |
| `spec.imageRegistryCredentials[].userAssignedIdentityId` | `string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |
| `spec.identity` | `AzureContainerInstanceIdentity` |  |  |  |
| `spec.identity.type` | `enum` | yes |  |  |
| `spec.identity.identityIds` | `[]string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |
| `spec.diagnosticsLogAnalytics` | `AzureContainerInstanceLogAnalytics` |  |  |  |
| `spec.diagnosticsLogAnalytics.workspaceId` | `string \| valueFrom` | yes |  | AzureLogAnalyticsWorkspace (`status.outputs.workspace_customer_id`) |
| `spec.diagnosticsLogAnalytics.workspaceKey` | `string \| valueFrom` (sensitive) | yes |  | AzureLogAnalyticsWorkspace (`status.outputs.primary_shared_key`) |
| `spec.diagnosticsLogAnalytics.logType` | `string` |  |  |  |
| `spec.diagnosticsLogAnalytics.metadata` | `map<string, string>` |  |  |  |
| `spec.dnsConfig` | `AzureContainerInstanceDnsConfig` |  |  |  |
| `spec.dnsConfig.nameservers` | `[]string` | yes |  |  |
| `spec.dnsConfig.searchDomains` | `[]string` |  |  |  |
| `spec.dnsConfig.options` | `[]string` |  |  |  |
| `spec.keyVaultKeyId` | `string \| valueFrom` |  |  | AzureKeyVaultKey (`status.outputs.key_id`) |
| `spec.keyVaultUserAssignedIdentityId` | `string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.resourceGroup

`string | valueFrom` · required

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.region

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.osType

`string` · required

- rule: {"required":true,"string":{"in":["Linux","Windows"]}}

### spec.restartPolicy

`string`

- rule: {"string":{"in":["","Always","Never","OnFailure"]}}

### spec.sku

`string`

- rule: {"string":{"in":["","Standard","Dedicated","Confidential","NotSpecified"]}}

### spec.priority

`string`

- rule: {"string":{"in":["","Regular","Spot"]}}

### spec.ipAddressType

`string`

- rule: {"string":{"in":["","Public","Private","None"]}}

### spec.dnsNameLabel

`string`

### spec.dnsNameLabelReusePolicy

`string`

- rule: {"string":{"in":["","Noreuse","ResourceGroupReuse","SubscriptionReuse","TenantReuse","Unsecure"]}}

### spec.subnetId

`string | valueFrom`

- references: AzureSubnet (`status.outputs.subnet_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.zones

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.exposedPorts

`[]AzureContainerInstancePort`

### spec.exposedPorts[].port

`int32`

- rule: port must be between 1 and 65535

### spec.exposedPorts[].protocol

`string`

- rule: {"string":{"in":["","TCP","UDP"]}}

### spec.containers

`[]AzureContainerInstanceContainer` · required

- rule: {"repeated":{"minItems":"1"}}

### spec.containers[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.containers[].image

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.containers[].cpu

`double` · required

- rule: {"required":true}

### spec.containers[].memory

`double` · required

- rule: {"required":true}

### spec.containers[].cpuLimit

`double` · optional (explicit presence)

- rule: {"double":{"gte":0}}

### spec.containers[].memoryLimit

`double` · optional (explicit presence)

- rule: {"double":{"gte":0}}

### spec.containers[].ports

`[]AzureContainerInstancePort`

### spec.containers[].ports[].port

`int32`

- rule: port must be between 1 and 65535

### spec.containers[].ports[].protocol

`string`

- rule: {"string":{"in":["","TCP","UDP"]}}

### spec.containers[].environmentVariables

`map<string, string>`

### spec.containers[].secureEnvironmentVariables

`map<string, string>` · sensitive

### spec.containers[].commands

`[]string`

### spec.containers[].volumes

`[]AzureContainerInstanceVolume`

- rule: set exactly one volume form: azure_file, empty_dir, git_repo, or secret

### spec.containers[].volumes[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.containers[].volumes[].mountPath

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.containers[].volumes[].readOnly

`bool`

### spec.containers[].volumes[].azureFile

`AzureContainerInstanceVolumeAzureFile`

### spec.containers[].volumes[].azureFile.shareName

`string | valueFrom` · required

- references: AzureStorageShare (`status.outputs.share_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageShare, name: <that resource's name>, fieldPath: status.outputs.share_name}} -- a bare string does not parse

### spec.containers[].volumes[].azureFile.storageAccountName

`string | valueFrom` · required

- references: AzureStorageAccount (`status.outputs.storage_account_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.storage_account_name}} -- a bare string does not parse

### spec.containers[].volumes[].azureFile.storageAccountKey

`string | valueFrom` · required · sensitive

- references: AzureStorageAccount (`status.outputs.primary_access_key`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.primary_access_key}} -- a bare string does not parse

### spec.containers[].volumes[].emptyDir

`bool`

### spec.containers[].volumes[].gitRepo

`AzureContainerInstanceVolumeGitRepo`

### spec.containers[].volumes[].gitRepo.url

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.containers[].volumes[].gitRepo.directory

`string`

### spec.containers[].volumes[].gitRepo.revision

`string`

### spec.containers[].volumes[].secret

`map<string, string>` · sensitive

### spec.containers[].security

`AzureContainerInstanceSecurity`

### spec.containers[].security.privilegeEnabled

`bool`

### spec.containers[].livenessProbe

`AzureContainerInstanceProbe`

### spec.containers[].livenessProbe.exec

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.containers[].livenessProbe.httpGet

`AzureContainerInstanceProbeHttpGet`

### spec.containers[].livenessProbe.httpGet.path

`string`

### spec.containers[].livenessProbe.httpGet.port

`int32`

- rule: port must be between 1 and 65535

### spec.containers[].livenessProbe.httpGet.scheme

`string`

- rule: {"string":{"in":["","http","https"]}}

### spec.containers[].livenessProbe.httpGet.httpHeaders

`map<string, string>`

### spec.containers[].livenessProbe.initialDelaySeconds

`int32`

- rule: {"int32":{"gte":0}}

### spec.containers[].livenessProbe.periodSeconds

`int32`

- rule: {"int32":{"gte":0}}

### spec.containers[].livenessProbe.failureThreshold

`int32`

- rule: {"int32":{"gte":0}}

### spec.containers[].livenessProbe.successThreshold

`int32`

- rule: {"int32":{"gte":0}}

### spec.containers[].livenessProbe.timeoutSeconds

`int32`

- rule: {"int32":{"gte":0}}

### spec.containers[].readinessProbe

`AzureContainerInstanceProbe`

### spec.containers[].readinessProbe.exec

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.containers[].readinessProbe.httpGet

`AzureContainerInstanceProbeHttpGet`

### spec.containers[].readinessProbe.httpGet.path

`string`

### spec.containers[].readinessProbe.httpGet.port

`int32`

- rule: port must be between 1 and 65535

### spec.containers[].readinessProbe.httpGet.scheme

`string`

- rule: {"string":{"in":["","http","https"]}}

### spec.containers[].readinessProbe.httpGet.httpHeaders

`map<string, string>`

### spec.containers[].readinessProbe.initialDelaySeconds

`int32`

- rule: {"int32":{"gte":0}}

### spec.containers[].readinessProbe.periodSeconds

`int32`

- rule: {"int32":{"gte":0}}

### spec.containers[].readinessProbe.failureThreshold

`int32`

- rule: {"int32":{"gte":0}}

### spec.containers[].readinessProbe.successThreshold

`int32`

- rule: {"int32":{"gte":0}}

### spec.containers[].readinessProbe.timeoutSeconds

`int32`

- rule: {"int32":{"gte":0}}

### spec.initContainers

`[]AzureContainerInstanceInitContainer`

### spec.initContainers[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.initContainers[].image

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.initContainers[].environmentVariables

`map<string, string>`

### spec.initContainers[].secureEnvironmentVariables

`map<string, string>` · sensitive

### spec.initContainers[].commands

`[]string`

### spec.initContainers[].volumes

`[]AzureContainerInstanceVolume`

- rule: set exactly one volume form: azure_file, empty_dir, git_repo, or secret

### spec.initContainers[].volumes[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.initContainers[].volumes[].mountPath

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.initContainers[].volumes[].readOnly

`bool`

### spec.initContainers[].volumes[].azureFile

`AzureContainerInstanceVolumeAzureFile`

### spec.initContainers[].volumes[].azureFile.shareName

`string | valueFrom` · required

- references: AzureStorageShare (`status.outputs.share_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageShare, name: <that resource's name>, fieldPath: status.outputs.share_name}} -- a bare string does not parse

### spec.initContainers[].volumes[].azureFile.storageAccountName

`string | valueFrom` · required

- references: AzureStorageAccount (`status.outputs.storage_account_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.storage_account_name}} -- a bare string does not parse

### spec.initContainers[].volumes[].azureFile.storageAccountKey

`string | valueFrom` · required · sensitive

- references: AzureStorageAccount (`status.outputs.primary_access_key`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.primary_access_key}} -- a bare string does not parse

### spec.initContainers[].volumes[].emptyDir

`bool`

### spec.initContainers[].volumes[].gitRepo

`AzureContainerInstanceVolumeGitRepo`

### spec.initContainers[].volumes[].gitRepo.url

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.initContainers[].volumes[].gitRepo.directory

`string`

### spec.initContainers[].volumes[].gitRepo.revision

`string`

### spec.initContainers[].volumes[].secret

`map<string, string>` · sensitive

### spec.initContainers[].security

`AzureContainerInstanceSecurity`

### spec.initContainers[].security.privilegeEnabled

`bool`

### spec.imageRegistryCredentials

`[]AzureContainerInstanceImageRegistryCredential`

### spec.imageRegistryCredentials[].server

`string | valueFrom` · required

- references: AzureContainerRegistry (`status.outputs.login_server`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureContainerRegistry, name: <that resource's name>, fieldPath: status.outputs.login_server}} -- a bare string does not parse

### spec.imageRegistryCredentials[].username

`string`

### spec.imageRegistryCredentials[].password

`string` · sensitive

### spec.imageRegistryCredentials[].userAssignedIdentityId

`string | valueFrom`

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.identity

`AzureContainerInstanceIdentity`

- rule: identity_ids is required for USER_ASSIGNED and SYSTEM_AND_USER_ASSIGNED and must be empty for SYSTEM_ASSIGNED

### spec.identity.type

`enum` · required

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_container_instance_identity_type_unspecified`
- `SYSTEM_ASSIGNED`
- `USER_ASSIGNED`
- `SYSTEM_AND_USER_ASSIGNED`

### spec.identity.identityIds

`[]string | valueFrom`

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.diagnosticsLogAnalytics

`AzureContainerInstanceLogAnalytics`

- rule: metadata requires log_type -- the provider only sends metadata alongside a log type and silently discards it otherwise

### spec.diagnosticsLogAnalytics.workspaceId

`string | valueFrom` · required

- references: AzureLogAnalyticsWorkspace (`status.outputs.workspace_customer_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureLogAnalyticsWorkspace, name: <that resource's name>, fieldPath: status.outputs.workspace_customer_id}} -- a bare string does not parse

### spec.diagnosticsLogAnalytics.workspaceKey

`string | valueFrom` · required · sensitive

- references: AzureLogAnalyticsWorkspace (`status.outputs.primary_shared_key`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureLogAnalyticsWorkspace, name: <that resource's name>, fieldPath: status.outputs.primary_shared_key}} -- a bare string does not parse

### spec.diagnosticsLogAnalytics.logType

`string`

- rule: {"string":{"in":["","ContainerInsights","ContainerInstanceLogs"]}}

### spec.diagnosticsLogAnalytics.metadata

`map<string, string>`

### spec.dnsConfig

`AzureContainerInstanceDnsConfig`

### spec.dnsConfig.nameservers

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.dnsConfig.searchDomains

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.dnsConfig.options

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.keyVaultKeyId

`string | valueFrom`

- references: AzureKeyVaultKey (`status.outputs.key_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureKeyVaultKey, name: <that resource's name>, fieldPath: status.outputs.key_id}} -- a bare string does not parse

### spec.keyVaultUserAssignedIdentityId

`string | valueFrom`

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.tags

`map<string, string>`

## Validation Rules

- `container_instance_spot_requires_ip_none`: priority "Spot" requires ip_address_type "None" -- Azure gives Spot groups no group IP (the provider enforces this before apply)
- `container_instance_subnet_conflicts_dns_label`: subnet_id and dns_name_label cannot be set together -- a subnet-joined group is private and has no public DNS name
- `container_instance_dns_label_requires_group_ip`: dns_name_label, dns_name_label_reuse_policy, and exposed_ports require an ip_address_type other than "None" -- without a group IP Azure silently discards them
- `container_instance_exposed_ports_declared_on_containers`: every exposed_ports entry must match a port and protocol declared by some container's ports -- Azure only exposes ports a container listens on (the provider enforces this before apply)

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureContainerInstance, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.container_group_id` | `string` |  |
| `status.outputs.container_group_name` | `string` |  |
| `status.outputs.ip_address` | `string` |  |
| `status.outputs.fqdn` | `string` |  |
| `status.outputs.identity_principal_id` | `string` |  |
| `status.outputs.identity_tenant_id` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.subnetId` | AzureSubnet | `status.outputs.subnet_id` |
| `spec.containers[].volumes[].azureFile.shareName` | AzureStorageShare | `status.outputs.share_name` |
| `spec.containers[].volumes[].azureFile.storageAccountName` | AzureStorageAccount | `status.outputs.storage_account_name` |
| `spec.containers[].volumes[].azureFile.storageAccountKey` | AzureStorageAccount | `status.outputs.primary_access_key` |
| `spec.initContainers[].volumes[].azureFile.shareName` | AzureStorageShare | `status.outputs.share_name` |
| `spec.initContainers[].volumes[].azureFile.storageAccountName` | AzureStorageAccount | `status.outputs.storage_account_name` |
| `spec.initContainers[].volumes[].azureFile.storageAccountKey` | AzureStorageAccount | `status.outputs.primary_access_key` |
| `spec.imageRegistryCredentials[].server` | AzureContainerRegistry | `status.outputs.login_server` |
| `spec.imageRegistryCredentials[].userAssignedIdentityId` | AzureUserAssignedIdentity | `status.outputs.identity_id` |
| `spec.identity.identityIds` | AzureUserAssignedIdentity | `status.outputs.identity_id` |
| `spec.diagnosticsLogAnalytics.workspaceId` | AzureLogAnalyticsWorkspace | `status.outputs.workspace_customer_id` |
| `spec.diagnosticsLogAnalytics.workspaceKey` | AzureLogAnalyticsWorkspace | `status.outputs.primary_shared_key` |
| `spec.keyVaultKeyId` | AzureKeyVaultKey | `status.outputs.key_id` |
| `spec.keyVaultUserAssignedIdentityId` | AzureUserAssignedIdentity | `status.outputs.identity_id` |

## See Also

- [Overview](../README.md)

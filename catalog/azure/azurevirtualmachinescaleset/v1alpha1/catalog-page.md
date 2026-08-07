# Azure Virtual Machine Scale Set

Deploys an Azure Virtual Machine Scale Set — a fleet of identical VMs managed as one resource. The spec is the instance template (image, size, OS profile, disks, network) plus the fleet controls (instance count, zones, upgrade policy, spot economics, instance repair) and the orchestration mode. ARM has exactly one scale-set resource type with an orchestration-mode property, and this component models it that way: one kind, an explicit `orchestrationMode`, and validation that gates each mode's capabilities.

## What Gets Created

When you deploy an AzureVirtualMachineScaleSet resource, Planton provisions:

- **Virtual Machine Scale Set** — a `Microsoft.Compute/virtualMachineScaleSets` resource in the specified region and resource group, in the chosen orchestration mode (FLEXIBLE by default, UNIFORM on request), stamping instances from the declared template
- **Azure Tags** — resource metadata tags applied to the scale set for tracking and governance, merged with any user-supplied tags (user tags win on key collision)

Nothing else is created here. Subnets, load balancers, NSGs, identities, and public IP prefixes are referenced, never created — the scale set consumes their outputs. Per-instance NICs, disks, and (optionally) public IPs are stamped from inline templates because they live and die with each instance.

## The Two Orchestration Modes

| | FLEXIBLE (default) | UNIFORM |
|---|---|---|
| **Philosophy** | A resilient VM group: instances are near-full VMs spread across fault domains | The classic large-fleet engine: identical instances managed as one unit |
| **Exclusive capabilities** | Mixed-SKU profiles (`skuName: Mix` + `skuProfile`), spot/on-demand `priorityMix`, per-OS `patchMode`/`patchAssessmentMode`/hotpatching, `networkApiVersion`, standalone-VM attach | `overprovision`, `scaleIn` policy, automatic OS image upgrades, spot `restore`, `galleryApplications`, trusted launch (`secureBoot`/`vtpm` + confidential disk encryption), LB health probe, host groups, edge zones, NAT-rule templates |
| **Health monitoring** | Application health extension only | Health extension or an AzureLoadBalancer probe (`healthProbeId`) |
| **Identity** | USER_ASSIGNED only (ARM's contract) | SYSTEM_ASSIGNED, USER_ASSIGNED, or both |
| **Required extras** | `platformFaultDomainCount` (1 with zones; the region's max for regional spreading) | — |

Every mode-gated field fails validation with a message naming the contract when used in the wrong mode — the spec tells the truth about what each mode supports.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An Azure Resource Group** where the scale set will be created (can reference an AzureResourceGroup resource)
- **An Azure Subnet** for the instances' NIC templates (can reference an AzureSubnet resource)
- **For load-balanced fleets**: an AzureLoadBalancer whose name-keyed `backend_pool_ids` output the IP configurations reference

## Quick Start

Create a file `fleet.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureVirtualMachineScaleSet
metadata:
  name: web-fleet
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureVirtualMachineScaleSet.web-fleet
spec:
  region: eastus
  resourceGroup:
    value: my-rg
  name: web-fleet
  skuName: Standard_D2s_v3
  instances: 3
  platformFaultDomainCount: 1
  osProfile:
    computerNamePrefix: web
    linux:
      adminUsername: azureuser
      sshPublicKeys:
        - publicKey: "ssh-ed25519 AAAA..."
  osDisk:
    caching: READ_WRITE
    storageAccountType: PREMIUM_LRS
  sourceImageReference:
    publisher: Canonical
    offer: ubuntu-24_04-lts
    sku: server
    version: latest
  networkInterfaces:
    - name: primary
      primary: true
      ipConfigurations:
        - name: internal
          primary: true
          subnetId:
            value: /subscriptions/xxx/resourceGroups/my-rg/providers/Microsoft.Network/virtualNetworks/my-vnet/subnets/app
```

Deploy:

```shell
planton apply -f fleet.yaml
```

This creates a three-instance FLEXIBLE Linux fleet in the subnet, SSH-key authenticated, on Premium SSDs.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region (e.g., `eastus`). Must match every referenced subnet and load balancer. | Required |
| `resourceGroup` | `StringValueOrRef` | Azure Resource Group name. Can reference an AzureResourceGroup via `valueFrom`. | Required |
| `name` | `string` | Scale-set name, unique within the resource group. Instance computer names derive from `osProfile.computerNamePrefix` (or this name). | Required, 1-64 chars |
| `skuName` | `string` | The VM size every instance uses (e.g., `Standard_D2s_v3`). On FLEXIBLE sets, `Mix` activates `skuProfile`. | Required |
| `osProfile` | `object` | Exactly one of `linux` / `windows`, carrying authentication and OS management. | Required, XOR enforced |
| `osDisk` | `object` | The OS disk template: `caching`, `storageAccountType`, optional size, ephemeral settings, encryption. | Required |
| `networkInterfaces[]` | `list` | NIC templates — at least one; the first must be `primary` when several are declared. | Required, min 1 |
| an image source | — | Exactly one of `sourceImageReference` (publisher/offer/sku/version) or `sourceImageId` (custom/gallery image ARM ID). | XOR enforced |

### Fleet Controls

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `orchestrationMode` | `enum` | `FLEXIBLE` | `FLEXIBLE` or `UNIFORM`. Fixed at creation; gates mode-specific capabilities. |
| `instances` | `int32` | platform-managed | Instance count (0-1000). Unset lets an autoscaler own the count. |
| `zones` | `list(string)` | none | Availability zones instances spread across (e.g., `["1","2","3"]`). Removing a zone replaces the set. |
| `zoneBalance` | `bool` | `false` | Strictly even zone distribution. Requires `zones`. Fixed at creation. |
| `platformFaultDomainCount` | `int32` | Azure picks (UNIFORM) | REQUIRED on FLEXIBLE sets: 1 with zones, or the region's max for regional spreading. Fixed at creation. |
| `upgradePolicy` | `object` | `MANUAL` | `mode` (`MANUAL`/`AUTOMATIC`/`ROLLING`), the `rolling` batch contract (required for ROLLING), UNIFORM `automaticOsUpgrade`, and the UNIFORM LB `healthProbeId`. |
| `automaticInstanceRepair` | `object` | off | Replace/restart/reimage unhealthy instances automatically. Requires health monitoring. |
| `terminationNotification` | `object` | off | A scheduled event up to 15 minutes before instance termination. |
| `spot` | `object` | on-demand | Presence makes the fleet spot: `evictionPolicy` (required), `maxBidPrice`, UNIFORM `restore`, FLEXIBLE `priorityMix`. |
| `scaleIn` | `object` | Azure default | UNIFORM only: which instances scale-in removes first (`DEFAULT`/`NEWEST_VM`/`OLDEST_VM`) and force-deletion. |
| `overprovision` | `bool` | `true` (UNIFORM) | UNIFORM only: provision extras during scale-out, keep the first healthy ones. |
| `skuProfile` | `object` | none | FLEXIBLE + `skuName: Mix` only: up to 5 candidate sizes with an allocation strategy (`LOWEST_PRICE`/`CAPACITY_OPTIMIZED`/`PRIORITIZED`). |

### Instance Template

| Field | Type | Description |
|-------|------|-------------|
| `osProfile.linux` | `object` | `adminUsername`, `sshPublicKeys[]` (the production path), optional `adminPassword` (sensitive), `disablePasswordAuthentication` (default true), FLEXIBLE per-OS patch modes. |
| `osProfile.windows` | `object` | `adminUsername`, `adminPassword` (sensitive, required), `automaticUpdatesEnabled`, `timezone`, `winrmListeners[]`, `additionalUnattendContents[]` (sensitive), `licenseType`, FLEXIBLE patch modes + `hotpatchingEnabled`. |
| `dataDisks[]` | `list` | Per-instance data-disk templates: LUN, caching, size, SKU (incl. UltraSSD/PremiumV2 with dialed IOPS/MBps), create option, CMK encryption. |
| `networkInterfaces[].ipConfigurations[]` | `list` | Subnet reference, IP version, load-balancer pool and (UNIFORM) NAT-rule references via the LB's name-keyed map outputs, Application Gateway pool references via the gateway's `backend_address_pool_ids` map output, ASG IDs, and an optional per-instance `publicIpAddress` template. |
| `extensions[]` | `list` | VM extensions — the `ApplicationHealthLinux`/`ApplicationHealthWindows` extension is the load-bearing one (rolling upgrades, instance repair, and hotpatching key on it). `protectedSettings` is sensitive; or source it from Key Vault. |
| `identity` | `object` | `SYSTEM_ASSIGNED` (UNIFORM only), `USER_ASSIGNED` (referencing AzureUserAssignedIdentity resources), or both. |
| `security` | `object` | UNIFORM trusted launch (`secureBootEnabled` + `vtpmEnabled`) and mode-independent `encryptionAtHostEnabled`. |
| `customData` / `userData` | `string` | Base64 cloud-init (sensitive, first boot only) / base64 IMDS-readable data (updatable; never secrets). |
| `secrets[]` | `list` | Key Vault certificates installed at provisioning (vault reference + certificate URLs; Windows adds the store). |
| `plan` | `object` | Marketplace purchase plan for third-party images. |
| `placement` | `object` | Proximity placement group, capacity reservation group, UNIFORM host group, `singlePlacementGroup`. |
| `bootDiagnostics` | `object` | Presence enables it; empty uses Azure's managed storage. |
| `tags` | `map` | Free-form tags merged over Planton-derived tags (user wins). |

## Examples

### Load-Balanced Rolling-Upgrade Fleet

The production shape: instances join an AzureLoadBalancer pool by its name-keyed map output, the health extension reports instance health, and template changes roll in health-checked batches:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureVirtualMachineScaleSet
metadata:
  name: web-fleet
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: prod-rg
  name: web-fleet
  skuName: Standard_D2s_v3
  instances: 6
  zones: ["1", "2", "3"]
  platformFaultDomainCount: 1
  osProfile:
    computerNamePrefix: web
    linux:
      adminUsername: azureuser
      sshPublicKeys:
        - publicKey: "ssh-ed25519 AAAA..."
  osDisk:
    caching: READ_ONLY
    storageAccountType: PREMIUM_LRS
    diffDiskSettings:
      placement: CACHE_DISK
  sourceImageReference:
    publisher: Canonical
    offer: ubuntu-24_04-lts
    sku: server
    version: latest
  networkInterfaces:
    - name: primary
      primary: true
      acceleratedNetworkingEnabled: true
      ipConfigurations:
        - name: internal
          primary: true
          subnetId:
            valueFrom:
              name: prod-app-subnet
          loadBalancerBackendAddressPoolIds:
            - valueFrom:
                name: prod-lb
                fieldPath: status.outputs.backend_pool_ids.web
  extensions:
    - name: health
      publisher: Microsoft.ManagedServices
      type: ApplicationHealthLinux
      typeHandlerVersion: "1.0"
      settings: '{"protocol":"http","port":80,"requestPath":"/healthz"}'
  upgradePolicy:
    mode: ROLLING
    rolling:
      maxBatchInstancePercent: 20
      maxUnhealthyInstancePercent: 20
      maxUnhealthyUpgradedInstancePercent: 20
      pauseTimeBetweenBatches: PT30S
  automaticInstanceRepair:
    enabled: true
```

### Spot Batch Fleet with an On-Demand Base

FLEXIBLE priority mixing guarantees two on-demand instances; everything above runs on spot at up to a quarter of the price:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureVirtualMachineScaleSet
metadata:
  name: batch-workers
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: prod-rg
  name: batch-workers
  skuName: Standard_D4s_v3
  instances: 20
  platformFaultDomainCount: 1
  spot:
    evictionPolicy: DELETE
    priorityMix:
      baseRegularCount: 2
      regularPercentageAboveBase: 0
  osProfile:
    linux:
      adminUsername: azureuser
      sshPublicKeys:
        - publicKey: "ssh-ed25519 AAAA..."
  osDisk:
    caching: READ_WRITE
    storageAccountType: STANDARD_SSD_LRS
  sourceImageReference:
    publisher: Canonical
    offer: ubuntu-24_04-lts
    sku: server
    version: latest
  networkInterfaces:
    - name: primary
      primary: true
      ipConfigurations:
        - name: internal
          primary: true
          subnetId:
            valueFrom:
              name: prod-batch-subnet
  terminationNotification:
    timeout: PT10M
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `scale_set_id` | `string` | ARM ID of the scale set — what an AzureVirtualMachine's `availability.virtualMachineScaleSetId` references (FLEXIBLE attach), and what monitoring/autoscale resources scope to |
| `scale_set_name` | `string` | The scale set's name as deployed |
| `unique_id` | `string` | The scale set's globally unique ARM-assigned identifier |
| `system_assigned_identity_principal_id` | `string` | The system-assigned identity's principal ID (UNIFORM sets with a SYSTEM_ASSIGNED identity) — the AzureRoleAssignment seam |

## Operational Notes

- **The orchestration mode is fixed at creation** — changing it replaces the fleet
- **Upgrade mode `MANUAL` (the default) applies template changes to NEW instances only** — existing instances wait for a manual upgrade; `ROLLING` is the production choice and requires health monitoring
- **Removing a zone from `zones` replaces the scale set**; adding one does not
- **`singlePlacementGroup` can go true→false but never back** — the false state is permanent
- **Spot fleets should handle the termination notification** — instances can be evicted with only the notification lead time

## Related Components

- [AzureResourceGroup](/docs/catalog/azure/azureresourcegroup) — provides the resource group for fleet placement
- [AzureSubnet](/docs/catalog/azure/azuresubnet) — where the instances' NIC templates deploy
- [AzureLoadBalancer](/docs/catalog/azure/azureloadbalancer) — exports the name-keyed `backend_pool_ids`, `nat_rule_ids`, and `probe_ids` maps the fleet references
- [AzureUserAssignedIdentity](/docs/catalog/azure/azureuserassignedidentity) — the identities instances authenticate to Azure services with
- [AzureVirtualMachine](/docs/catalog/azure/azurevirtualmachine) — attaches to a FLEXIBLE set via `availability.virtualMachineScaleSetId`
- [AzurePublicIpPrefix](/docs/catalog/azure/azurepublicipprefix) — the reserved CIDR per-instance public IPs can draw from

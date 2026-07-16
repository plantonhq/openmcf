---
title: "Virtual Machine"
description: "Virtual Machine deployment documentation"
icon: "package"
order: 100
componentName: "azurevirtualmachine"
---

# Azure Virtual Machine

Creates an Azure Virtual Machine -- the compute instance itself: its size, image, OS profile, OS disk, identity, placement, and security posture. The VM is deliberately just the machine: its network presence, data volumes, and identities are first-class referenced resources, matching Azure's own model.

## What Gets Created

When you deploy an AzureVirtualMachine resource, Planton provisions:

- **Virtual Machine** — an `azurerm_linux_virtual_machine` OR an `azurerm_windows_virtual_machine`, chosen by the spec's explicit OS discriminator (`osProfile.linux` XOR `osProfile.windows`); ARM models the two OSes as separate management surfaces, and the module deploys the matching one
- **Data Disk Attachments** — one `azurerm_virtual_machine_data_disk_attachment` per entry in `dataDiskAttachments`, each mounting a referenced `AzureManagedDisk` at a LUN

Everything else is referenced, never created here. The VM's network presence is one or more `AzureNetworkInterface` resources -- public IPs, NSG filtering, and subnet placement all live NIC-side. Its data volumes are `AzureManagedDisk` resources whose data outlives the machine. Only the OS disk is inline (born and dying with the VM), unless the VM boots from an existing referenced OS disk -- the golden-disk path.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A resource group** to create the VM in (an `AzureResourceGroup` in composed environments)
- **At least one network interface** (`AzureNetworkInterface`) in the same region; the first listed is the primary
- **Compute write rights**: `Microsoft.Compute/virtualMachines/write` (Virtual Machine Contributor, Contributor, or Owner)
- Optional: `AzureManagedDisk` resources for data volumes or an existing OS disk, `AzureUserAssignedIdentity` for user-assigned identity, `AzureKeyVault` for certificate installation

## Quick Start

Create a file `vm.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureVirtualMachine
metadata:
  name: app-vm
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureVirtualMachine.app-vm
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: prod-rg
  name: app-vm
  size: Standard_D2s_v3
  networkInterfaceIds:
    - valueFrom:
        name: app-vm-nic
  osProfile:
    linux:
      adminUsername: azureuser
      sshPublicKeys:
        - publicKey: ssh-ed25519 AAAA... you@workstation
  osDisk:
    caching: READ_WRITE
    storageAccountType: PREMIUM_LRS
  sourceImageReference:
    publisher: Canonical
    offer: ubuntu-24_04-lts
    sku: server
    version: latest
  bootDiagnostics: {}
```

Deploy:

```shell
planton apply -f vm.yaml
```

This is the canonical Linux production shape: an Ubuntu 24.04 LTS VM authenticated by SSH keys only (password authentication is disabled by default), attached to a referenced NIC that carries its entire network posture, with a Premium OS disk and managed-storage boot diagnostics -- the first tool for debugging a VM that will not boot.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region the VM runs in. Must match the region of every referenced NIC and disk. Changing it replaces the VM. | Required |
| `resourceGroup` | `StringValueOrRef` | Resource group name. Defaults to referencing an `AzureResourceGroup`'s name output. | Required |
| `name` | `string` | VM name, unique within the resource group. 1-64 characters for Linux, 1-15 for Windows (the OS computer name defaults to it, which is where the Windows limit bites). Changing it replaces the VM. | Required |
| `size` | `string` | VM size (SKU), e.g. `Standard_D2s_v3`. Determines vCPUs, memory, temp disk, accelerated-networking and ultra-disk support, and hourly cost. Resizing reboots in place. | Required |
| `networkInterfaceIds` | `list(StringValueOrRef)` | The attached NICs -- at least one; the FIRST is primary. References `AzureNetworkInterface` resources; subnet placement, public IPs, and NSG filtering all live NIC-side. | Required, min 1 |
| `osProfile` | `object` | Exactly one of `linux` or `windows`. See below. | Required |
| `osDisk` | `object` | The inline OS disk: `caching` and `storageAccountType` required. See below. | Required |

### OS Profile (exactly one of `linux` / `windows`)

`osProfile.computerName` optionally overrides the OS hostname (defaults to the VM name; useful when a VM name exceeds Windows' 15-character computer-name limit).

**Linux** (`osProfile.linux`) -- SSH-first authentication:

| Field | Description |
|-------|-------------|
| `adminUsername` | Admin account name. Required when booting from an image; must stay empty when booting from an existing OS disk. |
| `sshPublicKeys` | SSH public keys for the admin account -- the production authentication path. Each key's `username` defaults to `adminUsername`. |
| `adminPassword` | `StringValueOrRef` (secret). Only meaningful when `disablePasswordAuthentication` is explicitly `false`. |
| `disablePasswordAuthentication` | Defaults `true` (keys only) -- the right posture. Setting `false` requires `adminPassword`. |
| `patchMode` | `LINUX_IMAGE_DEFAULT` (Azure's default) or `LINUX_AUTOMATIC_BY_PLATFORM` (Azure Update Manager orchestrates; unlocks `patching.rebootSetting`). |
| `licenseType` | Bring-your-own-subscription for commercial distros: `RHEL_BYOS`, `SLES_BYOS`, `UBUNTU_PRO`, and the RHEL/SLES SAP variants. |

**Windows** (`osProfile.windows`) -- password + WinRM/unattend:

| Field | Description |
|-------|-------------|
| `adminUsername` / `adminPassword` | Both required when booting from an image (password 8-123 chars, 3 of 4 complexity classes -- ARM enforces). Source the password from a secret reference, never a manifest literal. |
| `patchMode` | `MANUAL`, `AUTOMATIC_BY_OS` (Azure's default: Windows Update as configured in the image), or `WINDOWS_AUTOMATIC_BY_PLATFORM` (prerequisite for hotpatching and reboot control). |
| `automaticUpdatesEnabled` | Windows Update automatic updates; Azure defaults `true`. Fixed at creation. |
| `hotpatchingEnabled` | Security updates without reboots, on supported Windows Server Azure Edition images only. Requires `WINDOWS_AUTOMATIC_BY_PLATFORM`. |
| `timezone` | Windows time zone, e.g. `Pacific Standard Time`. Unset uses UTC. |
| `winrmListeners` | WinRM remote-management listeners; `HTTPS` requires the certificate's Key Vault secret URL. |
| `additionalUnattendContents` | Raw unattend.xml fragments (`AUTO_LOGON` / `FIRST_LOGON_COMMANDS`) for pre-agent bootstrap; treated as secret material. |
| `licenseType` | Azure Hybrid Benefit: `WINDOWS_SERVER` / `WINDOWS_CLIENT` to bring an existing license. |

Booting from an existing OS disk (`osManagedDiskId`) forbids ALL authentication fields in either profile -- the disk already contains its users; spec-level validation enforces it.

### Image Source (exactly one)

| Field | Description |
|-------|-------------|
| `sourceImageReference` | Marketplace/platform image by its four coordinates (`publisher`/`offer`/`sku`/`version`). `version: latest` resolves at CREATION only -- the VM does not follow new image releases. |
| `sourceImageId` | A custom or gallery image by ARM ID (managed image, or Shared Image Gallery image/version -- community and direct-shared included). |
| `osManagedDiskId` | `StringValueOrRef` to an existing `AzureManagedDisk` that already carries an OS -- the disk-swap/golden-disk path. |

### OS Disk

| Field | Description |
|-------|-------------|
| `caching` | `NONE` / `READ_ONLY` / `READ_WRITE` (the general-purpose OS-disk mode). Required. |
| `storageAccountType` | `PREMIUM_LRS` (production default), `STANDARD_LRS`, `STANDARD_SSD_LRS`, or the ZRS variants (survive a zone outage but conflict with zone-pinning). PremiumV2/Ultra are data-disk-only. Required. |
| `diskSizeGb` | Unset inherits the image's size -- correct for almost everything. Can only grow. |
| `diffDiskSettings` | Presence makes the OS disk EPHEMERAL: it lives on local VM storage (cache/temp/NVMe placement), free and fast, WIPED on every stop/deallocate -- stateless fleets only. |
| `diskEncryptionSetId` / `secureVmDiskEncryptionSetId` | Customer-managed-key encryption (mutually exclusive; the secure variant is for confidential VMs and requires `securityEncryptionType`). |
| `securityEncryptionType` | Confidential-VM guest-state encryption: `VM_GUEST_STATE_ONLY` or `DISK_WITH_VM_GUEST_STATE`. Requires `security.vtpmEnabled` (the disk variant also requires `secureBootEnabled`). |
| `writeAcceleratorEnabled` | M-series + Premium + caching `NONE` -- sub-millisecond write latency. |

### Data Disk Attachments

Each `dataDiskAttachments` entry mounts a referenced first-class `AzureManagedDisk`:

| Field | Description |
|-------|-------------|
| `managedDiskId` | `StringValueOrRef` to the disk (defaults to referencing an `AzureManagedDisk`'s `disk_id` output). Required. |
| `lun` | Logical unit number 0-63, unique per VM -- the stable identity the OS addresses the disk by (`/dev/disk/azure/scsi1/lun{n}`). Keep LUNs stable across changes. |
| `caching` | `NONE` (required above 4 TiB; right for write-heavy volumes), `READ_ONLY`, or `READ_WRITE`. Required. |
| `writeAcceleratorEnabled` | Per-attachment Write Accelerator. |

Attachments are separate resources on both engines: disks attach and detach in place, and the disk -- and its data -- outlives the VM.

### Optional Fields (grouped)

| Group | Fields |
|-------|--------|
| Identity | `identity.type` (`SYSTEM_ASSIGNED` -- principal surfaced in outputs for `AzureRoleAssignment` grants; `USER_ASSIGNED` -- `identityIds` reference `AzureUserAssignedIdentity`; `SYSTEM_AND_USER_ASSIGNED`) |
| Spot | `spot` (presence = spot instance): `evictionPolicy` (`DEALLOCATE` keeps disks, restartable; `DELETE` removes VM and disks), `maxBidPrice` (default -1: pay up to on-demand, never price-evicted) |
| Placement | `availability`: `zone` (`"1"`/`"2"`/`"3"`) XOR `availabilitySetId`; `proximityPlacementGroupId`, `capacityReservationGroupId`, `dedicatedHostId` XOR `dedicatedHostGroupId`, `virtualMachineScaleSetId` (Flexible orchestration) + `platformFaultDomain` |
| Security | `security.secureBootEnabled` + `vtpmEnabled` (= trusted launch, the production posture for Gen2 images), `encryptionAtHostEnabled` (covers temp disks and caches; subscription feature must be registered) |
| Patching | `patching.assessmentMode` (`ASSESSMENT_IMAGE_DEFAULT` / `ASSESSMENT_AUTOMATIC_BY_PLATFORM`), `rebootSetting` (`ALWAYS`/`IF_REQUIRED`/`NEVER`), `bypassPlatformSafetyChecksOnUserScheduleEnabled` -- the latter two require the OS profile's patch mode to be AUTOMATIC_BY_PLATFORM |
| Diagnostics & events | `bootDiagnostics` (presence enables; empty = Azure-managed storage), `terminationNotification` (pre-termination scheduled event, up to PT15M drain), `osImageNotification` (pre-OS-image-upgrade event) |
| Boot & provisioning | `customData` (cloud-init, base64, secret, first boot only -- changing it replaces the VM), `userData` (base64, IMDS-readable, updatable -- never secrets), `extensionsTimeBudget` (PT15M-PT2H; Azure defaults PT1H30M), `provisionVmAgent` / `allowExtensionOperations` (both default true) |
| Applications & certs | `galleryApplications` (up to 100 VM Applications, ordered), `secrets` (Key Vault certificates installed at provisioning; references `AzureKeyVault`; Windows entries name a certificate store), `plan` (marketplace purchase plan for third-party images) |
| Hardware | `diskControllerType` (`SCSI` / `NVME` -- NVMe needs a supported size + Gen2 image), `additionalCapabilities` (`ultraSsdEnabled`, `hibernationEnabled`), `edgeZone` |
| Tags | `tags` merged over the Planton-derived resource tags (organization, environment, resource id); a user tag with the same key wins |

Secret-by-default: `customData`, both `adminPassword` fields, and unattend content are `sensitive` -- Planton forces them to managed-secret references.

## Examples

### Windows Server with Trusted Launch and Hybrid Benefit

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureVirtualMachine
metadata:
  name: win-vm
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureVirtualMachine.win-vm
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: prod-rg
  name: win-vm  # Windows computer names cap at 15 characters
  size: Standard_D2s_v3
  networkInterfaceIds:
    - valueFrom:
        name: win-vm-nic
  osProfile:
    windows:
      adminUsername: azureadmin
      adminPassword:
        valueFrom:
          name: win-vm-admin-password  # never a literal password in a manifest
      licenseType: WINDOWS_SERVER  # Azure Hybrid Benefit if you hold a license
  osDisk:
    caching: READ_WRITE
    storageAccountType: PREMIUM_LRS
    diskSizeGb: 128
  sourceImageReference:
    publisher: MicrosoftWindowsServer
    offer: WindowsServer
    sku: 2022-datacenter-g2
    version: latest
  availability:
    zone: "1"
  security:
    secureBootEnabled: true  # trusted launch on a Gen2 image
    vtpmEnabled: true
  bootDiagnostics: {}
```

### Spot Worker with a Persistent Data Disk

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureVirtualMachine
metadata:
  name: batch-worker
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureVirtualMachine.batch-worker
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: prod-rg
  name: batch-worker
  size: Standard_D4s_v3
  networkInterfaceIds:
    - valueFrom:
        name: batch-worker-nic
  osProfile:
    linux:
      adminUsername: azureuser
      sshPublicKeys:
        - publicKey: ssh-ed25519 AAAA... you@workstation
  osDisk:
    caching: READ_WRITE
    storageAccountType: STANDARD_SSD_LRS  # the OS disk is disposable on a spot worker
  sourceImageReference:
    publisher: Canonical
    offer: ubuntu-24_04-lts
    sku: server
    version: latest
  spot:
    evictionPolicy: DEALLOCATE  # disks survive eviction; the VM restarts when capacity returns
  identity:
    type: SYSTEM_ASSIGNED  # grant it queue/storage access via AzureRoleAssignment
  dataDiskAttachments:
    - managedDiskId:
        valueFrom:
          name: batch-checkpoints-disk
      lun: 0
      caching: READ_ONLY
  terminationNotification:
    timeout: PT15M  # maximum drain window before eviction
  bootDiagnostics: {}
```

### Booting from an Existing OS Disk (Golden-Disk Recovery)

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureVirtualMachine
metadata:
  name: recovered-vm
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureVirtualMachine.recovered-vm
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: prod-rg
  name: recovered-vm
  size: Standard_D2s_v3
  networkInterfaceIds:
    - valueFrom:
        name: recovered-vm-nic
  osProfile:
    linux: {}  # selects the OS; NO auth fields -- the disk already contains its users
  osDisk:
    caching: READ_WRITE
    storageAccountType: PREMIUM_LRS
  osManagedDiskId:
    valueFrom:
      name: app-vm-os-disk  # an AzureManagedDisk that already carries an OS
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `vm_id` | `string` | Full ARM ID of the VM -- what role assignments, diagnostics settings, and backup policies scope to |
| `vm_name` | `string` | The VM's name as deployed |
| `virtual_machine_guid` | `string` | The 128-bit unique GUID Azure assigns the VM -- what licensing and inventory systems key on (stable across restarts, unlike the ARM ID across recreate) |
| `private_ip_address` | `string` | The primary private IP across the VM's attached NICs -- a convenience echo of the primary `AzureNetworkInterface`'s address |
| `public_ip_address` | `string` | The primary public IP across the attached NICs, when any IP configuration is fronted by one (empty for private-only VMs) |
| `computer_name` | `string` | The OS hostname the VM booted with |
| `system_assigned_identity_principal_id` | `string` | Principal (object) ID of the system-assigned identity, populated only when the identity type includes SYSTEM_ASSIGNED -- what `AzureRoleAssignment` grants reference |

## Related Components

- [AzureNetworkInterface](/docs/catalog/azure/network-interface) — the VM's network presence: subnet placement, public IPs, NSG filtering, accelerated networking
- [AzureManagedDisk](/docs/catalog/azure/managed-disk) — first-class data volumes that outlive the VM, and the golden-disk boot source
- [AzureUserAssignedIdentity](/docs/catalog/azure/user-assigned-identity) — shareable managed identities attached via `identity.identityIds`
- [AzureRoleAssignment](/docs/catalog/azure/role-assignment) — grants against the VM's system-assigned principal or its user-assigned identities
- [AzureKeyVault](/docs/catalog/azure/key-vault) — certificates installed onto the VM at provisioning time
- [AzureResourceGroup](/docs/catalog/azure/resource-group) — provides the resource group for VM placement

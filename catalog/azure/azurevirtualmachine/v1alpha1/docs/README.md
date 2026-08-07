# Azure Virtual Machine Deployment: The Compute Shell, Done Honestly

## Introduction: A VM Is Less Than You Think

Ask most engineers what an Azure Virtual Machine is and they will describe a
bundle: a machine with its network card, its public IP, its firewall rules,
its disks, its credentials. That mental model is convenient and wrong -- and
infrastructure built on it pays for the error at the worst possible times:
during recovery, during scale-out, and during the 3 AM incident where a VM
must be rebuilt without losing its data or its address.

In Azure Resource Manager's actual model, a virtual machine is a **compute
shell**: vCPUs, memory, an OS disk, and a set of *references* to other
first-class resources. The network interface is its own resource with its own
lifecycle. The managed disk is its own resource that exists before and after
any VM it attaches to. The managed identity is its own principal in Entra ID.
The VM composes them; it does not own them.

`AzureVirtualMachine` models exactly that. The VM is deliberately just the
machine. Everything it composes with is referenced, never created invisibly
inside it. This document explains why that decomposition is the right
production shape, and then walks the decisions that actually matter when you
run VMs against production: OS choice, image sources, spot economics,
placement, security posture, patching, and the two very different channels
for boot-time data.

## The Compute-Shell Model: Why NICs and Disks Are Referenced

### Network presence lives on the NIC

A VM's entire network posture -- which subnet it sits in, whether it has a
public IP, which network security group filters its traffic, whether
accelerated networking is on -- is a property of the **network interface**,
not the VM. Azure has always modeled it this way; the portal just hides it.

`AzureVirtualMachine` takes one or more referenced `AzureNetworkInterface`
resources via `networkInterfaceIds` (required, at least one; the first is
the primary). The consequences are worth internalizing:

- **The address outlives the machine.** Replace the VM -- new image, new
  size family, disaster recovery -- and the NIC, its private IP, its public
  IP, and its NSG associations are untouched. The replacement VM attaches to
  the same NIC and answers at the same address.
- **Security is auditable where it acts.** NSG filtering is declared on the
  NIC (and the subnet), so reviewing a workload's exposure never requires
  opening the VM resource at all.
- **Multiple NICs are first-class.** Appliances and split management/data
  planes attach several NICs; the VM size caps how many. Each is its own
  resource with its own addressing and filtering.

### Data lives on referenced managed disks

The same logic applies to storage. A data volume's natural lifetime is the
lifetime of the *data*, not of any particular machine reading it. Azure
models data disks as standalone managed-disk resources; attachment to a VM
is a separate operation.

`AzureVirtualMachine` mirrors this with `dataDiskAttachments`: each entry
references a first-class `AzureManagedDisk` and mounts it at a **LUN**
(logical unit number, 0-63) with a caching mode. On both provisioning
engines the attachment is realized as its own attachment resource -- not as
an inline block on the VM -- so:

- Detaching a disk is deleting a spec entry, never touching the VM or the
  disk.
- The disk, its snapshots, and its data survive VM replacement and deletion.
- The LUN is the stable identity the OS addresses the disk by
  (`/dev/disk/azure/scsi1/lun{n}` on Linux). Keep LUNs stable across
  changes; the OS's device naming depends on it.

Caching per attachment matters operationally: `READ_ONLY` suits read-heavy
data, `NONE` is required for disks larger than 4 TiB and is right for
write-heavy volumes like database logs.

### Identity is a principal, not a credential

The third referenced seam is identity. A VM's workload should authenticate
to Azure services without stored credentials -- that is what managed
identities are for. The spec's `identity` block carries either the
**system-assigned** identity (created and rotated by Azure with the VM; its
principal ID surfaces in the stack outputs) or referenced
`AzureUserAssignedIdentity` resources (identities you manage and share
across resources), or both. Permissions are then composed with
`AzureRoleAssignment` against the principal -- grants are their own
resources too, reviewable and revocable without touching the VM.

### The one deliberate exception: the OS disk

Only the OS disk is inline (`osDisk`). It is born and dies with the VM by
definition -- there is nothing composable about a boot volume that only ever
serves one machine. The single exception to the exception: a VM can boot
from an **existing** referenced OS disk (`osManagedDiskId`), covered below.

## Linux XOR Windows: Two Surfaces, One Honest Choice

ARM does not model "a VM with an OS setting." It models **two separate
management surfaces** -- Linux VMs and Windows VMs -- with different
authentication contracts, different patch-mode vocabularies, and different
OS-management knobs. Any API that flattens them into one shape either
invents fields Azure doesn't have or silently drops half the surface.

`AzureVirtualMachine` keeps the split explicit: `osProfile` carries exactly
one of `linux` or `windows`, and the module deploys the matching ARM
resource.

### The Linux contract: SSH-first

- `adminUsername` plus at least one SSH public key is the production path.
- `disablePasswordAuthentication` defaults to **true** -- keys only. This is
  Azure's own default and the correct posture; a password is only accepted
  when you explicitly set the flag to false and provide `adminPassword`.
- SSH public keys are public material -- safe in a manifest, unlike any
  password.
- Patch modes use ARM's Linux vocabulary: `LINUX_IMAGE_DEFAULT` (whatever
  the image's own update configuration does) or
  `LINUX_AUTOMATIC_BY_PLATFORM` (Azure Update Manager orchestrates).
- `licenseType` covers bring-your-own-subscription for commercial distros:
  RHEL and SLES BYOS/SAP variants, Ubuntu Pro.

### The Windows contract: password + management surface

- `adminUsername` and `adminPassword` are both required when booting from an
  image; ARM enforces password complexity (8-123 characters, 3 of 4
  classes). The password is secret material -- source it from a secret
  reference, never a manifest literal.
- Windows carries the management knobs Linux has no equivalent for:
  `automaticUpdatesEnabled` (default true), `hotpatchingEnabled` (reboot-less
  security updates, supported Windows Server Azure Edition images only),
  `timezone`, WinRM listeners (HTTPS listeners reference a Key Vault
  certificate), and raw unattend.xml fragments for pre-agent bootstrap
  (treated as secret -- AutoLogon carries the admin password).
- Patch modes use ARM's Windows vocabulary: `MANUAL`, `AUTOMATIC_BY_OS`
  (Azure's default), or `WINDOWS_AUTOMATIC_BY_PLATFORM`.
- `licenseType: WINDOWS_SERVER` / `WINDOWS_CLIENT` is Azure Hybrid Benefit:
  bring an existing license and stop paying the image's Windows price --
  one of the largest single cost levers on a Windows fleet.

One naming trap: Windows computer names cap at **15 characters** and default
to the VM name. A longer VM name needs an explicit `osProfile.computerName`.

## Image Sources: Exactly One of Three

Every VM boots from exactly one image source, and the spec enforces the
choice up front rather than letting ARM reject it at deploy time.

**1. Marketplace coordinates** (`sourceImageReference`) -- publisher, offer,
SKU, version. The workhorse. One version behavior to burn in: `latest`
resolves **at creation only**. The VM does not follow new image releases
afterward; two VMs created a month apart from `latest` can run different
builds. Pin an explicit version for reproducible fleets.

**2. Custom or gallery image** (`sourceImageId`) -- a managed image or a
Shared Image Gallery image/version by ARM ID (community and direct-shared
gallery IDs included). The golden-image pipeline path: bake once, stamp
many.

**3. An existing OS disk** (`osManagedDiskId`) -- the disk-swap /
golden-disk path, and the most operationally interesting one. The VM boots
from a referenced `AzureManagedDisk` that already carries an operating
system: a disk restored from a snapshot, detached from a dead VM, or
prepared as a golden boot volume. Two rules follow from "the disk already
contains its users":

- The OS profile still selects the OS (Linux or Windows) but must carry
  **no authentication fields** -- no username, no password, no SSH keys.
  Spec-level validation enforces this pairing; the users that exist are the
  ones on the disk.
- Patching stays at the image default; provisioning-time patch orchestration
  assumes an image boot.

This is the recovery path that makes the compute-shell model pay off in
full: NIC survives, data disks survive, and now even the boot volume can
outlive its machine.

## Spot Capacity: Discounts with a Contract

Spot VMs run on Azure's spare capacity at discounts commonly in the 60-90%
range. The contract: Azure may **evict** the VM whenever it needs the
capacity back, with roughly 30 seconds of notice. Presence of the `spot`
block makes the VM a spot instance; absence means regular on-demand. The
setting is fixed at creation.

Two decisions define spot behavior:

**Eviction policy.** `DEALLOCATE` stops the VM: compute billing stops, the
disks persist, and the VM can restart when capacity returns -- right for
workers that checkpoint to a referenced data disk and resume. `DELETE`
removes the VM and its disks -- only for fully stateless fleets managed by
an orchestrator that replaces members automatically.

**Max bid price.** The default (-1) pays up to the on-demand price and is
never evicted on *price* -- only on capacity. Setting a dollar cap adds
price-based eviction; do it only when cost predictability genuinely beats
availability.

Pair spot with `terminationNotification` (up to `PT15M`): Azure emits a
scheduled event before termination, the workload polls the Instance Metadata
Service, sees the signal, and drains -- checkpointing to a data disk that,
because it is a referenced first-class resource, will still be there when
the replacement worker attaches it.

Spot is for interruption-tolerant workloads: batch, rendering, CI runners,
queue consumers. It is not a discount button for production web servers.

## Placement: Where the VM Runs Relative to Failure

The `availability` block positions the VM relative to Azure's fault
machinery. The options are not interchangeable; they answer different
questions.

| Option | Question it answers | Notes |
|--------|--------------------|-------|
| `zone` (`"1"`/`"2"`/`"3"`) | "Which datacenter building?" | The modern resilience unit. Zonal NIC public IPs and zonal disks must match the zone. Conflicts with availability sets. Fixed at creation. |
| `availabilitySetId` | "Spread across racks, pre-zones style?" | The classic fault/update-domain grouping. Prefer zones in zoned regions. |
| `proximityPlacementGroupId` | "As physically close as possible?" | Minimal inter-VM latency for HPC and chatty clusters -- the opposite trade of zone spreading. |
| `capacityReservationGroupId` | "Is capacity guaranteed?" | Consumes reserved capacity for burst/DR events -- capacity insurance, not fault isolation. |
| `dedicatedHostId` / `dedicatedHostGroupId` | "Whose hardware?" | Single-tenant physical isolation for compliance/licensing; pin a host or let Azure pick within the group (mutually exclusive). |
| `virtualMachineScaleSetId` (+ `platformFaultDomain`) | "Managed fault spreading for an individual VM?" | Attaches the VM to a FLEXIBLE-orchestration scale set (value-or-ref against AzureVirtualMachineScaleSet's `scale_set_id` output); optionally pin the fault domain. Fixed at creation. |

The spec enforces the real conflicts up front: zone XOR availability set;
capacity reservations cannot combine with availability sets or proximity
groups; a fault domain requires the scale-set attach. Leave the whole block
unset for regional placement with no constraints -- fine for dev, not for
anything with an availability target.

For real availability, run **two or more VMs in different zones behind a
load balancer**. One zonal VM has better isolation than one regional VM,
but it is still one machine.

## Security Posture: Trusted Launch vs Confidential VM

Two different threat models, often conflated:

**Trusted launch** (`security.secureBootEnabled` + `vtpmEnabled`) defends
the **boot chain**: UEFI secure boot ensures only signed components load;
the virtual TPM enables measured boot and attestation. It targets bootkits
and rootkits, costs nothing extra, and is the right default posture for
production VMs on Gen2 images. Both flags are fixed at creation -- decide
them on day one, because retrofitting means replacement.

**Confidential VMs** (`osDisk.securityEncryptionType`) defend against the
**host itself**: hardware-based encryption of the VM's guest state so that
even Azure's hypervisor operators cannot read memory or state.
`VM_GUEST_STATE_ONLY` encrypts the guest-state blob;
`DISK_WITH_VM_GUEST_STATE` also encrypts the OS disk. The spec enforces
ARM's prerequisites as paired validations: both modes require `vtpmEnabled`,
and the disk variant additionally requires `secureBootEnabled`. Confidential
VMs need confidential-capable sizes (DCasv5/ECasv5 families) and, for
customer-key guest-state encryption, `secureVmDiskEncryptionSetId` on the
OS disk.

Orthogonal to both: `encryptionAtHostEnabled` encrypts data on the compute
host itself, closing the gap platform encryption leaves -- temp disks and
disk caches. The subscription must have the `EncryptionAtHost` feature
registered. For most production fleets, trusted launch + encryption at host
is the sensible baseline; confidential VMs are for regulated workloads whose
threat model includes the cloud operator.

## Patch Orchestration: Who Drives Updates

Patching splits across the spec exactly where ARM splits it: the per-OS
patch **mode** lives in the OS profile (because the vocabularies differ per
OS), while the shared dials live in `patching`.

The mode question is: who orchestrates? `LINUX_IMAGE_DEFAULT` /
`AUTOMATIC_BY_OS` leave it to the image or Windows Update -- Azure's
defaults, reasonable for fleets with their own patch tooling. Setting the
mode to AUTOMATIC_BY_PLATFORM (per-OS value: `LINUX_AUTOMATIC_BY_PLATFORM` /
`WINDOWS_AUTOMATIC_BY_PLATFORM`) hands orchestration to **Azure Update
Manager**, which unlocks the shared dials:

- `patching.rebootSetting` -- when platform patching may reboot: `ALWAYS`,
  `IF_REQUIRED`, or `NEVER` (patches needing a reboot then wait for a manual
  one). Gated on AUTOMATIC_BY_PLATFORM, and the spec enforces the gate.
- `patching.bypassPlatformSafetyChecksOnUserScheduleEnabled` -- lets
  customer-scheduled patching bypass certain platform safety checks. Same
  gate.
- On Windows, AUTOMATIC_BY_PLATFORM is also the prerequisite for
  `hotpatchingEnabled`.

Independent of the mode, `patching.assessmentMode:
ASSESSMENT_AUTOMATIC_BY_PLATFORM` has Azure *assess* pending patches daily
-- visibility without ceding control of installation. Turning assessment on
fleet-wide is cheap and pays for itself the first time an auditor asks which
machines are behind.

## custom_data vs user_data: The Distinction That Prevents a Leak

Two fields deliver data to the VM, and they have opposite security
contracts. Confusing them is how bootstrap tokens end up readable by every
process on the machine.

**`customData`** is the cloud-init / provisioning channel: base64-encoded,
delivered **once at first boot**, not retrievable afterward through the
platform. Because bootstrap scripts routinely embed join tokens and initial
secrets, the platform treats the whole field as **secret material**, and it
is fixed at creation -- changing it replaces the VM (first-boot data is
meaningless to a machine that already booted).

**`userData`** is the runtime metadata channel: base64-encoded, readable
from *inside* the VM via the Instance Metadata Service **at any time**, by
any process that can reach `169.254.169.254` -- no privileges required. It
is updatable in place and readable back through ARM. That reachability is
the whole point (configuration a running workload re-reads) and the whole
constraint: **never put secrets in `userData`**.

The rule: secrets and one-shot bootstrap go in `customData`; anything a
running process should re-read goes in `userData`; long-lived secrets belong
in a vault fetched via the VM's managed identity -- which is exactly the
composition `identity` + `AzureRoleAssignment` exists for.

## Operational Lifecycle Notes

What replaces, what reboots, what updates in place:

- **Replacement** (the VM's identity): name, region, zone, image source,
  admin credentials, `customData`, and the security/confidential posture.
  The inline OS disk is replaced with it; NICs and data disks survive --
  which is exactly why they are referenced.
- **Reboot in place**: resizing (`size`); moving size families may
  deallocate first.
- **In place**: data-disk attach/detach, tags, `userData`, Windows license
  type, gallery applications.
- **Fixed at creation**: spot settings, ephemeral OS disk, edge zone,
  `provisionVmAgent`.

Two more defaults worth a deliberate look: `bootDiagnostics: {}` enables
serial-console output and boot screenshots on Azure-managed storage -- the
first tool for debugging a VM that will not boot, with no storage account to
operate; and the ephemeral OS disk (`osDisk.diffDiskSettings`) puts the boot
volume on the VM's local storage -- free and fast, but **wiped on every
stop/deallocate**, so it belongs only under stateless, image-driven fleets.

## The Planton Approach

The spec encodes ARM's real contracts as compile-time-checked validations
rather than deploy-time surprises: exactly one image source; exactly one OS
profile; Linux auth requirements (keys when password auth is disabled, the
default); Windows auth requirements; the no-auth rule when booting from an
existing disk; zone XOR availability set; the patch-mode gates on reboot
setting and safety-check bypass; and the vTPM/secure-boot prerequisites for
confidential-VM encryption.

Every composable seam is a typed reference with a sensible default target:
`networkInterfaceIds` resolves an `AzureNetworkInterface`'s
`network_interface_id` output, `dataDiskAttachments[].managedDiskId` and
`osManagedDiskId` resolve an `AzureManagedDisk`'s `disk_id`,
`identity.identityIds` resolve `AzureUserAssignedIdentity`, and
`secrets[].keyVaultId` resolves an `AzureKeyVault`. Fields with no
first-class kind yet (availability sets, proximity groups, dedicated hosts,
disk encryption sets) take plain ARM IDs -- stated plainly rather than
papered over.

Both provisioning engines deploy the same shape at behavioral parity, and
only explicit spec choices are ever sent -- an unspecified field and Azure's
default deploy identically. Tags merge user values over the Planton-derived
metadata tags, user wins on collision, so organizational governance
conventions hold.

A minimal honest production VM:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureVirtualMachine
metadata:
  name: app-vm
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
        - publicKey: ssh-ed25519 AAAA...
  osDisk:
    caching: READ_WRITE
    storageAccountType: PREMIUM_LRS
  sourceImageReference:
    publisher: Canonical
    offer: ubuntu-24_04-lts
    sku: server
    version: latest
  availability:
    zone: "1"
  security:
    secureBootEnabled: true
    vtpmEnabled: true
  identity:
    type: SYSTEM_ASSIGNED
  bootDiagnostics: {}
```

Zonal, trusted-launch, SSH-key-only, credential-less toward Azure services,
debuggable when it will not boot -- and every piece it composes with (the
NIC, and through it the subnet, public IP, and NSG; any data disks; any
role grants) is a first-class resource with its own lifecycle.

## Conclusion: Compose the Machine

The compute-shell model is not an abstraction preference; it is how Azure
actually works, surfaced instead of hidden. Model the VM as just the
machine, and replacement stops being scary: addresses persist on NICs, data
persists on disks, permissions persist on principals, and even the boot
volume can be swapped. Model it as a bundle, and every one of those
lifetimes is accidentally chained to the shortest-lived resource in the
system.

Choose the OS surface explicitly. Pin image versions for fleets. Give
interruption-tolerant work to spot with a drain window and a persistent
checkpoint disk. Pick placement for the failure you are defending against.
Turn on trusted launch day one. Decide who orchestrates patching. Keep
secrets out of `userData`.

And let the machine be just the machine.

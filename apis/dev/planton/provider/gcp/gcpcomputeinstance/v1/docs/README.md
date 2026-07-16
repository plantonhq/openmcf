# GcpComputeInstance — Deep Dive

## The problem this resource solves

The virtual machine is the escape hatch of every cloud: when a workload does not fit a managed runtime — a database that needs local NVMe, a licensed appliance, a GPU trainer, a bastion host — it lands on a VM. Compute Engine's instance API is correspondingly enormous, and most of its sharp edges (Spot's legacy preemptible coupling, live-migration restrictions for GPUs and confidential VMs, which properties silently force a replacement) are learned in production. This kind models the instance so those edges are taught up front: invalid combinations are rejected before deploy, replacement-inducing fields are documented as create-time decisions, and the resources a VM composes with — disks, addresses, networks, identities — are explicit references to first-class kinds rather than strings pasted between consoles.

## Where it sits in the composition

A VM is a compute node; almost everything durable about it lives in *other* resources it references:

- **GcpComputeDisk** — the boot disk can *be* one (`bootDisk.sourceDisk`), and every attached data disk *is* one (`attachedDisks[].source` → `status.outputs.self_link`). The disk owns its size, type, encryption, and lifecycle; the instance merely attaches it. A database VM's data survives the VM because the disk is a separate node in the graph.
- **GcpAddress** — static internal IPs (`networkInterfaces[].networkIp`) and static external IPs (`accessConfigs[].natIp`) reference a reserved address's `status.outputs.address`. The IP outlives the VM, so DNS and firewall allowlists stay stable across rebuilds.
- **GcpVpcNetwork / GcpSubnetwork** — each vNIC attaches by reference to a network (auto-mode VPCs) or a subnetwork (custom-mode VPCs, the production norm).
- **GcpServiceAccount** — the VM's workload identity (`serviceAccount.email` → `status.outputs.email`).
- **GcpKmsKey** — CMEK for the boot disk and for attached-disk attachments (`kmsKey` → `status.outputs.key_id`).
- **GcpProject** — the owning project (`projectId`).

Outbound, consumers reference the instance's `self_link`, `internal_ip`, and `external_ip` outputs — instance-group membership, DNS records, monitoring targets.

## Lifecycle contract

| Property | Behavior |
|---|---|
| `zone`, boot source, NIC count and their networks, `scratchDisks`, `hostname`, `confidentialInstanceConfig`, `reservationAffinity`, `resourceManagerTags`, `canIpForward` | Immutable (ForceNew) — changing them replaces the VM |
| `machineType`, `serviceAccount`, `shieldedInstanceConfig`, and several others | Updated by stopping and restarting the VM — the provider only does this when `allowStoppingForUpdate` is true; otherwise the update fails instead of causing surprise downtime |
| `desiredStatus` | Starts (`RUNNING`), suspends (`SUSPENDED`), or stops (`TERMINATED`) the VM in place. Compute billing stops while terminated; disks keep billing |
| `deletionProtection` | Guards the VM object only — deleting a protected instance fails until it is set back to false. Defaults to false because the data levers live elsewhere: the boot disk's `autoDelete` and each `GcpComputeDisk`'s own configuration |
| Boot disk `sizeGb` | Grows in place; never shrinks |

## Design decisions

### Exactly one boot source

A boot disk comes from exactly one of three places, and the spec enforces the XOR before deploy: `image` (fresh install — a family like `debian-cloud/debian-12` resolves to the newest image at create time), `sourceSnapshot` (restore), or `sourceDisk` (a pre-created bootable `GcpComputeDisk`, referenced like any other disk). When booting from an existing disk, size/type/encryption are ignored — the disk already owns them.

### Disks are references, not inline blocks

Attached data disks are deliberately *not* an inline create-a-disk convenience. Each is a `GcpComputeDisk` with its own spec, its own protections, and its own lifecycle; the instance holds only the attachment (`source`, `deviceName`, `mode`, and the matching `kmsKey` when the disk is CMEK-encrypted). This is what makes stateful VMs safe: destroying or replacing the VM never implies destroying the data.

### Spot is a single switch

GCP's API still requires Spot VMs to set the legacy `preemptible` flag and to disable automatic restart. The spec exposes only the modern lever — `scheduling.provisioningModel: SPOT` — and both engines derive the rest identically: preemptible is set, automatic restart is forced off. `instanceTerminationAction` chooses what preemption does: `STOP` keeps the stopped VM and its disks; `DELETE` removes the VM. A CEL rule rejects a termination action without Spot, and another rejects `automaticRestart: true` with Spot, so misconfigurations fail at validation rather than at the API.

### Maintenance coupling is validated

Confidential VMs and VMs with GPUs cannot live-migrate, so both require `scheduling.onHostMaintenance: TERMINATE`. The spec enforces both couplings pre-deploy — the two most common instance-creation errors in this space become manifest validation messages instead of opaque API rejections.

### Identity falls back sensibly

`instanceName` falls back to `metadata.name`, computed identically by both engines. `projectId` empty falls back to the provider's default project. An omitted `serviceAccount` block uses the Compute Engine default service account with its default scopes — fine for sandboxes, not for production (see best practices).

### Metadata surfaces stay distinct

`startupScript` rides the dedicated startup-script surface (re-runs on every boot), never plain metadata. `sshKeys` fold into the metadata `ssh-keys` key, newline-joined, byte-identically on both engines — and are ignored by VMs using OS Login, which is the recommended access path anyway.

## Coverage

The spec models the full GA surface of the instance resource: boot/attached/scratch disks, multi-NIC networking with IPv6 and alias ranges, Spot and standard scheduling with run-duration limits, sole tenancy, Shielded and Confidential VM, GPUs, advanced machine features, reservation affinity, egress bandwidth tiers, resource-manager tags, and the stop/suspend/start lever. The remaining tail is a handful of deliberate exclusions, each recorded below with its reason — nothing is missing by accident.

## Deliberately not modeled

- **Dynamic-NIC fields** (`igmp_query`, `vlan`, `parent_nic_name`) — absent from the released google `~> 6.x` provider line; revisit when the floor moves.
- **`partner_metadata` and `scheduling.graceful_shutdown`** — beta-only surfaces; the kind models the GA provider.
- **Raw CSEK encryption keys** (`disk_encryption_key_raw` and friends) — customer-supplied keys put secret material in manifests and state. CMEK via a `GcpKmsKey` reference is the modeled encryption path.
- **`network_attachment` (Private Service Connect interfaces)** — a composition-worthy surface that deserves a first-class kind of its own rather than a raw URL string here.
- **Legacy preemptible-without-Spot VMs** — the pre-Spot product. Spot supersedes it with the same discount and fewer restrictions, so only `provisioningModel: SPOT` is modeled.

## Best practices the presets encode

- **Dedicated least-privilege service account** — never ship production on the Compute Engine default account. Reference a `GcpServiceAccount` and grant it exactly the IAM roles the workload needs, with the single `cloud-platform` OAuth scope so IAM (not scopes) is the control plane.
- **Shielded VM on by default posture** — `enableSecureBoot: true` costs nothing on Google-provided images and blocks boot-chain tampering; vTPM and integrity monitoring are GCP defaults worth keeping.
- **No external IP unless the VM is deliberately public** — private VMs with Cloud NAT for egress, OS Login (`enable-oslogin` metadata) plus IAP for operator access.
- **Spot for stateless, interruption-tolerant work** — dev boxes, CI runners, batch workers. Pair with `instanceTerminationAction: DELETE` for fully disposable capacity, or `STOP` when the disks should survive preemption.
- **`allowStoppingForUpdate: true`** wherever a brief restart is acceptable — without it, machine-type and service-account changes fail instead of applying.
- **Stateful VMs**: `deletionProtection: true` on the instance, data on referenced `GcpComputeDisk` resources, `onHostMaintenance: MIGRATE` for zero-downtime host maintenance.

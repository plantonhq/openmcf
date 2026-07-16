# Azure Virtual Machine Scale Set: One Resource, Two Philosophies of Fleet

## Introduction: The Fleet Is the Unit

A virtual machine scale set inverts the single-VM mental model. With a VM, the machine is the unit of management: you name it, patch it, attach its disks, and mourn it when it dies. With a scale set, the **template** is the unit: you declare what every instance looks like — image, size, OS profile, disks, network shape — and fleet controls decide how many exist, where they spread, how they upgrade, and what happens when one becomes unhealthy. Instances are anonymous stampings of the template.

This document explains how Azure scale sets actually work — the orchestration-mode split, the template-versus-reference boundary, upgrade machinery, spot economics, and health — and records the design decisions behind this component's shape.

## One Kind, an Explicit Mode

ARM has exactly **one** scale-set resource type (`Microsoft.Compute/virtualMachineScaleSets`) with an orchestration-mode property. Some toolchains split it into separate per-mode (and per-OS) resource types; that split is tooling ergonomics, not Azure's model. This component models the ARM reality: one kind, an explicit `orchestrationMode` discriminator (FLEXIBLE by default — Azure's recommendation for new workloads), and an explicit `osProfile.linux` XOR `osProfile.windows` discriminator, with validation gating what each mode supports.

The mode split is real, not cosmetic — the capability sets genuinely bifurcate:

**UNIFORM-only capabilities** (the classic large-fleet engine): overprovisioning (provision extras during scale-out, keep the first healthy ones — free, because overprovisioned instances are not billed), scale-in selection policy, automatic OS image upgrades, spot restore, gallery applications, trusted launch (secure boot + vTPM) and confidential disk encryption, the load-balancer health probe, dedicated host groups, edge zones, and pool-style NAT-rule references in the NIC template.

**FLEXIBLE-only capabilities** (the resilient VM group): mixed-SKU profiles (`skuName: Mix` plus up to five candidate sizes and an allocation strategy — the capacity-resilience shape), spot/on-demand priority mixing (a guaranteed on-demand base plus a ratio above it), per-OS patch orchestration (`patchMode`, `patchAssessmentMode`, Windows hotpatching), the `networkApiVersion` selector, and standalone-VM attach — an `AzureVirtualMachine` can join the set's fault-spreading via `availability.virtualMachineScaleSetId`.

Two structural differences follow from the mode and are enforced as validation: FLEXIBLE sets **require** `platformFaultDomainCount` (1 when zones are the resilience unit, the region's max for regional spreading), and FLEXIBLE sets support **only user-assigned** managed identities (ARM's contract; UNIFORM supports system-assigned too).

Every one of these gates is a spec-level CEL rule with a message naming the contract — a `scaleIn` policy on a FLEXIBLE set fails validation immediately, not twenty minutes into a deploy.

## Template Versus Reference: Where the Line Is

The single most important boundary in the spec, and it deliberately differs from the standalone `AzureVirtualMachine`:

**Inline (template):** each instance's NICs, OS disk, data disks, and optional per-instance public IPs. These are declared inside the scale set because every instance stamps its own copies and they live and die with the instance. There is nothing to reference — the resources do not exist until an instance does, and there can be a thousand of them.

**Referenced (shared):** everything with an independent lifecycle that instances plug into — the subnet (`AzureSubnet`), load-balancer backend pools and NAT rules (the `AzureLoadBalancer`'s name-keyed `backend_pool_ids` / `nat_rule_ids` map outputs, e.g. `status.outputs.backend_pool_ids.web`), the NSG (`AzureNetworkSecurityGroup`), user-assigned identities (`AzureUserAssignedIdentity`), public IP prefixes (`AzurePublicIpPrefix`), Key Vaults for provisioning certificates and extension secrets (`AzureKeyVault`), and the UNIFORM health probe (the LB's `probe_ids` map).

The standalone VM references first-class NICs and disks because one machine's network identity and data have lives of their own. A fleet's per-instance NICs and disks have no independent life by definition — inlining them is honest composition, not a bundling shortcut.

## The Upgrade Machinery

`upgradePolicy.mode` decides how template changes reach existing instances:

- **MANUAL** (Azure's default): changes apply to NEW instances only; existing ones wait for an explicit upgrade. Safe, and surprising — a fleet can run two template versions indefinitely.
- **AUTOMATIC**: all instances upgrade immediately. Brief fleet-wide disruption; acceptable for dev.
- **ROLLING**: instances upgrade in health-checked batches — the production choice. Requires a `rolling` policy (batch percent, unhealthy thresholds, pause between batches, and optional zone-by-zone ordering, unhealthy-first prioritization, or surge-instead-of-in-place) and **health monitoring**.

Health monitoring is the load-bearing prerequisite of the whole health stack. Rolling upgrades, automatic instance repair, platform patch orchestration, and hotpatching all key on instance health, which comes from exactly two sources: the **application health extension** (`ApplicationHealthLinux` / `ApplicationHealthWindows`, declared in `extensions[]` with the port/path it probes) or, on UNIFORM sets only, a referenced **load-balancer health probe**. The spec enforces each dependency chain — ROLLING without health monitoring, repair without health, hotpatching without platform patching + the Windows health extension: all fail validation with the contract named.

**Automatic OS image upgrades** (UNIFORM) close the loop for image-driven fleets: with `sourceImageReference.version: latest`, Azure rolls new image releases across the fleet using the same rolling machinery — the fleet stays current with zero pipeline work, and automatic rollback restores the previous image if a batch fails health checks.

## Spot Economics

Presence of the `spot` message makes the fleet evictable and deeply discounted. The decisions:

- **`evictionPolicy`** (required): `DELETE` removes evicted instances and their disks — the usual fleet choice, letting count-based scaling replace capacity; `DEALLOCATE` stops them (disks persist, billing stops) so they can be restored.
- **`maxBidPrice`**: `-1` (default) pays up to the on-demand price and is never evicted on price grounds; a lower cap trades availability for a hard cost ceiling.
- **UNIFORM `restore`**: Azure automatically tries to restore evicted (deallocated) instances when capacity returns.
- **FLEXIBLE `priorityMix`**: the middle ground for fleets that cannot be all-spot — `baseRegularCount` instances are guaranteed on-demand, and `regularPercentageAboveBase` sets the on-demand ratio above the base. An all-spot batch tier with a guaranteed two-instance on-demand core is four lines of YAML.

Spot fleets should enable `terminationNotification` and drain on the scheduled event — eviction arrives with only the notification lead time (up to 15 minutes).

## Recorded Design Decisions

- **One kind with an `orchestrationMode` discriminator** rather than per-mode kinds: ARM has one resource type; the per-mode split elsewhere is tooling ergonomics. Same reasoning as the VM's single kind with a linux/windows discriminator.
- **FLEXIBLE is the default** (unspecified = FLEXIBLE): Azure's own recommendation for new workloads. UNIFORM stays fully modeled — its exclusive capabilities are real and an advanced org reaches them.
- **`osProfile.linux` XOR `osProfile.windows`**: an explicit OS discriminator; per-OS vocabularies (patch modes especially) stay in their own enums with ARM's exact values.
- **Per-instance NICs/disks inline, shared resources referenced** — the template-versus-reference line above.
- **Load-balancer composition through name-keyed map outputs**: pool membership, NAT-rule references, and the UNIFORM health probe all reference the `AzureLoadBalancer`'s `backend_pool_ids` / `nat_rule_ids` / `probe_ids` maps by sub-resource name (`fieldPath: status.outputs.backend_pool_ids.web`).
- **`overprovision` carries no default annotation** despite Azure defaulting it true on UNIFORM sets: the loader would materialize the default onto FLEXIBLE sets, where the field is illegal. Same reasoning for `doNotRunExtensionsOnOverprovisionedMachines`. The UNIFORM behavior is unchanged — unset means Azure's default (true).
- **`instances` is optional**: unset lets the platform manage the count, which is correct when an autoscale resource owns it; a fixed-size fleet sets it explicitly.
- **Extensions fold inline; no standalone scale-set-extension kind**: an extension has no life without its scale set and is not FK-referenced by anything. The provider's standalone extension resource exists for out-of-band attachment to Uniform sets — a workflow, not a composition seam.
- **`customData` is sensitive; `userData` is not**: cloud-init routinely embeds bootstrap secrets and is write-only; user data is IMDS-readable back from every instance by design, so treating it as secret would be a false promise — the field comment warns instead.
- **Windows `additionalUnattendContents[].content` and extension `protectedSettings` are sensitive**: AutoLogon fragments carry the admin password; protected settings carry credentials by definition.
- **Skipped, with reasons**: the deprecated standalone NAT-pool template surface (pool-style NAT rules on the load balancer superseded it); ResiliencyPolicy toggles (`resilientVmCreation/DeletionEnabled` — preview-adjacent, no composition seam); the legacy no-`skuName` orchestrated shape (the provider keeps it only for pre-2022 state compatibility).
- **Placement references are plain ARM IDs** (proximity placement groups, capacity reservation groups, host groups): those kinds are not modeled yet; the fields upgrade to references when they are.

## The Field-Coverage Floor

The spec covers the full azurerm v4.80 surface across all three provider resources (Uniform Linux, Uniform Windows, orchestrated/Flexible), unified under the mode and OS discriminators: identity, SKU + instances + mixed-SKU profiles, both OS profiles with per-OS patch vocabularies, OS/data disk templates (incl. ephemeral OS disks, UltraSSD/PremiumV2 dialed performance, CMK and confidential encryption), NIC + IP-configuration templates (accelerated networking, IP forwarding, auxiliary NVA acceleration, per-instance public IPs with prefixes and IP tags), upgrade policy + rolling + automatic OS upgrade, spot (eviction, bid, restore, priority mix), automatic instance repair, termination notification, extensions (incl. Key-Vault-sourced protected settings and provision ordering), boot diagnostics, secrets (Key Vault certificates), plan, gallery applications, placement (PPG, capacity reservation, host group, single placement group), zones + zone balance + fault domains, security (trusted launch, encryption at host), scale-in, capabilities (Ultra SSD), edge zone, and tags.

The conditional-contract matrix the provider enforces in Go is mirrored as ~30 spec-level CEL rules, so a manifest that would fail ARM fails validation first, with the contract named.

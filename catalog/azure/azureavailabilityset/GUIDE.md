# Azure Availability Set -- Operational Guide

Judgment calls that matter when you run availability sets in production.

## Zones first; sets where zones cannot go

An availability zone survives a DATACENTER failure; an availability set survives a rack/power/switch failure inside one datacenter. In zoned regions, new architectures should pin VMs to zones and skip the set entirely (a VM cannot have both). The set earns its place in regions without zones, in classic topologies being lifted as-is, and for tiers whose tooling predates zonal thinking. Do not retrofit sets onto zonal designs for extra credit -- they compose with zones not at all.

## The set is immutable; plan the domains once

Every setting except tags is create-only, and VMs join only at THEIR creation -- so the set's shape is decided before the first VM exists. The provider defaults (5 update domains, 3 fault domains, managed alignment on) are right for almost everything; lower fault domains only when the target region cannot provide 3 (Azure rejects the count at create). Rebuilding a set later means recreating every VM in it -- treat the set like a subnet, not like a tag.

## Two VMs minimum, or the set is theater

The classic 99.95% SLA starts at TWO VMs in the set -- a single-VM availability set provides exactly nothing except a future constraint. If a tier will only ever have one VM, skip the set; if it has two or more, put them ALL in (a tier half-in, half-out fails together anyway through the half that shares hardware).

## Managed alignment is not optional in practice

`managed: true` (the default) aligns fault domains with managed-disk storage so a storage-cluster failure does not cross your compute fault domains. The false setting exists for the unmanaged-disk era; every managed-disk VM (that is, every VM the catalog creates) belongs in a managed set. Leave the field unset.

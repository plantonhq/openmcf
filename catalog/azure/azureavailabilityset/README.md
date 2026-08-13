# Overview

The **AzureAvailabilitySet** component deploys an availability set -- the classic pre-zones placement grouping that spreads VMs across separate fault domains (power/network/rack) and update domains (planned-maintenance batches) so one hardware failure or maintenance window cannot take them all down.

## Purpose

- **Fault isolation without zones**: in regions without availability zones -- and for classic lift-and-shift topologies -- the availability set is the placement tool that keeps a multi-VM tier from sharing a single point of hardware failure.
- **The 99.95% SLA anchor**: two or more VMs in an availability set carry Azure's classic multi-VM SLA.
- **Join at creation**: VMs reference the set when they are created (AzureVirtualMachine's `availability.availability_set_id`); membership cannot change in place.

## Key Features

- Full azurerm v5 surface: fault/update domain counts, managed-disk alignment, proximity placement group, tags.
- Chart-ready: `resource_group` defaults its reference to AzureResourceGroup, and the `availability_set_id` output is exactly what AzureVirtualMachine's placement field references.
- Provider defaults ride through: unset domain counts and `managed` take the provider's values (5 update domains, 3 fault domains, managed alignment on).

## Use Cases

- **Classic web/app tiers**: two or more VMs behind a load balancer, spread across fault domains.
- **Regions without zones**: the only placement-level fault isolation available.
- **Latency-sensitive clusters**: pair the set with a proximity placement group for co-location plus fault spreading.

## Future Enhancements

- Proximity placement groups are a P2 catalog kind; the `proximity_placement_group_id` reference gains its typed default when that kind lands.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).

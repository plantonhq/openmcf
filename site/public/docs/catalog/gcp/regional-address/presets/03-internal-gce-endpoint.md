---
title: "Internal GCE Endpoint"
description: "This preset reserves a regional internal IP address with the `GCE_ENDPOINT` purpose within a specific subnetwork — the address type used for VM instances, alias IP ranges, and similar compute..."
type: "preset"
rank: "03"
presetSlug: "03-internal-gce-endpoint"
componentSlug: "regional-address"
componentTitle: "Regional Address"
provider: "gcp"
icon: "package"
order: 3
---

# Internal GCE Endpoint

This preset reserves a regional internal IP address with the `GCE_ENDPOINT` purpose within a specific subnetwork — the address type used for VM instances, alias IP ranges, and similar compute endpoints.

## When to Use

- VM instances that need a static internal IP in a specific subnetwork
- Alias IP ranges attached to VM network interfaces
- Any compute endpoint that requires a reserved private IP within a subnet's CIDR range

## Key Configuration Choices

- **INTERNAL address type** (`addressType: INTERNAL`) — reserves a private IP within your VPC
- **GCE_ENDPOINT purpose** — tells GCP this address is for a compute endpoint (validated: purpose requires INTERNAL)
- **Required subnetwork** (`subnetwork`) — GCE_ENDPOINT addresses must be scoped to a subnetwork; validated at the schema level
- **Required region** (`region`) — must match the subnetwork's region
- **Optional specific `address`** — set to reserve a particular IP within the subnetwork's range; omit to let GCP assign one

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<gcp-project-id>` | GCP project ID | GCP Console or `GcpProject` outputs |
| `<your-address-name>` | Name for this address resource | Choose a descriptive name (e.g., `vm-static-internal-ip`) |
| `<your-region>` | GCP region | Must match the subnetwork's region |
| `<subnetwork-self-link-or-name>` | Subnetwork for the reservation | `GcpSubnetwork` outputs (`subnetwork_self_link`) or GCP Console |

## Related Presets

- **01-external-nat-ip** — External static IP for Cloud NAT
- **02-internal-lb-vip** — Internal shared load balancer VIP

## Related Components

- [GcpSubnetwork](/docs/catalog/gcp/gcpsubnetwork) — provides the subnetwork referenced by `subnetwork`
- [GcpVpcNetwork](/docs/catalog/gcp/gcpvpcnetwork) — parent network for the subnetwork

# Internal (Private-Only) Gateway

This preset creates an east-west L7 gateway with no public exposure: a
single private frontend pinned to a static address in the gateway subnet,
fixed two-instance capacity, and connection draining so instances rotate
out without dropping in-flight requests.

## When to Use

- Internal service-to-service routing that needs L7 features (host/path
  routing, rewrites, session affinity) without internet exposure
- Hub-and-spoke architectures where an internal gateway fronts shared
  services for peered networks

## Key Configuration Choices

- **Private-only frontend** -- a subnet frontend with STATIC allocation;
  the pinned address is what internal DNS records point at (also exported
  as the `private_ip_address` output)
- **Fixed `capacity: 2`** -- internal gateways often have predictable
  load; capacity and autoscale are mutually exclusive, so pick one
- **Connection draining (60s)** -- backends leaving the pool keep serving
  existing connections for a minute before termination
- **Member-side pool joining** -- the `services` pool is left empty and
  joined by NICs/scale sets through
  `status.outputs.backend_address_pool_ids.services`

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the gateway in | The resource group's `status.outputs.resource_group_name` |
| `internal-gateway` | The gateway's name, unique within the resource group | Your naming convention |
| `<gateway-subnet-arm-id>` | The DEDICATED gateway subnet | The subnet's `status.outputs.subnet_id` |
| `<static-private-ip>` | An unassigned address in the gateway subnet | Your IP address plan |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |

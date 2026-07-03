# Private Service Connect Backend

A PSC network endpoint group fronting a published producer service or a Google API — the way to reach a private, PSC-published backend from a load balancer without exposing it to the internet.

## When to Use

- A load balancer backend that is a Private Service Connect endpoint (a producer's published service attachment)
- Fronting a regional Google API bundle through PSC

## Remix Notes

- `pscTargetService` is required — it is the producer's service-attachment URL or a Google API bundle name.
- Reference a `GcpVpc` and `GcpSubnetwork` under `network.valueFrom` / `subnetwork.valueFrom` to wire the NEG to networks Planton manages.
- `pscData.producerPort` selects a specific producer port; omit it to connect to the first port in the producer's advertised range.

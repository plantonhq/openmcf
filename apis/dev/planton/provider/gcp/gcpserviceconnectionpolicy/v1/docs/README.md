# GCP Service Connection Policy — Deep Dive

## The problem this resource solves

Google's newer managed services stopped peering with customer networks. Private services access (the VPC-peering model behind Cloud SQL and AlloyDB private IP) requires the customer to reserve address ranges and accept a peering — address planning the producer then owns forever. The PSC-first generation (Memorystore for Valkey, Redis Cluster) inverts the model: the producer stays in its own network, and *endpoints* into your VPC are created on demand by Google's service connectivity automation.

But automation that can place forwarding rules in your subnets needs explicit, scoped authorization. That authorization is the service connection policy: a per-network, per-region, per-service-class grant naming the subnets endpoint IPs may come from and (optionally) how many connections may exist. Without the policy, creating a Memorystore instance on the network fails with a connectivity error. With it, endpoints appear automatically and the instance's IPs surface on the instance itself.

## Where it sits in the composition

Three first-class nodes, each with an independent lifecycle and owner:

- **GcpVpcNetwork** — the consumer network, typically owned by a platform team.
- **GcpSubnetwork** — the address space endpoints draw from. Regular-purpose subnets work; no special PSC purpose is needed for service connection policies (unlike PSC NAT subnets for published services).
- **GcpServiceConnectionPolicy** — the authorization joining a service class to the network and subnets.

The managed-service instance (e.g. `GcpMemorystoreInstance`) does not reference the policy — it references the network, and the automation finds the policy by (network, service class, region). Composition ordering is therefore a deployment-dependency concern: deploy the policy before the first instance, keep it alive until the last instance is gone.

## Cardinality: one policy per (network, service class, region)

The API rejects a second policy for the same triple. A shared VPC serving Memorystore in three regions needs three policies. Memorystore for Valkey (`gcp-memorystore`) and Redis Cluster (`gcp-memorystore-redis`) are distinct service classes — each needs its own policy even on the same network and region.

## The connection limit as a governance lever

`pscConfig.limit` caps how many PSC connections the automation may create under the policy. In a shared network this converts unbounded self-service into reviewed growth: teams create instances freely until the cap, then the network owner deliberately raises it (an in-place update). Zero/unset leaves GCP's default behavior.

## Producer location allowlisting

By default the automation connects to whatever producer instance the consumer project provisions. `CUSTOM_RESOURCE_HIERARCHY_LEVELS` plus `allowedGoogleProducersResourceHierarchyLevels` restricts producers to explicit `projects/…`, `folders/…`, or `organizations/…` locations — the lever regulated environments use to guarantee endpoints only route inside the company's own hierarchy. The spec enforces the two fields move together, because an allowlist without the mode flag is silently ignored by the API.

## Mutability profile

| Surface | Mutability |
|---|---|
| `location`, `network`, `serviceClass`, policy name | Immutable (ForceNew) |
| `pscConfig.subnetworks` | Mutable — add subnets as regions grow |
| `pscConfig.limit` | Mutable — raise the cap in place |
| producer allowlist | Mutable |
| `description`, `labels` | Mutable |

The mutable `pscConfig` is the operational point: capacity and governance changes never recreate the policy, so they never disturb existing endpoints.

## Address formats (a real trap)

The Service Connectivity API requires relative resource paths — `projects/{p}/global/networks/{n}` and `projects/{p}/regions/{r}/subnetworks/{s}` — and rejects full `https://` self-link URLs. Compute-family APIs accept both, so values sourced from compute outputs are frequently self-links. Both engines normalize by stripping the compute self-link prefix, making references and literals safe in either form.

## Teardown ordering

Deleting a policy does not delete the endpoints already placed under it, but it strands them operationally and blocks new connections. Destroy the managed-service instances first; destroy the policy last. In E2E and ephemeral environments this ordering is enforced by the dependency graph.

## IAM and API prerequisites

- `networkconnectivity.googleapis.com` and `compute.googleapis.com` enabled (both modules enable them with `disable_on_destroy=false`).
- `roles/networkconnectivity.consumerNetworkAdmin` (or the granular `networkconnectivity.serviceConnectionPolicies.*` permissions) on the project owning the network.
- In Shared VPC topologies the policy lives in the host project with the network.

## Failure modes an operator will actually meet

- **Instance creation fails with a connectivity/precondition error** — no policy exists for the service class in that region. Create the policy first.
- **"Policy already exists"** — the (network, service class, region) triple is taken; grow the existing policy instead of creating another.
- **Endpoint creation stalls** — the subnets listed are exhausted or in the wrong region; add a subnet (mutable) in the policy's region.
- **Connection limit reached** — the automation refuses new connections; raise `limit` in place.

## Deliberately not modeled (recorded reasons)

- **`deletion_policy`** — client-side Terraform lever (ABANDON drops the policy from state without deleting it); conflicts with Planton-managed destroy (catalog-wide decision).

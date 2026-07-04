# GCP Service Networking Connection — Deep Dive

## The problem this resource solves

Google's managed data services (Cloud SQL, AlloyDB, Memorystore, Filestore)
do not run inside your VPC — they run in producer projects Google operates.
"Private IP" for these services therefore cannot mean "an IP in your subnet";
it means a **VPC peering** between your network and the producer's network,
plus an agreed block of address space the producer is allowed to allocate
service subnets from. That pair — the peering and the reserved space — is
private services access, and this resource is the peering half.

Without it, a `GcpCloudSql` with `privateIpEnabled: true` fails to create
with a service-networking error. With it, every private-IP managed service
in the network works — one connection serves all of them.

## The three-node composition

Private services access decomposes into three first-class resources, each
independently owned and referenceable:

1. **GcpVpcNetwork** — the consumer network.
2. **GcpGlobalAddress** — an `INTERNAL` address with purpose `VPC_PEERING`
   and a `prefixLength`: the reserved block. Reserving it does nothing on
   its own; it only marks space the producer may use.
3. **GcpServiceNetworkingConnection** — hands the reserved range names to
   the producer and creates the peering.

The connection references ranges by **name** — the API's own addressing for
peering ranges — which is why `GcpGlobalAddress` exports its `name` as a
stack output alongside the address itself.

## Cardinality: one connection per (network, service) pair

GCP enforces exactly one connection for a given network + producer pair, and
the create call for a duplicate fails with "Cannot modify allocated ranges".
Two practical consequences:

- **Growth is in-place.** More capacity means appending another reserved
  range to `reservedPeeringRanges` on the existing resource. GCP keeps every
  service subnet it already carved, so the update is additive and safe.
- **Adoption is explicit.** If a connection for the pair already exists
  outside Planton (a console flow, gcloud), `updateOnCreationFail: true`
  converts the duplicate-create failure into an in-place update that adopts
  it. It is off by default so an accidental collision surfaces as an error
  instead of silently rewriting someone else's ranges.

## Sizing reserved ranges

The producer requests a block (typically a /24) per service subnet, per
region, per service — and instances multiply. A /16 is the common default
reservation; fragmented or undersized ranges surface later as
range-exhaustion errors on instance create. Because ranges can be appended
but never resized in place, generous first allocations beat expansion
churn.

## Teardown ordering

GCP refuses to delete the connection while the producer still holds subnets
in the reserved ranges. Destroy order is therefore: private-IP service
instances first (Cloud SQL, AlloyDB, Memorystore, ...), then this
connection, then the reserved ranges and network. Planton's dependency graph
(instances reference the network; the connection references the network and
ranges) yields this order naturally when resources are composed by
reference.

## Mutability profile

`reservedPeeringRanges` updates in place. `network` and `service` are
ForceNew: changing either destroys and recreates the connection — severing
private connectivity for every producer resource riding it — so treat both
as permanent once instances exist.

## How teams set this up today (deployment methods compared)

The same three-step dance — reserve a range, create the connection, point
services at private IP — exists in every toolchain. The differences are in
who remembers the ordering and the cardinality rule.

| Method | How it looks | Where it bites |
|--------|--------------|----------------|
| **gcloud** | `gcloud compute addresses create ... --purpose=VPC_PEERING` then `gcloud services vpc-peerings connect --ranges=...` | Two imperative commands with an implicit ordering; `vpc-peerings connect` silently *updates* an existing connection (the CLI's `connect` is an upsert), so a second team's ranges can be overwritten without an error. Nothing records which services depend on the connection at teardown time. |
| **Cloud Console** | The "Private services access" tab under VPC network → Private service connection | Fine for a first exploration; produces exactly the unmanaged pre-existing connection that later forces `updateOnCreationFail: true` when the network moves under IaC. |
| **Raw Terraform** | `google_compute_global_address` + `google_service_networking_connection` in one stack | Correct, but every team re-derives the same pitfalls: the ranges must be passed by *name*, `deletion_policy = "ABANDON"` gets cargo-culted to dodge teardown-order errors, and the one-connection-per-pair rule surfaces as a confusing "Cannot modify allocated ranges" apply failure. |
| **Planton composition** | `GcpVpcNetwork` + `GcpGlobalAddress` + `GcpServiceNetworkingConnection`, wired by references | The ordering and addressing rules are encoded once: references resolve the self-link and range names, the dependency graph yields create/destroy order, and adoption of a console-era connection is an explicit spec field instead of an accident. |

The upsert behavior of `gcloud services vpc-peerings connect` is worth
internalizing: it is the reason the API returns "Cannot modify allocated
ranges" to *declarative* tools — the underlying `services.connections.create`
call refuses to double-create, and the CLI papers over it by falling back to
update. `updateOnCreationFail` is this module's explicit, opt-in version of
that fallback.

## IAM and API prerequisites

- The deploying identity needs `servicenetworking.services.addPeering` on
  the network's project — bundled in **Service Networking Admin**
  (`roles/servicenetworking.networksAdmin`) — plus **Compute Network Admin**
  (`roles/compute.networkAdmin`) for the peering objects, and
  `serviceusage.services.enable` (Service Usage Admin or Editor) for the
  in-module API enablement.
- The module enables `servicenetworking.googleapis.com` and
  `compute.googleapis.com` itself, and deliberately leaves them enabled on
  destroy (`disable_on_destroy = false` semantics on both engines): APIs are
  project-wide, and tearing down one connection must never break unrelated
  workloads that share the project.
- For a **Shared VPC** host network, the connection must be created in the
  host project (address the network by full self-link); service-project
  identities cannot peer the host network.

## Failure modes an operator will actually meet

- **"Cannot modify allocated ranges" on create** — a connection for this
  (network, service) pair already exists (console, gcloud, or another
  stack). Either adopt it with `updateOnCreationFail: true` or import/delete
  the stray. Never create a second resource for the same pair.
- **Range exhaustion on instance create** (Cloud SQL/AlloyDB error citing
  the allocated range) — append a new `GcpGlobalAddress` to
  `reservedPeeringRanges`; the in-place update is additive and safe.
- **"Failed to delete connection" on destroy** — a producer still holds
  service subnets in the reserved ranges. Destroy the private-IP instances
  first; the connection cannot (and should not) force-delete them.
- **Peering exists but instances are unreachable** — check the VPC's
  peering route exchange: producer-side routes arrive over the
  `servicenetworking-googleapis-com` peering, and custom-route
  export/import (configured on the VPC peering, not on this resource) is
  needed only for transitive on-prem/VPN reachability.

## Deliberately not modeled (recorded reasons)

- **`deletion_policy`** — a client-side Terraform lever (`ABANDON` removes
  the connection from state without deleting it in GCP) that conflicts with
  Planton-managed destroy; catalog-wide decision. The teardown-ordering
  requirement it usually papers over is documented above instead.
- **`google_service_networking_peered_dns_domain`** — a sibling provider
  resource that forwards a private DNS suffix from the producer network over
  this peering. Real but second-order; Tier-2 candidate on concrete pull.
- **`google_service_networking_vpc_service_controls`** — the VPC Service
  Controls enablement toggle for this connection; belongs with a broader
  VPC-SC perimeter story (Tier-2).

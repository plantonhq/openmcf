# GcpServerlessVpcConnector — Research and Design Documentation

## 1. Introduction

### What Is a Serverless VPC Access Connector?

Serverless compute (Cloud Functions, Cloud Run, App Engine) runs outside your VPC. When a function needs to reach a private IP — a Cloud SQL private address, a Memorystore instance, an internal load balancer — something has to bridge the serverless environment into the network. The Serverless VPC Access connector is that bridge: a GCP-managed autoscaling fleet of small forwarding instances that live inside the VPC and relay serverless egress.

The connector is **regional shared infrastructure**: one connector serves every serverless workload in its region that references it. It is not per-service plumbing, which is why it is a first-class composable node rather than a field folded into each serverless kind.

### The Composition Boundary

- **GcpServerlessVpcConnector** owns the bridge (placement, machine type, scaling band).
- **GcpCloudFunction / GcpCloudRun / GcpCloudRunJob** attach by reference — their `vpcConnector`/`vpcAccess.connector` fields take the connector's fully qualified name (the `selfLink` output).
- **GcpVpcNetwork / GcpSubnetwork** are the placement targets.

Cloud Run (services and jobs) additionally supports **Direct VPC egress** — instances get IPs in a subnetwork with no connector infrastructure. That is the recommended Cloud Run path; the connector remains the only mechanism for Cloud Functions and App Engine, and stays relevant where org constraints mandate it.

## 2. Placement: The Two Arms

The API attaches a connector to exactly one of:

1. **network + ip_cidr_range** — the connector carves a dedicated `/28` out of the network's address space. The range must overlap no existing subnet, peered range, or route. Simple, and right for single-project VPCs.
2. **subnet** — the connector occupies an EXISTING `/28` subnetwork created exclusively for it. **Required on Shared VPC** (the range lives in the host project, named via `subnet.project_id`); also preferred where network admins want the range managed like any other subnet.

The provider expresses this as `ExactlyOneOf(network, subnet.name)` plus `RequiredWith(ip_cidr_range, network)`; the Planton spec enforces the same shape pre-deploy with three CEL rules (exactly-one placement, network-requires-cidr, cidr-requires-network).

## 3. Capacity Model

- **machine_type** — per-instance throughput class (`f1-micro` ~100 Mbps, `e2-micro` ~200 Mbps, `e2-standard-4` ~1 Gbps class). Mutable in place: a fleet-wide capacity lever with no replacement.
- **min_instances (2–9) / max_instances (3–10)** — the autoscaling band, min strictly below max. Two sharp edges the spec and modules document:
  - The fleet **never scales in on its own** — after a burst it stays at the high-water mark until the band is manually reduced.
  - The provider applies band **increases in place** but **replaces the connector on any decrease** (a `CustomizeDiff` that forces new only on shrinkage) — a brief egress outage for every attached workload.

### Deliberately Not Modeled

- **min_throughput / max_throughput** — the legacy scaling contract. The provider's own descriptions discourage both in favor of the instance fields, they `ConflictsWith` the instance fields, and they are ForceNew (every change replaces the connector). Modeling one honest scaling contract keeps the spec free of a trap arm.
- **deletion_policy** — absent from the released 6.x provider line (verified via schema probe against 6.50.0; present only on the unreleased main line and in the bridged Pulumi SDK). Modeling it would create a one-engine field.
- **connected_projects** (output-only) — an informational read-side list, not a composition key.

## 4. Terraform Provider Floor

Designed from `google_vpc_access_connector` on the released Terraform Google provider 6.x line (`~> 6.0`, schema-probed at 6.50.0). The resource has **no labels surface**, so this is one of the rare GCP kinds where the modules attach no platform attribution labels — both engines skip them identically. Both engines enable `vpcaccess.googleapis.com` before creating the connector.

## 5. Registry

- **Enum:** 721 (second entry in the 720–729 GCP serverless overflow block)
- **ID prefix:** `gcpvpcconn`
- **Prerequisites:** `GcpVpcNetwork`, `GcpSubnetwork` (one per placement arm)

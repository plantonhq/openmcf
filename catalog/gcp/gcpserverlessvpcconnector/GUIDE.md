# GcpServerlessVpcConnector Guide

The judgment this guide protects: the connector is shared regional
infrastructure that looks per-workload. Sizing it, placing it, and
tearing it down all have consequences for every function and Cloud Run
service in the region that egresses through it.

## Placement: network vs subnet

Two modes, exactly one:

- **Network placement** (`network` + `ipCidrRange`): the connector carves
  a dedicated /28 out of the VPC. The simple path for a single-project
  VPC — pick an unused /28 that no subnet, peered range, or route claims.
- **Subnet placement** (`subnet`): the connector occupies an existing
  dedicated /28 subnetwork. **Required on Shared VPC** (the range lives
  in the host project — set `subnet.projectId`), and the right choice
  wherever network admins manage all address space as subnets.

The placement is immutable; switching modes is a replacement.

## Fleet economics

Throughput = instances × per-instance rate for the machine type
(`f1-micro` ~100 Mbps, `e2-micro` ~200 Mbps, `e2-standard-4` ~1 Gbps
class). Two behaviors dominate real bills and real incidents:

- **The fleet never scales in.** After a burst it stays at the high-water
  mark until an operator lowers the band. A connector that spiked once
  keeps billing for the spike.
- **Decreasing `minInstances` or `maxInstances` REPLACES the connector**
  — a brief egress outage for every workload using it. Increases apply
  in place. Start low and grow; shrinking is the disruptive direction.

`machineType` is the exception: it changes in place, making it the
cheapest capacity lever when the instance band is already right.

Throughput-based sizing (`min_throughput`/`max_throughput`) is
deliberately not modeled: the provider discourages it in favor of the
instance band, the two conflict, and it forces replacement on change.

## Create and destroy are slow

GCP builds and tears down a managed instance fleet: budget 3–5 minutes
each way. In pipelines, do not treat a slow connector create as a hang.

## One connector serves the region

Cloud Functions, Cloud Run, and App Engine in the same region share one
connector — attach by the `self_link` output. Resist per-service
connectors: each one occupies a /28 and runs its own always-on fleet.

## Teardown discipline

Destroying the connector silently breaks private egress for every
serverless workload configured to use it — the workloads keep running
and start failing on VPC-internal calls. `deletionPolicy: PREVENT` is
the right posture once anything in production egresses through it;
`ABANDON` leaves the fleet running unmanaged (and billing).

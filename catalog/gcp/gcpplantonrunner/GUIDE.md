# GcpPlantonRunner Guide

The judgment this guide carries: this appliance exists for exactly one
reason -- a deploy target that only the inside of your network can reach.
The mistakes it protects against are deploying it when the hosted fleet
would do, treating the join token as the runner's identity, and tearing
it down before the workloads it deploys.

## When to use it (and when not)

Reach for this kind only when a target is invisible from outside the
network: a private GKE control plane, a private-IP Cloud SQL instance.
If every endpoint you deploy to is public, Planton's hosted runner fleet
already covers you -- a self-hosted runner then buys you nothing except
one more piece of standing infrastructure to size, patch (by version
pin), and pay for. The moment a private endpoint enters the
architecture, the calculus flips: no hosted runner will ever reach it,
and this appliance is the only deploy path.

One runner covers a network perimeter, not a workload. Deploy one per
VPC-region you operate in, not one per cluster -- whatever the
subnetwork's routes can reach, the runner can deploy to.

## The token is not the identity

Enrollment is token-first: the manifest carries a join token, and the
runner trades it for its own individually revocable identity on first
boot. Two operational consequences that are easy to get backwards:

- **Rotation is calm.** The service reads the token secret's `latest`
  version at instance start only. Rotating the token requires no service
  update and interrupts nothing -- the running instance keeps serving on
  its minted identity; the next instance replacement joins with the new
  value.
- **Revocation is precise.** Revoking a token never touches the runners
  it already admitted; revoke the runner's own identity to cut off a
  runner. Instance replacement re-joins with the same token -- the
  token's lineage re-admits the runner it originally admitted, so a
  Cloud Run instance being replaced is a non-event, not a re-enrollment
  ceremony.

## Why the CPU is always allocated

Cloud Run's default throttles CPU between requests -- correct for
request-serving apps, fatal for this one. The runner is a pull-based
worker: it polls its work queue and executes long-running IaC operations
with **no inbound request in flight**. Under the default, a
Terraform apply would freeze mid-operation the moment Cloud Run decided
no request was being served. The module pins `cpu_idle = false`
(instance-based, always-allocated CPU) for exactly this reason -- never
"optimize" it back; the symptom is operations that hang or time out
midway with nothing useful in the logs.

## VPC egress cannot lock you out

`vpcAccess` uses Direct VPC egress with private-ranges-only routing:
only RFC-1918-bound traffic rides the VPC; the runner's control-plane
dial-out keeps its normal internet path. This is a deliberate safety
property -- a misconfigured VPC (missing NAT, botched firewall, wrong
subnetwork) can break the runner's reach into private endpoints, but it
can never sever the runner from the control plane that manages it. The
runner stays visible in `planton runner list` and diagnosable while you
fix the network, instead of going dark exactly when you need it.

## The identity posture

When `serviceAccount` is unset, the module creates a dedicated
permissionless account -- deliberately never the project's Compute
Engine default, which typically carries broad project access the runner
should not inherit silently. The only grant the module ever makes is
`secretAccessor` on the runner's own token secret, scoped to exactly
that one secret and exactly that one principal. Everything beyond that
is yours to grant, deliberately, on the account the outputs expose as
`service_account_email` -- the seam keyless cloud access rides. Never
grant the runner roles "while you're in there"; the permissionless
default is the point.

## Deploys fail loudly here

Cloud Run gates service creation on the first revision becoming ready.
If the token secret cannot be resolved at instance start, or the
container never comes up, the apply itself fails -- you find out at
deploy time, with the error in the deploy output. This is unlike
substrates that report infrastructure success and leave a crash-looping
container for you to discover later. Treat a failed apply of this kind
as a real signal: check the token reference and the accessor grant
before retrying.

## Teardown ordering

The runner is the deploy path for the private workloads behind it, so it
is destroyed **last**: in-cluster workloads are destroyed through the
runner, the cluster over the GCP path, and the runner at the end.
Destroying the runner first strands everything it deploys -- nothing can
reach the private endpoints to tear them down. The appliance itself
tears down cleanly (no deletion protection): the token makes it
re-mintable standing infrastructure, so destroy-and-recreate is a valid
recovery move.

## On the diagram

The runner renders as its own node, with edges to the `GcpProject`,
`GcpVpcNetwork`, `GcpSubnetwork`, and `GcpServiceAccount` it references
-- the architecture shows *how* private deploys happen, not just that
they do. A referenced runtime service account renders as a first-class
node whose permissions are visible in the composition; the created
permissionless default renders as nothing.

## Pairs well with

- GcpGkeCluster -- the canonical private target: pair a
  private-control-plane cluster with a runner in the same VPC-region.
- GcpVpcNetwork / GcpSubnetwork -- the placement that defines the
  runner's reach; the subnetwork must be in the runner's region.
- GcpServiceAccount -- compose the runtime identity first-class when
  keyless cloud access needs real permissions, instead of granting roles
  to the module-created account out of band.

# AzurePlantonRunner Guide

The judgment this guide carries: this appliance exists for exactly one
reason -- a deploy target that only the inside of your network can reach.
The mistakes it protects against are deploying it when the hosted fleet
would do, pointing it at an environment that cannot reach the private
endpoints it was deployed for, and tearing it down before the workloads
it deploys.

## When to use it (and when not)

Reach for this kind only when a target is invisible from outside the
network: a private AKS API server, a private-endpoint database. If every
endpoint you deploy to is public, Planton's hosted runner fleet already
covers you -- a self-hosted runner then buys you nothing except one more
piece of standing infrastructure to size, patch (by version pin), and
pay for. The moment a private endpoint enters the architecture, the
calculus flips: no hosted runner will ever reach it, and this appliance
is the only deploy path.

One runner covers a network perimeter, not a workload. Deploy one per
VNet you operate in, not one per cluster -- whatever the environment's
network can reach, the runner can deploy to.

## The environment decides everything about reach

The runner spec carries no network knobs at all -- deliberately. The
referenced Container App Environment IS the network decision: a
VNet-integrated environment gives the runner reach into that network's
private endpoints; a plain one leaves it with only public reach. If the
runner cannot reach the private AKS API it was deployed for, the fix is
never on the runner -- recompose or replace the environment. The module
references the environment and the resource group, never creates or
mutates them, so the runner can be destroyed and recreated freely
without touching the network layer.

## The token is not the identity

Enrollment is token-first: the manifest carries a join token, and the
runner trades it for its own individually revocable identity on first
boot. Two operational consequences that are easy to get backwards:

- **Rotation is calm.** The token is only read at join. Rotating it
  interrupts nothing -- the running replica keeps serving on its minted
  identity, and the next replica replacement joins with the new value.
  The token rides the app's configuration as a secret (never a plain env
  var), so rotation is an ordinary apply, not an emergency.
- **Revocation is precise.** Revoking a token never touches the runners
  it already admitted; revoke the runner's own identity to cut off a
  runner. Replica replacement re-joins with the same token -- the
  token's lineage re-admits the runner it originally admitted, so a
  Container Apps replica being replaced is a non-event, not a
  re-enrollment ceremony.

## The pairing law

On the Consumption plan, memory is always cpu x 2 -- 0.5 vCPU pairs with
1Gi, 1 with 2Gi, and so on. Azure enforces this only at deploy time, so
the spec validates it up front: an invalid pairing fails instantly with
an explainable message instead of minutes into an apply. The practical
consequence: you cannot size memory independently. If IaC operations
fail mid-apply under memory pressure, the move is up a whole pairing
step, which buys CPU you may not need -- that is the plan's trade, not
the component's.

## The no-ingress posture

The app has no ingress block at all -- not "ingress disabled", none
configured. The runner initiates every connection it uses (control
plane, work queue, image pulls), so there is nothing to expose and no
listener to protect. Never add ingress to "check on" the runner; its
observability path is the control plane (`planton runner list`) and the
app's console logs, not an HTTP endpoint.

## The image pull is part of the create (proven live)

Container Apps validates the container image's manifest while
provisioning the app's first revision, DURING the create: an image
reference the registry cannot serve fails the whole deployment in
seconds with `ContainerAppOperationError ... MANIFEST_UNKNOWN` on both
engines -- there is no half-created app to inspect. Once the pull
succeeds, the boundary flips: the app provisions independently of
replica health (a container that starts and crash-loops -- say, a
runner whose join is refused -- still leaves the app provisioned and
manageable; measured creates ~36-61s). Two operational consequences:

- A create failing with `MANIFEST_UNKNOWN` is a REGISTRY problem
  (unpublished tag, wrong repository, single-platform image missing
  linux/amd64 -- Container Apps runs amd64), never a spec or module
  problem. Check the reference with
  `docker buildx imagetools inspect <image>` before touching anything
  else.
- The default `image_repository` resolves the official image's
  `latest` tag, which every runner release moves. If you override it
  to a mirror, keeping the mirror's tags populated -- for BOTH
  platforms -- becomes your responsibility, and this create-time
  validation is where a stale mirror announces itself.

## The singleton is a law, not a default

Exactly one replica, single-revision mode. A runner's identity is minted
for one live replica -- a second replica joining under the same name
would revoke the first's key (token lineage: re-admission re-mints and
revokes). Revision rollover briefly overlaps two replicas; the draining
one's revoked key dies with it, and a rollback self-heals by re-joining.
Never scale this app, and never switch it to multiple-revision mode.

## Teardown ordering

The runner is the deploy path for the private workloads behind it, so it
is destroyed **last**: in-cluster workloads are destroyed through the
runner, the cluster over the Azure path, and the runner at the end.
Destroying the runner first strands everything it deploys -- nothing can
reach the private endpoints to tear them down. The environment and
resource group outlive the runner (they are referenced, not owned), so
destroying the runner never disturbs neighbors sharing them.

## On the diagram

The runner renders as its own node, with edges to the
`AzureResourceGroup` and `AzureContainerAppEnvironment` it references --
the architecture shows *how* private deploys happen, not just that they
do. Because the environment carries the network decision, composing it
first-class (rather than passing a literal ID) is what makes the private
reach visible on the diagram.

## Pairs well with

- AzureAksCluster -- the canonical private target: pair a
  private-API-server cluster with a runner in an environment on the same
  VNet.
- AzureContainerAppEnvironment -- the placement that defines the
  runner's reach; VNet-integrate it for private endpoints.
- AzureResourceGroup -- compose the group first-class so the appliance's
  lifecycle boundary is explicit in the architecture.

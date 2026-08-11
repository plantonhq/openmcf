# GcpVertexAiDeployedIndex Guide

The judgment this guide protects: a deployment is a long-lived
placement, not a stateless join. Deploys take tens of minutes (the
provider allows 45), and the ONLY fields that update in place afterward
are the replica bounds inside the sizing arm — everything else
undeploys and redeploys, which means minutes of serving gap unless you
blue/green with a second deployment.

## Size with the arm that matches your certainty

`automaticResources` is the zero-knowledge start: Vertex AI picks
machine types and scales replicas between your bounds. Move to
`dedicatedResources` when you know the load shape — a pinned machine
type gives predictable latency and cost, but the type must be
compatible with the INDEX's `shardSize` (small shards accept
e2-standard-2 and up; large shards need e2-highmem-16 or
n2d-standard-32 class machines), and the pairing is checked at deploy
time, half an hour in. Min replicas below 2 carries no SLA on either
arm.

## The deployment_group ↔ reserved ranges pairing is HELD forever

Pairing `deploymentGroup` with `reservedIpRanges` gives each group a
predictable IP space on a peered endpoint — and the API remembers: a
non-default group, once used with a set of reserved ranges, can only
ever be used with exactly that set again. Choose group names and their
ranges as a permanent allocation, not a label you can rewrite.

## deployed_index_id outlives the deployment — and its parent

The user-chosen `deployedIndexId` is held by GCP while a deployment is
failed or still undeploying — with a 400 ("retry later or use a
different ID") — and the hold SURVIVES DELETING THE PARENT ENDPOINT.
Never hardcode an ID you may need to recreate quickly; derive it with a
run/version suffix so a redeploy never waits on an invisible hold. The
ID also forbids hyphens (letters, numbers, underscores only).

## Auth applies to the private path only

`authConfig` (JWT via allowed service-account issuers) guards the
private query endpoint on peered/PSC endpoints. A public endpoint's
auth is IAM on the public domain — the block does nothing there.

## deletionPolicy speaks undeploy, not delete

Empty/`DELETE` undeploys the index from the endpoint on destroy —
queries against this deployment stop; the index itself is untouched.
`PREVENT` makes destroy fail: the right posture when this deployment IS
the production query path. `ABANDON` keeps it serving (and billing for
replicas) outside management — the escape hatch for handing traffic to
another control plane without a serving gap.

## What is deliberately absent

No labels and no project field exist on this resource class in the GCP
API (the deployment lives inside the endpoint), so the platform's
attribution labels are impossible here and none are faked. Cost
attribution happens on the endpoint and the index.

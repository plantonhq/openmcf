# Where a Service Answers — URLs and Rollout Verification

Every successful deployment record carries two answers a person otherwise digs for in a cloud console: the addresses the deployed service serves at, and whether anything actually answered there. Read this when someone asks where their service is running, whether a deploy is really up, or what a rollout verdict on a deployment record means.

## The URLs on a deployment record

A deployment record lists every address its deployed resources declared — the serverless platform's generated URL, the load balancer's DNS name, the ingress host, a function's endpoint. Each entry names its source (which resource and which of its outputs produced it), so you can always say WHERE an address came from.

Two things to hold when relaying them:

- The list is **discovery, not endorsement**. An address appears because the deployment target declared it, including addresses that did not answer the probe — the address that did not answer is exactly what a person debugging needs, so never hide it.
- A bare hostname (a load balancer's DNS name) appears under both https and http, because the scheme is not knowable from the target's outputs. Whichever one answers is the one that serves.

Some targets declare no address at all — a Kubernetes workload with no ingress has only in-cluster DNS, and container-orchestrator services expose their address through the load balancer deployed beside them. A record with no URLs is honest, not broken; say what would add one (an ingress, a load balancer) rather than treating it as a failure.

## The rollout verdict, and how to read it

Verification is staged the way a person checks by hand, and each step reports itself: **resources created on the provider** (the apply landed), and **endpoint answering** (the probe's report, with per-address outcomes). Deeper steps — the workload actually online, in each platform's own terms — join the list as the platform grows them. The overall verdict is one of three words, and relaying them accurately matters more than anything else in this file:

- **verified** — at least one declared address answered. The probe runs from the organization's own runner, inside their infrastructure, so private addresses count; the verdict's summary names which runner probed.
- **failed** — addresses existed and none answered from the runner's vantage. Say this carefully: **the deployment itself succeeded** — resources applied, the record exists, later environments still deployed. A failed verdict is a report that nothing answered at the declared addresses, with the per-address transport errors in the check's detail. Common honest explanations: the workload takes longer to come online than the probe window, DNS for the hostname is not set up yet, or the runner has no network route to that address. Investigate; never rephrase it as "the deploy failed".
- **unverifiable** — nothing could be probed, with the reason recorded: no declared address, or no runner was available to probe. Never treat it as a warning to fix urgently; treat it as the platform declining to guess.

Any HTTP status code counts as an answer — a 404 from a load balancer's default action still proves the address serves. The status code is in the check's detail; if someone expected a healthy page and the probe recorded a 500, that distinction is yours to surface.

## Answering "where is my service?"

The one-command answer is `planton service urls <service>` — each environment's current deployment with its addresses and verdict. For the full story on one deployment (per-address outcomes, the probe's vantage, the staged checks), read the deployment record itself. Promotions and rollbacks produce records with the same URLs and verdicts as any push-triggered deploy — no special case.

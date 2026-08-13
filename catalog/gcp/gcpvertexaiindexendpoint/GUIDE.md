# GcpVertexAiIndexEndpoint Guide

The judgment this guide protects: connectivity is THE create-time
decision. Public, VPC-peered, and Private Service Connect are mutually
exclusive and every one of them is immutable — changing your mind
replaces the endpoint, which undeploys every index serving on it.

## Choose the connectivity arm by who queries, not by habit

**Public** (`publicEndpointEnabled: true`) is right when callers live
outside your VPCs or you want zero network plumbing: queries go to a
managed domain name (the `public_endpoint_domain_name` output), and
auth is IAM. **VPC-peered** (`network`) keeps queries inside the VPC
but requires Private Services Access on the network FIRST — the
`GcpGlobalAddress` (VPC_PEERING) + `GcpServiceNetworkingConnection`
composition, deployed before this endpoint. **PSC**
(`privateServiceConnectConfig`) is the modern private model: no
peering; consumers connect through a service attachment that surfaces
on the DEPLOYED INDEX's outputs, not here — the attachment exists only
once an index is deployed.

## PSC has two consumer stories; only the allowlist works today

`projectAllowlist` lets consumer projects create their own forwarding
rules against the service attachment — explicit, auditable, and the
proven path. `pscAutomationConfigs` asks Vertex AI to create the
consumer-side endpoints itself (per project/network pair), but on
index endpoints the live API accepts the create and then silently
drops the configs: the stored endpoint omits them and no consumer-side
endpoint is ever provisioned (verified against the live API; the
provider documents the field as used by online inference endpoints
only). Because the PSC block is immutable, that silent drop would
otherwise surface as a perpetual replacement diff on every re-plan —
so the spec refuses the field on this kind. Wire consumers through the
allowlist; the field unlocks if Google extends automation to vector
search.

## One endpoint, many deployments — pool deliberately

Deployments share the endpoint's connectivity and (with automatic
resources) its serving infrastructure. Grouping related indexes on one
endpoint pools cost; isolating a latency-critical index on its own
endpoint isolates noisy neighbors. This is the only axis where you can
re-decide later — deployments move between endpoints by
undeploy/redeploy, the endpoint itself never changes shape.

## CMEK is a create-time decision with an IAM prerequisite

`kmsKeyName` (a `GcpKmsKey` reference) must name a key in the
endpoint's region, and the Vertex AI service agent needs
`roles/cloudkms.cryptoKeyEncrypterDecrypter` on it BEFORE the create.
Immutable — encrypting an existing endpoint means replacing it and
redeploying every index.

## deletionPolicy guards the whole serving surface

Empty/`DELETE` deletes the endpoint and every deployment on it stops
serving. `PREVENT` makes destroy fail — the right posture for an
endpoint with production deployments, because the blast radius is
every index on it, not one resource. `ABANDON` hands the running
endpoint to another management plane (replicas keep billing).

## What is deliberately absent

The endpoint serves nothing by itself: compute sizing, auth on the
private query path, and reserved-range pinning all live on
`GcpVertexAiDeployedIndex` — the resource that joins an index to this
endpoint.

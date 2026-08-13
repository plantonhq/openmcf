# GcpVertexAiEndpoint Guide

The judgment this guide protects: the endpoint is the STABLE half of
model serving — a durable URL and network boundary that outlives every
model deployed onto it. Keep the endpoint boring and long-lived; let
models and traffic move.

## The endpoint outlives the model

Deploying models onto the endpoint is an operational act (Vertex AI API,
console, CI) — deliberately outside this resource. What this node owns
is everything a deployment inherits: the region, the network posture,
the CMEK key, the logging tap, and the ID applications call. Create it
once per serving surface and let model versions come and go behind it.

## trafficSplit moves traffic, not models

The split maps a DeployedModel's ID (assigned at deploy time) to its
percentage; values must sum to exactly 100 and every ID must be
currently deployed — so a fresh endpoint keeps it EMPTY, and the first
entry arrives only after the first model deploy. Canary flows are a
spec edit: 90/10, watch, 0/100. Keep the IaC split as the source of
truth after each deploy, or leave it unmanaged and empty — mixing
console-edited splits with a managed field is how rollouts fight
their own automation.

## Three network postures, chosen at create

Public (default), VPC-peered (`network`, needs Private Services Access),
or PSC (`privateServiceConnectConfig`) — mutually exclusive and
immutable. PSC is the modern private path: `projectAllowlist` says who
may connect, and `pscAutomationConfigs` goes further — Vertex AI creates
the consumer-side endpoints (forwarding rule + IP) in each listed
project/network itself. Both PSC lists are mutable, so onboarding a new
consumer project is an edit, not a rebuild.

## The numeric-ID reservation trap

Endpoint IDs are numeric and GCP RESERVES a deleted endpoint's ID —
destroy-then-recreate with the same identity fails with 409 until the
reservation lapses. The module derives a stable ID from the resource
identity, so the trap only bites on delete/recreate cycles: recreate
under a new name (or explicit `endpointName`) instead of waiting.

## Destroy refuses while models serve

GCP rejects deleting an endpoint that still has deployed models — undeploy
first, always. `deletionPolicy: PREVENT` suits the endpoint whose URL is
baked into applications; `ABANDON` keeps it serving (deployed models keep
billing) while dropping it from management.

## Logging is drift insurance

`requestResponseLoggingConfig` samples real prediction traffic into
BigQuery — the raw material for drift detection and incident forensics.
Sample low (0.01–0.1) on high-QPS endpoints; the cost lives in BigQuery,
not Vertex. Turn it on before the incident, not after.

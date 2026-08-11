# GcpRegionNetworkEndpointGroup Guide

The judgment this guide protects: a regional NEG is the immutable
adapter between a load balancer and everything that is not a VM group.
Every field is ForceNew and an in-use NEG cannot be deleted — so the
whole operational story is naming and replacement discipline.

## The endpoint-type decision

- **`SERVERLESS`** (the default): put Cloud Run, Cloud Functions, or App
  Engine behind an external ALB — custom domains, CDN, Cloud Armor, IAP
  in front of serverless. Set exactly one of `cloudRun` / `cloudFunction`
  / `appEngine`. The NEG must live in the SAME region as the workload.
- **`PRIVATE_SERVICE_CONNECT`**: front a published producer service or a
  Google API bundle — set `pscTargetService` (and `network`/`subnetwork`).
- **`INTERNET_IP_PORT` / `INTERNET_FQDN_PORT`**: an external origin
  behind a Google load balancer — the hybrid/third-party-backend pattern.
- **`GCE_VM_IP_PORTMAP`**: PSC port mapping to VM IP:port targets.

## Everything is ForceNew — name for replacement

Any change destroys and recreates the NEG, and GCP refuses to delete a
NEG a backend service still references (`resourceInUseByAnotherResource`).
The working pattern is versioned names (`api-neg-v2`): create the
replacement NEG, repoint the backend service's `backends[].group`, then
destroy the original. Editing a live NEG's spec in place WILL fail its
apply at the delete step — that is the API talking, not the module.

## urlMask is the fan-out lever

One NEG can front many Cloud Run services (or functions) by parsing the
target out of the request URL — `urlMask: "<service>.example.com"` routes
per-host, `"/api/<service>"` per-path. It replaces N per-service NEGs
with one, at the cost of coupling routing to URL structure. GCP resolves
targets at serving time, so masked services need not exist when the NEG
is created — a typo surfaces as 404s, not as a deploy failure.

## Serverless targets resolve lazily

`cloudRun.service` (and friends) are resolved at request time. Creating
the NEG before the service exists is legal and useful in charts — but it
also means a deleted service leaves the NEG "healthy" while requests
fail. The NEG is wiring, not a health check.

## Teardown discipline

`deletionPolicy: PREVENT` suits a NEG referenced by a production backend
service — the in-use refusal protects the attached case, PREVENT also
covers the window where the backend service was destroyed first.
`ABANDON` leaves the NEG in GCP unmanaged.

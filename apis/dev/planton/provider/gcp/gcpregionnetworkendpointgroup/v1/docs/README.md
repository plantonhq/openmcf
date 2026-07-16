# GCP Region Network Endpoint Group — Deep Dive

## Where this sits in the load balancing family

A Google Cloud external Application Load Balancer is assembled from first-class,
independently-referenceable resources: a forwarding rule and IP address at the
front, a target proxy, a URL map that routes, and one or more backend services
that decide how traffic reaches compute. A **backend service** points at
**backends**, and each backend's `group` is either an instance group (Compute
Engine VMs) or a **network endpoint group**.

The regional NEG is what makes non-VM targets first-class in that graph. Without
it, a Cloud Run service could not sit behind a shared external load balancer
with a custom domain, Cloud CDN, Cloud Armor, and IAP. The NEG is the node that
says "this backend is a Cloud Run service in us-central1" — and because it is
its own resource with its own self-link, it composes into a backend service the
same way an instance group does.

## The endpoint type is the whole design

One provider resource covers five endpoint types, and the type decides which
fields are meaningful. Modeling them as one kind (rather than five) matches the
GCP resource exactly and keeps the mental model small: choose a type, fill the
one block it needs.

- **SERVERLESS** — the common case. Exactly one of `cloudRun`, `cloudFunction`,
  or `appEngine`. Each can name its target directly (`service` / `function`) or
  derive it per-request from a `urlMask` template, which is how one NEG can fan
  a wildcard domain out to many services. The App Engine block may be empty to
  route to the project's default application.
- **PRIVATE_SERVICE_CONNECT** — front a published producer service or a Google
  API through a PSC endpoint (`pscTargetService`, with an optional `network`,
  `subnetwork`, and `pscData.producerPort`).
- **INTERNET_IP_PORT / INTERNET_FQDN_PORT** — front an external origin (an
  on-prem or third-party backend) reached over the internet.
- **GCE_VM_IP_PORTMAP** — PSC port mapping to VM IP:port targets.

The spec's CEL rules enforce the type/block coherence GCP would otherwise reject
at apply time: exactly one serverless block for SERVERLESS (none otherwise),
`pscTargetService` required for PSC, and `network`/`subnetwork`/`pscData` kept
off serverless NEGs. Catching these before deploy turns a slow round-trip
failure into an instant, explained validation error.

## Immutability and the create-before-destroy hazard

Every field of a regional NEG is immutable (ForceNew): changing the region, the
endpoint type, or the serverless target destroys and recreates the NEG. Two
consequences matter operationally:

1. A NEG that a backend service references cannot be deleted while in use — GCP
   returns `resourceInUseByAnotherResource`. Recreating one therefore requires
   creating the replacement first and repointing the backend service before the
   old NEG is destroyed. When you change the referencing backend service through
   Planton, the dependency ordering handles this.
2. The serverless target need not exist when the NEG is created — endpoints are
   resolved at serving time. This is why a NEG can be provisioned in the same
   change as (or before) the Cloud Run service it fronts, and why a backend
   service with only serverless NEG backends needs no health check.

## Composition

The NEG's `self_link` is the composition handle. A `GcpBackendService` backend
sets `group` to it; because a backend group can be either an instance group or a
NEG, `backends[].group` defaults to this kind but accepts any group producer by
explicit reference. A single backend service can combine serverless NEGs from
several regions to serve one global anycast IP.

## Deliberate scope boundaries

- The **zonal** NEG (`google_compute_network_endpoint_group`) and the **global**
  internet NEG (`google_compute_global_network_endpoint_group`) are separate GCP
  resources with different schemas and endpoint-type enums; they are distinct
  kinds, not folded here.
- `serverlessDeployment` (API Gateway and similar) is a beta-only block on the
  released provider line and is not modeled; the five GA endpoint types cover
  the overwhelming majority of real use.

---
kinds:
  - GcpCloudRun
  - GcpCloudRunDomainMapping
  - GcpDnsRecord
  - GcpDnsZone
---

# Cloud Run Service-to-Service: the Stable Hostname Two Services Need

Two Cloud Run services that call each other — a frontend that talks to a
backend API is the canonical shape — cannot be wired the way the rest of
the catalog wires producers to consumers. The obvious move is to feed the
backend's `status.outputs.url` into the frontend as an environment
variable. It does not work, and the reason is a hard limit of the Cloud
Run schema, not a missing feature of the platform.

## The trap: env vars take no reference

`GcpCloudRun.spec.containers[].env[].value` is a plain `string`. Its only
non-literal form is `valueFromSecret` (a Secret Manager reference) — there
is NO `valueFrom` arm on an env var. So this does not parse:

```yaml
# WRONG — env vars cannot reference another resource's output
env:
  - name: BACKEND_URL
    valueFrom:                       # not a field on env[]
      kind: GcpCloudRun
      name: my-backend
      fieldPath: status.outputs.url
```

The failure is quiet if you reach for it by habit: the field is dropped or
rejected, and the frontend ships with no way to find the backend. The
backend's own `run.app` URL is real and reachable (a public service with
`allowUnauthenticated: true`), but its hostname carries a deploy-time hash
and is ugly to hardcode — and hardcoding it couples the frontend to a URL
that can change.

## The fix: give each service a stable custom domain

Map a custom domain onto the backend, then point the frontend at that
fixed hostname with a plain literal env var. The hostname never changes
across redeploys, works identically for server-side and browser-side
calls, and is the shape you want in production anyway.

```yaml
# Backend gets its own domain — same shape as the frontend's.
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpCloudRunDomainMapping
metadata:
  name: staging-backend-domain
spec:
  region: asia-south1
  domain: api.staging.example.com
  route:
    valueFrom:
      kind: GcpCloudRun
      name: staging-backend
      fieldPath: status.outputs.service_name
  certificateMode: AUTOMATIC
---
# Frontend targets the backend's fixed hostname — a literal, not a reference.
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpCloudRun
metadata:
  name: staging-frontend
spec:
  region: asia-south1
  containers:
    - name: web
      image: us-docker.pkg.dev/cloudrun/container/hello:latest
      env:
        - name: BACKEND_URL
          value: https://api.staging.example.com
```

For a browser-side SPA the same hostname is what the browser must reach,
so the custom domain is not optional there — the backend has to be
publicly resolvable on a real name regardless. Keep the backend private
(no domain, VPC-internal ingress) ONLY when every call is server-to-server
from inside the frontend container.

## The DNS record usually crosses environments

Each domain needs a `GcpDnsRecord` (a CNAME to `ghs.googlehosted.com.`) in
the managed zone. In practice the zone is created ONCE in a shared
environment, while the app chart deploys into a per-environment slot
(`staging`, `dev`, `prod`). The record therefore references the zone
ACROSS environments — the `valueFrom` block carries an explicit `env:`:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpDnsRecord
metadata:
  name: staging-backend-cname
spec:
  managedZone:
    valueFrom:
      kind: GcpDnsZone
      name: example-com-zone       # the zone's resource name in the shared env
      fieldPath: status.outputs.zone_name
      env: shared                  # the zone is not in this app's environment
  type: CNAME
  name:
    value: api.staging.example.com.
  values:
    - value: ghs.googlehosted.com.
  ttlSeconds: 300
```

Never expose the zone name as a bare param the user hand-copies — it is a
value the platform already knows, so wire it. Two consequences of the
cross-environment reference:

1. **No deploy-order edge is created across the boundary.** A `valueFrom`
   inside one chart orders the deploy; a reference into another
   environment does not. The shared zone must ALREADY be deployed when the
   app chart deploys — compose and deploy the shared environment first.
2. **The build validates the reference's shape, not the zone's
   existence.** A green build does not prove the zone is there; the deploy
   resolves it. Name that assumption when handing off.

## The one-time prerequisite no manifest performs

`GcpCloudRunDomainMapping` refuses to create until the domain is VERIFIED
for the deploying GCP identity (Search Console / `gcloud domains verify`).
Verifying the parent domain once (`example.com`) covers every subdomain —
so `staging.example.com` and `api.staging.example.com` both map with no
extra step. No IaC resource does this; it is out-of-band and per parent
domain.

## On the diagram

Both services render as Cloud Run nodes; each domain mapping and DNS
record renders as its own node wired to its service. The
frontend→backend dependency, however, is INVISIBLE — it lives inside a
literal env-var string, drawing no edge. A reviewer confirming that a
frontend can reach its backend must read the env var, not trust the
picture: the graph shows two disconnected services even when the wiring is
correct.

## See also

- `GcpCloudRunDomainMapping`'s `reference.md` — the immutability and
  `resource_records` output that DNS consumes.
- `dependencies` guidance on cross-environment `valueFrom` (`env:`) and the
  ordering duty it does not create.

# KubernetesExternalDns: Research and Design

## Introduction

ExternalDNS is the Kubernetes-SIGs controller that makes DNS follow
workloads: it watches **sources** (Services, Ingresses, Gateway API routes,
CRDs, ...) for desired hostnames, translates them into records, and writes
them to a **provider** (Route 53, Cloud DNS, Azure DNS, Cloudflare, ...)
through a **registry** that tracks which records it owns. This component
installs that controller from the official Helm chart. The three-part
upstream model — sources, providers, registry — is exactly the shape of the
typed spec.

## Upstream Architecture

**Sources** are pluggable watchers. The chart default is `service` +
`ingress`; Gateway API route sources (`gateway-httproute` and friends),
mesh sources (Istio, Contour, Kong, ...), and the `crd` source (declarative
DNSEndpoint objects) are all selected by name. The spec's `sources` list is
a validated enum of upstream's source names.

**Providers** come in two tiers. In-tree providers (AWS, Google, Azure,
Cloudflare, and a few others) compile into the controller binary. Everything
else rides the **webhook provider architecture**: the controller runs with
`--provider=webhook` and the actual provider implementation runs as a
sidecar container serving a small HTTP API on localhost. Upstream is
actively moving providers OUT of tree onto this architecture — it is the
extension point, not a workaround. The spec models the four major in-tree
providers as typed arms plus a `webhook` arm that takes the sidecar image,
args, env, and resources.

**The registry solves ownership.** DNS zones are shared, mutable stores
with no native concept of "who created this record." The default `txt`
registry writes a companion TXT record per managed name carrying an owner
ID; on every sync the controller only touches records whose TXT says
`txt_owner_id`. This is what makes the `sync` policy (which deletes) safe
in shared zones, and why every instance sharing a zone MUST have a distinct
owner ID. `txt_prefix`/`txt_suffix` reshape the TXT names (a TXT record
cannot coexist with a CNAME of the same name — prefixing sidesteps the
collision). The `dynamodb` registry moves ownership state into a DynamoDB
table for zones that must stay free of TXT metadata; `noop` disables
ownership entirely (exclusive zones only); `aws-sd` targets AWS Cloud Map.

**Policies dial destructiveness.** `upsert-only` creates and updates but
never deletes — the safe default for zones with records managed elsewhere.
`sync` fully reconciles (including deletes of owned records) — correct for
dedicated zones. `create-only` never touches a record twice.

## The Dual-Provider Split

Kubernetes components run IN one environment and often talk TO another.
ExternalDNS is the sharpest example: the cluster's host (EKS/GKE/AKS/
self-managed) and the DNS provider (Route 53/Cloud DNS/Azure DNS/
Cloudflare) are independent axes. The spec keeps them independent:
`dns_provider` selects the destination; `workload_identity` (the shared
per-cloud oneof) plus each arm's static credentials select authentication.
EKS + Cloudflare needs only a token; GKE + Route 53 needs static AWS keys
(no federation exists); same-cloud pairings go keyless through workload
identity or the node's ambient identity. Static credentials never land in
chart values: the modules materialize them as deterministically-named
Kubernetes Secrets (`<name>-cloudflare-credentials`, `<name>-aws-credentials`,
`<name>-gcp-credentials`, `<name>-azure-config`) and wire env/volume
references. Azure is config-file-shaped upstream: the module renders
`azure.json` (identity mode included: service principal > Workload Identity
> managed identity) and mounts it at the controller's default config path.

## Multiple Instances Are the Norm

One installation manages one provider, so multi-provider clusters (or
clusters with distinct ownership boundaries — public vs. internal zones,
per-team instances via `namespaced`) run several installations. The release
name, chart fullname, and controller ServiceAccount name are all pinned to
`metadata.name` — instances coexist without name collisions, and every
chart object carries a deterministic, manifest-derived name that
verification and cloud-side identity bindings key off. This is the opposite
of single-installation components (cert-manager, ESO) whose release names
are fixed; the difference is architectural: ExternalDNS has no cluster-
scoped CRDs or webhooks to fight over.

## Typed Surface vs Escape Hatch

The typed spec covers: provider arms (with per-arm zone filters as
StringValueOrRef lists referencing the zone kinds), workload identity,
sources, policy, registry + TXT identity + DynamoDB registry settings,
domain/annotation/label filtering, managed record types, interval and
event-triggered reconciliation, namespaced mode, logging, resources,
scheduling (node selector, tolerations, priority class), Prometheus
(opt-in ServiceMonitor), and image overrides (air-gapped mirrors,
controller tag pinning).

`helm_values` merges LAST with Helm `-f` semantics on both engines
(Terraform natively via the two-document values list; Pulumi module-side
with the same deep-merge). Deliberately unmodeled as typed fields:

- **Out-of-tree providers** — they ride the `webhook` arm exactly as
  upstream prescribes; typing each provider's flags would chase a moving
  target upstream itself is decomposing
- **DNSEndpoint / the `crd` source machinery** — expressible via `sources:
  [crd]`; the DNSEndpoint objects themselves are workload-side resources,
  not installation configuration
- **extraArgs/extraEnv/extraVolumes** — the modules use them internally for
  provider wiring; user-level additions belong in `helm_values`
- **The chart's securityContext blocks** — upstream defaults are correct;
  overriding them is an expert move that belongs in `helm_values`

## Install Semantics

Both engines install a REAL Helm release and wait for the controller
Deployment to become Available (atomic + cleanup-on-fail: a failed install
never leaves a half-deployed controller). One deliberate boundary: the
controller validates provider CREDENTIALS at first zone sync, not at
startup — an install with wrong credentials goes green and surfaces in
controller logs, matching upstream behavior. Liveness is the install's
contract; credential validity is the sync loop's.

## Upgrade Posture

`chart_version` pins the chart (default: the validated pin). Chart releases
are cut separately from the controller: chart 1.21.x ships controller
v0.21.x, with the chart's appVersion deciding the image tag unless
`image_tag` overrides it — the knob for holding the controller back or
rolling it forward independently. Upgrades re-run the release with the new
chart; there are no CRDs to migrate on the default arms (the `crd` source's
DNSEndpoint CRD is the one exception, and it ships with the chart when that
source is enabled).

## Outputs as Composition Seams

`service_account_name` (pinned to `metadata.name`) — what the cloud side
must trust for keyless access: the IRSA trust policy's
`system:serviceaccount:<namespace>:<name>` subject, the GKE Workload
Identity member, the Azure federated credential subject. `release_name` —
the handle for verification and Helm-level operations. `namespace` — where
the controller and its credential Secrets live.

## E2E

The kind-cluster lanes prove install machinery: chart-default and tuned
installs on the in-memory provider (full watching/sync surface without a
real zone), the keyless AWS arm (provider=aws, no credentials — the
controller wraps Route 53 failures as soft errors, so Available + retrying
is the asserted state). Real record writes ride the real-cluster lanes with
cloud identity.

# KubernetesClusterIssuer: Research and Design

## Introduction

A ClusterIssuer is the cert-manager custom resource that represents a signing
authority serving Certificates in EVERY namespace of a cluster. This
component creates one ClusterIssuer, named after the resource
(`metadata.name`), plus the credential Secrets its configuration requires.

## Design Authority

Designed field-by-field from the pinned cert-manager upstream (CRD schema
`cert-manager.io_clusterissuers.yaml` and the ACME/issuer Go types). The
upstream IssuerSpec is IDENTICAL for Issuer and ClusterIssuer — the Planton
specs share one `CertManagerIssuerConfig` proto and one Pulumi spec builder,
so the two kinds are structurally incapable of drifting.

## The Grain: One Issuer, Named After the Resource

The issuer's identity is `metadata.name` — the same convention every other
kind follows. Certificates select it by name (`issuerRef`), and the
ingress-shim annotation (`cert-manager.io/cluster-issuer`) carries the same
name. Nothing constrains an issuer to one DNS domain: solver SELECTORS
(`dns_zones`, `dns_names`, `match_labels`) scope challenge strategies within
a single issuer, which is upstream's own model — one "letsencrypt-production"
with per-zone solvers, not one issuer per zone.

## Backend Coverage (and the one deliberate exclusion)

| Backend | Coverage |
|---|---|
| ACME | Full: multiple solvers with selectors; HTTP-01 (Ingress + Gateway HTTPRoute); DNS-01 across Cloudflare, Route53, Azure DNS, Cloud DNS, DigitalOcean, RFC 2136, acme-dns, Akamai, and the `webhook` extension point; EAB; profiles; preferred chain; account-key controls |
| CA | Full: Secret-held keypair (FK default: a KubernetesCertificate's secret output), CRL/OCSP/AIA URLs |
| SelfSigned | Full |
| Vault | Full: token / AppRole / Kubernetes auth methods, enterprise namespaces, private CA bundles |

**Venafi is deliberately not modeled**: it drives a proprietary commercial
platform (TPP / Venafi-as-a-Service) whose users are better served by that
platform's own tooling — a desirability exclusion, recorded here rather than
silently dropped.

**ACME requires at least one solver** (stricter than the live API, which
accepts a solver-less issuer that can never satisfy any challenge — an
always-misconfiguration rejected at validation instead of hanging at
issuance).

## Credential Model

Wherever upstream expects a `secretRef` to a pre-created Secret, the spec
carries the credential VALUE, marked sensitive. Modules materialize each
credential as a Secret named `<resource-name>-<purpose>` (per-solver:
`<resource-name>-solver<N>-<provider>`) in cert-manager's cluster-resource
namespace — the only namespace cert-manager reads Secrets from for
cluster-scoped resources — and wire the CR's secretRefs to them. Identical
names and data keys on both engines. Keyless paths (Route53 IRSA, Cloud DNS
Workload Identity, Azure Managed Identity) leave credential fields empty and
inherit the controller identity configured on KubernetesCertManager.

## Engine Mechanics

- **Pulumi**: the shared spec builder renders the CRD-JSON spec map; the CR
  applies as a CustomResource. cert-manager's validating webhook checks the
  applied spec strictly, and the kind-cluster E2E lanes exercise the arms
  live — shape errors fail loudly.
- **Terraform**: `kubectl_manifest` (alekc/kubectl) applies the CR — no
  cluster connection at plan time, so an issuer can be PLANNED before
  cert-manager's CRDs exist (single-run infra charts, offline plan proofs).
  The hashicorp `kubernetes_manifest` resource was deliberately migrated
  away from: it requires the CRD's schema from the live cluster at plan.
- **Neither engine waits for Ready**: readiness depends on external
  reachability (ACME server, Vault, DNS) that is not part of applying the
  resource — the never-block-on-a-controller posture.

## Validation Highlights

Exactly-one contracts (backend, solver challenge type, DNS provider,
Cloudflare credential form, Vault auth method), paired-field contracts (EAB
key id + HMAC key, TSIG name + secret, Azure SP triplet), and vocabulary
checks (cname_strategy, Azure zone type/environment) — each with a teaching
message naming the fix.

## E2E

Self-signed reaches Ready with zero external dependencies — the
deterministic kind-cluster lane. ACME/DNS-01 lanes need a real domain plus
DNS credentials and ride the batched real-cluster/Cloudflare lanes.

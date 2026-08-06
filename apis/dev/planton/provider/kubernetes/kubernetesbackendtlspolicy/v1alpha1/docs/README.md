# KubernetesBackendTlsPolicy -- Research Documentation

This document explains why `KubernetesBackendTlsPolicy` exists, how it maps
the upstream Gateway API `BackendTLSPolicy` into the Planton component
model, and the reasoning behind its spec, validation, and IaC. It is the
source-of-truth companion to the user-facing `README.md` (getting started)
and `catalog-page.md` (catalog listing).

## Table of Contents

1. Introduction
2. Standard Channel, v1
3. Where BackendTLSPolicy Sits
4. Anatomy of BackendTLSPolicySpec
5. The Trust-Anchor Model
6. Hostname and Subject Alternative Names
7. Validation Rules
8. Faithful-Projection Details
9. Why BackendTLSPolicy Is a First-Class Planton Component
10. Composing in Infra Charts
11. Implementation Landscape
12. Common Pitfalls
13. Conclusion
14. References

## Introduction

The Gateway API models north-south traffic in layers: Gateways terminate
client TLS, routes decide where requests go -- and `BackendTLSPolicy`
decides how the LAST hop, gateway to backend, is secured. It instructs the
gateway implementation to originate TLS toward the backend Service and
VERIFY the certificate the backend presents, closing the plaintext gap
between the gateway and the workload. It is a Direct policy attachment:
rather than being referenced by routes, it targets the backend Service
itself (Core support), optionally narrowed to one named Service port.

`KubernetesBackendTlsPolicy` brings that resource into Planton as a
first-class deployment component with 100% fidelity to the upstream
Gateway API v1.6.1 `BackendTLSPolicySpec`
(`kubernetes-sigs/gateway-api` `apis/v1/backendtlspolicy_types.go`),
standard channel. There are NO deliberately unmodeled surfaces: the full
CRD spec surface -- targetRefs with sectionName, both trust-anchor arms,
hostname, subjectAltNames, and options -- is typed.

## Standard Channel, v1

BackendTLSPolicy is served as `gateway.networking.k8s.io/v1` in the
standard channel (the CRD ships in the standard-channel manifest set). The
`v1alpha3` version is deprecated upstream and no longer served -- the CRD
still declares it for schema history, with `served: false`. Consequences:

- The default standard-channel install of `KubernetesGatewayApiCrds`
  (Gateway API v1.6.1) includes the BackendTLSPolicy CRD; no experimental
  channel is required.
- The typed crd2pulumi resource is `gatewayv1.NewBackendTLSPolicy`, and
  the Terraform module applies apiVersion
  `gateway.networking.k8s.io/v1`.

## Where BackendTLSPolicy Sits

```
KubernetesGatewayApiCrds        (install standard-channel CRDs)
        |
KubernetesGateway               (terminates client TLS)
        |
KubernetesHttpRoute / routes    (decide WHERE traffic goes)
        |
backend Service  <---- KubernetesBackendTlsPolicy (this component:
        |                       HOW the gateway-to-backend hop is secured)
   backend pods (serving TLS)
```

The policy is deployed NEXT TO the backend, not next to the gateway: a
BackendTLSPolicy can only target Services in ITS OWN namespace (upstream
forbids cross-namespace targetRefs for this policy), and its CA references
are same-namespace too.

## Anatomy of BackendTLSPolicySpec

The Planton spec flattens the upstream `BackendTLSPolicySpec` after the
standard namespaced envelope (`namespace`):

- `target_refs` -- 1 to 16 references to the backend Services the policy
  applies to (upstream `LocalPolicyTargetReferenceWithSectionName`). Core
  support targets a Service; `section_name` names a Service PORT and
  narrows the policy to connections to that port (omitted, the policy
  covers the entire resource). Upstream notes implementations SHOULD
  support a single targetRef -- multiple entries are accepted by the API,
  but one is the safest portable shape.
- `validation` -- how the gateway validates the backend's TLS handshake:
  the trust anchor, the hostname, and optional SANs (upstream
  `BackendTLSPolicyValidation`).
- `options` -- up to 16 implementation-specific TLS options (upstream
  AnnotationKey/AnnotationValue constraints: domain-prefixable keys up to
  253 chars, values up to 4096).

Reference names are foreign keys: `target_refs[].name` defaults to
`KubernetesService` (field path `status.outputs.service_name`) and
`ca_certificate_refs[].name` defaults to `KubernetesConfigMap` (field path
`status.outputs.configmap_name`), both as `StringValueOrRef`. In infra
charts these wire with `valueFrom` and give the resource graph real
dependency edges; when the referent is not Planton-managed, the literal
name is passed with `value:`.

## The Trust-Anchor Model

The trust anchor comes from exactly one of two arms (the CRD enforces
exactly-one-of with two CEL rules; the spec mirrors both):

- **`ca_certificate_refs`** (bring your own CA, 1-8 refs): same-namespace
  objects carrying the PEM CA bundle that signs the backend's serving
  certificate. Core support is ONE ConfigMap with the bundle in a key
  named `ca.crt` -- exactly what a cert-manager CA chain materializes
  (the root Certificate's ConfigMap or a trust-manager Bundle target is
  the natural referent). More refs, other kinds, or multi-certificate
  bundles are implementation-specific.
- **`well_known_ca_certificates`**: trust a well-known CA set instead.
  `System` (the one upstream-defined value) trusts the implementation's
  system certificate store -- the arm for backends serving
  publicly-issued certificates. Implementations may define their own
  domain-prefixed sets; upstream support for this arm is
  Implementation-specific.

## Hostname and Subject Alternative Names

`validation.hostname` is mandatory and does double duty (upstream Core):
it is sent as the SNI (RFC 6066) for the backend connection, and -- unless
`subject_alt_names` is set -- it is the identity the backend certificate
must prove.

`subject_alt_names` (up to 5, Extended support) covers the case where the
certificate identity differs from the SNI hostname -- the SPIFFE-ID
pattern in mTLS meshes is the common one. Each entry is a closed-enum
`type` (`Hostname` | `URI`) with exactly the matching value field; with
SANs set, the hostname only SELECTS the certificate, and must be added as
a `Hostname` SAN if it should also authenticate.

## Validation Rules

Every upstream `XValidation` and kubebuilder marker is translated to
`buf.validate`, so manifests fail at validation time with the same rules
the API server would apply:

- `target_refs`: `min_items: 1`, `max_items: 16`, plus the CRD's own two
  same-target CEL rules -- refs to the same target must all carry a
  `section_name` (or none of them), and no two refs to the same target
  may carry the same `section_name`.
- `validation`: the two trust-anchor CEL rules (must not both be set; one
  must be set); `ca_certificate_refs` max 8; `subject_alt_names` max 5.
- SAN entries: the CRD's four type/value pairing rules -- `hostname`
  required for type `Hostname` and forbidden otherwise; `uri` required
  for type `URI` and forbidden otherwise -- plus the upstream hostname
  (wildcard-first-label allowed) and URI (scheme mandatory) patterns.
- `options`: max 16 pairs, AnnotationKey pattern on keys, 4096-char cap
  on values.
- Group/Kind/SectionName/Hostname fields carry the upstream character
  patterns and length bounds verbatim.

## Faithful-Projection Details

Two spec decisions exist purely so the rendered CR is byte-faithful to
what the CRD demands:

- **`group` is presence-required but may be empty.** The CRD requires the
  `group` KEY on targetRefs and caCertificateRefs but its Group type
  explicitly allows the empty value (Services and ConfigMaps live in the
  core API group -- the empty string). The proto models this as an
  `optional` string with a presence `required` rule: it must be SET, and
  because it is `optional`, protojson does not drop the empty string from
  the projection. Each engine guarantees emission its own way: the Pulumi
  module always sets `Group:` from the resolved value, and the Terraform
  module's converter emits `group: ""` which the module's null-prune (not
  empty-prune) passes through to the CR.
- **`wellKnownCACertificates` casing is pinned via `json_name`.**
  protojson would derive `wellKnownCaCertificates` from the field name,
  but the CRD's key is `wellKnownCACertificates` (capital CA) and the API
  server rejects the miscased key as undeclared. The spec pins the exact
  key with `json_name` -- the faithful-projection contract includes
  acronym casing.

One structural choice follows the same logic: the target and CA reference
messages are deliberately LOCAL to this kind rather than reusing the
shared `KubernetesGatewayApiLocalObjectReference`, whose plain-string name
serves arbitrary extension objects -- a policy target's name is a real
foreign key to the Service it secures, and a CA reference's name to the
ConfigMap carrying the bundle.

## Why BackendTLSPolicy Is a First-Class Planton Component

Without it, gateway-to-backend TLS forces customers back to raw
`KubernetesManifest` YAML: no proto validation (the same-target and
trust-anchor rules would only fail at the API server), no typed SDKs, no
FK wiring to the Service and ConfigMap, no InfraChart composability, no UI
wizards. `KubernetesBackendTlsPolicy` closes the end-to-end-encryption gap
in the Gateway API family alongside the route components.

## Composing in Infra Charts

Every neighbor reference is a `StringValueOrRef` foreign key:

- `namespace` (default kind `KubernetesNamespace`),
- `target_refs[].name` (default kind `KubernetesService`),
- `ca_certificate_refs[].name` (default kind `KubernetesConfigMap`).

Wiring them with `valueFrom` creates real dependency edges, so the policy
deploys after the Service it secures and the ConfigMap carrying its trust
bundle -- no manual relationship declarations. A cert-manager CA chain
composes naturally: the chain materializes the CA bundle ConfigMap, and
the policy's `ca_certificate_refs` points at it. When the target or
referent is not Planton-managed, the literal name is passed with `value:`
-- no edge is created, matching reality.

## Implementation Landscape

Support for BackendTLSPolicy varies more than for the core route kinds;
the Gateway API conformance suite tracks it as its own feature
(`BackendTLSPolicy`). Facts verifiable from implementation sources:

- **Istio** implements BackendTLSPolicy in its Gateway API controller,
  including `sectionName` for port-specific TLS and `ServiceEntry` as an
  additional targetRef kind; its agentgateway controller translates
  BackendTLSPolicy as well.
- **Cilium's** Gateway API implementation does NOT declare
  BackendTLSPolicy support -- the feature is on its GatewayClass
  supported-features exempt list.

Confirm your controller before relying on the policy: an unimplemented
policy applies cleanly and secures nothing. The per-ancestor
Accepted/ResolvedRefs conditions appear when a Gateway implementation
reconciles the resource, which is observed via kubectl rather than awaited
at apply time.

## Common Pitfalls

- **Same-namespace only.** Cross-namespace targetRefs and CA references
  are upstream-invalid for this policy -- create it in the backend's
  namespace.
- **The backend must actually serve TLS.** The policy makes the gateway
  originate and verify TLS; a plaintext backend then fails the handshake.
- **`hostname` must match the certificate** (unless SANs are set): it is
  both the SNI and the identity check, so it must appear in the backend
  certificate's DNS SANs.
- **With `subjectAltNames`, the hostname no longer authenticates.** Add it
  as a `Hostname` SAN entry if it should.
- **`sectionName` must name a real Service port** -- otherwise the policy
  fails to attach (surfaced through ResolvedRefs).
- **`group` must be present-but-empty** for core-group referents; the CRD
  rejects a missing key.
- **Controller support varies.** Verify the implementation honors the
  policy -- and for `wellKnownCACertificates` and `subjectAltNames`,
  which are Implementation-specific/Extended, verify those arms too.
- **Literal `value:` names create no dependency edge.** In infra charts,
  prefer `valueFrom` for Planton-managed Services and ConfigMaps so
  ordering is guaranteed.

## Conclusion

`KubernetesBackendTlsPolicy` completes the end-to-end-encryption slice of
the Planton Gateway API layer: the full upstream v1 surface, the CRD's own
validation rules enforced at manifest time, and both trust-anchor arms
composable with cert-manager chains and mesh trust bundles through real
foreign keys.

## References

- [Gateway API BackendTLSPolicy](https://gateway-api.sigs.k8s.io/api-types/backendtlspolicy/)
- Upstream types: `kubernetes-sigs/gateway-api` `apis/v1/backendtlspolicy_types.go` (v1.6.1)
- Upstream CRD: `config/crd/standard/gateway.networking.k8s.io_backendtlspolicies.yaml`
- Sibling components: `KubernetesGateway`, `KubernetesHttpRoute`, `KubernetesGatewayApiCrds`

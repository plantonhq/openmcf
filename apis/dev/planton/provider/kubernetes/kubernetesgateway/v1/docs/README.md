# KubernetesGateway -- Research Documentation

Comprehensive background on the Kubernetes Gateway API `Gateway` resource, the
upstream specification this component mirrors, and the reasoning behind the
Planton `KubernetesGateway` component.

## Table of Contents

1. [Introduction](#introduction)
2. [The Gateway API Role-Oriented Model](#the-gateway-api-role-oriented-model)
3. [Where Gateway Sits](#where-gateway-sits)
4. [Anatomy of GatewaySpec](#anatomy-of-gatewayspec)
5. [Listener Distinctness and TLS Rules](#listener-distinctness-and-tls-rules)
6. [Why Gateway Is a First-Class Planton Component](#why-gateway-is-a-first-class-planton-component)
7. [Design Notes](#design-notes)
8. [Controller Landscape](#controller-landscape)
9. [80/20 Scoping](#8020-scoping)
10. [Common Pitfalls](#common-pitfalls)
11. [Conclusion](#conclusion)
12. [References](#references)

## Introduction

The Kubernetes Gateway API is the successor to the Ingress API. It is a
standards-based, role-oriented, expressive specification for configuring L4/L7
traffic routing into a cluster, developed by SIG-Network and implemented by 20+
controllers (Istio, Envoy Gateway, NGINX, Cilium, Traefik, Kong, and the major
cloud load balancers). It reached GA (`v1`) in October 2023 and continues to
evolve through quarterly releases; this component tracks **v1.6.1**.

A `Gateway` is the central object: it represents an instance of traffic-handling
infrastructure (typically a load balancer or proxy deployment) and declares the
**listeners** -- the ports, protocols, and TLS settings -- on which it accepts
connections. Routes then attach to a Gateway to describe how matching traffic is
dispatched to backends.

## The Gateway API Role-Oriented Model

The Gateway API deliberately separates concerns across three personas, each
owning a distinct resource:

```
Infrastructure Provider          Cluster Operator              Application Developer
        │                              │                               │
        ▼                              ▼                               ▼
   GatewayClass  ───selects───►     Gateway      ◄───parentRef───   HTTPRoute / GRPCRoute
   (controller)                  (listeners,                        TLSRoute / TCPRoute
                                  addresses, TLS)                   (host/path/backends)
```

- **GatewayClass** (infrastructure provider): names the controller implementation.
- **Gateway** (cluster operator): owns the public entry point -- ports, protocols, TLS, and which routes/namespaces may attach.
- **Routes** (application developer): describe routing rules and attach to a Gateway via `parentRefs`.

`KubernetesGateway` is the Planton representation of the middle layer. It pairs
with `KubernetesGatewayClass` and the route components (`KubernetesHttpRoute`,
`KubernetesGrpcRoute`, `KubernetesTlsRoute`, `KubernetesTcpRoute`,
`KubernetesUdpRoute`), plus `KubernetesListenerSet` for delegated listener
management.

## Where Gateway Sits

```
KubernetesGatewayApiCrds        (install CRDs + controller)
        │
        ▼
KubernetesGatewayClass          (controller class)  ──gatewayClassName FK──┐
        │                                                                  │
KubernetesNamespace ──namespace FK──┐                                      │
                                    ▼                                      ▼
                              KubernetesGateway  (listeners, TLS, addresses)
                                    ▲          ▲
        KubernetesCertificate ──────┘          └────── HTTPRoute / GRPCRoute / TLSRoute / TCPRoute
        (TLS Secret for HTTPS                          (attach via parentRefs)
         listener certificateRefs)
```

## Anatomy of GatewaySpec

The upstream `GatewaySpec` (standard channel, v1.6.1) has the following fields,
all mirrored in `KubernetesGatewaySpec` after the Planton envelope
(`namespace`):

| Upstream field | Planton field | Notes |
|----------------|---------------|-------|
| `gatewayClassName` | `gateway_class_name` | Foreign key to `KubernetesGatewayClass`. |
| `listeners` | `listeners` | 1-64 listeners; the heart of the spec. |
| `addresses` | `addresses` | Up to 16 requested addresses (Extended support). |
| `infrastructure` | `infrastructure` | Labels (max 8) / annotations (max 16) / parametersRef for created resources. |
| `allowedListeners` | `allowed_listeners` | Which ListenerSets may attach (defaults to none). |
| `tls` | `tls` | Gateway-wide frontend (mutual TLS) and backend client-cert config. |
| `defaultScope` | *(excluded)* | Experimental; absent from the standard CRD. See Design Notes. |

### Listener

A listener is the core unit: `name`, `hostname`, `port`, `protocol`, `tls`, and
`allowedRoutes`. The protocol determines which fields are meaningful:

- **HTTP** -- cleartext; no TLS, hostname-aware.
- **HTTPS** -- TLS terminated at the Gateway; requires a certificate.
- **TLS** -- TLS either terminated or passed through (for TLSRoute); requires `mode`.
- **TCP / UDP** -- connection/datagram forwarding; no hostname, no TLS.

### Per-listener TLS vs gateway-wide TLS

There are two distinct TLS concepts, and the component models both:

- **Listener `tls`** (`ListenerTLSConfig`): how an individual HTTPS/TLS listener
  terminates or passes through TLS, including the certificate Secrets it serves.
- **Gateway `tls`** (`GatewayTLSConfig`): gateway-wide **frontend** client
  certificate validation (mutual TLS for inbound connections, with optional
  per-port overrides) and **backend** client-certificate material the Gateway
  presents when connecting to upstreams.

## Listener Distinctness and TLS Rules

The upstream spec enforces several cross-cutting rules via CEL `XValidation`.
`KubernetesGateway` translates each one faithfully into `buf.validate` so the
same errors surface at author time in Planton UIs and CLIs:

- Each listener `name` is unique within the Gateway.
- The combination of `port`, `protocol`, and `hostname` is unique across listeners.
- `tls` must not be set when `protocol` is HTTP, TCP, or UDP.
- An HTTPS listener may only `Terminate` (mode unset or `Terminate`).
- A TLS listener must declare its `tls.mode`.
- `hostname` must not be set for TCP/UDP listeners.
- A `Terminate` listener must provide `certificateRefs` or `options`.
- Requested IPAddress and Hostname address values are each unique.
- Per-port frontend TLS overrides target unique ports.

## Why Gateway Is a First-Class Planton Component

Without this component, customers wanting Gateway API ingress are forced to use
`KubernetesManifest` (raw YAML), which sacrifices:

1. **Proto validation** -- the distinctness and TLS rules above are enforced before apply.
2. **Foreign-key wiring** -- `namespace` and `gateway_class_name` reference other Planton resources, enabling InfraChart DAG ordering and Planton UI resource pickers.
3. **Typed IaC** -- both the Pulumi and Terraform modules construct the CRD from typed inputs, catching structural errors at compile/plan time.
4. **Composability** -- the Gateway's outputs (`gateway_name`, `namespace`) feed Route components that attach to it.

## Design Notes

- **Complete standard-channel surface.** Every field of the standard-channel
  `GatewaySpec` is represented. The only omission is the experimental
  `defaultScope` field, which is absent from the standard CRD and from the typed
  Pulumi resource Planton provisions with; including it would have no
  deployable target.
- **Value fields are strings, validated per upstream kind.** Gateway API models
  its enums as string type aliases, and the proto mirrors that. Open-set values
  (`protocol`, address `type`) are validated with the upstream **regex**,
  preserving custom domain-prefixed values. Closed enums (`tls.mode`, namespace
  `from`, frontend validation `mode`, selector `operator`) are validated with
  CEL membership checks. Both keep exact upstream casing, so no case-mapping is
  needed in the IaC layer, CEL string comparisons translate verbatim, and
  controller-specific values keep working.
- **No baked-in Planton defaults for upstream fields.** Upstream kubebuilder
  defaults (e.g. `tls.mode=Terminate`, address `type=IPAddress`, route
  `from=Same`) are controller/CRD-enforced and are documented in field comments
  rather than set as Planton defaults, so controller behavior is never
  second-guessed.
- **Reference names are foreign keys.** Each `certificate_refs[].name` is a
  `StringValueOrRef` defaulting to `KubernetesSecret` -- the Secret is
  typically produced by a `KubernetesCertificate` (the cert-manager seam), so
  `valueFrom` against the Certificate's `status.outputs.secret_name` wires the
  listener to the issued certificate with a real dependency edge. The frontend
  mTLS `ca_certificate_refs[].name` similarly defaults to
  `KubernetesConfigMap`. Referents outside Planton pass literal names with
  `value:`.
- **Typed crd2pulumi resources.** The Pulumi module uses `gatewayv1.NewGateway`,
  not an untyped `CustomResource`. The Terraform module applies through
  `kubectl_manifest` (alekc/kubectl) with server-side apply, plannable before
  the Gateway API CRDs exist.

## Controller Landscape

`gatewayClassName` ultimately resolves to a controller. Common choices:

| Controller | Typical GatewayClass controllerName |
|------------|-------------------------------------|
| Istio | `istio.io/gateway-controller` |
| Envoy Gateway | `gateway.envoyproxy.io/gatewayclass-controller` |
| NGINX Gateway Fabric | `gateway.nginx.org/nginx-gateway-controller` |
| Cilium | `io.cilium/gateway-controller` |

The Gateway itself is controller-agnostic; the GatewayClass selects the
implementation, and implementation-specific tuning flows through
`infrastructure.parametersRef` or listener `tls.options`.

## 80/20 Scoping

- **In scope:** the full standard-channel `GatewaySpec` -- listeners (all five
  core protocols), per-listener TLS, allowed routes, requested addresses,
  infrastructure metadata, gateway-wide frontend/backend TLS, and allowed
  listeners (the opt-in gate for `KubernetesListenerSet` attachment).
- **Out of scope:** the experimental `defaultScope` field; CRD `status`
  (reconciled asynchronously by the controller and observed via kubectl, not
  stored in stack outputs).

## Common Pitfalls

- **Setting `tls` on an HTTP/TCP/UDP listener.** TLS only applies to HTTPS/TLS
  listeners; the spec rejects it otherwise.
- **Forgetting a certificate on a Terminate listener.** A terminating listener
  must reference at least one TLS Secret (or supply implementation `options`).
- **Duplicate listeners.** Two listeners may not share the same name, nor the
  same port/protocol/hostname combination.
- **Expecting status in outputs.** Assigned addresses and listener conditions
  are controller-managed; query them with `kubectl get gateway`. Stack outputs
  expose only the stable identifiers (`gateway_name`, `namespace`,
  `gateway_class_name`).
- **Missing prerequisites.** The Gateway will not program until the Gateway API
  CRDs and a controller-backed GatewayClass are present.

## Conclusion

`KubernetesGateway` brings the central Gateway API object into Planton at full
standard-channel fidelity, with validation, typed IaC, and foreign-key
composability. Together with `KubernetesGatewayClass` and the route components,
it completes a declarative, controller-agnostic ingress layer that customers can
compose in InfraCharts without dropping to raw YAML.

## References

- [Gateway API: Gateway](https://gateway-api.sigs.k8s.io/api-types/gateway/)
- [Gateway API: TLS configuration](https://gateway-api.sigs.k8s.io/guides/tls/)
- [Gateway API v1.6.1 specification](https://github.com/kubernetes-sigs/gateway-api/tree/v1.6.1)
- [Pulumi Kubernetes Provider](https://www.pulumi.com/registry/packages/kubernetes/)

## Composing in Infra Charts

`KubernetesGateway` is the mid-tier hub of the ingress DAG: it depends on a
GatewayClass and (for TLS) a certificate Secret, and routes attach to it. Every
neighbor reference is a `StringValueOrRef` foreign key, so `valueFrom` wiring
creates real dependency edges and the platform builds the DAG automatically:

- `namespace` -> `KubernetesNamespace.spec.name`
- `gateway_class_name` -> `KubernetesGatewayClass.status.outputs.gateway_class_name`
- `listeners[].tls.certificate_refs[].name` -> `KubernetesSecret`, typically via
  a `KubernetesCertificate`'s `status.outputs.secret_name` (the cert-manager
  seam: the listener terminates with the issued certificate the moment it
  exists)
- `tls.frontend.*.validation.ca_certificate_refs[].name` -> `KubernetesConfigMap`

```yaml
metadata:
  name: "{{ values.env }}-gateway"
spec:
  namespace:
    valueFrom:
      kind: KubernetesNamespace
      name: "{{ values.env }}-ns"
      fieldPath: spec.name
  gateway_class_name:
    valueFrom:
      kind: KubernetesGatewayClass
      name: "{{ values.env }}-gateway-class"
      fieldPath: status.outputs.gateway_class_name
  listeners:
    - name: https
      port: 443
      protocol: HTTPS
      tls:
        mode: Terminate
        certificate_refs:
          - name:
              valueFrom:
                kind: KubernetesCertificate
                name: "{{ values.domain }}-cert"
                fieldPath: status.outputs.secret_name
```

When a referenced Secret or ConfigMap is not Planton-managed, pass the literal
name with `value:` instead.

Full ingress stack:
`CertManager -> ClusterIssuer -> Certificate -> (Secret) -> Gateway -> HTTPRoute / GRPCRoute`.

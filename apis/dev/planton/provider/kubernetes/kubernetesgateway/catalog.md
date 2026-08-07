# Gateway on Kubernetes

Creates a namespaced Kubernetes Gateway API `Gateway` -- an instance of traffic-handling infrastructure that binds a set of **listeners** (logical endpoints with a port, protocol, and optional TLS) to addresses, programmed by the controller behind a `GatewayClass`. The Gateway is the role-oriented successor to Ingress: platform teams own the `GatewayClass`, infrastructure teams own the `Gateway`, and application teams attach `Routes` (HTTP/gRPC/TLS/TCP) to its listeners. This component mirrors the upstream Gateway API v1 `Gateway` spec with full fidelity while adding proto validation, typed SDKs, and InfraChart composability.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **A namespaced Gateway** named after `metadata.name` in `spec.namespace`, attached to the `spec.gatewayClassName` GatewayClass, with the listeners (and any requested addresses, infrastructure attributes, and gateway-wide TLS) you configure.
- **Kubernetes Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking.

The controller behind the GatewayClass reconciles the Gateway asynchronously: it assigns addresses and reports per-listener conditions in the Gateway's status, which you observe with `kubectl` (these controller-managed values are intentionally not stored as stack outputs).

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with kubeconfig credentials for the target cluster. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline kubeconfig authentication.

### Kubernetes Cluster

- **Gateway API CRDs installed** -- deploy the `KubernetesGatewayApiCrds` component first. Gateway is a Gateway API type and will not register without the CRDs.
- **A controller-backed GatewayClass exists** -- `spec.gatewayClassName` must resolve to a `KubernetesGatewayClass` whose controller (Istio, Envoy Gateway, NGINX, ...) is installed and running.
- **The target namespace exists** -- `spec.namespace` should resolve to a real `KubernetesNamespace`.
- **Referenced TLS Secrets exist** -- for HTTPS/TLS listeners in `Terminate` mode, the `certificateRefs` Secrets must exist in the Gateway's namespace (or be authorized cross-namespace by a `KubernetesReferenceGrant`). They are commonly produced by a cert-manager `KubernetesCertificate`.

## Deploy

### Console

Open the deployment store, find **Gateway on Kubernetes**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and seven spec steps: **Namespace** (immutable) and **Gateway Class** (immutable) first, then the required **Listeners**, followed by the optional **Addresses**, **Infrastructure**, **Allowed Listeners**, and **Gateway TLS** steps. Start from the **HTTPS / TLS Terminate** or **Multi-Protocol** preset in the [Presets](#presets) tab for a directly deployable configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesGateway
metadata:
  name: my-https-gateway
  org: acme-corp
  env: prod
spec:
  namespace:
    value: istio-ingress
  gatewayClassName:
    value: istio
  listeners:
    - name: https
      hostname: app.example.com
      port: 443
      protocol: HTTPS
      tls:
        mode: Terminate
        certificateRefs:
          - name: app-tls
      allowedRoutes:
        namespaces:
          from: Same
        kinds:
          - kind: HTTPRoute
```

```shell
planton apply -f gateway.yaml
```

This creates a Gateway in `istio-ingress` with one HTTPS listener on port 443 that terminates TLS using the `app-tls` Secret and accepts `HTTPRoute`s from its own namespace.

## Key Configuration

These are the most important decisions when configuring a Gateway. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Namespace** -- The `namespace` field is where the Gateway runs. It is **immutable**: Routes and TLS Secrets in this namespace attach without a `ReferenceGrant`; those elsewhere require one. Reference an existing `KubernetesNamespace` or type the name directly.

**Gateway Class** -- The `gatewayClassName` field selects the `GatewayClass`, and therefore the controller that programs the Gateway. It is **immutable** (a Gateway cannot switch controllers in place) and is a foreign key to a `KubernetesGatewayClass` output, so the class is deployed before the Gateway.

**Listeners** -- One or more `listeners` (1-64). Each is a `name` + `port` + `protocol` (HTTP, HTTPS, TLS, TCP, UDP, or a domain-prefixed custom protocol), with an optional virtual `hostname`. HTTPS/TLS listeners add **TLS termination** (`mode` Terminate/Passthrough + `certificateRefs` + `options`), and any listener may restrict attachment via **allowedRoutes** (namespaces `from` All/Selector/Same + a label selector, and the allowed Route `kinds`).

**Addresses** -- Optional requested addresses (`type` IPAddress/Hostname/NamedAddress + `value`, up to 16). When omitted, the controller assigns addresses; set them to pin a reserved static IP or hostname.

**Infrastructure** -- Optional `labels` and `annotations` (up to 8 each) propagated to the resources the controller creates -- the standard way to tune cloud load-balancer behavior -- plus an optional per-Gateway `parametersRef`.

**Allowed Listeners** -- Optional ListenerSet attachment policy (`from` All/Selector/Same/None). Defaults to None; retained for spec fidelity and forward compatibility (ListenerSet is not yet a Planton resource).

**Gateway TLS** -- Optional gateway-wide TLS: **frontend** client-certificate validation (mTLS), with a default and optional per-port overrides, and a **backend** client certificate the Gateway presents to upstreams. Per-listener HTTPS termination is configured on each listener instead.

## Outputs and Dependencies

### What This Component Consumes

This component takes foreign-key references to a `KubernetesNamespace` (via `spec.namespace`), a `KubernetesGatewayClass` (via `spec.gatewayClassName`), `KubernetesSecret` (listener `certificateRefs` and the backend client certificate -- typically a cert-manager `KubernetesCertificate`'s issued Secret through its `secret_name` output), and `KubernetesConfigMap` (frontend `caCertificateRefs` CA bundles), so an InfraChart deploys those targets before the Gateway and the resource graph carries the edges. Literal names cover material created outside Planton; cross-namespace references require a `KubernetesReferenceGrant`.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `gateway_name` | Name of the created Gateway (equals `metadata.name`) | Routes reference this in their `parentRefs`; orders the Gateway before its Routes in an InfraChart |
| `namespace` | The resolved namespace the Gateway was created in | Same-namespace / ReferenceGrant rules for attaching Routes and Secrets |
| `gateway_class_name` | The resolved GatewayClass the Gateway belongs to | Confirming which controller programs the Gateway |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**HTTPS with TLS termination** -- A single HTTPS listener on 443 that terminates TLS using a Secret from cert-manager, accepting same-namespace `HTTPRoute`s. Start from the **HTTPS / TLS Terminate** preset.

**Multi-protocol Gateway** -- An HTTP listener (often redirecting to HTTPS) alongside an HTTPS listener, sharing the Gateway's address. Start from the **Multi-Protocol** preset.

**TLS passthrough** -- A TLS listener in `Passthrough` mode that forwards the encrypted stream to a `KubernetesTlsRoute` without decrypting at the Gateway.

## Works With

- **KubernetesGatewayApiCrds** -- installs the Gateway API CRDs (prerequisite, install first).
- **KubernetesGatewayClass** -- the class (`spec.gatewayClassName`) that selects the controller (install first).
- **KubernetesNamespace** -- the namespace (`spec.namespace`) the Gateway runs in.
- **KubernetesCertificate** -- commonly produces the TLS Secrets referenced by HTTPS listeners.
- **KubernetesReferenceGrant** -- authorizes cross-namespace TLS Secret / CA references from this Gateway.
- **KubernetesHttpRoute / KubernetesGrpcRoute / KubernetesTlsRoute / KubernetesTcpRoute** -- the Routes that attach to this Gateway's listeners.

# Reference Grant on Kubernetes

Creates a namespaced Kubernetes Gateway API `ReferenceGrant` -- a runtime authorization that permits resources in *other* namespaces to reference specified kinds of resources in *this* grant's namespace. In the Gateway API, every cross-namespace reference (a Gateway's TLS `certificateRefs`, a Route's `backendRefs`, and similar) is denied by default; a ReferenceGrant placed in the *referenced* ("to") namespace is what explicitly authorizes it. This component mirrors the upstream Gateway API v1 `ReferenceGrant` spec with full fidelity while adding proto validation, typed SDKs, and InfraChart composability.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **A namespaced ReferenceGrant** named after `metadata.name` in `spec.namespace`, declaring the trusted sources (`from`) that may reference the permitted targets (`to`) in that namespace.
- **Kubernetes Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking.

ReferenceGrant has no controller-managed status upstream (the Gateway API project deliberately omitted it), so there is no `Accepted`/`Programmed` condition to wait on -- the grant takes effect as soon as it exists.

## The Trust Direction (read this first)

A ReferenceGrant is non-obvious because it lives in the **target** namespace, not the source:

- **`spec.namespace`** -- the namespace being referenced *into* (the "to" side). The grant must be created here; deleting it revokes the access.
- **`from`** -- the trusted sources: each `{ group, kind, namespace }` names a kind of resource in a source namespace that is allowed to reference in. Entries combine with OR.
- **`to`** -- the referenceable targets in this grant's namespace: each `{ group, kind, name? }`. Omit `name` to allow all resources of that kind, or set it to narrow the grant to one.

Example: a `Gateway` in `istio-ingress` needs a TLS `Secret` in `cert-manager`. The grant is created in `cert-manager` (`spec.namespace`), with `from: [{ kind: Gateway, group: gateway.networking.k8s.io, namespace: istio-ingress }]` and `to: [{ kind: Secret, group: "" }]`.

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with kubeconfig credentials for the target cluster. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline kubeconfig authentication.

### Kubernetes Cluster

- **Gateway API CRDs installed** -- deploy the `KubernetesGatewayApiCrds` component first. ReferenceGrant is a Gateway API type and will not register without the CRDs.
- **The target namespace exists** -- `spec.namespace` should resolve to a real `KubernetesNamespace`. Reference an existing one or type the namespace name directly.

## Deploy

### Console

Open the deployment store, find **Reference Grant on Kubernetes**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and three spec steps: **Namespace** (the immutable "to" namespace), **Sources** (the trusted `from` entries), and **Targets** (the referenceable `to` entries). Start from the **Allow Gateway Secret Ref** or **Allow Route Backend Ref** preset in the [Presets](#presets) tab for a directly deployable configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesReferenceGrant
metadata:
  name: allow-gateway-secret-ref
  org: acme-corp
  env: prod
spec:
  namespace:
    value: cert-manager
  from:
    - group: gateway.networking.k8s.io
      kind: Gateway
      namespace: istio-ingress
  to:
    - group: ""
      kind: Secret
```

```shell
planton apply -f reference-grant.yaml
```

This creates a ReferenceGrant in `cert-manager` that authorizes Gateways in `istio-ingress` to reference Secrets there.

## Key Configuration

These are the most important decisions when configuring a ReferenceGrant. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Namespace** -- The `namespace` field is the "to" namespace where the grant is created and which it protects. It is **immutable**: the grant boundary is defined relative to this namespace, so changing it means creating a different grant. Reference an existing `KubernetesNamespace` or type the name directly.

**Sources (`from`)** -- One or more `{ group, kind, namespace }` tuples (1-16). `kind` and `namespace` are required; `group` is the API group (use `gateway.networking.k8s.io` for Gateways and Routes, or empty `""` for core kinds). Each entry is an additional source allowed to reference in.

**Targets (`to`)** -- One or more `{ group, kind, name? }` tuples (1-16). `kind` is required; `group` is `""` for the core kinds `Secret` and `Service`. Leave `name` empty to allow all resources of the kind, or set it to restrict the grant to a single named resource.

## Outputs and Dependencies

### What This Component Consumes

This component takes a foreign-key reference to a `KubernetesNamespace` via `spec.namespace`. The `from`/`to` entries are kind-level trust assertions (not pointers to specific resource instances), so they create no deploy-ordering dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `reference_grant_name` | Name of the created ReferenceGrant (equals `metadata.name`) | Lets Gateways/Routes making the cross-namespace reference order themselves after the grant in an InfraChart |
| `namespace` | The resolved "to" namespace the grant was created in | Confirming which namespace's resources the grant authorizes inbound references to |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Allow a Gateway to reference TLS Secrets** -- A Gateway terminates TLS using a certificate Secret in another namespace (typically the cert-manager namespace). Create the grant in the Secret's namespace, with `from` the Gateway's namespace/kind and `to` `kind: Secret`. Start from the **Allow Gateway Secret Ref** preset.

**Allow Routes to reference backend Services** -- HTTP/gRPC routes in an application namespace forward traffic to backend Services in another namespace. Create the grant in the backend namespace, with one `from` entry per trusted route kind and `to` `kind: Service`. Start from the **Allow Route Backend Ref** preset.

## Works With

- **KubernetesGatewayApiCrds** -- installs the Gateway API CRDs (prerequisite, install first).
- **KubernetesNamespace** -- the "to" namespace (`spec.namespace`) and the source namespaces named in `from`.
- **KubernetesGateway** -- a common `from` source when referencing cross-namespace TLS Secrets.
- **KubernetesHttpRoute / KubernetesGrpcRoute / KubernetesTlsRoute / KubernetesTcpRoute** -- common `from` sources when referencing cross-namespace backend Services.

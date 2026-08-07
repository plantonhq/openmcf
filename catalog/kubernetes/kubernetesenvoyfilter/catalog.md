# Envoy Filter on Kubernetes

Defines an Istio EnvoyFilter: a namespaced, expert-only customization of the raw Envoy proxy configuration that Istio generates for the workloads it selects. It applies an ordered list of patches -- merge, add, remove, insert, or replace -- to low-level Envoy objects (listeners, filter chains, network and HTTP filters, route configurations, virtual hosts, routes, and clusters). This is the escape hatch you reach for only when a capability is not yet exposed by a higher-level Istio API.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **An EnvoyFilter resource** -- a namespaced Istio policy that patches the generated Envoy configuration for the workloads it selects.
- **Kubernetes Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking.

## What It Controls

- **Attachment** -- which workloads the patches apply to: namespace-wide (or mesh-wide from the Istio root namespace), a label selector for a set of pods, or target references attached to a Gateway, Service, or ServiceEntry. At most one of a selector and target references is set.
- **Priority** -- the order in which patch sets apply within a context. Root-namespace patch sets always apply before workload-namespace ones; within a set, patches apply in list order.
- **Config patches** -- the ordered list of changes. Each picks an Envoy object by **apply-to** (the object type) and an optional **match** (context plus listener / route-configuration / cluster details), then applies an **operation** with a free-form **value** and, for ADD, an optional **filter class**.

## Important Behavior

EnvoyFilter patches Envoy's internal xDS API directly. The patch value is **free-form** configuration (a `google.protobuf.Struct`) that istiod merges into the generated config with **no schema validation** -- a malformed patch can break the patched workload's traffic. Match objects as specifically as possible, and test on a narrow selector before going namespace- or mesh-wide. `REMOVE` deletes the matched object and takes no value; every other operation takes a value (the fragment to merge or insert). Prefer the ADD operation with a `filterClass` (AUTHN/AUTHZ/STATS) over the `INSERT_*` operations, which rely on potentially-unstable filter names. The workload selector is a plain runtime label match (no dependency edge) -- order this resource after the workloads it patches with `metadata.relationships`. A target reference's **name** is a foreign key defaulting to a Gateway reference: wiring it with `valueFrom` orders this resource after the gateway it patches, while a literal `value:` covers Services, ServiceEntries, and resources created outside Planton. Waypoint proxies require target references; label selectors are ignored by waypoints.

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with kubeconfig credentials for the target cluster. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline kubeconfig authentication.

### Kubernetes Cluster

- **Istio installed** -- the Istio CRDs (Istio Base CRDs on Kubernetes) must be present and the Istio control plane (istiod) running. Patches are only honored where istiod is active.
- **Target namespace exists** -- the resource is created in a specific namespace; reference an existing one or create it first.

## Deploy

### Console

Open the deployment store, find **Envoy Filter on Kubernetes**, and click **Deploy**. The creation wizard walks you through the namespace, the attachment, and the config patches, with guidance at each step.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesEnvoyFilter
metadata:
  name: reviews-idle-timeout
  org: acme-corp
  env: prod
spec:
  namespace:
    value: prod-apps
  workloadSelector:
    labels:
      app: reviews
  configPatches:
    - applyTo: CLUSTER
      match:
        context: SIDECAR_OUTBOUND
        cluster:
          service: reviews.prod.svc.cluster.local
      patch:
        operation: MERGE
        value:
          common_http_protocol_options:
            idle_timeout: 30s
```

```shell
planton apply -f envoy-filter.yaml
```

This sets a 30-second idle timeout on the outbound cluster for the `reviews` service, applied only to the `reviews` workload's sidecars in `prod-apps`.

## Key Configuration

- **Namespace** -- the namespace the resource is created in. It is fixed once created and sets the blast radius (mesh-wide from the Istio root namespace).
- **Attachment** -- All Workloads (namespace-wide), a Workload Selector (match pod labels), or Target References (attach to a Gateway/Service/ServiceEntry), plus an optional **priority**.
- **Config patches** -- each patch sets **apply-to** (LISTENER, HTTP_FILTER, CLUSTER, ...), an optional **match** (context, plus a listener / route-configuration / cluster object with its own details), and a **patch** (operation, the free-form value, and the ADD filter class).

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Description |
|------------|-------------|
| `namespace` | The namespace the resource is created in. Reference an existing Namespace on Kubernetes or supply the name directly. |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources and operators can reference:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `envoy_filter_name` | Name of the created EnvoyFilter resource (equals `metadata.name`) | Ordering resources that depend on the patch being in place |
| `namespace` | The namespace the resource was created in | Confirming where the patches apply |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

- **Insert a stock Envoy HTTP filter** -- add a native Envoy filter (CORS, ext_authz, rate limiting) into a gateway's HTTP connection manager when no first-class Istio API exists for it yet on your version.
- **Tune an upstream cluster** -- MERGE connection-pool or timeout settings onto a generated cluster that the typed DestinationRule API does not expose.
- **Remove a generated object** -- REMOVE a filter or listener Istio adds by default when you need a leaner proxy configuration.

## Works With

EnvoyFilter is the advanced member of the Istio family. It requires the Istio Base CRDs on Kubernetes and a running Istio control plane. Its attachment can target a **Gateway** or **Service Entry** via target references, or match application workloads by label. Prefer the first-class typed Istio APIs -- **Telemetry**, **Authorization Policy**, **Destination Rule**, **Request Authentication**, **Peer Authentication** -- wherever they cover your need, and reach for EnvoyFilter only for what they do not. To order this resource after the workloads, Gateway, or Service it patches within an infra chart, express the dependency through `metadata.relationships`.

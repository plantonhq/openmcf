# Service Entry on Kubernetes

Defines an Istio ServiceEntry: a namespaced resource that adds a host into Istio's internal service registry. Use it to make a service that platform service discovery does not know about -- an external API reached over the public internet, a SaaS endpoint, or a VM/legacy service -- routable, observable, and policy-addressable from inside the mesh by its hostname.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **A ServiceEntry** -- a namespaced Istio registry entry describing the service's hosts, ports, location (external or internal), how its endpoints are resolved, and (optionally) the static endpoints or in-mesh workloads that back it.
- **Kubernetes Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking.

## What It Controls

- **Identity** -- the hostnames the entry answers to, the ports it exposes (with protocols), and whether the service is external to the mesh or part of it.
- **Resolution and backing** -- how the proxy turns a host into endpoints: forward to the original destination, resolve by DNS, use static endpoint IPs, or select in-mesh workloads by label.
- **Reach and visibility** -- optional virtual IP/CIDR addresses for IP-based matching, and the namespaces the entry is exported (visible) to.
- **TLS verification** -- subject alternate names the proxy checks when originating TLS to the service.

## Important Behavior

A ServiceEntry teaches the mesh how to reach **out** to a service; it is not an ingress. Static endpoints and a workload selector are mutually exclusive -- at most one backs the service, and both omitted is valid (e.g. an external host resolved by DNS). Static endpoints require Static, DNS, or DNS Round Robin resolution (never None), and DNS Round Robin allows at most one endpoint. CIDR addresses are honored only with None or Static resolution.

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with kubeconfig credentials for the target cluster. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline kubeconfig authentication.

### Kubernetes Cluster

- **Istio installed** -- the Istio CRDs (Istio Base CRDs on Kubernetes) must be present and the Istio control plane (istiod) running. The entry is only honored where istiod is active.
- **Target namespace exists** -- the entry is created in a specific namespace; reference an existing one or create it first.

## Deploy

### Console

Open the deployment store, find **Service Entry on Kubernetes**, and click **Deploy**. The creation wizard walks you through the namespace, the service identity (hosts, location, ports), and how the service is resolved and backed, with guidance at each step.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesServiceEntry
metadata:
  name: external-https-api
  org: acme-corp
  env: prod
spec:
  namespace:
    value: prod-apps
  hosts:
    - api.stripe.com
  location: MESH_EXTERNAL
  resolution: DNS
  ports:
    - number: 443
      name: https
      protocol: TLS
```

```shell
planton apply -f service-entry.yaml
```

This registers `api.stripe.com` as an external HTTPS service that workloads in `prod-apps` can reach by name, with TLS routed by SNI and the host resolved via DNS.

## Key Configuration

- **Namespace** -- the namespace the entry is created in. It is fixed once created; the entry is registered relative to it.
- **Hosts** -- the hostnames the entry matches (Host/Authority header, TLS SNI, or the DNS name to resolve). A bare `*` is not allowed; partial wildcards like `*.example.com` are.
- **Location** -- **External** (a service outside the mesh, the default) or **Internal** (part of the mesh, for VMs/legacy services brought in).
- **Ports** -- the ports the service is reached on, each with a unique name and number and an optional protocol (HTTP, HTTPS, GRPC, HTTP2, MONGO, TCP, TLS).
- **Resolution & backing** -- the resolution mode (None, Static, DNS, DNS Round Robin) plus the backing: passthrough/DNS, **static endpoints**, or a **workload selector**. At most one backing is set.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Description |
|------------|-------------|
| `namespace` | The namespace the entry is created in. Reference an existing Namespace on Kubernetes or supply the name directly. |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources and operators can reference:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `service_entry_name` | Name of the created ServiceEntry (equals `metadata.name`) | Ordering resources that depend on the entry being in place |
| `namespace` | The namespace the entry was created in | Confirming where the entry applies |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

- **Reach an external HTTPS API** -- register a SaaS/partner endpoint so mesh workloads can call it by hostname, with TLS routed by SNI and the host resolved via DNS. The canonical starting point. Start from the **external-https-api** preset.
- **Bring static endpoints into the mesh** -- register a VM-hosted database or legacy service as an internal destination with explicit backing IPs and Static resolution. Start from the **static-mesh-internal-endpoints** preset.

## Works With

ServiceEntry is part of the Istio networking family. It requires the Istio Base CRDs on Kubernetes and a running Istio control plane, and it pairs naturally with a Destination Rule -- the ServiceEntry makes the host *reachable and addressable*, while a DestinationRule configures *how* the mesh talks to it (load balancing, outlier detection, TLS origination). To order the entry after the workloads it fronts within an infra chart, express the dependency through `metadata.relationships`.

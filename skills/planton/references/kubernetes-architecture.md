# Kubernetes Architecture Judgment

The catalog's Kubernetes kinds compose into a platform, and some roads through
them are paved — deeper modules, presets, first-class integration — while
others are technically present but shallow. Recommending well means knowing
which is which. This reference is judgment for what runs ON the cluster;
`aws-architecture.md` covers the cloud around it, `environments.md` covers how
many clusters and environments, and `kubernetes-on-cluster.md` covers the
wiring mechanics.

## The paved road: traffic, DNS, and TLS

When an app needs to be reachable on a hostname, this is the stack — end to
end, in dependency order:

1. **`KubernetesGatewayApiCrds`** — the Gateway API resource definitions.
2. **`KubernetesIstio`** — the platform's deepest-supported gateway stack.
   The kind bundles istio/base, istiod, AND the ingress gateway in one
   resource — do not also add `KubernetesIstioBaseCrds` (that kind exists
   only for advanced split installs).
3. **`KubernetesGatewayClass`** — registers Istio as the Gateway API
   controller.
4. **`KubernetesGateway`** — the listener (typically one, with HTTPS and the
   cert).
5. **One `KubernetesHttpRoute` per hostname/app** — "route my app on
   `app.example.com`" is exactly one of these.
6. **`KubernetesCertManager`** — TLS certificates, automatically.
7. **`KubernetesExternalDns`** — watches the routes/gateways and writes the
   DNS records into the Route 53 zone automatically; its EKS config takes the
   zone by reference (`valueFrom` the chart's `AwsRoute53Zone`, or the zone id
   of an existing zone).

**Hard rule: never hand-wire load-balancer IPs or Elastic IPs into DNS
records when external-dns can be part of the architecture.** Static-IP
annotation plumbing is exactly the manual-work class this platform exists to
delete. If you catch yourself proposing "create an EIP and inject it via
annotation," stop and add external-dns instead.

`KubernetesIngressNginx` exists in the catalog but is the unpaved
alternative — offer it only when the user explicitly wants nginx.

## Educate at moments of leverage

When a platform capability erases work the user was bracing for, say so in
one sentence at the moment it lands — "external-dns will create that DNS
record automatically whenever you add a route; you never touch Route 53" is
the canonical example. These moments are where developers discover what the
platform (and Kubernetes) can do, and they are the single best use of
education. Everything else about the machinery stays in reserve per
`deployment-model.md`.

## The two-chart pattern: shared infrastructure + environment chart

When an ask mixes platform components and app workloads, propose the split —
by name, up front:

- **The shared-infrastructure chart** (deployed once): VPC/network, the EKS
  cluster, Gateway API CRDs, Istio, GatewayClass, cert-manager, external-dns,
  and the Route 53 zone. This is Scenario 1 of `kubernetes-on-cluster.md` —
  it publishes the cluster's connection for everything that follows.
- **The environment chart** (deployed once per environment — dev, prod, …):
  the app's namespace, its `KubernetesGateway`/`KubernetesHttpRoute` with the
  hostname as a param, and a **placeholder `KubernetesDeployment`** whose
  image is a param. Scenario 2 — it consumes the shared cluster's connection;
  `values.env` differentiates the deployments so one chart serves every
  environment.

**The placeholder-first philosophy:** the goal is the fastest *validated
base*. A placeholder answering on the real hostname proves the entire road —
cluster, gateway, DNS, TLS — and the user's CI/CD then updates the image to
ship the real app. Say this arc out loud when proposing the split, so the
user knows what "done" unlocks and where CI/CD picks up.

**Application-first:** ask about the APP — what it is, what port it listens
on, what hostname it should answer on. The infrastructure exists in service
of it; a composition that never asked what the user is deploying has failed
regardless of how clean the charts are.

**Sequencing:** propose both charts up front, then compose and finish the
shared-infrastructure chart before the environment chart — producers before
consumers, each chart driven green before the next begins. In your
workspace, each chart is its own top-level subfolder (the identity check in
the skill), so both live side by side in one conversation; when the folder
you were given IS a single chart, finish it and have the user open or
create the second chart's folder — the grounding duty (`discovery.md`)
means the next conversation discovers the deployed cluster automatically.
Either way, one folder holds one chart: never mix two charts' files into
one chart folder.

## Cost shape on the cluster

The cluster itself is the big always-on charge (control plane + nodes + NAT —
see `cost-transparency.md`); everything in the paved road above is free
software running on nodes the user already pays for. One more reason the
shared-infrastructure/environment-chart split is cost-honest: environments
share the paid platform instead of multiplying it (`environments.md`).

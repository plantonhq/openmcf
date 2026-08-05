# Kubernetes Workloads and Their Cluster

Every resource whose kind starts with `Kubernetes` (`KubernetesIstio`,
`KubernetesCertManager`, `KubernetesDeployment`, `KubernetesGatewayApiCrds`,
…) deploys INTO a cluster, and it reaches that cluster through a **Kubernetes
provider connection** — a platform record carrying the cluster's endpoint and
credentials. This is the single most-missed wiring in composed charts: a
Kubernetes resource without a resolvable connection fails the moment its
pipeline node starts, with no cloud call ever made.

How a resource's connection resolves, in order:

1. Explicit spec value (`spec.provider_info.connection`) — rarely used in charts.
2. The `planton.dev/connection` **annotation** — the chart mechanism. Value is
   the org-scoped connection slug.
3. The environment's default Kubernetes connection.
4. The organization's default Kubernetes connection.

Fresh instances and new organizations usually have NO Kubernetes defaults, so
in practice: **every Kubernetes-kind resource in a chart must carry the
`planton.dev/connection` annotation.** Decide which of the two scenarios below
applies before writing any Kubernetes manifest.

## Scenario 1 — the chart creates the cluster too (one-run composition)

The platform automatically creates ("materializes") a Kubernetes provider
connection when a cluster resource (e.g. `AwsEksCluster`) finishes deploying —
before any dependent node is released. The chart's job is to make the
connection's name and the workloads' annotation agree, and the pattern that
makes disagreement impossible is ONE values param used on BOTH ends:

```yaml
# values.yaml
params:
  - name: connection_name
    description: Name of the Kubernetes connection published for this cluster
    value: my-cluster

# templates/cluster.yaml — the PRODUCER declares the connection's name
apiVersion: aws.planton.dev/v1
kind: AwsEksCluster
metadata:
  name: "{{ values.env }}-cluster"
  annotations:
    planton.dev/connection-name: "{{ values.env }}-{{ values.connection_name }}"
spec:
  …

# templates/kubernetes/istio.yaml — every CONSUMER references the same expression
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesIstio
metadata:
  name: "{{ values.env }}-istio"
  annotations:
    planton.dev/connection: "{{ values.env }}-{{ values.connection_name }}"
  relationships:
    - kind: AwsEksCluster
      name: "{{ values.env }}-cluster"
      type: runs_on
spec:
  …
```

Rules of the pattern:

- `planton.dev/connection-name` (on the cluster) sets the published
  connection's name **literally** — nothing is prepended. Always embed
  `{{ values.env }}` in the value: connections are org-scoped, and two
  environments deploying the same chart must not collide. A collision with a
  connection the cluster did not create fails the deploy loudly.
- `planton.dev/connection` (on each workload) must render to the **identical
  string**. Using the same values expression on both ends guarantees it.
- Without the `connection-name` annotation, the platform publishes the
  connection at `<env>-<cluster resource metadata.name>` — you will see this
  legacy formula in older charts, where the workload annotations re-derive it
  (`"{{ values.env }}-{{ values.cluster_name }}"` against a cluster named
  without an env prefix). Prefer the explicit `connection-name` pattern in new
  charts; it composes cleanly with env-prefixed resource names.

## Scenario 2 — the cluster already exists

When the user deploys onto a cluster that is NOT in this chart, no cluster
resource and no `connection-name` annotation exist here. The workloads
annotate the **existing** connection's slug, taken as a param so the user can
point the chart at any cluster:

```yaml
params:
  - name: cluster_connection
    description: Slug of the Kubernetes provider connection for the target cluster
    value: ""

# every Kubernetes workload:
metadata:
  annotations:
    planton.dev/connection: "{{ values.cluster_connection }}"
```

Ground the actual slug instead of guessing: if the cluster was deployed
through Planton (e.g. by another chart), its connection exists at the name
that chart declared — look it up (`planton` CLI: list the org's Kubernetes
provider connections) and propose it as the param default. In this scenario
the workloads need no `runs_on` relationship to a cluster — the cluster is
not in this chart's graph.

## Ordering — a separate concern from credentials

`metadata.relationships` entries create REAL dependency edges in the deploy
graph: a resource with a `runs_on` or `depends_on` relationship does not start
until its target succeeds. The connection annotation carries credentials only;
it creates no edge. Both are required in Scenario 1:

- **Every Kubernetes workload gets a `runs_on` relationship** to the in-chart
  cluster resource (`kind` + the cluster's exact `metadata.name` expression).
- **When the chart has a node group** (e.g. `AwsEksNodeGroup`), anything that
  schedules pods (operators, deployments, Helm-backed addons) should `runs_on`
  the NODE GROUP instead — a cluster with zero nodes accepts the install but
  pods never start and the deploy times out. CRD-only kinds (e.g.
  `KubernetesGatewayApiCrds`) may target the cluster directly.
- **Chain the Kubernetes resources among themselves with `depends_on`** in
  their real order: CRDs → the operator that serves them → the instances that
  use them. Example: `KubernetesGatewayApiCrds` ← `KubernetesIstio` ←
  `KubernetesGatewayClass`.

## Never render a connection manifest inside a chart

Charts render cloud resources only. A `KubernetesProviderConnection` document
in templates fails the build ("UNSUPPORTED CLOUD RESOURCE KIND") — connections
are org-scoped, authorization-bearing records the platform materializes or
users create; they never belong in templates.

## Self-check before finishing any chart with Kubernetes kinds

1. Every `Kubernetes*` resource carries `planton.dev/connection`.
2. Scenario 1: the cluster carries `planton.dev/connection-name`, and every
   workload's annotation uses the same values expression.
3. Every workload has a `runs_on` relationship (node group when pods are
   scheduled and one exists; cluster otherwise) — Scenario 1 only.
4. Kubernetes resources are `depends_on`-chained in dependency order
   (CRDs → operators → instances).
5. `planton chart build` is green — but note the build cannot verify the
   annotation matches a connection that will exist at deploy time; the
   self-check above is what protects the deploy.

# KubernetesClusterAutoscaler: Research and Design

## Introduction

The Kubernetes Cluster Autoscaler is the node-group scaler: when pods
are unschedulable it raises the desired size of a matching pre-defined
group (an EC2 Auto Scaling group, an Azure VMSS, a GCE managed instance
group, a Cluster API MachineDeployment), and it scales groups back down
when nodes sit underutilized. This component installs it from the
official Helm chart (`cluster-autoscaler` at
`https://kubernetes.github.io/autoscaler`; pinned default chart 9.59.0,
which ships autoscaler 1.35.0 — chart and app versions move SEPARATELY,
and the chart pin governs). Upstream guidance is to keep the
autoscaler's minor aligned with the cluster's Kubernetes minor; older
clusters override the image tag via `helm_values`.

## Where This Component Belongs — and Where It Does Not

The autoscaler's model is pre-defined groups with size bounds. That
matches EKS-with-ASGs, Cluster API and self-managed clusters, and
providers without a managed autoscaler. Two deliberate boundaries:

- **GKE and AKS ship a MANAGED autoscaler** configured as a toggle on
  the node pool itself, through the cluster kinds. Deploying this
  component there is the exception, not the rule — the `gce` arm exists
  for SELF-MANAGED clusters on GCE (managed instance groups by name
  prefix), and the `azure` arm for AKS clusters that opt out of the
  managed autoscaler and for self-managed clusters on VMSS. Deploy this
  only where self-managed autoscaling is the real posture.
- **Karpenter is the on-demand alternative**: instead of resizing
  groups, it launches right-sized machines per pending pod. For AWS
  clusters wanting that model, KubernetesKarpenter is the modern choice;
  this component remains the fit where groups are the deliberate
  capacity model.

One installation per cluster: the autoscaler leader-elects and owns the
cluster-wide scaling decision, so the release name is fixed to
`cluster-autoscaler`. The fixed release name plus the chart's naming
scheme (`<release>-<cloudProvider>-<chartName>`) makes the
service-account name deterministic per arm — e.g.
`cluster-autoscaler-aws-cluster-autoscaler` — which is why it is a stack
output: it is the subject cloud-side keyless bindings are written
against.

## The Provider Oneof

The autoscaler binary ships dozens of cloud providers (the upstream
source tree's `cloudprovider/` directory spans alicloud through utho);
one installation talks to exactly one. The spec models the arms the
chart itself gives first-class values for, as a REQUIRED oneof:

- **`aws`** — region (required by the chart), tag-based
  `auto_discovery` (the recommended mode: every ASG carrying
  `k8s.io/cluster-autoscaler/enabled` +
  `k8s.io/cluster-autoscaler/<cluster_name>` enrolls by tagging alone)
  XOR static `node_groups` (each rendering one `--nodes=min:max:name`
  flag), and IRSA XOR static access keys.
- **`azure`** — subscription and resource group (for AKS the NODE
  resource group, the MC_* group holding the VMSS instances), tag-based
  discovery by `cluster_name` XOR static `node_groups`, and an
  `identity` message enforcing exactly one credential posture: workload
  identity, managed identity (optionally naming a user-assigned
  identity), or a service principal.
- **`gce`** — managed instance groups by NAME PREFIX with size bounds
  (the chart's GCE contract — no tagging), plus an optional Workload
  Identity service-account email.
- **`cluster_api`** — the self-managed / multi-cloud arm: `mode` picks
  where the workload cluster and the management objects live
  (`incluster-incluster` chart default through `single-kubeconfig`);
  kubeconfig modes require `kubeconfig_secret` — the NAME of an
  in-cluster Secret, never the kubeconfig itself (a CEL rule enforces
  the pairing). `namespace` narrows discovery, and
  `namespace_scoped_rbac` switches the chart to namespaced
  Role/RoleBinding — the least-privilege posture the chart documents as
  most useful for Cluster API.
- **`civo`** — managed node pools on Civo: cluster ID, region, and an
  API key stored as a managed secret.
- **`kwok`** — the simulation provider: nodes are FAKE, created by the
  KWOK controller (which must be installed on the cluster). For testing
  scaling policies and evaluating the autoscaler without a cloud
  account; never a production arm. `config_map_name` (chart default
  `kwok-provider-config`) names the provider ConfigMap of simulated node
  templates.

A chart-mechanics note the modules encode: the chart's Deployment
template only renders when a discovery value or a static group list is
set. The GCE contract (`autoscalingGroupsnamePrefix` — the chart key
really does carry a lowercase "n") and the kwok arm do not satisfy that
gate on their own, so the modules render a benign
`autoDiscovery.clusterName` (the resource's metadata name) for those two
arms — without it the release would "succeed" while installing no
autoscaler pod. The civo arm's node pools ride `helm_values`
(`autoscalingGroups`), which satisfies the gate the same way.

## The Typed Scaling Block and the extra_args Contract

The autoscaler is configured almost entirely through CLI flags — over a
hundred of them. The spec splits the surface deliberately:

- **`scaling`** types the flags every real installation ends up tuning:
  `expander` (comma-separated ordered list from random/most-pods/
  least-waste/price/priority/grpc; the priority expander additionally
  reads the cluster-autoscaler-priority-expander ConfigMap, shipped via
  `helm_values` expanderPriorities), `balance_similar_node_groups` (the
  multi-AZ pattern where one group per zone is the norm),
  `scan_interval` (upstream default 10s), `max_node_provision_time`
  (15m), the two skip-nodes switches (both upstream-default TRUE, so
  they are presence-aware optionals — an explicit false must render),
  and the `scale_down` sub-block (master switch, utilization threshold
  "0.5", unneeded time 10m, and the add/delete/failure cool-downs).
- **`extra_args`** is the long tail — flag-name → value pairs rendered
  into the chart's extraArgs without the leading `--`. This IS the
  chart's own contract for arbitrary flags; names are validated for
  shape only, and unknown flags fail at pod start. Precedence is
  defined: the typed block renders first, `extra_args` merges OVER it,
  so user entries win on key collision; the chart's own extraArgs
  defaults (logtostderr/stderrthreshold/v) survive per-key.

## Deployment Posture

`deployment.replicas` (chart default 1) leader-elect — extras are warm
standbys, not horizontal capacity. `priority_class_name` chart-defaults
to `system-cluster-critical`: the component that adds capacity must not
be evicted for lack of it. The node selector deserves its own emphasis:
the autoscaler should not run on a node it may delete — a management
node group is the standard home. `prometheus.service_monitor` requires
the Prometheus operator CRDs (the release fails without them), and
`service_monitor_selector_release` replaces the chart's default
`release: prometheus-operator` selector label to match the actual
Prometheus installation.

## Typed Surface vs Escape Hatch

The typed spec covers namespace and lifecycle, chart version, the six
provider arms with their credential postures, the scaling block, the
extraArgs contract, deployment sizing/scheduling, and own telemetry.

Deliberately unmodeled as typed fields (all reachable via
`helm_values`):

- **The priority-expander ConfigMap** (`expanderPriorities`) — data for
  one expander choice, not installation shape
- **Civo node-pool lists** (`autoscalingGroups` for the civo arm) — the
  chart README's contract for civo; a niche arm's group list
- **Image overrides** — required when pinning the autoscaler minor to an
  older cluster minor
- **PodDisruptionBudget tuning, VPA, extra volumes/mounts,
  secretKeyRefNameOverride** — the chart's operational long tail, each
  with a correct default
- **Azure workload-identity pod labels** — the chart sets the ARM
  workload-identity extension env; clusters relying on the
  azure-workload-identity webhook add `podLabels` via `helm_values`

## Install Semantics

Both engines install a REAL Helm release with the chart pinned, values
rendered from the typed spec, and the `helm_values` document merged last
with Helm `-f` semantics (maps deep-merge with the later document
winning, lists replace). Secret-bearing arms (AWS access keys, Azure
service principal, Civo API key) render into the chart's own credential
Secret and reach the pod via secretKeyRef — the secret value never lands
in the pod spec. The module (not Helm) owns namespace creation via
`create_namespace`; the flag stays false for the conventional
`kube-system` home.

## Outputs

`namespace`, `release_name` (fixed `cluster-autoscaler`), and
`service_account_name` (the chart-derived
`cluster-autoscaler-<provider>-cluster-autoscaler` — the subject IRSA
trust policies, GCP WI bindings, and Entra federated credentials are
written against).

## E2E

The behavioral facts are properties of the platform, not of any one test
run:

- The kwok arm is the deterministic behavioral proof: with the KWOK
  controller installed, an unschedulable pod drives the autoscaler to
  "launch" fake nodes from the provider ConfigMap — scaling logic
  end-to-end with no cloud account.
- Real-cloud arms prove out against pre-existing groups: the autoscaler
  raises an ASG/VMSS desired size for pending pods and shrinks it after
  `unneeded_time` under the utilization threshold.
- The Deployment render gate means a mis-declared provider arm can
  produce a release with no pod — the modules' gate values exist
  precisely to keep that failure impossible for the typed arms.
- The ServiceMonitor arm fails the release on clusters without the
  Prometheus operator CRDs, by design.
- Unknown `extra_args` flags fail at pod start, not at install — the
  release reports ready only if the pod survives its flag parsing.

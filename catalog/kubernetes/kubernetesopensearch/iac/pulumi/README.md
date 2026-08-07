# KubernetesOpenSearch Pulumi Module

Deploys one operator-managed OpenSearch cluster: the optional namespace
and the `opensearch.opster.io/v1` OpenSearchCluster resource. The custom
resource renders through typed crd2pulumi SDK bindings pinned to the
OpenSearch Kubernetes Operator's CRD (v2.8.0) — field or structure drift
against the pinned CRD fails at COMPILE time, not at apply time.

Prerequisites at deploy time: the OpenSearch Kubernetes Operator
(`KubernetesOpenSearchOperator`) on the cluster; the Prometheus
Operator's ServiceMonitor CRD when `monitoring` is enabled.

## What the Module Creates

1. **Namespace** (optional) — created with the standard governance
   labels when `create_namespace` is true
2. **OpenSearchCluster** — the ONLY custom resource. Everything else is
   operator-created from it: one StatefulSet per node pool
   (`<name>-<pool>`), the main Service `<name>` and discovery Service
   `<name>-discovery`, generated TLS Secrets, the security-plugin admin
   bootstrap with its `<name>-admin-password` credentials Secret, and
   the optional Dashboards deployment + Service `<name>-dashboards`.

There are NO ingress resources in the module by design — exposure
composes from first-class kinds (KubernetesIngress, Gateway API kinds)
referencing the exported service handles.

## CR Mapping (spec → OpenSearchCluster)

- `metadata.name` = resource name; the namespace resolves from the
  spec's value-or-ref.
- **`general`** — `serviceName` is ALWAYS `metadata.name`
  (module-owned; the operator names the main Service after it) and
  `vendor` is always `"opensearch"`. `version`, `httpPort` (default
  9200), `additionalConfig`, service `annotations`,
  `setVMMaxMapCount` (spec default true; rendered only when true —
  the CRD default is false), `drainDataNodes`, `pluginsList`,
  `keystore` (secret name + optional key mappings),
  `snapshotRepositories`, `monitoring` (when enabled), and
  `additionalVolumes` map 1:1. The CRD's ImageSpec takes ONE image
  string, so the shared ContainerImage folds into `repo:tag`;
  `imagePullSecrets` come from the spec's own list.
- **`bootstrap`** — resources, jvm, additionalConfig (rendered only
  when declared).
- **`nodePools`** — one entry per `node_pools` pool: `component` is the
  pool name, plus replicas, roles, resources, jvm, diskSize,
  persistence (PVC arm with the CRD's `storageClass` key and pinned
  `accessModes: [ReadWriteOnce]`, or emptyDir with optional
  sizeLimit), per-pool additionalConfig, nodeSelector, tolerations,
  and pdb — PDB bounds follow intstr semantics (integer strings render
  as YAML numbers, percentages as strings, identical on both engines).
- **`security`** — only declared blocks render: `tls.transport`
  (generate/perNode spec-default true and rendered only when true —
  the CRD default is false — plus secret/caSecret/nodesDn/adminDn),
  `tls.http` (generate + secret), and `config`
  (securityConfigSecret/adminSecret/adminCredentialsSecret).
- **`dashboards`** — rendered only when enabled: replicas (default 1),
  version (defaults to the CLUSTER version — module-owned; the CRD
  requires it and Dashboards refuses mismatched clusters), resources,
  tls, basePath, additionalConfig, opensearchCredentialsSecret,
  service (type + loadBalancerSourceRanges), pluginsList.

## Operator Contracts the Outputs Encode

- **`http_endpoint` is always https** — the operator's HTTP layer
  serves TLS in EVERY posture: with `security.tls.http` it uses
  operator-generated (or provided) certificates; with `spec.security`
  entirely absent the TLS reconciler generates nothing
  (`pkg/reconcilers/tls.go`) and the opensearchproject image's demo
  security configuration serves TLS instead — the operator itself
  always connects over https (`pkg/builders/cluster.go` URLForCluster,
  and the node readiness probe curls `https://localhost`).
- **`admin_credentials_secret_name`** is the operator-generated
  `<name>-admin-password` Secret (fields `username`/`password`,
  `pkg/builders/cluster.go` PasswordSecret) — exported EMPTY when
  `security.config` is set: a custom security config replaces the
  operator bootstrap and the user's `admin_credentials_secret` is
  authoritative.
- **Dashboards** — Service `<name>-dashboards` on the fixed port 5601;
  the endpoint scheme flips to https when `dashboards.tls.enable`.

## Usage

```shell
planton pulumi up --manifest e2e/manifest.yaml --module-dir <path-to-this-module>
```

## Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the cluster runs in |
| `cluster_name` | Name of the OpenSearchCluster resource (equals `metadata.name`) |
| `service_name` | The cluster's main Service (= `metadata.name`) |
| `http_endpoint` | In-cluster HTTP API endpoint (`https://<name>.<namespace>.svc.cluster.local:<http_port>`) |
| `admin_credentials_secret_name` | `<name>-admin-password` — empty when a custom security config is declared |
| `dashboards_service_name` | `<name>-dashboards` — empty when dashboards are not enabled |
| `dashboards_endpoint` | In-cluster Dashboards endpoint on port 5601 — empty when dashboards are not enabled |
| `port_forward_command` | Port-forward command for workstation access |

## Module Structure

- `main.go`: entrypoint that calls the module
- `module/main.go`: namespace → OpenSearchCluster → output exports
- `module/locals.go`: the naming/endpoint contract (service names, the
  https endpoint rationale, the admin-secret handle) — kept in lockstep
  with the Terraform module's `locals.tf`
- `module/cluster.go`: the OpenSearchCluster resource (general,
  bootstrap, node pools, security, dashboards) with the shared
  ContainerResources/WorkloadToleration translations and the intstr
  PDB-bound rendering
- `module/vars.go`: the module-owned constants (vendor `"opensearch"`,
  the fixed Dashboards port 5601)

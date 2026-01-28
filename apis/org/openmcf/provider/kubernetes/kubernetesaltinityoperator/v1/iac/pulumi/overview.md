# Pulumi Module Overview: Kubernetes Altinity Operator

## Architecture

The Pulumi module for the Kubernetes Altinity Operator follows OpenMCF's standard Kubernetes add-on deployment pattern, with specific design decisions tailored to operator deployment requirements.

## High-Level Architecture

```
KubernetesAltinityOperatorStackInput
    ↓
  Locals (Data Transformations)
    ↓
  Kubernetes Provider Setup
    ↓
  Resource Creation
    ├── Namespace (dedicated for operator)
    └── Helm Release (operator deployment)
         ├── CRDs (ClickHouseInstallation, ClickHouseKeeperInstallation)
         ├── Operator Deployment
         └── ServiceAccount + RBAC
    ↓
  Stack Outputs
```

## Key Design Decisions

### 1. Helm-Based Deployment

**Decision:** Deploy the operator via the official Altinity Helm chart rather than raw Kubernetes manifests.

**Rationale:**
- The Altinity operator has a complex installation with multiple components (CRDs, Deployment, ServiceAccount, RBAC)
- The Helm chart is the officially maintained installation method with guaranteed correctness
- Helm provides atomic rollback capabilities, critical for operator upgrades
- The chart encapsulates best practices from Altinity's production deployments at scale

**Tradeoff:** Adds a dependency on Helm, but this is acceptable as Helm is ubiquitous in Kubernetes deployments.

### 2. CRD Management Strategy

**Decision:** Enable CRD installation via the Helm chart (`operator.createCRD: true`).

**Rationale:**
- Ensures CRDs and operator version are always in sync
- Simplifies initial deployment (one-step installation)
- Helm 3's improved CRD handling makes this safe for most scenarios

**Alternative Considered:** Two-step installation (CRDs first via kubectl, then operator via Helm). This is the production best practice for environments with strict change management, but adds complexity for most users.

**Future Enhancement:** Could add a `spec.install_crds_separately` flag for advanced users who want two-step installation.

### 3. Cluster-Wide vs. Namespace-Scoped Operation

**Decision:** Configure the operator to watch all namespaces (`watch.namespaces: [".*"]`).

**Rationale:**
- Enables a single operator instance to manage ClickHouse deployments across all namespaces
- Matches the common use case: centralized platform team managing operator, application teams deploying ClickHouse instances
- Reduces resource overhead (no need for operator-per-namespace)

**Alternative Available:** For multi-tenant environments requiring strict namespace isolation, users can override the Helm values to set `rbac.namespaceScoped: true` and `watch.namespaces: []` to watch only the installation namespace.

### 4. Resource Allocation Pattern

**Decision:** Make container resources configurable via the spec, with sensible defaults.

**Rationale:**
- The operator itself is lightweight (manages cluster metadata, not data plane), but resource needs scale with the number of managed ClickHouse instances
- Defaults (100m CPU, 256Mi memory) are appropriate for small-to-medium deployments (1-10 ClickHouse clusters)
- Production deployments managing 50+ clusters may need higher limits

**Default Values:**
```
Requests: 100m CPU, 256Mi memory
Limits: 1000m CPU, 1Gi memory
```

### 5. Namespace Isolation

**Decision:** Create a dedicated namespace for the operator (default: `kubernetes-altinity-operator`).

**Rationale:**
- Separates operator infrastructure from application workloads
- Enables RBAC policies specific to operator management
- Makes it easy to identify and monitor operator resources

**Note:** This namespace is for the **operator itself**, not for ClickHouse deployments. ClickHouse instances are typically deployed in separate application namespaces.

## Module Structure

### `vars.go`
Defines constants for:
- Default namespace
- Helm chart name (`altinity-clickhouse-operator`)
- Helm chart repository URL
- Pinned chart version (ensures reproducible deployments)

### `locals.go`
Data transformations and computed values:
- **Namespace resolution:** Uses spec value or falls back to default
- **Helm values preparation:** Constructs the Helm values map with:
  - CRD installation flag
  - Resource limits from spec
  - Namespace watch configuration

**Why separate locals?** Separating data transformations from resource creation improves testability and makes the code easier to reason about.

### `kubernetes_altinity_operator.go`
Main resource creation logic:
1. **Provider setup:** Configures Kubernetes provider from stack input credentials
2. **Namespace creation:** Creates dedicated namespace for operator
3. **Helm release:** Deploys operator with computed Helm values

**Key Pulumi patterns:**
- Uses `pulumi.Provider(kubeProvider)` to ensure all resources use the target cluster
- Sets `pulumi.Parent(ns)` on Helm release for logical resource hierarchy
- Enables `Atomic: true` for rollback on failure
- Configures `WaitForJobs: true` to block until operator is ready

### `outputs.go`
Exports stack outputs matching `KubernetesAltinityOperatorStackOutputs`:
- `namespace`: The resolved namespace where operator is installed

**Why minimal outputs?** The operator is infrastructure—consumers typically only need to know where it's installed. The operator's own status is visible via its CRDs.

## Resource Relationships

```
Namespace (kubernetes-altinity-operator)
  ↓ (parent relationship)
Helm Release (altinity-clickhouse-operator)
  ↓ (creates)
  ├── CRDs
  │   ├── ClickHouseInstallation
  │   └── ClickHouseKeeperInstallation
  ├── Deployment (operator pods)
  ├── ServiceAccount
  ├── ClusterRole (if cluster-wide)
  └── ClusterRoleBinding (if cluster-wide)
```

**Dependency Flow:**
1. Namespace must exist before Helm release
2. Helm release atomically creates all child resources
3. Operator deployment waits for CRDs to be registered
4. Operator enters ready state, begins watching for ClickHouseInstallation CRs

## Data Flow

### Input → Locals → Resources

1. **Input:** `KubernetesAltinityOperatorStackInput`
   - Target cluster credentials
   - Namespace (optional, empty string = use default)
   - Container resources

2. **Locals Transformation:**
   - Resolve namespace: `spec.namespace || "kubernetes-altinity-operator"`
   - Construct Helm values map with resource limits and watch config

3. **Resource Creation:**
   - Create namespace using resolved value
   - Deploy Helm chart with transformed values
   - Wait for operator to reach ready state

4. **Output:**
   - Export namespace for downstream use

## Comparison to Manual Installation

### Manual Approach (via kubectl + Helm)
```bash
# Step 1: Create namespace
kubectl create namespace clickhouse-operator

# Step 2: Install CRDs (optional, for strict change control)
kubectl apply -f https://raw.githubusercontent.com/Altinity/clickhouse-operator/master/deploy/operator/clickhouse-operator-install-bundle.yaml

# Step 3: Install operator via Helm
helm repo add altinity https://docs.altinity.com/clickhouse-operator/
helm install clickhouse-operator altinity/altinity-clickhouse-operator \
  --namespace clickhouse-operator \
  --set operator.createCRD=true \
  --set operator.resources.requests.cpu=100m \
  --set operator.resources.requests.memory=256Mi
```

### Pulumi Module Approach
```yaml
# manifest.yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesAltinityOperator
metadata:
  name: clickhouse-operator
spec:
  target_cluster:
    credential_id: my-k8s-cluster
  namespace: kubernetes-altinity-operator  # optional
  container:
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 1000m
        memory: 1Gi
```

```bash
planton pulumi up --manifest manifest.yaml
```

**Advantages of Pulumi approach:**
1. **Declarative:** Entire config in manifest, no imperative commands
2. **Idempotent:** Safe to rerun, detects drift
3. **Integrated:** Works with OpenMCF's provider credential system
4. **Auditable:** Changes tracked in manifest version history
5. **Consistent:** Same deployment method as other Kubernetes add-ons

## Operational Considerations

### Upgrading the Operator

The operator version is pinned in `vars.HelmChartVersion`. To upgrade:

1. Update `HelmChartVersion` in `vars.go`
2. Test in non-production environment
3. Run `planton pulumi up`

Pulumi will:
- Detect the chart version change
- Perform Helm upgrade with `Atomic: true` (rollback on failure)
- Wait for new operator pods to become ready

**Important:** Always check Altinity's release notes for breaking changes, especially CRD schema changes.

### Monitoring and Observability

The operator exposes Prometheus metrics on port 8888. To enable monitoring:

1. Install Prometheus Operator or equivalent
2. Create a `ServiceMonitor` to scrape operator metrics
3. Import Altinity's Grafana dashboard (ID: 11159)

**Key Metrics:**
- `clickhouse_operator_reconciles_total`: Number of reconciliations
- `clickhouse_operator_reconcile_errors_total`: Failed reconciliations
- `clickhouse_operator_managed_chi_total`: Number of managed ClickHouse instances

### Troubleshooting

**Operator pod not starting:**
```bash
kubectl describe pod -n kubernetes-altinity-operator -l app=clickhouse-operator
```

Check for:
- CRD registration failures (in older Kubernetes versions)
- RBAC permission issues (if namespace-scoped mode is enabled)
- Resource limits too low for number of managed instances

**CRDs not found:**
The chart should install CRDs automatically. If missing:
```bash
kubectl get crd | grep clickhouse
```

Expected CRDs:
- `clickhouseinstallations.clickhouse.altinity.com`
- `clickhousekeeperinstallations.clickhouse.altinity.com`

**Operator managing wrong namespaces:**
Check Helm values:
```bash
helm get values clickhouse-operator -n kubernetes-altinity-operator
```

Verify `configs.files.config.yaml.watch.namespaces` is set correctly.

## Future Enhancements

### Potential Module Improvements

1. **Custom Helm Values Override:**
   Add `spec.helm_values_override` to allow advanced users to pass arbitrary Helm values without forking the module.

2. **Two-Step CRD Installation:**
   Add `spec.install_crds_separately: bool` to support strict change management environments.

3. **ServiceMonitor Integration:**
   Auto-create ServiceMonitor if Prometheus Operator is detected in cluster.

4. **Version Pinning Strategy:**
   Consider allowing `spec.operator_version` to override the default pinned version for testing pre-release versions.

5. **Multi-Operator Deployments:**
   Support deploying multiple namespace-scoped operators for extreme multi-tenancy scenarios.

## References

- [Altinity Operator GitHub](https://github.com/Altinity/clickhouse-operator)
- [Altinity Helm Chart](https://github.com/Altinity/clickhouse-operator/tree/master/deploy/helm)
- [Operator Documentation](https://github.com/Altinity/clickhouse-operator/tree/master/docs)
- [ClickHouse Operator Manifests](https://github.com/Altinity/clickhouse-operator/tree/master/docs/chi-examples)


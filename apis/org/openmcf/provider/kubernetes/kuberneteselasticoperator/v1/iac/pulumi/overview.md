# Kubernetes Elastic Operator - Pulumi Module Overview

## Purpose

This Pulumi module automates the deployment of the Elastic Cloud on Kubernetes (ECK) operator to Kubernetes clusters. The ECK operator provides automated lifecycle management for the Elastic Stack (Elasticsearch, Kibana, APM Server, Enterprise Search, Beats, Agent, Logstash) on Kubernetes.

## Architecture

### Component Hierarchy

```
KubernetesElasticOperator (API Resource)
    │
    ├── Kubernetes Namespace (elastic-system)
    │   └── Labels: resource, resource_id, resource_kind, org, env
    │
    └── Helm Release (eck-operator)
        ├── ECK Operator Deployment
        │   ├── Operator Pod
        │   │   ├── Container: eck-operator
        │   │   ├── Resources: CPU/Memory limits and requests
        │   │   └── Service Account: elastic-operator
        │   └── Reconciliation Controller
        │
        ├── Custom Resource Definitions (CRDs)
        │   ├── elasticsearch.k8s.elastic.co
        │   ├── kibana.k8s.elastic.co
        │   ├── apmserver.k8s.elastic.co
        │   ├── enterprisesearch.k8s.elastic.co
        │   ├── beat.k8s.elastic.co
        │   ├── agent.k8s.elastic.co
        │   └── logstash.k8s.elastic.co
        │
        └── RBAC Resources
            ├── ClusterRole: elastic-operator
            ├── ClusterRoleBinding: elastic-operator
            └── ServiceAccount: elastic-operator
```

### Data Flow

1. **Input**: KubernetesElasticOperator API resource (protobuf)
2. **Transformation**: Pulumi module converts to Kubernetes/Helm resources
3. **Deployment**: Helm chart deploys ECK operator to cluster
4. **Operation**: ECK operator watches for Elastic Stack CRDs
5. **Management**: Operator creates and manages Elastic Stack components

## Design Decisions

### Why Helm?

The module uses Helm for ECK deployment because:

1. **Official Support**: Elastic provides and maintains the official Helm chart
2. **Versioning**: Helm simplifies operator version management and upgrades
3. **Configuration**: Helm values provide structured configuration without YAML templating
4. **Rollback**: Helm's rollback capability provides safety net for failed upgrades
5. **Community Standard**: Helm is the de facto package manager for Kubernetes

### Label Propagation

The module configures ECK to inherit Planton labels:

```go
"configKubernetes": pulumi.Map{
    "inherited_labels": pulumi.ToStringArray([]string{
        "resource", "organization", "environment", 
        "resource_kind", "resource_id",
    }),
}
```

**Benefits**:
- All Elastic Stack resources (Elasticsearch, Kibana, etc.) automatically receive these labels
- Enables filtering and grouping across the entire Elastic Stack deployment
- Supports cost allocation and resource tracking
- Maintains organizational metadata throughout the stack

### Resource Management

The module allows configuring operator pod resources:

```go
values["resources"] = pulumi.Map{
    "limits": pulumi.StringMap{
        "cpu":    pulumi.String(cr.GetLimits().Cpu),
        "memory": pulumi.String(cr.GetLimits().Memory),
    },
    "requests": pulumi.StringMap{
        "cpu":    pulumi.String(cr.GetRequests().Cpu),
        "memory": pulumi.String(cr.GetRequests().Memory),
    },
}
```

**Rationale**:
- Prevents operator from consuming excessive cluster resources
- Ensures operator has guaranteed minimum resources for reliable operation
- Allows scaling operator resources based on number of managed clusters

### Namespace Isolation

The operator deploys to a dedicated `elastic-system` namespace:

**Advantages**:
- Clear separation from user workloads
- Simplified RBAC policy management
- Easy identification of operator resources
- Follows Kubernetes best practices for addon operators

### Default Values

The module uses production-ready defaults from `vars.go`:

```go
Namespace:        "elastic-system",
HelmChartName:    "eck-operator",
HelmChartRepo:    "https://helm.elastic.co",
HelmChartVersion: "2.14.0",
```

**Philosophy**:
- Sensible defaults minimize configuration burden
- Version pinning ensures reproducible deployments
- Central constants file (`vars.go`) simplifies version upgrades

## Module Workflow

### 1. Input Processing

```go
func Resources(ctx *pulumi.Context, stackInput *KubernetesElasticOperatorStackInput) error {
    locals := initializeLocals(ctx, stackInput)
    // ... process input
}
```

**Actions**:
- Parse `KubernetesElasticOperatorStackInput` protobuf
- Extract metadata, spec, and target cluster configuration
- Initialize local variables and labels

### 2. Kubernetes Provider Setup

```go
k8sProvider, err := kubernetesprovider.GetWithKubernetesClusterCredential(
    ctx, stackInput.Target.Spec.TargetCluster, "kubernetes")
```

**Actions**:
- Retrieve Kubernetes cluster credentials
- Configure Pulumi Kubernetes provider
- Establish connection to target cluster

### 3. Namespace Creation

```go
ns, err := corev1.NewNamespace(ctx, vars.Namespace, &corev1.NamespaceArgs{
    Metadata: metav1.ObjectMetaPtrInput(&metav1.ObjectMetaArgs{
        Name:   pulumi.String(vars.Namespace),
        Labels: pulumi.ToStringMap(locals.KubeLabels),
    }),
}, pulumi.Provider(k8sProvider))
```

**Actions**:
- Create `elastic-system` namespace
- Apply Planton labels for tracking and organization
- Set namespace as parent for subsequent resources

### 4. Helm Values Construction

```go
values := pulumi.Map{
    "configKubernetes": pulumi.Map{
        "inherited_labels": pulumi.ToStringArray([...]string{...}),
    },
}
// Add resources if specified
if cr := locals.KubernetesElasticOperator.Spec.GetContainer().GetResources(); cr != nil {
    values["resources"] = constructResourceMap(cr)
}
```

**Actions**:
- Build Helm values map from spec
- Configure label inheritance
- Set operator pod resource limits/requests

### 5. Helm Release Deployment

```go
helm.NewRelease(ctx, "kubernetes-elastic-operator", &helm.ReleaseArgs{
    Name:            pulumi.String(vars.HelmChartName),
    Namespace:       ns.Metadata.Name(),
    Chart:           pulumi.String(vars.HelmChartName),
    Version:         pulumi.String(vars.HelmChartVersion),
    RepositoryOpts:  helm.RepositoryOptsArgs{Repo: pulumi.String(vars.HelmChartRepo)},
    Values:          values,
    // ... additional options
}, pulumi.Parent(ns))
```

**Actions**:
- Fetch ECK operator Helm chart from Elastic repository
- Apply Helm values for configuration
- Deploy operator with specified version
- Wait for deployment to complete

### 6. Output Export

```go
ctx.Export(OpNamespace, ns.Metadata.Name())
```

**Actions**:
- Export namespace name as stack output
- Make namespace available to other Pulumi stacks or automation

## Error Handling

The module implements robust error handling:

```go
if err := validator.Validate(stackInput); err != nil {
    return errors.Wrap(err, "failed to validate stack-input")
}

if _, err := helm.NewRelease(...); err != nil {
    return errors.Wrap(err, "install helm chart")
}
```

**Strategy**:
- Validate inputs before processing
- Wrap errors with context for debugging
- Return errors to Pulumi for proper reporting
- No silent failures - all errors propagate

## Key Features

### Idempotency

The module is fully idempotent:
- Re-running `pulumi up` with same inputs makes no changes
- Helm manages state and detects configuration drift
- Namespace and RBAC resources are declarative

### Upgrade Support

Operator upgrades are handled via:
1. Update `HelmChartVersion` in `vars.go`
2. Run `pulumi up`
3. Helm performs rolling upgrade of operator

### Rollback Capability

If an upgrade fails:
```bash
pulumi rollback
```
Or use Helm directly:
```bash
helm rollback -n elastic-system eck-operator
```

## Integration Points

### OpenMCF

- Uses shared Kubernetes provider utilities
- Integrates with Planton credential management
- Follows Planton labeling conventions
- Supports Planton's GitOps workflows

### Elastic Stack

After operator deployment, users can create:
- **Elasticsearch** clusters with StatefulSets
- **Kibana** instances connected to Elasticsearch
- **APM Server** for application performance monitoring
- **Beats** for data collection
- **Enterprise Search** for search applications

### Kubernetes Ecosystem

Integrates with:
- **Cert-Manager**: For TLS certificate management
- **Prometheus**: Operator metrics available at `:9090/metrics`
- **External Secrets Operator**: For secret injection
- **Network Policies**: For pod-level traffic control

## Testing

### Unit Testing

The module includes protobuf validation tests:

```bash
cd ../../
go test -v
```

### Integration Testing

Deploy to test cluster:

```bash
pulumi stack init test
pulumi up --stack test
kubectl get all -n elastic-system
pulumi destroy --stack test
```

## Performance Considerations

### Resource Sizing

Operator resource needs scale with:
- Number of managed Elasticsearch clusters
- Size of managed clusters (node count)
- Update frequency
- Custom resource complexity

### Recommended Scaling

| Managed Clusters | CPU Request | Memory Request |
|------------------|-------------|----------------|
| 1-2 (small) | 50m | 100Mi |
| 3-5 (medium) | 100-200m | 256-512Mi |
| 6-10 (large) | 200-500m | 512Mi-1Gi |
| 10+ (enterprise) | 500m-1000m | 1-2Gi |

## Limitations

### Single Namespace

The operator is deployed to a fixed namespace (`elastic-system`). Multiple operators in the same cluster require namespace customization (not currently supported).

### Version Pinning

The Helm chart version is pinned in code. Automatic updates are not supported to ensure stability and predictability.

### Cluster-Wide Scope

The ECK operator has cluster-wide permissions by default. Namespace-scoped deployments require additional configuration.

## Future Enhancements

Potential improvements:
- Support for High Availability (multiple operator replicas)
- Configurable namespace deployment
- Webhook configuration for validating/mutating admission
- Integration with Planton's monitoring stack
- Support for enterprise license key injection

## References

- [ECK Operator Source](https://github.com/elastic/cloud-on-k8s)
- [ECK Documentation](https://www.elastic.co/guide/en/cloud-on-k8s/current/index.html)
- [Helm Chart Repository](https://helm.elastic.co)
- [Pulumi Kubernetes Provider](https://www.pulumi.com/registry/packages/kubernetes/)


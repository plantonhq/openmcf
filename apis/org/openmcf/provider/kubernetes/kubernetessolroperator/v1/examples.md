# KubernetesSolrOperator Examples

This document provides practical examples for deploying the Apache Solr Operator on Kubernetes clusters using the KubernetesSolrOperator resource.

## Example 1: Basic Operator Deployment with Defaults

This is the simplest deployment using all default values. Suitable for development and testing.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesSolrOperator
metadata:
  name: solr-operator-basic
spec:
  targetCluster:
    clusterName: "my-gke-cluster"
  namespace:
    value: "solr-operator-system"
  create_namespace: true
  container: {}
```

**What this creates:**
- Deploys operator version v0.9.1 (default) with default resource limits (50m CPU, 100Mi memory requests)
- Registers SolrCloud, SolrBackup, and SolrPrometheusExporter CRDs
- Operator manages all SolrCloud resources in the cluster

**When to use:**
- Quick proof-of-concept deployments
- Development environments
- Small-scale deployments managing 1-5 Solr clusters

---

## Example 2: Operator with Specific Version and Custom Resources

Specify a specific operator version and increase resources for production environments.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesSolrOperator
metadata:
  name: solr-operator-production
spec:
  targetCluster:
    clusterName: "production-gke-cluster"
  namespace:
    value: "solr-operator-system"
  create_namespace: true
  # Specify operator version - see https://github.com/apache/solr-operator/releases
  operator_version: "v0.9.1"
  container:
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 2000m
        memory: 2Gi
```

**What this creates:**
- Operator with increased baseline resources (100m CPU, 256Mi memory)
- Higher burst capacity (up to 2 CPUs, 2Gi memory)
- Better performance when reconciling many SolrCloud resources

**When to use:**
- Production environments managing 10+ Solr clusters
- High-frequency scale operations
- Environments with rapid SolrCloud creation/deletion

---

## Example 3: Operator Deployment to Specific Cluster

Deploy the operator to a specific Kubernetes cluster.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesSolrOperator
metadata:
  name: solr-operator-prod-cluster
spec:
  targetCluster:
    clusterName: "prod-k8s-cluster-01"
  namespace:
    value: "solr-operator-system"
  create_namespace: false
  container:
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 1000m
        memory: 1Gi
```

**What this creates:**
- Operator deployed to cluster identified by `prod-k8s-cluster-01` credential
- Uses specified resource allocations
- Manages SolrCloud resources only in that target cluster

**When to use:**
- Multi-cluster environments
- Deploying to specific production cluster
- When credentials are managed separately (e.g., in secrets management system)

---

## Example 4: Operator Deployment with Cluster Selector (GKE)

Deploy to a Google Kubernetes Engine (GKE) cluster in the same environment.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesSolrOperator
metadata:
  name: solr-operator-gke
spec:
  targetCluster:
    clusterName: "prod-gke-cluster"
  namespace:
    value: "solr-operator-system"
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 50m
        memory: 100Mi
      limits:
        cpu: 1000m
        memory: 1Gi
```

**What this creates:**
- Operator deployed to GKE cluster in the same environment
- Automatic cluster discovery based on kind selector
- Standard resource allocations

**When to use:**
- Environment-based deployments (dev/staging/prod)
- When cluster credentials are environment-scoped
- GCP-native deployments with GKE

**Valid cluster kinds:**
- `400`: Azure AKS Cluster
- `615`: GCP GKE Cluster Core

---

## Example 5: Minimal Resource Operator (Cost-Optimized)

Ultra-minimal operator deployment for non-production or heavily constrained environments.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesSolrOperator
metadata:
  name: solr-operator-minimal
spec:
  targetCluster:
    clusterName: "dev-cluster"
  namespace:
    value: "solr-operator-system"
  create_namespace: false
  container:
    resources:
      requests:
        cpu: 25m
        memory: 64Mi
      limits:
        cpu: 500m
        memory: 512Mi
```

**What this creates:**
- Operator with minimal resource footprint
- Reduced burst capacity
- Lower baseline cost

**When to use:**
- Resource-constrained clusters
- Cost optimization in non-critical environments
- Managing 1-2 small Solr clusters

**Warning:** This configuration may struggle with large-scale reconciliations or rapid changes. Not recommended for production.

---

## Example 6: High-Performance Operator (Enterprise)

Maximum performance configuration for large-scale enterprise deployments.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesSolrOperator
metadata:
  name: solr-operator-enterprise
spec:
  targetCluster:
    clusterName: "enterprise-prod-cluster"
  namespace:
    value: "solr-operator-system"
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 250m
        memory: 512Mi
      limits:
        cpu: 4000m
        memory: 4Gi
```

**What this creates:**
- Operator with substantial guaranteed resources
- High burst capacity for intensive reconciliation workloads
- Optimized for managing 50+ Solr clusters

**When to use:**
- Large enterprise deployments
- Multi-tenant platforms offering Solr-as-a-Service
- High-throughput environments with frequent scale operations

---

## Example 7: Using Existing Namespace

Deploy operator to a pre-existing namespace managed separately.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesSolrOperator
metadata:
  name: solr-operator-existing-ns
spec:
  targetCluster:
    clusterName: "my-cluster"
  namespace:
    value: "solr-operator-system"
  create_namespace: false  # Namespace must already exist
  container:
    resources:
      requests:
        cpu: 50m
        memory: 100Mi
      limits:
        cpu: 1000m
        memory: 1Gi
```

**What this creates:**
- Deploys operator to existing `solr-operator-system` namespace
- Namespace must be created beforehand (manually or via KubernetesNamespace resource)
- Ideal for GitOps workflows where namespaces are managed separately

**When to use:**
- Namespace is managed by platform team
- Using KubernetesNamespace resource for namespace lifecycle
- RBAC policies require separate namespace management
- Multi-component deployments sharing a namespace

---

## Post-Deployment: Creating SolrCloud Clusters

After deploying the operator, you can create SolrCloud clusters:

### Simple SolrCloud Example

```yaml
apiVersion: solr.apache.org/v1beta1
kind: SolrCloud
metadata:
  name: example-solr
  namespace: solr-production
spec:
  replicas: 3
  solrImage:
    repository: solr
    tag: 8.11.3
  solrJavaMem: "-Xms1g -Xmx2g"
  zookeeperRef:
    provided:
      replicas: 3
      persistence:
        spec:
          resources:
            requests:
              storage: 5Gi
  dataStorage:
    persistent:
      reclaimPolicy: Retain
      pvcTemplate:
        spec:
          resources:
            requests:
              storage: 20Gi
```

### Production SolrCloud with Backup

```yaml
apiVersion: solr.apache.org/v1beta1
kind: SolrCloud
metadata:
  name: production-solr
  namespace: solr-production
spec:
  replicas: 5
  solrImage:
    repository: solr
    tag: 8.11.3
  solrJavaMem: "-Xms4g -Xmx8g"
  solrGCTune: "-XX:+UseG1GC"
  
  # External ZooKeeper (recommended for production)
  zookeeperRef:
    connectionInfo:
      zkConnectionString: "zk-1.example.com:2181,zk-2.example.com:2181,zk-3.example.com:2181/solr"
  
  # Persistent storage
  dataStorage:
    persistent:
      reclaimPolicy: Retain
      pvcTemplate:
        spec:
          storageClassName: fast-ssd
          resources:
            requests:
              storage: 100Gi
  
  # Resource limits
  customSolrKubeOptions:
    podOptions:
      resources:
        requests:
          cpu: "2"
          memory: 8Gi
        limits:
          cpu: "4"
          memory: 16Gi
      
      # Pod anti-affinity for high availability
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchExpressions:
                  - key: technology
                    operator: In
                    values:
                      - solr-cloud
              topologyKey: kubernetes.io/hostname
```

---

## Common Operations

### Checking Operator Status

```bash
# Get operator pod
kubectl get pods -n solr-operator-system

# View operator logs
kubectl logs -n solr-operator-system <operator-pod-name>

# Check installed CRDs
kubectl get crds | grep solr
```

### Managing SolrCloud Clusters

```bash
# List all SolrCloud instances
kubectl get solrcloud --all-namespaces

# Describe a specific SolrCloud
kubectl describe solrcloud production-solr -n solr-production

# Scale SolrCloud replicas
kubectl patch solrcloud production-solr -n solr-production \
  --type='json' \
  -p='[{"op": "replace", "path": "/spec/replicas", "value": 7}]'
```

### Backup Operations

```yaml
apiVersion: solr.apache.org/v1beta1
kind: SolrBackup
metadata:
  name: daily-backup
  namespace: solr-production
spec:
  solrCloud: production-solr
  repositoryName: s3-backup-repo
  recurrence:
    schedule: "0 2 * * *"  # Daily at 2 AM
  location: "s3://my-backup-bucket/solr-backups"
```

---

## Resource Planning

### Operator Resource Usage

| Scenario | Clusters Managed | CPU Request | Memory Request | CPU Limit | Memory Limit |
|----------|------------------|-------------|----------------|-----------|--------------|
| Small    | 1-5              | 50m         | 100Mi          | 1000m     | 1Gi          |
| Medium   | 5-20             | 100m        | 256Mi          | 2000m     | 2Gi          |
| Large    | 20-50            | 250m        | 512Mi          | 4000m     | 4Gi          |
| Enterprise| 50+             | 500m        | 1Gi            | 8000m     | 8Gi          |

### Total Cluster Resource Requirements

Remember to account for:
- **Operator Pod**: Resources specified in KubernetesSolrOperator spec
- **SolrCloud Pods**: Resources specified in each SolrCloud spec
- **ZooKeeper Ensemble**: If using operator-provided ZooKeeper (3+ pods)
- **Monitoring**: Prometheus exporters if enabled

---

## Troubleshooting

### Operator Not Creating Resources

```bash
# Check operator logs for errors
kubectl logs -n solr-operator-system <operator-pod>

# Check SolrCloud status
kubectl get solrcloud <name> -n <namespace> -o yaml

# Look for events
kubectl get events -n <namespace> --sort-by='.lastTimestamp'
```

### High Operator Resource Usage

If the operator pod is hitting resource limits:

1. Check number of managed SolrCloud instances
2. Review reconciliation frequency in operator config
3. Increase resource limits in operator spec
4. Consider splitting clusters across multiple operators

---

## Best Practices

1. **One Operator Per Cluster**: Deploy a single operator to manage all SolrCloud instances in a cluster
2. **Namespace Isolation**: Deploy operator in dedicated namespace (e.g., `solr-operator-system`)
3. **Resource Monitoring**: Set up Prometheus monitoring for operator metrics
4. **Version Pinning**: Specify exact Solr and operator versions for reproducibility
5. **Test Upgrades**: Test operator upgrades in non-production environments first
6. **Backup Strategy**: Implement regular backups using SolrBackup CRD
7. **External ZooKeeper**: Use external ZooKeeper ensemble for production SolrCloud clusters

---

## Additional Resources

- **Apache Solr Operator Documentation**: https://apache.github.io/solr-operator/
- **SolrCloud CRD Reference**: https://apache.github.io/solr-operator/docs/solr-cloud/
- **Backup & Restore Guide**: https://apache.github.io/solr-operator/docs/solr-backup/
- **Research Documentation**: [docs/README.md](docs/README.md) - Detailed deployment patterns and decision framework
- **Component README**: [README.md](README.md) - Full API reference and architecture overview


# KubernetesElasticOperator - Usage Examples

This document provides practical examples for deploying the Elastic Cloud on Kubernetes (ECK) operator using the **KubernetesElasticOperator** API resource.

> **Note:** After creating a YAML file with your configuration, apply it using:
> ```shell
> planton apply -f <yaml-file>
> ```

---

## 1. Basic Installation with Namespace Creation

Deploy ECK operator with default resource configuration, creating a new namespace.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesElasticOperator
metadata:
  name: eck-operator
spec:
  target_cluster:
    kubernetes_credential_id: "my-k8s-cluster"
  namespace:
    value: "elastic-system"
  create_namespace: true
  container:
    resources:
      requests:
        cpu: "50m"
        memory: "100Mi"
      limits:
        cpu: "1000m"
        memory: "1Gi"
```

**What this does:**
- Creates the `elastic-system` namespace with appropriate labels
- Installs ECK operator version 2.14.0 in the created namespace
- Allocates 50m CPU and 100Mi memory as baseline (requests)
- Limits operator to 1 CPU core and 1Gi memory maximum
- Uses Kubernetes credential ID to access the target cluster

**When to use:** Standard production deployment for managing moderate-sized Elastic Stack installations (1-5 Elasticsearch clusters) when you want the component to manage the namespace lifecycle.

---

## 2. Using Existing Namespace

Deploy ECK operator into an existing namespace that's managed elsewhere.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesElasticOperator
metadata:
  name: eck-operator
spec:
  target_cluster:
    kubernetes_credential_id: "my-k8s-cluster"
  namespace:
    value: "platform-operators"
  create_namespace: false
  container:
    resources:
      requests:
        cpu: "50m"
        memory: "100Mi"
      limits:
        cpu: "1000m"
        memory: "1Gi"
```

**What this does:**
- Uses the existing `platform-operators` namespace
- Does NOT create or manage the namespace
- Installs ECK operator version 2.14.0 in the existing namespace
- Allocates standard resource requests and limits

**When to use:**
- When namespace is managed by another component (e.g., KubernetesNamespace)
- When namespace has specific resource quotas or policies managed externally
- When deploying multiple operators into a shared namespace
- When namespace lifecycle should be independent of this component

**Prerequisites:**
- The namespace must exist before applying this manifest
- You must have permissions to deploy resources in the namespace

---

## 3. High-Availability Production Configuration

For large-scale production environments managing multiple Elasticsearch clusters.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesElasticOperator
metadata:
  name: eck-operator-prod
  labels:
    environment: production
    team: platform
spec:
  target_cluster:
    kubernetes_credential_id: "production-cluster"
  namespace:
    value: "elastic-system"
  create_namespace: true
  container:
    resources:
      requests:
        cpu: "200m"
        memory: "512Mi"
      limits:
        cpu: "2000m"
        memory: "2Gi"
```

**What this does:**
- Increases resource allocation for operator to handle higher workload
- 200m CPU and 512Mi memory baseline guarantees responsive operation
- Allows operator to scale up to 2 CPU cores and 2Gi memory under load
- Adds labels for organizational tracking

**When to use:** Production environments with:
- 5+ Elasticsearch clusters
- Frequent scaling operations
- Large cluster sizes (10+ nodes per Elasticsearch cluster)
- High update frequency

---

## 3. Development/Testing Environment

Minimal resource allocation for local development or testing.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesElasticOperator
metadata:
  name: eck-operator-dev
spec:
  target_cluster:
    kubernetes_credential_id: "dev-cluster"
  namespace:
    value: "elastic-dev"
  create_namespace: true
  container:
    resources:
      requests:
        cpu: "25m"
        memory: "64Mi"
      limits:
        cpu: "500m"
        memory: "512Mi"
```

**What this does:**
- Creates the `elastic-dev` namespace for development isolation
- Reduces resource footprint for development/testing scenarios
- 25m CPU and 64Mi memory minimal baseline
- Limited to 500m CPU (0.5 cores) and 512Mi memory
- Suitable for managing 1-2 small Elasticsearch clusters

**When to use:**
- Local development with minikube or kind
- CI/CD test environments
- Proof-of-concept deployments
- Learning and experimentation

---

## 4. Using Cluster Selector

Select the target Kubernetes cluster using a selector instead of credential ID.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesElasticOperator
metadata:
  name: eck-operator-selector
spec:
  target_cluster:
    kubernetes_cluster_selector:
      name: "staging-cluster"
  namespace:
    value: "elastic-system"
  create_namespace: true
  container:
    resources:
      requests:
        cpu: "50m"
        memory: "100Mi"
      limits:
        cpu: "1000m"
        memory: "1Gi"
```

**What this does:**
- Uses cluster selector to identify the target Kubernetes cluster
- Creates the `elastic-system` namespace
- Useful when cluster is in the same environment as the operator
- Provides dynamic cluster selection based on naming

**When to use:**
- Multi-cluster environments with naming conventions
- When credential management is handled separately
- GitOps workflows with cluster auto-discovery

---

## 5. Multi-Environment Deployment Pattern

Example showing how to deploy ECK operator across multiple environments with environment-specific configurations.

### Production

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesElasticOperator
metadata:
  name: eck-operator-prod
  labels:
    environment: production
spec:
  target_cluster:
    kubernetes_credential_id: "prod-cluster"
  namespace:
    value: "elastic-system"
  create_namespace: true
  container:
    resources:
      requests:
        cpu: "200m"
        memory: "512Mi"
      limits:
        cpu: "2000m"
        memory: "2Gi"
```

### Staging

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesElasticOperator
metadata:
  name: eck-operator-staging
  labels:
    environment: staging
spec:
  target_cluster:
    kubernetes_credential_id: "staging-cluster"
  namespace:
    value: "elastic-system"
  create_namespace: true
  container:
    resources:
      requests:
        cpu: "100m"
        memory: "256Mi"
      limits:
        cpu: "1000m"
        memory: "1Gi"
```

### Development

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesElasticOperator
metadata:
  name: eck-operator-dev
  labels:
    environment: development
spec:
  target_cluster:
    kubernetes_credential_id: "dev-cluster"
  namespace:
    value: "elastic-dev"
  create_namespace: true
  container:
    resources:
      requests:
        cpu: "25m"
        memory: "64Mi"
      limits:
        cpu: "500m"
        memory: "512Mi"
```

**Pattern Benefits:**
- **Environment Parity**: Same deployment process across all environments
- **Resource Optimization**: Right-sized resources per environment
- **Cost Efficiency**: Lower resource usage in non-production
- **Testing Confidence**: Test changes in dev/staging before production

---

## Post-Installation: Creating Elastic Stack Resources

After the ECK operator is installed, you can create Elasticsearch, Kibana, and other Elastic Stack components.

### Example: Elasticsearch Cluster

```yaml
apiVersion: elasticsearch.k8s.elastic.co/v1
kind: Elasticsearch
metadata:
  name: production-es
  namespace: default
spec:
  version: 8.11.0
  nodeSets:
  - name: default
    count: 3
    config:
      node.store.allow_mmap: false
    volumeClaimTemplates:
    - metadata:
        name: elasticsearch-data
      spec:
        accessModes:
        - ReadWriteOnce
        resources:
          requests:
            storage: 100Gi
        storageClassName: fast-ssd
    podTemplate:
      spec:
        containers:
        - name: elasticsearch
          resources:
            requests:
              cpu: 1
              memory: 4Gi
            limits:
              cpu: 2
              memory: 8Gi
```

### Example: Kibana Instance

```yaml
apiVersion: kibana.k8s.elastic.co/v1
kind: Kibana
metadata:
  name: production-kibana
  namespace: default
spec:
  version: 8.11.0
  count: 2
  elasticsearchRef:
    name: production-es
  podTemplate:
    spec:
      containers:
      - name: kibana
        resources:
          requests:
            cpu: 500m
            memory: 1Gi
          limits:
            cpu: 1
            memory: 2Gi
```

### Example: APM Server

```yaml
apiVersion: apm.k8s.elastic.co/v1
kind: ApmServer
metadata:
  name: production-apm
  namespace: default
spec:
  version: 8.11.0
  count: 1
  elasticsearchRef:
    name: production-es
  podTemplate:
    spec:
      containers:
      - name: apm-server
        resources:
          requests:
            cpu: 200m
            memory: 512Mi
          limits:
            cpu: 500m
            memory: 1Gi
```

---

## Resource Sizing Guidelines

### Operator Resource Needs

The ECK operator's resource requirements depend on:

| Scenario | Clusters Managed | Recommended CPU | Recommended Memory |
|----------|------------------|-----------------|-------------------|
| Small (Dev/Test) | 1-2 clusters | 25-50m | 64-100Mi |
| Medium (Production) | 3-5 clusters | 100-200m | 256-512Mi |
| Large (Enterprise) | 6-10 clusters | 200-500m | 512Mi-1Gi |
| Extra Large | 10+ clusters | 500m-1000m | 1-2Gi |

### Factors Affecting Resource Usage

- **Cluster Size**: Larger Elasticsearch clusters require more operator resources
- **Update Frequency**: Frequent rolling updates increase operator workload
- **CRD Complexity**: Complex custom resources (many sidecars, volumes) use more memory
- **Monitoring Overhead**: Operator continuously watches all managed resources

---

## Verification

After deployment, verify the ECK operator is running:

```bash
# Check operator pod
kubectl get pods -n elastic-system

# View operator logs
kubectl logs -n elastic-system -l control-plane=elastic-operator

# Verify CRDs are installed
kubectl get crds | grep elastic

# Check operator version
kubectl get deployment -n elastic-system elastic-operator -o jsonpath='{.spec.template.spec.containers[0].image}'
```

Expected CRDs installed:
- `elasticsearch.k8s.elastic.co`
- `kibana.k8s.elastic.co`
- `apmserver.k8s.elastic.co`
- `enterprisesearch.k8s.elastic.co`
- `beat.k8s.elastic.co`
- `agent.k8s.elastic.co`
- `logstash.k8s.elastic.co`

---

## Troubleshooting

### Operator Pod Not Starting

```bash
# Check pod events
kubectl describe pod -n elastic-system -l control-plane=elastic-operator

# Check resource availability
kubectl top nodes
kubectl describe node <node-name>
```

### CRDs Not Installing

```bash
# Verify CRDs
kubectl get crds | grep elastic

# Check operator logs for CRD errors
kubectl logs -n elastic-system -l control-plane=elastic-operator | grep CRD
```

### Permission Issues

```bash
# Check operator service account
kubectl get serviceaccount -n elastic-system elastic-operator

# Verify RBAC
kubectl get clusterrole elastic-operator
kubectl get clusterrolebinding elastic-operator
```

---

## Next Steps

1. **Deploy Elasticsearch**: Create an Elasticsearch cluster using the `Elasticsearch` CRD
2. **Deploy Kibana**: Connect Kibana to your Elasticsearch cluster
3. **Configure Monitoring**: Set up monitoring for your Elastic Stack deployments
4. **Set Up Backups**: Configure snapshot repositories for Elasticsearch backups
5. **Implement Security**: Configure RBAC, network policies, and TLS

For more information:
- [Elastic ECK Documentation](https://www.elastic.co/guide/en/cloud-on-k8s/current/index.html)
- [ECK Quickstart Guide](https://www.elastic.co/guide/en/cloud-on-k8s/current/k8s-quickstart.html)
- [OpenMCF Docs](https://openmcf.org)


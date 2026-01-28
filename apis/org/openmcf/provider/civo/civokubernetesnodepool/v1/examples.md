# Civo Kubernetes Node Pool Examples

Real-world node pool configurations for heterogeneous Kubernetes cluster workloads.

## Table of Contents

1. [Multi-Pool Cluster Architecture](#1-multi-pool-cluster-architecture)
2. [Autoscaling Web Application Pool](#2-autoscaling-web-application-pool)
3. [Fixed-Size Database Pool](#3-fixed-size-database-pool)
4. [Batch Processing Pool (Scale to Zero)](#4-batch-processing-pool-scale-to-zero)
5. [GPU Workload Pool (Future)](#5-gpu-workload-pool-future)

---

## 1. Multi-Pool Cluster Architecture

**Use Case:** Separate workloads by resource requirements across multiple node pools.

### General-Purpose Workers

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoKubernetesNodePool
metadata:
  name: general-workers
spec:
  nodePoolName: general-workers
  cluster:
    value: "prod-cluster"
  size: g4s.kube.medium
  nodeCount: 5
  autoScale: true
  minNodes: 3
  maxNodes: 10
  tags:
    - workload:general
    - tier:workers
```

### Memory-Intensive Pool

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoKubernetesNodePool
metadata:
  name: memory-intensive-pool
spec:
  nodePoolName: memory-intensive-workers
  cluster:
    value: "prod-cluster"
  size: g4s.kube.xlarge
  nodeCount: 2
  autoScale: true
  minNodes: 1
  maxNodes: 5
  tags:
    - workload:memory-intensive
    - team:data-engineering
```

### Batch Processing Pool

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoKubernetesNodePool
metadata:
  name: batch-pool
spec:
  nodePoolName: batch-workers
  cluster:
    value: "prod-cluster"
  size: g4s.kube.large
  nodeCount: 0
  autoScale: true
  minNodes: 0
  maxNodes: 20
  tags:
    - workload:batch
    - autoscale:aggressive
```

**Cost estimate:**
- General: 3-10 medium nodes = $65-$217/month
- Memory: 1-5 xlarge nodes = $87-$435/month  
- Batch: 0-20 large nodes = $0-$869/month (scales to zero when idle)
- **Total range**: $152-$1,521/month

**Deploy all pools:**

```bash
planton apply -f general-workers.yaml
planton apply -f memory-intensive-pool.yaml
planton apply -f batch-pool.yaml

# Verify
kubectl get nodes
```

**Workload placement:**

```yaml
# General workload (default pool)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-app
spec:
  template:
    spec:
      nodeSelector:
        pool: general-workers

# Memory-intensive workload
apiVersion: apps/v1
kind: Deployment
metadata:
  name: analytics-engine
spec:
  template:
    spec:
      nodeSelector:
        pool: memory-intensive-workers
      resources:
        requests:
          memory: "8Gi"
          cpu: "2000m"

# Batch job
apiVersion: batch/v1
kind: Job
metadata:
  name: data-processing
spec:
  template:
    spec:
      nodeSelector:
        pool: batch-workers
      tolerations:
        - key: workload
          value: batch
          effect: NoSchedule
```

---

## 2. Autoscaling Web Application Pool

**Use Case:** Handle variable web traffic with automatic node scaling.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoKubernetesNodePool
metadata:
  name: web-autoscale-pool
spec:
  nodePoolName: web-workers
  cluster:
    value: "web-cluster"
  size: g4s.kube.medium
  nodeCount: 3  # Initial size
  autoScale: true
  minNodes: 2   # Minimum availability
  maxNodes: 12  # Peak capacity
  tags:
    - workload:web
    - environment:production
```

**Deployment with HPA:**

```yaml
# Deploy web application
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-app
spec:
  replicas: 10  # Will autoscale via HPA
  template:
    spec:
      containers:
        - name: app
          resources:
            requests:
              cpu: "200m"
              memory: "256Mi"
            limits:
              cpu: "500m"
              memory: "512Mi"
---
# Horizontal Pod Autoscaler
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: web-app-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: web-app
  minReplicas: 10
  maxReplicas: 50
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
```

**How it works:**
1. Traffic increases → CPU usage rises
2. HPA scales pods from 10 to 50
3. Pods can't fit on 3 nodes
4. Cluster autoscaler adds nodes (up to 12)
5. Pods are scheduled on new nodes
6. Traffic decreases → HPA scales down pods
7. Nodes become underutilized
8. Cluster autoscaler removes nodes (down to 2)

---

## 3. Fixed-Size Database Pool

**Use Case:** Dedicated nodes for stateful databases with predictable resources.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoKubernetesNodePool
metadata:
  name: database-pool
spec:
  nodePoolName: database-workers
  cluster:
    value: "prod-cluster"
  size: g4s.kube.large
  nodeCount: 3  # Fixed size - no autoscaling
  tags:
    - workload:database
    - stateful:true
```

**Why no autoscaling?**
- Stateful workloads need stable nodes
- Database pods use local storage (PVs)
- Scaling nodes might orphan data
- Predictable cost and capacity

**Taint nodes for databases only:**

```bash
# After pool creation, taint nodes
kubectl taint nodes -l pool=database-workers workload=database:NoSchedule

# Deploy database with toleration
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgresql
spec:
  template:
    spec:
      tolerations:
        - key: workload
          value: database
          effect: NoSchedule
      nodeSelector:
        pool: database-workers
      volumes:
        - name: data
          persistentVolumeClaim:
            claimName: postgres-pvc
```

---

## 4. Batch Processing Pool (Scale to Zero)

**Use Case:** Run scheduled batch jobs without idle node costs.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoKubernetesNodePool
metadata:
  name: batch-processing-pool
spec:
  nodePoolName: batch-workers
  cluster:
    value: "data-cluster"
  size: g4s.kube.large
  nodeCount: 0  # Start with zero nodes
  autoScale: true
  minNodes: 0   # Scale to zero when idle
  maxNodes: 20  # Scale up to 20 for large jobs
  tags:
    - workload:batch
    - cost-optimize:true
```

**Deploy batch jobs:**

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: nightly-reports
spec:
  schedule: "0 2 * * *"  # 2 AM daily
  jobTemplate:
    spec:
      template:
        spec:
          nodeSelector:
            pool: batch-workers
          containers:
            - name: report-generator
              image: myapp/report-generator:latest
              resources:
                requests:
                  cpu: "1000m"
                  memory: "2Gi"
          restartPolicy: OnFailure
```

**Cost savings:**

- **Without scale-to-zero**: 2 nodes × $43 × 730 hrs = $63/month idle cost
- **With scale-to-zero**: Pay only during job execution (~2 hours/day = $5.60/month)
- **Savings**: $57/month (~90% reduction)

---

## 5. GPU Workload Pool (Future)

**Note**: Civo doesn't currently offer GPU instances. When available:

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoKubernetesNodePool
metadata:
  name: gpu-pool
spec:
  nodePoolName: gpu-workers
  cluster:
    value: "ml-cluster"
  size: g4s.kube.gpu.large  # Hypothetical GPU instance
  nodeCount: 1
  autoScale: true
  minNodes: 0
  maxNodes: 5
  tags:
    - workload:ml-training
    - gpu:enabled
```

---

## Advanced Patterns

### Pool per Team

Isolate teams with separate pools:

```yaml
# Team A
- nodePoolName: team-a-workers
  size: g4s.kube.medium
  nodeCount: 3
  tags: ["team:a"]

# Team B
- nodePoolName: team-b-workers
  size: g4s.kube.large
  nodeCount: 2
  tags: ["team:b"]
```

### Development vs Production Pools

```yaml
# Dev pool (smaller, more flexible)
- nodePoolName: dev-workers
  size: g4s.kube.small
  nodeCount: 2
  autoScale: true
  minNodes: 1
  maxNodes: 5

# Prod pool (larger, stable)
- nodePoolName: prod-workers
  size: g4s.kube.large
  nodeCount: 5
  autoScale: false  # Fixed size for stability
```

---

## Additional Resources

- [Main README](README.md) - Component overview
- [Research Documentation](docs/README.md) - Deep dive into node pool management
- [Pulumi Module](iac/pulumi/README.md) - Direct Pulumi usage
- [Civo API](https://www.civo.com/api/kubernetes)

---

## Need Help?

- Check [README.md](README.md#troubleshooting) for troubleshooting
- Open issue on [GitHub](https://github.com/plantonhq/openmcf/issues)
- Contact Civo: support@civo.com


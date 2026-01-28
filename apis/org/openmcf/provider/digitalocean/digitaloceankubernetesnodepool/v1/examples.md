# DigitalOcean Kubernetes Node Pool Examples

This document provides practical examples for deploying node pools to DigitalOcean Kubernetes (DOKS) clusters using the OpenMCF API. Each example is a complete, deployable manifest demonstrating different use cases and configuration patterns.

## Table of Contents

- [Example 1: Basic Fixed-Size Node Pool](#example-1-basic-fixed-size-node-pool)
- [Example 2: Autoscaling Production Pool](#example-2-autoscaling-production-pool)
- [Example 3: GPU Pool with Taints](#example-3-gpu-pool-with-taints)
- [Example 4: Multi-Pool Architecture](#example-4-multi-pool-architecture)
- [Example 5: Labeled Pool for Specific Workloads](#example-5-labeled-pool-for-specific-workloads)

---

## Example 1: Basic Fixed-Size Node Pool

**Use Case:** Simple, fixed-size node pool for stable workloads.

**Configuration:**
- **Pool Name:** app-workers
- **Node Count:** 3 (fixed, no autoscaling)
- **Instance Size:** s-2vcpu-4gb (~$20/month per node = ~$60/month total)
- **No special labels or taints**

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanKubernetesNodePool
metadata:
  name: app-workers
spec:
  nodePoolName: app-workers
  cluster:
    value: "your-cluster-id-here"  # Replace with your cluster ID
  size: s-2vcpu-4gb
  nodeCount: 3
  tags:
    - env:dev
    - team:engineering
```

**Deployment:**

```bash
# Save the manifest
cat > basic-pool.yaml << 'EOF'
[paste manifest above]
EOF

# Deploy using OpenMCF CLI
openmcf apply -f basic-pool.yaml

# Verify pool creation
openmcf get digitaloceankubernetesnodepool app-workers
```

---

## Example 2: Autoscaling Production Pool

**Use Case:** Production workload pool that scales based on demand.

**Configuration:**
- **Pool Name:** prod-autoscale-workers
- **Initial Nodes:** 5
- **Autoscaling:** Enabled (min: 3, max: 10)
- **Instance Size:** s-4vcpu-8gb (~$43/month per node)
- **Labels:** workload=application, env=production
- **Tags:** Cost attribution tags

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanKubernetesNodePool
metadata:
  name: prod-autoscale-workers
spec:
  nodePoolName: prod-autoscale-workers
  cluster:
    value: "your-cluster-id-here"
  size: s-4vcpu-8gb
  nodeCount: 5
  autoScale: true
  minNodes: 3
  maxNodes: 10
  labels:
    workload: application
    env: production
    tier: backend
  tags:
    - env:production
    - cost-center:engineering
    - autoscaling:enabled
```

**Deployment:**

```bash
# Save and deploy
cat > autoscale-pool.yaml << 'EOF'
[paste manifest above]
EOF

openmcf apply -f autoscale-pool.yaml

# Monitor autoscaling behavior
kubectl get nodes -l workload=application

# Watch pod scheduling
kubectl get pods -o wide
```

---

## Example 3: GPU Pool with Taints

**Use Case:** Specialized GPU pool for ML/AI workloads, isolated from general workloads.

**Configuration:**
- **Pool Name:** gpu-workers
- **Node Count:** 2 (expensive, keep minimal)
- **Instance Size:** g-4vcpu-16gb (GPU-enabled Droplet)
- **Taints:** Prevents non-GPU pods from scheduling
- **Labels:** hardware=gpu for pod node selection

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanKubernetesNodePool
metadata:
  name: gpu-workers
spec:
  nodePoolName: gpu-workers
  cluster:
    value: "your-cluster-id-here"
  size: g-4vcpu-16gb
  nodeCount: 2
  autoScale: true
  minNodes: 0  # Scale to zero when not in use
  maxNodes: 4
  labels:
    hardware: gpu
    workload: ml
  taints:
    - key: nvidia.com/gpu
      value: "true"
      effect: NoSchedule
  tags:
    - hardware:gpu
    - cost-center:ml-team
```

**Deployment:**

```bash
# Save and deploy
cat > gpu-pool.yaml << 'EOF'
[paste manifest above]
EOF

openmcf apply -f gpu-pool.yaml

# Verify taints are applied
kubectl get nodes -l hardware=gpu -o json | jq '.items[].spec.taints'
```

**Using the GPU Pool:**

Pods must tolerate the GPU taint to schedule on these nodes:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: gpu-training-job
spec:
  tolerations:
    - key: nvidia.com/gpu
      operator: "Equal"
      value: "true"
      effect: NoSchedule
  nodeSelector:
    hardware: gpu
  containers:
    - name: training
      image: tensorflow/tensorflow:latest-gpu
      resources:
        limits:
          nvidia.com/gpu: 1
```

---

## Example 4: Multi-Pool Architecture

**Use Case:** Complete production architecture with separate pools for system services, applications, and batch jobs.

### Pool 1: System Services (Sacrificial Default Pool)

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanKubernetesNodePool
metadata:
  name: system-default
spec:
  nodePoolName: system-default
  cluster:
    value: "your-cluster-id-here"
  size: s-1vcpu-2gb  # Minimal size
  nodeCount: 1  # Never touch this pool after creation
  labels:
    workload: system
    pool-type: sacrificial-default
  taints:
    - key: system-only
      value: "true"
      effect: NoSchedule
  tags:
    - pool:system-default
    - managed:do-not-modify
```

### Pool 2: Application Workers

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanKubernetesNodePool
metadata:
  name: app-workers
spec:
  nodePoolName: app-workers
  cluster:
    value: "your-cluster-id-here"
  size: s-4vcpu-8gb
  nodeCount: 5
  autoScale: true
  minNodes: 3
  maxNodes: 10
  labels:
    workload: application
    tier: web
  tags:
    - pool:application
    - env:production
```

### Pool 3: Batch Processing

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanKubernetesNodePool
metadata:
  name: batch-workers
spec:
  nodePoolName: batch-workers
  cluster:
    value: "your-cluster-id-here"
  size: c-8vcpu-16gb  # CPU-optimized
  nodeCount: 2
  autoScale: true
  minNodes: 0  # Scale to zero when no jobs
  maxNodes: 8
  labels:
    workload: batch
    compute-type: cpu-optimized
  taints:
    - key: workload
      value: batch
      effect: PreferNoSchedule  # Prefer but don't require
  tags:
    - pool:batch
    - autoscaling:aggressive
```

**Deployment:**

```bash
# Deploy all pools
openmcf apply -f system-pool.yaml
openmcf apply -f app-pool.yaml
openmcf apply -f batch-pool.yaml

# Verify multi-pool setup
kubectl get nodes --show-labels
kubectl get nodes -o custom-columns=NAME:.metadata.name,POOL:.metadata.labels.dopl-cloud-node-pool
```

---

## Example 5: Labeled Pool for Specific Workloads

**Use Case:** Pool dedicated to specific application teams or services.

**Configuration:**
- **Pool Name:** frontend-workers
- **Labels:** Multiple labels for precise targeting
- **Node Count:** 4 nodes
- **Autoscaling:** Enabled

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanKubernetesNodePool
metadata:
  name: frontend-workers
spec:
  nodePoolName: frontend-workers
  cluster:
    value: "your-cluster-id-here"
  size: s-2vcpu-4gb
  nodeCount: 4
  autoScale: true
  minNodes: 2
  maxNodes: 8
  labels:
    workload: frontend
    team: web
    app: ecommerce
    tier: presentation
  tags:
    - team:web
    - app:ecommerce
    - tier:frontend
```

**Using Labels for Pod Placement:**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: frontend-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: frontend
  template:
    metadata:
      labels:
        app: frontend
    spec:
      nodeSelector:
        workload: frontend
        team: web
      containers:
        - name: web
          image: nginx:latest
          resources:
            requests:
              cpu: "500m"
              memory: "512Mi"
```

---

## Common Operations

### Get Node Pool Status

```bash
# Get detailed node pool information
openmcf get digitaloceankubernetesnodepool <pool-name> -o yaml

# Get node pool outputs (ID, node IDs)
openmcf outputs digitaloceankubernetesnodepool <pool-name>
```

### Update Node Pool

```bash
# Edit the manifest file to change configuration
vim my-pool.yaml

# Apply changes (e.g., increase max_nodes for autoscaling)
openmcf apply -f my-pool.yaml
```

### Scale Node Pool Manually

```bash
# Edit manifest to change nodeCount
# Example: Change nodeCount from 3 to 5
vim my-pool.yaml

# Apply
openmcf apply -f my-pool.yaml
```

### Delete Node Pool

```bash
# Delete node pool (drains nodes first)
openmcf delete -f my-pool.yaml

# Or delete by name
openmcf delete digitaloceankubernetesnodepool <pool-name>
```

---

## Best Practices

### Pool Sizing

1. **Start Small**: Begin with minimal nodes and enable autoscaling
2. **Separate Workloads**: Use multiple pools for different workload types
3. **Right-Size**: Choose Droplet sizes based on actual resource needs
4. **Monitor Utilization**: Use `kubectl top nodes` to track resource usage

### Autoscaling

1. **Set Appropriate Boundaries**: Don't set max_nodes too high (cost control)
2. **Use Min > 0 for Critical Workloads**: Ensures capacity is always available
3. **Scale-to-Zero for Batch**: Use min_nodes=0 for non-critical batch workloads
4. **Pod Resource Requests**: Always set CPU/memory requests for autoscaler to work

### Labels and Taints

1. **Use Labels for Selection**: Pods use nodeSelector or node affinity
2. **Use Taints for Isolation**: Prevents unwanted pods on specialized nodes
3. **Document Conventions**: Maintain consistent label naming across pools
4. **Avoid Over-Tainting**: Too many taints can make scheduling difficult

### Cost Optimization

1. **Autoscaling**: Let pools scale down during off-hours
2. **Right-Size Droplets**: Don't over-provision CPU/memory
3. **Use Tags**: Track costs per pool, team, or application
4. **Scale-to-Zero**: Use min_nodes=0 for batch/dev pools

### Security

1. **Network Policies**: Implement pod-to-pod security with NetworkPolicy
2. **RBAC**: Restrict who can create/modify node pools
3. **Pod Security Standards**: Enforce security contexts and admission policies
4. **Regular Updates**: Keep nodes updated (handled by DOKS automatically)

---

## Troubleshooting

### Pool Creation Fails

**Symptom:** Pool creation times out or fails

**Possible Causes:**
- Cluster ID is invalid or cluster doesn't exist
- Droplet size not available in cluster region
- Account quota exceeded
- Invalid taint effect (must be NoSchedule, PreferNoSchedule, or NoExecute)

**Solution:**
```bash
# Verify cluster exists
openmcf get digitaloceankubernetescluster <cluster-name>

# Check available Droplet sizes
doctl kubernetes options sizes

# Check account limits
doctl account get
```

### Pods Not Scheduling on Pool

**Symptom:** Pods remain in Pending state despite pool having capacity

**Possible Causes:**
- Pod doesn't tolerate pool taints
- Pod nodeSelector doesn't match pool labels
- Insufficient resources (CPU/memory)

**Solution:**
```bash
# Check pod scheduling events
kubectl describe pod <pod-name>

# Verify node labels
kubectl get nodes --show-labels

# Check node taints
kubectl describe node <node-name> | grep Taints

# Verify pod tolerations and nodeSelector
kubectl get pod <pod-name> -o yaml
```

### Autoscaling Not Working

**Symptom:** Pool doesn't scale despite pod resource pressure

**Possible Causes:**
- Autoscaling not enabled (autoScale: false)
- min_nodes and max_nodes not configured
- Pods don't have resource requests set

**Solution:**
```bash
# Verify autoscaling is enabled
openmcf get digitaloceankubernetesnodepool <pool-name> -o yaml | grep autoScale

# Check cluster autoscaler logs
kubectl logs -n kube-system -l app=cluster-autoscaler

# Verify pod resource requests
kubectl get pods -o json | jq '.items[] | .spec.containers[] | .resources'
```

---

## Reference Links

- [DigitalOcean Kubernetes Documentation](https://docs.digitalocean.com/products/kubernetes/)
- [DOKS Node Pools](https://docs.digitalocean.com/products/kubernetes/how-to/add-node-pools/)
- [Kubernetes Node Affinity](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/)
- [Kubernetes Taints and Tolerations](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/)
- [OpenMCF Documentation](https://docs.openmcf.org/)

---

**Note:** Replace placeholder values like `your-cluster-id-here` with your actual DOKS cluster IDs before deploying. You can find cluster IDs using `doctl kubernetes cluster list` or `openmcf get digitaloceankubernetescluster`.


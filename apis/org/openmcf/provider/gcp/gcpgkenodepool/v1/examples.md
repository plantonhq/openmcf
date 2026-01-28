# GCP GKE Node Pool Examples

This document provides production-ready YAML manifests for common GKE node pool configurations. All examples assume an existing `GcpGkeCluster` resource has been deployed.

## Table of Contents

1. [Minimal Configuration](#1-minimal-configuration)
2. [Dev/Test Pool with Spot VMs](#2-devtest-pool-with-spot-vms)
3. [Production General-Purpose Pool](#3-production-general-purpose-pool)
4. [High-Memory Pool for Caches](#4-high-memory-pool-for-caches)
5. [GPU Pool for Machine Learning](#5-gpu-pool-for-machine-learning)
6. [Compute-Optimized Pool for CPU-Intensive Workloads](#6-compute-optimized-pool-for-cpu-intensive-workloads)

---

## 1. Minimal Configuration

The simplest possible node pool with only required fields. Suitable for quick testing.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeNodePool
metadata:
  name: minimal-pool
spec:
  # Required: Reference to parent GKE cluster
  clusterProjectId:
    value: my-gcp-project
  clusterName:
    value: my-cluster
  clusterLocation:
    value: us-central1
  
  # Fixed size: 3 nodes (no autoscaling)
  nodeCount: 3
```

**Defaults applied:**
- `machine_type`: `e2-medium` (2 vCPU, 4 GB RAM)
- `disk_size_gb`: 100
- `disk_type`: `pd-standard`
- `image_type`: `COS_CONTAINERD`
- Auto-upgrade and auto-repair: **enabled**

**Deploy:**
```bash
openmcf apply -f minimal-pool.yaml
```

---

## 2. Dev/Test Pool with Spot VMs

Cost-optimized pool for development and testing environments. Uses Spot VMs for up to 91% cost savings and supports scale-to-zero.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeNodePool
metadata:
  name: dev-pool
  org: my-org
  env: dev
spec:
  clusterProjectId:
    value: dev-gcp-project
  clusterName:
    value: dev-cluster
  clusterLocation:
    value: us-central1
  
  # Cost-optimized machine type
  machineType: e2-medium
  diskSizeGb: 50
  diskType: pd-standard
  
  # Enable Spot VMs for cost savings
  spot: true
  
  # Autoscaling with scale-to-zero
  autoscaling:
    minNodes: 0      # Scale to zero during off-hours
    maxNodes: 5
    locationPolicy: ANY  # Hunt for Spot capacity across all zones
  
  # Labels for workload scheduling
  nodeLabels:
    environment: dev
    spot: "true"
    cost-tier: low
  
  # Management settings
  management:
    disableAutoUpgrade: false
    disableAutoRepair: false
```

**Use case:** Development clusters, CI/CD test environments, batch processing jobs that can tolerate interruptions.

**Deploy:**
```bash
openmcf apply -f dev-pool.yaml
```

---

## 3. Production General-Purpose Pool

High-availability pool for production web applications with SSD disks and custom service account.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeNodePool
metadata:
  name: prod-general
  org: my-org
  env: prod
spec:
  clusterProjectId:
    value: prod-gcp-project
  clusterName:
    value: prod-cluster
  clusterLocation:
    value: us-central1
  
  # Production-grade machine type
  machineType: n2-standard-4
  diskSizeGb: 100
  diskType: pd-ssd  # SSD for better performance
  imageType: COS_CONTAINERD
  
  # Custom service account with minimal permissions
  serviceAccount: gke-prod-sa@prod-gcp-project.iam.gserviceaccount.com
  
  # High-availability autoscaling
  autoscaling:
    minNodes: 3      # Minimum for multi-zone HA
    maxNodes: 10
    locationPolicy: BALANCED  # Even distribution across zones
  
  # Labels for workload scheduling
  nodeLabels:
    environment: production
    workload-tier: general
    disk-type: ssd
  
  # Management settings (keep auto-upgrade/repair enabled)
  management:
    disableAutoUpgrade: false
    disableAutoRepair: false
```

**Use case:** Production web applications, microservices, general-purpose workloads requiring reliability and performance.

**Deploy:**
```bash
openmcf apply -f prod-general-pool.yaml
```

---

## 4. High-Memory Pool for Caches

Memory-optimized pool for Redis, Memcached, or in-memory databases.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeNodePool
metadata:
  name: high-memory-pool
  org: my-org
  env: prod
spec:
  clusterProjectId:
    value: prod-gcp-project
  clusterName:
    value: prod-cluster
  clusterLocation:
    value: us-central1
  
  # High-memory machine type (M2 series)
  machineType: m2-ultramem-208  # 208 vCPU, 5888 GB RAM
  diskSizeGb: 200
  diskType: pd-ssd
  
  # Fixed size for consistent memory capacity
  nodeCount: 2
  
  # Labels for targeted scheduling
  nodeLabels:
    workload-type: high-memory
    cache-tier: primary
    environment: production
  
  # Management settings
  management:
    disableAutoUpgrade: false
    disableAutoRepair: false
```

**Use case:** In-memory caches (Redis, Memcached), analytics workloads, large Java heap applications.

**Note:** For smaller memory needs, consider `n2-highmem-32` (32 vCPU, 256 GB RAM) or similar.

**Deploy:**
```bash
openmcf apply -f high-memory-pool.yaml
```

---

## 5. GPU Pool for Machine Learning

Specialized pool with GPU accelerators for ML training and inference workloads.

**Note:** GPU accelerator configuration is a roadmap item. This example shows the current configuration pattern with labels for GPU scheduling. Once GPU support is added to the spec, this example will be updated.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeNodePool
metadata:
  name: gpu-pool
  org: my-org
  env: prod
spec:
  clusterProjectId:
    value: ml-gcp-project
  clusterName:
    value: ml-cluster
  clusterLocation:
    value: us-central1
  
  # N1 series required for GPU attachment
  machineType: n1-standard-8
  diskSizeGb: 200
  diskType: pd-ssd
  
  # Autoscaling for cost efficiency
  autoscaling:
    minNodes: 1
    maxNodes: 5
    locationPolicy: ANY  # GPUs may have limited zone availability
  
  # Labels for GPU workload scheduling
  nodeLabels:
    workload-type: gpu
    accelerator: nvidia-tesla-t4
    ml-workload: "true"
    environment: production
  
  # Management settings
  management:
    disableAutoUpgrade: false
    disableAutoRepair: false
```

**Use case:** ML training, inference, video processing, CUDA-enabled applications.

**Workload scheduling:** Use Kubernetes `nodeSelector` or `affinity` to target these nodes:

```yaml
# In your Pod/Deployment spec
nodeSelector:
  workload-type: gpu
  accelerator: nvidia-tesla-t4

resources:
  limits:
    nvidia.com/gpu: 1  # Request 1 GPU
```

**Deploy:**
```bash
openmcf apply -f gpu-pool.yaml
```

---

## 6. Compute-Optimized Pool for CPU-Intensive Workloads

Pool optimized for CPU-bound workloads like batch processing, scientific computing, or high-performance web servers.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeNodePool
metadata:
  name: compute-optimized-pool
  org: my-org
  env: prod
spec:
  clusterProjectId:
    value: prod-gcp-project
  clusterName:
    value: prod-cluster
  clusterLocation:
    value: us-central1
  
  # Compute-optimized machine type (C2 series)
  machineType: c2-standard-16  # 16 vCPU, 64 GB RAM, 3.8 GHz all-core turbo
  diskSizeGb: 100
  diskType: pd-ssd
  
  # Autoscaling for burst capacity
  autoscaling:
    minNodes: 2
    maxNodes: 8
    locationPolicy: BALANCED
  
  # Labels for workload scheduling
  nodeLabels:
    workload-type: compute-intensive
    cpu-tier: high-performance
    environment: production
  
  # Management settings
  management:
    disableAutoUpgrade: false
    disableAutoRepair: false
```

**Use case:** Video encoding, scientific simulations, high-frequency trading, CPU-bound web applications.

**Deploy:**
```bash
openmcf apply -f compute-optimized-pool.yaml
```

---

## Foreign Key References

All examples use **literal values** for cluster references:

```yaml
clusterProjectId:
  value: my-gcp-project
clusterName:
  value: my-cluster
clusterLocation:
  value: us-central1
```

### Using References to Other Resources

If you're managing the cluster with OpenMCF, you can **reference** the cluster resource directly:

```yaml
spec:
  clusterProjectId:
    resource:
      kind: GcpGkeCluster
      name: my-cluster
      fieldPath: spec.project_id
  clusterName:
    resource:
      kind: GcpGkeCluster
      name: my-cluster
      fieldPath: metadata.name
  clusterLocation:
    resource:
      kind: GcpGkeCluster
      name: my-cluster
      fieldPath: spec.location
```

This creates a dependency: the node pool waits for the cluster to be created before provisioning.

---

## Validation and Deployment

### Validate Manifests

```bash
openmcf validate -f pool.yaml
```

### Deploy Node Pool

```bash
openmcf apply -f pool.yaml
```

### Check Status

```bash
openmcf get gcpgkenodepool prod-general
```

### Update Node Pool

Modify the YAML and re-apply:

```bash
openmcf apply -f pool.yaml
```

OpenMCF computes the diff and updates only changed fields. Node configuration changes (machine type, disk type) trigger node recreation using GKE's cordon-and-drain pattern.

### Delete Node Pool

```bash
openmcf delete gcpgkenodepool prod-general
```

**Warning:** This terminates all nodes and evicts all pods in the pool. Ensure workloads are drained or migrated first.

---

## Best Practices Summary

1. **Use multiple specialized pools** instead of one large homogeneous pool
2. **Enable autoscaling** for cost efficiency and burst capacity
3. **Use Spot VMs** (`spot: true`) for fault-tolerant workloads with `location_policy: ANY`
4. **Keep auto-upgrade and auto-repair enabled** for security and reliability
5. **Use custom service accounts** with minimal IAM permissions
6. **Apply node labels** for precise workload scheduling with Kubernetes `nodeSelector`
7. **Set `location_policy: BALANCED`** for HA workloads, `ANY` for Spot VMs
8. **Use SSD disks** (`pd-ssd`) for production workloads requiring low-latency I/O

---

## Related Documentation

- **[Component README](./README.md)**: Overview and quick start
- **[Research Document](./docs/README.md)**: Deep dive into GKE node pools and deployment methods
- **[Pulumi Module](./iac/pulumi/README.md)**: Standalone Pulumi usage
- **[Terraform Module](./iac/tf/README.md)**: Standalone Terraform/OpenTofu usage


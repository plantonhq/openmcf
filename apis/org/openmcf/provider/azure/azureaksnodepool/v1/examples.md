# Azure AKS Node Pool Examples

This document provides comprehensive examples for deploying Azure Kubernetes Service (AKS) node pools using the AzureAksNodePool API resource.

## Table of Contents

- [Minimal Example](#minimal-example)
- [Production User Pool](#production-user-pool)
- [Additional System Pool](#additional-system-pool)
- [Spot Instance Pool](#spot-instance-pool)
- [Windows Node Pool](#windows-node-pool)
- [GPU Node Pool](#gpu-node-pool)
- [Memory-Optimized Pool](#memory-optimized-pool)

---

## Minimal Example

The simplest node pool with only required fields.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureAksNodePool
metadata:
  name: basic-pool
spec:
  clusterName:
    value: my-aks-cluster
  vmSize: Standard_D4s_v3
  initialNodeCount: 2
  availabilityZones:
    - "1"
    - "2"
```

**Deploy:**
```shell
planton apply -f basic-pool.yaml
```

**What you get:**
- User mode node pool (default)
- Linux OS (default)
- 2 nodes
- Regular (non-Spot) instances
- No autoscaling (fixed size)
- Spread across 2 availability zones

---

## Production User Pool

Production-ready user pool with autoscaling and multi-AZ.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureAksNodePool
metadata:
  name: production-app-pool
  labels:
    environment: production
    team: platform
spec:
  clusterName:
    value: production-aks-cluster
  
  vmSize: Standard_D8s_v3  # 8 vCPU, 32 GB RAM
  initialNodeCount: 2
  
  autoscaling:
    minNodes: 2
    maxNodes: 10
  
  availabilityZones:
    - "1"
    - "2"
    - "3"
  
  mode: USER
  osType: LINUX
  spotEnabled: false
```

**Deploy:**
```shell
planton apply -f production-pool.yaml
```

**What you get:**
- User mode pool for application workloads
- Autoscaling from 2 to 10 nodes
- Multi-AZ (3 zones) for high availability
- D8s_v3 VMs (8 vCPU, 32 GB RAM)
- Regular instances for reliability

---

## Additional System Pool

Add a secondary system pool to an existing cluster.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureAksNodePool
metadata:
  name: system-pool-secondary
spec:
  clusterName:
    value: production-aks-cluster
  
  vmSize: Standard_D4s_v5  # 4 vCPU, 16 GB RAM
  initialNodeCount: 3
  
  availabilityZones:
    - "1"
    - "2"
    - "3"
  
  mode: SYSTEM  # System mode for critical components
  osType: LINUX  # System pools must be Linux
  spotEnabled: false  # Spot not allowed for system pools
```

**Deploy:**
```shell
planton apply -f system-pool.yaml
```

**What you get:**
- System mode pool for cluster components
- Fixed size (3 nodes, no autoscaling)
- Multi-AZ for reliability
- Linux only (required for system pools)
- Regular instances (Spot not supported for system mode)

**Note:** System pools cannot use Spot instances and cannot scale to zero.

---

## Spot Instance Pool

Cost-optimized pool for fault-tolerant workloads.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureAksNodePool
metadata:
  name: spot-batch-pool
spec:
  clusterName:
    value: production-aks-cluster
  
  vmSize: Standard_D8s_v3
  initialNodeCount: 1  # Start with 1, can scale to 0
  
  autoscaling:
    minNodes: 0  # Allow scale to zero
    maxNodes: 20
  
  availabilityZones:
    - "1"
    - "2"
    - "3"
  
  mode: USER
  osType: LINUX
  spotEnabled: true  # Enable Spot for 30-90% cost savings
```

**Deploy:**
```shell
planton apply -f spot-pool.yaml
```

**What you get:**
- User mode pool with Spot instances
- 30-90% cost savings vs regular VMs
- Can scale to zero when idle
- Nodes can be evicted with 30 seconds notice
- Ideal for batch jobs, CI/CD, dev/test

**Best for:**
- Batch processing jobs
- CI/CD workloads
- Development/staging environments
- Stateless, fault-tolerant applications

**Not suitable for:**
- Critical production services
- Stateful applications
- Workloads requiring guaranteed availability

---

## Windows Node Pool

Windows container support for legacy .NET applications.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureAksNodePool
metadata:
  name: windows-pool
spec:
  clusterName:
    value: production-aks-cluster
  
  vmSize: Standard_D4s_v3
  initialNodeCount: 2
  
  autoscaling:
    minNodes: 2
    maxNodes: 5
  
  availabilityZones:
    - "1"
    - "2"
  
  mode: USER  # Windows pools must be User mode
  osType: WINDOWS
  spotEnabled: false
```

**Deploy:**
```shell
planton apply -f windows-pool.yaml
```

**What you get:**
- Windows Server nodes
- Support for Windows containers
- User mode only (Windows cannot be System)
- Higher resource overhead than Linux

**Requirements:**
- Cluster must use Azure CNI networking
- Windows Server 2019 or 2022
- Larger VM sizes recommended (≥4 vCPU)

**Use cases:**
- Legacy .NET Framework applications
- Windows-specific dependencies
- Migration from Windows Server

---

## GPU Node Pool

GPU-accelerated pool for AI/ML workloads.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureAksNodePool
metadata:
  name: gpu-ml-pool
spec:
  clusterName:
    value: production-aks-cluster
  
  vmSize: Standard_NC6s_v3  # 6 vCPU, 112 GB RAM, 1x V100 GPU
  initialNodeCount: 1
  
  autoscaling:
    minNodes: 0  # Scale to zero when not in use
    maxNodes: 5
  
  availabilityZones:
    - "1"
  
  mode: USER
  osType: LINUX
  spotEnabled: false  # Or true for cost savings (but risk eviction during training)
```

**Deploy:**
```shell
planton apply -f gpu-pool.yaml
```

**What you get:**
- NVIDIA V100 GPUs
- Can scale to zero to save costs
- Ideal for ML training and inference

**GPU VM sizes:**
- **Standard_NC6s_v3**: 1x V100 GPU (16GB) - training
- **Standard_NC12s_v3**: 2x V100 GPUs - distributed training
- **Standard_ND40rs_v2**: 8x V100 GPUs - large-scale training
- **Standard_NC4as_T4_v3**: 1x T4 GPU - cost-effective inference

**Note:** GPU VMs are expensive. Enable autoscaling with `minNodes: 0` to scale to zero when idle.

---

## Memory-Optimized Pool

High-memory pool for caching and in-memory databases.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureAksNodePool
metadata:
  name: memory-pool
spec:
  clusterName:
    value: production-aks-cluster
  
  vmSize: Standard_E8s_v5  # 8 vCPU, 64 GB RAM
  initialNodeCount: 2
  
  autoscaling:
    minNodes: 2
    maxNodes: 8
  
  availabilityZones:
    - "1"
    - "2"
    - "3"
  
  mode: USER
  osType: LINUX
  spotEnabled: false
```

**Deploy:**
```shell
planton apply -f memory-pool.yaml
```

**What you get:**
- High memory-to-CPU ratio (8 GB RAM per vCPU)
- E-series VMs optimized for memory-intensive workloads
- Multi-AZ for HA

**Memory-optimized VM sizes:**
- **Standard_E4s_v5**: 4 vCPU, 32 GB RAM
- **Standard_E8s_v5**: 8 vCPU, 64 GB RAM
- **Standard_E16s_v5**: 16 vCPU, 128 GB RAM
- **Standard_E32s_v5**: 32 vCPU, 256 GB RAM

**Use cases:**
- Redis, Memcached (caching)
- In-memory databases
- Large heap Java/JVM applications
- Data analytics workloads

---

## Deployment Tips

### Verify Deployment

After deployment, check the node pool status:

```shell
kubectl get nodes
kubectl get nodes -l agentpool=<pool-name>
kubectl describe node <node-name>
```

### Scale Node Pool Manually

If autoscaling is disabled, scale manually:

```shell
az aks nodepool scale \
  --resource-group rg-<cluster-name> \
  --cluster-name <cluster-name> \
  --name <pool-name> \
  --node-count <count>
```

### Upgrade Node Pool

Upgrade Kubernetes version:

```shell
az aks nodepool upgrade \
  --resource-group rg-<cluster-name> \
  --cluster-name <cluster-name> \
  --name <pool-name> \
  --kubernetes-version 1.30
```

### Delete Node Pool

Remove a node pool (drains nodes gracefully):

```shell
planton delete -f pool.yaml
```

Or via Azure CLI:

```shell
az aks nodepool delete \
  --resource-group rg-<cluster-name> \
  --cluster-name <cluster-name> \
  --name <pool-name>
```

---

## Common Configurations

### Fixed-Size Pool (No Autoscaling)

```yaml
spec:
  initialNodeCount: 3
  # No autoscaling field = fixed size
```

### Autoscaling Pool

```yaml
spec:
  initialNodeCount: 2
  autoscaling:
    minNodes: 2
    maxNodes: 10
```

### Scale-to-Zero Pool (User Mode Only)

```yaml
spec:
  initialNodeCount: 1
  autoscaling:
    minNodes: 0  # Can scale to zero
    maxNodes: 20
  mode: USER  # Only User pools can scale to zero
```

### Single-Zone Pool (Dev/Test)

```yaml
spec:
  availabilityZones:
    - "1"  # Single zone acceptable for dev/test
```

### Multi-Zone Pool (Production)

```yaml
spec:
  availabilityZones:
    - "1"
    - "2"
    - "3"  # 3 zones for production HA
```

---

## Best Practices

### Production Recommendations

✅ **Use multi-AZ for production**
```yaml
availabilityZones: ["1", "2", "3"]
```

✅ **Enable autoscaling for user pools**
```yaml
autoscaling:
  minNodes: 2  # Baseline capacity
  maxNodes: 10  # Peak capacity
```

✅ **Separate pools by workload type**
- General-purpose: Standard_D8s_v3
- Compute-intensive: Standard_F16s_v2
- Memory-intensive: Standard_E8s_v5
- Batch/dev: Spot instances

✅ **Use Spot for non-critical workloads**
```yaml
spotEnabled: true
autoscaling:
  minNodes: 0
  maxNodes: 20
```

### Anti-Patterns

❌ **Don't mix critical and non-critical workloads**
- Use separate pools with different priorities

❌ **Don't use Spot for production databases**
- Spot instances can be evicted with 30 seconds notice

❌ **Don't forget availability zones in production**
- Single-zone deployments are single points of failure

❌ **Don't use Windows for system pools**
- System pools must be Linux

---

## Next Steps

After deploying node pools:

1. **Label nodes**: Use Kubernetes labels to route workloads
2. **Set taints**: Prevent unwanted pods from scheduling
3. **Configure pod disruption budgets**: Protect during upgrades
4. **Monitor node health**: Use Container Insights
5. **Optimize costs**: Use Spot for appropriate workloads

## Support

For issues or questions:
- Check the [main README](./README.md) for component overview
- Review the [research documentation](./docs/README.md) for architecture details
- Consult the [Azure AKS Documentation](https://docs.microsoft.com/en-us/azure/aks/)


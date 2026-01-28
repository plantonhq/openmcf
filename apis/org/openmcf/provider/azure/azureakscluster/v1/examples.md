# Azure AKS Cluster Examples

This document provides comprehensive examples for deploying production-ready Azure Kubernetes Service (AKS) clusters using the AzureAksCluster API resource.

## Table of Contents

- [Minimal Example](#minimal-example)
- [Production Example](#production-example)
- [Private Cluster Example](#private-cluster-example)
- [Development/Test Example](#developmenttest-example)
- [Multi-Node Pool Example](#multi-node-pool-example)
- [With Spot Instances Example](#with-spot-instances-example)
- [Using Foreign Key Reference](#using-foreign-key-reference)

---

## Minimal Example

The simplest production-ready AKS cluster with required fields.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureAksCluster
metadata:
  name: minimal-aks-cluster
spec:
  region: eastus
  vnetSubnetId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-network/providers/Microsoft.Network/virtualNetworks/app-vnet/subnets/aks-nodes
  
  # Required: System node pool configuration
  systemNodePool:
    vmSize: Standard_D4s_v5
    autoscaling:
      minCount: 3
      maxCount: 5
    availabilityZones:
      - "1"
      - "2"
      - "3"
```

**Deploy:**
```shell
planton apply -f minimal-example.yaml
```

**What you get:**
- Standard tier control plane (99.95% SLA with AZs)
- Azure CNI with Overlay mode (default)
- Kubernetes version 1.30 (default)
- System node pool across 3 availability zones
- Azure AD RBAC enabled
- System-Assigned Managed Identity

---

## Production Example

Complete production configuration with all recommended features.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureAksCluster
metadata:
  name: prod-app-aks
  labels:
    environment: production
    team: platform
spec:
  # Region and networking
  region: eastus
  vnetSubnetId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-network/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/aks-nodes
  
  # Pin Kubernetes version for production stability
  kubernetesVersion: "1.30"
  
  # Standard tier provides financially-backed 99.95% uptime SLA
  controlPlaneSku: STANDARD
  
  # Use Azure CNI with Overlay mode (modern, IP-efficient)
  networkPlugin: AZURE_CNI
  networkPluginMode: OVERLAY
  
  # API server access - restrict to known networks
  privateClusterEnabled: false
  authorizedIpRanges:
    - "203.0.113.0/24"   # Corporate office network
    - "198.51.100.0/24"  # CI/CD agents
    - "192.0.2.0/24"     # Operations team VPN
  
  # Azure AD RBAC for centralized access control
  disableAzureAdRbac: false
  
  # System node pool (critical components)
  systemNodePool:
    vmSize: Standard_D4s_v5  # 4 vCPU, 16 GB RAM
    autoscaling:
      minCount: 3  # HA requirement
      maxCount: 5
    availabilityZones:
      - "1"
      - "2"
      - "3"  # Multi-AZ for 99.95% SLA
  
  # User node pools (application workloads)
  userNodePools:
    - name: general
      vmSize: Standard_D8s_v5  # 8 vCPU, 32 GB RAM
      autoscaling:
        minCount: 2
        maxCount: 10
      availabilityZones:
        - "1"
        - "2"
        - "3"
      spotEnabled: false
  
  # Production add-ons
  addons:
    enableContainerInsights: true
    logAnalyticsWorkspaceId: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-monitoring/providers/Microsoft.OperationalInsights/workspaces/prod-logs
    enableKeyVaultCsiDriver: true
    enableAzurePolicy: true
    enableWorkloadIdentity: true
```

**Deploy:**
```shell
planton apply -f production-example.yaml
```

**What you get:**
- Standard tier control plane with 99.95% SLA
- Multi-AZ deployment across 3 availability zones
- Separate system and user node pools
- Restricted API server access
- Azure AD RBAC integration
- Container Insights monitoring
- Key Vault CSI driver for secrets
- Azure Policy for governance
- Workload Identity for secret-less authentication

---

## Private Cluster Example

Fully private AKS cluster with no public API endpoint.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureAksCluster
metadata:
  name: private-aks-cluster
spec:
  region: westus2
  vnetSubnetId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-network/providers/Microsoft.Network/virtualNetworks/secure-vnet/subnets/aks-private
  
  kubernetesVersion: "1.30"
  controlPlaneSku: STANDARD
  networkPlugin: AZURE_CNI
  networkPluginMode: OVERLAY
  
  # Private cluster - no public API endpoint
  privateClusterEnabled: true
  # No authorized IP ranges needed (no public endpoint)
  authorizedIpRanges: []
  
  disableAzureAdRbac: false
  
  # System node pool
  systemNodePool:
    vmSize: Standard_D4s_v5
    autoscaling:
      minCount: 3
      maxCount: 5
    availabilityZones: ["1", "2", "3"]
  
  # User node pool
  userNodePools:
    - name: apps
      vmSize: Standard_D8s_v5
      autoscaling:
        minCount: 2
        maxCount: 8
      availabilityZones: ["1", "2", "3"]
      spotEnabled: false
  
  # Add-ons
  addons:
    enableContainerInsights: true
    logAnalyticsWorkspaceId: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-monitoring/providers/Microsoft.OperationalInsights/workspaces/private-logs
    enableKeyVaultCsiDriver: true
    enableAzurePolicy: true
    enableWorkloadIdentity: true
```

**Deploy:**
```shell
planton apply -f private-example.yaml
```

**What you get:**
- Private API server endpoint (accessible only from VNet)
- Ideal for highly secure environments
- Requires VPN or bastion host for kubectl access
- All production features enabled

**Access:**
```shell
# From a VM in the same VNet or via VPN
az aks get-credentials --resource-group rg-private-aks-cluster --name private-aks-cluster
kubectl get nodes
```

---

## Development/Test Example

Cost-optimized configuration for development and testing.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureAksCluster
metadata:
  name: dev-aks-cluster
spec:
  region: eastus
  vnetSubnetId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-dev-network/providers/Microsoft.Network/virtualNetworks/dev-vnet/subnets/aks-dev
  
  kubernetesVersion: "1.30"
  
  # FREE tier - no SLA, saves $73/month (dev/test only)
  controlPlaneSku: FREE
  
  networkPlugin: AZURE_CNI
  networkPluginMode: OVERLAY
  
  # Public cluster for easier access during development
  privateClusterEnabled: false
  # No IP restrictions (not recommended for production!)
  authorizedIpRanges: []
  
  # Disable Azure AD RBAC for simpler dev setup
  disableAzureAdRbac: true
  
  # Smaller, single-AZ system node pool
  systemNodePool:
    vmSize: Standard_D2s_v3  # 2 vCPU, 8 GB RAM - smaller for cost
    autoscaling:
      minCount: 1  # Single node acceptable for dev
      maxCount: 2
    availabilityZones: ["1"]  # Single AZ for cost savings
  
  # Optional: Add user node pool if needed
  userNodePools:
    - name: dev
      vmSize: Standard_D4s_v5
      autoscaling:
        minCount: 1
        maxCount: 3
      availabilityZones: ["1"]
      spotEnabled: true  # Use spot for 70% cost savings
  
  # Minimal add-ons (skip monitoring to reduce costs)
  addons:
    enableContainerInsights: false
    logAnalyticsWorkspaceId: ""
    enableKeyVaultCsiDriver: false
    enableAzurePolicy: false
    enableWorkloadIdentity: false
```

**Deploy:**
```shell
planton apply -f dev-example.yaml
```

**What you get:**
- Free tier control plane (no SLA)
- Single-AZ deployment for cost savings
- Smaller VM sizes
- Spot instances for user workloads
- Public API access
- No monitoring (cost optimization)

**⚠️ Warning:** This configuration is for development only. Never use in production!

---

## Multi-Node Pool Example

Multiple specialized node pools for different workload types.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureAksCluster
metadata:
  name: multi-pool-aks
spec:
  region: eastus
  vnetSubnetId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-network/providers/Microsoft.Network/virtualNetworks/app-vnet/subnets/aks-nodes
  
  kubernetesVersion: "1.30"
  controlPlaneSku: STANDARD
  networkPlugin: AZURE_CNI
  networkPluginMode: OVERLAY
  privateClusterEnabled: false
  authorizedIpRanges: ["203.0.113.0/24"]
  disableAzureAdRbac: false
  
  # System node pool
  systemNodePool:
    vmSize: Standard_D4s_v5
    autoscaling:
      minCount: 3
      maxCount: 5
    availabilityZones: ["1", "2", "3"]
  
  # Multiple specialized user node pools
  userNodePools:
    # General-purpose applications
    - name: general
      vmSize: Standard_D8s_v5  # 8 vCPU, 32 GB RAM
      autoscaling:
        minCount: 2
        maxCount: 10
      availabilityZones: ["1", "2", "3"]
      spotEnabled: false
    
    # Compute-intensive workloads
    - name: compute
      vmSize: Standard_F16s_v2  # 16 vCPU, 32 GB RAM (compute-optimized)
      autoscaling:
        minCount: 0  # Scale to zero when not needed
        maxCount: 5
      availabilityZones: ["1", "2", "3"]
      spotEnabled: false
    
    # Memory-intensive workloads
    - name: memory
      vmSize: Standard_E8s_v5  # 8 vCPU, 64 GB RAM (memory-optimized)
      autoscaling:
        minCount: 1
        maxCount: 5
      availabilityZones: ["1", "2", "3"]
      spotEnabled: false
    
    # Batch/background jobs (spot instances)
    - name: batch
      vmSize: Standard_D8s_v5
      autoscaling:
        minCount: 0
        maxCount: 20
      availabilityZones: ["1", "2", "3"]
      spotEnabled: true  # 30-90% cost savings, can be evicted
  
  # Add-ons
  addons:
    enableContainerInsights: true
    logAnalyticsWorkspaceId: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-monitoring/providers/Microsoft.OperationalInsights/workspaces/logs
    enableKeyVaultCsiDriver: true
    enableAzurePolicy: true
    enableWorkloadIdentity: true
```

**Deploy:**
```shell
planton apply -f multi-pool-example.yaml
```

**What you get:**
- Four specialized node pools for different workload types
- General pool for standard applications
- Compute pool for CPU-intensive workloads
- Memory pool for memory-intensive workloads
- Batch pool with spot instances for cost-effective background jobs
- Ability to scale pools independently

**Use with pod selectors:**
```yaml
# Deploy pod to compute pool
apiVersion: v1
kind: Pod
metadata:
  name: compute-intensive-app
spec:
  nodeSelector:
    pool-name: compute
  containers:
    - name: app
      image: my-compute-app:latest
```

---

## With Spot Instances Example

Cost-optimized cluster using spot instances for fault-tolerant workloads.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureAksCluster
metadata:
  name: spot-aks-cluster
spec:
  region: eastus
  vnetSubnetId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-network/providers/Microsoft.Network/virtualNetworks/app-vnet/subnets/aks-nodes
  
  kubernetesVersion: "1.30"
  controlPlaneSku: STANDARD
  networkPlugin: AZURE_CNI
  networkPluginMode: OVERLAY
  disableAzureAdRbac: false
  
  # System node pool (regular instances for reliability)
  systemNodePool:
    vmSize: Standard_D4s_v5
    autoscaling:
      minCount: 3
      maxCount: 5
    availabilityZones: ["1", "2", "3"]
  
  # User node pools with spot instances
  userNodePools:
    # Regular pool for critical apps
    - name: critical
      vmSize: Standard_D8s_v5
      autoscaling:
        minCount: 2
        maxCount: 5
      availabilityZones: ["1", "2", "3"]
      spotEnabled: false  # Regular instances
    
    # Spot pool for fault-tolerant workloads
    - name: spot
      vmSize: Standard_D8s_v5
      autoscaling:
        minCount: 0
        maxCount: 30  # Scale up spot instances aggressively
      availabilityZones: ["1", "2", "3"]
      spotEnabled: true  # Spot instances - 30-90% savings
  
  addons:
    enableContainerInsights: true
    logAnalyticsWorkspaceId: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-monitoring/providers/Microsoft.OperationalInsights/workspaces/logs
    enableKeyVaultCsiDriver: true
    enableAzurePolicy: true
    enableWorkloadIdentity: true
```

**Deploy:**
```shell
planton apply -f spot-example.yaml
```

**What you get:**
- Mixed pool strategy: regular + spot instances
- 30-90% cost savings on spot workloads
- Critical apps on regular instances
- Fault-tolerant apps on spot instances

**Best practices for spot:**
- Use for stateless, fault-tolerant workloads
- Implement graceful shutdown (handle eviction notices)
- Use pod disruption budgets
- Design for horizontal scaling

---

## Using Foreign Key Reference

Reference an AzureVpc resource to automatically resolve the subnet ID.

```yaml
# First, create the VNet
apiVersion: azure.openmcf.org/v1
kind: AzureVpc
metadata:
  name: my-azure-vnet
spec:
  region: eastus
  cidrBlock: 10.0.0.0/16
  # ... other VNet configuration

---

# Then create AKS cluster referencing the VNet
apiVersion: azure.openmcf.org/v1
kind: AzureAksCluster
metadata:
  name: aks-with-ref
spec:
  region: eastus
  
  # Reference the VNet's nodes subnet output
  vnetSubnetId:
    ref:
      kind: AzureVpc
      name: my-azure-vnet
      path: status.outputs.nodes_subnet_id
  
  kubernetesVersion: "1.30"
  controlPlaneSku: STANDARD
  networkPlugin: AZURE_CNI
  networkPluginMode: OVERLAY
  disableAzureAdRbac: false
  
  systemNodePool:
    vmSize: Standard_D4s_v5
    autoscaling:
      minCount: 3
      maxCount: 5
    availabilityZones: ["1", "2", "3"]
  
  addons:
    enableContainerInsights: false
    enableKeyVaultCsiDriver: true
    enableAzurePolicy: true
    enableWorkloadIdentity: true
```

**Deploy:**
```shell
# Deploy both resources
planton apply -f azure-vpc.yaml
planton apply -f aks-with-ref.yaml
```

**What you get:**
- Automatic subnet ID resolution
- Type-safe foreign key references
- Simplified configuration management
- Clear dependencies between resources

---

## Deployment Tips

### Getting Credentials

After deployment, get cluster credentials:

```shell
az aks get-credentials --resource-group rg-<cluster-name> --name <cluster-name>
```

### Verify Deployment

Check cluster status:

```shell
kubectl get nodes
kubectl get pods -n kube-system
kubectl cluster-info
```

### Node Pool Labels

System node pool automatically gets:
- `node-role=system`
- `only_critical_addons_enabled=true`

User node pools get:
- `node-role=user`
- `pool-name=<pool-name>`

Use these labels with pod nodeSelectors.

### Common Configurations

**Production Checklist:**
- ✅ `controlPlaneSku`: `STANDARD` (99.95% SLA)
- ✅ `networkPluginMode`: `OVERLAY` (modern, IP-efficient)
- ✅ `kubernetesVersion`: Pinned to specific version
- ✅ `systemNodePool`: 3+ nodes across 3 AZs
- ✅ `authorizedIpRanges`: Configured or private cluster enabled
- ✅ `disableAzureAdRbac`: `false` (Azure AD RBAC enabled)
- ✅ `addons.enableContainerInsights`: `true`
- ✅ `addons.enableWorkloadIdentity`: `true`

**Development Configuration:**
- ⚠️ `controlPlaneSku`: `FREE` (no SLA, cost savings)
- ⚠️ Single AZ deployment acceptable
- ⚠️ Smaller VM sizes (Standard_D2s_v3)
- ⚠️ Public access without IP restrictions
- ⚠️ Monitoring disabled for cost savings

---

## Advanced Networking Example

Custom networking configuration with specific CIDRs.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureAksCluster
metadata:
  name: custom-network-aks
spec:
  region: eastus
  vnetSubnetId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-network/providers/Microsoft.Network/virtualNetworks/app-vnet/subnets/aks-nodes
  
  kubernetesVersion: "1.30"
  controlPlaneSku: STANDARD
  networkPlugin: AZURE_CNI
  networkPluginMode: OVERLAY
  disableAzureAdRbac: false
  
  systemNodePool:
    vmSize: Standard_D4s_v5
    autoscaling:
      minCount: 3
      maxCount: 5
    availabilityZones: ["1", "2", "3"]
  
  # Custom networking CIDRs
  advancedNetworking:
    podCidr: "10.244.0.0/16"      # Pod CIDR for Overlay mode
    serviceCidr: "10.10.0.0/16"   # Custom service CIDR
    dnsServiceIp: "10.10.0.10"    # Must be within serviceCidr
    customDnsServers:              # Optional custom DNS
      - "8.8.8.8"
      - "8.8.4.4"
```

**Use when:**
- You need custom network CIDRs to avoid conflicts
- You have specific DNS server requirements
- You're integrating with existing network infrastructure

---

## Next Steps

After deploying your AKS cluster:

1. **Configure kubectl**: Get credentials and verify connectivity
2. **Configure node pools**: Use labels/taints to route workloads appropriately
3. **Install cluster add-ons**: Deploy ingress controller, cert-manager, etc.
4. **Set up RBAC**: Assign Azure AD users/groups to Kubernetes roles
5. **Deploy workloads**: Start deploying your applications
6. **Configure monitoring**: Set up alerts and dashboards in Azure Monitor
7. **Implement pod security**: Use Azure Policy or Pod Security Standards

## Support

For issues or questions:
- Check the [main README](./README.md) for component overview
- Review the [research documentation](./docs/README.md) for architecture details
- Consult the [Azure AKS Documentation](https://docs.microsoft.com/en-us/azure/aks/)

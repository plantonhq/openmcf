# DigitalOcean Kubernetes Cluster Examples

This document provides practical examples for deploying DigitalOcean Kubernetes (DOKS) clusters using the OpenMCF API. Each example is a complete, deployable manifest that demonstrates different use cases and configuration patterns.

## Table of Contents

- [Example 1: Minimal Development Cluster](#example-1-minimal-development-cluster)
- [Example 2: Staging Cluster with Autoscaling](#example-2-staging-cluster-with-autoscaling)
- [Example 3: Production Cluster with High Availability](#example-3-production-cluster-with-high-availability)
- [Example 4: Secure Production Cluster with Firewall](#example-4-secure-production-cluster-with-firewall)

---

## Example 1: Minimal Development Cluster

**Use Case:** Quick testing environment for local development. Cost optimization is paramount.

**Configuration:**
- **Region:** nyc1 (New York)
- **Nodes:** 2 nodes
- **Instance Size:** s-1vcpu-2gb (~$12/month per node = ~$24/month total)
- **HA Control Plane:** false (saves $40/month)
- **Auto-Upgrade:** true (accept patches automatically)
- **Autoscaling:** false (fixed 2-node size)

**Monthly Cost:** ~$24 (nodes only, control plane free)

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanKubernetesCluster
metadata:
  name: dev-cluster
spec:
  clusterName: dev-cluster
  region: nyc1
  kubernetesVersion: "1.29"
  vpc:
    value: "your-vpc-id-here"  # Replace with your VPC UUID
  highlyAvailable: false
  autoUpgrade: true
  disableSurgeUpgrade: false
  tags:
    - env:dev
    - team:engineering
  defaultNodePool:
    size: s-1vcpu-2gb
    nodeCount: 2
    autoScale: false
    minNodes: 0
    maxNodes: 0
```

**Deployment:**

```bash
# Save the manifest to a file
cat > dev-cluster.yaml << 'EOF'
[paste manifest above]
EOF

# Deploy using OpenMCF CLI
openmcf apply -f dev-cluster.yaml
```

---

## Example 2: Staging Cluster with Autoscaling

**Use Case:** Multi-node cluster mirroring production configuration for integration testing.

**Configuration:**
- **Region:** sfo3 (San Francisco)
- **Nodes:** 3 initial nodes (minimum for HA simulation)
- **Instance Size:** s-2vcpu-4gb (~$20/month per node)
- **HA Control Plane:** false (acceptable downtime for staging)
- **Auto-Upgrade:** true
- **Autoscaling:** true (min: 3, max: 6)
- **Maintenance Window:** Sunday 02:00 UTC

**Monthly Cost:** ~$60 baseline (nodes only), scaling up to ~$120 at max capacity

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanKubernetesCluster
metadata:
  name: staging-cluster
spec:
  clusterName: staging-cluster
  region: sfo3
  kubernetesVersion: "1.29"
  vpc:
    value: "your-vpc-id-here"  # Replace with your VPC UUID
  highlyAvailable: false
  autoUpgrade: true
  disableSurgeUpgrade: false
  maintenanceWindow: "sunday=02:00"
  tags:
    - env:staging
    - team:qa
    - project:core-platform
  defaultNodePool:
    size: s-2vcpu-4gb
    nodeCount: 3
    autoScale: true
    minNodes: 3
    maxNodes: 6
```

**Deployment:**

```bash
# Save the manifest
cat > staging-cluster.yaml << 'EOF'
[paste manifest above]
EOF

# Deploy
openmcf apply -f staging-cluster.yaml

# Verify cluster status
openmcf get digitaloceankubernetescluster staging-cluster

# Get kubeconfig
openmcf kubeconfig digitaloceankubernetescluster staging-cluster > ~/.kube/staging-config
export KUBECONFIG=~/.kube/staging-config

# Verify cluster connectivity
kubectl get nodes
```

---

## Example 3: Production Cluster with High Availability

**Use Case:** Robust cluster for production traffic with HA control plane and autoscaling.

**Configuration:**
- **Region:** nyc1 (New York)
- **Nodes:** 5 initial nodes, autoscaling enabled
- **Instance Size:** s-4vcpu-8gb (~$43/month per node)
- **HA Control Plane:** true (critical for uptime)
- **Auto-Upgrade:** true (with maintenance window set)
- **Autoscaling:** true (min: 5, max: 10)
- **Registry Integration:** true (DOCR)
- **Maintenance Window:** Sunday 02:00 UTC

**Monthly Cost (baseline):**
- Control plane (HA): $40 (may be waived with 5+ nodes)
- Nodes: 5 × $43 = $215
- **Total:** ~$255/month (before autoscaling)

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanKubernetesCluster
metadata:
  name: prod-cluster
spec:
  clusterName: prod-cluster
  region: nyc1
  kubernetesVersion: "1.29"
  vpc:
    value: "your-vpc-id-here"  # Replace with your production VPC UUID
  highlyAvailable: true
  autoUpgrade: true
  disableSurgeUpgrade: false
  maintenanceWindow: "sunday=02:00"
  registryIntegration: true
  tags:
    - env:production
    - criticality:high
    - team:platform
    - cost-center:engineering
  defaultNodePool:
    size: s-4vcpu-8gb
    nodeCount: 5
    autoScale: true
    minNodes: 5
    maxNodes: 10
```

**Deployment with Validation:**

```bash
# Save the manifest
cat > prod-cluster.yaml << 'EOF'
[paste manifest above]
EOF

# Validate manifest before applying
openmcf validate -f prod-cluster.yaml

# Deploy (will take ~3-5 minutes)
openmcf apply -f prod-cluster.yaml

# Monitor deployment progress
openmcf status digitaloceankubernetescluster prod-cluster

# Once deployed, get kubeconfig
openmcf kubeconfig digitaloceankubernetescluster prod-cluster > ~/.kube/prod-config
export KUBECONFIG=~/.kube/prod-config

# Verify cluster health
kubectl get nodes
kubectl get pods -A

# Check cluster version
kubectl version --short
```

---

## Example 4: Secure Production Cluster with Firewall

**Use Case:** Maximum security production cluster with control plane firewall, VPC isolation, and all production features.

**Configuration:**
- **Region:** sfo3 (San Francisco)
- **Nodes:** 5 initial nodes, autoscaling to 10
- **Instance Size:** s-4vcpu-8gb
- **HA Control Plane:** true
- **Control Plane Firewall:** Restricted to office and VPN IPs
- **Auto-Upgrade:** true
- **Registry Integration:** true
- **Maintenance Window:** Sunday 03:00 UTC

**Security Features:**
- Control plane API access restricted to whitelisted IPs
- VPC isolation for cluster nodes
- Container registry integration for secure image pulls
- Automatic security patch upgrades

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanKubernetesCluster
metadata:
  name: secure-prod-cluster
spec:
  clusterName: secure-prod-cluster
  region: sfo3
  kubernetesVersion: "1.29"
  vpc:
    value: "your-vpc-id-here"  # Replace with your production VPC UUID
  highlyAvailable: true
  autoUpgrade: true
  disableSurgeUpgrade: false
  maintenanceWindow: "sunday=03:00"
  registryIntegration: true
  controlPlaneFirewallAllowedIps:
    - "203.0.113.10/32"     # Office IP
    - "203.0.113.20/32"     # VPN endpoint
    - "198.51.100.0/24"     # CI/CD runner subnet
  tags:
    - env:production
    - security:restricted
    - compliance:required
    - team:platform
  defaultNodePool:
    size: s-4vcpu-8gb
    nodeCount: 5
    autoScale: true
    minNodes: 5
    maxNodes: 10
```

**Deployment with Security Verification:**

```bash
# Save the manifest
cat > secure-prod-cluster.yaml << 'EOF'
[paste manifest above]
EOF

# Validate manifest
openmcf validate -f secure-prod-cluster.yaml

# Deploy
openmcf apply -f secure-prod-cluster.yaml

# Get kubeconfig
openmcf kubeconfig digitaloceankubernetescluster secure-prod-cluster > ~/.kube/secure-prod-config

# Note: You must be connecting from one of the whitelisted IPs to access the cluster
export KUBECONFIG=~/.kube/secure-prod-config

# Verify access (will fail if not from whitelisted IP)
kubectl get nodes

# Deploy network policies for pod-to-pod security
kubectl apply -f - <<EOF
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-ingress
  namespace: default
spec:
  podSelector: {}
  policyTypes:
  - Ingress
EOF

# Verify firewall is working (test from non-whitelisted IP should fail)
```

---

## Common Operations

### Get Cluster Status

```bash
# Get detailed cluster information
openmcf get digitaloceankubernetescluster <cluster-name> -o yaml

# Get cluster outputs (ID, kubeconfig, API endpoint)
openmcf outputs digitaloceankubernetescluster <cluster-name>
```

### Update Cluster Configuration

```bash
# Edit the manifest file to change configuration
vim my-cluster.yaml

# Apply changes (e.g., increase max_nodes for autoscaling)
openmcf apply -f my-cluster.yaml
```

### Scale Node Pool Manually

```bash
# Edit manifest to change nodeCount
# Example: Change nodeCount from 3 to 5
vim my-cluster.yaml

# Apply
openmcf apply -f my-cluster.yaml
```

### Delete Cluster

```bash
# Delete cluster (will destroy all resources)
openmcf delete -f my-cluster.yaml

# Or delete by name
openmcf delete digitaloceankubernetescluster <cluster-name>
```

---

## Best Practices

### Cost Optimization

1. **Dev Clusters:** Use small node sizes (s-1vcpu-2gb) and disable HA
2. **Staging Clusters:** Use moderate sizes (s-2vcpu-4gb) without HA
3. **Production Clusters:** Enable HA only for business-critical workloads
4. **Autoscaling:** Set appropriate min/max bounds to prevent runaway costs
5. **Teardown:** Delete dev/staging clusters outside working hours to save costs

### Security

1. **Always use control plane firewall in production** - Restrict API access to known IPs
2. **Enable VPC isolation** - Run clusters in dedicated VPCs
3. **Use registry integration** - Securely pull private images from DOCR
4. **Apply NetworkPolicies** - Implement least-privilege pod-to-pod communication
5. **Enable auto-upgrade** - Receive security patches automatically

### High Availability

1. **Enable HA control plane for production** - Multiple masters for redundancy
2. **Use 3+ nodes minimum** - Survive single node failures
3. **Set maintenance windows** - Control when updates occur
4. **Enable surge upgrades** - Maintain capacity during updates
5. **Distribute across regions** - Consider multi-cluster deployments for critical apps

### Monitoring

1. **Deploy Prometheus + Grafana** - Comprehensive cluster monitoring
2. **Enable Hubble** - Network flow observability (included with Cilium CNI)
3. **Configure alerts** - Get notified of node failures, resource exhaustion
4. **Log aggregation** - Use Loki or external SaaS for centralized logs
5. **Regular audits** - Run kube-bench for security compliance

---

## Troubleshooting

### Cluster Creation Fails

**Symptom:** Cluster creation times out or fails

**Possible Causes:**
- VPC UUID is invalid or doesn't exist
- Kubernetes version is not supported by DigitalOcean
- Region is invalid
- Insufficient quota in your DigitalOcean account

**Solution:**
```bash
# Verify VPC exists
doctl vpcs list

# Check available Kubernetes versions
doctl kubernetes options versions

# Check account limits
doctl account get
```

### Cannot Access Cluster API

**Symptom:** `kubectl` commands timeout or fail with connection refused

**Possible Causes:**
- Control plane firewall is enabled but your IP is not whitelisted
- Kubeconfig is not set up correctly
- Cluster is still provisioning

**Solution:**
```bash
# Check if your IP is whitelisted (if firewall is enabled)
curl ifconfig.me  # Get your current IP

# Update firewall rules if needed
# Edit manifest to add your IP to controlPlaneFirewallAllowedIps

# Re-fetch kubeconfig
openmcf kubeconfig digitaloceankubernetescluster <cluster-name> > ~/.kube/config

# Verify cluster status
openmcf status digitaloceankubernetescluster <cluster-name>
```

### Autoscaling Not Working

**Symptom:** Node pool doesn't scale despite pod resource requests

**Possible Causes:**
- Autoscaling not enabled (`autoScale: false`)
- Min/max nodes not configured
- Pod resource requests not set

**Solution:**
```bash
# Verify autoscaling is enabled in manifest
# Set autoScale: true, minNodes, and maxNodes

# Ensure pods have resource requests
kubectl get pods -o json | jq '.items[] | .spec.containers[] | .resources'

# Check cluster autoscaler logs
kubectl logs -n kube-system -l app=cluster-autoscaler
```

---

## Reference Links

- [DigitalOcean Kubernetes Documentation](https://docs.digitalocean.com/products/kubernetes/)
- [DOKS Pricing](https://www.digitalocean.com/pricing/kubernetes)
- [Kubernetes Version Support Policy](https://docs.digitalocean.com/products/kubernetes/details/versions/)
- [OpenMCF Documentation](https://docs.openmcf.org/)

---

**Note:** Replace placeholder values like `your-vpc-id-here` with your actual DigitalOcean VPC UUIDs before deploying. You can find VPC UUIDs in the DigitalOcean control panel or using `doctl vpcs list`.


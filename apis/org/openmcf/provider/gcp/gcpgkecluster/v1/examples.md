# GCP GKE Cluster Examples

This document provides real-world configuration examples for different GKE cluster deployment scenarios.

---

## Example 1: Minimal Production Cluster

A production-ready regional GKE cluster with all security features enabled.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeCluster
metadata:
  name: prod-cluster
  org: acme-corp
  env:
    id: production
spec:
  # GCP project (foreign key reference)
  projectId:
    valueRef:
      kind: GcpProject
      name: acme-prod
      path: status.outputs.project_id
  
  # Regional cluster for high availability (control plane across 3+ zones)
  location: us-central1
  
  # VPC networking configuration (foreign key references)
  subnetworkSelfLink:
    valueRef:
      kind: GcpSubnetwork
      name: prod-gke-subnet
      path: status.outputs.self_link
  
  # Pod IP secondary range
  clusterSecondaryRangeName:
    valueRef:
      kind: GcpSubnetwork
      name: prod-gke-subnet
      path: status.outputs.pods_secondary_range_name
  
  # Service IP secondary range
  servicesSecondaryRangeName:
    valueRef:
      kind: GcpSubnetwork
      name: prod-gke-subnet
      path: status.outputs.services_secondary_range_name
  
  # Control plane private endpoint CIDR (must be /28)
  masterIpv4CidrBlock: 172.16.0.0/28
  
  # Private nodes (no public IPs)
  enablePublicNodes: false
  
  # Regular release channel (recommended for production)
  releaseChannel: REGULAR
  
  # Enable network policies (Calico for microsegmentation)
  disableNetworkPolicy: false
  
  # Enable Workload Identity (IAM for pods)
  disableWorkloadIdentity: false
  
  # Cloud NAT for private node egress
  routerNatName:
    valueRef:
      kind: GcpRouterNat
      name: prod-nat
      path: metadata.name
```

**What this creates:**
- Regional cluster in `us-central1` (HA control plane)
- Private nodes with Cloud NAT for internet egress
- VPC-native networking (IP aliasing for pods and services)
- Workload Identity enabled for secure pod-level IAM
- Network policies enabled for microsegmentation
- Auto-upgrades on REGULAR channel

**Use when:**
- Production workloads requiring high availability
- Security is a priority (private nodes, Workload Identity)
- You want auto-upgrades with a balanced release schedule

---

## Example 2: Development Cluster (Cost-Optimized)

A zonal cluster for development with relaxed security settings and bleeding-edge Kubernetes versions.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeCluster
metadata:
  name: dev-cluster
  org: acme-corp
  env:
    id: development
spec:
  # GCP project (direct value for dev)
  projectId:
    value: acme-dev-12345
  
  # Zonal cluster (cheaper, single control plane replica)
  location: us-central1-a
  
  # VPC networking configuration (direct values)
  subnetworkSelfLink:
    value: https://www.googleapis.com/compute/v1/projects/acme-dev-12345/regions/us-central1/subnetworks/dev-gke-subnet
  
  clusterSecondaryRangeName:
    value: dev-pods
  
  servicesSecondaryRangeName:
    value: dev-services
  
  # Control plane CIDR (different /28 range from prod)
  masterIpv4CidrBlock: 172.16.0.16/28
  
  # Private nodes
  enablePublicNodes: false
  
  # Rapid release channel (bleeding-edge for dev)
  releaseChannel: RAPID
  
  # Disable network policies (simplify dev debugging)
  disableNetworkPolicy: true
  
  # Disable Workload Identity (not needed for dev)
  disableWorkloadIdentity: true
  
  # Cloud NAT reference
  routerNatName:
    value: dev-nat
```

**What this creates:**
- Zonal cluster in `us-central1-a` (lower cost, no HA)
- Private nodes with Cloud NAT
- No network policies (pods can reach all pods)
- No Workload Identity (simplified IAM)
- Rapid release channel (new Kubernetes versions quickly)

**Use when:**
- Development and testing environments
- Cost is a concern (zonal cluster, single control plane)
- You want the latest Kubernetes features
- Security requirements are relaxed

---

## Example 3: High-Security Production Cluster (Stable Channel)

A risk-averse production cluster with the most conservative upgrade strategy.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeCluster
metadata:
  name: secure-prod
  org: financial-corp
  env:
    id: production
spec:
  projectId:
    valueRef:
      kind: GcpProject
      name: financial-prod
      path: status.outputs.project_id
  
  # Regional cluster for HA
  location: us-east1
  
  subnetworkSelfLink:
    valueRef:
      kind: GcpSubnetwork
      name: secure-subnet
      path: status.outputs.self_link
  
  clusterSecondaryRangeName:
    valueRef:
      kind: GcpSubnetwork
      name: secure-subnet
      path: status.outputs.pods_secondary_range_name
  
  servicesSecondaryRangeName:
    valueRef:
      kind: GcpSubnetwork
      name: secure-subnet
      path: status.outputs.services_secondary_range_name
  
  masterIpv4CidrBlock: 172.16.0.32/28
  
  # Private nodes (maximum security)
  enablePublicNodes: false
  
  # Stable release channel (conservative upgrades)
  releaseChannel: STABLE
  
  # Enable network policies (required for compliance)
  disableNetworkPolicy: false
  
  # Enable Workload Identity (required for compliance)
  disableWorkloadIdentity: false
  
  routerNatName:
    valueRef:
      kind: GcpRouterNat
      name: secure-nat
      path: metadata.name
```

**What this creates:**
- Regional cluster with HA control plane
- Private nodes with all security features enabled
- Network policies for microsegmentation
- Workload Identity for pod-level IAM
- Stable release channel (most conservative, tested Kubernetes versions)

**Use when:**
- Highly regulated industries (finance, healthcare)
- Risk-averse production environments
- Compliance requirements mandate network policies and Workload Identity
- You prefer battle-tested Kubernetes versions

---

## Example 4: Multi-Cluster Setup (Same Region)

Multiple clusters in the same region for different teams or environments, sharing VPC infrastructure.

### Cluster 1: Frontend

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeCluster
metadata:
  name: frontend-cluster
  org: acme-corp
  env:
    id: production
spec:
  projectId:
    valueRef:
      kind: GcpProject
      name: acme-prod
      path: status.outputs.project_id
  location: us-central1
  subnetworkSelfLink:
    valueRef:
      kind: GcpSubnetwork
      name: frontend-subnet
      path: status.outputs.self_link
  clusterSecondaryRangeName:
    valueRef:
      kind: GcpSubnetwork
      name: frontend-subnet
      path: status.outputs.pods_secondary_range_name
  servicesSecondaryRangeName:
    valueRef:
      kind: GcpSubnetwork
      name: frontend-subnet
      path: status.outputs.services_secondary_range_name
  masterIpv4CidrBlock: 172.16.0.0/28  # First /28 range
  enablePublicNodes: false
  releaseChannel: REGULAR
  disableNetworkPolicy: false
  disableWorkloadIdentity: false
  routerNatName:
    valueRef:
      kind: GcpRouterNat
      name: frontend-nat
      path: metadata.name
```

### Cluster 2: Backend

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeCluster
metadata:
  name: backend-cluster
  org: acme-corp
  env:
    id: production
spec:
  projectId:
    valueRef:
      kind: GcpProject
      name: acme-prod
      path: status.outputs.project_id
  location: us-central1
  subnetworkSelfLink:
    valueRef:
      kind: GcpSubnetwork
      name: backend-subnet
      path: status.outputs.self_link
  clusterSecondaryRangeName:
    valueRef:
      kind: GcpSubnetwork
      name: backend-subnet
      path: status.outputs.pods_secondary_range_name
  servicesSecondaryRangeName:
    valueRef:
      kind: GcpSubnetwork
      name: backend-subnet
      path: status.outputs.services_secondary_range_name
  masterIpv4CidrBlock: 172.16.0.16/28  # Second /28 range (no overlap)
  enablePublicNodes: false
  releaseChannel: REGULAR
  disableNetworkPolicy: false
  disableWorkloadIdentity: false
  routerNatName:
    valueRef:
      kind: GcpRouterNat
      name: backend-nat
      path: metadata.name
```

**Why multiple clusters?**
- **Team isolation**: Frontend and backend teams have independent clusters
- **Blast radius containment**: Issues in one cluster don't affect the other
- **Independent scaling**: Frontend and backend can scale node pools independently
- **Security boundaries**: Network policies within clusters, VPC firewall rules between clusters

**Key considerations:**
- Use non-overlapping master CIDR blocks (`172.16.0.0/28`, `172.16.0.16/28`, etc.)
- Use separate subnets for each cluster (or shared subnet with careful IP planning)
- Each cluster needs its own Cloud NAT (or shared NAT across clusters)

---

## Example 5: Shared VPC (Enterprise Setup)

A GKE cluster in a service project consuming a subnet from a host project (Shared VPC).

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeCluster
metadata:
  name: service-cluster
  org: enterprise-corp
  env:
    id: production
spec:
  # Service project where cluster is created
  projectId:
    value: service-project-12345
  
  location: us-central1
  
  # Subnet from host project (Shared VPC)
  subnetworkSelfLink:
    value: https://www.googleapis.com/compute/v1/projects/host-project-99999/regions/us-central1/subnetworks/shared-gke-subnet
  
  # Secondary ranges from host project subnet
  clusterSecondaryRangeName:
    value: shared-pods
  
  servicesSecondaryRangeName:
    value: shared-services
  
  masterIpv4CidrBlock: 172.16.0.48/28
  enablePublicNodes: false
  releaseChannel: REGULAR
  disableNetworkPolicy: false
  disableWorkloadIdentity: false
  
  # Cloud NAT in host project
  routerNatName:
    value: shared-nat
```

**What this enables:**
- **Centralized network management**: Platform team manages VPC in host project
- **Service project isolation**: Application teams deploy clusters in service projects
- **Shared infrastructure**: Subnets, firewalls, NAT, VPN managed centrally
- **Cost efficiency**: Shared network resources across multiple service projects

**Prerequisites:**
- Shared VPC enabled between host project and service project
- Service project's GKE service account has Shared VPC Network User role in host project
- Subnet exists in host project with secondary ranges for pods and services

---

## Complete Manifest with Dependencies

Here's a complete manifest showing a GKE cluster with all prerequisite resources:

```yaml
---
apiVersion: gcp.openmcf.org/v1
kind: GcpProject
metadata:
  name: my-project
spec:
  projectName: my-gcp-project
  billingAccountId: "012345-6789AB-CDEF01"
  organizationId: "123456789012"

---
apiVersion: gcp.openmcf.org/v1
kind: GcpVpc
metadata:
  name: my-vpc
spec:
  projectId:
    valueRef:
      kind: GcpProject
      name: my-project
      path: status.outputs.project_id
  autoCreateSubnetworks: false  # Custom mode

---
apiVersion: gcp.openmcf.org/v1
kind: GcpSubnetwork
metadata:
  name: my-subnet
spec:
  projectId:
    valueRef:
      kind: GcpProject
      name: my-project
      path: status.outputs.project_id
  vpcSelfLink:
    valueRef:
      kind: GcpVpc
      name: my-vpc
      path: status.outputs.self_link
  region: us-central1
  ipCidrRange: 10.0.0.0/24  # Primary range for nodes
  secondaryIpRanges:
    - range_name: pods
      ipCidrRange: 10.1.0.0/16  # Pod IPs
    - range_name: services
      ipCidrRange: 10.2.0.0/20  # Service IPs

---
apiVersion: gcp.openmcf.org/v1
kind: GcpRouter
metadata:
  name: my-router
spec:
  projectId:
    valueRef:
      kind: GcpProject
      name: my-project
      path: status.outputs.project_id
  vpcSelfLink:
    valueRef:
      kind: GcpVpc
      name: my-vpc
      path: status.outputs.self_link
  region: us-central1

---
apiVersion: gcp.openmcf.org/v1
kind: GcpRouterNat
metadata:
  name: my-nat
spec:
  projectId:
    valueRef:
      kind: GcpProject
      name: my-project
      path: status.outputs.project_id
  routerName:
    valueRef:
      kind: GcpRouter
      name: my-router
      path: metadata.name
  region: us-central1

---
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeCluster
metadata:
  name: my-cluster
spec:
  projectId:
    valueRef:
      kind: GcpProject
      name: my-project
      path: status.outputs.project_id
  location: us-central1
  subnetworkSelfLink:
    valueRef:
      kind: GcpSubnetwork
      name: my-subnet
      path: status.outputs.self_link
  clusterSecondaryRangeName:
    value: pods
  servicesSecondaryRangeName:
    value: services
  masterIpv4CidrBlock: 172.16.0.0/28
  enablePublicNodes: false
  releaseChannel: REGULAR
  disableNetworkPolicy: false
  disableWorkloadIdentity: false
  routerNatName:
    valueRef:
      kind: GcpRouterNat
      name: my-nat
      path: metadata.name
```

**Deploy order:**
1. GcpProject
2. GcpVpc
3. GcpSubnetwork
4. GcpRouter
5. GcpRouterNat
6. GcpGkeCluster

OpenMCF automatically handles dependency ordering based on foreign key references.

---

## Next Steps

After creating a GKE cluster, you'll typically want to:

1. **Provision Node Pools**: Use `GcpGkeNodePool` resources to add compute nodes
2. **Configure kubectl**: Access the cluster with `gcloud container clusters get-credentials`
3. **Set Up Workload Identity**: Bind Kubernetes Service Accounts to Google Service Accounts
4. **Deploy Applications**: Apply Kubernetes manifests for your workloads
5. **Implement Network Policies**: Define microsegmentation rules

For more examples and patterns, see the [research documentation](docs/README.md).


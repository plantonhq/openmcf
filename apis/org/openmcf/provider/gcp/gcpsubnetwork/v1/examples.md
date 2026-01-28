# GCP Subnetwork Examples

This document provides practical examples for common GCP subnetwork deployment scenarios.

## Table of Contents

1. [Basic Subnet](#example-1-basic-subnet)
2. [Subnet with Secondary Ranges for GKE](#example-2-subnet-with-secondary-ranges-for-gke)
3. [Multi-Region Subnets](#example-3-multi-region-subnets)
4. [Private Google Access Enabled](#example-4-private-google-access-enabled)
5. [VPC Connector Subnet (Cloud Run)](#example-5-vpc-connector-subnet-cloud-run)
6. [Large GKE Cluster Subnet](#example-6-large-gke-cluster-subnet)

---

## Example 1: Basic Subnet

A simple subnet for general compute workloads.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpSubnetwork
metadata:
  name: basic-subnet
spec:
  projectId: my-project-123
  vpcSelfLink: projects/my-project-123/global/networks/my-vpc
  subnetworkName: basic-subnet
  region: us-central1
  ipCidrRange: 10.0.0.0/24
```

**Use Case**: Small development environment or dedicated service subnet.

**IP Capacity**: 256 IPs (252 usable after GCP reserves 4)

---

## Example 2: Subnet with Secondary Ranges for GKE

Standard production subnet for a GKE cluster with pod and service ranges.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpSubnetwork
metadata:
  name: gke-subnet-us-central1
spec:
  projectId: prod-project
  vpcSelfLink: projects/prod-project/global/networks/prod-vpc
  subnetworkName: gke-subnet-us-central1
  region: us-central1
  ipCidrRange: 10.100.0.0/20          # 4,096 IPs for nodes
  secondaryIpRanges:
    - rangeName: pods
      ipCidrRange: 10.100.16.0/18     # 16,384 IPs for pods
    - rangeName: services
      ipCidrRange: 10.100.80.0/24     # 256 IPs for services
  privateIpGoogleAccess: true
```

**Use Case**: Production GKE cluster with ~100 nodes, 11,000 pods, 200 services.

**Calculation**:
- **Nodes**: 100 nodes × 1 IP = 100 IPs → /20 (4,096 IPs) provides ample room
- **Pods**: 100 nodes × 110 pods/node = 11,000 IPs → /18 (16,384 IPs)
- **Services**: 200 services → /24 (256 IPs)

---

## Example 3: Multi-Region Subnets

Deploy subnets across multiple regions for high availability.

```yaml
# US Central
apiVersion: gcp.openmcf.org/v1
kind: GcpSubnetwork
metadata:
  name: app-subnet-us-central1
spec:
  projectId: multi-region-project
  vpcSelfLink: projects/multi-region-project/global/networks/global-vpc
  subnetworkName: app-subnet-us-central1
  region: us-central1
  ipCidrRange: 10.0.0.0/20
  privateIpGoogleAccess: true
---
# US East
apiVersion: gcp.openmcf.org/v1
kind: GcpSubnetwork
metadata:
  name: app-subnet-us-east1
spec:
  projectId: multi-region-project
  vpcSelfLink: projects/multi-region-project/global/networks/global-vpc
  subnetworkName: app-subnet-us-east1
  region: us-east1
  ipCidrRange: 10.0.16.0/20
  privateIpGoogleAccess: true
---
# Europe West
apiVersion: gcp.openmcf.org/v1
kind: GcpSubnetwork
metadata:
  name: app-subnet-europe-west1
spec:
  projectId: multi-region-project
  vpcSelfLink: projects/multi-region-project/global/networks/global-vpc
  subnetworkName: app-subnet-europe-west1
  region: europe-west1
  ipCidrRange: 10.0.32.0/20
  privateIpGoogleAccess: true
---
# Asia Southeast
apiVersion: gcp.openmcf.org/v1
kind: GcpSubnetwork
metadata:
  name: app-subnet-asia-southeast1
spec:
  projectId: multi-region-project
  vpcSelfLink: projects/multi-region-project/global/networks/global-vpc
  subnetworkName: app-subnet-asia-southeast1
  region: asia-southeast1
  ipCidrRange: 10.0.48.0/20
  privateIpGoogleAccess: true
```

**Use Case**: Global application with regional GKE clusters or compute instances.

**IP Planning**: Each region gets a /20 (4,096 IPs) carved from 10.0.0.0/18 (total 16,384 IPs).

---

## Example 4: Private Google Access Enabled

Subnet optimized for internal-only workloads accessing Google APIs.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpSubnetwork
metadata:
  name: private-subnet
spec:
  projectId: secure-project
  vpcSelfLink: projects/secure-project/global/networks/secure-vpc
  subnetworkName: private-subnet
  region: us-west1
  ipCidrRange: 10.200.0.0/22
  privateIpGoogleAccess: true
```

**Use Case**: VMs without external IPs that need to:
- Pull images from Google Container Registry (GCR)
- Access Cloud Storage (GCS) buckets
- Query BigQuery datasets
- Call other Google APIs

**Security Benefit**: No public internet exposure; all Google API traffic stays on Google's internal network.

---

## Example 5: VPC Connector Subnet (Cloud Run)

Dedicated /28 subnet for Serverless VPC Access connector.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpSubnetwork
metadata:
  name: vpc-connector-subnet
spec:
  projectId: serverless-project
  vpcSelfLink: projects/serverless-project/global/networks/app-vpc
  subnetworkName: vpc-connector-subnet
  region: us-central1
  ipCidrRange: 10.8.0.0/28           # Exactly 16 IPs (VPC connector requirement)
```

**Use Case**: Serverless VPC Access connector for Cloud Run or Cloud Functions.

**Requirements**:
- **Must be /28** (16 IPs) for VPC connectors
- **Cannot share** with other resources
- **Separate subnet per connector** (one per region)

**Next Step**: Create a `VpcAccessConnector` referencing this subnet.

---

## Example 6: Large GKE Cluster Subnet

High-capacity subnet for very large GKE clusters (500+ nodes).

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpSubnetwork
metadata:
  name: mega-gke-subnet
spec:
  projectId: enterprise-project
  vpcSelfLink: projects/enterprise-project/global/networks/enterprise-vpc
  subnetworkName: mega-gke-subnet
  region: us-central1
  ipCidrRange: 10.50.0.0/16          # 65,536 IPs for nodes
  secondaryIpRanges:
    - rangeName: pods
      ipCidrRange: 10.51.0.0/16      # 65,536 IPs for pods
    - rangeName: services
      ipCidrRange: 10.52.0.0/24      # 256 IPs for services
  privateIpGoogleAccess: true
```

**Use Case**: Enterprise-scale GKE cluster with 500+ nodes.

**Calculation**:
- **Nodes**: 500 nodes → /16 (65,536 IPs) provides massive headroom
- **Pods**: 500 nodes × 110 pods/node = 55,000 IPs → /16 (65,536 IPs)
- **Services**: 200-500 services → /24 (256 IPs)

**Note**: This is an extremely large allocation. Most production clusters use /20 primary and /18 pods.

---

## Example 7: Environment-Specific Subnets

Separate IP ranges per environment (dev, staging, prod).

```yaml
# Development
apiVersion: gcp.openmcf.org/v1
kind: GcpSubnetwork
metadata:
  name: dev-subnet
spec:
  projectId: dev-project
  vpcSelfLink: projects/dev-project/global/networks/dev-vpc
  subnetworkName: dev-subnet
  region: us-central1
  ipCidrRange: 10.2.0.0/20
  privateIpGoogleAccess: true
---
# Staging
apiVersion: gcp.openmcf.org/v1
kind: GcpSubnetwork
metadata:
  name: staging-subnet
spec:
  projectId: staging-project
  vpcSelfLink: projects/staging-project/global/networks/staging-vpc
  subnetworkName: staging-subnet
  region: us-central1
  ipCidrRange: 10.1.0.0/20
  privateIpGoogleAccess: true
---
# Production
apiVersion: gcp.openmcf.org/v1
kind: GcpSubnetwork
metadata:
  name: prod-subnet
spec:
  projectId: prod-project
  vpcSelfLink: projects/prod-project/global/networks/prod-vpc
  subnetworkName: prod-subnet
  region: us-central1
  ipCidrRange: 10.0.0.0/20
  privateIpGoogleAccess: true
```

**Use Case**: Isolate environments with separate VPCs and non-overlapping IP ranges.

**IP Planning**:
- **Production**: 10.0.0.0/16 (highest priority)
- **Staging**: 10.1.0.0/16
- **Development**: 10.2.0.0/16

---

## Example 8: Referencing Outputs in Other Resources

Use subnet outputs in downstream resources like GKE clusters.

```yaml
# Step 1: Create the subnet
apiVersion: gcp.openmcf.org/v1
kind: GcpSubnetwork
metadata:
  name: gke-prod-subnet
spec:
  projectId: prod-gcp-project
  vpcSelfLink: projects/prod-gcp-project/global/networks/prod-vpc
  subnetworkName: gke-prod-subnet
  region: us-west1
  ipCidrRange: 10.50.0.0/20
  secondaryIpRanges:
    - rangeName: pods
      ipCidrRange: 10.50.16.0/18
    - rangeName: services
      ipCidrRange: 10.50.80.0/24
  privateIpGoogleAccess: true
---
# Step 2: Reference in GKE cluster (pseudo-example)
apiVersion: gcp.openmcf.org/v1
kind: GkeCluster
metadata:
  name: prod-cluster
spec:
  projectId: prod-gcp-project
  region: us-west1
  networkConfig:
    # Reference the subnet by name or self-link output
    subnetworkRef:
      kind: GcpSubnetwork
      name: gke-prod-subnet
      fieldPath: status.outputs.subnetwork_self_link
    clusterSecondaryRangeName: pods
    servicesSecondaryRangeName: services
```

**Use Case**: Modular infrastructure where subnets are created separately from clusters.

**Benefit**: Reuse subnets across multiple resources, simplify dependency management.

---

## Example 9: CIDR Planning Template

A template for planning IP allocations across environments and regions.

```yaml
# IP Allocation Plan: 10.0.0.0/8

# Production Environment: 10.0.0.0/16
# ├─ us-central1:   10.0.0.0/20   (app subnet)
# ├─ us-east1:      10.0.16.0/20  (app subnet)
# ├─ us-west1:      10.0.32.0/18  (GKE: primary + secondaries)
# └─ Reserved:      10.0.96.0/19  (future expansion)

# Staging Environment: 10.1.0.0/16
# ├─ us-central1:   10.1.0.0/20
# └─ Reserved:      10.1.16.0/20

# Development Environment: 10.2.0.0/16
# └─ us-central1:   10.2.0.0/20

# Reserved for Future: 10.3.0.0/16 - 10.255.0.0/16
```

**Production us-central1 App Subnet**:

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpSubnetwork
metadata:
  name: prod-app-us-central1
spec:
  projectId: prod-project
  vpcSelfLink: projects/prod-project/global/networks/prod-vpc
  subnetworkName: prod-app-us-central1
  region: us-central1
  ipCidrRange: 10.0.0.0/20
  privateIpGoogleAccess: true
```

**Use Case**: Systematic IP planning for large organizations.

**Benefit**: Avoid CIDR conflicts, enable VPC peering, simplify hybrid connectivity.

---

## CIDR Cheat Sheet

Quick reference for common subnet sizes:

| CIDR | Total IPs | Usable IPs | Use Case |
|------|-----------|------------|----------|
| /28  | 16        | 12         | VPC Connector (required for Cloud Run) |
| /26  | 64        | 60         | Small dev environment |
| /24  | 256       | 252        | Standard app subnet (100-200 VMs) |
| /22  | 1,024     | 1,020      | Medium GKE cluster (50-100 nodes) |
| /20  | 4,096     | 4,092      | Production GKE cluster (100-300 nodes) |
| /18  | 16,384    | 16,380     | GKE pod range (100 nodes × 110 pods) |
| /17  | 32,768    | 32,764     | Large GKE pod range (200+ nodes) |
| /16  | 65,536    | 65,532     | Enterprise environment or mega-cluster |

**Remember**: GCP reserves 4 IPs per subnet (network, broadcast, gateway, DNS).

---

## Common Pitfalls

### ❌ Under-Sized Pod Range

```yaml
# BAD: Too small for 100-node cluster
secondaryIpRanges:
  - rangeName: pods
    ipCidrRange: 10.0.16.0/20  # Only 4,096 IPs
```

**Problem**: 100 nodes × 110 pods = 11,000 IPs needed, but /20 only provides 4,096.

**Fix**: Use /18 (16,384 IPs) or /17 (32,768 IPs).

### ❌ Overlapping CIDR Ranges

```yaml
# BAD: Overlapping with existing subnet
# Existing subnet: 10.0.0.0/20
# New subnet:      10.0.8.0/22  # Overlaps!
```

**Problem**: CIDR conflict prevents subnet creation.

**Fix**: Use non-overlapping range like 10.0.16.0/22.

### ❌ Missing Secondary Ranges for GKE

```yaml
# BAD: No secondary ranges
spec:
  ipCidrRange: 10.0.0.0/20
  # Missing secondary_ip_ranges!
```

**Problem**: GKE cluster creation will fail or force routes-based mode (deprecated).

**Fix**: Always define secondary ranges for GKE subnets.

---

## Next Steps

- **Deep Dive**: Read [docs/README.md](docs/README.md) for comprehensive CIDR planning guidance
- **API Reference**: See [README.md](README.md) for full field descriptions
- **GKE Integration**: Check `GkeCluster` examples for cluster configuration
- **VPC Setup**: See `GcpVpc` for custom-mode VPC creation

---

**Questions or Issues?** Consult the [README.md](README.md) or [research documentation](docs/README.md).


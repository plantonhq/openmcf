# GCP Cloud Router NAT Examples

This document provides comprehensive examples of `GcpRouterNat` resource configurations, ranging from minimal setups to production-grade configurations with manual IP allocation and specific subnet coverage.

---

## Example 1: Minimal Configuration (Auto-Allocation, All Subnets)

The simplest configuration for development or staging environments. Uses auto-allocated NAT IPs and covers all subnets in the region.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpRouterNat
metadata:
  name: dev-nat-uscentral1
spec:
  projectId: "my-project"
  vpcSelfLink: "https://www.googleapis.com/compute/v1/projects/my-project/global/networks/my-vpc"
  routerName: dev-router-uscentral1
  natName: dev-nat-uscentral1
  region: us-central1
```

**Features:**
- Auto-allocated NAT IPs (Google manages IP allocation and scaling)
- All subnets in `us-central1` automatically covered
- Default logging: `ERRORS_ONLY` (recommended for production)
- Uses VPC network by direct self-link

**Use Case:** Quick setup for development environments or private GKE clusters needing internet access.

---

## Example 2: Reference Existing GcpVpc Resource

Reference an existing `GcpVpc` resource instead of providing a direct VPC self-link. This is the recommended approach when using OpenMCF resources.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpRouterNat
metadata:
  name: staging-nat-uswest1
  org: acme-corp
  env: staging
spec:
  projectId:
    ref:
      kind: GcpProject
      name: my-gcp-project
      fieldPath: status.outputs.project_id
  vpcSelfLink:
    ref:
      kind: GcpVpc
      name: staging-vpc
      fieldPath: status.outputs.network_self_link
  routerName: staging-router-uswest1
  natName: staging-nat-uswest1
  region: us-west1
```

**Features:**
- References `GcpVpc` resource named `staging-vpc`
- Automatically uses the VPC's network self-link from its outputs
- Inherits organization and environment labels from metadata
- Auto-allocated IPs and all-subnets coverage (defaults)

**Use Case:** Standard staging environment with existing VPC resource.

---

## Example 3: Manual IP Allocation for Partner Allowlisting

Production environment where external partners require allowlisting specific egress IPs.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpRouterNat
metadata:
  name: prod-nat-uscentral1
  org: acme-corp
  env: prod
  labels:
    cost-center: cc-4510
    compliance: pci-dss
spec:
  vpcSelfLink:
    ref:
      kind: GcpVpc
      name: prod-vpc
      fieldPath: status.outputs.network_self_link
  routerName: prod-router-uscentral1
  natName: prod-nat-uscentral1
  region: us-central1
  natIpNames:
    - "prod-nat-ip-1-uscentral1"
    - "prod-nat-ip-2-uscentral1"
  logFilter: ERRORS_ONLY
```

**Features:**
- Manual allocation mode: uses specified static external IPs
- Two static IPs for redundancy and load distribution
- Production logging enabled (`ERRORS_ONLY`)
- Governance labels for cost tracking and compliance

**Prerequisites:**
Create static IPs before deploying NAT:
```bash
gcloud compute addresses create prod-nat-ip-1-uscentral1 --region=us-central1
gcloud compute addresses create prod-nat-ip-2-uscentral1 --region=us-central1
```

**Use Case:** Production environment where external partners (payment processors, APIs) require allowlisting your egress IPs in their firewalls.

---

## Example 4: Specific Subnet Coverage

Restrict NAT to specific subnets instead of all subnets in the region. Useful for fine-grained egress control.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpRouterNat
metadata:
  name: staging-nat-specific-subnets
  org: acme-corp
  env: staging
spec:
  vpcSelfLink:
    ref:
      kind: GcpVpc
      name: staging-vpc
      fieldPath: status.outputs.network_self_link
  routerName: staging-router-useast1
  natName: staging-nat-specific-subnets
  region: us-east1
  subnetworkSelfLinks:
    - "https://www.googleapis.com/compute/v1/projects/my-project/regions/us-east1/subnetworks/private-subnet-1"
    - "https://www.googleapis.com/compute/v1/projects/my-project/regions/us-east1/subnetworks/private-subnet-2"
  logFilter: ERRORS_ONLY
```

**Features:**
- NAT only applies to specified subnets (`private-subnet-1` and `private-subnet-2`)
- Other subnets in `us-east1` are excluded from NAT
- Auto-allocated IPs (default)
- Production logging enabled

**Use Case:** Security requirement to restrict internet access to specific subnets (e.g., GKE node subnets) while excluding others (e.g., database subnets).

---

## Example 5: Disable Logging for Cost Optimization (Non-Production)

Development environment where logging is disabled to reduce Cloud Logging costs.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpRouterNat
metadata:
  name: dev-nat-uswest2
  org: acme-corp
  env: dev
spec:
  vpcSelfLink:
    ref:
      kind: GcpVpc
      name: dev-vpc
      fieldPath: status.outputs.network_self_link
  routerName: dev-router-uswest2
  natName: dev-nat-uswest2
  region: us-west2
  logFilter: DISABLED
```

**Features:**
- Logging disabled (`DISABLED`) to reduce costs
- Auto-allocated IPs and all-subnets coverage
- Minimal configuration for development

**Use Case:** Development or sandbox environments where NAT logging is not required and cost optimization is prioritized.

---

## Example 6: Full Logging for Security Auditing

Production environment with comprehensive logging for security auditing or troubleshooting.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpRouterNat
metadata:
  name: prod-nat-audit-uscentral1
  org: acme-corp
  env: prod
  labels:
    compliance: sox
    audit: true
spec:
  vpcSelfLink:
    ref:
      kind: GcpVpc
      name: prod-vpc
      fieldPath: status.outputs.network_self_link
  routerName: prod-router-audit-uscentral1
  natName: prod-nat-audit-uscentral1
  region: us-central1
  logFilter: ALL
```

**Features:**
- Full logging (`ALL`): logs every NAT translation
- Useful for security auditing and detailed troubleshooting
- **Warning:** Generates significant log volume and costs

**Use Case:** Production environment with strict compliance requirements (SOX, PCI-DSS) requiring full egress traffic auditing.

**Cost Consideration:** `ALL` logging can generate substantial log volume. Use sparingly and ensure log retention policies are configured to manage costs.

---

## Example 7: Multi-Region Deployment (Active-Active)

Deploy NAT gateways in multiple regions for active-active architecture.

**us-central1:**
```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpRouterNat
metadata:
  name: prod-nat-uscentral1
  org: acme-corp
  env: prod
spec:
  vpcSelfLink:
    ref:
      kind: GcpVpc
      name: prod-vpc
      fieldPath: status.outputs.network_self_link
  routerName: prod-router-uscentral1
  natName: prod-nat-uscentral1
  region: us-central1
  natIpNames:
    - "prod-nat-ip-uscentral1"
```

**europe-west1:**
```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpRouterNat
metadata:
  name: prod-nat-europewest1
  org: acme-corp
  env: prod
spec:
  vpcSelfLink:
    ref:
      kind: GcpVpc
      name: prod-vpc
      fieldPath: status.outputs.network_self_link
  routerName: prod-router-europewest1
  natName: prod-nat-europewest1
  region: europe-west1
  natIpNames:
    - "prod-nat-ip-europewest1"
```

**Features:**
- Separate NAT gateway in each region
- Each region uses its own static IP for predictable egress
- Supports active-active or active-passive architectures

**Use Case:** Multi-region production deployment where each region needs independent NAT egress.

**Partner Allowlisting:** Provide NAT IPs from all active regions to partners for firewall rules.

---

## Example 8: Private GKE Cluster with NAT

Deploy NAT for a private GKE cluster to enable container image pulls and external API access.

**Step 1: Create VPC and Subnets**
```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpVpc
metadata:
  name: gke-vpc
spec:
  projectId: my-gke-project
  autoCreateSubnetworks: false
```

**Step 2: Deploy NAT Gateway**
```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpRouterNat
metadata:
  name: gke-nat-uscentral1
spec:
  vpcSelfLink:
    ref:
      kind: GcpVpc
      name: gke-vpc
      fieldPath: status.outputs.network_self_link
  routerName: gke-router-uscentral1
  natName: gke-nat-uscentral1
  region: us-central1
  logFilter: ERRORS_ONLY
```

**Step 3: Deploy Private GKE Cluster**
```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGkeCluster
metadata:
  name: private-gke-cluster
spec:
  vpcSelfLink:
    ref:
      kind: GcpVpc
      name: gke-vpc
      fieldPath: status.outputs.network_self_link
  region: us-central1
  privateCluster: true
  # NAT is automatically used by private nodes
```

**Features:**
- Private GKE cluster nodes have no public IPs
- NAT gateway provides egress for:
  - Container image pulls (Docker Hub, GCR, Artifact Registry)
  - Package manager (apt, yum)
  - External API calls from applications
- Enable Private Google Access on subnets for efficient Google API access

**Use Case:** Production GKE cluster with enhanced security (no public IPs on nodes) while maintaining internet egress capability.

---

## Example 9: Hybrid Cloud with VPN and NAT

NAT gateway for private resources that need both on-premises connectivity (via VPN) and internet access.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpRouterNat
metadata:
  name: hybrid-nat-uscentral1
  org: acme-corp
  env: prod
spec:
  vpcSelfLink:
    ref:
      kind: GcpVpc
      name: hybrid-vpc
      fieldPath: status.outputs.network_self_link
  routerName: hybrid-router-uscentral1
  natName: hybrid-nat-uscentral1
  region: us-central1
  subnetworkSelfLinks:
    - "https://www.googleapis.com/compute/v1/projects/my-project/regions/us-central1/subnetworks/app-subnet"
  natIpNames:
    - "hybrid-nat-ip-uscentral1"
  logFilter: ERRORS_ONLY
```

**Features:**
- NAT covers only application subnet (excludes VPN gateway subnet)
- Static IP for external partner integration
- Private resources can reach:
  - On-premises systems via VPN (uses VPN gateway, not NAT)
  - Internet services via NAT (uses NAT gateway)

**Use Case:** Hybrid cloud deployment during migration or for long-term on-premises integration.

---

## Example 10: Complete Production Configuration

A comprehensive production configuration with all features enabled.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpRouterNat
metadata:
  name: prod-payment-nat-uscentral1
  org: acme-corp
  env: prod
  labels:
    cost-center: cc-4510
    compliance: pci-dss
    component: payment-processing
    criticality: high
spec:
  vpcSelfLink:
    ref:
      kind: GcpVpc
      name: prod-payment-vpc
      fieldPath: status.outputs.network_self_link
  region: us-central1
  subnetworkSelfLinks:
    - ref:
        kind: GcpSubnet
        name: payment-app-subnet
        fieldPath: status.outputs.subnet_self_link
    - ref:
        kind: GcpSubnet
        name: payment-worker-subnet
        fieldPath: status.outputs.subnet_self_link
  natIpNames:
    - "prod-payment-nat-ip-1-uscentral1"
    - "prod-payment-nat-ip-2-uscentral1"
    - "prod-payment-nat-ip-3-uscentral1"
  logFilter: ERRORS_ONLY
```

**Features:**
- References VPC and subnet resources (fully declarative)
- Manual IP allocation with 3 static IPs (high-throughput workload)
- Specific subnet coverage (only payment-related subnets)
- Production logging (`ERRORS_ONLY`)
- Comprehensive governance labels

**Use Case:** Critical production payment processing workload with:
- External payment processor requiring allowlisted IPs
- High egress throughput (3 IPs for capacity)
- Security isolation (only specific subnets)
- Compliance tracking (PCI-DSS labels)

---

## Field Descriptions

### Required Fields

- **`metadata.name`**: Unique name for the Cloud Router and NAT gateway.
- **`spec.project_id`**: GCP project ID where the Cloud Router and NAT will be created (string value or reference to `GcpProject` resource).
- **`spec.vpc_self_link`**: VPC network reference (self-link or reference to `GcpVpc` resource).
- **`spec.region`**: GCP region for the Cloud Router and NAT (must match where your private resources are located).
- **`spec.router_name`**: Name of the Cloud Router to create (1-63 lowercase characters, must start with a letter).
- **`spec.nat_name`**: Name of the NAT configuration on the router (1-63 lowercase characters, must start with a letter).

### Optional Fields

- **`spec.subnetwork_self_links`**: List of specific subnets to cover with NAT.
  - **Default:** Empty list (all subnets in the region are covered).
  - **Use:** Security requirement to restrict NAT to specific subnets.
  
- **`spec.nat_ip_names`**: List of static external IP names for manual allocation.
  - **Default:** Empty list (auto-allocation mode).
  - **Use:** External partners require allowlisting specific egress IPs.
  - **Note:** Static IPs must be created in advance in the same region.
  
- **`spec.log_filter`**: NAT translation logging level.
  - **Default:** `ERRORS_ONLY` (recommended for production).
  - **Options:**
    - `DISABLED`: No logging (use for non-production to reduce costs).
    - `ERRORS_ONLY`: Log translation errors only (recommended).
    - `ALL`: Log all translations (security auditing, generates high log volume).

---

## Best Practices

### 1. Use Auto-Allocation Unless Allowlisting is Required

**Default (Recommended):**
```yaml
natIpNames: []  # or omit field entirely
```

**Manual Allocation (Only if needed):**
```yaml
natIpNames:
  - "nat-ip-1"
  - "nat-ip-2"
```

**Why:** Auto-allocation is simpler, scales automatically, and avoids paying for unused static IPs.

### 2. Cover All Subnets by Default

**Default (Recommended):**
```yaml
subnetworkSelfLinks: []  # or omit field entirely
```

**Specific Subnets (Only if needed):**
```yaml
subnetworkSelfLinks:
  - "https://www.googleapis.com/compute/v1/projects/my-project/regions/us-central1/subnetworks/subnet-1"
```

**Why:** All-subnets coverage is future-proof—new subnets automatically get NAT without config updates.

### 3. Enable ERRORS_ONLY Logging in Production

**Production (Recommended):**
```yaml
logFilter: ERRORS_ONLY  # or omit for default
```

**Development (Cost Optimization):**
```yaml
logFilter: DISABLED
```

**Security Auditing (High Cost):**
```yaml
logFilter: ALL
```

**Why:** `ERRORS_ONLY` detects NAT issues (port exhaustion, connection failures) without excessive log volume.

### 4. Deploy One NAT per Region

Cloud NAT is regional—deploy separate NAT gateways in each region:

- `us-central1` → `prod-nat-uscentral1`
- `europe-west1` → `prod-nat-europewest1`
- `asia-east1` → `prod-nat-asiaeast1`

### 5. Combine with Private Google Access

Enable Private Google Access on subnets to keep Google API traffic internal:

- Traffic to GCS, GCR, Cloud SQL Admin API bypasses NAT (no egress costs)
- NAT only handles true external traffic (third-party APIs, package repos)

### 6. Create Static IPs Before Deploying Manual Allocation

If using manual allocation, create static IPs first:

```bash
gcloud compute addresses create nat-ip-1 --region=us-central1
gcloud compute addresses create nat-ip-2 --region=us-central1
```

Then reference them:
```yaml
natIpNames:
  - "nat-ip-1"
  - "nat-ip-2"
```

### 7. Plan Maintenance Windows for Allocation Changes

Switching from auto to manual allocation (or vice versa) briefly disrupts existing connections:

- Plan during maintenance window
- Notify stakeholders
- Monitor connection recovery

---

## Next Steps

After deploying Cloud Router NAT:

1. **Verify Connectivity**: Test that private resources can reach the internet.
2. **Monitor Metrics**: Set up Cloud Monitoring dashboards for port allocation and dropped connections.
3. **Configure Alerts**: Alert on port exhaustion or connection failures.
4. **Review Logs**: Check Cloud Logging for NAT translation errors.
5. **Partner Integration**: Provide static NAT IPs to partners for allowlisting.

---

## Related Resources

- **GcpVpc**: Create custom VPC networks with subnets.
- **GcpSubnet**: Define subnets within VPC networks.
- **GcpStaticIp**: Create static external IPs for manual NAT allocation.
- **GcpGkeCluster**: Deploy private GKE clusters that use Cloud NAT.
- **GcpFirewall**: Define egress firewall rules.

---

## Common Patterns

### Pattern: Private GKE + NAT
```yaml
# 1. VPC
kind: GcpVpc
metadata:
  name: gke-vpc

# 2. NAT
kind: GcpRouterNat
metadata:
  name: gke-nat
spec:
  vpcSelfLink: { ref to gke-vpc }
  region: us-central1

# 3. Private GKE Cluster
kind: GcpGkeCluster
metadata:
  name: private-gke
spec:
  vpcSelfLink: { ref to gke-vpc }
  privateCluster: true
```

### Pattern: Multi-Region with Manual IPs
```yaml
# Region 1: us-central1
kind: GcpRouterNat
metadata:
  name: nat-uscentral1
spec:
  region: us-central1
  natIpNames: ["nat-ip-uscentral1"]

# Region 2: europe-west1
kind: GcpRouterNat
metadata:
  name: nat-europewest1
spec:
  region: europe-west1
  natIpNames: ["nat-ip-europewest1"]
```

### Pattern: Selective Subnet Coverage
```yaml
kind: GcpRouterNat
metadata:
  name: selective-nat
spec:
  region: us-central1
  subnetworkSelfLinks:
    - { ref to app-subnet }
    - { ref to worker-subnet }
    # database-subnet is excluded
```

---

For more details, see [README.md](README.md) and [docs/README.md](docs/README.md).


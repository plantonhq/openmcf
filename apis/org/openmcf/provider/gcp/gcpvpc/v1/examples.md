# GCP VPC Examples

This document provides comprehensive examples for deploying GCP VPCs using OpenMCF. Each example includes the manifest YAML and explains the use case and architecture considerations.

## Table of Contents

- [Example 1: Basic Development VPC](#example-1-basic-development-vpc)
- [Example 2: Production VPC with Global Routing](#example-2-production-vpc-with-global-routing)
- [Example 3: Multi-Environment Setup with Consistent Naming](#example-3-multi-environment-setup-with-consistent-naming)
- [Example 4: Shared VPC Host Network](#example-4-shared-vpc-host-network)
- [Example 5: VPC with Foreign Key Reference](#example-5-vpc-with-foreign-key-reference)
- [Example 6: VPC with Private Services Access](#example-6-vpc-with-private-services-access)

---

## Example 1: Basic Development VPC

### Use Case

A simple, custom-mode VPC for development and testing workloads. This is the most common starting point—a clean VPC with no automatic subnets, allowing you to explicitly create subnets in only the regions you need.

### Architecture

- **Custom mode**: No automatic subnet creation
- **Regional routing**: Default routing mode (simplest for single-region or independent multi-region workloads)
- **Single project**: Directly deployed to a specific GCP project

### Manifest

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpVpc
metadata:
  name: dev-network
  id: gcpvpc-dev-001
  org: engineering
  env: development
spec:
  projectId:
    value: my-dev-project-123
  networkName: dev-network
  autoCreateSubnetworks: false  # Custom mode (recommended)
```

### What Gets Created

- A VPC named `dev-network` in project `my-dev-project-123`
- No subnets (you create them separately as needed)
- Regional routing mode (Cloud Routers advertise routes only within their region)
- GCP labels for `resource`, `resource-id`, `resource-org`, and `env` for tracking

### Next Steps After Deployment

Create subnets in the regions you need:

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpSubnetwork
metadata:
  name: dev-uswest1-subnet
spec:
  vpcSelfLink:
    ref:
      kind: GcpVpc
      name: dev-network
  region: us-west1
  ipCidrRange: 10.10.0.0/16
```

---

## Example 2: Production VPC with Global Routing

### Use Case

A production VPC with global routing enabled, suitable for multi-region architectures with hybrid connectivity (Cloud VPN or Cloud Interconnect). This setup allows Cloud Routers in any region to advertise routes to all other regions, which is essential for on-premises connectivity.

### Architecture

- **Custom mode**: Explicit subnet creation
- **Global routing**: Cloud Routers advertise routes across all regions
- **Production-grade**: Suitable for Shared VPC and hybrid connectivity

### Manifest

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpVpc
metadata:
  name: prod-network
  id: gcpvpc-prod-001
  org: platform-team
  env: production
spec:
  projectId:
    value: prod-network-host-project
  networkName: prod-network
  autoCreateSubnetworks: false
  routingMode: GLOBAL  # Required for multi-region hybrid connectivity
```

### What Gets Created

- A VPC named `prod-network` in the host project
- Global routing mode enabled (routes advertised across all regions)
- Ready to be configured as a Shared VPC host

### When to Use Global Routing

**Use GLOBAL routing when:**
- You have Cloud VPN or Cloud Interconnect connecting to on-premises networks
- You need multi-region failover for hybrid connectivity
- You want Cloud Routers in any region to advertise routes to all regions

**Use REGIONAL routing (default) when:**
- Your workloads are entirely cloud-native
- Each region operates independently
- You don't need cross-region route advertisement

### Cost Consideration

Global routing does not incur additional charges, but it does change routing behavior. Ensure your firewall rules and route priorities are configured correctly to avoid unintended traffic paths.

---

## Example 3: Multi-Environment Setup with Consistent Naming

### Use Case

A consistent VPC setup across development, staging, and production environments, each in its own GCP project. This pattern ensures IP ranges don't overlap and environments are isolated while maintaining a consistent architecture.

### Architecture

- **Environment isolation**: Each environment in a separate project
- **Non-overlapping IP ranges**: Plan subnets to avoid conflicts (important for peering or future hybrid connectivity)
- **Consistent naming**: Same structure across environments for operational clarity

### Development VPC

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpVpc
metadata:
  name: app-network
  id: gcpvpc-dev-app-001
  org: product-team
  env: development
spec:
  projectId:
    value: myapp-dev-project
  networkName: app-network
  autoCreateSubnetworks: false
```

**IP Plan**: `10.10.0.0/16` for dev subnets

### Staging VPC

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpVpc
metadata:
  name: app-network
  id: gcpvpc-staging-app-001
  org: product-team
  env: staging
spec:
  projectId:
    value: myapp-staging-project
  networkName: app-network
  autoCreateSubnetworks: false
```

**IP Plan**: `10.20.0.0/16` for staging subnets

### Production VPC

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpVpc
metadata:
  name: app-network
  id: gcpvpc-prod-app-001
  org: product-team
  env: production
spec:
  projectId:
    value: myapp-prod-project
  networkName: app-network
  autoCreateSubnetworks: false
  routingMode: GLOBAL  # Production may need hybrid connectivity
```

**IP Plan**: `10.30.0.0/16` for prod subnets

### Best Practices

1. **Use the same VPC name** (`app-network`) across environments for consistency
2. **Non-overlapping IP ranges** prevent conflicts when peering or connecting to on-premises
3. **Environment labels** (`env: development|staging|production`) enable cost tracking and resource organization
4. **Global routing in prod** only if needed—keep dev/staging regional for simplicity

---

## Example 4: Shared VPC Host Network

### Use Case

A Shared VPC architecture where a central "host" project owns the VPC, and multiple "service" projects (e.g., different teams or applications) attach to it. This pattern centralizes network administration while allowing teams to deploy resources in their own projects.

### Architecture

- **Host project**: Owns the VPC and subnets
- **Service projects**: Attach to the host VPC and deploy resources (GKE, VMs, etc.)
- **Global routing**: Typically enabled for multi-region service projects

### Host VPC Manifest

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpVpc
metadata:
  name: shared-vpc-host
  id: gcpvpc-shared-001
  org: platform-team
  env: shared
spec:
  projectId:
    value: network-host-project-123
  networkName: shared-vpc-host
  autoCreateSubnetworks: false
  routingMode: GLOBAL
```

### What Gets Created

- A VPC in the `network-host-project-123` (designated as a Shared VPC host)
- Global routing mode for multi-region service projects

### Post-Deployment Steps

After creating the VPC, enable Shared VPC and attach service projects:

```bash
# Enable Shared VPC on the host project
gcloud compute shared-vpc enable network-host-project-123

# Attach service projects
gcloud compute shared-vpc associated-projects add service-project-1 \
  --host-project=network-host-project-123

gcloud compute shared-vpc associated-projects add service-project-2 \
  --host-project=network-host-project-123
```

### When to Use Shared VPC

**Use Shared VPC when:**
- Multiple teams or applications need centralized network management
- You want to enforce consistent security policies (firewall rules, IAM) across projects
- Cost tracking requires separate projects, but network isolation is not needed

**Don't use Shared VPC when:**
- Each project needs completely isolated networks (use separate VPCs instead)
- Network administrators are decentralized (Shared VPC requires central control)

---

## Example 5: VPC with Foreign Key Reference

### Use Case

Reference another OpenMCF resource (like a `GcpProject`) instead of hardcoding the project ID. This enables clean dependency management and ensures the VPC is created only after the project exists.

### Architecture

- **Declarative dependencies**: VPC depends on the referenced project
- **Automatic resolution**: OpenMCF resolves the project ID from the `GcpProject` resource

### Project Resource

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpProject
metadata:
  name: myapp-project
spec:
  projectId: myapp-prod-123
  billingAccountId: 012345-ABCDEF-678910
  folderId: "123456789012"
```

### VPC Resource (Referencing the Project)

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpVpc
metadata:
  name: myapp-vpc
spec:
  projectId:
    ref:
      kind: GcpProject
      name: myapp-project  # References the project defined above
  networkName: myapp-vpc
  autoCreateSubnetworks: false
```

### How It Works

1. OpenMCF creates the `GcpProject` resource first
2. Once the project exists, it extracts `status.outputs.project_id`
3. The `GcpVpc` resource uses that project ID automatically
4. Ensures proper dependency ordering (VPC created after project)

### When to Use Foreign Key References

**Use references when:**
- You're managing both the project and VPC via OpenMCF
- You want explicit dependency ordering in your infrastructure code
- You're building reusable templates that reference other resources

**Use direct values when:**
- The project already exists and is managed outside OpenMCF
- You need to deploy the VPC independently without dependency tracking

---

## Example 6: VPC with Private Services Access

### Use Case

A VPC configured for private connectivity to Google managed services (Cloud SQL, Memorystore, Filestore, etc.). This enables private IP addresses for managed databases and caches without exposing them to the public internet.

### Architecture

- **Private Services Access**: VPC peering with Google's service network
- **IP allocation**: Reserved IP range for Google managed services
- **Secure by default**: Managed services accessible only from within the VPC

### Prerequisites

The `servicenetworking.googleapis.com` API must be enabled on the GCP project. This is handled at the project level (not VPC level) to avoid issues when multiple VPCs share the same project:

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpProject
metadata:
  name: my-prod-project
spec:
  projectId: my-prod-project-123
  billingAccountId: 012345-ABCDEF-678910
  enabledApis:
    - compute.googleapis.com
    - servicenetworking.googleapis.com  # Required for Private Services Access
```

### Manifest

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpVpc
metadata:
  name: prod-vpc
  id: gcpvpc-prod-001
  org: platform-team
  env: production
spec:
  projectId:
    ref:
      kind: GcpProject
      name: my-prod-project
  networkName: prod-vpc
  autoCreateSubnetworks: false
  privateServicesAccess:
    enabled: true
    ipRangePrefixLength: 16  # /16 = 65,536 IPs for managed services
```

### What Gets Created

- VPC network in custom mode
- IP range allocated for Google managed services (e.g., `10.100.0.0/16`)
- VPC peering with Google's service network (`servicenetworking.googleapis.com`)

### IP Range Prefix Length

The `ipRangePrefixLength` determines how many IPs are reserved for Google managed services:

| Prefix | IPs Available | Use Case |
|--------|---------------|----------|
| `/24` | 256 | Tiny - just a few instances |
| `/20` | 4,096 | Small - dozens of instances |
| `/16` | 65,536 | **Recommended** - hundreds of instances |
| `/8` | 16.7M | Enterprise - massive scale |

**Default**: `/16` (65,536 IPs) - sufficient for most use cases.

### Next Steps: Deploy Cloud SQL with Private IP

After deploying the VPC with Private Services Access, Cloud SQL can use private IP:

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudSql
metadata:
  name: prod-postgres
  env: production
spec:
  projectId:
    ref:
      kind: GcpProject
      name: my-prod-project
  databaseEngine: POSTGRESQL
  databaseVersion: POSTGRES_15
  tier: db-custom-2-4096
  storageGb: 100
  region: us-central1
  network:
    privateIpEnabled: true
    vpcId:
      ref:
        kind: GcpVpc
        name: prod-vpc  # References the VPC with Private Services Access
```

### Why Private Services Access at VPC Level?

**API enablement belongs at the GcpProject level**, not VPC level, because:
- If VPC A and VPC B both enable `servicenetworking.googleapis.com`
- When VPC A is destroyed, the API might be disabled
- This would break VPC B and all dependent resources

The VPC module:
- **Assumes** `servicenetworking.googleapis.com` is already enabled
- **Only creates** the IP allocation and VPC peering connection

---

## Advanced Pattern: GKE-Ready VPC Foundation

### Use Case

A VPC designed to host Google Kubernetes Engine (GKE) clusters, with subnets that include secondary IP ranges for pods and services.

### VPC Manifest

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpVpc
metadata:
  name: gke-network
  id: gcpvpc-gke-001
  org: platform-team
  env: production
spec:
  projectId:
    value: gke-prod-project
  networkName: gke-network
  autoCreateSubnetworks: false
  routingMode: GLOBAL  # Useful for multi-region GKE clusters
```

### GKE-Ready Subnet (Separate Resource)

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpSubnetwork
metadata:
  name: gke-uswest1-subnet
spec:
  vpcSelfLink:
    ref:
      kind: GcpVpc
      name: gke-network
  region: us-west1
  ipCidrRange: 10.40.0.0/16  # Primary range for node IPs
  secondaryIpRanges:
    - rangeName: gke-pods
      ipCidrRange: 10.41.0.0/16  # Secondary range for pod IPs
    - rangeName: gke-services
      ipCidrRange: 10.42.0.0/16  # Secondary range for service IPs
  privateGoogleAccess: true  # Enable Private Google Access for GKE
```

### Why Separate VPC and Subnet Resources?

- **Modularity**: VPC defines routing behavior; subnets define IP allocation per region
- **Reusability**: One VPC, multiple subnets in different regions
- **Clarity**: Each resource has a single, focused responsibility

### GKE IP Planning Best Practices

- **Primary range**: Node IPs (sized for maximum expected nodes)
- **Pod secondary range**: Pod IPs (sized for `max_pods_per_node * max_nodes`)
- **Service secondary range**: ClusterIP service IPs (sized for expected services)
- **Avoid overlap**: Ensure ranges don't overlap with other VPCs, on-premises, or peered networks

---

## Deployment Commands

### Using Pulumi

```bash
cd iac/pulumi
pulumi stack init dev
pulumi config set gcp:project <your-project-id>
pulumi up
```

### Using Terraform

```bash
cd iac/tf
terraform init
terraform plan -var-file=dev.tfvars
terraform apply -var-file=dev.tfvars
```

---

## Common Pitfalls and How to Avoid Them

### Pitfall 1: Using Auto-Mode VPCs

**What happens:**
```yaml
spec:
  autoCreateSubnetworks: true  # ❌ Don't do this
```

**Result**: GCP creates subnets in all regions with fixed IP ranges (`10.128.0.0/9`). These ranges overlap with other auto-mode VPCs, preventing peering and causing conflicts with on-premises networks.

**Fix**: Always use custom mode (`autoCreateSubnetworks: false`) and create subnets explicitly.

---

### Pitfall 2: IP Range Conflicts

**Scenario**: You create three VPCs (dev, staging, prod) and later want to peer them or connect to on-premises.

**What happens**: If you didn't plan IP ranges, they overlap and peering fails.

**Fix**: Plan non-overlapping IP ranges upfront:
- Dev: `10.10.0.0/16`
- Staging: `10.20.0.0/16`
- Prod: `10.30.0.0/16`
- On-premises: `192.168.0.0/16`

---

### Pitfall 3: Unnecessary Global Routing

**What happens**: You set `routing_mode: GLOBAL` without understanding the implications.

**Result**: Cloud Routers in any region advertise routes globally. If you don't have hybrid connectivity, this adds complexity without benefit.

**Fix**: Use `REGIONAL` (default) unless you have multi-region Cloud VPN/Interconnect.

---

## Next Steps

After deploying your VPC:

1. **Create subnets** in the regions you need (via `GcpSubnetwork` resources)
2. **Configure firewall rules** to control traffic (via `GcpFirewallRule` resources)
3. **Enable Private Google Access** if workloads need to reach Google APIs without public IPs
4. **Set up Cloud NAT** if private instances need outbound internet access
5. **Deploy workloads** (GKE clusters, VMs, Cloud Run services, etc.)

---

## Additional Resources

- [GCP VPC Best Practices](https://cloud.google.com/architecture/best-practices-vpc-design)
- [IP Address Planning for Large-Scale GKE Deployments](https://medium.com/google-cloud/ip-address-planning-for-large-scale-gke-deployments-48fdee0f7722)
- [Shared VPC Overview](https://cloud.google.com/vpc/docs/shared-vpc)
- [Custom Mode vs Auto Mode VPCs](https://cloud.google.com/vpc/docs/vpc#subnet-ranges)

---

For more details, see the [main README](README.md) and [research documentation](docs/README.md).


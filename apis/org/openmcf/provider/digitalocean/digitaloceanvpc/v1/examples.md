# DigitalOcean VPC Examples

This document provides practical examples for deploying DigitalOcean Virtual Private Cloud (VPC) networks using the OpenMCF API. Each example demonstrates different use cases following the **80/20 principle**: auto-generated CIDR blocks for simplicity (80%) vs. explicit IP planning for production control (20%).

## Table of Contents

- [Example 1: Minimal Development VPC (Auto-Generated CIDR)](#example-1-minimal-development-vpc-auto-generated-cidr)
- [Example 2: Staging VPC with Explicit /20 CIDR](#example-2-staging-vp-with-explicit-20-cidr)
- [Example 3: Production VPC with /16 CIDR and Description](#example-3-production-vpc-with-16-cidr-and-description)
- [Example 4: Multi-Environment Setup with Non-Overlapping Ranges](#example-4-multi-environment-setup-with-non-overlapping-ranges)

---

## Example 1: Minimal Development VPC (Auto-Generated CIDR)

**Use Case:** Quick development environment. Let DigitalOcean handle IP address planning automatically.

**Configuration:**
- **Region:** nyc3 (New York)
- **IP Range:** Auto-generated /20 CIDR block (4,096 IPs)
- **Cost:** $0/month (VPCs are free on DigitalOcean)

**80/20 Principle:** This is the **80% use case**. No IP planning required—just specify name and region.

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanVpc
metadata:
  name: dev-vpc
spec:
  region: nyc3
  # ip_range_cidr is intentionally omitted
  # DigitalOcean will auto-generate a non-conflicting /20 block
```

**Deployment:**

```bash
# Save the manifest
cat > dev-vpc.yaml << 'EOF'
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanVpc
metadata:
  name: dev-vpc
spec:
  region: nyc3
EOF

# Deploy using OpenMCF CLI
openmcf apply -f dev-vpc.yaml

# Check the auto-generated IP range
openmcf get digitaloceanvpc dev-vpc -o yaml | grep ip_range
# Example output: ip_range: 10.116.0.0/20
```

**Result:** DigitalOcean automatically assigns a /20 CIDR block (e.g., `10.116.0.0/20`) that doesn't conflict with your existing VPCs.

---

## Example 2: Staging VPC with Explicit /20 CIDR

**Use Case:** Staging environment with controlled IP range to avoid future peering conflicts.

**Configuration:**
- **Region:** sfo3 (San Francisco)
- **IP Range:** 10.100.16.0/20 (explicit, non-overlapping with dev and prod)
- **Description:** Added for team documentation

**20% Use Case:** Explicit IP planning for controlled networking.

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanVpc
metadata:
  name: staging-vpc
spec:
  description: "Staging environment VPC for pre-production testing"
  region: sfo3
  ipRangeCidr: "10.100.16.0/20"
```

**Deployment:**

```bash
cat > staging-vpc.yaml << 'EOF'
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanVpc
metadata:
  name: staging-vpc
spec:
  description: "Staging environment VPC for pre-production testing"
  region: sfo3
  ipRangeCidr: "10.100.16.0/20"
EOF

openmcf apply -f staging-vpc.yaml

# Verify CIDR is exactly as specified
openmcf get digitaloceanvpc staging-vpc
```

---

## Example 3: Production VPC with /16 CIDR and Description

**Use Case:** Large production environment needing maximum IP capacity.

**Configuration:**
- **Region:** nyc1 (New York)
- **IP Range:** 10.101.0.0/16 (65,536 IPs for growth)
- **Description:** Full documentation for production operations
- **Default:** false (explicit)

**Best Practice:** Use /16 for production DOKS clusters to support:
- Large node pools (hundreds of nodes)
- Multiple clusters in same VPC
- Future expansion without migration pain

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanVpc
metadata:
  name: prod-vpc
spec:
  description: "Main production VPC for all services - immutable after creation"
  region: nyc1
  ipRangeCidr: "10.101.0.0/16"
  isDefaultForRegion: false
```

**Deployment:**

```bash
cat > prod-vpc.yaml << 'EOF'
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanVpc
metadata:
  name: prod-vpc
spec:
  description: "Main production VPC for all services - immutable after creation"
  region: nyc1
  ipRangeCidr: "10.101.0.0/16"
  isDefaultForRegion: false
EOF

openmcf apply -f prod-vpc.yaml

# Get VPC ID for use in cluster/database deployments
VPC_ID=$(openmcf outputs digitaloceanvpc prod-vpc | grep vpc_id | awk '{print $2}')
echo "VPC ID: $VPC_ID"
```

**Usage with DOKS Cluster:**

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
    value: "your-prod-vpc-id"  # From VPC outputs above
  defaultNodePool:
    size: s-4vcpu-8gb
    nodeCount: 5
```

---

## Example 4: Multi-Environment Setup with Non-Overlapping Ranges

**Use Case:** Complete dev/staging/prod setup with carefully planned, non-overlapping CIDR blocks for future VPC peering.

**IP Allocation Strategy:**
- Development: 10.100.0.0/20 (4,096 IPs)
- Staging: 10.100.16.0/20 (4,096 IPs)
- Production: 10.101.0.0/16 (65,536 IPs)

No overlap ensures future VPC peering is possible.

### Development VPC

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanVpc
metadata:
  name: dev-vpc
spec:
  description: "Development environment - auto-destroy after hours"
  region: nyc3
  ipRangeCidr: "10.100.0.0/20"
```

### Staging VPC

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanVpc
metadata:
  name: staging-vpc
spec:
  description: "Staging environment - mirrors production architecture"
  region: nyc3
  ipRangeCidr: "10.100.16.0/20"
```

### Production VPC

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanVpc
metadata:
  name: prod-vpc
spec:
  description: "Production environment - HIGH AVAILABILITY"
  region: nyc3
  ipRangeCidr: "10.101.0.0/16"
  isDefaultForRegion: true
```

**Deployment Sequence:**

```bash
# Deploy in order (VPC-first principle)
openmcf apply -f dev-vpc.yaml
openmcf apply -f staging-vpc.yaml
openmcf apply -f prod-vpc.yaml

# Verify no overlapping CIDRs
openmcf get digitaloceanvpc -o yaml | grep -E 'name:|ipRangeCidr:'

# Expected output:
# name: dev-vpc
# ipRangeCidr: 10.100.0.0/20
# name: staging-vpc
# ipRangeCidr: 10.100.16.0/20
# name: prod-vpc
# ipRangeCidr: 10.101.0.0/16
```

**VPC Peering (Optional):**

If you need to connect staging and production for testing:

```bash
# Create VPC peering (use doctl or Terraform)
doctl vpcs peerings create \
  --name staging-to-prod \
  --vpc-ids $(openmcf outputs digitaloceanvpc staging-vpc | grep vpc_id | awk '{print $2}'),$(openmcf outputs digitaloceanvpc prod-vpc | grep vpc_id | awk '{print $2}')

# Verify peering
doctl vpcs peerings list
```

---

## Common Operations

### Get VPC Status

```bash
# Get detailed VPC information
openmcf get digitaloceanvpc <vpc-name> -o yaml

# Get VPC outputs (ID, URN, IP range)
openmcf outputs digitaloceanvpc <vpc-name>
```

### Get Auto-Generated IP Range

```bash
# For VPCs with auto-generated CIDR
VPC_CIDR=$(openmcf get digitaloceanvpc dev-vpc -o json | jq -r '.status.outputs.ip_range_computed')
echo "Auto-generated CIDR: $VPC_CIDR"
```

### Use VPC ID in Other Resources

```bash
# Get VPC ID
VPC_ID=$(openmcf outputs digitaloceanvpc prod-vpc | grep vpc_id | awk '{print $2}')

# Use in DOKS cluster manifest
cat > cluster.yaml << EOF
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanKubernetesCluster
metadata:
  name: my-cluster
spec:
  vpc:
    value: "$VPC_ID"
  # ... other config
EOF
```

### Delete VPC

**Warning:** VPC deletion fails if resources still exist in the VPC. Delete all Droplets, DOKS clusters, databases, and load balancers first.

```bash
# List resources in VPC
doctl vpcs resources get <vpc-id>

# Delete after cleaning up resources
openmcf delete -f my-vpc.yaml

# Or delete by name
openmcf delete digitaloceanvpc <vpc-name>
```

---

## Best Practices

### IP Address Planning

**80% Use Case (Recommended for Dev/Test):**
- Omit `ipRangeCidr` field
- Let DigitalOcean auto-generate /20 block
- Fastest deployment, zero planning overhead

**20% Use Case (Production with IPAM):**
- Use /16 for production (65,536 IPs)
- Use /20 for staging/dev (4,096 IPs)
- Plan non-overlapping ranges for future peering
- Document IP allocation in description field

**T-Shirt Sizing Strategy:**

| Environment | Recommended CIDR | IPs Available | Use Case |
|-------------|------------------|---------------|----------|
| **Dev/Test** | Auto-generated /20 | ~4,000 | Quick experimentation |
| **Small Prod** | 10.x.0.0/20 | ~4,000 | Small production workloads |
| **Medium Prod** | 10.x.0.0/18 | ~16,000 | Growing production |
| **Large Prod** | 10.x.0.0/16 | ~65,000 | DOKS clusters, databases, scaling |

### VPC-First Deployment

**Critical Rule:** Create VPCs **before** creating:
- DOKS clusters
- Load balancers
- Managed databases (can be migrated, but disruptive)

**Anti-Pattern:**
```bash
# ❌ WRONG: Create cluster first, then VPC
openmcf apply -f cluster.yaml
openmcf apply -f vpc.yaml

# Now cluster cannot be migrated to VPC!
```

**Correct Pattern:**
```bash
# ✅ RIGHT: VPC-first
openmcf apply -f vpc.yaml
# Wait for VPC creation (~30 seconds)
openmcf apply -f cluster.yaml
```

### Avoiding CIDR Conflicts

**Reserved Ranges (Do NOT Use):**
- `10.244.0.0/16` - Reserved by DigitalOcean
- `10.245.0.0/16` - Reserved by DigitalOcean
- `10.246.0.0/24` - Reserved by DigitalOcean
- `10.229.0.0/16` - Reserved by DigitalOcean
- `10.10.0.0/16` in nyc1 - Region-specific reservation

**Avoid Common Home Network Ranges:**
- `192.168.0.0/24`, `192.168.1.0/24` - Common router defaults
- Causes VPN and SSH tunnel conflicts

**Recommended Allocation:**
- Use `10.100.0.0/16` - `10.199.0.0/16` range
- Subdivide for environments within this space

### Multi-Region Architecture

```yaml
# US East VPC
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanVpc
metadata:
  name: us-east-vpc
spec:
  region: nyc3
  ipRangeCidr: "10.100.0.0/16"

---
# US West VPC
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanVpc
metadata:
  name: us-west-vpc
spec:
  region: sfo3
  ipRangeCidr: "10.101.0.0/16"

---
# Asia VPC
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanVpc
metadata:
  name: asia-vpc
spec:
  region: blr1
  ipRangeCidr: "10.102.0.0/16"
```

**VPC Peering:** Non-overlapping ranges allow future peering for cross-region communication.

---

## Troubleshooting

### VPC Creation Fails

**Symptom:** VPC creation returns error

**Possible Causes:**
- CIDR block overlaps with existing VPC
- Invalid region
- Invalid CIDR format (must be /16, /20, or /24)
- Using reserved DigitalOcean range

**Solution:**
```bash
# List existing VPCs and their CIDRs
doctl vpcs list --format Name,Region,IPRange

# Choose non-overlapping CIDR
# Verify region is valid
doctl vpcs list-regions
```

### Cannot Delete VPC

**Symptom:** VPC deletion fails with "resources still attached"

**Possible Causes:**
- Droplets still running in VPC
- DOKS clusters attached
- Managed databases in VPC
- Load balancers in VPC

**Solution:**
```bash
# List all resources in VPC
doctl vpcs resources get <vpc-id>

# Delete resources first
# Then delete VPC
openmcf delete digitaloceanvpc <vpc-name>
```

### DOKS Cluster Won't Use VPC

**Symptom:** Created cluster is not in desired VPC

**Cause:** VPC ID not specified during cluster creation

**Solution:**
```bash
# Get VPC ID first
VPC_ID=$(openmcf outputs digitaloceanvpc my-vpc | grep vpc_id | awk '{print $2}')

# Specify in cluster manifest
cat > cluster.yaml << EOF
spec:
  vpc:
    value: "$VPC_ID"
EOF
```

### IP Range Validation Fails

**Symptom:** "invalid IP range" error

**Valid Formats:**
- `/16` - 65,536 IPs (e.g., `10.100.0.0/16`)
- `/20` - 4,096 IPs (e.g., `10.100.16.0/20`)
- `/24` - 256 IPs (e.g., `10.100.16.0/24`)

**Invalid:**
- `/8`, `/12`, `/18`, `/22`, `/28` - Not supported by DigitalOcean
- Non-RFC1918 ranges - Must use 10.0.0.0/8, 172.16.0.0/12, or 192.168.0.0/16

---

## Migration Scenarios

### Migrating Resources to VPC

**Managed Databases (Easy):**
```bash
# Databases can be migrated with zero downtime
# Use DigitalOcean Control Panel:
# Database → Settings → Cluster Datacenter → Edit → Select VPC
```

**Droplets (Requires Downtime):**
```bash
# 1. Power off Droplet
doctl compute droplet-action shutdown <droplet-id>

# 2. Create snapshot
doctl compute droplet-action snapshot <droplet-id> --snapshot-name "pre-vpc-migration"

# 3. Create new Droplet from snapshot in VPC
# 4. Update DNS/firewall rules
# 5. Delete old Droplet
```

**DOKS Clusters (Cannot Migrate):**
- DOKS clusters **cannot be migrated** to VPC
- Must destroy and recreate cluster
- Plan VPC deployment before creating clusters

---

## Advanced Patterns

### Reserved IP Blocks for Services

Plan CIDR allocation within VPC:

```
Production VPC: 10.101.0.0/16

Subdivisions (manual, not enforced by platform):
- 10.101.0.0/20   - DOKS cluster 1
- 10.101.16.0/20  - DOKS cluster 2  
- 10.101.32.0/20  - Droplets (web servers)
- 10.101.48.0/20  - Managed databases
- 10.101.64.0/18  - Reserved for future growth
```

**Note:** DigitalOcean VPCs don't support subnets. This is documentation-only planning.

### Cost Attribution

VPCs are free, but tag them for organization:

```yaml
metadata:
  name: prod-vpc
  tags:
    - env:production
    - team:platform
    - cost-center:engineering
```

---

## Reference Links

- [DigitalOcean VPC Documentation](https://docs.digitalocean.com/products/networking/vpc/)
- [VPC Pricing](https://www.digitalocean.com/pricing) (VPCs are free)
- [RFC1918 Private Address Space](https://tools.ietf.org/html/rfc1918)
- [CIDR Calculator](https://www.subnet-calculator.com/cidr.php)
- [OpenMCF Documentation](https://docs.openmcf.org/)

---

**Important Notes:**

1. **Immutability:** VPC IP ranges cannot be changed after creation. Plan for growth.
2. **VPC-First:** Always create VPCs before DOKS clusters and load balancers.
3. **80/20 Principle:** Use auto-generated CIDR for dev/test, explicit CIDR for production.
4. **No Overlapping:** Plan non-overlapping ranges if you might need VPC peering.
5. **Free Traffic:** All internal VPC traffic is free and unlimited.


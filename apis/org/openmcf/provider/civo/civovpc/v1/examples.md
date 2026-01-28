# CivoVpc Configuration Examples

This document provides practical configuration examples for different scenarios and use cases.

## Table of Contents

- [Development: Auto-Allocated CIDR](#development-auto-allocated-cidr)
- [Staging: Explicit CIDR for Predictability](#staging-explicit-cidr-for-predictability)
- [Production: Single Region](#production-single-region)
- [Production: Multi-Region with Explicit CIDRs](#production-multi-region-with-explicit-cidrs)
- [Multi-Tenant: Isolation in One Region](#multi-tenant-isolation-in-one-region)
- [Complete Environment Stack](#complete-environment-stack)
- [CIDR Planning Reference](#cidr-planning-reference)

---

## Development: Auto-Allocated CIDR

**Use Case:** Quick dev network for testing. Let Civo handle CIDR allocation.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoVpc
metadata:
  name: dev-test-network
spec:
  civoCredentialId: civo-cred-123
  networkName: dev-test-network
  region: LON1
  description: "Development test network"
```

**Characteristics:**
- No `ip_range_cidr` specified—Civo auto-allocates from available address pools
- Simplest possible configuration for quick experimentation
- Not set as default network (`is_default_for_region` defaults to `false`)
- Suitable for throwaway dev environments

**When to Use:**
- Initial platform exploration
- Temporary testing environments
- Non-critical development workloads
- When CIDR planning is not required

**When NOT to Use:**
- Production or staging environments
- When networks might need VPN connectivity later
- When predictable address spaces are required

---

## Staging: Explicit CIDR for Predictability

**Use Case:** Staging environment with planned CIDR for future VPN integration.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoVpc
metadata:
  name: staging-main-network
spec:
  civoCredentialId: civo-cred-123
  networkName: staging-main-network
  region: NYC1
  ipRangeCidr: "10.20.2.0/24"
  description: "Staging environment network (NYC1)"
```

**Rationale:**
- Explicit CIDR (`10.20.2.0/24`) following hierarchical schema:
  - `10` = Civo namespace
  - `20` = NYC1 region ID
  - `2` = Staging environment ID
- Prevents conflicts if we later add production (`10.20.1.0/24`) or dev (`10.20.3.0/24`) in the same region
- Makes it safe to connect networks via self-managed VPN without overlapping address spaces
- Production-like configuration for realistic testing

**When to Use:**
- Pre-production environments that mirror production architecture
- When multiple environments exist in the same region
- When planning for future VPN or mesh networking
- When testing network connectivity patterns

---

## Production: Single Region

**Use Case:** Production network in a single region with explicit CIDR and documentation.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoVpc
metadata:
  name: prod-main-network
spec:
  civoCredentialId: civo-cred-123
  networkName: prod-main-network
  region: NYC1
  ipRangeCidr: "10.20.1.0/24"
  description: "Production network (NYC1)"
```

**Rationale:**
- Production-ready configuration with explicit CIDR
- Hierarchical address scheme (`10.20.1.0/24` = NYC1 region, prod environment ID 1)
- Environment ID `1` reserved for production across all regions
- Clear documentation via `description` field
- No risk of CIDR conflicts with other environments

**Best Practices:**
- Always use explicit CIDRs for production
- Document the network purpose in the `description` field
- Reserve environment ID `1` for production
- Version control this configuration
- Review changes through code review before applying

---

## Production: Multi-Region with Explicit CIDRs

**Use Case:** Production networks in two regions (LON1 and FRA1) for global traffic distribution and high availability.

```yaml
# Production - London
---
apiVersion: civo.openmcf.org/v1
kind: CivoVpc
metadata:
  name: prod-lon1-network
spec:
  civoCredentialId: civo-cred-123
  networkName: prod-lon1-network
  region: LON1
  ipRangeCidr: "10.10.1.0/24"
  description: "Production network (London)"

# Production - Frankfurt
---
apiVersion: civo.openmcf.org/v1
kind: CivoVpc
metadata:
  name: prod-fra1-network
spec:
  civoCredentialId: civo-cred-123
  networkName: prod-fra1-network
  region: FRA1
  ipRangeCidr: "10.30.1.0/24"
  description: "Production network (Frankfurt)"
```

**Rationale:**
- Separate networks per region (Civo's architectural constraint)
- Non-overlapping CIDRs using hierarchical schema:
  - LON1: `10.10.1.0/24` (region ID 10, environment ID 1)
  - FRA1: `10.30.1.0/24` (region ID 30, environment ID 1)
- Both use environment ID `1` for production
- Enables future multi-region connectivity via self-managed VPN or mesh networking (Tailscale, WireGuard)

**Use Cases:**
- Global application deployment for low-latency access
- High availability across geographic regions
- Disaster recovery with active-active or active-passive patterns
- Data sovereignty requirements (EU data in FRA1, US data in NYC1)

**Connectivity Considerations:**
- No native inter-region peering in Civo
- For cross-region communication, use:
  - Public IPs with TLS encryption
  - Self-managed VPN (WireGuard, OpenVPN)
  - Mesh networking (Tailscale, Nebula)
- Non-overlapping CIDRs enable any of these patterns

---

## Multi-Tenant: Isolation in One Region

**Use Case:** Isolate different end-customers from each other (e.g., SaaS provider hosting multiple client applications).

```yaml
# Customer A
---
apiVersion: civo.openmcf.org/v1
kind: CivoVpc
metadata:
  name: customer-a-net
spec:
  civoCredentialId: civo-cred-123
  networkName: customer-a-net
  region: FRA1
  ipRangeCidr: "10.30.1.0/24"
  description: "Customer A production network"

# Customer B
---
apiVersion: civo.openmcf.org/v1
kind: CivoVpc
metadata:
  name: customer-b-net
spec:
  civoCredentialId: civo-cred-123
  networkName: customer-b-net
  region: FRA1
  ipRangeCidr: "10.30.2.0/24"
  description: "Customer B production network"

# Customer C
---
apiVersion: civo.openmcf.org/v1
kind: CivoVpc
metadata:
  name: customer-c-net
spec:
  civoCredentialId: civo-cred-123
  networkName: customer-c-net
  region: FRA1
  ipRangeCidr: "10.30.3.0/24"
  description: "Customer C production network"
```

**Rationale:**
- Complete Layer 3 isolation between customers
- Each customer gets their own network in the same region
- Prevents cross-tenant access (security requirement for SaaS)
- Sequential CIDR blocks for organizational clarity (`10.30.1.0/24`, `10.30.2.0/24`, `10.30.3.0/24`)

**Use Cases:**
- SaaS providers hosting multiple client applications
- Digital agencies managing client workloads
- Managed service providers (MSPs)
- Platform-as-a-Service offerings

**Security Benefits:**
- True network isolation (Civo's OpenStack-based networking)
- No risk of Customer A accessing Customer B's resources
- Separate firewall rules per customer network
- Clear audit trail per customer environment

**Scaling Considerations:**
- Civo's /24 limit = 256 IPs per customer network
- Sufficient for most Kubernetes clusters (3-10 nodes + load balancers)
- If a customer needs more than 256 host IPs, create multiple networks
  - Example: `customer-a-net-1` (10.30.1.0/24), `customer-a-net-2` (10.30.2.0/24)

---

## Complete Environment Stack

**Use Case:** Full environment stack (dev, staging, prod) in a single region with proper CIDR planning.

```yaml
# Development
---
apiVersion: civo.openmcf.org/v1
kind: CivoVpc
metadata:
  name: dev-nyc1-network
spec:
  civoCredentialId: civo-cred-123
  networkName: dev-nyc1-network
  region: NYC1
  ipRangeCidr: "10.20.3.0/24"
  description: "Development environment (NYC1)"

# Staging
---
apiVersion: civo.openmcf.org/v1
kind: CivoVpc
metadata:
  name: stage-nyc1-network
spec:
  civoCredentialId: civo-cred-123
  networkName: stage-nyc1-network
  region: NYC1
  ipRangeCidr: "10.20.2.0/24"
  description: "Staging environment (NYC1)"

# Production
---
apiVersion: civo.openmcf.org/v1
kind: CivoVpc
metadata:
  name: prod-nyc1-network
spec:
  civoCredentialId: civo-cred-123
  networkName: prod-nyc1-network
  region: NYC1
  ipRangeCidr: "10.20.1.0/24"
  description: "Production environment (NYC1)"
```

**Rationale:**
- Three isolated networks in the same region (NYC1)
- Hierarchical CIDR scheme:
  - `10` = Civo namespace
  - `20` = NYC1 region ID
  - `1` = Production (highest priority)
  - `2` = Staging
  - `3` = Development
- Environment ID ordering (prod = 1, stage = 2, dev = 3) is a convention for consistency
- Clear separation prevents dev mistakes from affecting prod

**Best Practices:**
- Always create separate networks for each environment
- Never mix dev/staging/prod in the same network (security anti-pattern)
- Document environment purpose in `description` field
- Use consistent CIDR schemes across all regions

**Deployment Workflow:**
1. Create networks (one-time setup)
2. Deploy Kubernetes clusters or instances into each network
3. Attach appropriate firewalls to each network
4. Apply Kubernetes Network Policies within clusters

**Zero Egress Fees Benefit:**
- Civo doesn't charge for data transfer between networks
- Encourages best-practice isolation without financial penalty
- Unlike AWS/GCP where multi-VPC architectures incur inter-AZ/inter-VPC costs

---

## CIDR Planning Reference

### Recommended Hierarchical Schema

```
10.A.B.0/24
  ├─ A = Region ID (00-99)
  │   ├─ 10 = LON1 (London)
  │   ├─ 20 = NYC1 (New York)
  │   ├─ 30 = FRA1 (Frankfurt)
  │   ├─ 40 = PHX1 (Phoenix)
  │   └─ 50 = SYD1 (Sydney)
  │
  └─ B = Environment / Tenant ID (00-99)
      ├─ 1 = Production (reserved for production across all regions)
      ├─ 2 = Staging
      ├─ 3 = Development
      └─ 10-99 = Customer/Tenant IDs (for multi-tenant scenarios)
```

### Complete CIDR Allocation Table

| Network | CIDR | Region | Environment | Use Case |
|---------|------|--------|-------------|----------|
| `prod-lon1-network` | `10.10.1.0/24` | LON1 | Production | Production (London) |
| `stage-lon1-network` | `10.10.2.0/24` | LON1 | Staging | Staging (London) |
| `dev-lon1-network` | `10.10.3.0/24` | LON1 | Development | Dev (London) |
| `prod-nyc1-network` | `10.20.1.0/24` | NYC1 | Production | Production (New York) |
| `stage-nyc1-network` | `10.20.2.0/24` | NYC1 | Staging | Staging (New York) |
| `dev-nyc1-network` | `10.20.3.0/24` | NYC1 | Development | Dev (New York) |
| `prod-fra1-network` | `10.30.1.0/24` | FRA1 | Production | Production (Frankfurt) |
| `stage-fra1-network` | `10.30.2.0/24` | FRA1 | Staging | Staging (Frankfurt) |
| `dev-fra1-network` | `10.30.3.0/24` | FRA1 | Development | Dev (Frankfurt) |
| `customer-a-prod` | `10.30.10.0/24` | FRA1 | Customer A | Multi-tenant: Customer A |
| `customer-b-prod` | `10.30.11.0/24` | FRA1 | Customer B | Multi-tenant: Customer B |
| `customer-c-prod` | `10.30.12.0/24` | FRA1 | Customer C | Multi-tenant: Customer C |

### Region ID Mapping

Civo currently supports these regions (assign IDs as needed):

| Region Code | Region Name | Suggested Region ID | Example CIDR Base |
|-------------|-------------|---------------------|-------------------|
| `LON1` | London 1 | 10 | `10.10.x.0/24` |
| `NYC1` | New York 1 | 20 | `10.20.x.0/24` |
| `FRA1` | Frankfurt 1 | 30 | `10.30.x.0/24` |
| `PHX1` | Phoenix 1 | 40 | `10.40.x.0/24` |
| `SYD1` | Sydney 1 | 50 | `10.50.x.0/24` |

### RFC1918 Private Address Spaces

Civo networks use RFC1918 private address spaces. Available ranges:

- `10.0.0.0/8` (16,777,216 addresses) - **Recommended for Civo**
- `172.16.0.0/12` (1,048,576 addresses)
- `192.168.0.0/16` (65,536 addresses)

**Why use 10.0.0.0/8 for Civo?**
- Largest address space (plenty of room for growth)
- Easy to distinguish from on-premise networks (which often use 192.168.x.x)
- Industry standard for cloud networking

### Planning for VPN Connectivity

If you plan to connect Civo networks via VPN (WireGuard, Tailscale, OpenVPN):

1. **Ensure non-overlapping CIDRs** (this schema guarantees it)
2. **Document your allocation** in version control
3. **Reserve address space** for VPN endpoints (e.g., `10.0.1.0/24` for VPN gateway networks)
4. **Test connectivity** in staging before production rollout

### Scaling Beyond 256 IPs per Network

If a single environment needs more than 256 host IPs (Civo's /24 limit):

**Option 1: Multiple Networks per Environment**
```yaml
# Production - Network 1
---
apiVersion: civo.openmcf.org/v1
kind: CivoVpc
metadata:
  name: prod-nyc1-network-1
spec:
  networkName: prod-nyc1-network-1
  region: NYC1
  ipRangeCidr: "10.20.1.0/24"

# Production - Network 2
---
apiVersion: civo.openmcf.org/v1
kind: CivoVpc
metadata:
  name: prod-nyc1-network-2
spec:
  networkName: prod-nyc1-network-2
  region: NYC1
  ipRangeCidr: "10.20.11.0/24"
```

**Option 2: Kubernetes Overlay Networking**
- Kubernetes pods don't consume VPC IPs (they use overlay networks)
- A single /24 Civo network can support clusters with thousands of pods
- Only nodes, load balancers, and standalone instances consume VPC IPs

---

## Common Patterns Summary

| Pattern | CIDR Strategy | Use Case |
|---------|---------------|----------|
| **Auto-Allocated** | Omit `ip_range_cidr` | Quick dev/test environments |
| **Single Environment** | Explicit CIDR (e.g., `10.20.1.0/24`) | Production in one region |
| **Multi-Environment** | Hierarchical (env ID 1, 2, 3) | Dev, staging, prod in same region |
| **Multi-Region** | Hierarchical (region ID 10, 20, 30) | Global deployment |
| **Multi-Tenant** | Sequential in same region | SaaS, MSP, agency workloads |
| **Mixed** | Region + Environment IDs | All of the above combined |

---

## Next Steps

- Review [README.md](./README.md) for architectural details and best practices
- Read [docs/README.md](./docs/README.md) for comprehensive deployment guide
- Check [Civo documentation](https://www.civo.com/docs/networking/private-networks) for platform details

---

**Remember:**
- Always use explicit CIDRs for production
- Plan your CIDR scheme upfront
- Document your allocation in version control
- Test network changes in staging before production
- Embrace the "many small networks" pattern


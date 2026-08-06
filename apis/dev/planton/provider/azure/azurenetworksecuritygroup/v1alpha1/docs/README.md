# AzureNetworkSecurityGroup -- Research & Design Documentation

## Overview

Azure Network Security Groups (NSGs) are stateful, Layer 3-4 packet-filtering firewalls
that control traffic flow for Azure resources. They are the primary network access control
mechanism in Azure, analogous to:

- **AWS Security Groups** -- stateful, but AWS SGs are allow-only (no deny rules)
- **GCP Firewall Rules** -- similar priority-based evaluation, but GCP uses VPC-level rules
- **Traditional firewalls** -- NSGs are per-subnet or per-NIC, not per-network-boundary

NSGs are foundational to Azure networking. Every enterprise Azure deployment uses NSGs
to implement network segmentation, enforce the principle of least privilege, and meet
compliance requirements.

## Azure NSG Architecture

### Rule Evaluation Model

Azure evaluates NSG rules using a priority-based, first-match model:

1. **Direction separation** -- Inbound and outbound rules are evaluated independently
2. **Priority ordering** -- Rules are evaluated from lowest priority number (highest priority)
   to highest priority number (lowest priority)
3. **First match wins** -- The first rule whose 5-tuple matches the traffic determines the
   access decision (Allow or Deny)
4. **Default rules** -- If no user rule matches, Azure's implicit default rules apply

### 5-Tuple Matching

Each rule matches traffic based on:

| Field | Description | Example |
|-------|-------------|---------|
| Source IP | Source address/CIDR/tag | `10.0.1.0/24`, `VirtualNetwork`, `*` |
| Source Port | Source port/range | `*`, `1024-65535` |
| Destination IP | Destination address/CIDR/tag | `10.0.2.0/24`, `Internet` |
| Destination Port | Destination port/range | `443`, `80`, `22` |
| Protocol | Transport protocol | `TCP`, `UDP`, `ICMP`, `AH`, `ESP`, `ANY` |

### Implicit Default Rules

Every NSG has six immutable default rules (three inbound, three outbound):

**Inbound defaults:**

| Priority | Name | Action | Description |
|----------|------|--------|-------------|
| 65000 | AllowVnetInBound | Allow | VNet-to-VNet traffic |
| 65001 | AllowAzureLoadBalancerInBound | Allow | Load Balancer health probes |
| 65500 | DenyAllInBound | Deny | All other inbound traffic |

**Outbound defaults:**

| Priority | Name | Action | Description |
|----------|------|--------|-------------|
| 65000 | AllowVnetOutBound | Allow | VNet-to-VNet traffic |
| 65001 | AllowInternetOutBound | Allow | All outbound internet traffic |
| 65500 | DenyAllOutBound | Deny | All other outbound traffic |

### Stateful Behavior

NSGs are stateful: if inbound traffic is allowed, the return (response) traffic is
automatically permitted regardless of outbound rules, and vice versa. This means you
only need to define rules in one direction for established connections.

### Association Model

NSGs can be associated with:

- **Subnets** -- Rules apply to all resources in the subnet
- **Network Interfaces (NICs)** -- Rules apply to a specific VM/resource

When both subnet-level and NIC-level NSGs exist, traffic is evaluated against both:
inbound traffic hits the subnet NSG first, then the NIC NSG. Outbound traffic hits
the NIC NSG first, then the subnet NSG.

## Deployment Landscape

### Azure Terraform Provider Resources

The AzureRM Terraform provider offers two approaches for NSG rules:

1. **Inline rules** -- `security_rule` blocks within `azurerm_network_security_group`
2. **Separate rules** -- `azurerm_network_security_rule` as standalone resources

**Important:** These two approaches conflict. Using both inline and separate rules for
the same NSG causes Terraform state conflicts. Planton's Terraform module uses separate
rules exclusively, giving each rule its own state identity and plan line.

### Azure Pulumi Provider Resources

The Pulumi Azure Classic SDK (v6) mirrors the Terraform resources:

- `network.NetworkSecurityGroup` -- NSG resource
- `network.NetworkSecurityRule` -- Individual rule resource

Planton's Pulumi module manages rules inline on the group instead: the pinned SDK's
standalone rule resource flattens the application-security-group ID lists to a single
value, while the inline form carries the full lists. Both engines put the identical
rule set into ARM, and removing the last rule removes it on both.

### NSG vs Application Security Groups (ASGs)

Azure also offers Application Security Groups (ASGs) for grouping VMs by application
role. ASGs can be used as source or destination in NSG rules instead of IP addresses --
identity-based addressing that follows workloads as they scale instead of pinning CIDRs.
Rules reference ASGs via `source_application_security_group_ids` and
`destination_application_security_group_ids` as plain ARM IDs (up to 10 each):
application security groups are not yet modeled as a Planton kind.

### NSG Flow Logs

Azure supports NSG Flow Logs for traffic auditing (via Network Watcher). This is
a separate resource (`azurerm_network_watcher_flow_log`) and is not bundled into
the NSG component. Enabling flow logs is an operational concern that can be handled
by a separate Planton component or infra chart configuration.

## Design Rationale

### Why Closed Proto Enums (Not Strings)

The `direction`, `access`, and `protocol` fields are closed protobuf enums
(`INBOUND`/`OUTBOUND`, `ALLOW`/`DENY`, `ANY`/`TCP`/`UDP`/`ICMP`/`AH`/`ESP`) validated
with `defined_only`:

1. **Invalid values cannot reach Azure** -- a typo like `Tpc` is rejected at manifest
   validation time instead of failing mid-deployment
2. **The value set is part of the API contract** -- what the field accepts is visible
   in the schema itself, not buried in a validation expression
3. **The ARM mapping is trivial and contained** -- each IaC module maps enum values to
   ARM's mixed-case strings (`INBOUND` becomes `Inbound`, `ANY` becomes `*`) in one place

### Why Rules Are Folded into the NSG

Security rules live inside the NSG spec rather than as their own Planton resource:

1. **No independent life** -- a rule has no meaning outside its group and is never
   referenced by anything else
2. **Reviewed as one unit** -- a group's rule set is designed and audited together;
   splitting it across resources hides the whole picture
3. **Empty is meaningful** -- an NSG with no rules is still useful: Azure's implicit
   default rules then govern all traffic

### Why No Subnet Association

The NSG-to-subnet association is deliberately excluded from this component:

1. **Independent lifecycle** -- An NSG may be created, reviewed, and approved before
   being associated with a subnet
2. **Reuse** -- The same NSG could potentially be associated with multiple subnets
   (same rules for peer subnets)
3. **Infra chart responsibility** -- The enterprise-network-foundation chart creates
   both NSGs and subnets, then wires the associations as a separate step
4. **Separation of concerns** -- NSG defines "what traffic to allow/deny";
   association defines "where to enforce it"

### Why Rules Have a Description Field

Each rule carries an optional `description` (max 140 characters) because:

1. Azure supports it natively
2. It serves a real operational need -- when reviewing NSGs in Azure Portal or via CLI,
   descriptions explain the intent behind each rule
3. It has zero cost (optional field, no IaC module complexity)
4. Auditors and future operators read descriptions before they read the tuples

### 80/20 Omissions

The following Azure NSG features are deliberately omitted:

| Feature | Reason for Omission |
|---------|---------------------|
| Application Security Groups as a Planton kind | Rules accept plain ARM IDs; a dedicated kind can be added later without breaking changes |
| NSG Flow Logs | Separate operational concern, different lifecycle |

All omissions can be added later without breaking changes.

## Enterprise Deployment Patterns

### Three-Tier Architecture

The most common pattern -- web, app, and data tiers each with their own NSG:

```
Internet
  │
  ▼
[Web NSG] ── Allow 80, 443 from Internet
  │          Deny all other inbound
  ▼
[App NSG] ── Allow 8080 from web subnet
  │          Allow health probes from LB
  │          Deny all other inbound
  ▼
[Data NSG] ── Allow 5432 from app subnet
              Allow 6380 from app subnet
              Allow 22 from mgmt subnet
              Deny all other inbound
```

### Hub-and-Spoke Network

In hub-and-spoke topologies, NSGs control traffic between spokes via the hub:

- **Hub NSG** -- Controls traffic to/from shared services (DNS, NVA, bastion)
- **Spoke NSGs** -- Per-workload rules specific to each spoke's requirements
- **Gateway NSG** -- Controls traffic entering/leaving the Azure environment

### AKS Cluster NSG

AKS clusters require specific NSG rules for control plane communication:

- Allow TCP 443 (API server) from authorized IP ranges
- Allow TCP 9000 (tunnelfront) from AzureCloud service tag
- Allow UDP 1194 (tunnel) from AzureCloud service tag

## References

- [Azure NSG Documentation](https://learn.microsoft.com/en-us/azure/virtual-network/network-security-groups-overview)
- [NSG Default Rules](https://learn.microsoft.com/en-us/azure/virtual-network/network-security-groups-overview#default-security-rules)
- [Service Tags](https://learn.microsoft.com/en-us/azure/virtual-network/service-tags-overview)
- [Terraform azurerm_network_security_group](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/network_security_group)
- [Terraform azurerm_network_security_rule](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/network_security_rule)

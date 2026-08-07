# Network Firewall on OCI

Deploys an Oracle Cloud Infrastructure Network Firewall with an inline firewall policy defining address lists, service definitions, URL pattern lists, and security rules for stateful traffic inspection. The firewall appliance is deployed into a subnet and inspects traffic based on IP addresses, TCP/UDP ports, URL patterns, and optional intrusion detection/prevention. Integrates with Planton's Provider Connections for OCI credential management and ValueFromRef for wiring to compartments, subnets, and security groups.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Network Firewall** -- the firewall appliance in the specified compartment and subnet with configurable static IP, availability domain placement, NAT configuration, and optional NSG bindings
- **Firewall Policy** -- a policy container attached to the firewall that holds all inspection sub-resources
- **Address Lists** -- one per entry in `policy.addressLists`; named collections of IP CIDRs or FQDNs referenced by security rules
- **Services** -- one per entry in `policy.services`; named TCP or UDP port range definitions referenced by security rules
- **Service Lists** -- one per entry in `policy.serviceLists`; named groups of services for reuse across multiple rules
- **URL Lists** -- one per entry in `policy.urlLists`; named URL pattern collections for L7 HTTP(S) traffic inspection
- **Security Rules** -- one per entry in `policy.securityRules`; ordered rules matching traffic against address lists, services, and URL lists with allow/drop/reject/inspect actions
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the firewall and policy

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the firewall in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- A subnet for the firewall appliance. The firewall inspects traffic flowing through this subnet. Provide the subnet OCID directly or reference an OciSubnet Cloud Resource via ValueFromRef. Changing the subnet forces recreation.
- Route table rules directing traffic through the firewall's private IP address. After deployment, configure route rules on other subnets to route traffic via the firewall's `ipv4Address` output.

## Deploy

### Console

Open the deployment store, find **Network Firewall on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Web Perimeter** preset in the [Presets](#presets) tab to pre-populate a firewall allowing HTTP/HTTPS inbound with a default-deny rule.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciNetworkFirewall
metadata:
  name: perimeter-fw
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  subnetId:
    value: "ocid1.subnet.oc1..example"
  policy:
    displayName: perimeter-policy
    addressLists:
      - name: any-ipv4
        type: ip
        addresses:
          - 0.0.0.0/0
    services:
      - name: https
        type: tcp_service
        portRanges:
          - minimumPort: 443
    securityRules:
      - name: allow-https
        action: allow
        condition:
          sourceAddresses:
            - any-ipv4
          services:
            - https
      - name: deny-all
        action: drop
        condition:
          sourceAddresses:
            - any-ipv4
          destinationAddresses:
            - any-ipv4
```

```shell
planton apply -f network-firewall.yaml
```

This creates a firewall with a policy allowing HTTPS traffic and dropping everything else. No URL filtering, IDS/IPS inspection, or NAT configuration is included.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the firewall to a compartment, subnet, and security group deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: security
      fieldPath: status.outputs.compartmentId
  subnetId:
    valueFrom:
      kind: OciSubnet
      name: firewall-subnet
      fieldPath: status.outputs.subnetId
  networkSecurityGroupIds:
    - valueFrom:
        kind: OciSecurityGroup
        name: firewall-nsg
        fieldPath: status.outputs.networkSecurityGroupId
```

The InfraPipeline resolves the dependency graph, deploys the compartment, subnet, and security group first, then provisions the firewall with the resolved values.

## Key Configuration

These are the most important decisions when configuring a network firewall. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Inline policy model** -- The entire firewall policy (address lists, services, service lists, URL lists, security rules) is defined inline in the manifest. Sub-resources reference each other by name. This declarative model means a single `planton apply` provisions the complete firewall configuration atomically.

**Security rule ordering** -- Rules are evaluated in list order; the first matching rule determines the action. Place specific allow/reject rules before broad deny-all rules. Priority is derived from the position in the `securityRules` list (first rule = highest priority). All sub-resource names are immutable after creation.

**Address list types** -- Choose `ip` for CIDR-based matching (e.g., `10.0.0.0/8`, `0.0.0.0/0`) or `fqdn` for domain-based matching (e.g., `malware.example.com`). FQDN lists enable DNS-aware filtering without hard-coding IP addresses.

**Inspection actions** -- Set a rule's `action` to `inspect` with `inspection: intrusion_detection` for passive monitoring or `intrusion_prevention` for active blocking of known attack signatures. Inspection rules require OCI Threat Intelligence integration.

**Immutable placement** -- The `subnetId`, `ipv4Address`, `ipv6Address`, and `availabilityDomain` fields are immutable after creation. Plan subnet placement and IP addressing before deploying.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |
| **OciSubnet** | `subnetId` | `status.outputs.subnetId` |
| **OciSecurityGroup** (optional) | `networkSecurityGroupIds` | `status.outputs.networkSecurityGroupId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `firewallId` | OCID of the network firewall | Monitoring, IAM policy scoping, resource management |
| `ipv4Address` | IPv4 address of the firewall appliance | Route table rules directing traffic through the firewall |
| `policyId` | OCID of the firewall policy | Policy management, audit |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Web perimeter** -- A firewall allowing inbound HTTP/HTTPS and all outbound traffic with a default-deny rule. Address lists separate internal networks from any-IPv4 for clear rule targeting. Start from the **Web Perimeter** preset.

**IDS with URL filtering** -- A firewall with intrusion detection on inbound web traffic and URL pattern blocking for known malicious domains. Combines IDS inspection rules with FQDN-based address lists and URL pattern lists. Start from the **IDS with URL Filtering** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this firewall and its policy
- [**Subnet on OCI**](/cloud-catalog/oci-subnet) -- provides the subnet where the firewall appliance is deployed
- [**Network Security Group on OCI**](/cloud-catalog/oci-security-group) -- provides network security rules applied to the firewall appliance
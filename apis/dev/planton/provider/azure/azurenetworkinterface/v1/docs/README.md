# Azure Network Interface Deployment: The Resource Your VM Actually Plugs Into

## Introduction: The Invisible Resource

Most engineers meet the Azure Network Interface by accident. They create a virtual machine in the portal, and somewhere in the deployment summary a resource named `myvm123` with type "Network interface" appears — created implicitly, named automatically, and promptly forgotten. It works, so nobody looks at it again. Until the day they need a second IP address on the VM, or a static next-hop for a route table, or a security rule that applies to one machine and not its neighbors — and discover that every one of those capabilities lives on the NIC, not the VM.

This is not an implementation detail. In Azure's resource model, the NIC is a **first-class resource**: a VM does not contain its network configuration, it *references* one or more NICs (a required list, of which the first is primary). The NIC owns the subnet placement, the private addresses, the public IP frontings, the per-workload security attachment, the DNS behavior, and the performance features. The VM is a compute shell wired to it.

Treating the NIC as an afterthought works right up until it doesn't. Treating it as the deliberate, independently-managed resource it actually is unlocks patterns — multi-NIC appliances, dual-stack addressing, per-workload firewalling, stable next-hop IPs that survive VM rebuilds — that the implicit model cannot express. This document explains the NIC's anatomy, the decisions that matter in production, and how Planton models the resource declaratively.

## The NIC's Place in Azure's Model

Three relationships define the NIC:

**The VM references the NIC, never contains it.** An Azure VM carries a list of NIC IDs; the first is its primary interface. Because the NIC is independent, it can outlive the VM (delete and recreate the machine, keep the address), and a VM can carry several NICs — a management interface on one subnet and a data interface on another, or the inside/outside arms of a network virtual appliance. The VM size caps how many NICs it can carry (from 2 on small sizes to 8 on large ones).

**The NIC deploys into a subnet.** Each of the NIC's IP configurations places a private address in a subnet, and all of a NIC's configurations must address subnets of the same virtual network. The NIC is therefore the joint where compute meets network: its region must match the virtual network's, and its dependency chain runs Subnet → Virtual Network → Resource Group.

**Traffic-shaping attachments hang off the NIC.** A network security group can attach at the NIC level. Application security group memberships are declared NIC-side. And load-balancer membership is expressed on the NIC, not on the load balancer: an IP configuration declares which backend pools it joins and which inbound NAT rules it completes, each realized as its own association resource.

## Anatomy of an IP Configuration

The unit of addressing on a NIC is the **IP configuration**: a private address in a subnet, optionally fronted by a public IP. Every NIC carries at least one, and most carry exactly one. Each configuration answers four questions:

**Where does the private address live?** A subnet, referenced by ID. IPv4 configurations must name their subnet (ARM's contract); an IPv6 configuration inherits the NIC's subnet placement.

**How is the address assigned?** *Dynamic* allocation lets Azure pick a free address from the subnet at creation — and, crucially, that address is **stable for the NIC's lifetime**. Dynamic does not mean DHCP-style churn; the address only changes if the NIC itself is recreated. This is the right choice for virtually all workloads. *Static* allocation pins an exact address, and is warranted precisely when something *outside* Azure's knowledge depends on the value: a route table's next-hop IP, a partner's firewall rule, an appliance peering configuration.

**Which address family?** IPv4 by default. A dual-stack NIC carries an IPv4 configuration and an IPv6 configuration side by side, in a dual-stack subnet.

**Is there a public face?** A configuration may be fronted by a public IP. The production shape for fleets is *no* — workloads sit behind a load balancer, Application Gateway, or NAT gateway, and only those front doors hold public addresses. Instance-level public IPs are for the workloads that are legitimately one machine: bastion hosts, appliance outside arms, single-tenant demos.

### Why Multiple Configurations?

Two scenarios justify more than one configuration on a NIC:

- **Dual-stack**: an IPv4 and an IPv6 configuration side by side, for workloads that must speak both families natively.
- **Multi-IP**: many IPv4 addresses on one NIC — one per TLS site on a web server, or a pool of source addresses for an appliance's NAT.

When a NIC has multiple configurations, ARM requires exactly one to be primary — and requires the *first* in the list to be it. Getting this wrong is a deploy-time ARM error in the raw API; a well-designed spec surfaces it at validation time instead.

## NIC-Level vs. Subnet-Level Security: The Two-Layer Model

Azure lets a network security group attach at two scopes, and understanding how they compose is essential to not locking yourself out (or leaving a door open):

- **Subnet-level NSG**: attached to the subnet, filters every workload in it. This is the *default posture* — shared rules, managed once, applied uniformly.
- **NIC-level NSG**: attached to one NIC, filters that workload specifically.

When both are attached, they layer, and the ordering matters:

- **Inbound** traffic must pass the **subnet NSG first, then the NIC NSG**. Both must allow it.
- **Outbound** traffic evaluates in the reverse order: NIC NSG first, then subnet NSG.

The practical guidance: rely on subnet-level filtering alone for the common case. Attach a NIC-level NSG when *one workload needs rules its subnet neighbors must not share* — the internet-exposed bastion in a subnet of private machines, the appliance arm with deliberately permissive forwarding rules. And remember the layering when debugging: an inbound flow that the NIC NSG clearly allows can still be dropped by the subnet NSG upstream, and the failure gives no hint which layer denied it.

### Application Security Groups: Rules That Target Roles

Writing NSG rules against IP ranges couples the security policy to the network's address plan — every re-IP, every scale-out, every migration risks breaking a rule. **Application security groups** (ASGs) decouple them: a NIC declares membership in named groups ("web-servers", "databases"), and NSG rules target the groups. "Allow web-servers → databases on 5432" survives any amount of address churn.

ASG membership is declared on the NIC, and — like the NSG attachment — it is its own ARM operation, not a NIC property. Joining or leaving a group never touches the NIC itself.

## Load-Balancer Membership: Declared on the Member

In Azure's model, a load balancer declares its backend *pools*, but never their *members*. Which instances belong to a pool is each member's declaration — a NIC-side association, exactly the shape of the NSG and ASG attachments. This API mirrors that seam with three fields on the IP configuration:

**`loadBalancerBackendAddressPoolIds`** joins the configuration to load-balancer backend pools. The load balancer exports its pools as a name-keyed map, so the reference reads exactly as the pool was named — `valueFrom` fieldPath `status.outputs.backend_pool_ids.web` — instead of a raw ARM ID. A configuration can join several pools, and each membership is its own association resource: scaling a backend fleet in or out is purely NIC-side change, and the load balancer's spec never moves.

**`loadBalancerInboundNatRuleIds`** completes single-target inbound NAT rules. The load balancer declares the port forward (frontend port → backend port); the NIC-side association picks which instance receives it — referenced the same way, e.g. fieldPath `status.outputs.nat_rule_ids.ssh-admin`. This is the per-instance admin-access pattern: the rule lives with the load balancer, the attachment lives with the machine.

**`applicationGatewayBackendAddressPoolIds`** is the Layer 7 counterpart: joining an Application Gateway's backend pool, also realized as association resources. Pools are referenced through the gateway's name-keyed `backend_address_pool_ids` map output (e.g. `status.outputs.backend_address_pool_ids.web`) -- the same member-side seam as the load balancer's.

The division of labor is worth internalizing: the load balancer owns the routing topology (frontends, pools, probes, rules), the NIC owns its own memberships, and the name-keyed output maps are the seam between them. Neither resource ever edits the other.

## Accelerated Networking: On by Deliberate Choice

**Accelerated networking** (SR-IOV, single-root I/O virtualization) lets the NIC bypass the host's virtual switch: packets move between the VM and the physical NIC hardware directly, cutting latency dramatically, raising packets-per-second ceilings, and lowering jitter and CPU overhead. There is no price for it and virtually no downside.

Azure's default is nonetheless **off**, for one reason: not every VM size supports it. The eligibility rule of thumb:

- **Supported**: most current-generation general-purpose and compute-optimized sizes with 2+ vCPUs (and most constrained sizes with 4+).
- **Not supported**: the smallest burstable sizes and some older families.

The constraint is the **VM size, not the workload** — there is no workload that prefers the slow path. The production stance is therefore simple: enable accelerated networking on every NIC destined for a supported size, and treat a NIC without it as a finding. The one operational caveat: enabling or disabling it on an attached NIC requires the VM to be stopped/deallocated, so it is cheapest to get right at creation.

## IP Forwarding and the Appliance Pattern

By default, Azure's virtual network fabric enforces a simple invariant: a NIC only receives traffic addressed to its own IPs. Anything else is dropped at the fabric level — before the VM's operating system ever sees the packet.

**Network virtual appliances** (NVAs) — firewall VMs, routers, SD-WAN concentrators — break this invariant on purpose. The pattern:

1. A route table on the workload subnets sets the next hop to the appliance's private IP (`VirtualAppliance` next-hop type).
2. The fabric delivers workload traffic — addressed to arbitrary destinations — to the appliance's NIC.
3. The appliance inspects, filters, or transforms it, and forwards it onward.

Step 2 only works if the appliance NIC has **IP forwarding enabled**. This is the load-bearing flag of the entire pattern, and its failure mode is vicious: a route table pointing at an appliance whose NIC does not forward **silently blackholes** every routed packet. No error, no log, no NSG deny — traffic simply vanishes. It is one of the classic Azure networking incidents, and it is entirely preventable at provisioning time.

Two supporting choices complete the appliance NIC:

- **Static private address.** Route tables reference the appliance by exact IP. A dynamic address would break every route pointing at it if the NIC were ever recreated.
- **Accelerated networking.** Appliances are packets-per-second workloads; SR-IOV matters more here than anywhere.

Note the flag's scope: IP forwarding disables the *fabric's* destination check. The guest OS must also be configured to route (e.g. `net.ipv4.ip_forward=1` on Linux) — the NIC setting and the OS setting are independent, and both are required.

Conversely: enabling IP forwarding on an ordinary workload NIC is a small but real security regression — it permits the NIC to source and sink traffic it shouldn't. Leave it off everywhere except appliances.

### Auxiliary NVA Acceleration (Preview)

For appliances whose bottleneck is *connection rate* rather than raw bandwidth, Azure offers an auxiliary acceleration path — a **preview feature requiring subscription enrollment**. It comes in two knobs that must be set together:

- **Auxiliary mode**: `AcceleratedConnections` (optimizes connections-per-second), `MaxConnections` (optimizes for very large numbers of simultaneous connections), or `Floating` (floating-IP support on the auxiliary path).
- **Auxiliary SKU**: `A1` through `A8`, sizing the acceleration tier.

Two things to internalize: first, this is for network virtual appliances — an ordinary application NIC has no use for it. Second, because the subscription must be enrolled in the preview, sending these properties on a non-enrolled subscription fails the deployment. The correct default for every non-appliance NIC is to send *nothing*.

## DNS on the NIC: Two Small Levers

The NIC carries two DNS settings, both narrow in scope:

**Custom DNS servers** override the virtual network's DNS settings *for this NIC only*. The strong preference is to configure DNS on the virtual network so every workload inherits it uniformly; the per-NIC override exists for appliances that legitimately need different resolution than their network (a DNS-inspection appliance, a machine that must resolve against a partner's servers).

**Internal DNS name label** gives the NIC a name other VMs in the same virtual network can resolve its private IP by. The label combines with a VNet-internal suffix (surfaced as an output after creation) into a full FQDN. Useful for lightweight intra-VNet addressing without running private DNS zones; leave unset for IP-only addressing.

## Lifecycle: What Replaces, What Updates

Operating a NIC safely requires knowing which changes are destructive:

- **Identity — replaces the NIC**: name, region, and edge zone. Because a replacement NIC is a new resource, changing any of these *detaches the NIC from its VM*. Plan these as maintenance events.
- **Everything else — updates in place**: IP configurations, DNS settings, accelerated networking, IP forwarding, the NSG attachment, ASG memberships, load-balancer and Application Gateway memberships, and tags.

Two more lifecycle facts worth knowing:

- **The MAC address does not exist at creation.** Azure assigns it when the NIC attaches to a *running* VM. Anything keyed on the MAC — license servers, appliance registrations — must wait for the attachment, not the provisioning.
- **All memberships are separate ARM operations**, not NIC properties: the NSG attachment, ASG memberships, load-balancer pool and NAT-rule memberships, and Application Gateway pool memberships. Detaching any of them is a small, independent change that never risks the NIC itself.

## What Is Deliberately Not Modeled (Yet)

A declarative API earns trust by being explicit about its edges.

**Application security groups as a first-class kind**: ASG memberships are accepted today as plain ARM IDs, because the ASG itself is not yet a modeled resource. The membership mechanism (association resources) will not change when it becomes one — only the reference will.

Nothing else notable is withheld. The spec covers the NIC's full production surface: multi-configuration addressing, dual-stack, static pinning, public fronting, both DNS levers, accelerated networking, IP forwarding, the preview auxiliary acceleration, edge-zone placement, gateway load-balancer chaining, NIC-level NSG, ASG memberships, load-balancer pool and NAT-rule memberships, Application Gateway pool memberships, and tags.

## The Planton Approach

Planton provides a declarative, protobuf-based API for deploying Azure Network Interfaces. The design philosophy mirrors Azure's own model honestly: **the NIC is the attachment point**, a first-class resource that everything composes with by reference.

### Composition by Reference

- **The VM references the NIC.** `AzureVirtualMachine`'s `network_interface_ids` consumes this resource's `network_interface_id` output — a required list whose first entry is the primary interface, so a VM can carry several NICs without the API pretending otherwise.
- **Each IP configuration references its subnet** (`subnetId` → an `AzureSubnet`'s `subnet_id` output) **and optionally a public IP** (`publicIpAddressId` → an `AzurePublicIp`'s `public_ip_id` output). The public address stays visible in the resource graph, allowlistable, and reusable if the NIC is replaced.
- **The NIC-level NSG is referenced** (`networkSecurityGroupId` → an `AzureNetworkSecurityGroup`'s output) and realized as an association resource — Azure's own model — so filtering changes never touch the NIC.
- **Load-balancer memberships reference the load balancer's name-keyed maps**: `loadBalancerBackendAddressPoolIds` → `status.outputs.backend_pool_ids.<pool-name>` and `loadBalancerInboundNatRuleIds` → `status.outputs.nat_rule_ids.<rule-name>`, one association resource per membership.
- **ASG memberships** take plain ARM IDs (one association resource each) until application security groups become a modeled kind. **Application Gateway pool memberships** reference the gateway's name-keyed `backend_address_pool_ids` map output.

### Validation Where ARM Would Only Fail at Deploy Time

The spec encodes ARM's contracts as validation rules, so misconfigurations surface before anything reaches Azure:

- When a NIC declares multiple IP configurations, the **first must be marked primary** (and at most one may be).
- **Static allocation requires a pinned address**, and dynamic allocation forbids one.
- **IPv4 configurations require a subnet**; IPv6 configurations inherit the NIC's placement.
- **Auxiliary mode and SKU must be set together** — both or neither.

### Defaults That Match Azure's

Unset enums apply Azure's own defaults — dynamic allocation, IPv4 — so an unspecified spec and Azure's default deploy identically on both provisioning engines. The preview auxiliary properties send *nothing* when unspecified, the correct shape for every non-appliance NIC on any subscription.

### Example: The Standard Workload NIC

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureNetworkInterface
metadata:
  name: app-nic
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: app-rg
  name: app-nic
  ipConfigurations:
    - name: primary
      subnetId:
        valueFrom:
          name: app-subnet
  acceleratedNetworkingEnabled: true
```

This configuration:
- Places one dynamic IPv4 address in the referenced subnet (stable for the NIC's lifetime)
- Enables SR-IOV — right for every supported VM size
- Carries no public IP and no NIC-level NSG: the production posture for a workload behind a load balancer or NAT gateway, filtered at the subnet level

### Example: The Appliance Inside Arm

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureNetworkInterface
metadata:
  name: fw-inside-nic
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: hub-rg
  name: fw-inside-nic
  ipConfigurations:
    - name: primary
      subnetId:
        valueFrom:
          name: fw-inside-subnet
      privateIpAllocation: STATIC
      privateIpAddress: 10.0.254.4
  ipForwardingEnabled: true
  acceleratedNetworkingEnabled: true
```

This configuration:
- Pins the static address route tables reference as their `VirtualAppliance` next hop
- Enables IP forwarding — without it, every routed packet is silently dropped
- Pairs naturally with a public-facing outside NIC on the same appliance VM (the VM references both)

## Conclusion: Make the Attachment Point Explicit

The network interface is where a virtual machine's network life actually happens — its addresses, its subnet presence, its public exposure, its per-workload security, its performance path. Azure's model has always said so; the implicit, portal-generated NIC just made it easy not to notice.

Modeling the NIC deliberately pays off in exactly the situations that matter: the appliance whose next-hop IP must survive a rebuild, the bastion whose firewall rules must not leak to its neighbors, the dual-stack service, the VM whose NIC — and address — outlives it. Planton's API keeps the resource honest: the VM references the NIC, the NIC references its subnet, public IP, security group, and the load-balancer pools it serves, and every ARM contract that would otherwise fail at deploy time is enforced at validation time.

Give your VMs a NIC that was designed, not generated.

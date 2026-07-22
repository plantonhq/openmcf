# Hub-Spoke Network Foundation

The enterprise landing-zone network core for Azure. One deploy produces a
hub virtual network carrying a zone-redundant Azure Firewall, isolated
workload spokes whose only way out is through that firewall, internal DNS
across the whole estate, and the log pipeline that makes every allow/deny
decision inspectable. It is the network a platform team builds first --
and the one every workload that follows plugs into.

## The architecture

Traffic flows through exactly one control point:

- **The hub** is a small, stable network that carries shared
  infrastructure only: the firewall in its ARM-mandated
  `AzureFirewallSubnet`, and a management subnet for operations tooling
  added later. Workloads never deploy into the hub.
- **Each spoke** is an isolated workload network peered to the hub. Its
  workload subnet attaches a route table whose default route
  (`0.0.0.0/0`) forwards to the firewall's private IP -- referenced
  directly from the firewall's output, never hand-copied -- and disables
  BGP route propagation so no future on-premises connection can advertise
  a path around the firewall. Implicit outbound internet access is off;
  the declared route is the only exit.
- **The firewall policy** carries WHAT is enforced, separately from the
  firewall that enforces it. Rules live in a rule collection group per
  concern: this chart ships the platform baseline (Azure Monitor
  telemetry, NTP, and a curated FQDN allowlist for package repositories),
  and each workload team extends the estate by deploying additional rule
  collection groups against the same policy -- never by editing the
  baseline. Rules address spokes through IP Groups (one per spoke), so a
  spoke's address plan changes in one place.
- **Peerings** run in both directions per spoke with forwarded traffic
  allowed, because the firewall relays flows that originate in other
  networks; spoke-to-spoke traffic therefore transits the firewall and is
  denied until a rule allows it.
- **Private DNS** links one internal zone to the hub and every spoke,
  with VM auto-registration on -- machines are resolvable by hostname
  from first boot.
- **Telemetry**: the firewall's full log stream (rule hits, threat
  intelligence, DNS proxy) and the policy's traffic analytics land in a
  Log Analytics workspace deployed with the foundation.

## What is on by default

- **Deny-by-default egress.** Anything not matched by the baseline
  allowlist (or a team's own rules) is dropped and logged. The allowlist
  ships with package mirrors only (`allowed_egress_fqdns`).
- **Threat intelligence in ALERT** (`threat_intelligence_mode`): known-
  malicious flows are logged but not yet blocked -- review the alert
  stream, then raise to `DENY`.
- **Zone-redundant firewall and public IP** (`zone_redundant`): free on
  Azure Firewall and the production posture; disable only in regions
  without availability zones.
- **No implicit outbound anywhere.** Workload and management subnets set
  `defaultOutboundAccessEnabled: false`; egress exists only where a route
  declares it.
- **STANDARD firewall tier** (`firewall_tier`): the full rule engine and
  threat intelligence. `PREMIUM` additionally enables IDPS signature
  inspection (deployed in alert mode); TLS inspection becomes available
  after adding a CA certificate and identity to the policy.

## Parameters worth understanding

- **Address plan** (`hub_vnet_cidr`, `spoke_*_vnet_cidr`, subnet CIDRs):
  the one set of decisions that is expensive to change later -- CIDRs
  must not overlap each other or any network this estate will ever peer
  or VPN to. The firewall subnet must be at least /26.
- **`spoke_2_enabled`**: start with one spoke or two. A third spoke is
  the same stamped pattern -- VNet, NSG, route table, subnet, two
  peerings, an IP Group, a DNS link -- composed from the same first-class
  resources.
- **`allowed_egress_fqdns`**: the curated platform allowlist. Keep it
  minimal; workload-specific destinations belong in that workload's own
  rule collection group.
- **`dns_proxy_enabled` + `spoke_dns_servers`**: enable the firewall's
  DNS proxy when you need FQDN-based *network* rules (non-HTTP protocols
  filtered by name). The proxy only helps once spokes resolve through the
  firewall -- after the first deploy, set `spoke_dns_servers` to the
  firewall's `private_ip_address` output and redeploy.

## After deployment

The firewall takes 10-20 minutes to provision; the full foundation
typically lands in under 30.

- **Verify the egress posture** from any VM in a workload subnet:
  `curl https://packages.microsoft.com` succeeds; a domain outside the
  allowlist times out -- and appears as a deny in the logs.
- **Query the decisions** in the workspace: firewall logs land in
  resource-specific tables (`AZFWApplicationRule`, `AZFWNetworkRule`),
  and policy analytics answers "which rule matched this flow".
- **Hand out the egress IP**: the firewall public IP's `ip_address`
  output is the one address external partners see from this estate.
- **Natural next steps**: deploy workloads into the spoke subnets; add a
  rule collection group per team; add `privatelink.*` DNS zones (zone +
  per-network links, registration off) as private endpoints arrive; and
  when on-premises connectivity comes, land the gateway in the hub --
  the spokes' routing already assumes nothing bypasses the firewall.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).

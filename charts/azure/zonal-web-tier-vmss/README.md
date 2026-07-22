# Azure Zonal Web Tier (VMSS)

The honest IaaS baseline: a Flexible virtual machine scale set spread across three availability zones behind a zone-redundant Standard Load Balancer — SSH-key-only Linux bootstrapped by cloud-init, a user-assigned identity for keyless Azure access, an NSG that admits web traffic and nothing else, explicit (never implicit) SNAT egress, and the CPU alert that says when to scale. Deployed with defaults it serves HTTP from nginx on every instance; replace one parameter (the cloud-init) and it serves your workload.

## Who this is for

A team lifting-and-shifting, running software that needs OS-level control, or simply choosing VMs on purpose. Everyone can start a VM in the portal; the production shape — zonal spread, pool membership wired at NIC creation, probe-gated traffic, outbound rules instead of deprecated implicit egress, no management ports on the Internet — is a day of careful wiring. This chart is that day.

## Architecture

```
                 Internet ── 80/443 ──▶ zone-redundant public IP (zones 1-3)
                                              │
   ┌──────────────────────────────────────────▼───────────────────┐
   │  Standard Load Balancer                                      │
   │   frontend "public" · pool "web" · probe HTTP :80 /          │
   │   rules 80/443 (disableOutboundSnat) · outbound rule "egress"│
   └───────────────┬──────────────────────────────────────────────┘
                   │ NIC-side membership: backend_pool_ids.web
   ┌───────────────▼──────────────────────────────────────────────┐
   │  VNet 10.50.0.0/16 · subnet "web" (NSG: 80/443 in, deny-else)│
   │  FLEXIBLE VMSS · zones 1-3 · Ubuntu 24.04 LTS                │
   │  cloud-init bootstrap · SSH-key-only · UAI attached          │
   └───────────────┬──────────────────────────────────────────────┘
     platform metrics (VMSS + LB) ──▶ Log Analytics
     alert: avg CPU > 80% (15m) ──▶ ops email
```

Design decisions worth knowing:

- **Membership is wired from the member side.** The scale set's NIC references the pool through the load balancer's name-keyed `backend_pool_ids.web` map output — every instance joins the pool at creation, including every future scale-out instance. The LB template deliberately carries no member list.
- **Egress is explicit.** Azure's implicit outbound access is deprecated and unpredictable at scale; the outbound rule gives every instance deterministic SNAT through the tier's own address (1024 ports per instance). The load-balancing rules set `disableOutboundSnat` because Azure rejects a pool where rule-driven SNAT overlaps an outbound rule.
- **No management port exists.** The NSG admits 443 and 80 from the Internet; Azure's implicit rules already allow probe and VNet traffic and deny the rest. SSH happens through Azure Bastion or a peered management network — a deliberate omission, not a gap (inbound NAT to Flexible-mode instances is not modeled NIC-side; Uniform mode owns that pattern).
- **Flexible mode, and what it implies.** Instances are real VMs (individually inspectable, mixable sizes later); fault-domain count 1 is the zonal contract; and identities must be user-assigned — which the chart treats as a feature, since the tier's identity is a first-class node other resources grant roles to.

## Resources

| Kind | Name | Purpose |
| --- | --- | --- |
| AzureResourceGroup | `{env}-web-tier` | One container for the estate |
| AzureLogAnalyticsWorkspace | `{env}-web-tier-logs` | Platform metrics destination |
| AzureUserAssignedIdentity | `{env}-web-tier-identity` | The fleet's keyless credential |
| AzureVirtualNetwork | `{env}-web-vnet` | The tier's network |
| AzureNetworkSecurityGroup | `{env}-web-nsg` | 80/443 in, deny-else (subnet-attached) |
| AzureSubnet | `{env}-web-subnet` | The web subnet, NSG enforced |
| AzurePublicIp | `{env}-web-pip` | Zone-redundant frontend address |
| AzureLoadBalancer | `{env}-web-lb` | Frontend, pool, probe, rules, outbound |
| AzureVirtualMachineScaleSet | `{env}-web-vmss` | The zonal fleet |
| AzureMonitorActionGroup | `{env}-web-tier-ops` | Alert routing |
| AzureMonitorDiagnosticSetting | vmss + lb | Platform metrics into the workspace |
| AzureMonitorMetricAlert | `{env}-web-cpu-alert` | The scale signal |

## Parameters

| Parameter | Description | Default | Must change |
| --- | --- | --- | --- |
| `region` | Azure region (needs availability zones) | `centralus` | |
| `instances` | Fleet size | `3` | |
| `vm_sku` | Instance size | `Standard_D2s_v3` | |
| `admin_username` | Linux admin account | `webadmin` | |
| `ssh_public_key` | The only interactive door (default is a discarded throwaway) | throwaway key | yes |
| `cloud_init_base64` | First-boot bootstrap (default installs nginx) | nginx cloud-init | eventually |
| `vnet_cidr` / `web_subnet_cidr` | Network layout | `10.50.0.0/16` / `10.50.1.0/24` | |
| `ops_email` | Alert recipient | `ops@example.com` | yes |
| `log_retention_days` | Workspace retention | `30` | |

## After deploying

1. **Browse to the tier** — the public IP's `ip_address` output serves nginx's default page over HTTP from first boot (allow a few minutes for cloud-init's package install on each instance).
2. **Replace the bootstrap** — base64 your own `#cloud-config` (installing your app and its TLS termination) into `cloud_init_base64`. New instances pick it up at creation; existing ones recycle on scale-in/out or manual reimage.
3. **Reach an instance when you must** — deploy Azure Bastion into the VNet (or peer a management network); the NSG deliberately exposes no SSH to the Internet, and your `ssh_public_key` is the only accepted credential.
4. **Grant the identity what the app needs** — `{env}-web-tier-identity` is the fleet's credential; add role assignments (storage, Key Vault, databases) against its `principal_id` output instead of shipping secrets in cloud-init.

## Day 2

- **TLS at the instances** — the LB is a pass-through L4; terminate TLS in nginx/caddy via cloud-init (the 443 rule and probe are already in place), or put Application Gateway in front when you need L7 routing and WAF.
- **Guest telemetry** — install the Azure Monitor agent through cloud-init to add syslog and in-guest metrics to the workspace; the platform metrics flowing already cover CPU, disk, and network from the host's view.
- **Autoscale** — the CPU alert tells a human; when the signal is trusted, add an autoscale setting against the scale set so Azure acts on it directly (manage it alongside the chart, then fold it in).
- **Zone drills** — with three zones and three instances, killing one instance simulates a zone loss: the probe pulls it from rotation and the tier keeps serving. Run that drill before production depends on it.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).

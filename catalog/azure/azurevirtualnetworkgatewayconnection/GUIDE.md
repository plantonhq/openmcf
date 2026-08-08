# Azure Virtual Network Gateway Connection -- Operational Guide

Judgment that does not fit in field references: what "working" means,
where tunnels actually fail, and how to structure multi-site
deployments.

## Provisioned is not connected -- internalize this

The single most common confusion with this resource: ARM `Succeeded`
means Azure ACCEPTED the tunnel definition, not that the tunnel is up.
The tunnel reaches `Connected` only when the far side negotiates
successfully. Diagnose the live state with:

```bash
az network vpn-connection show -g <rg> -n <name> --query connectionStatus
```

`Connecting` forever = the two ends disagree. In practice the cause is
one of exactly four things, in this order of likelihood:

1. **Shared key mismatch** (or the on-premises device was never
   configured at all)
2. **IKE version mismatch** (`connectionProtocol` vs the device)
3. **Proposal mismatch** -- the device demands algorithms Azure's
   default set does not offer: pin them via `ipsecPolicy`
4. **The on-premises firewall** blocking UDP 500/4500 or ESP

## One connection per site, forever

Do not be tempted to "reuse" a connection by editing its far side. The
far-side references (`localNetworkGatewayId` is the exception -- it
updates in place; the peer gateway and circuit are ForceNew) and most
identity fields replace the connection. Replacement is cheap (seconds),
but it drops the tunnel -- schedule it.

## Shared keys

- Omit `sharedKey` and Azure generates a strong one -- read it back
  with `az network vpn-connection shared-key show` and configure the
  device with it. This is the least-leaky flow.
- When the device side dictates the key, reference a secret. Never
  paste keys into manifests -- the field is sensitive by contract and
  secret-resolution exists exactly for this.
- Re-keying: update the secret and re-apply; Azure applies it via the
  shared-key API without replacing the connection (brief renegotiation
  drop).

## Custom IPsec proposals

Azure's default proposal set is broad and modern devices negotiate
fine without a pinned policy -- start WITHOUT one. Pin `ipsecPolicy`
only when the device documentation demands exact parameters, and pin
ALL of it (the six algorithms are required together). Two ARM quirks
the spec cannot fully protect you from:

- GCM IKE encryption requires the MATCHING GCM integrity value
  (GCMAES256 with GCMAES256).
- Basic-SKU gateways reject custom policies entirely.

## VNet-to-VNet honestly

Same-region VNets almost always want VNet PEERING instead: cheaper, no
gateway hops, no bandwidth ceiling. Choose VNet-to-VNet tunnels when
you need encryption in transit across regions or transitive routing a
peering mesh cannot express. Remember BOTH gateways need a connection
(this resource, once per direction) with the SAME shared key.

## BGP vs static routing

Static (the default) routes the site's `addressSpaces` -- fine for a
handful of stable sites. Switch to BGP when sites multiply or prefixes
churn: enable it on the gateway, the site description (its
`bgpSettings`), AND this connection (`bgpEnabled`). All three, or
routes silently do not flow.

## NAT opt-in

The gateway owns NAT rules; each tunnel opts in via
`egressNatRuleIds`/`ingressNatRuleIds`. The ids come from the gateway's
`nat_rule_ids` output (a name-to-id map) -- supply them as literals or
explicit references. A rule no connection references does nothing.

## Engine note

The Pulumi engine models exactly ONE `trafficSelectorPolicies` entry
(the classic SDK's shape) and fails loudly on more; multi-selector
connections deploy via the Terraform engine. Most connections need no
selectors at all.

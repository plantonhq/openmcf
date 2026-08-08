# Azure ExpressRoute Circuit -- Operational Guide

Judgment that does not fit in field references.

## The meter starts before the cable exists

Azure bills the circuit from the moment the service key is issued --
NOT from when your provider completes the cross-connect. A circuit
created "to get ahead of procurement" burns its monthly fee while
sitting `NotProvisioned`. Create circuits when the provider order is
actually in flight, and delete abandoned ones promptly.

## The handoff is manual and out-of-band

The `service_key` output is what your connectivity provider needs to
provision the physical link. Nothing in Azure automates that handoff:
you (or your procurement flow) deliver the key to the provider's
portal or NOC, then watch `service_provider_provisioning_state` move
NotProvisioned → Provisioning → Provisioned. Expect days-to-weeks of
provider lead time -- the circuit object itself provisions in minutes.

## Nothing routes until peerings exist

A `Provisioned` circuit still moves no traffic: private peering (for
your VNets) or Microsoft peering (for Microsoft public services) must
be configured on it, and ARM REJECTS peering configuration while the
circuit is unprovisioned. Sequence deployments circuit → provider
handoff → peering → gateway connection.

## Choosing tier and family honestly

- **LOCAL** is the sleeper deal: no egress fees, priced below
  STANDARD, but reaches only the Azure regions in the circuit's own
  metro. Perfect for "our datacenter is across the street from the
  Azure region".
- **STANDARD** covers the geopolitical area -- the default.
- **PREMIUM** buys global reach and larger route tables; buy it for
  the topology, not "just in case" (it is a meaningful price step).
- **UNLIMITED_DATA** beats METERED_DATA at roughly two-thirds
  sustained utilization of the circuit's bandwidth; below that,
  metered is cheaper.

## Bandwidth is a ratchet

ARM grows a provisioned circuit's bandwidth in place but can NEVER
shrink it -- the engines replace the circuit on a decrease, which
means a new service key and a new provider provisioning cycle. Size
at the low end of plausible; growing later is cheap.

## Authorizations are credentials, not metadata

Each `authorizations` entry generates a key that lets a gateway in
ANOTHER subscription consume this circuit's bandwidth. Treat the
`authorization_keys` output like any credential store: issue one
per consuming team/subscription (named for the consumer), and delete
the entry to revoke access when the consumer is decommissioned.

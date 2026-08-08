# Azure ExpressRoute Circuit Peering -- Operational Guide

Judgment that does not fit in field references.

## The circuit must be PROVISIONED first

ARM rejects peering configuration on a circuit whose provider state is
`NotProvisioned` -- the error names the provisioning state, not your
configuration. This is a sequencing fact, not a bug: circuit first,
provider handoff second (days-to-weeks of carrier lead time), peering
third. Watch the circuit's `service_provider_provisioning_state`
output; deploy peerings only after it reads `Provisioned`.

## The type is the identity

A circuit carries AT MOST one peering of each type -- the type string
IS the ARM child's name. Deploying a second "private peering" does not
add one; it would overwrite the first's configuration. Model each
circuit's private peering as exactly one Planton resource, and treat
type changes as delete-and-recreate.

## Session addressing comes from the provider

The VLAN id and the /30 pairs are assigned or confirmed by your
connectivity provider -- they configure the same values on their side
of the cross-connect. Your router takes the first usable address of
each /30, Microsoft's edge the second. A mismatch here is the classic
"BGP never comes up" cause; verify against the provider's handoff
document, not against what deploys cleanly (ARM accepts any coherent
addressing).

## Microsoft peering has a human validation step

`microsoftPeeringConfig.advertisedPublicPrefixes` must be public
prefixes REGISTERED to you (or to `customerAsn`) in an internet
routing registry. Microsoft validates ownership out-of-band: a peering
with unvalidated prefixes deploys but sits in a validation-needed
state until Microsoft support approves it. And remember the route
filter: WITHOUT `routeFilterId` selecting service communities,
Microsoft peering advertises nothing to you at all.

## Global Reach spans subscriptions with keys

`connections` link private peerings across circuits. Same
subscription: just the far peering's ARM id and a non-overlapping /29.
Different subscription: the far circuit's owner issues an
authorization (the circuit kind's `authorizations` list) and you
redeem its key here. The /29 must not overlap either side's address
space -- ARM checks late, so plan it like any other prefix.

## shared_key is write-only

The BGP MD5 key is never returned by ARM. Expect an imported peering
to plan an in-place update on it, and keep the source of truth in your
secret manager -- the deployed state cannot tell you what key is live.

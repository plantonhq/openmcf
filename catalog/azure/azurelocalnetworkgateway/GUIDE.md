# Azure Local Network Gateway -- Operational Guide

Judgment that does not fit in field references.

## It describes; it does not deploy

Nothing about this resource touches the on-premises device. Azure never
probes the endpoint at provisioning time -- a description with a wrong
address deploys `Succeeded` and the failure surfaces later as a tunnel
stuck in `Connecting`. Verify the device facts (public IP, prefixes,
BGP ASN) against the device itself, not against what deploys cleanly.

## Address vs FQDN

Use `gatewayAddress` for anything with a static public IP -- it is the
predictable, debuggable form. Reserve `gatewayFqdn` for sites that
genuinely change addresses (consumer broadband, DHCP-assigned WAN):
Azure re-resolves the name periodically, which means DNS TTL and
propagation become part of your tunnel's failure surface.

## Keep prefixes honest

`addressSpaces` is what Azure ROUTES into the tunnel -- it is a routing
statement, not documentation. Declaring 10.0.0.0/8 "to be safe" routes
all of it at the site, breaking any other tunnel or peering that
carries part of that space. Declare exactly what lives behind the
device; widen deliberately.

Overlaps with the VNet's own space (or another site's) are the classic
multi-site failure: ARM sometimes accepts them and traffic blackholes.
When spaces genuinely overlap, the fix is NAT rules on the virtual
network gateway, not creative prefix declarations.

## BGP sites

- The peering address is the device's TUNNEL-INTERIOR interface IP
  (often an APIPA or loopback), never its public address -- the single
  most common BGP misconfiguration.
- With BGP carrying the routes, leave `addressSpaces` empty so learned
  routes win unambiguously; mixing both is legal but makes route
  provenance hard to reason about during incidents.
- The update path for a BGP-only site has an ARM quirk (a site cannot
  hold both an empty prefix list and no BGP): the provider handles the
  ordering, but expect prefix-list edits on BGP sites to take two ARM
  round-trips.

## Lifecycle

Everything except name/region/resource-group updates in place, and the
object is free -- prefer editing a site's description over
delete-and-recreate, because deleting a description that a connection
references breaks the tunnel until the reference is restored.

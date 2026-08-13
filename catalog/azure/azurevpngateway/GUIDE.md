# Azure VPN Gateway -- Operational Guide

Judgment that does not fit in field references.

## This is the expensive, slow object of the branch family

The site and the connection are free, second-class citizens of the
lifecycle; the GATEWAY bills from creation (~$0.36/hr per scale unit
class) and takes 32-36 minutes to create and 11-13 to delete
(measured live, one scale unit, eastus). Design
everything else around not recreating it: its ForceNew set (hub,
region, routing preference, BGP asn/peer_weight) deserves a decision
BEFORE the first deploy, because changing any of them later replaces
the gateway and every connection on it.

## One gateway per hub -- occupancy, not quota

ARM allows exactly one VPN gateway in a hub. A "second gateway"
request is really a second-hub decision (and hubs are one per region
per WAN). If a create fails with a conflict, look for a
half-deleted predecessor still occupying the hub, not a quota issue.

## Scale units size the pair, not one instance

Each scale unit buys 500 Mbps AGGREGATE across the active-active
pair, and a single tunnel is further capped by its connection's
`bandwidthMbps`. Before buying units, check which limit you are
actually hitting -- one hot tunnel saturates its link cap long before
the gateway's aggregate.

## The ASN is 65515 until Microsoft says otherwise

Virtual WAN VPN gateways currently pin ASN 65515 (the spec lets you
set it because ARM's API does, and Azure may lift the restriction).
Branch devices must NOT use 65515-65520; a branch that insists on the
gateway changing its ASN is a design conversation, not a config file.

## Custom APIPA addresses exist for one reason: interop

The instance_0/instance_1 `customIps` (169.254.21.0-169.254.22.255)
are for far sides that dictate the tunnel-inside addresses -- AWS
site-to-site VPN being the canonical case. Azure applies them AFTER
the gateway exists (they update in place), so adding interop later
never replaces the gateway.

## NAT rules do nothing until a tunnel opts in

A NAT rule on the gateway is inert configuration until a connection
link lists it in `egressNatRuleIds`/`ingressNatRuleIds`. Deploying
rules ahead of branches is safe and cheap; the moment to be careful is
the OPT-IN, which changes live traffic. With BGP branches, remember
`bgpRouteTranslationForNatEnabled` -- untranslated advertisements
re-leak the overlap NAT was hiding.

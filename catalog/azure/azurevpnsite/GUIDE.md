# Azure VPN Site -- Operational Guide

Judgment that does not fit in field references.

## Name links for their meaning, not their order

The link's NAME is the key everything else uses: the `link_ids` output,
the connection's `vpnLinks` references, your runbooks. "primary-isp" /
"backup-isp" survive an ISP swap; "link1" / "link2" invite the wrong
tunnel being debugged at 2am.

## The site is cheap to edit -- the connected link is not

Site fields update in place, but a connection pins each tunnel to a
link's ARM ID and CANNOT be repointed (ARM replaces the tunnel).
Removing or renaming a link that a connection uses is therefore a
connection change first, a site change second. Adding links is always
safe.

## Static prefixes and BGP answer different questions

`addressCidrs` says "route these prefixes into the tunnels" -- for the
whole site. Per-link `bgp` says "learn routes from this speaker." With
BGP on every link, leave `addressCidrs` empty so learned routes win;
mixing both is legal but makes the effective routing the union, which
surprises more often than it helps.

## FQDN endpoints trade stability for a re-resolve window

A link by FQDN survives the branch's public IP changing -- but Azure
re-resolves on its own schedule, so the tunnel is down between the IP
change and the re-resolve. Branches with genuinely static IPs should
use `ipAddress`; FQDN is for consumer-grade uplinks where the IP WILL
change.

## O365 breakout is declarative only

The `o365Policy` categories are read by SD-WAN partner automation --
declaring them changes no Azure routing by itself. Without a partner
device consuming the policy, expect nothing to happen.

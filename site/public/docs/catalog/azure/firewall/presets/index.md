---
title: "Presets"
description: "Ready-to-deploy configuration presets for Firewall"
type: "preset-list"
componentSlug: "firewall"
componentTitle: "Firewall"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-hub-spoke-egress"
    rank: "01"
    title: "Hub-Spoke Egress Firewall"
    excerpt: "This preset creates the production hub firewall: zone-redundant, STANDARD tier, policy-attached, deployed into the hub VNet's `AzureFirewallSubnet` with one Standard static public IP. It is the..."
  - slug: "02-forced-tunneling"
    rank: "02"
    title: "Forced-Tunneling Firewall (Private Data Path)"
    excerpt: "This preset creates a firewall whose DATA path carries no public IP at all -- outbound traffic is forced on-premises (via ExpressRoute/VPN and a 0.0.0.0/0 route on the AzureFirewallSubnet) instead of..."
---

# Firewall Presets

Ready-to-deploy configuration presets for Firewall. Each preset is a complete manifest you can copy, customize, and deploy.

---
title: "Presets"
description: "Ready-to-deploy configuration presets for Bastion Host"
type: "preset-list"
componentSlug: "bastion-host"
componentTitle: "Bastion Host"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-basic-host"
    rank: "01"
    title: "Basic Host"
    excerpt: "This preset deploys the default Bastion shape: dedicated infrastructure at fixed capacity, browser-based RDP/SSH sessions to every VM the network can reach, no public IPs on the machines themselves."
  - slug: "02-standard-tunneling"
    rank: "02"
    title: "Standard with Tunneling"
    excerpt: "This preset deploys a Standard host shaped for engineering teams: native-client sessions from a local terminal (`az network bastion ssh/rdp/tunnel`), file transfer, IP-based connections across..."
  - slug: "03-developer"
    rank: "03"
    title: "Developer Host"
    excerpt: "This preset deploys the free Developer tier: Azure-shared infrastructure attached straight to a virtual network -- no `AzureBastionSubnet` to carve, no public IP to allocate, no hourly bill."
---

# Bastion Host Presets

Ready-to-deploy configuration presets for Bastion Host. Each preset is a complete manifest you can copy, customize, and deploy.

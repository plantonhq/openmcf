# Azure Bastion Host -- Operational Guide

Judgment calls that matter when you run Bastion in production.

## Choose the long-term SKU up front -- downgrades replace the host

Upgrading (Basic -> Standard -> Premium) is an in-place, ~10-minute configuration change. Downgrading has NO path: the host is deleted and recreated, dropping every active session. The cost difference between Standard and Premium is marginal -- when session recording or private-only deployment is plausibly in your future, start at Premium.

## The subnet is a contract, not a suggestion

Dedicated-infrastructure hosts deploy ONLY into a subnet named exactly `AzureBastionSubnet`, /26 or larger, carrying nothing else. ARM enforces the name at deploy time -- the reference resolves after offline validation, so a wrongly-named subnet surfaces as a deploy-time failure, not a manifest error. Carve the subnet when you design the network, not when you need the host.

## The public IP is the host's alone

Bastion binds its Standard static public IP exclusively -- sharing the address with a NAT gateway, load balancer, or gateway fails at deploy. Give the host its own address and treat it as part of the host's lifecycle.

## Developer is a real tier with real limits

The Developer SKU is free and skips the subnet/IP ceremony entirely -- but it is shared infrastructure: one connection per user, no NSG support on the shared path, no virtual-network peering reach, a limited region list, and zero feature knobs. It is the right answer for a developer poking at a dev VM, and the wrong answer the moment two people need concurrent sessions or traffic crosses a peering.

## Plan Kerberos at create time

`kerberos_enabled` is honored at CREATE only -- the provider silently ignores later changes (no update path exists). If domain-joined sign-in is on the roadmap, enable it when the host is born; retrofitting means replacing the host.

## Session hygiene is a feature choice

Clipboard (`copy_paste_enabled`, on by default) and file copy are exfiltration channels as much as conveniences -- compliance environments commonly disable the clipboard and leave file copy off. Shareable links bypass Azure RBAC for the link holder: treat every link as a credential, scope it to one VM, and prefer them off unless a workflow demands them.

## Scale units are a capacity dial, not an afterthought

Each scale unit carries roughly 20 concurrent sessions; Standard/Premium scale 2-50 in place. Watch session counts and scale BEFORE users queue -- scaling is an update, not a replacement, so there is no reason to over-provision on day one.

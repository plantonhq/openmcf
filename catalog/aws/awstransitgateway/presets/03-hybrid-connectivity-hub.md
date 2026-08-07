# Hybrid Connectivity Hub

A hub prepared for on-premises connectivity: a deliberately chosen Amazon-side ASN that will not collide with your data center's BGP, ECMP enabled so parallel VPN tunnels aggregate bandwidth, and cross-VPC security group referencing for precise east-west rules.

## When to Use

- Connecting on-premises networks to multiple VPCs through VPN or Direct Connect
- Topologies where firewall rules should reference security groups across VPCs instead of broad CIDRs

## What It Configures

- **`amazonSideAsn: 64620`** — a non-default private ASN; BGP sessions to your on-premises routers need the two sides to differ
- **`vpnEcmpSupport: true`** — multiple VPN tunnels advertising the same routes load-balance instead of failing over
- **`securityGroupReferencingSupport: true`** — SG rules in one attached VPC can reference groups in another

## What to Customize

- Replace `<aws-region>` and pick an ASN coordinated with your network team (changing it later replaces the gateway)
- VPN connections and Direct Connect gateways attach outside this resource; their attachments join route tables by literal attachment ID

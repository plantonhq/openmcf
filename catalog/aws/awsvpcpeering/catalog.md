# AWS VPC Peering

Private networking between two VPCs — same account or across accounts and regions — over AWS's backbone, with no gateway, VPN, or single point of failure. Non-transitive by design: each pair peers explicitly.

## What Gets Managed

- The request side: your VPC, the peer VPC (with owner account and region for cross-boundary peers), same-account auto-acceptance, and DNS-resolution options so instances resolve each other's private hostnames.
- Or the accept side: adopt a pending connection by its `pcx-` id and accept it, with the accepter's DNS-resolution option.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with EC2 VPC permissions (per side, in cross-account topologies).

### AWS Prerequisites

- Two VPCs with NON-OVERLAPPING CIDR blocks — AWS rejects overlapping peers outright.
- For DNS resolution across the peering: `enableDnsHostnames` on both VPCs.

## After You Deploy

- Same-account auto-accepted peerings go ACTIVE in seconds. Cross-account/cross-region requests sit `pending-acceptance` until the peer side's accept-arm instance runs.
- Peering moves no traffic by itself: add routes toward the peer CIDR in both VPCs' route tables, and open security groups.

## Common Changes

- Enable private DNS across the peering: flip the DNS-resolution options (needs an ACTIVE connection).
- Never re-declare an existing pair: AWS hands back the same `pcx-` id, and two instances then fight over one connection.
- Decommission from the requester: the accept arm's destroy abandons management without deleting; the request arm's destroy is the real delete.

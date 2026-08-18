# AWS Network ACL

The subnet's firewall: ordered allow AND deny rules evaluated at the subnet boundary, stateless in both directions. Security groups say who may talk; network ACLs are where you say who may NOT — block a bad CIDR, quarantine a subnet, enforce tier boundaries.

## What Gets Managed

- The ACL in its VPC, its inbound and outbound rules (rule number, allow/deny, protocol by name or number, IPv4 or IPv6 CIDR, port range, ICMP type/code), and which subnets it filters.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with EC2 VPC permissions.

### AWS Prerequisites

- The VPC (reference an AwsVpc's `vpc_id` output) and the subnets to associate.

## After You Deploy

- Listed subnets are atomically re-parented onto this ACL (a subnet has exactly one ACL — AWS replaces, never attaches). Removing a subnet hands it back to the VPC's default NACL.
- Traffic matching no rule hits AWS's invisible catch-all DENY — an empty ACL blocks everything.

## Common Changes

- Remember the replies: stateless means inbound 443 needs an OUTBOUND ephemeral-port allow (1024–65535) for responses — the classic NACL outage.
- Leave numbering gaps (100, 200, ...) so rules can be inserted without renumbering; lower numbers win.
- Block a range: a deny rule numbered BELOW the allows it must beat.

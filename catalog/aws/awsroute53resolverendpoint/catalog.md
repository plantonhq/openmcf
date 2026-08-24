# AWS Route 53 Resolver Endpoint

Hybrid DNS between AWS and everywhere else: inbound endpoints let on-prem resolvers answer queries INTO your VPCs, outbound endpoints forward VPC queries OUT to on-prem (or any) name servers — with the forwarding rules and their VPC associations managed alongside the endpoint.

## What Gets Managed

- The endpoint: direction (inbound / outbound / inbound-delegation), 2–10 ENI placements across subnets (optionally with pinned IPs), security groups, IP family, DNS transport protocols (Do53 / DoH), and CloudWatch metrics toggles.
- Forwarding rules keyed by name: the domain each rule steers, FORWARD targets (IP, port, protocol per name server), SYSTEM overrides that restore recursive resolution for subdomains, and each rule's VPC associations.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Route 53 Resolver and EC2 (ENI) permissions.

### AWS Prerequisites

- Two or more subnets (different AZs recommended) with free IP capacity for the endpoint ENIs.
- A security group allowing DNS (TCP+UDP 53) between the endpoint and the resolvers or targets it talks to.

## After You Deploy

- The endpoint reaches OPERATIONAL in a few minutes (ENI provisioning). Inbound: point on-prem conditional forwarders at the `ip_addresses` output. Outbound: queries for each rule's domain start forwarding from every associated VPC.
- Endpoint ENIs bill per hour whether or not queries flow — the endpoint is the family's one always-on cost.

## Common Changes

- Add or remove a forwarding rule, retarget a rule's name servers, or associate a rule to more VPCs — all in-place list edits.
- Add ip_addresses for capacity (AWS adds before it removes, so the two-address floor never breaks).
- Changing direction or security groups replaces the endpoint; detaching a rule from its endpoint replaces the rule.

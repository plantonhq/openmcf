# AWS Managed Prefix List

Named CIDR sets as infrastructure: define "our offices" or "partner ranges" once as a prefix list, reference the one `pl-` id from security groups, NACLs, and route tables everywhere — then update the list once when the network changes instead of hunting rules across accounts.

## What Gets Managed

- The prefix list: its address family (IPv4 or IPv6, fixed for life), its capacity (`max_entries`), and its described CIDR entries as the complete set.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with EC2 VPC permissions.

### AWS Prerequisites

- None — prefix lists are standalone and free.

## After You Deploy

- Reference the `prefix_list_id` output from security-group rules and route tables; AWS resolves the current entries at enforcement time.
- Every entry change bumps the `version` output — AWS's audit trail for the list.

## Common Changes

- Network change: edit the entries — every referencing rule follows without touching it.
- Out of room: raise `max_entries` (in place; but each referencing security-group rule then consumes more rule quota — the quota cost is `max_entries`, not actual entries).
- Keep descriptions current: they are the only record of WHY a CIDR is trusted.

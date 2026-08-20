# Web Front-Door IP

This preset reserves an IPv4 address and assigns it to your web droplet, so DNS points at an address that survives droplet replacements. Assigned reservations are free.

## When to Use

- Any droplet serving public traffic whose DNS should never chase droplet churn
- Blue/green droplet swaps: rebuild the standby, then re-point this one field

## Key Configuration Choices

- **Droplet by reference** -- the assignment follows your droplet resource; replacing the droplet in a later apply re-points the address automatically.
- **Same region as the droplet** -- assignment only works intra-region; keep the two manifests' regions aligned.
- **IPv4 (the default)** -- what public DNS A records need.

## What You Get

A stable public address in the `reserved_ip_address` output for your DNS A record -- re-pointable between droplets in place, without the address ever changing.

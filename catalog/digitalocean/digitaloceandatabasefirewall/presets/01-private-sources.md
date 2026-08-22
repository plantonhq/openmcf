# Private Network + Tagged Fleet

This preset trusts exactly two source classes: the private VPC range and every Droplet carrying the `backend` tag. Nothing on the public internet can reach the cluster while this rule set is deployed.

## When to Use

- Production clusters that only internal workloads should reach
- Dynamic Droplet fleets, where tag membership tracks instances automatically
- Replacing the open-by-default posture every new cluster starts with

## Key Configuration Choices

- **CIDR over single IPs** -- the VPC range survives instance churn; single IPs are for fixed bastions.
- **Tag over droplet_ids** -- new tagged Droplets gain access with no manifest change; prefer references (`dropletIds`) only for a small fixed set.
- **No 0.0.0.0/0** -- adding it would deliberately re-admit the internet; leave it out.

## What You Get

A cluster reachable only from the declared sources -- and one standing warning: DESTROYING this resource clears the rule set, after which the cluster accepts connections from anywhere again.

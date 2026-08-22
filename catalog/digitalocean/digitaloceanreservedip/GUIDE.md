# DigitalOcean Reserved IP -- Operational Guide

What experience with this component teaches that the field reference cannot.

## The idle state is the expensive state

Most orphaned cloud resources cost nothing; a forgotten UNASSIGNED reserved IPv4 bills ~$5/month forever. Assignment suspends the charge entirely. If a reservation stops being useful, destroy it -- "keep it around just in case" is the one posture this kind punishes.

## The address is the asset -- and destroy releases it permanently

The whole point of reserving is that the address survives droplet churn: DNS points at the reservation, droplets come and go behind it. The flip side: destroying the reservation releases the address to DigitalOcean's pool, and recreating gets a DIFFERENT one -- every DNS record naming the old address goes stale. Destroy is a DNS event; plan it as one.

## Same region or no assignment

A reserved IP assigns only to droplets in ITS region. Cross-region failover is not what this kind does (that is DNS or a global load balancer); this is intra-region droplet swapping. Pin the reservation's region to where its droplets actually live.

## Re-pointing is the failover move -- but not through the six-phase runner

Assign, re-point, and unassign all apply in place on IPv4 (the address never changes). In practice: run a standby droplet, and failover is a one-field manifest change. Note the v6 flavor re-points by replacing its assignment object -- same one-field edit in the manifest, one extra resource turn under the hood.

## IPv6 delete errors are swallowed upstream

The provider's v6 delete ignores every error except 404 (an inverted error check, verified in source at the pin) -- a failed release can look like a success in IaC output. The proof lane's live absence check is the trustworthy signal; operators cleaning up by hand should verify with `doctl compute reserved-ipv6 list`.

## What is deliberately NOT here

`ip_address` as an input (the provider declares it but never sends it -- the address is an OUTPUT of reservation); the standalone v4 assignment resource (all-ForceNew with a timestamped id that can never round-trip an import -- the reservation's own mutable argument is strictly better); and project placement (the DigitalOceanProject kind's membership list carries the reservation's `urn`).

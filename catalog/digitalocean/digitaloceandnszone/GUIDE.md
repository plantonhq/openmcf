# DigitalOcean DNS Zone -- Operational Guide

Judgment calls that matter when you run DNS zones on DigitalOcean.

## The zone works before the delegation does

Adding a domain to DigitalOcean hosts it instantly on ns1/ns2/ns3.digitalocean.com — queries against those servers answer immediately, which is what the E2E lanes verify. The public internet, though, resolves through whatever nameservers the registrar advertises, so nothing changes for real users until the registrar's NS delegation is updated to DigitalOcean's set (the `name_servers` output) and the old delegation's TTL expires. Plan cutovers in that order: create the zone, populate the records, verify against DigitalOcean's nameservers directly (`dig @ns1.digitalocean.com`), then flip the registrar.

## One zone name per all of DigitalOcean

Domain names are unique across every DigitalOcean account, not just yours. An "already exists" error on create usually means another account (a colleague's personal account, a previous company account) holds the domain — DigitalOcean support, not a retry, is the fix. The reverse also matters: deleting a zone releases the name for anyone to claim.

## Inline records vs. the standalone record kind

This kind's `records` list is for zones whose records ship as one unit under one owner — the company website, a product domain. The standalone `DigitalOceanDnsRecord` kind exists for records with different owners or lifecycles — an application chart adding its hostname to a shared zone. Both write the same DigitalOcean records; pick one home per record so ownership is never split.

## Multi-value entries are real round-robin

An entry with several `values` creates one record per value with the same name and type. For A records that is DNS round-robin — resolvers rotate through the addresses. It is not health-checked failover: a dead address keeps being served. Round-robin is fine for stateless redundancy; use a load balancer when you need health-aware traffic steering.

## Leave the apex NS and SOA alone

Every zone comes with DigitalOcean's own NS records at the apex and one SOA. The API lets you write both types, and there is almost never a reason to: rewriting apex NS records breaks delegation inside the zone, and the SOA is DigitalOcean's operational record. The one legitimate NS use is delegating a subdomain (`sub.example.com`) to other nameservers.

## `ipAddress` seeds a record nobody manages

The create-only `ipAddress` field makes DigitalOcean create an apex A record at zone creation — and that record is invisible to this component forever after: it is not in `records`, later edits do not see it, and deleting the value from the manifest changes nothing live. It exists for migration compatibility. New zones should declare the apex record in `records`, where it is tracked, diffed, and updatable.

## Destroying a zone destroys every record in it

Including records created by the standalone record kind and records added by hand in the control panel. Before destroying a shared zone, enumerate what lives there (`doctl compute domain records list <domain>`) — the blast radius is the whole domain's resolution.

# Azure Private DNS Record -- Operational Guide

Judgment calls that matter when you run private DNS records in production.

## Declare all of a name's values in ONE record, always

Azure stores every value for a (name, type) pair as one record set. Two AzurePrivateDnsRecord resources declaring the same name and type do not merge -- the second conflicts with the first, and on the import path one silently owns values the other believes it manages. When a service gains a second address, edit the existing record's list.

## The record answers nothing until the zone is linked

Records live in the zone; RESOLUTION lives in the zone's virtual-network links. A record that "does not resolve" from some network is almost always a missing AzurePrivateDnsZoneVirtualNetworkLink for that network, not a record problem. Check links first, records second.

## Watch for collisions with auto-registered VM records

Zone links with auto-registration enabled write A records for every VM in the linked network, managed by the service. A declared record with the same name as a registered VM hostname fights the registration lifecycle. Keep declared names out of the auto-registration namespace (or keep registration links separate from resolution links).

## Name, zone, and type are fixed -- plan renames as add-then-remove

Changing a record's name, zone, or payload type replaces the record set -- a resolution gap for that name while the apply runs. For a rename, ADD the new record first, then remove the old one in a second apply; both names answer during the transition and caches drain on their own schedule.

## Size the TTL to the record's change cadence

TTL is how long resolvers cache an answer -- and therefore how long clients keep using an address after you change it. Failover-sensitive names (database primaries, active/passive pairs) want 60 or lower; stable infrastructure names are fine at 3600+. The platform default is 300. Remember the cache drains from the moment of the change, not the deploy: a TTL lowered in the same apply as an address change does not speed that change up.

## Private DNS has no alias records -- reference outputs instead

The public-DNS alias mechanism (records tracking an Azure resource's address) does not exist in private zones. The equivalent pattern here: reference the owning component's output (`a` values from a NIC's private IP output, CNAME values from a service's FQDN output) so redeployments flow the new value through the reference at the next apply.

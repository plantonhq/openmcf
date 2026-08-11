# GcpAddress Guide

The judgment this guide protects: an address reservation is a promise to
the outside world. Everything except labels recreates it, and a recreated
EXTERNAL address is a DIFFERENT IP — reserve once, reference everywhere,
and guard the destroy path.

## Reserve first, attach second

Create the reservation as its own node before the thing that uses it
(Cloud NAT, a regional forwarding rule, a VM). Letting a consumer
auto-allocate an ephemeral IP and "promoting" it later is how teams end up
with addresses no manifest owns. The reservation's `address` output is the
value DNS records and allow-lists should reference.

## The recreate trap

Every field except `labels` is ForceNew. Renaming the address, switching
its type, or adjusting its purpose destroys and re-reserves — and for
EXTERNAL addresses the replacement IP is whatever GCP hands out next.
Anything pinned to the old IP outside GCP (registrar DNS, a partner's
firewall) breaks silently. Treat any plan that shows this resource being
replaced as a production event, not a refactor.

## Purpose selects the wiring, the CELs enforce it

`GCE_ENDPOINT` and `DNS_RESOLVER` live in a subnetwork; `VPC_PEERING` and
`IPSEC_INTERCONNECT` live in a network with a `prefixLength`. The spec
rejects mismatched combinations before deploy, so a validation failure
here means the design is wrong, not the syntax. PSC addresses are
global-only — that workflow belongs to GcpGlobalAddress.

## BYOIP is an ownership statement

`ipCollection` reserves out of a customer-owned PublicDelegatedPrefix
instead of Google's pool. It only makes sense when the organization
already runs a BYOIP program (the PDP must exist and be in the right
mode); for everyone else the default pool is correct. It is create-time
only — moving an address between collections is a re-reservation.

## Teardown discipline

`DELETE` releases the IP back to Google — for an EXTERNAL address that is
irreversible in practice (someone else can claim it). `PREVENT` suits any
address that external systems point at: destroy fails instead of
releasing. `ABANDON` keeps the reservation (and its idle-address billing)
while dropping it from management — the right escape hatch when handing
an IP to another team's tooling. GCP refuses to release an address a
resource is still using, but PREVENT also covers the window after the
consumer is gone.

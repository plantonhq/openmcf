# GcpCertManagerDnsAuthorization Guide

The judgment this guide protects: a DNS authorization is the piece that
makes certificate issuance independent of serving traffic. It looks like
a footnote next to the certificate, but it is the node that enables
zero-downtime migration and wildcard issuance — and the one whose quiet
deletion breaks renewals months later.

## The zero-downtime pattern, in order

1. Create the authorization for the domain.
2. Compose a `GcpDnsRecord` from its outputs — `dns_record_name`,
   `dns_record_type`, `dns_record_data` are exactly the CNAME the zone
   needs.
3. Once the record propagates, create the `GcpCertManagerCert`
   referencing this authorization — it validates and issues WITHOUT the
   domain ever pointing at the new load balancer.
4. Cut DNS over to the new load balancer only when the certificate is
   ACTIVE.

The certificate never gates on serving traffic — the gap that makes
classic managed certificates a migration hazard simply does not exist.

## One authorization per domain, wildcard included

Authorizing `example.com` issues certificates for both `example.com`
and `*.example.com` — do not create a second authorization for the
wildcard. Subdomains are separate: `api.example.com` needs its own
authorization (or its certificate rides the parent's wildcard).

## FIXED_RECORD vs PER_PROJECT_RECORD

`FIXED_RECORD` (the global-location default) is the classic one
record per authorization. `PER_PROJECT_RECORD` scopes the record per
(domain, project) so multiple projects can validate the same domain
with one shared CNAME — the multi-project platform pattern, and the
default for non-global locations. The choice is ForceNew, and GCP picks
a location-dependent default when unset — set it explicitly in anything
multi-project.

## The authorization must outlive issuance

Renewal re-validates through the same CNAME. Deleting the authorization
(or its DNS record) breaks nothing immediately — certificates keep
serving — but renewals for domains not yet serving through the load
balancer stop validating, which surfaces as an expiry incident months
later. Give authorizations under live certificates
`deletionPolicy: PREVENT`; the validation chain is exactly the kind of
resource that gets swept in a cleanup because nothing visibly depends
on it day to day.

## Conventions and gotchas

- Location must match the certificates it validates: global
  authorizations for global certificates, regional for regional.
- The domain is bare — no `*.` prefix, no trailing dot; the wildcard is
  implicit.
- Labels ride beneath the platform's attribution labels, identically on
  both engines.

## Pairs well with

- `GcpDnsRecord` — carries the validation CNAME into the zone.
- `GcpDnsZone` — the zone that record lives in.
- `GcpCertManagerCert` — the consumer; its managed arm references this
  kind by `authorization_id`.

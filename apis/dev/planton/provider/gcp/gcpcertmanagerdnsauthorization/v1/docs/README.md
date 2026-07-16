# GcpCertManagerDnsAuthorization: Design Notes

## What This Component Models

One `google_certificate_manager_dns_authorization` — Certificate Manager's
standalone proof-of-domain-control resource. It is first-class in the API:
it exists independently of any certificate, multiple certificates
reference the same authorization, and it survives certificate rotation.
The component models exactly that node, so the certificate-issuance flow
composes from explicit, individually-ownable pieces:

```
GcpCertManagerDnsAuthorization ──(dns_record_* outputs)──▶ GcpDnsRecord
        ▲
        │ (authorization_id)
GcpCertManagerCert.managed.dns_authorizations
```

## Why the Validation Record Is an Output, Not a Side Effect

The authorization's only observable artifact is a CNAME record GCP asks
you to serve. Creating that record from inside this component would mean
this kind silently writes into a zone that `GcpDnsRecord` owns — a hidden
cross-kind write with ordering and drift hazards, and it would force this
kind to know which zone is authoritative (it cannot: the zone may live in
another project or another provider). Exporting the record's three fields
keeps ownership honest and lets the record ride the same reference
machinery as everything else.

## Domain Semantics

An authorization is for ONE bare domain and implicitly covers that
domain's wildcard — `example.com` validates both `example.com` and
`*.example.com`. The spec therefore rejects `*.` prefixes (they would
create a `\*.example.com` authorization covering `*.*.example.com`,
almost never the intent) and trailing dots (the API takes bare domains,
unlike Cloud DNS record names).

## FIXED_RECORD vs PER_PROJECT_RECORD

- `FIXED_RECORD` — the classic DNS-01 style record, unique per
  authorization. Default for global authorizations.
- `PER_PROJECT_RECORD` — one record per (domain, project), letting many
  Certificate Manager resources across projects share a single validation
  record. Default for non-global locations.

The spec leaves the field optional so GCP's location-appropriate default
applies; setting it is only needed for cross-project record sharing.

## Scope Boundaries

- `deletion_policy` is not modeled — absent from the released provider
  line the modules pin (`google ~> 6.0`).

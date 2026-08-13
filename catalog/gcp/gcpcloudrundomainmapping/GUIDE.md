# GcpCloudRunDomainMapping Guide

The judgment this guide protects: a domain mapping is a pointer, not a
deployment. It is free, immutable, and re-creates in seconds — the real
work lives outside it, in domain verification (one-time, human) and DNS
(the records it emits). Treat the mapping as the cheap middle of a
three-step story and the two ends as the parts that need planning.

## Verify the domain BEFORE the first apply

GCP refuses to create a mapping for a domain the deploying identity has
not verified — the error arrives at apply time, after everything else
worked. Verification is one-time per domain (Search Console or
`gcloud domains verify`), out-of-band by design: no Terraform or Pulumi
resource performs it. Two operational corollaries:

- **Verify as the identity that deploys.** A domain verified by your
  user account does not cover the service account a runner deploys with —
  add robot accounts as additional verified owners in Search Console.
- **Verify the apex once**: subdomains inherit verification, so
  `example.com` verified covers `app.example.com`, `api.example.com`,
  and every future subdomain.

## The mapping is immutable — plan for replacement, not update

Every field except `deletionPolicy` is create-only at the provider: a
changed domain, route, certificate mode, or even a label replaces the
mapping. Replacement is cheap (seconds, free) but not invisible — a
managed certificate re-issues after the replacement, so expect a brief
TLS gap on a live domain. Batch spec changes rather than trickling them,
and use `deletionPolicy: PREVENT` on production mappings so a casual
destroy cannot un-map a serving domain.

## Created ≠ serving: DNS closes the loop

A successful apply means the mapping EXISTS — the domain serves only
after the `resource_records` output is published in the domain's zone. A
root domain receives A/AAAA sets; a subdomain one CNAME
(`ghs.googlehosted.com.`). Wire the records into GcpDnsRecord (same
graph, zero copy-paste) or publish them at your external DNS host. The
mapping rests in `CertificatePending` until DNS propagates — a state the
provider itself treats as create-success, so IaC never hangs on
unpublished DNS. Certificate issuance follows minutes after DNS lands.

## Migrations: NONE first, AUTOMATIC after cutover

Mapping a domain that currently serves elsewhere? `certificateMode:
AUTOMATIC` would attempt issuance against DNS still pointing at the old
host. Create with `NONE`, stage the emitted records (lower TTLs ahead of
time), cut DNS over, then flip to `AUTOMATIC` — the flip is a
replacement, seconds. The migration preset carries this shape with
`deletionPolicy: PREVENT` armed for the cutover window.

## One domain, one mapping — force_override is a last resort

GCP fails a create for a domain already mapped elsewhere with a conflict
error; `forceOverride: true` silently steals the domain instead. The safe
workflow is to hit the conflict error first and set the flag only once
that error confirmed the takeover is intended — never leave it armed in a
standing manifest.

## When to graduate to the load balancer

The mapping serves the "one service, one domain" story. Multi-service
routing under one domain, Cloud Armor, CDN, cross-region failover, or
IPv6 static addressing all belong to the LB composition (serverless NEG →
backend service → URL map → HTTPS proxy → forwarding rule). The two
coexist fine during a migration — the mapping can keep serving while the
LB stack builds up, and DNS decides which front door traffic uses.

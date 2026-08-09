# GcpManagedSslCertificate Guide

The judgment this guide protects: a Google-managed classic certificate
finishes provisioning only AFTER the domain's DNS points at the load
balancer that serves it. Deploy success means the object exists, not
that TLS works — the sequencing around that gap is the whole operational
story.

## Provisioning is asynchronous and DNS-gated

Creation returns immediately with the certificate in PROVISIONING.
Google can issue it only once each domain's public DNS resolves to the
load balancer IP this certificate serves — until then, clients hitting
the domain get Google's default certificate and a browser warning. The
working order for a NEW domain:

1. Deploy the full serving chain (backend, proxy with this certificate,
   forwarding rule) — provisioning pends, which is fine.
2. Point the domain's A/AAAA records at the forwarding rule's IP.
3. Wait — issuance typically completes within the hour; `expire_time`
   stays empty until it does.

For a domain that is ALREADY serving traffic through another
certificate, this gap means downtime. That migration is exactly what
`GcpCertManagerCert` with a `GcpCertManagerDnsAuthorization` exists
for — validation completes before cutover. Prefer it for anything live.

## The domain list is the identity

`domains` is ForceNew: adding or removing a domain replaces the
certificate, and the replacement re-enters PROVISIONING per domain.
Rotate the same way as a self-managed certificate — create the new
certificate alongside, repoint the proxy, then destroy the old one
(GCP refuses to delete a certificate a proxy still references, so the
destroy fails rather than dropping TLS). `deletionPolicy: ABANDON`
hands the old certificate off mid-rotation; `PREVENT` guards one whose
replacement is not yet ACTIVE.

## What this kind cannot do

- **Wildcards** — classic managed certificates validate through the
  serving load balancer and cannot prove control of `*.example.com`;
  use `GcpCertManagerCert` with a DNS authorization.
- **Regional load balancers** — this resource is global-only.
- **Pre-issue validation** — there is no DNS-authorization arm here;
  issuance always waits on serving DNS.

## Conventions and gotchas

- Renewal is automatic and invisible — but it uses the same serving
  check, so a domain whose DNS moves away silently breaks the NEXT
  renewal, not the current certificate. Treat DNS moves as certificate
  events.
- Managed and self-managed certificates share one name namespace.
- Up to 100 domains per certificate; every one of them gates
  provisioning — one unpointed domain holds the whole certificate in
  PROVISIONING.

## Pairs well with

- `GcpTargetHttpsProxy` — consumes this certificate's self_link.
- `GcpGlobalForwardingRule` — the IP the domains must point at.
- `GcpCertManagerCert` + `GcpCertManagerDnsAuthorization` — the
  migration-safe, wildcard-capable successor stack.

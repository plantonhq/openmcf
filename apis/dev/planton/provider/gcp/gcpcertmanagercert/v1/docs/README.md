# GcpCertManagerCert: Design Notes

## What This Component Models

One `google_certificate_manager_certificate` — GCP's modern certificate
resource — with its two mutually exclusive arms modeled exactly as the API
defines them:

- **managed**: Google provisions and renews. The arm's inputs are the
  domain list and the domain-control proof.
- **self_managed**: uploaded PEM chain + private key. The key carries the
  secret-material annotation and is masked everywhere; the certificate
  chain is public by definition.

GCP has two certificate generations, and the catalog models both honestly
as separate kinds: the classic compute certificates
(`GcpManagedSslCertificate`, `GcpSslCertificate`) attach to proxies via
`sslCertificates`, while Certificate Manager certificates attach via
`certificateManagerCertificates` or a certificate map, add regional
serving, `scope` control, wildcard support, and pre-serving validation.
One kind per provider resource — a certificate kind that silently switches
between different underlying resources would leak its implementation as
behavioral surprises (different attachment fields, different validation
semantics).

## Why DNS Authorizations Are a Separate Kind

Certificate Manager's DNS authorization is its own API resource with its
own lifecycle: it can exist before any certificate, multiple certificates
can reference the same authorization, and it survives certificate
rotation. Modeling it inside the certificate would:

- hide a resource that other certificates need to share,
- force the certificate to also own DNS records (crossing into
  `GcpDnsRecord`'s ownership), and
- make the pre-serving validation flow — create authorization → serve its
  record → THEN create the certificate — inexpressible.

So `GcpCertManagerDnsAuthorization` is first-class, the certificate takes
`dns_authorizations` as references, and the validation record composes into
the zone as an explicit `GcpDnsRecord`.

## The Three Validation Modes

1. **DNS authorization** (references) — validates without serving traffic;
   required for wildcards; the production default.
2. **Issuance config** — private-PKI signing via a
   CertificateIssuanceConfig path. Mutually exclusive with DNS
   authorizations, mirroring the API.
3. **Load-balancer authorization** (both omitted) — GCP validates through
   the serving load balancer itself. Zero extra setup, but the certificate
   can only turn ACTIVE after traffic already reaches the proxy, and
   wildcards are not supported (enforced pre-deploy).

## Provisioning Semantics

A managed certificate is created immediately but stays in
`managed_state = PROVISIONING` until domain validation completes — which
requires the validation record to be publicly resolvable. The component
treats creation as success and exports the state honestly; waiting for
ACTIVE is a composition concern (attach after ACTIVE when downtime
matters), not a create-time gate.

## Scope Boundaries

- `deletion_policy` is not modeled — absent from the released provider
  line the modules pin (`google ~> 6.0`).
- `certificate_map` / `certificate_map_entry` (SNI-scale certificate
  serving) are separate API resources, recorded as future kinds; the
  proxy's `certificate_map` field accepts a literal map path today.
- The deprecated `certificate_pem`/`private_key_pem` provider aliases are
  skipped in favor of the current `pem_certificate`/`pem_private_key`.

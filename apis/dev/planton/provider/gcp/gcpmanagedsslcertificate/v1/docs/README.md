# GCP Managed SSL Certificate — Deep Dive

## Where this sits in the load balancing family

A Google Cloud external Application Load Balancer terminates TLS at a **target
HTTPS proxy**, which holds a list of SSL certificates. A Google-managed SSL
certificate is the hands-off option: Google generates the private key, submits
the CSR, validates domain ownership via DNS, issues the cert, and renews it
before expiry. You never see or store key material.

The certificate is its own first-class resource with a `self_link` — the
composition handle a target HTTPS proxy references in `ssl_certificates`. That
separation matters for rotation: you create a new certificate, repoint the
proxy, then destroy the old one.

## DNS-gated asynchronous provisioning

Creating the certificate object returns immediately, but the certificate stays
**PROVISIONING** until each domain's DNS points at the load balancer's IP (the
same forwarding rule the proxy serves). Until provisioning completes:

- The domain serves Google's default certificate, not yours.
- `expire_time` in stack outputs stays empty.
- Console status shows PROVISIONING rather than ACTIVE.

This is normal — the resource exists and is usable in composition (you can
attach it to a proxy) long before Google finishes issuance.

## Immutability and the create-before-destroy hazard

Every field of a Google-managed SSL certificate is immutable (ForceNew): changing
the name or the domain list destroys and recreates the certificate. Two
consequences matter operationally:

1. A certificate that a target HTTPS proxy references **cannot be deleted**
   while in use — GCP returns `resourceInUseByAnotherResource`. Recreating one
   therefore requires creating the replacement first and repointing the proxy's
   `ssl_certificates` list before the old certificate is destroyed
   (create-before-destroy). When you change the referencing proxy through
   Planton, the dependency ordering handles this.
2. Wildcard domains (`*.example.com`) are **not supported** by Google-managed
   certificates — list each hostname explicitly.

## Deliberate scope boundaries

- **Self-managed certificates** (`google_compute_ssl_certificate` where you
  supply the key) are a separate GCP resource and a distinct kind.
- **Certificate Manager** (`google_certificate_manager_certificate`) is the
  newer DNS-authorization-based flow for complex multi-cloud DNS; this kind
  covers the classic Compute Engine managed certificate attached directly to
  HTTPS proxies.
- No fields in this spec are sensitive — Google holds the private key.

## Composition

The certificate's `self_link` is what a target HTTPS proxy lists in
`ssl_certificates`. Pair it with a `GcpGlobalAddress` (for the DNS A record),
a `GcpUrlMap`, and `GcpBackendService` backends to assemble a full HTTPS load
balancer front end.

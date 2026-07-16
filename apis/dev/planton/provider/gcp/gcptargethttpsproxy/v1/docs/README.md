# GCP Target HTTPS Proxies: Where TLS Actually Terminates

## The Node That Owns the Handshake

In Google's decomposed global load balancer, exactly one resource decides how the client-facing TLS handshake behaves: the target HTTPS proxy. Certificates, minimum TLS version (via SSL policy), QUIC/HTTP-3 negotiation, TLS 1.3 0-RTT early data, and mutual-TLS client authentication all attach here — not on the forwarding rule (which only owns IP+port) and not on the backend service (which owns the LB→backend leg).

That separation has a practical payoff: certificates rotate, TLS policies tighten, and QUIC toggles — all without touching the VIP, DNS, or routing. GCP exposes dedicated in-place calls (`setSslCertificates`, `setSslPolicy`, `setUrlMap`, `setQuicOverride`) and the provider uses them, so these are true zero-downtime updates.

## Three Certificate Mechanisms, One Choice

The proxy accepts certificates through exactly one of three mutually exclusive mechanisms — the spec enforces the choice pre-deploy because GCP rejects combinations at the API:

1. **`sslCertificates`** — up to 15 classic compute SSL certificates. The workhorse: reference `GcpManagedSslCertificate` resources (Google handles issuance and renewal) or self-managed compute certificates by self-link. The load balancer picks the cert matching the client's SNI hostname.
2. **`certificateManagerCertificates`** — Certificate Manager certificates, honored ONLY by the cross-region internal ALB (`INTERNAL_MANAGED` scheme).
3. **`certificateMap`** — a Certificate Manager map that resolves the certificate by SNI at scale. The SaaS-custom-domains mechanism: thousands of customer domains behind one proxy, with domain onboarding happening entirely in Certificate Manager.

Traffic Director (`INTERNAL_SELF_MANAGED`) is the odd one out: it ignores `sslCertificates` entirely and drives TLS through `serverTlsPolicy`.

## The Chicken-and-Egg That Isn't

A Google-managed certificate stays PROVISIONING until its domains' DNS points at the load balancer's IP — but the load balancer only exists once the proxy (and forwarding rule) exist. GCP resolves this deliberately: **a PROVISIONING certificate attaches to a proxy just fine**, and attachment is in fact a prerequisite for activation. The correct bring-up order is: create the cert → attach it here → create the forwarding rule → point DNS at the VIP → the cert activates. Until then the domain serves Google's default certificate.

## server_tls_policy: The mTLS Lever

`serverTlsPolicy` references a network security ServerTlsPolicy that can demand and validate client certificates — mutual TLS. Two properties are worth knowing:

- It composes WITH server certificates on external ALBs (server cert from `sslCertificates`, client validation from the policy), and REPLACES them on Traffic Director.
- The provider deliberately PATCHes `null` when the field is cleared, so removing mTLS from a live frontend is a clean in-place downgrade to plain TLS.

## tls_early_data: Latency vs Replay Safety

TLS 1.3 0-RTT lets a resuming client send its first HTTP request inside the handshake — zero effective round trips, over TCP and QUIC. Early data is replayable by design, so the four modes are a security dial: `STRICT` (safe methods without query parameters only), `PERMISSIVE` (all requests), `UNRESTRICTED` (even non-idempotent replays — only for replay-tolerant services), `DISABLED` (the GCP default). Unlike everything else TLS on this proxy, the mode is immutable — changing it recreates the proxy.

## The 90/10 Coverage Decision

| Provider field | Modeled | Notes |
|---|---|---|
| `name` | ✅ `proxyName` | Defaults to `metadata.name`; RFC1035 validated |
| `project` | ✅ `projectId` | `StringValueOrRef` → GcpProject; empty → provider default |
| `description` | ✅ | |
| `url_map` | ✅ `urlMap` | Required ref → GcpUrlMap self_link; in-place |
| `ssl_certificates` | ✅ `sslCertificates` | Repeated ref → GcpManagedSslCertificate self_link; max 15 |
| `certificate_manager_certificates` | ✅ | Repeated ref → GcpCertManagerCert certificate_name |
| `certificate_map` | ✅ `certificateMap` | Plain URI (no Planton kind models cert maps yet) |
| `ssl_policy` | ✅ `sslPolicy` | Un-defaulted ref until an SSL-policy kind lands |
| `server_tls_policy` | ✅ `serverTlsPolicy` | Un-defaulted ref (network security policy) |
| `quic_override` | ✅ `quicOverride` | NONE default via middleware |
| `tls_early_data` | ✅ `tlsEarlyData` | Immutable; empty = GCP default |
| `http_keep_alive_timeout_sec` | ✅ | 5-1200 CEL; 0 = GCP default |
| `proxy_bind` | ✅ `proxyBind` | Traffic Director binding |
| `proxy_id` / `fingerprint` (computed) | outputs only | |
| `deletion_policy` | ❌ | Absent from the released 6.x line |
| `timeouts` | ❌ | Operation plumbing |

## Composition

1. **GcpManagedSslCertificate** — issues and renews the certificate; its `self_link` lands in `sslCertificates`.
2. **GcpUrlMap** — routes decrypted requests; its `self_link` is `urlMap`.
3. **GcpTargetHttpsProxy** (this component) — terminates TLS.
4. **GcpGlobalForwardingRule** — binds the VIP on port 443 to this proxy's `self_link`.

The kind registry declares `GcpUrlMap` and `GcpManagedSslCertificate` as prerequisites, so the E2E harness (and composed charts) installs both before the proxy.

## Operational Notes

- **Rotation runbook**: attach the replacement cert (list update, in-place) → wait ACTIVE → detach the old one. Never detach-then-attach; GCP requires ≥1 certificate on a proxy serving certificate-list traffic.
- **A certificate attached to a proxy cannot be deleted** — swap the proxy's reference first (create-before-destroy), or the destroy fails.
- The module enables `compute.googleapis.com` before creating the proxy, so a fresh project works on the first deploy.

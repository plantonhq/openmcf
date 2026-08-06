# GCP SSL Certificate (Self-Managed) — Deep Dive

## Where this sits in the load balancing family

A target HTTPS proxy presents certificates to clients from its
`ssl_certificates` list. GCP offers two ways to put a classic Compute Engine
certificate in that list: Google-managed (Google issues and renews;
`GcpManagedSslCertificate`) and **self-managed** (this kind — you upload a
PEM chain and private key). Both are the same API collection: they share one
name namespace per scope, carry the same kind of `self_link`, and attach to
proxies identically. The proxy cannot tell them apart.

Choose self-managed for what managed certificates cannot do:

- **Wildcards** — `*.example.com` is not issuable as a Google-managed cert.
- **Your own CA** — private/corporate CAs, EV/OV certificates, pinned chains.
- **Internal load balancers** — managed issuance validates via public DNS,
  which internal ALBs do not have; private-CA self-managed certs are the
  standard TLS story for internal HTTPS.
- **TLS before DNS cutover** — a managed cert stays PROVISIONING until DNS
  points at the load balancer; a self-managed cert is ACTIVE on creation.

The trade: **nothing renews itself.** Google never touches your expiry. The
`expire_time` stack output (parsed from the uploaded chain) is the rotation
clock.

## The private key is the only secret

The `private_key` field is annotated sensitive: encrypted in Pulumi state,
write-only in GCP (the API never returns it; the provider tracks it by
hash), and never surfaced in stack outputs. The `certificate` chain is
deliberately NOT treated as a secret — it is public handshake material
presented to every connecting client; masking it would only hurt review and
debugging. (The Terraform provider marks the chain "sensitive" purely to
keep multi-kilobyte PEM out of plan diffs — a display concern, not a secrecy
contract.)

Spec validation checks PEM framing (correct `-----BEGIN ...-----` headers,
including catching a key pasted into the certificate field and vice versa);
GCP validates the actual cryptographic material at deploy time — a
mismatched key/chain pair or an encrypted key is rejected there.

## Immutability and the rotation sequence

Every argument is ForceNew. A certificate a proxy references cannot be
deleted (`resourceInUseByAnotherResource`), which makes naive in-place
rotation impossible and is exactly the safety you want. The rotation
sequence:

1. Create the replacement certificate under a new name (version the
   `certificateName`, e.g. `prod-app-cert-2027`).
2. Update the proxy's `ssl_certificates` to the new reference — GCP swaps
   the list in place (setSslCertificates) with zero downtime.
3. Destroy the old certificate once nothing references it.

## One kind, two scopes

GCP models global and regional SSL certificates as separate API collections
with an identical surface, so this kind folds both: empty `region` creates
the global certificate, a set region creates the regional one. Scope is
permanent, and regional proxies only accept certificates in their own
region; the `region` stack output lets composition confirm compatibility.

## Deliberately not modeled (recorded reasons)

- **`private_key_wo` / `private_key_wo_version`** — Terraform's write-only
  argument flow; absent from the released 6.x schema, and Planton's
  sensitive-field handling already keeps the key encrypted and out of
  outputs.
- **`name_prefix`** — a Terraform-side create-before-destroy naming trick;
  Planton's metadata-driven naming owns resource names, and the rotation it
  serves is expressed as a new Planton resource with a versioned
  `certificateName`.
- **`deletion_policy`** — a client-side Terraform lever that conflicts with
  Planton-managed destroy (catalog-wide decision).

## Composition

The certificate's `self_link` goes into a target HTTPS proxy's
`ssl_certificates` list — the list's default reference kind is the managed
certificate, so self-managed entries use an explicit reference kind. Pair
with a `GcpSslPolicy` on the same proxy for a fully hardened TLS frontend.

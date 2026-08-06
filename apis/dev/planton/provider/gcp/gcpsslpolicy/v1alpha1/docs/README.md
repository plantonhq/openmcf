# GCP SSL Policy — Deep Dive

## Where this sits in the load balancing family

TLS termination at a Google Cloud Application Load Balancer has two halves:
**what** the load balancer presents (certificates, attached to the target
HTTPS proxy) and **how** it negotiates (protocol versions and cipher suites —
this resource). An SSL policy is referenced by the proxy's `ssl_policy`
field; one policy is commonly shared by every proxy in an organization,
making it the single lever for the fleet's TLS floor.

If a proxy has no SSL policy, GCP applies its default: minimum TLS 1.0 with
the COMPATIBLE cipher set. That default is deliberately permissive — any
compliance regime (PCI DSS, ISO 27001 audits, internal security review)
effectively requires attaching a policy.

## Profiles, and how CUSTOM pairs with custom_features

- **COMPATIBLE** — widest client range, including legacy ciphers.
- **MODERN** — drops broken ciphers, keeps broad reach. The right default
  for internet-facing production.
- **RESTRICTED** — only suites with modern security guarantees (forward
  secrecy, AEAD). Drops some older clients. Also the profile GCP requires
  when raising the floor beyond what other profiles allow.
- **CUSTOM** — an explicit allowlist in `custom_features`. GCP enforces the
  pairing both ways: CUSTOM without features is an error, and features on
  any other profile are an error. The spec validates both directions before
  deploy — the same rule the provider enforces at plan time.

Two facts auditors care about: GCP exposes no maximum-version control, and
TLS 1.3 suites are not listable in `custom_features` — TLS 1.3 is always
negotiable when the client supports it, whatever the profile says. The
`enabled_features` stack output reports the exact suites GCP computed from
the profile, which is the artifact a security review asks for.

## One kind, two scopes

GCP models global and regional SSL policies as separate API collections with
an identical configuration surface, so this kind folds both: empty `region`
creates the global policy (for global external HTTPS proxies), a set region
creates the regional one (for regional external and internal ALB proxies).
Scope is permanent — a policy cannot move between scopes, and regional
proxies only accept policies in their own region. The `region` stack output
lets downstream composition confirm scope compatibility.

## Mutability profile

`profile`, `min_tls_version`, and `custom_features` update **in place** and
take effect on the next client handshake of every referencing proxy — this
is what makes the policy a fleet-wide hardening lever. `name`, `project`,
and (unusually) `description` are ForceNew: changing them destroys and
recreates the policy, briefly breaking every proxy that references the old
`self_link`.

## Deliberately not modeled (recorded reasons)

- **`post_quantum_key_exchange`** (X25519MLKEM768 negotiation control),
  **`FIPS_202205` profile**, and **`TLS_1_3` as a minimum version** — all
  present only on the provider's unreleased main branch; the released 6.x
  schema (verified on both GA and beta) accepts none of them. Model when the
  released line ships them.
- **`deletion_policy`** — a client-side Terraform lever that conflicts with
  Planton-managed destroy (catalog-wide decision).

## Composition

The policy's `self_link` is what a target HTTPS proxy references in
`ssl_policy`. Pair it with certificates (`GcpManagedSslCertificate` or
self-managed `GcpSslCertificate`), a `GcpUrlMap`, and a
`GcpGlobalForwardingRule` to assemble a hardened HTTPS load balancer
frontend.

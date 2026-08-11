# GcpSslPolicy Guide

The judgment this guide protects: GCP's default TLS posture is
permissive on purpose (TLS 1.0, COMPATIBLE ciphers), and an SSL policy
is shared configuration with fleet-wide blast radius — one in-place edit
changes the handshake of every proxy referencing it, on the next client
hello.

## The floor ladder, and where to stand

- **No policy / COMPATIBLE + TLS 1.0** — maximum client reach, fails
  modern compliance scans. Acceptable only where ancient clients are a
  documented requirement.
- **MODERN + TLS_1_2** — the production default: drops broken ciphers,
  keeps broad reach, satisfies PCI DSS. Start here.
- **RESTRICTED** — only ciphers with modern security guarantees, and the
  only profile that accepts a `TLS_1_3` floor. Some older clients drop
  off; that is the point.
- **FIPS_202205** — the FIPS 140-2/3 validated suite set for FedRAMP-class
  regimes; requires exactly a `TLS_1_2` floor (both directions are
  validated before deploy).
- **CUSTOM** — hand-picked suites for regimes that name ciphers
  explicitly. A too-narrow list locks out real clients silently; prefer
  RESTRICTED unless an auditor hands you the list.

TLS 1.3 is always negotiable regardless of the floor — GCP has no
maximum-version control, and TLS 1.3 suites are never listable in
`customFeatures`. The floor governs what OLD protocols are refused.

## Post-quantum is a stance about dates, not a feature flag

`postQuantumKeyExchange` controls the X25519MLKEM768 hybrid group's
rollout: `DEFAULT` follows GCP's timeline (disallowed before October
2026, allowed after), `ENABLED` opts in now, `DEFERRED` opts out until
the later mandatory date (October 2027). Choose `ENABLED` when client
libraries are known-good; choose `DEFERRED` only with a recorded reason
to revisit — it is a countdown, not an off switch.

## Shared policy, shared blast radius

profile, minTlsVersion, customFeatures, and the post-quantum stance all
update IN PLACE and apply to every referencing proxy. That makes
tightening a fleet's floor a one-line change — and makes a typo in that
line a fleet-wide incident. Treat edits to a widely-referenced policy
with change-control gravity, and give compliance-mandated policies
`deletionPolicy: PREVENT` so a teardown cannot silently return the
fleet to GCP defaults.

## Conventions and gotchas

- Only name, project, and description are immutable — and yes, the
  description being ForceNew is a GCP quirk on this resource; write it
  once, well.
- A policy cannot move between global and regional scope, and regional
  proxies reference only policies in their own region.
- The `enabled_features` output is GCP's computed truth of what is
  actually served — audit against it, not against the spec.

## Pairs well with

- `GcpTargetHttpsProxy` — consumes this policy's self_link; the
  certificate is identity, this policy is negotiation.
- `GcpSslCertificate` / `GcpManagedSslCertificate` — served under the
  rules this policy sets.

# Kubernetes identity tier live-proven: Keycloak Operator, Keycloak, OpenBao and OpenFGA enter the green E2E matrix

## What changed

- **Four kinds proven against a live cluster, both engines** —
  KubernetesKeycloakOperator, KubernetesKeycloak, KubernetesOpenBao,
  KubernetesOpenFga. Every behavioral promise ran with verifier-output
  evidence: the Keycloak operator installed with all four CRDs
  established and NO server deployed (the two-kind grain's invariant),
  and its destroy removed workloads AND CRDs (the release-manifest
  bundle posture); a real admin logged into Keycloak with the
  bootstrap credentials and read the admin API on every lane, a
  verifier-owned realm survived a pod replacement because configuration
  lives in the composed PostgreSQL, and a two-instance deployment
  clustered live through the discovery Service; OpenBao's designed seal
  lifecycle ran end to end on every non-dev lane — fresh pods asserted
  NotReady-by-design and 501-uninitialized, initialized and unsealed by
  the verifier, readiness FLIPPING only then, with a KV v2 round-trip —
  and the Raft lane proved the Shamir restart truth (a replaced pod
  comes back SEALED, re-unseals with the run's key, and the marker
  survives on the Raft volume); OpenFGA answered the Zanzibar loop
  live (store → model → tuple → Check BOTH ways: the granted user
  ALLOWED, the ungranted DENIED) behind an asserted 401 auth gate.
  Blind import round-trips proved every kind's recipes, including the
  conditional inline-keys Secret and both multi-document bundle maps.
  All four entered the green E2E CI matrix.

- **Keycloak hostname pairings that only fail at boot are now
  validation errors.** The server refuses to start when
  `hostname-backchannel-dynamic` is enabled without the hostname
  declared as a FULL URL, and when a separate admin URL is set without
  a full-URL hostname — both surface on Kubernetes only as
  CrashLoopBackOff with the reason buried in the crashed container's
  log (verified live, source-verified at the pin). The spec now rejects
  both pairings at validation time with messages naming the server's
  own rules; the field comments, module comments and docs teach the
  full-URL requirement.

- **Pod-replacement proofs now gate on a NEW pod UID via a shared
  verifier helper.** A terminating pod keeps reporting Running/Ready
  for its whole grace period, so a wait keyed on phase or readiness
  straight after a pod delete can pass against the dying pod and run
  the "post-replacement" assertions against the old server. The shared
  helper waits for the same-named pod with a DIFFERENT uid —
  optionally at phase Running only, for servers whose readiness is
  deliberately gated on out-of-band bootstrap (a sealed vault) — and
  the Keycloak, OpenBao and SeaweedFS durability proofs all ride it.

## Why it matters

Identity is the front door of every platform: these four kinds cover
login/SSO, secrets management and fine-grained authorization. Each can
now be deployed through either engine with the confidence of live
proofs — including the operational truths a first-time operator must
know (a fresh OpenBao is sealed by design; a replaced Raft pod needs
re-unsealing; Keycloak's clustering needs nothing beyond the operator's
discovery Service) — and misconfigurations that previously produced
opaque crash-loops fail at validation with actionable messages.

# KubernetesCertificate: Research and Design

## Introduction

A Certificate is cert-manager's request contract: declare what certificate
should exist (names, lifetime, key parameters, outputs), point at who signs
it (issuerRef), and cert-manager keeps it issued and renewed into a TLS
Secret for as long as the resource exists. This component covers the
complete cert-manager.io/v1 Certificate surface.

## Design Authority

Designed field-by-field from the pinned upstream CRD schema
(`cert-manager.io_certificates.yaml`). Every constraint the CRD or the
controller enforces that a user could trip is mirrored in validation with a
teaching message: literal-subject exclusivity, renew-before vs
renew-before-percentage, per-algorithm key sizes, the x509 usages
vocabulary, signature algorithms, keystore password requirements, PKCS#12
profiles, output-format vocabulary.

## Coverage Decisions

- **All SAN types**: DNS (wildcards), IP, URI (the SPIFFE mTLS pattern),
  email, and otherName (OID + UTF-8 value — Microsoft UPN client certs)
- **`literal_subject`** carries the full LDAP RFC 4514 DN when attribute
  ORDER matters (LDAP auth certificates); mutually exclusive with
  `subject`/`common_name` because it embeds them
- **`renew_before_percentage`** is modeled alongside `renew_before`:
  the percentage form scales when the issuer overrides the requested
  duration (Let's Encrypt always issues 90 days regardless of request)
- **Keystores**: JKS and PKCS#12 with INLINE sensitive passwords (upstream
  v1.15+ supports literal passwords natively — no Secret pre-creation);
  PKCS#12 profiles modeled with the legacy default documented
- **`is_ca` + `name_constraints`**: the delegated-internal-CA guardrail —
  constraints on what a bootstrapped CA may sign
- **Issuer selection**: FK-backed ClusterIssuer/Issuer arms (chart
  composition) plus an `external` arm (group/kind/name) for third-party
  issuer controllers (AWS PCA, Google CAS) — full upstream issuerRef
  fidelity
- **Deliberately unmodeled**: feature-gated surfaces at the pinned upstream
  until graduation

## Engine Mechanics

- **Pulumi**: typed crd2pulumi Certificate resource (types regenerated from
  the pinned CRDs); enum vocabularies translate from the proto's lowercase
  to the CRD's exact casing in one builder
- **Terraform**: `kubectl_manifest` renders the identical CR (the enum
  translation tables are mirrored in locals) — plannable before the CRDs
  exist, which single-run infra charts require
- **Neither engine waits for issuance**: issuance time belongs to the
  issuer; an unreachable CA would block forever. The E2E verifier polls
  Ready and checks the signed material landed in the Secret — the module
  does not.

## E2E

Real issuance proven on the kind cluster: minimal and full-surface
scenarios issue from a live self-signed ClusterIssuer fixture (resolved
through the FK), and the verifier requires Ready PLUS tls.crt/tls.key data
in the target Secret. The root-CA fixture this kind ships powers the
KubernetesIssuer CA-chain composition proof. ACME issuance (public CA, real
DNS) rides the batched real-cluster lanes.

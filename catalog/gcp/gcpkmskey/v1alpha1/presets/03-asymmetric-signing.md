# Asymmetric Signing Key

The code- and artifact-signing primitive: an `ASYMMETRIC_SIGN` key whose
private material never leaves Cloud KMS, with public keys fetched by
verifiers per version.

## What this preset creates

An EC P-256 signing key. Signing happens inside KMS via the
`asymmetricSign` API; consumers download the per-version public key to
verify. There is no rotation period — verifiers pin exact key versions,
so new versions are minted deliberately, not on a schedule.

## Prerequisites

- A `GcpKmsKeyRing` named `signing-keys` (a `global` ring is common for
  signing keys consumed from many regions — see the `GcpKmsKeyRing`
  presets).

## Remix ideas

- Switch `algorithm` to `RSA_SIGN_PSS_4096_SHA512` for RSA-ecosystem
  verifiers.
- Pair with Binary Authorization attestors: the attestor references this
  key's versions to verify container signatures at admission.
- Set `versionTemplate.protectionLevel: HSM` for hardware-held signing
  material.

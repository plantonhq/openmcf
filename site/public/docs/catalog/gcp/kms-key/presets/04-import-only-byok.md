---
title: "Import-Only BYOK Key"
description: "The bring-your-own-key container: a key that can only ever hold imported material, for organizations whose policy requires key generation outside Google's infrastructure."
type: "preset"
rank: "04"
presetSlug: "04-import-only-byok"
componentSlug: "kms-key"
componentTitle: "KMS Key"
provider: "gcp"
icon: "package"
order: 4
---

# Import-Only BYOK Key

The bring-your-own-key container: a key that can only ever hold imported
material, for organizations whose policy requires key generation outside
Google's infrastructure.

## What this preset creates

An empty, import-only symmetric key. GCP will never generate a version
for it — `importOnly` is the immutable guarantee, and
`skipInitialVersionCreation` (required with it) prevents the initial
auto-generated version. The key has no usable versions until material is
imported through a key ring import job, so CMEK consumers cannot use it
until the import ceremony completes.

## Prerequisites

- A `GcpKmsKeyRing` named `byok-keys` (see the `GcpKmsKeyRing` presets).
- An externally generated key, wrapped and imported via
  `gcloud kms keys versions import` (the import-job ceremony lives
  outside infrastructure-as-code by design — it handles raw key
  material).

## Remix ideas

- Set an `EXTERNAL` protection level instead when the material should
  stay in an external key manager entirely (Cloud EKM) rather than being
  imported into Google's HSM/software fleet.
- Add `destroyScheduledDuration` to extend the recovery window —
  re-importing destroyed BYOK material requires the original source.

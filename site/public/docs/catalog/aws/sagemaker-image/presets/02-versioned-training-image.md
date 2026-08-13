---
title: "Versioned Training Image"
description: "This preset registers a GPU training environment as a fully-annotated version — ECR image, compatibility metadata, and the `latest` and `stable` aliases — so training jobs select it by name instead..."
type: "preset"
rank: "02"
presetSlug: "02-versioned-training-image"
componentSlug: "sagemaker-image"
componentTitle: "SageMaker Image"
provider: "aws"
icon: "package"
order: 2
---

# Versioned Training Image

This preset registers a GPU training environment as a fully-annotated
version — ECR image, compatibility metadata, and the `latest` and
`stable` aliases — so training jobs select it by name instead of by
raw ECR path.

## When to Use

- A team-standard training container that jobs should reference by
  alias, not by ECR digest
- The first versioned entry in a registry that will grow with each
  image build

## What You Get

- One AWS-numbered version pointing at the ECR image, marked
  `TRAINING` / `GPU` / `PyTorch 2.4` / `python 3.12` with a `STABLE`
  vendor signal
- `latest` and `stable` aliases that move freely between versions —
  promote a new build by moving them, in place

## Customize

- Append the next build as a new entry at the END of `versions` —
  entries are keyed by position and a changed `baseImage` replaces the
  version under a new AWS-assigned number
- Add `releaseNotes` (max 255 characters) to carry the build's
  changelog into the registry
- Flip `vendorGuidance` to `TO_BE_ARCHIVED` when a version is on its
  way out — it updates in place

<p align="center">
  <img src="logo.svg" alt="AWS SageMaker Image" width="80"/>
</p>

# AWS SageMaker Image

Create and manage [Amazon SageMaker AI images](https://docs.aws.amazon.com/sagemaker/latest/dg/studio-byoi.html)
— the named registry entries that make YOUR container images (custom
kernels, training environments) selectable in Studio and notebook
surfaces, together with their folded versions, each pointing at a
concrete ECR image in the same account and region.

## What Gets Created

- **An image** — the named registry entry, carrying the IAM role
  SageMaker assumes to pull from ECR plus Studio display metadata
  (`display_name`, `description`) that updates in place.
- **Versions** (folded satellites) — each registers one ECR image
  (`base_image`) as an AWS-numbered version, with movable `aliases`
  (e.g. `latest`, `stable`) and in-place compatibility metadata:
  `job_type`, `ml_framework`, `processor`, `programming_lang`,
  `release_notes`, `vendor_guidance`.

## Versions Are Append-Only

Version numbers are AWS-assigned, monotonic, and never reused: changing
a version's `base_image` replaces it under a NEW number (the old number
stays retired). Entries in `versions` are keyed by position — append
new versions at the end and never reorder existing entries. The
provider sleeps ~1 minute before create for IAM propagation, and image
versions carry no tags upstream.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).

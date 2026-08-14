# AwsSagemakerImage — Terraform/OpenTofu module

Deploys an Amazon SageMaker AI image (`aws_sagemaker_image`) with its
folded version satellites (`aws_sagemaker_image_version`).

Module facts worth knowing before editing:

- **Versions are keyed by position** — the `for_each` keys are the
  entries' indices ("0", "1", ...), the append-only contract taught on
  the spec; reordering entries rewires satellites.
- **`base_image` is the version's identity** — changing it replaces
  the version under a new AWS-assigned number (numbers are monotonic
  and never reused; the old number stays retired).
- **The provider sleeps ~1 minute before `CreateImage`** (IAM
  propagation) — every create is at least a minute.
- **AWS serializes version creation per image** (the provider holds a
  mutex) — versions attach one at a time.
- **`aws_sagemaker_image_version` carries no tags** (by provider
  design) — only the image is tagged.

Outputs mirror the Pulumi module key-for-key: `image_name`,
`image_arn`, `version_numbers`.

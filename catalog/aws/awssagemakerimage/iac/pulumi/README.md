# AwsSagemakerImage — Pulumi module (Go)

Deploys an Amazon SageMaker AI image (`sagemaker.Image`) with its
folded version satellites (`sagemaker.ImageVersion`).

Module facts worth knowing before editing:

- **Versions are keyed by position** — resource names use the entries'
  indices ("version-0", "version-1", ...), identical keys to the
  Terraform module's `for_each` and the append-only contract taught on
  the spec; reordering entries rewires satellites.
- **`BaseImage` is the version's identity** — changing it replaces the
  version under a new AWS-assigned number (numbers are monotonic and
  never reused; the old number stays retired).
- **The provider sleeps ~1 minute before `CreateImage`** (IAM
  propagation) — every create is at least a minute.
- **AWS serializes version creation per image** (the provider holds a
  mutex) — versions attach one at a time; the module orders them after
  the image with `DependsOn`.
- **`sagemaker.ImageVersion` carries no tags** (by provider design) —
  only the image is tagged.

Outputs mirror the Terraform module key-for-key: `image_name`,
`image_arn`, `version_numbers`.

# AwsLambdaLayer — Pulumi module

Publishes one Lambda layer version (`lambda.LayerVersion`) from its S3 archive, with share grants (`lambda.LayerVersionPermission`), one per spec entry.

Module facts worth knowing before editing:

- **Everything is ForceNew** — any change publishes a NEW version; `SkipDestroy` keeps the old one available in AWS.
- **The grant action is pinned** — `lambda:GetLayerVersion` is the only action AWS supports on layers; the module never exposes it.
- **The version bridges string-to-int** — the layer reports `Version` as a string, the permission takes a number; the module converts via ApplyT.
- **Neither resource is taggable at AWS** — hence no tags anywhere in this module.

Outputs mirror the Terraform module key-for-key: `layer_arn`, `layer_version_arn` (import ID), `version`, `code_sha256`, `permission_revision_ids`.

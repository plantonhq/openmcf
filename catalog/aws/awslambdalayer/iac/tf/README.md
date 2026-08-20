# AwsLambdaLayer — Terraform/OpenTofu module

Publishes one Lambda layer version (`aws_lambda_layer_version`) from its S3 archive, with share grants (`aws_lambda_layer_version_permission`, keyed by statement_id).

Module facts worth knowing before editing:

- **Everything is ForceNew** — any change publishes a NEW version; `skip_destroy` keeps the old one available in AWS.
- **The grant action is pinned** — `lambda:GetLayerVersion` is the only action AWS supports on layers; the module never exposes it.
- **Grants reference the folded layer** — `layer_name` takes the layer ARN and `version_number` the published version (a string at the layer, a number at the grant).
- **Neither resource is taggable at AWS** — hence no tags anywhere in this module.

Outputs mirror the Pulumi module key-for-key: `layer_arn`, `layer_version_arn` (import ID), `version`, `code_sha256`, `permission_revision_ids`.

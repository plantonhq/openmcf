# AwsBedrockInvocationLogging — Terraform/OpenTofu module

Manages the region's Bedrock model invocation logging configuration
(`aws_bedrock_model_invocation_logging_configuration`) — a settings
singleton delivering model call logs to CloudWatch Logs and/or S3.

Module facts worth knowing before editing:

- **One configuration per account+region** (the resource id IS the
  region). Two instances against the same region fight over it.
- **Data-type toggles are presence-typed.** Null passes through and
  the provider applies AWS's enabled-by-default; explicit false is
  sent. Do not collapse them to plain bools — that would make
  "unset" and "false" indistinguishable.
- **The spec's CEL guarantees at least one destination**, so the two
  dynamic blocks never render an empty logging_config.
- **Bedrock writes to S3 as its service principal** — the bucket
  policy, not an IAM role, is what authorizes the S3 arm. The
  CloudWatch arm DOES use the role.
- **Destroy deletes the configuration** (unlike some siblings in the
  settings-singleton class that only reset).
- **No tags.** The upstream resource carries no tags argument.

Outputs mirror the Pulumi module key-for-key: `configured_region`.

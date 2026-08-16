# AwsSsmParameter — Terraform/OpenTofu module

Manages one Parameter Store entry (`aws_ssm_parameter`).

Module facts worth knowing before editing:

- **The name is `spec.parameter_name`, never `metadata.name`** —
  hierarchical names carry slashes metadata.name cannot.
- **The value arms cross-map**: `spec.secure_value` renders as the
  provider's sensitive `value`; `spec.value` renders as
  `insecure_value` (readable plans — that argument's purpose). The
  spec guarantees exactly one arm.
- **`overwrite` renders only when true** — an explicit false would
  break the provider's own update path (unset means fail-on-foreign
  at create, overwrite own updates).
- **Optional strings render null-when-empty** so the module never
  fights provider defaults (tier Standard, the aws/ssm key,
  data_type text).

Outputs mirror the Pulumi module key-for-key: `parameter_name`,
`parameter_arn`, `version` (stringified), `tier` (the AWS-resolved
tier).

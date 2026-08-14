# AwsSsmParameter — Pulumi module (Go)

Manages one Parameter Store entry (`ssm.Parameter`).

Module facts worth knowing before editing:

- **The name is `spec.parameter_name`, never `metadata.name`** —
  hierarchical names carry slashes metadata.name cannot.
- **The value arms cross-map**: `spec.secure_value` renders as the
  provider's sensitive `Value`; `spec.value` renders as
  `InsecureValue` (readable previews — that argument's purpose). The
  spec guarantees exactly one arm.
- **`Overwrite` renders only when true** — an explicit false would
  break the provider's own update path (unset means fail-on-foreign
  at create, overwrite own updates).
- **Optional strings render only when set** so the module never
  fights provider defaults (tier Standard, the aws/ssm key,
  data_type text).

Outputs mirror the Terraform module key-for-key: `parameter_name`,
`parameter_arn`, `version` (stringified), `tier` (the AWS-resolved
tier).

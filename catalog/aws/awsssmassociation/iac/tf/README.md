# AwsSsmAssociation — Terraform/OpenTofu module

Manages one State Manager association (`aws_ssm_association`).

Module facts worth knowing before editing:

- **`name` is the DOCUMENT reference** (`spec.document_name` — a
  literal AWS-managed name or a resolved AwsSsmDocument reference) and
  forces replacement; everything else versions in place.
- **AWS identifies the association by a generated UUID** — the
  `association_id` output, never the name.
- **Optional strings render null-when-empty** so the module never
  fights provider defaults (document version `$DEFAULT`, sync
  compliance AUTO).
- **`wait_for_success_timeout_seconds` renders only when set** — a
  create-time gate that would always time out on target-less
  associations.

Outputs mirror the Pulumi module key-for-key: `association_id`,
`association_arn`, `document_name`.

# AwsSsmAssociation — Pulumi module (Go)

Manages one State Manager association (`ssm.Association`).

Module facts worth knowing before editing:

- **`Name` is the DOCUMENT reference** (`spec.document_name` — a
  literal AWS-managed name or a resolved AwsSsmDocument reference) and
  forces replacement; everything else versions in place.
- **AWS identifies the association by a generated UUID** — the
  `association_id` output, never the name.
- **Optional values render only when set** so the module never fights
  provider defaults (document version `$DEFAULT`, sync compliance
  AUTO).
- **`WaitForSuccessTimeoutSeconds` renders only when set** — a
  create-time gate that would always time out on target-less
  associations.

Outputs mirror the Terraform module key-for-key: `association_id`,
`association_arn`, `document_name`.

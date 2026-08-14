# AwsSsmDocument — Pulumi module (Go)

Manages one customer-owned SSM document (`ssm.Document`).

Module facts worth knowing before editing:

- **The document name is `metadata.name`** (materialized as
  `Locals.DocumentName`) and forces replacement.
- **Sharing renders the provider's flat Permissions map** from the
  spec's typed account-ID list: `{type: "Share", account_ids:
  "<comma-joined>"}` — `Share` is the only legal type.
- **Content updates create a new document version** and the provider
  promotes it to the default; schema-1.x documents only update when
  content itself changes (an AWS rule the module cannot soften).
- **Attachment metadata never round-trips** — no SSM API returns it;
  the import map declares `attachments_source` config-only.

Outputs mirror the Terraform module key-for-key: `document_name`,
`document_arn`, `default_version`, `latest_version`, `document_hash`,
`status`.

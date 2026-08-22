# AwsOrganizationalUnit — Pulumi module

Manages one organizational unit
(`organizations.OrganizationalUnit`).

Module facts worth knowing before editing:

- **`Name` renders `spec.ou_name`** — the explicit name field (OU
  names allow spaces `metadata.name` cannot carry); renames apply in
  place.
- **`ParentId` arrives resolved** — the platform resolves the
  value-or-reference before the module runs; the module reads
  `spec.ParentId.GetValue()`.
- **The parent forces replacement** — AWS moves accounts, never OUs.
- **Creation rides the provider's finalization retry** — up to four
  minutes of FinalizingOrganizationException tolerance right after
  CreateOrganization, so org-then-OU chains compose cleanly.

Outputs mirror the Terraform module key-for-key: `ou_id` (the import
ID) and `arn`.

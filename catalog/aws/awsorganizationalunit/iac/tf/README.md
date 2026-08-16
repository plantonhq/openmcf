# AwsOrganizationalUnit — Terraform/OpenTofu module

Manages one organizational unit
(`aws_organizations_organizational_unit`).

Module facts worth knowing before editing:

- **`name` renders `spec.ou_name`** — the explicit name field (OU
  names allow spaces `metadata.name` cannot carry); renames apply in
  place.
- **`parent_id` arrives resolved** — the platform resolves the
  value-or-reference (an AwsOrganization's `root_id`, a parent OU's
  `ou_id`, or a literal) before the module runs; the module passes the
  string through.
- **The parent is ForceNew** — a parent change replaces the OU (AWS
  moves accounts, never OUs).
- **Creation rides the provider's finalization retry** — up to four
  minutes of FinalizingOrganizationException tolerance right after
  CreateOrganization, so org-then-OU chains compose cleanly.

Outputs mirror the Pulumi module key-for-key: `ou_id` (the import ID)
and `arn`.

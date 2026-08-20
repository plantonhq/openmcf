# AwsOrganization — Terraform/OpenTofu module

Manages THE organization (`aws_organizations_organization`) with its
folded delegated-administrator registrations
(`aws_organizations_delegated_administrator`, for_each keyed
`{account_id}//{service_principal}`) and the optional singleton
resource policy (`aws_organizations_resource_policy`, count-gated).

Module facts worth knowing before editing:

- **`feature_set` renders only on an explicit choice** (unset = ALL,
  the provider default). The ALL → CONSOLIDATED_BILLING downgrade
  forces replacement — delete-and-recreate of the ENTIRE organization.
- **Service access and policy types are diffed enable/disable calls**
  on update (the provider applies disables first) — never a resource
  replacement.
- **The delegated-admin for_each key uses "//"** — the import
  machinery's segment delimiter; each half of the provider's
  `{account_id}/{service_principal}` import composite derives from its
  own key segment.
- **The resource policy is a per-organization singleton** — the
  count-gated arm IS that object; the org resource itself is
  untaggable, so the identity tag map lands here.
- **Destroy calls DeleteOrganization** — AWS requires zero members,
  OUs, and policies first.

Outputs mirror the Pulumi module key-for-key: `organization_id`,
`arn`, `management_account_id`, `management_account_arn`,
`management_account_email`, `root_id` (the OU/policy wiring feed), and
`resource_policy_id` (the folded singleton's import ID, empty when the
arm is absent).

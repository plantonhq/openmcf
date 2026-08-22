# AwsOrganization — Pulumi module

Manages THE organization (`organizations.Organization`) with its
folded delegated-administrator registrations
(`organizations.DelegatedAdministrator`, one per spec entry, parented
to the organization) and the optional singleton resource policy
(`organizations.ResourcePolicy`).

Module facts worth knowing before editing:

- **`FeatureSet` is sent only on an explicit choice** (unset = ALL,
  the provider default). The ALL → CONSOLIDATED_BILLING downgrade
  forces replacement — delete-and-recreate of the ENTIRE organization.
- **Service access and policy types are diffed enable/disable calls**
  on update — never a resource replacement.
- **The resource policy is a per-organization singleton** — the
  conditional arm IS that object; the org resource itself is
  untaggable, so the identity tag map lands here.
- **`root_id` exports `Roots[0].Id`** — the provider reads the root
  list with the organization; first-level OUs and root-scoped policy
  attachments wire to it.
- **Destroy calls DeleteOrganization** — AWS requires zero members,
  OUs, and policies first.

Outputs mirror the Terraform module key-for-key: `organization_id`,
`arn`, `management_account_id`, `management_account_arn`,
`management_account_email`, `root_id`, and `resource_policy_id`
(empty when the arm is absent).

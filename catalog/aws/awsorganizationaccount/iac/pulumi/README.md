# AwsOrganizationAccount — Pulumi module

Manages a member account (`organizations.Account`) with its folded
alternate contacts (`account.AlternativeContact` — NOTE the bridge
renames the provider's "alternate" to "Alternative"), primary contact
(`account.PrimaryContact`), and opt-in region enablements
(`account.Region`), all parented to the account.

Module facts worth knowing before editing:

- **The account carries `pulumi.IgnoreChanges(["roleName"])`** — AWS
  has no read API for the bootstrap role; without the ignore,
  importing an existing account plans a destructive replacement (the
  Terraform module mirrors this with `lifecycle.ignore_changes`).
- **`CloseOnDeletion` and `CreateGovcloud` are always sent
  explicitly** (both engines pin the engine-behavior booleans).
- **Contact satellites wire `AccountId` from the pivot's ID** — never
  configurable.
- **Primary-contact optional leaves are sent only when set** — the
  Put API has no unset semantics; empties mean "leave the last value
  on file".
- **Region enablement sends `Enabled` as configured** — `false`
  actively disables; destroying the entry is a provider NO-OP.

Outputs mirror the Terraform module key-for-key: `account_id`, `arn`,
`state`, and `govcloud_id`.

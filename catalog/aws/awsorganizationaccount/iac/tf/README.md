# AwsOrganizationAccount — Terraform/OpenTofu module

Manages a member account (`aws_organizations_account`) with its folded
alternate contacts (`aws_account_alternate_contact`, for_each keyed by
type token), primary contact (`aws_account_primary_contact`,
count-gated), and opt-in region enablements (`aws_account_region`,
for_each keyed by region name).

Module facts worth knowing before editing:

- **`role_name` carries `lifecycle.ignore_changes`** — AWS has no read
  API for it; without the ignore, importing an existing account plans
  a destructive replacement (the Pulumi module mirrors this with
  `IgnoreChanges`).
- **`close_on_deletion` and `create_govcloud` are always sent
  explicitly** (both engines pin the engine-behavior booleans) — which
  destroy AWS performs is modeled configuration here.
- **Contact satellites wire `account_id` from the pivot** — never
  configurable; the same APIs' manage-my-own-account form is
  deliberately out of scope.
- **Primary-contact optional leaves render only when set** — the Put
  API has no unset semantics, so an empty leaf means "leave the
  last value on file", and the module never sends empties.
- **Region enablement sends `enabled` as configured** — an entry with
  `enabled = false` actively disables; destroy of the entry is a
  provider NO-OP (the region keeps its last state).

Outputs mirror the Pulumi module key-for-key: `account_id` (the
import ID; contact and region composites build on it), `arn`, `state`,
and `govcloud_id`.

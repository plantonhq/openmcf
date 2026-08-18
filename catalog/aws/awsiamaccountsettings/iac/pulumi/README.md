# AwsIamAccountSettings — Pulumi module (Go)

Manages IAM's account-level settings: the sign-in alias
(`iam.AccountAlias`), the password policy
(`iam.AccountPasswordPolicy`), and the STS global-endpoint token
version (`iam.SecurityTokenServicePreferences`).

Module facts worth knowing before editing:

- **This is a settings singleton for a GLOBAL service** — one object
  per ACCOUNT (not per region); `metadata.name` never reaches AWS.
- **Each arm renders only when present** — an omitted arm leaves the
  account's current setting untouched, and that omission is
  meaningful.
- **Destroy semantics DIFFER per arm**: the alias truly deletes, the
  password policy resets to AWS defaults, the STS preference's delete
  is a NO-OP (the last-applied version persists).
- **The password policy is replaced WHOLE on every apply** — an unset
  field is AWS's default, never "keep the current setting". Plain
  bools render unconditionally (false == absent at AWS); the
  presence-typed knobs (minimum length, self-service changes, age,
  reuse) render only when set.
- **Nothing here is taggable at AWS** — the module deliberately
  carries no tag map (mirrored in the Terraform module).

Outputs mirror the Terraform module key-for-key: `account_id`,
`account_alias`, `expire_passwords` (both empty when their arm is
unset).

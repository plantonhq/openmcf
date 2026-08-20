# AwsIamAccountSettings — Terraform/OpenTofu module

Manages IAM's account-level settings: the sign-in alias
(`aws_iam_account_alias`), the password policy
(`aws_iam_account_password_policy`), and the STS global-endpoint token
version (`aws_iam_security_token_service_preferences`).

Module facts worth knowing before editing:

- **This is a settings singleton for a GLOBAL service** — one object
  per ACCOUNT (not per region); `metadata.name` never reaches AWS.
- **Each arm renders only when present** (`count` on arm presence) —
  an omitted arm leaves the account's current setting untouched, and
  that omission is meaningful.
- **Destroy semantics DIFFER per arm**: the alias truly deletes, the
  password policy resets to AWS defaults, the STS preference's delete
  is a NO-OP (the last-applied version persists).
- **The password policy is replaced WHOLE on every apply** — an unset
  field is AWS's default, never "keep the current setting". Plain
  bools render unconditionally (false == absent at AWS); the
  presence-typed knobs (minimum length, self-service changes, age,
  reuse) pass through as null when unset.
- **Nothing here is taggable at AWS** — the module deliberately
  carries no tag map (mirrored in the Pulumi module).

Outputs mirror the Pulumi module key-for-key: `account_id`,
`account_alias`, `expire_passwords` (both empty when their arm is
unset).

<p align="center">
  <img src="logo.svg" alt="AWS IAM Account Settings" width="80"/>
</p>

# AWS IAM Account Settings

Manage [IAM's account-level settings](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_passwords_account-policy.html)
— the sign-in alias, the password policy, and the STS global-endpoint
token version.

This is a **settings singleton** for a GLOBAL service: IAM keeps
exactly one of each setting per ACCOUNT (not per region), so deploy at
most one instance per account. `metadata.name` never reaches AWS.

## What Gets Managed

- **accountAlias** — the console sign-in URL
  (`https://<alias>.signin.aws.amazon.com/console`). An account has
  exactly ONE alias, globally unique across all of AWS: applying this
  arm REPLACES whatever alias existed, changing the sign-in URL
  everyone uses.
- **passwordPolicy** — IAM user console-password rules (length,
  character classes, expiry, reuse prevention, self-service changes).
  Replaced WHOLE on every apply: an unset field is AWS's default,
  never "keep the current setting".
- **sts** — which token version the global STS endpoint issues
  (v2Token works in ALL regions including opt-in ones).

Destroy semantics DIFFER per arm: the alias truly deletes, the
password policy resets to AWS defaults, and the STS preference
persists (a no-op delete).

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).

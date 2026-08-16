# Sign-In Alias

This preset sets the account's console sign-in alias — the friendly
URL (`https://<alias>.signin.aws.amazon.com/console`) humans use
instead of the 12-digit account ID.

## When to Use

- Every human-facing account, once — a readable sign-in URL per
  account naming convention
- Multi-account setups where "which account is this" must be obvious
  at the sign-in page

## What You Get

- The alias applied account-wide (3-63 characters: lowercase letters,
  digits, single hyphens)
- The alias echoed as an output for docs and onboarding material

## Customize

- **An account has exactly ONE alias and applying this REPLACES it**
  — check `aws iam list-account-aliases` first; every bookmarked
  sign-in URL changes with it
- Aliases are globally unique across ALL of AWS — namespace them
  (`<company>-<env>`)
- Destroying this arm deletes the alias; sign-in reverts to the bare
  account ID

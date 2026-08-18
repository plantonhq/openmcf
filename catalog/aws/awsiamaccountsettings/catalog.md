# AWS IAM Account Settings

The account's IAM-wide switches in one place: the console sign-in
alias, the password policy for IAM users, and the STS global-endpoint
token version. One instance per account — IAM is global and keeps
exactly one of each setting.

## What Gets Managed

- The sign-in alias (the friendly console URL — an account has
  exactly one, and applying it replaces whatever existed).
- The password policy: length, character classes, expiry, reuse
  prevention, self-service changes — replaced whole on every apply.
- The STS token version: v2 tokens work in every region, including
  opt-in ones.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with IAM permissions.

### AWS Account

- Nothing — the settings objects always exist; this component manages
  their values. Check `aws iam list-account-aliases` before applying
  the alias arm: applying replaces the current alias and changes the
  sign-in URL everyone uses.

## Deploy

### Console

Create the resource from the AWS catalog, set the arms to manage, and
deploy.

### CLI

```bash
planton apply -f iam-account-settings.yaml
```

## After Deploy

- IAM user sign-ins are held to the policy immediately (existing
  passwords age against maxPasswordAge from their last change).
- Destroy contracts differ per arm: alias deletes, password policy
  resets to AWS defaults, the STS version persists.

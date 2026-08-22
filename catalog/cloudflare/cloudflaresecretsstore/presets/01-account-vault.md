# Account vault

The account's single Secrets Store under a boring, permanent name. Cloudflare allows exactly one store per account and both fields are create-only -- this is permanent infrastructure, created once and never renamed.

**When to use it:** the first time any Worker binding or AI Gateway in the account needs a shared secret.

**What to change:** only the `account_id`. Keep the name unremarkable and stable -- renaming replaces the store and destroys every secret inside it. If your account already has a store, adopt it by import instead of applying this preset.

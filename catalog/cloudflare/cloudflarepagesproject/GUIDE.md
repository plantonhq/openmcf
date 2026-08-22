# CloudflarePagesProject guide

Operational judgment for Pages projects. The README covers what each field is; this covers how the pieces interact.

## The name is the identity

Pages has no separate id. The project name is unique in the account, CEL-constrained to lowercase alnum/hyphens (max 58), and is what every API path and import ID uses. Renaming is a new project.

## Preview and production are a pair

Cloudflare rejects a project whose environments disagree on knobs like `fail_open`. If you only set one environment, the module mirrors it to both so a single config just works. Set both explicitly only when they must differ.

## Secrets do not survive import

`deployment_configs.*.secrets` are `secret_text` env vars. The API never returns the value, so a post-import plan always wants to re-assert them. Provide them as managed-secret references so the platform resolves them just-in-time; never put a literal secret in the manifest.

## Direct-upload vs git source

A project with only `name` + `productionBranch` is a direct-upload project (`wrangler pages deploy` later). A git source (`spec.source`) connects a GitHub/GitLab repo and needs the Cloudflare Git integration already authorized — that is a dashboard/OAuth step this kind cannot perform. Bindings (KV, D1, R2, queues, …) are per-environment and reference already-deployed resources.

## Custom domains initialize slowly

Attaching a domain creates a `cloudflare_pages_domain` that sits `initializing` / `pending` until DNS and SSL catch up. On a throwaway hostname that wait never ends. Use a real hostname in a zone on this account, or skip domains on the live path and plan-prove the send path only.

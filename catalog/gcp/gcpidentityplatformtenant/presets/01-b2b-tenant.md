# B2B Tenant

One isolated user pool for one customer organization — the standard
B2B SaaS shape. Users sign up with email/password inside the tenant and
never mix with any other customer's pool.

## What it configures

- `displayName` — the customer organization's name, the only naming
  input (the tenant's resource ID is server-generated; read the
  `tenant_id` output).
- `allowPasswordSignup: true` — email/password sign-up enabled for the
  tenant.
- `deletionPolicy: PREVENT` — destroy FAILS rather than deleting the
  tenant, because deletion removes every user account unrecoverably.

## Adjust before deploying

- **displayName** — one tenant per customer organization; name it after
  the customer.
- **The project** — its GcpIdentityPlatformConfig must already set
  `multiTenant.allowTenants: true`, or tenant creation is rejected.
- Wire the client SDK's `tenantId` from the `tenant_id` output — an SDK
  without it authenticates against the project-level pool instead.

## When to choose something else

If the customer brings their own corporate IdP, start from the
**SSO Tenant** preset instead — their users authenticate against their
own directory and no passwords live in the tenant.

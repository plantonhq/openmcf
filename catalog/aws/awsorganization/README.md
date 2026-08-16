<p align="center">
  <img src="logo.svg" alt="AWS Organization" width="80"/>
</p>

# AWS Organization

Manage [THE AWS Organization](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_introduction.html)
of the deploying account — creating it makes that account the
organization's management account — together with trusted service
access, delegated administrators, and the org's resource-based
delegation policy.

## What Gets Managed

- **The organization** (`o-...` — the import ID): its **featureSet**
  (ALL when unset — the level every advanced arm requires; the
  downgrade to CONSOLIDATED_BILLING replaces the ENTIRE organization).
- **Trusted service access** folded as
  `awsServiceAccessPrincipals` — the ONE home for it (the provider
  warns the standalone resource fights this argument with a perpetual
  diff).
- **Policy types** folded as `enabledPolicyTypes` — a type must be
  enabled here before any [AWS Organization Policy](../awsorganizationpolicy)
  of that type can attach.
- **Delegated administrators** folded as immutable
  `{accountId, servicePrincipal}` registrations (each imports as
  `{account_id}/{service_principal}`).
- **The resource policy** folded as a structured document — AWS keeps
  exactly ONE per organization (`rp-...`).

Destroy deletes the whole organization — AWS requires every member
account, OU, and policy gone first, by design.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).

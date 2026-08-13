# GcpIdentityPlatformTenant Guide

The judgment this guide protects: a tenant IS a customer's user base.
Every account, credential, and IdP connection for that organization lives
inside it, and the destructive operations are absolute.

## One tenant per customer organization

The multi-tenancy model that works for B2B SaaS: each customer
organization gets its own tenant — separated users, separated identity
providers, separated policies — inside one shared project. A user in one
tenant does not exist in another; a customer's SAML connection touches
only their pool. Resist the temptation to share a tenant between
customers "temporarily": user pools cannot be split later.

## The tenant ID is not yours to choose

`displayName` is the only naming input. GCP generates the resource ID at
create time, and everything downstream — the client SDK's `tenantId`,
tenant-scoped API calls — must consume the `tenant_id` output rather
than a predicted value. Wire it through references, never literals.

## disableAuth is the kill switch

Setting it true stops ALL authentication in the tenant and existing
sessions stop refreshing — effectively instant lockout of one customer
without touching their data or configuration. The right tool for
offboarding disputes, incident containment, or suspending a delinquent
account; flipping it back restores service.

## PREVENT is the production posture

`deletionPolicy: DELETE` removes the tenant with every user account in
it, unrecoverable — there is no undelete for a user pool. Any tenant
whose users are real people should carry `PREVENT`, making a destroy
fail loudly instead of erasing a customer. Reserve `DELETE` for
ephemeral and test tenants; `ABANDON` parks the tenant outside
management while it keeps serving sign-ins.

## Client apps must scope to the tenant

An SDK initialized without `tenantId` authenticates against the
project-level pool, not the tenant — a silent wrong-pool bug. Every
client surface serving a tenant's users sets `tenantId` from this
resource's `tenant_id` output; the sign-in flows, providers, and users
it sees are then exactly the tenant's own.

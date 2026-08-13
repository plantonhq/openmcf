# Azure Mongo Cluster User -- Operational Guide

Judgment calls that matter when you run Entra grants on Mongo vCore clusters in production.

## Grants are how apps should connect; the admin is break-glass

The native administrator credential on the cluster is one shared password with total power. The grant model inverts that: each application's managed identity gets its own binding, authenticates with short-lived Entra tokens, and can be revoked alone. Wire every application through a grant and keep the admin credential in a vault for incidents -- if an app ships with the admin connection string, rotating that password becomes an every-consumer outage.

## The object ID is the identity; get the RIGHT one

Managed identities carry two UUIDs -- the client ID and the PRINCIPAL (object) ID. The grant wants the principal ID (reference the identity's principal output rather than copying UUIDs by hand). For app registrations, use the ENTERPRISE APPLICATION's object ID (the service principal in your tenant), not the app registration's object ID -- the wrong one creates a grant nobody can authenticate against, which looks like a driver bug, not a config bug.

## Replacement is the update model, and that is fine

Every field is create-only: role changes, principal changes, anything -- the provider drops and re-adds the grant. That is safe for an access binding (no data rides on it), but it does mean a moment without access during the replace. For zero-gap role changes, add a second grant first, then remove the old one.

## Entra auth is a cluster-level switch someone must have thrown

A grant against a cluster whose `authenticationMethods` lacks "MicrosoftEntraID" fails at deploy time. When adopting Entra auth on an existing NativeAuth-only cluster, update the cluster first (an in-place update), then land the grants. Keep "NativeAuth" in the list during migration -- removing it cuts over every consumer at once.

# GcpAlloydbUser Guide

The judgment this guide protects: database users are credentials with a
lifecycle, not rows in a setup script. One user per application, IAM
where possible, and a destroy path that understands PostgreSQL ownership.

## IAM users first, passwords second

`ALLOYDB_IAM_USER` authenticates through IAM — no stored password, no
rotation calendar, revocation through the same grants as everything else.
It is the right default for services running as GCP service accounts
(the userId is the account's email). `ALLOYDB_BUILT_IN` with a password
is for tools and drivers that cannot speak IAM; the spec refuses a
password on an IAM user because a stored credential on an IAM identity
is a contradiction.

## Rotation is an in-place update

For built-in users, changing `password` rotates it live — no recreate,
no connection reset for sessions already open. Rotate on a schedule by
editing the spec; the sensitive-field handling keeps the value out of
rendered plans.

## Grant roles, not superuser

`databaseRoles` is additive PostgreSQL role membership. Most application
users need only their schema's roles; `alloydbsuperuser` is for the one
migration/admin identity, not the app fleet. IAM users typically carry
`alloydbiamuser` plus their application roles.

## Dropping a user is not dropping its objects

PostgreSQL keeps ownership rows when a role disappears. Before a
`DELETE` destroy of a user that ever created tables, reassign ownership
inside the database (REASSIGN OWNED BY) — otherwise the next migration
trips over orphaned ownership. `PREVENT` suits the credential a
production service still authenticates with; `ABANDON` leaves the user
live on the cluster while dropping it from management — the handoff
path, not the cleanup path.

## Immutable identity

`cluster`, `userId`, and `userType` recreate the user when changed.
Renaming a user is a new user plus an application cutover — treat the
userId as an API contract with every connection string that embeds it.

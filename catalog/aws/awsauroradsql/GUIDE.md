# AwsAuroraDsql — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## IAM tokens are the only door

DSQL has no CREATE USER WITH PASSWORD — every connection authenticates with a 15-minute IAM auth token as the password (`aws dsql generate-db-connect-admin-auth-token`, or the non-admin variant for roles you map inside the database). Applications need token-refresh plumbing, not credential storage; there is no secret to rotate.

## Design multi-region at day zero

The pairing window is PENDING_SETUP — a freshly created cluster before its first activation completes. The modules order the peering correctly, but the PEERS must already exist, so a multi-region database is deployed as one instance of this kind per region, each naming the others' `cluster_arn` and the same witness region. A live single-region cluster cannot be upgraded — plan the topology before data lands.

## The witness region is a third region, deliberately

The witness stores transaction logs and arbitrates during a region failure — it runs no queries and holds no full data. Pick a third region distinct from both peers; changing it later replaces the cluster.

## Not all of PostgreSQL is here

DSQL speaks the PostgreSQL wire protocol but is a distributed engine underneath: no extensions, no foreign keys enforcement quirks aside, optimistic concurrency instead of long locks — transactions that conflict RETRY rather than block. Port schema and load tests before committing a migration; "PostgreSQL-compatible" is a dialect claim, not a drop-in guarantee.

## Deletes are gated twice by design

`deletion_protection_enabled` refuses deletes at AWS; `force_destroy` makes the module disable it first (a wall becomes a speed bump - keep it false in production). The DELETED record lingers briefly after destroy; the E2E verifier treats DELETING/DELETED as absent.

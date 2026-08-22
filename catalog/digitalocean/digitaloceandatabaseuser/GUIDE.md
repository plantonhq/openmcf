# DigitalOcean Database User -- Operational Guide

What experience with this component teaches that the field reference cannot.

## One user per service, always

The cluster's built-in `doadmin` user works everywhere, which is exactly why production should not use it: one leaked credential exposes everything, and rotation breaks every consumer at once. Create one user per service; each gets its own server-generated password that rotates (by replacing the user) and revokes (by deleting it) independently.

## Passwords rotate by replacement

DigitalOcean generates passwords; there is no "set password" surface here. To rotate a credential, create a replacement user (new name), move the service to it, then delete the old user. Deleting a user revokes its access immediately.

## ACLs: declare everything, trust the manifest

DigitalOcean returns Kafka/OpenSearch ACLs only at create time -- the console and API reads never show them again. Treat the manifest as the single source of truth and review permission changes in code review, not in the console. ACL edits apply in place (no replacement).

## MySQL auth plugin: leave it unset

Unset means DigitalOcean's `caching_sha2_password` -- the modern plugin. Set `mysql_native_password` only for clients too old to speak it, and treat that as a dated compatibility decision to revisit. Clearing the field later resets the user to the modern default in place.

## MongoDB quirk

MongoDB clusters return the user's password ONLY in the create response. The outputs still carry it (captured at create), but an imported MongoDB user has no recoverable password -- rotate by replacement instead.

## Serialized operations on busy clusters

User creates and deletes serialize per cluster (an API constraint the provider enforces with a lock). A chart creating ten users on one cluster deploys them one at a time -- expect linear time, not a hang.

## What is deliberately NOT here

In-database GRANTs for PostgreSQL/MySQL (data-plane SQL, not IaC), password values in specs (server-generated only), and per-user settings MongoDB exposes in the raw API but the provider does not bridge.

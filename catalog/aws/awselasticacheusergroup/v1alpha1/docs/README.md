# AWS ElastiCache User Group: The RBAC Attachment Unit

## What This Component Is

An ElastiCache user group collects RBAC users and attaches, as one object,
to replication groups (`user_group_ids`) and serverless caches
(`user_group_id`). Users define WHO may connect and WHAT they may do; the
group defines WHERE those identities apply. `AwsElasticacheUserGroup` models
one group — caches reference the group's id, not individual users, so
granting an application access is a membership edit here, never a change to
the cache resource.

Create-time immutable in AWS: the user group id (`metadata.name`) and
`engine`. Membership (`userIds`) updates in place.

## Why Groups Exist Separately from Users

AWS's RBAC model has three layers, and conflating them hides the graph:

1. **Users** — identities with ACL access strings and authentication modes
2. **Groups** — membership sets that attach to caches as one object
3. **Caches** — reference group ids, not user ids directly

Modeling the group as its own node makes the architecture graph honest: you
see exactly which users apply to which caches, and you revoke access by
removing a user id from the group — the cache's spec never changes.

## The Mandatory Default User

AWS refuses to create a group unless its membership includes a user whose
user NAME is exactly "default". That user defines what unauthenticated
clients may do. The standard production pattern pairs a locked-down default
user (`accessString: "off ~* +@all"`, `no-password-required`) with
per-application password or IAM users. CEL cannot prove name-data across
resources, so this constraint surfaces at deploy time on the group — the
spec comment and presets steer users to include the default user, but
validation cannot reject a missing one upfront.

## Membership as References, Not Glue Resources

The Terraform provider exposes `aws_elasticache_user_group_association` as
a separate resource for adding users to a group. That glue has no
independent lifecycle — membership belongs on the group. The spec's
`userIds` list is the single surface; both engines update membership in
place on the group resource.

## Deliberately Skipped (with reasons)

- **Outpost-only fields**: ElastiCache on Outposts is a separate deployment
  surface; deferred until real demand appears.
- **`aws_elasticache_user_group_association`**: glue with no independent
  lifecycle — folded into the spec's `userIds` repeated reference field,
  updated in place on the group.

## Operational Notes

- A group only accepts users of its own engine (`redis` or `valkey`) and
  region.
- Memcached has no RBAC; user groups apply to Redis and Valkey only.
- Adding a user to a group takes effect without recreating the group or
  the cache — existing connections keep their session until they reconnect.
- A cache can reference multiple groups (replication groups) or one group
  (serverless caches) depending on the cache resource type.

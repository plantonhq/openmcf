# AWS ElastiCache User: Per-Application Cache Identity

## What This Component Is

An ElastiCache user is one identity in AWS's Redis/Valkey RBAC system: a
user id, an AUTH user name, an ACL access string, and an authentication
mode. RBAC replaces the old shared AUTH token model — each application gets
its own credentials and permissions, and revoking one application never
disturbs the others. `AwsElasticacheUser` models one user; user groups
(collecting users) and caches (referencing groups) are separate nodes so
access control composes through the graph.

Create-time immutable in AWS: the user id (`metadata.name`), `userName`,
and `engine`. The access string and authentication mode update in place.

## The RBAC Graph: User → Group → Cache

- **User (this component)** — WHO and WHAT: identity plus ACL access string
- **User group** — WHERE membership is declared: which users apply together
- **Cache** — references the group's id in `user_group_ids` / `user_group_id`

Adding an application to a cache is a membership edit on the group — the
cache resource never changes. Credential rotation is an in-place update on
the user (add a second password, roll clients, remove the old).

## Authentication Modes

- **Password** — the client sends `AUTH <userName> <password>`. One or two
  passwords (16–128 characters) enable zero-downtime rotation.
- **IAM** — the client signs a short-lived token with its AWS IAM identity.
  Requires `userName` to equal the user id, transit encryption on the cache,
  and `elasticache:Connect` on both the user ARN and the cache ARN.
- **No-password-required** — for the mandatory "default" user when switched
  "off" in its access string; never for a production user that is "on".

## The Mandatory "Default" User

Every user group must contain a user whose user NAME is exactly "default".
It defines what unauthenticated clients may do. The production pattern is
`accessString: "off ~* +@all"` with `no-password-required` — anonymous
connections are rejected outright, and every client must authenticate as a
named user. CEL cannot prove name-data across resources, so this constraint
surfaces at deploy time on the group, not at user validation time.

## Deliberately Skipped (with reasons)

- **Outpost-only fields** (`user_id` overrides, Outpost ARNs): ElastiCache
  on Outposts is a separate deployment surface; deferred until real demand
  appears.
- **Legacy flat auth arms** on the provider resource (`passwords` /
  `no_password_required` as top-level fields): folded into the nested
  `authenticationMode` message — one honest shape with a type discriminator.
- **Terraform write-only `passwords_wo`**: the spec carries passwords inline
  (marked sensitive); write-only fields are a Terraform-provider convenience
  with no protobuf equivalent.

## Operational Notes

- Users are regional — a group only accepts users from its own region.
- Memcached has no RBAC; these users apply to Redis and Valkey only.
- Tightening an access string takes effect on new connections without
  recreating the user.
- Several users may share a `userName`; AWS unions their credentials at
  authentication time.

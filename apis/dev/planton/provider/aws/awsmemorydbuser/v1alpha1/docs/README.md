# AwsMemorydbUser — Design Notes

## Provider mapping

Maps 1:1 onto `aws_memorydb_user` (Terraform) / `memorydb.User` (Pulumi).
The full provider surface is modeled: access string, authentication mode
(type + passwords), and the name.

## Design decisions

- **The user name derives from `metadata.name`.** In MemoryDB the user name
  IS the user's single AWS identity — Required + ForceNew in the provider,
  unique per region, and the value ACLs reference. There is no
  ElastiCache-style split between a user id and a reusable AUTH name, so a
  separate spec field would duplicate `metadata.name` with no added meaning.
  AWS caps the name at 40 characters and enforces it at create.
- **Only `password` and `iam` authentication types.** The provider's input
  enum (`InputAuthenticationType`) accepts exactly these two. The broader
  read-side enum also carries `no-password`, but that state is not creatable
  through the API — modeling it would invent surface AWS rejects.
  Unauthenticated access exists only through the built-in "open-access" ACL.
- **Passwords are sensitive.** 1–2 entries, 16–128 characters each
  (provider-validated bounds), marked `(sensitive)` so they never appear in
  rendered manifests; AWS never returns them on read.

## Deliberately skipped (with reasons)

- `name_prefix` — the naming basis is `metadata.name`; random-suffix naming
  would diverge the engines (no AWS proto uses it).
- Tags beyond the identity set — the AWS surface has no per-kind user tags
  field by convention; identity tags derive from metadata.

## Update semantics

`access_string` and `authentication_mode` update in place via UpdateUser
(the provider waits for the user to return to active). The user name forces
replacement.

<p align="center">
  <img src="logo.svg" alt="AWS IAM Group" width="80"/>
</p>

# AWS IAM Group

Manage an [IAM group](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_groups.html)
— the container that grants a set of users a shared permission set —
with its declarative membership and its policies.

## What Gets Managed

- **The group** (named from `metadata.name`; renames update in place —
  the ARN recomputes but members and policies persist) and its IAM
  **path**.
- **Membership** — the DECLARATIVE `users` list (references to
  [AWS IAM User](../awsiamuser) outputs or literal names): this list
  is authoritative, so out-of-band additions are removed on the next
  apply. The users must already exist.
- **Permissions**: **managedPolicyArns** (references to
  [AWS IAM Policy](../awsiampolicy) outputs or literal ARNs —
  AWS-managed policies included), each attachment reconciling
  individually; and **inlinePolicies**, a name-keyed map of JSON
  documents that live and die with the group.

IAM groups are untaggable at AWS — the one deliberate absence against
the catalog's tag convention.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).

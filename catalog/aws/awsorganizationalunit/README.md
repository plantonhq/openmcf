<p align="center">
  <img src="logo.svg" alt="AWS Organizational Unit" width="80"/>
</p>

# AWS Organizational Unit

Manage an [AWS Organizations organizational unit](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_manage_ous.html)
— a container in the organization's OU tree that member accounts are
grouped into and organization policies attach to.

## What Gets Managed

- **The OU** (`ou-...` — the import ID): its **ouName** (an explicit
  spec field — OU names allow spaces and arbitrary characters
  `metadata.name` cannot carry; renames apply in place) and its
  **parentId** — the organization root for a first-level OU (the
  default wiring resolves an [AWS Organization](../awsorganization)
  reference's `root_id`), a parent OU for nesting, or a literal
  `r-...`/`ou-...` ID.

The parent is IMMUTABLE — AWS moves accounts between OUs, never OUs
themselves; a parent change replaces the OU.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).

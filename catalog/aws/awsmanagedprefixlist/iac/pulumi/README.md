# AwsManagedPrefixList — Pulumi module (Go)

Manages one customer-managed prefix list (`ec2.ManagedPrefixList`) with its entries in-line.

Module facts worth knowing before editing:

- **`AddressFamily` replaces the list** — every referencing security-group rule, NACL rule, and route breaks with the old `pl-` id; everything else updates in place.
- **In-line entries are the single declarative owner** — the standalone `ec2.ManagedPrefixListEntry` resource is the identical payload and fights this form; this module never uses it.
- **Resizes are safe** — the provider orders `MaxEntries` increases before entry changes and decreases after, so a resize never transiently strands entries.
- **Description-only edits cost two API round trips** — the provider removes and re-adds the entry; expected, not drift.
- **`MaxEntries` is the quota contract** — a security-group rule referencing this list consumes `MaxEntries` rule-quota slots regardless of actual entry count.

Outputs mirror the Terraform module key-for-key: `prefix_list_id` (import ID), `prefix_list_arn`, `owner_id`, `version`.

# AwsManagedPrefixList — Terraform/OpenTofu module

Manages one customer-managed prefix list (`aws_ec2_managed_prefix_list`) with its entries in-line.

Module facts worth knowing before editing:

- **`address_family` replaces the list** — every referencing security-group rule, NACL rule, and route breaks with the old `pl-` id; everything else updates in place.
- **In-line entries are the single declarative owner** — the standalone `aws_ec2_managed_prefix_list_entry` resource is the identical payload and fights this form; this module never uses it.
- **Resizes are safe** — the provider orders `max_entries` increases before entry changes and decreases after, so a resize never transiently strands entries.
- **Description-only edits cost two API round trips** — the provider removes and re-adds the entry; expected, not drift.
- **`max_entries` is the quota contract** — a security-group rule referencing this list consumes `max_entries` rule-quota slots regardless of actual entry count.

Outputs mirror the Pulumi module key-for-key: `prefix_list_id` (import ID), `prefix_list_arn`, `owner_id`, `version`.

# AwsManagedPrefixList

One customer-managed prefix list: a named, versioned set of CIDR blocks that security-group rules, NACL rules, and route tables reference as a single `pl-` id. Update the list once and every referencing rule follows. Entries fold in-line as the single declarative owner.

## Highlights

- **The catalog's reusable CIDR vocabulary**: office networks, partner ranges, scanner fleets — defined once, referenced everywhere, updated in one place.
- **`max_entries` is the honest capacity contract**: a security-group rule referencing the list consumes `max_entries` rule-quota slots regardless of actual entry count — taught on the field, sized deliberately in every preset.
- **Provider quirks stay in the modules**: resize ordering (increases before entries, decreases after) and the two-round-trip description edit never leak to authors; the family-match and unique-CIDR CELs catch mistakes at validation.
- **`address_family` replaces the list** — and every reference to the old `pl-` id with it; taught on the field.

## Both Engines

Both modules render the single resource with in-line entries identically and export the same outputs: `prefix_list_id` (import ID), `prefix_list_arn`, `owner_id`, `version` (AWS's per-change concurrency token).

## Chart Wiring

The `prefix_list_id` output is what AwsSecurityGroup rules, AwsNetworkAcl consumers, and route-table surfaces reference. Standalone by design — the list itself references nothing.

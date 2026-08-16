# Team Chargeback

This preset creates the one-rule category most organizations need
first: every team's spend grouped by the `team` tag's own values, with
untagged spend surfaced as "unallocated".

## When to Use

- Chargeback/showback by team, product, or any tag your resources
  already carry
- Finding untagged spend — the "unallocated" bucket IS the tagging
  gap report

## What You Get

- One INHERITED_VALUE rule fanning out to a value per distinct tag
  value — new teams appear automatically, no rule edits
- A default bucket that makes untagged spend visible instead of
  hiding it in "Uncategorized"

## Customize

- The tag key must be an ACTIVATED cost-allocation tag in Billing or
  the rule inherits nothing
- Pin specific values first with REGULAR rules above the inherited
  one (first match wins) — e.g. route a shared account wholesale
  before the tag fan-out
- Split the "unallocated" bucket across teams with a
  `splitChargeRules` entry (see the shared-platform-split preset)

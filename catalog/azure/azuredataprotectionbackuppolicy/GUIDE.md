# Azure Data Protection Backup Policy -- Operational Guide

Judgment calls that matter when you run Data Protection policies in production.

## The policy is immutable -- plan replacements, not edits

Every field on every variant is fixed at creation: the provider ships no update path at all. Changing a schedule or retention value REPLACES the policy, and the backup instances bound to it re-bind to the replacement. Name policies by their contract (e.g. `daily-p7d`) rather than by their consumer, so a changed contract reads as the new object it actually is.

## Match the variant to the datasource family

Each variant maps to exactly one Azure datasource type, and the stores differ by design: disks and AKS retain on the OPERATIONAL store (snapshots near the source), databases and Data Lake on the VAULT store (isolated copies), and blob storage is the only dual-tier variant (continuous in-account retention AND scheduled vaulted copies). The spec vocabularies only accept what each datasource's service actually supports today -- a rejected store value is the service's boundary, not a catalog gap.

## Read the schedule grammar once, carefully

`R/2024-01-01T02:00:00+00:00/P1D` -- from that instant, repeat daily. The date anchors the phase (backups run at 02:00 UTC because the anchor says so); the `P` duration sets the cadence (`P1D` daily, `P1W` weekly, `PT4H` every four hours). The time zone field shifts interpretation for two variants Azure validates strictly (MySQL, Data Lake) -- the Windows-style names ("India Standard Time"), not IANA names ("Asia/Kolkata").

## Retention rules are tags, not filters

A named rule (criteria + duration + priority) TAGS matching backups at creation time -- the first backup of the week tagged `weekly` lives out that rule's duration. When several rules match one backup, the LOWEST priority number wins; the unnamed default layer (priority 99 internally) catches everything untagged. Data Lake is the exception to keep in mind: its rules have no priority field -- ORDER in the list is priority.

## Blob's two tiers are different products

The operational tier (`operationalDefaultRetentionDuration` alone) is continuous point-in-time restore inside the storage account -- no schedule, no vault copy, and no named rules. The vault tier (`vaultDefaultRetentionDuration` + intervals) is scheduled copies into the vault that survive account deletion. Configure both for defense in depth; configure only operational when in-account restore is all you need.

## Keep the install-manifest policy boring

When a policy exists to serve backup-instance testing or chart scaffolding, use the disk variant with a plain daily schedule -- the simplest shape with the fewest service-side moving parts. Exercise exotic grammar (calendar criteria, dual tiers, month-end selectors) in policies dedicated to it.

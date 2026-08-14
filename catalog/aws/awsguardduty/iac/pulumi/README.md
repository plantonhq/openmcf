# AwsGuardDuty — Pulumi module

Manages the region's GuardDuty posture:
`aws:guardduty/detector:Detector` plus its satellites — detector
features, filters, IP/threat-intel sets, the publishing destination,
the organization surface (admin account, configuration, org
features), members with per-member features, and the member-side
invite accepter.

Module facts worth knowing before editing:

- **Feature resources are patches.** Create/Update are the same
  UpdateDetector call and Delete is a no-op upstream — features
  removed from the spec are NOT reverted by AWS, and upstream
  serializes feature writes per detector under a global mutex (every
  feature resource parents on the detector).
- **The org-configuration delete is a no-op too** — destroy leaves
  the org posture as last applied (taught in the GUIDE).
- **Resource names are the spec entry names** (features and filters
  by name, sets by name, members by account id, member features by
  "account-feature") — the same keys the Terraform module's for_each
  uses and the output maps carry. Iteration is name-sorted for
  deterministic previews.
- **The publishing destination is deliberately untagged** — tags are
  ForceNew on it upstream (a tag edit would REPLACE findings export).

Outputs mirror the Terraform module key-for-key: `detector_id`,
`detector_arn`, `account_id`, `ip_set_ids`, `threat_intel_set_ids`,
`publishing_destination_id`.

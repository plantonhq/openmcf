# Windows Default Baseline

This preset claims the account/region DEFAULT baseline for Windows:
security and critical updates auto-approve after a 3-day soak, and
available-but-unapproved security updates report as NON_COMPLIANT so
nothing hides.

## When to Use

- Accounts where every Windows node — grouped or not — should patch
  against your policy instead of AWS's predefined default
- Compliance postures where "a security update exists and is not yet
  approved" must be visible

## What You Get

- A Windows rule on SecurityUpdates/CriticalUpdates at MSRC
  Critical/Important with `approveAfterDays: 3`
- The default designation: ungrouped Windows nodes patch against this
  baseline; deleting the component RESTORES AWS's own default
- The Windows-only available-security-updates posture set to
  NON_COMPLIANT

## Customize

- Add `patchGroups` to also govern named fleets explicitly
- Add `approvedPatchesEnableNonSecurity`-style breadth only on Linux —
  Windows patches are always security-classified
- At most one baseline per OS should set `setAsDefaultBaseline` —
  claiming displaces the previous holder silently

# Region Opt-Ins

This preset sets the region's resource-type opt-ins to the common
posture: the core data stores on, the niche types deliberately off —
every type listed explicitly.

## When to Use

- Before the region's first backup plan — opt-ins gate what plans can
  protect, silently
- Making the opt-in posture reviewable code instead of console state

## What You Get

- Explicit true/false for every common resource type (AWS returns the
  full set on read — listing everything is what keeps plans clean of
  perpetual differences)
- One instance per region; the region IS the identity

## Customize

- Flip types as your estate grows — a type left `false` is silently
  skipped by every backup plan selection in the region
- Add `resourceTypeManagementPreference` to let AWS Backup fully
  manage a type's advanced features (one-way at AWS: flippable per
  type, never clearable)
- Destroy reverts NOTHING (a no-op at AWS) — to opt a type out, apply
  it as `false`

# AWS SSM Patch Baseline

The patching policy for one operating system: which patches are
auto-approved (and when), which are explicitly blocked, which patch
groups it governs, and whether it is the account's default for that
OS.

## What Gets Managed

- Approval rules: filter patches by classification/severity/product,
  then approve days after release (a soak period) or everything up to
  a fixed date.
- Explicit approve/reject lists, with rejected-dependency behavior
  (install anyway, or block and report).
- Patch groups: bind the baseline to fleets carrying the matching
  `Patch Group` tag.
- The default designation for the OS — with delete restoring AWS's own
  default.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with SSM permissions.

### AWS Account

- Nothing — baselines stand alone. Patching later needs managed nodes
  tagged into the governed patch groups.

## Deploy

### Console

Create the resource from the AWS catalog, pick the OS, add an approval
rule, name the patch groups, and deploy.

### CLI

```bash
planton apply -f patch-baseline.yaml
```

## After Deploy

- Nodes tagged into a governed patch group evaluate against this
  baseline when `AWS-RunPatchBaseline` runs — schedule it with
  [AWS SSM Association](/cloud-catalog/aws-ssm-association) (scan) and
  [AWS SSM Maintenance Window](/cloud-catalog/aws-ssm-maintenance-window)
  (install).
- Baselines, groups, and the designation are free.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).

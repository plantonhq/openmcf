<p align="center">
  <img src="logo.svg" alt="AWS SSM Patch Baseline" width="80"/>
</p>

# AWS SSM Patch Baseline

Manage an [AWS Systems Manager patch baseline](https://docs.aws.amazon.com/systems-manager/latest/userguide/patch-manager-predefined-and-custom-patch-baselines.html)
— which patches get auto-approved for one operating system, the patch
groups the baseline governs, and the account/region default
designation.

## What Gets Managed

- **The baseline** (`metadata.name` is the baseline name; AWS
  identifies it as `pb-...`): its **operatingSystem** (15 supported,
  WINDOWS when unset), **approval rules** (filters plus exactly one
  approval arm — days-after-release or a fixed date, mutually
  exclusive by CEL), **global filters**, explicit
  **approved/rejected patch lists** with the rejected-dependency
  action, alternative **repo sources** (Linux), and Windows'
  available-security-updates compliance posture.
- **Patch groups** folded as name registrations — each entry binds the
  named group (the `Patch Group` tag value on managed nodes) to this
  baseline.
- **The default designation** folded as `setAsDefaultBaseline`:
  claiming makes this baseline the OS default; un-setting or deleting
  RESTORES AWS's own predefined default (a true revert).

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).

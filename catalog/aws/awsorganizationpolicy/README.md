<p align="center">
  <img src="logo.svg" alt="AWS Organization Policy" width="80"/>
</p>

# AWS Organization Policy

Manage an [AWS Organizations policy](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_manage_policies.html)
— a service control policy or any of its twelve sibling types —
together with its attachments to roots, organizational units, and
member accounts.

## What Gets Managed

- **The policy** (`p-...` — the import ID): its **policyName** (an
  explicit spec field — spaces are legal; renames apply in place), its
  **type** (SERVICE_CONTROL_POLICY when unset; immutable — the 2026
  set spans RCPs, declarative EC2, SecurityHub, Bedrock, S3, and
  more), the structured **content** document, and a description.
- **Attachments** folded as immutable per-target entries — an
  [AWS Organizational Unit](../awsorganizationalunit) reference by
  default, or literal root/OU/account IDs (each imports as
  `{target_id}:{policy_id}`).

The policy's type must be enabled on the organization
([AWS Organization](../awsorganization)'s `enabledPolicyTypes`) before
attachments succeed. AWS-managed policies (FullAWSAccess, ...) can
never be adopted.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).

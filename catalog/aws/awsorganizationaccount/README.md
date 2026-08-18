<p align="center">
  <img src="logo.svg" alt="AWS Organization Account" width="80"/>
</p>

# AWS Organization Account

Manage a [member account of an AWS Organization](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_manage_accounts.html)
— creating the account, placing it in the OU tree, and its
account-level settings: alternate contacts, primary contact
information, and opt-in region enablement.

## What Gets Managed

- **The account** (its 12-digit ID is the import ID): **accountName**
  (an explicit spec field — spaces are legal; renames apply in place),
  the root user's **email** (unique across ALL of AWS; immutable),
  **parentId** placement (an [AWS Organizational Unit](../awsorganizationalunit)
  reference; changes MOVE the account), the pre-created bootstrap
  **roleName** (write-once, no read API), **iamUserAccessToBilling**,
  and **createGovcloud**.
- **Alternate contacts** folded as three typed arms
  (billing/operations/security; each imports as
  `{account_id}/{TYPE}`).
- **The primary contact** folded (delete is a NO-OP — the last-written
  contact stays on file).
- **Opt-in regions** folded as per-region enablement entries
  (~60-minute operations; delete is a NO-OP).

**The delete contract**: `closeOnDeletion: false` (default) REMOVES
the account from the organization — it survives standalone;
`closeOnDeletion: true` CLOSES it (~90-day suspension, quota-limited).
Neither is a clean delete — treat member accounts as long-lived.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).

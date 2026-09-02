# CloudflareEmailRoutingAddress guide

Operational judgment for destination addresses. The README covers what each field is; this covers how the pieces interact.

## Verification is a human step, by design

Creating the address sends a verification email to the mailbox; until its owner clicks the link, the `verified` output stays empty and NO rule or catch-all can forward to it (Cloudflare rejects the configuration). Plan rollouts accordingly: create addresses first, have owners verify, then wire rules. There is no API to force-verify — the `status` field cannot do it for non-admin callers.

## Addresses are account-scoped and shared

One verified address serves every zone in the account. Model each real mailbox ONCE and reference it from all zones' rules — duplicating the same email as multiple resources in one account fails on create (Cloudflare enforces uniqueness per account).

## The status override is for un-verifying only

`status` exists for one real operation: flipping a verified address back to `unverified` (e.g. after a mailbox changes hands) to cut forwarding without deleting the address's history. Setting `verified` requires Cloudflare account admin privileges — do not reach for it as a verification shortcut. Note the Pulumi engine cannot send this field yet (SDK gap, documented in the module); use the Terraform engine when the override matters.

## A fresh address cannot be deleted for ~10 minutes

Cloudflare refuses to delete a destination address for roughly ten minutes after it is created — the API answers 400 with code 2032 "Destination address has been created too recently" (measured live 2026-08-26: a delete was refused at 9m14s and accepted at 10m15s after create). A destroy that runs inside that window fails with that error no matter which engine issued it; the address is not stuck — simply re-run the destroy once the window passes. Plan create-then-delete automation (previews, ephemeral environments) with this guard in mind.

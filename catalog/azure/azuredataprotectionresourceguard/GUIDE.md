# Azure Data Protection Resource Guard -- Operational Guide

Judgment calls that matter when you run Resource Guards in production.

## The guard is only as strong as its scope separation

Multi-User Authorization works because the approver is SOMEONE ELSE. Deploy the guard into a resource group (or better, a subscription) administered by a different person or team than the vaults it protects. An attacker -- or a mistaken script -- with rights over the vault's scope then still cannot approve its own destructive operations. A guard the vault's admin also controls adds a click, not a control.

## Empty exclusion list is the strong posture

With no exclusions, EVERY critical vault operation (disabling soft delete, deleting protected items, reducing retention) requires an approval through the guard. Add entries to `vaultCriticalOperationExclusionList` only for operations the org has deliberately decided to leave ungated -- and treat each entry as a standing security decision worth a comment in the manifest.

## Some operations can never be excluded

Azure keeps a mandatory always-guarded set and rejects any exclusion list naming one of them at create time (`BMSUserErrorInvalidCriticalOperationExclusionList`): the operations that would disarm the guard itself -- `backupconfig/write`, both vault families' `backupResourceGuardProxies/delete`, `backupVaults/write#reduceSoftDeleteSecurity` -- and the cross-tenant vault-mapping operations. The Terraform provider does not check this; Planton's validation mirrors ARM's list so a bad manifest fails at admission instead of mid-deploy. Note the asymmetry while composing a list: `backupconfig/write` is mandatory but `backupconfig/delete` is excludable, and `backupResourceGuardProxies/write` is excludable while its `/delete` sibling is not.

## One guard serves the estate

Vaults reference the guard by ARM ID; there is no per-vault guard object. One well-placed guard per environment (or per compliance boundary) keeps the approval path singular and auditable. Creating the guard changes nothing by itself -- protection starts when the first vault references it.

## Deleting the guard removes the gate

The guard deletes cleanly even while vaults reference it -- and with it goes the approval requirement. Treat guard deletion as itself a privileged operation: gate it with resource locks or pipeline approvals, because Azure will not gate the guard with itself.

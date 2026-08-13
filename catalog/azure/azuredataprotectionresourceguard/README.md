# Overview

The **AzureDataProtectionResourceGuard** component creates a Data Protection Resource Guard -- the approval gate behind Multi-User Authorization (MUA) for Azure Backup. Once a backup vault references a guard, privileged vault operations (disabling soft delete, deleting protection, reducing retention) require an approval THROUGH the guard before they execute. The guard's protection comes entirely from scope separation: place it where a DIFFERENT administrator controls it than the vaults it guards. The guard itself is a free configuration object.

## Purpose

- **A compromised admin cannot delete the backups**: destructive vault operations need a second party's approval when the guard lives in another administrator's scope.
- **One guard, many vaults**: vaults opt in by referencing the guard's ARM ID -- a single approval gate serves a whole environment's backup estate.
- **Governance as declarative infrastructure**: the exclusion list (which critical operations bypass the gate) is reviewed and versioned like everything else.

## Key Features

- Full azurerm v5 surface: the guard object with its critical-operation exclusion list and tags.
- Secure-by-default: an empty exclusion list guards EVERY critical operation -- exclusions are the deliberate exception, not the default.
- Chart-ready: publishes `resource_guard_id`, the reference backup vaults consume to enable Multi-User Authorization.

## Use Cases

- **Ransomware-resistant backup estates**: pair with immutable vaults so neither data nor its retention can be quietly destroyed.
- **Separation of duties**: backup operations owned by one team, destructive-operation approval owned by another.
- **Compliance regimes** that mandate dual authorization on data-destruction paths.

## Future Enhancements

- Typed guard references on the vault kinds tighten as the family's proof lanes complete.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).

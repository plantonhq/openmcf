# Overview

The **AzureBackupProtectedFileShare** component puts one Azure Files share under a backup policy's protection in a Recovery Services vault -- the binding that turns a policy from configuration into actual backups. Creating it only REGISTERS protection: the first backup runs on the policy's schedule, not immediately.

## Purpose

- **Protection as declarative infrastructure**: which share, which policy, which vault -- reviewed and versioned, never clicked together.
- **The full chain made explicit**: the share's storage account must be REGISTERED with the vault first (AzureBackupContainerStorageAccount); the spec's default reference wires through the registration so the order is automatic.
- **Typed references end to end**: vault by name, storage account through its registration, share by name, policy by ARM ID -- chart-ready.

## Key Features

- Full azurerm v5 surface (five arguments -- the binding is deliberately small).
- The provider's discovery reality on the record: creation runs an Inquire pass to find protectable shares inside the registered container, and fails loudly when the account is not registered.
- Honest lifecycle semantics: `backupPolicyId` is the ONLY updatable field; destroy deletes the backup data (vault soft delete may hold it 14 days); creates and deletes run in the provider's 80-minute timeout class.

## Use Cases

- **Protecting team and application file shares**: one binding per share, all riding one policy and one registration.
- **Charts**: an app chart that provisions its share and its protection together -- restore-ready from day one.
- **Selective protection**: protect the shares that matter in an account without touching the rest.

## Future Enhancements

- Restore flows (point-in-time to original or alternate location) are operational actions, not declarative state -- they stay in the portal/CLI by design.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).

# Azure Backup Protected VM -- Operational Guide

Judgment that saves real time when protecting VMs. The field reference lives in the API Explorer; this is the operational layer above it.

## The first backup is NOT at create time

Creating the protection registers the VM and nothing more -- the item sits in `IRPending` (initial replication pending) until the policy's next scheduled window runs the first full backup. A VM protected at 09:00 against a 23:00 policy is UNPROTECTED for fourteen hours. When onboarding critical machines, trigger the first backup manually (`az backup protection backup-now`) instead of waiting for the schedule.

## Destroy means "delete my backups" -- decide if that is what you want

Destroying this resource stops protection AND deletes the backup data (the engines' defaults, kept deliberately -- a binding whose destroy silently strands paid storage would be worse). Two escape hatches, both engine-level features rather than spec fields: `vm_backup_stop_protection_and_retain_data_on_destroy` keeps the data and stops backups; the `suspend` variant keeps the data under an immutable vault's rules. And remember the vault's own guard: it refuses to delete while protected items remain -- teardown order is protections first, vault last.

## Soft delete holds backup data for 14 days after deletion

Deleting protection soft-deletes the recovery points -- they hold vault storage (unbilled) and the VM's registration for 14 days. Re-protecting the same VM inside that window collides with the ghost: the provider recovers it automatically only when its `recover_soft_deleted_backup_protected_vm` feature is on; otherwise undelete manually (`az backup protection undelete`) or wait out the window.

## Disk filters: exclude for savings, include for surgical protection

`excludeDiskLuns` backs up everything except the listed data disks -- right for rebuildable content (caches, scratch, tempdb). `includeDiskLuns` inverts it: the OS disk plus ONLY the listed disks. They are mutually exclusive, and the OS disk is always backed up either way. Databases with native backup tooling are the classic exclude case -- do not pay twice for the same bytes.

## protection_state is a dial, not a status display

Leave it unset in normal operation (Azure manages it; transient states read back as Protected). `ProtectionStopped` halts backups while keeping existing data -- the "pause" for a VM being decommissioned gradually. `BackupsSuspended` does the same but ONLY on an immutable vault (Unlocked/Locked) -- Azure enforces that against the LIVE vault at apply, so flipping a mutable vault's protection to BackupsSuspended fails at deploy time, not manifest time.

## One VM, one vault -- Azure's rule, not ours

A VM can be protected by exactly one Recovery Services vault. Re-pointing the POLICY updates in place (`backupPolicyId` is the supported re-policy path); re-pointing the VAULT means delete-protection-here, wait out soft delete, protect-there. Plan vault topology before mass onboarding, not after.

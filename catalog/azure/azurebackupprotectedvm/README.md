# Overview

The **AzureBackupProtectedVm** component registers one virtual machine under a backup policy's protection in a Recovery Services vault -- the binding that turns "we have a backup policy" into "this VM is actually backed up". Creating it only registers protection: the first backup runs on the policy's schedule.

## Purpose

- **Protection as declarative infrastructure**: which VM, which policy, which disks -- reviewed and versioned, auditable at a glance.
- **Typed references end-to-end**: the vault, the policy, and the VM itself all wire by reference -- a chart can stand up VM + vault + policy + protection in one apply.
- **Honest destroy semantics**: destroying the binding stops protection AND deletes the backup data (the engines' default); the retain-on-destroy alternatives are documented where operators look.

## Key Features

- Full azurerm v5 surface: source VM, backup policy (re-pointable in place), disk-LUN include/exclude filters, and the protection-state dial (Protected / ProtectionStopped / BackupsSuspended).
- Service contracts recorded where they bite: BackupsSuspended requires an immutable vault (apply-time, checked against the live vault); transient Azure states (IRPending) read back as Protected; soft-deleted ghosts hold the registration.
- Selective disk backup: exclude data disks that hold rebuildable content, or include only the disks that matter -- both directions modeled, mutually exclusive.

## Use Cases

- **Every production VM**: the standard binding to the environment's daily policy.
- **Selective protection**: databases whose data disks are backed up natively -- exclude them, protect the OS disk.
- **Chart-composed stacks**: VM + protection created together so nothing ships unprotected.

## Future Enhancements

- The file-share protection kind extends the same binding pattern to Azure Files as its contract lands in the catalog.

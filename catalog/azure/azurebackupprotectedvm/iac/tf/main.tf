# Register the VM under the policy's protection. Creation only
# REGISTERS protection -- the first backup runs on the policy's
# schedule, not immediately. ARM derives the protected item's own name
# from the VM's group and name (VM;iaasvmcontainerv2;{vm-rg};{vm-name}).
#
# Destroy semantics kept deliberately at the engines' defaults:
# destroying stops protection AND deletes the backup data (the
# provider's retain-on-destroy features stay off in provider.tf) --
# recorded on the spec.
resource "azurerm_backup_protected_vm" "main" {
  resource_group_name = var.spec.resource_group
  recovery_vault_name = var.spec.recovery_vault_name
  source_vm_id        = var.spec.source_vm_id
  backup_policy_id    = var.spec.backup_policy_id

  # Mutually exclusive on the spec (CEL) and the provider
  # (ConflictsWith): the disks to skip OR the disks to keep.
  exclude_disk_luns = length(var.spec.exclude_disk_luns) > 0 ? var.spec.exclude_disk_luns : null
  include_disk_luns = length(var.spec.include_disk_luns) > 0 ? var.spec.include_disk_luns : null

  # Optional+Computed on the provider: unset lets Azure manage the
  # posture (transient states like IRPending read back as Protected).
  # BackupsSuspended additionally requires the VAULT to be immutable --
  # an apply-time contract Azure checks against the live vault
  # (recorded on the spec field).
  protection_state = var.spec.protection_state != "" ? var.spec.protection_state : null
}

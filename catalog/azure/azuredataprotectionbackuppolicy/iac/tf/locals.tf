locals {
  # The kubernetes_cluster variant is the one provider resource that
  # addresses the vault by NAME + resource group instead of by ID.
  # ARM vault IDs are structured
  # (/subscriptions/{sub}/resourceGroups/{rg}/providers/
  #  Microsoft.DataProtection/backupVaults/{name}), so both derive
  # deterministically from the spec's single vault_id reference.
  vault_id_parts            = split("/", var.spec.vault_id)
  vault_resource_group_name = local.vault_id_parts[4]
  vault_name                = local.vault_id_parts[8]
}

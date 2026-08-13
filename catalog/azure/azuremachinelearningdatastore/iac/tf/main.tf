# Create the Machine Learning datastore -- the saved connection that
# tells the workspace where data lives. ONE kind covers the three
# storage flavors; exactly one variant block is set (spec CEL) and
# selects which provider resource is created below. All three write
# the same ARM child collection
# (.../workspaces/{ws}/dataStores/{name}), so exactly one of the three
# resources exists per deployment.
#
# Credentials are WRITE-ONLY: ARM never returns account keys, SAS
# tokens or client secrets -- the provider echoes them from
# configuration (recorded in the import catalog as write-normalized).

# The blob-container variant. The only variant where is_default is
# settable; auth is account key or SAS unless a workspace-identity
# mode covers service-side access (spec CEL).
resource "azurerm_machine_learning_datastore_blobstorage" "main" {
  count = var.spec.blob_storage != null ? 1 : 0

  name                 = var.spec.name
  workspace_id         = var.spec.workspace_id
  storage_container_id = var.spec.blob_storage.storage_container_id

  description = var.spec.description != "" ? var.spec.description : null
  is_default  = var.spec.blob_storage.is_default

  # Enum name -> wire value; unspecified ("") omits the property so the
  # provider applies its default, "None".
  service_data_auth_identity = lookup(local.service_data_identity_wire, var.spec.service_data_identity, null)

  # Sensitive -- resolved from secret references, masked in plan
  # output. When both are set the provider sends the SAS token (its
  # own precedence).
  account_key             = var.spec.blob_storage.account_key != "" ? var.spec.blob_storage.account_key : null
  shared_access_signature = var.spec.blob_storage.shared_access_signature != "" ? var.spec.blob_storage.shared_access_signature : null

  tags = local.final_tags
}

# The Data Lake Gen2 variant. Auth is the service-principal triad
# (all-or-none, spec CEL) or workspace identity / none; no account
# key or SAS on this variant.
resource "azurerm_machine_learning_datastore_datalake_gen2" "main" {
  count = var.spec.data_lake_gen2 != null ? 1 : 0

  name                 = var.spec.name
  workspace_id         = var.spec.workspace_id
  storage_container_id = var.spec.data_lake_gen2.storage_container_id

  description = var.spec.description != "" ? var.spec.description : null

  service_data_identity = lookup(local.service_data_identity_wire, var.spec.service_data_identity, null)

  tenant_id     = var.spec.data_lake_gen2.tenant_id != "" ? var.spec.data_lake_gen2.tenant_id : null
  client_id     = var.spec.data_lake_gen2.client_id != "" ? var.spec.data_lake_gen2.client_id : null
  client_secret = var.spec.data_lake_gen2.client_secret != "" ? var.spec.data_lake_gen2.client_secret : null
  authority_url = var.spec.data_lake_gen2.authority_url != "" ? var.spec.data_lake_gen2.authority_url : null

  tags = local.final_tags
}

# The Azure Files variant. The provider's schema requires exactly one
# of account key / SAS here regardless of identity mode (spec CEL
# mirrors it). The share id uses the v5 format
# (.../fileServices/default/shares/{name}).
resource "azurerm_machine_learning_datastore_fileshare" "main" {
  count = var.spec.file_share != null ? 1 : 0

  name                 = var.spec.name
  workspace_id         = var.spec.workspace_id
  storage_fileshare_id = var.spec.file_share.storage_fileshare_id

  description = var.spec.description != "" ? var.spec.description : null

  service_data_identity = lookup(local.service_data_identity_wire, var.spec.service_data_identity, null)

  account_key             = var.spec.file_share.account_key != "" ? var.spec.file_share.account_key : null
  shared_access_signature = var.spec.file_share.shared_access_signature != "" ? var.spec.file_share.shared_access_signature : null

  tags = local.final_tags
}

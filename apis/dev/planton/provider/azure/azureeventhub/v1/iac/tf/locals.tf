locals {
  # Event hubs carry no Azure tags: ARM does not support tags on Event
  # Hubs entities (hubs/consumer groups/rules), so the platform's
  # identity tags live on the parent namespace.

  # Gate-state wire values. The tfvars wire format carries the FULL
  # proto enum value name; unset deploys Active.
  status_map = {
    "ACTIVE"        = "Active"
    "DISABLED"      = "Disabled"
    "SEND_DISABLED" = "SendDisabled"
  }
  status = (
    var.spec.status != null && var.spec.status != ""
    ? local.status_map[var.spec.status]
    : "Active"
  )

  # Cleanup-policy wire values -- the spec enum requires an explicit
  # choice when retention_description is declared, so no unspecified
  # fallback row is needed.
  cleanup_policy_map = {
    "DELETE"  = "Delete"
    "COMPACT" = "Compact"
  }

  # Capture encoding wire values -- required/explicit in the spec, no
  # fallback row needed.
  capture_encoding_map = {
    "AVRO"         = "Avro"
    "AVRO_DEFLATE" = "AvroDeflate"
  }

  # Capture storage-auth wire values. Unset keeps Azure's default
  # (service-managed SAS) -- "StorageSAS" is the provider's marker for
  # "send no identity", not an ARM value.
  capture_auth_map = {
    "STORAGE_SAS"     = "StorageSAS"
    "SYSTEM_ASSIGNED" = "SystemAssigned"
    "USER_ASSIGNED"   = "UserAssigned"
  }
  capture_auth = (
    var.spec.capture_description != null && var.spec.capture_description.destination.storage_authentication_type != null && var.spec.capture_description.destination.storage_authentication_type != ""
    ? local.capture_auth_map[var.spec.capture_description.destination.storage_authentication_type]
    : "StorageSAS"
  )
}

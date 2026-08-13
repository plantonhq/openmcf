locals {
  resource_id = (
    var.metadata.id != null && var.metadata.id != ""
    ? var.metadata.id
    : var.metadata.name
  )

  # PARITY-EXCEPTION: resource_kind here is the family-wide snake-case
  # literal and resource_id falls back to metadata.name, while the Pulumi
  # module emits the lowered CloudResourceKind enum string and omits
  # resource_id when metadata.id is empty. Output-neutral (tags never feed
  # stack outputs); aligning the two shapes is a family-wide convention
  # change, not a per-kind fix.
  base_tags = {
    "resource"      = "true"
    "resource_id"   = local.resource_id
    "resource_kind" = "azure_express_route_port"
    "resource_name" = var.metadata.name
  }

  org_tag = (
    var.metadata.org != null && var.metadata.org != ""
  ) ? { "organization" = var.metadata.org } : {}

  env_tag = (
    var.metadata.env != null && var.metadata.env != ""
  ) ? { "environment" = var.metadata.env } : {}

  # Metadata-derived tags first, then the user's spec tags merged over them:
  # user tags deliberately win so an org's governance conventions (cost
  # center, owner) can override the derived values where they collide.
  final_tags = merge(local.base_tags, local.org_tag, local.env_tag, var.spec.tags)

  # The spec's enum NAMES mapped onto ARM's vocabulary.
  encapsulation_wire = {
    "DOT1Q" = "Dot1Q"
    "QINQ"  = "QinQ"
  }
  billing_type_wire = {
    "METERED_DATA"   = "MeteredData"
    "UNLIMITED_DATA" = "UnlimitedData"
  }
  macsec_cipher_wire = {
    "GCM_AES_128"     = "GcmAes128"
    "GCM_AES_256"     = "GcmAes256"
    "GCM_AES_XPN_128" = "GcmAesXpn128"
    "GCM_AES_XPN_256" = "GcmAesXpn256"
  }

  # Unset (null) applies ARM's default -- MeteredData -- mirroring the
  # Pulumi module's nil handling.
  billing_type = (
    var.spec.billing_type == null
    ? "MeteredData"
    : lookup(local.billing_type_wire, var.spec.billing_type, var.spec.billing_type)
  )

  # Map the identity type enum's name string to ARM's comma-separated
  # value.
  identity_type = (
    var.spec.identity == null ? null :
    var.spec.identity.type == "SYSTEM_ASSIGNED" ? "SystemAssigned" :
    var.spec.identity.type == "USER_ASSIGNED" ? "UserAssigned" :
    var.spec.identity.type == "SYSTEM_AND_USER_ASSIGNED" ? "SystemAssigned, UserAssigned" : null
  )
}

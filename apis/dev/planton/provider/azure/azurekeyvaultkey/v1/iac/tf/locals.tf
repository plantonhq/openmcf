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
    "resource_kind" = "azure_key_vault_key"
    "resource_name" = var.metadata.name
  }

  org_tag = (
    var.metadata.org != null && var.metadata.org != ""
  ) ? { "organization" = var.metadata.org } : {}

  env_tag = (
    var.metadata.env != null && var.metadata.env != ""
  ) ? { "environment" = var.metadata.env } : {}

  # Metadata-derived tags first, then the user's spec tags merged over them:
  # user tags deliberately win so an org's governance conventions can
  # override the derived values where they collide.
  final_tags = merge(local.base_tags, local.org_tag, local.env_tag, var.spec.tags)

  # The spec enums arrive as FULL proto value names; each map below is the
  # complete verbatim vocabulary, translated to the exact strings Azure's
  # API expects (which is case-sensitive about all of them).
  key_type_map = {
    "RSA"     = "RSA"
    "RSA_HSM" = "RSA-HSM"
    "EC"      = "EC"
    "EC_HSM"  = "EC-HSM"
  }

  curve_map = {
    "P_256"  = "P-256"
    "P_256K" = "P-256K"
    "P_384"  = "P-384"
    "P_521"  = "P-521"
  }

  key_ops_map = {
    "DECRYPT"    = "decrypt"
    "ENCRYPT"    = "encrypt"
    "SIGN"       = "sign"
    "UNWRAP_KEY" = "unwrapKey"
    "VERIFY"     = "verify"
    "WRAP_KEY"   = "wrapKey"
  }

  key_type = local.key_type_map[var.spec.key_type]

  # null when unset so Azure applies its own curve default (P-256) --
  # identical behavior on both engines.
  curve = var.spec.curve != null ? local.curve_map[var.spec.curve] : null

  key_opts = [for op in var.spec.key_opts : local.key_ops_map[op]]
}

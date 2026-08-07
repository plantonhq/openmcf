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
    "resource_kind" = "azure_key_vault_certificate"
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

  # The spec enums arrive as FULL proto value names (the tfvars wire format
  # never strips prefixes); each map below is the complete verbatim
  # vocabulary, translated to the exact strings Azure's certificate API
  # expects (which is case-sensitive about all of them -- note the
  # lowerCamel key-usage extensions vs the UpperCamel action types).
  key_type_map = {
    "RSA"     = "RSA"
    "RSA_HSM" = "RSA-HSM"
    "EC"      = "EC"
    "EC_HSM"  = "EC-HSM"
    "OCT"     = "oct"
  }

  curve_map = {
    "P_256"  = "P-256"
    "P_256K" = "P-256K"
    "P_384"  = "P-384"
    "P_521"  = "P-521"
  }

  action_type_map = {
    "AUTO_RENEW"     = "AutoRenew"
    "EMAIL_CONTACTS" = "EmailContacts"
  }

  # The secret face's media type: what consumers reading the certificate's
  # secret get back.
  content_type_map = {
    "PKCS12" = "application/x-pkcs12"
    "PEM"    = "application/x-pem-file"
  }

  key_usage_map = {
    "CRL_SIGN"          = "cRLSign"
    "DATA_ENCIPHERMENT" = "dataEncipherment"
    "DECIPHER_ONLY"     = "decipherOnly"
    "DIGITAL_SIGNATURE" = "digitalSignature"
    "ENCIPHER_ONLY"     = "encipherOnly"
    "KEY_AGREEMENT"     = "keyAgreement"
    "KEY_CERT_SIGN"     = "keyCertSign"
    "KEY_ENCIPHERMENT"  = "keyEncipherment"
    "NON_REPUDIATION"   = "nonRepudiation"
  }
}

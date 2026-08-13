locals {
  # metadata.name IS the zone's domain name (ForceNew — a hosted zone cannot
  # be renamed).
  zone_name = var.metadata.name

  # Resource-identity tags follow the catalog convention.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsRoute53Zone"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Private-zone VPC set — vpc_region defaults to the zone's region so
  # single-region graphs never repeat it. Public zones carry no VPCs
  # (CEL-enforced), so this list is empty for them.
  vpc_associations = var.spec.is_private ? [
    for association in var.spec.vpc_associations : {
      vpc_id     = association.vpc_id
      vpc_region = association.vpc_region != "" ? association.vpc_region : var.spec.region
    }
  ] : []

  # Null-when-unset scalars so the provider applies its own defaults instead
  # of this module freezing them.
  comment           = var.spec.comment != "" ? var.spec.comment : null
  delegation_set_id = var.spec.delegation_set_id != "" ? var.spec.delegation_set_id : null

  # DNSSEC — the KSK name defaults to a zone-derived name when omitted, and
  # the KSK status defaults to ACTIVE (the provider default; INACTIVE is the
  # documented diagnostics lever).
  dnssec_enabled = var.spec.dnssec != null
  ksk_name = local.dnssec_enabled ? (
    var.spec.dnssec.key_signing_key_name != "" ? var.spec.dnssec.key_signing_key_name : "${replace(local.zone_name, ".", "-")}-ksk"
  ) : null
  ksk_status = local.dnssec_enabled ? (
    var.spec.dnssec.key_signing_key_status != "" ? var.spec.dnssec.key_signing_key_status : "ACTIVE"
  ) : null

  query_logging_enabled = var.spec.query_logging != null
}

# Create one DNS record set in an Azure PRIVATE DNS zone.
#
# The record type is whichever typed payload the spec carries (validation
# guarantees exactly one), so exactly one of the count-gated resources
# below materializes. Private record sets are addressed by the ZONE's ARM
# id plus the record name -- unlike the public DNS resources, which
# address by (resource group, zone name).
#
# Private DNS has no alias records (that is a public-DNS concept) and
# supports exactly these seven types -- no CAA, no NS (private zones
# cannot delegate subdomains).

# IPv4 address record -- Azure caps an A record set at 20 addresses
# (spec-validated).
resource "azurerm_private_dns_a_record" "main" {
  count = length(var.spec.a) > 0 ? 1 : 0

  name                = var.spec.name
  private_dns_zone_id = var.spec.private_dns_zone_id
  ttl                 = local.ttl
  records             = var.spec.a
  tags                = local.final_tags
}

# IPv6 address record.
resource "azurerm_private_dns_aaaa_record" "main" {
  count = length(var.spec.aaaa) > 0 ? 1 : 0

  name                = var.spec.name
  private_dns_zone_id = var.spec.private_dns_zone_id
  ttl                 = local.ttl
  records             = var.spec.aaaa
  tags                = local.final_tags
}

# Canonical-name record -- one target hostname. cname is a
# StringValueOrRef in the spec; the tfvars converter flattens it to the
# resolved literal, so the module reads a plain string.
resource "azurerm_private_dns_cname_record" "main" {
  count = var.spec.cname != null && var.spec.cname != "" ? 1 : 0

  name                = var.spec.name
  private_dns_zone_id = var.spec.private_dns_zone_id
  ttl                 = local.ttl
  record              = var.spec.cname
  tags                = local.final_tags
}

# Mail-exchange record set -- each entry carries its own preference, so
# multi-server mail setups (10 primary / 20 secondary) express exactly.
resource "azurerm_private_dns_mx_record" "main" {
  count = length(var.spec.mx) > 0 ? 1 : 0

  name                = var.spec.name
  private_dns_zone_id = var.spec.private_dns_zone_id
  ttl                 = local.ttl
  tags                = local.final_tags

  dynamic "record" {
    for_each = var.spec.mx
    content {
      preference = record.value.preference
      exchange   = record.value.exchange
    }
  }
}

# Pointer record set -- reverse DNS (IP-to-name) in private in-addr.arpa
# / ip6.arpa zones.
resource "azurerm_private_dns_ptr_record" "main" {
  count = length(var.spec.ptr) > 0 ? 1 : 0

  name                = var.spec.name
  private_dns_zone_id = var.spec.private_dns_zone_id
  ttl                 = local.ttl
  records             = var.spec.ptr
  tags                = local.final_tags
}

# Service-locator record set -- priority/weight/port/target per endpoint.
resource "azurerm_private_dns_srv_record" "main" {
  count = length(var.spec.srv) > 0 ? 1 : 0

  name                = var.spec.name
  private_dns_zone_id = var.spec.private_dns_zone_id
  ttl                 = local.ttl
  tags                = local.final_tags

  dynamic "record" {
    for_each = var.spec.srv
    content {
      priority = record.value.priority
      weight   = record.value.weight
      port     = record.value.port
      target   = record.value.target
    }
  }
}

# Text record set -- each value up to 1,024 characters (Azure's private
# DNS cap). Each value is a StringValueOrRef in the spec; the tfvars
# converter flattens the list to resolved literals, so the module reads
# plain strings.
resource "azurerm_private_dns_txt_record" "main" {
  count = length(var.spec.txt) > 0 ? 1 : 0

  name                = var.spec.name
  private_dns_zone_id = var.spec.private_dns_zone_id
  ttl                 = local.ttl
  tags                = local.final_tags

  dynamic "record" {
    for_each = var.spec.txt
    content {
      value = record.value
    }
  }
}

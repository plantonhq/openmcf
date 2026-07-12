# Create one DNS record set in an Azure public DNS zone.
#
# The record type is whichever typed payload the spec carries (validation
# guarantees exactly one), so exactly one of the count-gated resources
# below materializes. Azure's management plane addresses record sets by
# (resource group, zone name, type, record name) -- there is no ARM-id
# addressing mode for record sets on either engine.
#
# Alias records (A/AAAA/CNAME only): when the payload carries
# target_resource_id instead of literal values, Azure keeps the answer in
# sync with the referenced resource -- no drift window when a Public IP's
# address changes, and a way to point the zone APEX at an Azure resource
# where DNS itself forbids CNAME. The provider requires exactly one of
# records/target_resource_id, which spec validation already guarantees;
# empty collections are passed as null so the provider never sees an
# empty-but-present argument.

# IPv4 address record -- literal addresses or an Azure-resource alias.
resource "azurerm_dns_a_record" "main" {
  count = var.spec.a != null ? 1 : 0

  name                = var.spec.name
  zone_name           = var.spec.zone_name
  resource_group_name = var.spec.resource_group
  ttl                 = local.ttl
  records             = length(var.spec.a.addresses) > 0 ? var.spec.a.addresses : null
  target_resource_id  = var.spec.a.target_resource_id != "" ? var.spec.a.target_resource_id : null
  tags                = local.final_tags
}

# IPv6 address record -- literal addresses or an Azure-resource alias.
resource "azurerm_dns_aaaa_record" "main" {
  count = var.spec.aaaa != null ? 1 : 0

  name                = var.spec.name
  zone_name           = var.spec.zone_name
  resource_group_name = var.spec.resource_group
  ttl                 = local.ttl
  records             = length(var.spec.aaaa.addresses) > 0 ? var.spec.aaaa.addresses : null
  target_resource_id  = var.spec.aaaa.target_resource_id != "" ? var.spec.aaaa.target_resource_id : null
  tags                = local.final_tags
}

# Canonical-name record -- one target hostname or an Azure-resource alias.
resource "azurerm_dns_cname_record" "main" {
  count = var.spec.cname != null ? 1 : 0

  name                = var.spec.name
  zone_name           = var.spec.zone_name
  resource_group_name = var.spec.resource_group
  ttl                 = local.ttl
  record              = var.spec.cname.value != "" ? var.spec.cname.value : null
  target_resource_id  = var.spec.cname.target_resource_id != "" ? var.spec.cname.target_resource_id : null
  tags                = local.final_tags
}

# Mail-exchange record set -- each entry carries its own preference, so
# multi-server mail setups (10 primary / 20 secondary) express exactly.
# The provider's preference attribute is string-typed; the spec's integer
# is converted here.
resource "azurerm_dns_mx_record" "main" {
  count = length(var.spec.mx) > 0 ? 1 : 0

  name                = var.spec.name
  zone_name           = var.spec.zone_name
  resource_group_name = var.spec.resource_group
  ttl                 = local.ttl
  tags                = local.final_tags

  dynamic "record" {
    for_each = var.spec.mx
    content {
      preference = tostring(record.value.preference)
      exchange   = record.value.exchange
    }
  }
}

# Service-locator record set -- priority/weight/port/target per endpoint.
resource "azurerm_dns_srv_record" "main" {
  count = length(var.spec.srv) > 0 ? 1 : 0

  name                = var.spec.name
  zone_name           = var.spec.zone_name
  resource_group_name = var.spec.resource_group
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

# Certificate-authority-authorization record set -- which CAs may issue
# certificates for this name.
resource "azurerm_dns_caa_record" "main" {
  count = length(var.spec.caa) > 0 ? 1 : 0

  name                = var.spec.name
  zone_name           = var.spec.zone_name
  resource_group_name = var.spec.resource_group
  ttl                 = local.ttl
  tags                = local.final_tags

  dynamic "record" {
    for_each = var.spec.caa
    content {
      flags = record.value.flags
      tag   = local.caa_tag_map[record.value.tag]
      value = record.value.value
    }
  }
}

# Text record set -- SPF, DKIM, DMARC, domain verification. Values up to
# 4096 characters are legal: the provider transparently splits each into
# the 254-character strings DNS requires and reassembles them on read.
resource "azurerm_dns_txt_record" "main" {
  count = length(var.spec.txt) > 0 ? 1 : 0

  name                = var.spec.name
  zone_name           = var.spec.zone_name
  resource_group_name = var.spec.resource_group
  ttl                 = local.ttl
  tags                = local.final_tags

  dynamic "record" {
    for_each = var.spec.txt
    content {
      value = record.value
    }
  }
}

# Name-server record set -- delegates a CHILD subdomain to another zone's
# name servers. The zone's own apex NS records are Azure-managed.
resource "azurerm_dns_ns_record" "main" {
  count = length(var.spec.ns) > 0 ? 1 : 0

  name                = var.spec.name
  zone_name           = var.spec.zone_name
  resource_group_name = var.spec.resource_group
  ttl                 = local.ttl
  records             = var.spec.ns
  tags                = local.final_tags
}

# Pointer record set -- reverse DNS (IP-to-name) in in-addr.arpa /
# ip6.arpa zones.
resource "azurerm_dns_ptr_record" "main" {
  count = length(var.spec.ptr) > 0 ? 1 : 0

  name                = var.spec.name
  zone_name           = var.spec.zone_name
  resource_group_name = var.spec.resource_group
  ttl                 = local.ttl
  records             = var.spec.ptr
  tags                = local.final_tags
}

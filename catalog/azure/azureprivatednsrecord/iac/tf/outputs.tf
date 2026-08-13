# Exactly one of the seven count-gated record resources materializes, so
# each output coalesces PER ATTRIBUTE across the variants (whole-object
# coalescing would taint the outputs with the union of every variant's
# sensitivity).

output "record_id" {
  description = "The Azure Resource Manager ID of the record set. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/privateDnsZones/{zone}/{TYPE}/{name}"
  value = coalesce(
    try(azurerm_private_dns_a_record.main[0].id, null),
    try(azurerm_private_dns_aaaa_record.main[0].id, null),
    try(azurerm_private_dns_cname_record.main[0].id, null),
    try(azurerm_private_dns_mx_record.main[0].id, null),
    try(azurerm_private_dns_ptr_record.main[0].id, null),
    try(azurerm_private_dns_srv_record.main[0].id, null),
    try(azurerm_private_dns_txt_record.main[0].id, null),
  )
}

output "fqdn" {
  description = "The fully qualified domain name of the record set, with a trailing dot (e.g. db.internal.contoso.com.) -- resolvable only from virtual networks linked to the zone."
  value = coalesce(
    try(azurerm_private_dns_a_record.main[0].fqdn, null),
    try(azurerm_private_dns_aaaa_record.main[0].fqdn, null),
    try(azurerm_private_dns_cname_record.main[0].fqdn, null),
    try(azurerm_private_dns_mx_record.main[0].fqdn, null),
    try(azurerm_private_dns_ptr_record.main[0].fqdn, null),
    try(azurerm_private_dns_srv_record.main[0].fqdn, null),
    try(azurerm_private_dns_txt_record.main[0].fqdn, null),
  )
}

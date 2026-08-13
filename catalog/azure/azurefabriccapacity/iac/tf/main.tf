# Create the Microsoft Fabric capacity -- the billing and compute
# anchor of Microsoft Fabric. A running capacity bills PER HOUR from
# the moment it exists; the F-SKU scales up and down in place. "Fabric"
# is the SKU tier's only legal value at v5 (deliberately not part of
# the spec); the administrators list is required non-empty by the spec
# (Azure rejects a capacity created with none, and clearing it on a
# running capacity is a lockout, not a configuration).
resource "azurerm_fabric_capacity" "main" {
  name                = var.spec.name
  resource_group_name = var.spec.resource_group
  location            = var.spec.region

  administration_members = var.spec.administration_members

  sku {
    name = var.spec.sku_name
    tier = "Fabric"
  }

  tags = local.final_tags
}

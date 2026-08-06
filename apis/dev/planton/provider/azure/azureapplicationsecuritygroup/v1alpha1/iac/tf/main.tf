# Create the application security group -- an empty, named grouping of
# network interfaces that NSG rules can target by name instead of by IP.
#
# The group carries no members: membership is declared from the member side
# (a network interface lists the ASGs it joins; an NSG rule references
# source/destination ASGs), which is what makes the ASG a stable
# composition anchor. Everything except tags is fixed at creation; changing
# name or region replaces the group and every rule/NIC referencing it must
# be re-pointed.
resource "azurerm_application_security_group" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  tags = local.final_tags
}

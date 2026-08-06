# Create the network security group -- the stateful firewall that filters
# traffic for every subnet and NIC it guards.
#
# Lifecycle notes worth knowing before operating this resource:
# - Rules and tags update IN PLACE and take effect immediately for every
#   subnet and NIC the group guards. Name, region, and resource group are
#   the group's ARM identity; changing any of them replaces it, detaching
#   it from every subnet until re-attached.
# - The group itself is a shell: with no rules, Azure's implicit defaults
#   govern (allow VNet-internal and load-balancer traffic, deny other
#   inbound, allow all outbound).
# - The subnet-side attachment is deliberately not modeled here: a subnet
#   declares which NSG guards it (matching Azure's model), so one group
#   serves many subnets without listing them.
resource "azurerm_network_security_group" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group
  tags                = local.final_tags
}

# Each rule is a standalone resource under the group: the standalone form
# enforces the singular/plural conflicts at plan time and gives each rule
# its own state identity and plan line. (The Pulumi module manages the same
# rules inline on the group -- its standalone rule type flattens the
# application-security-group lists -- but both engines put the IDENTICAL
# rule set into ARM, and removing the last rule removes it on both.)
resource "azurerm_network_security_rule" "rules" {
  for_each = { for rule in var.spec.security_rules : rule.name => rule }

  name                        = each.value.name
  resource_group_name         = var.spec.resource_group
  network_security_group_name = azurerm_network_security_group.main.name

  description = each.value.description
  priority    = each.value.priority
  direction   = local.direction_to_arm[each.value.direction]
  access      = local.access_to_arm[each.value.access]
  protocol    = local.protocol_to_arm[each.value.protocol]

  # Ports: the spec guarantees at most one form is set (exactly one for
  # destination). An unset source means any -- "*" is sent so both engines
  # deploy the identical rule.
  source_port_range = (
    length(each.value.source_port_ranges) > 0 ? null :
    each.value.source_port_range != null ? each.value.source_port_range : "*"
  )
  source_port_ranges      = length(each.value.source_port_ranges) > 0 ? each.value.source_port_ranges : null
  destination_port_range  = length(each.value.destination_port_ranges) > 0 ? null : each.value.destination_port_range
  destination_port_ranges = length(each.value.destination_port_ranges) > 0 ? each.value.destination_port_ranges : null

  # Addressing: the spec guarantees at most one style per side (single
  # prefix, prefix list, or application security groups). All unset means
  # any -- "*" is sent so both engines deploy the identical rule.
  source_address_prefix = (
    length(each.value.source_application_security_group_ids) > 0 ? null :
    length(each.value.source_address_prefixes) > 0 ? null :
    each.value.source_address_prefix != null ? each.value.source_address_prefix : "*"
  )
  source_address_prefixes               = length(each.value.source_address_prefixes) > 0 ? each.value.source_address_prefixes : null
  source_application_security_group_ids = length(each.value.source_application_security_group_ids) > 0 ? each.value.source_application_security_group_ids : null

  destination_address_prefix = (
    length(each.value.destination_application_security_group_ids) > 0 ? null :
    length(each.value.destination_address_prefixes) > 0 ? null :
    each.value.destination_address_prefix != null ? each.value.destination_address_prefix : "*"
  )
  destination_address_prefixes               = length(each.value.destination_address_prefixes) > 0 ? each.value.destination_address_prefixes : null
  destination_application_security_group_ids = length(each.value.destination_application_security_group_ids) > 0 ? each.value.destination_application_security_group_ids : null
}

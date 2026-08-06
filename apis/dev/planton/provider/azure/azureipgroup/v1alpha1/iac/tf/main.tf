# Create the IP Group -- a named, reusable set of IP addresses and CIDR
# ranges that Azure Firewall and Firewall Policy rules reference by ARM id
# instead of repeating literal address lists.
#
# The group is passive: it carries no rules of its own. Consumption is
# declared from the rule's side (a firewall policy rule lists
# source/destination IP Groups), which is what makes the group a stable
# composition anchor -- its address list updates in place and every
# referencing rule follows immediately, without touching the rules
# themselves. Renaming or moving the group replaces it and every rule that
# referenced it must be re-pointed.
resource "azurerm_ip_group" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # Single addresses ("203.0.113.7") and CIDR blocks ("10.10.0.0/16") are
  # both legal entries; Azure caps a group at 5,000 entries.
  cidrs = var.spec.cidrs

  tags = local.final_tags
}

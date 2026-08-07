locals {
  # Enum wire maps. tfvars carries FULL proto enum value names; the maps
  # translate them to azurerm's exact (case-sensitive) vocabulary.
  filter_action_wire = {
    "ALLOW" = "Allow"
    "DENY"  = "Deny"
  }

  # Network/DNAT protocols: Azure spells the wildcard "Any" with
  # TCP/UDP/ICMP uppercase -- an irregular casing the provider validates
  # case-sensitively.
  rule_protocol_wire = {
    "ANY"  = "Any"
    "TCP"  = "TCP"
    "UDP"  = "UDP"
    "ICMP" = "ICMP"
  }

  # Application L7 protocol types: mixed-case wire values.
  application_protocol_type_wire = {
    "HTTP"  = "Http"
    "HTTPS" = "Https"
    "MSSQL" = "Mssql"
  }

  # No tags: ARM does not support tags on rule collection groups -- they
  # are child documents of the policy, not top-level tracked resources.
}

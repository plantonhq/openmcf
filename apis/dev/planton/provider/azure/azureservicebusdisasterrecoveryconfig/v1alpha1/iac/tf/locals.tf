locals {
  # Geo-DR configs carry no Azure tags: ARM does not support tags on
  # disasterRecoveryConfigs, so the platform's identity tags live on the
  # paired namespaces.

  # The optional alias rule -- sent only when non-empty so the pairing
  # defaults to the namespace's built-in root rule (Azure's own
  # behavior) when unset.
  alias_authorization_rule_id = (
    var.spec.alias_authorization_rule_id == null || var.spec.alias_authorization_rule_id == ""
    ? null
    : var.spec.alias_authorization_rule_id
  )
}

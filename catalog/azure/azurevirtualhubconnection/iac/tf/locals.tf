# The hub connection carries no tags and no resource group of its own:
# ARM addresses it as a child of the hub, and the provider's schema has
# no tags argument -- so this module needs none of the family's usual
# tag-merging locals.
locals {
  # The spec's enum NAMES mapped onto ARM's vocabulary. Unset applies
  # ARM's default (Contains) explicitly so the rendered plan shows the
  # real value -- mirroring the Pulumi module's nil handling. ARM fixes
  # the criteria once the connection is created.
  override_criteria_wire = {
    "CONTAINS" = "Contains"
    "EQUAL"    = "Equal"
  }
}

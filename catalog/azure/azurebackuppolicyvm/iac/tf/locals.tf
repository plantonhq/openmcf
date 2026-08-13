# The backup policy resource carries NO tags argument on the provider
# (ARM backup policies are untagged), so this module derives no tag
# map -- deliberately unlike its vault sibling.
locals {
  # Whether the tiering policy's archive rule is configured -- used to
  # render the two-level tiering_policy block.
  has_tiering_policy = var.spec.tiering_policy != null
}

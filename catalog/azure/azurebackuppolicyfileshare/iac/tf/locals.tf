# The backup policy resource carries NO tags argument on the provider
# (ARM backup policies are untagged), so this module derives no tag
# map. Reference fields arrive pre-resolved (the platform middleware
# resolves valueFrom references before IaC modules run), so the spec
# values are literal names here.
locals {
  # Kept for symmetry with the family's modules; nothing derived.
}

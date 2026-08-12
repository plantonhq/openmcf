# The association carries NO tags argument on the provider (ARM
# extension resources are untagged), so this module derives no tag map.
# Reference fields arrive pre-resolved (the platform middleware
# resolves valueFrom references before IaC modules run), so the spec
# values are literal ARM ids here.
locals {
  # Kept for symmetry with the family's modules; nothing derived.
}

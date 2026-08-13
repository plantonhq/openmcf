# The virtual network link carries NO tags argument on the provider
# (ARM gives ruleset links a free-form metadata map instead, modeled
# as spec.metadata), so this module derives no tag map. Reference
# fields arrive pre-resolved (the platform middleware resolves
# valueFrom references before IaC modules run), so the spec values are
# literal ARM IDs here.
locals {
  # Kept for symmetry with the family's modules; nothing derived.
}

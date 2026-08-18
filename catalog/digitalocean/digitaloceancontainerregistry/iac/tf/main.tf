# DigitalOcean Container Registry (DOCR).
#
# A DigitalOcean account holds exactly ONE container registry, and registry
# names are globally unique across ALL DigitalOcean accounts. The registry
# name and region are create-only; only the subscription tier can change
# after creation.
resource "digitalocean_container_registry" "registry" {
  name                   = var.spec.name
  subscription_tier_slug = var.spec.subscription_tier
  region                 = local.region
}

# Docker credentials, minted only when the spec asks for them. Neither knob is
# recoverable from the DigitalOcean API afterwards (the API only ever returns
# a freshly minted credential), and import of this resource is defective at
# the current provider pin -- see the kind's import map.
resource "digitalocean_container_registry_docker_credentials" "credentials" {
  count = local.create_docker_credentials ? 1 : 0

  registry_name = digitalocean_container_registry.registry.name
  write         = var.spec.docker_credentials.write
  # Unset defers to the provider default: the API maximum (~50 years).
  expiry_seconds = var.spec.docker_credentials.expiry_seconds
}

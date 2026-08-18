locals {
  # The proto enum value names ARE the DigitalOcean tier slugs
  # (starter/basic/professional) and region slugs (nyc3/sfo3/...), so both
  # pass through unmodified. Region is Optional+Computed at the provider:
  # unset must be sent as null (an empty string fails provider validation)
  # so DigitalOcean picks the region.
  region = var.spec.region != "" ? var.spec.region : null

  create_docker_credentials = var.spec.docker_credentials != null
}

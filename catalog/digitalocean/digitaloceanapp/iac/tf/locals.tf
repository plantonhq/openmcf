locals {
  # Proto enum identifiers (docker_hub, pre_deploy) become the provider's
  # uppercase tokens. Unspecified values are omitted so the provider default
  # applies.
  region = var.spec.region
}

# Convert one env list into the provider's env block shape.
# secret wins over plaintext when both somehow appear.

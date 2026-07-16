locals {
  # Private-link target types arrive as the spec enum's FULL value
  # names; the Terraform provider wants snake_case secondaries
  # (blob_secondary/web_secondary) where the pulumi bridge wants
  # camelCase -- each engine maps to its own provider's dialect, landing
  # on the same ARM group id. "Gateway" is capitalized in Azure's own
  # vocabulary.
  private_link_target_type_map = {
    "SITES"                = "sites"
    "BLOB"                 = "blob"
    "BLOB_SECONDARY"       = "blob_secondary"
    "WEB"                  = "web"
    "WEB_SECONDARY"        = "web_secondary"
    "MANAGED_ENVIRONMENTS" = "managedEnvironments"
    "GATEWAY"              = "Gateway"
  }

  # The provider requires certificate_name_check_enabled explicitly;
  # absent in tfvars means the spec's documented default (true).
  certificate_name_check_enabled = coalesce(var.spec.certificate_name_check_enabled, true)

  # No Azure tags: ARM does not support tags on Front Door origins, so
  # the platform's identity tags live on the profile.
}

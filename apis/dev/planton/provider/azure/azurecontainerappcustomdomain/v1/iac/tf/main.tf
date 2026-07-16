# Bind a custom domain to a Container App.
#
# Lifecycle notes worth knowing before operating this resource:
# - CREATE BLOCKS ON DOMAIN VALIDATION: Azure verifies ownership of the
#   hostname during the operation, so the asuid TXT record (carrying the
#   app's custom_domain_verification_id) and the CNAME/A routing record
#   must resolve publicly BEFORE this resource deploys, or the apply
#   fails. The app must have ingress enabled.
# - Every field replaces the binding when changed -- Azure models it as an
#   entry in the app's ingress configuration with no update surface.
#
# The binding is dispatched between two variants because Terraform's
# lifecycle block is static per resource and the two TLS flows need
# opposite drift handling:
# - BYO variant: the certificate reference is user-declared state -- drift
#   on it is real and must plan, so no ignore_changes.
# - MANAGED variant: Azure attaches the issued managed certificate to the
#   binding asynchronously, out of band; without ignore_changes on the
#   certificate fields every refresh after that attachment would read as
#   drift and plan a replacement. Ignoring these fields is safe here ONLY
#   because the spec declares neither.
# Exactly one variant materializes; outputs coalesce per attribute.

# Bring-your-own-certificate binding: serve TLS with a certificate stored
# on the app's environment (AzureContainerAppEnvironmentCertificate).
resource "azurerm_container_app_custom_domain" "byo" {
  count = local.is_byo_certificate ? 1 : 0

  name                                     = var.spec.domain_name
  container_app_id                         = var.spec.container_app_id
  container_app_environment_certificate_id = var.spec.container_app_environment_certificate_id
  certificate_binding_type                 = local.certificate_binding_type_map[var.spec.certificate_binding_type]
}

# Managed-certificate binding: deploy certificate-less; Azure fills the
# certificate binding in asynchronously once the matching
# AzureContainerAppEnvironmentManagedCertificate issues.
resource "azurerm_container_app_custom_domain" "managed" {
  count = local.is_byo_certificate ? 0 : 1

  name             = var.spec.domain_name
  container_app_id = var.spec.container_app_id

  lifecycle {
    # Azure mutates these two attributes out of band when it attaches the
    # managed certificate -- the provider documents exactly this ignore
    # for the managed flow.
    ignore_changes = [certificate_binding_type, container_app_environment_certificate_id]
  }
}

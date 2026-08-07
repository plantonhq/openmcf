# Provision an Azure-managed TLS certificate on a Container App
# Environment -- free, domain-validated, renewed by Azure.
#
# Lifecycle notes worth knowing before operating this resource:
# - CREATE BLOCKS ON DOMAIN VALIDATION: Azure only completes the operation
#   once it has proven you control subject_name, so the apply polls until
#   the required DNS records resolve publicly (the asuid TXT record with
#   the app's custom_domain_verification_id, plus the CNAME/HTTP routing
#   record) -- or fails around the 30-minute mark. Publish the records
#   BEFORE deploying this resource.
# - Only tags update in place; every other change re-issues the
#   certificate.
# - Azure attaches the issued certificate to the matching
#   AzureContainerAppCustomDomain binding asynchronously -- the binding
#   deploys certificate-less first, then Azure fills it in.
resource "azurerm_container_app_environment_managed_certificate" "main" {
  name                         = var.spec.certificate_name
  container_app_environment_id = var.spec.container_app_environment_id
  subject_name                 = var.spec.subject_name
  domain_control_validation    = local.domain_control_validation

  tags = local.final_tags
}

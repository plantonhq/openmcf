# The environment is the secure boundary every app, job, storage
# registration, and Dapr component in it shares. Name, region, resource
# group, subnet, ILB, and zone redundancy are all ForceNew -- recreating
# the environment takes every workload in it down.
resource "azurerm_container_app_environment" "main" {
  name                = var.spec.environment_name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # Logging destination: explicit choice honored as-is; unset with a
  # workspace deploys log-analytics; unset without one is streaming-only
  # (resolution in locals.tf).
  log_analytics_workspace_id = var.spec.log_analytics_workspace_id
  logs_destination           = local.logs_destination

  # The Dapr telemetry connection string is write-only in ARM (never
  # returned on read), which is why it is ForceNew.
  dapr_application_insights_connection_string = var.spec.dapr_application_insights_connection_string

  # VNet injection: the subnet must be /21 or larger; ILB and zone
  # redundancy only exist for VNet-injected environments -- the provider
  # requires the subnet whenever either is SPECIFIED (even as false), so
  # both are sent only alongside the subnet.
  infrastructure_subnet_id       = var.spec.infrastructure_subnet_id
  internal_load_balancer_enabled = var.spec.infrastructure_subnet_id != null ? var.spec.internal_load_balancer_enabled : null
  zone_redundancy_enabled        = var.spec.infrastructure_subnet_id != null ? var.spec.zone_redundancy_enabled : null

  # Only valid alongside workload profiles (spec-enforced); when omitted
  # Azure generates the platform resource-group name itself.
  infrastructure_resource_group_name = var.spec.infrastructure_resource_group_name

  # Unset lets Azure derive the value from the network configuration
  # (Enabled externally, Disabled behind an ILB).
  public_network_access = local.public_network_access

  # mTLS encrypts and authenticates all app-to-app traffic in the
  # environment; costs latency and peak throughput under load.
  mutual_tls_enabled = var.spec.mutual_tls_enabled

  # The environment-level managed identity (used by platform operations,
  # e.g. Key Vault certificate reads). The spec's CEL guarantees identity
  # ids are present exactly when the type includes UserAssigned.
  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_map[identity.value.type]
      identity_ids = identity.value.user_assigned_identity_ids
    }
  }

  # Dedicated / GPU compute pools. Azure always includes the standard
  # Consumption profile itself, so only the declared profiles are sent.
  # The Consumption-family profiles are serverless: Azure rejects
  # instance counts on them, which the spec's CEL front-loads.
  dynamic "workload_profile" {
    for_each = var.spec.workload_profiles
    content {
      name                  = workload_profile.value.name
      workload_profile_type = local.workload_profile_type_map[workload_profile.value.workload_profile_type]
      minimum_count         = workload_profile.value.minimum_count
      maximum_count         = workload_profile.value.maximum_count
    }
  }

  tags = local.final_tags
}

# The custom DNS suffix is a property of the environment itself (ARM
# patches the environment's CustomDomainConfiguration; the resource's id
# IS the environment id), so the spec folds it in and the module realizes
# it through the association resource.
resource "azurerm_container_app_environment_custom_domain" "main" {
  count = var.spec.custom_domain != null ? 1 : 0

  container_app_environment_id = azurerm_container_app_environment.main.id
  dns_suffix                   = var.spec.custom_domain.dns_suffix
  certificate_blob_base64      = var.spec.custom_domain.certificate_blob_base64
  certificate_password         = var.spec.custom_domain.certificate_password
}

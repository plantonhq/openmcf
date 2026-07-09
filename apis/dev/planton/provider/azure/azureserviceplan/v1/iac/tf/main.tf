# The App Service Plan: the compute tier every Web App / Function App on
# it shares. The SKU is deliberately NOT ForceNew -- plans scale up, down,
# and across tiers in place (Azure rejects the few impossible moves, like
# Consumption <-> dedicated, at apply time). Name, OS type, region, and
# resource group ARE ForceNew, and recreating the plan takes every app on
# it down with it.
resource "azurerm_service_plan" "main" {
  name                = var.spec.service_plan_name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group
  os_type             = local.os_type
  sku_name            = local.sku_name

  # Placing the plan inside an App Service Environment v3 (single-tenant
  # compute). The spec gates this to Isolated SKUs, mirroring Azure's own
  # creation-time rule.
  app_service_environment_id = var.spec.app_service_environment_id

  # Unset lets Azure apply the SKU's default capacity (typically 1); the
  # serverless tiers (Y1/FC1/EP*) manage instance count themselves.
  worker_count = var.spec.worker_count

  # Premium-plan automatic HTTP-load scaling; the ceiling comes from
  # maximum_elastic_worker_count. Both gates live in the spec, mirroring
  # the provider's own SKU checks.
  premium_plan_auto_scale_enabled = var.spec.premium_plan_auto_scale_enabled
  maximum_elastic_worker_count    = var.spec.maximum_elastic_worker_count

  # Flipping zone balancing on with fewer than 2 workers forces the plan
  # to be recreated -- keep worker_count a multiple of the region's zone
  # count (typically 3).
  zone_balancing_enabled = var.spec.zone_balancing_enabled

  per_site_scaling_enabled = var.spec.per_site_scaling_enabled

  tags = local.final_tags
}

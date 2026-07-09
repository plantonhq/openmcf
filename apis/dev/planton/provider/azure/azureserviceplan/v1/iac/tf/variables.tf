variable "metadata" {
  description = "Metadata for the resource, including name and labels"
  type = object({
    name    = string,
    id      = optional(string),
    org     = optional(string),
    env     = optional(string),
    labels  = optional(map(string)),
    tags    = optional(list(string)),
    version = optional(object({ id = string, message = string }))
  })
}

variable "spec" {
  description = "Azure Service Plan specification"
  type = object({
    # The Azure region where the Service Plan will be created. ForceNew.
    region = string

    # The Azure Resource Group name. References are resolved to a literal
    # name by the platform before the module runs. ForceNew.
    resource_group = string

    # The name of the Service Plan (1-60 alphanumerics/hyphens/
    # underscores; unique within the resource group). ForceNew.
    service_plan_name = string

    # The operating system type, as the spec enum's name string (LINUX /
    # WINDOWS / WINDOWS_CONTAINER). Absent means LINUX. ForceNew.
    os_type = optional(string)

    # The SKU, as the spec enum's name string (e.g. PREMIUM_P1V3,
    # ELASTIC_PREMIUM_EP1, CONSUMPTION_Y1). One value picks the tier's
    # capabilities and the VM size. NOT ForceNew -- plans re-tier in
    # place.
    sku_name = string

    # The ARM ID of the App Service Environment v3 hosting the plan.
    # Only Isolated SKUs may set this (spec-enforced).
    app_service_environment_id = optional(string)

    # Number of VM instances. Unset lets Azure apply the SKU's default
    # capacity. Consumption/Flex/Elastic tiers manage this themselves.
    worker_count = optional(number)

    # Automatic HTTP-load scaling for Premium plans (spec gates it to
    # Premium SKUs).
    premium_plan_auto_scale_enabled = optional(bool, false)

    # The scale-out ceiling for Elastic Premium / Workflow plans (and
    # Premium plans with auto-scale) -- the serverless cost-control
    # lever.
    maximum_elastic_worker_count = optional(number)

    # Spread instances across availability zones. The spec gates this to
    # the tiers Azure supports (Premium and above).
    zone_balancing_enabled = optional(bool, false)

    # Let individual apps scale independently of the plan's instance
    # count.
    per_site_scaling_enabled = optional(bool, false)

    # Free-form user tags, merged over the metadata-derived tags (user
    # tags win on key collision).
    tags = optional(map(string), {})
  })
}

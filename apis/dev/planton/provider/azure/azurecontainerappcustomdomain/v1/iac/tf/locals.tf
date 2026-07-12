locals {
  # The binding-type enum arrives as the proto value name; Azure's wire
  # vocabulary is mixed-case and the provider validates it
  # case-sensitively.
  certificate_binding_type_map = {
    "SNI_ENABLED" = "SniEnabled"
    "DISABLED"    = "Disabled"
    "AUTO"        = "Auto"
  }

  # The two lifecycle flows, dispatched below: bring-your-own certificate
  # (certificate id present) vs Azure-managed certificate (both
  # certificate fields absent). Spec validation guarantees the
  # certificate id and binding type travel together.
  is_byo_certificate = var.spec.container_app_environment_certificate_id != ""

  # No tag map: Azure models the binding as an entry in the app's ingress
  # configuration, not a taggable ARM resource.
}

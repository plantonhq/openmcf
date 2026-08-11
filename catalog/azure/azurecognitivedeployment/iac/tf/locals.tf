# The deployment is an ARM child of its account: the provider's schema
# carries no location, resource group, or tags for it (ARM derives all
# three through the account), so this module derives no tag map.
locals {
  # The spec's SKU tiers to the provider's wire values. Unspecified
  # (the enum's zero value renders as "") maps to null so ARM derives
  # the tier from the SKU name.
  sku_tier_wire = {
    "FREE"       = "Free"
    "BASIC"      = "Basic"
    "STANDARD"   = "Standard"
    "PREMIUM"    = "Premium"
    "ENTERPRISE" = "Enterprise"
  }

  # The spec's version-upgrade options to the provider's wire values.
  # Unspecified maps to null and the provider applies its default,
  # "OnceNewDefaultVersionAvailable".
  version_upgrade_option_wire = {
    "ONCE_CURRENT_VERSION_EXPIRED"       = "OnceCurrentVersionExpired"
    "ONCE_NEW_DEFAULT_VERSION_AVAILABLE" = "OnceNewDefaultVersionAvailable"
    "NO_AUTO_UPGRADE"                    = "NoAutoUpgrade"
  }
}

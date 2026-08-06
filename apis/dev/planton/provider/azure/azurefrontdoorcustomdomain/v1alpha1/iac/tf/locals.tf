locals {
  # Enums arrive as the spec enum's FULL value names; ARM wants its own
  # casing. certificate_type absent means ManagedCertificate (Azure's
  # default), materialized in main.tf.
  certificate_type_map = {
    "MANAGED_CERTIFICATE"  = "ManagedCertificate"
    "CUSTOMER_CERTIFICATE" = "CustomerCertificate"
  }

  # The year-versioned names are ARM's own vocabulary; only CUSTOMIZED
  # needs a casing change.
  cipher_suite_set_type_map = {
    "TLS12_2022" = "TLS12_2022"
    "TLS12_2023" = "TLS12_2023"
    "CUSTOMIZED" = "Customized"
  }

  # No Azure tags: ARM does not support tags on Front Door custom
  # domains, so the platform's identity tags live on the profile.
}

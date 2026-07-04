# A Key Vault certificate is a data-plane object: the provider talks to the
# vault's {name}.vault.azure.net endpoint, not ARM -- which is why creation
# fails with a 403 when the deploying credential lacks data-plane
# certificate permissions on the vault, even if it owns the subscription.
#
# Enrollment against issuer "Self" completes synchronously; a CA issuer
# keeps the operation pending until the CA responds, so expect longer
# creates on integrated-CA policies.
resource "azurerm_key_vault_certificate" "main" {
  name         = var.spec.name
  key_vault_id = var.spec.key_vault_id

  # Import path: bring an existing PFX/PEM bundle (the contents carry the
  # private key). Changing the contents imports a NEW VERSION of the
  # certificate rather than replacing the object.
  dynamic "certificate" {
    for_each = var.spec.certificate != null ? [var.spec.certificate] : []
    content {
      contents = certificate.value.contents
      password = certificate.value.password
    }
  }

  # Generate path (and governance for imports that carry an explicit
  # policy). Everything except lifetime_action creates a new certificate
  # version when changed; lifetime_action updates in place.
  dynamic "certificate_policy" {
    for_each = var.spec.certificate_policy != null ? [var.spec.certificate_policy] : []
    content {
      issuer_parameters {
        name = certificate_policy.value.issuer_name
      }

      # The private key's shape. key_size stays null for EC keys so Azure
      # derives it from the curve -- identical behavior on both engines.
      key_properties {
        exportable = certificate_policy.value.key_properties.exportable
        key_type   = local.key_type_map[certificate_policy.value.key_properties.key_type]
        key_size   = certificate_policy.value.key_properties.key_size
        curve      = certificate_policy.value.key_properties.curve != null ? local.curve_map[certificate_policy.value.key_properties.curve] : null
        reuse_key  = certificate_policy.value.key_properties.reuse_key
      }

      # Renewal/notification actions as expiry approaches. Exactly one
      # trigger field per action (spec validation enforces Azure's
      # contract).
      dynamic "lifetime_action" {
        for_each = certificate_policy.value.lifetime_actions
        content {
          action {
            action_type = local.action_type_map[lifetime_action.value.action_type]
          }
          trigger {
            days_before_expiry  = lifetime_action.value.trigger.days_before_expiry
            lifetime_percentage = lifetime_action.value.trigger.lifetime_percentage
          }
        }
      }

      secret_properties {
        content_type = local.content_type_map[certificate_policy.value.secret_properties.content_type]
      }

      # X.509 content -- present only when the vault generates the
      # certificate (spec validation requires it then; imports derive it
      # from the bundle).
      dynamic "x509_certificate_properties" {
        for_each = certificate_policy.value.x509_certificate_properties != null ? [certificate_policy.value.x509_certificate_properties] : []
        content {
          subject            = x509_certificate_properties.value.subject
          key_usage          = [for u in x509_certificate_properties.value.key_usage : local.key_usage_map[u]]
          extended_key_usage = x509_certificate_properties.value.extended_key_usage
          validity_in_months = x509_certificate_properties.value.validity_in_months

          dynamic "subject_alternative_names" {
            for_each = x509_certificate_properties.value.subject_alternative_names != null ? [x509_certificate_properties.value.subject_alternative_names] : []
            content {
              dns_names = subject_alternative_names.value.dns_names
              emails    = subject_alternative_names.value.emails
              upns      = subject_alternative_names.value.upns
            }
          }
        }
      }
    }
  }

  tags = local.final_tags
}

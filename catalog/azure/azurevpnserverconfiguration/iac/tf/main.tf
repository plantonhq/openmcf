# Create the VPN server configuration -- the reusable point-to-site
# authentication policy (Entra ID / certificate / RADIUS, trusted
# certificates, tunnel protocols) point-to-site VPN gateways attach
# to. The object is free, deploys in seconds, and gateways pick up
# in-place changes without redeploying.
#
# The spec's CEL contracts already guarantee each enabled
# authentication type brings its block (AAD -> aad_authentication,
# Certificate -> client_root_certificates, Radius -> radius) -- the
# provider enforces the same three rules at apply time.
resource "azurerm_vpn_server_configuration" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # The wire values are the spec's own vocabulary ("AAD",
  # "Certificate", "Radius") -- no mapping needed.
  vpn_authentication_types = var.spec.vpn_authentication_types

  dynamic "azure_active_directory_authentication" {
    for_each = var.spec.aad_authentication != null ? [var.spec.aad_authentication] : []
    content {
      audience = azure_active_directory_authentication.value.audience
      issuer   = azure_active_directory_authentication.value.issuer
      tenant   = azure_active_directory_authentication.value.tenant
    }
  }

  dynamic "client_root_certificate" {
    for_each = var.spec.client_root_certificates
    content {
      name             = client_root_certificate.value.name
      public_cert_data = client_root_certificate.value.public_cert_data
    }
  }

  dynamic "client_revoked_certificate" {
    for_each = var.spec.client_revoked_certificates
    content {
      name       = client_revoked_certificate.value.name
      thumbprint = client_revoked_certificate.value.thumbprint
    }
  }

  # The spec requires every field of a configured proposal (no partial
  # pinning); the vocabularies are already wire values.
  dynamic "ipsec_policy" {
    for_each = var.spec.ipsec_policy != null ? [var.spec.ipsec_policy] : []
    content {
      dh_group               = ipsec_policy.value.dh_group
      ike_encryption         = ipsec_policy.value.ike_encryption
      ike_integrity          = ipsec_policy.value.ike_integrity
      ipsec_encryption       = ipsec_policy.value.ipsec_encryption
      ipsec_integrity        = ipsec_policy.value.ipsec_integrity
      pfs_group              = ipsec_policy.value.pfs_group
      sa_lifetime_seconds    = ipsec_policy.value.sa_lifetime_seconds
      sa_data_size_kilobytes = ipsec_policy.value.sa_data_size_kilobytes
    }
  }

  dynamic "radius" {
    for_each = var.spec.radius != null ? [var.spec.radius] : []
    content {
      # Sensitive: the secret never appears in plan output, and ARM
      # never returns it on reads (the import round-trip declares the
      # matching tolerance).
      dynamic "server" {
        for_each = radius.value.servers
        content {
          address = server.value.address
          secret  = server.value.secret
          score   = server.value.score
        }
      }

      dynamic "client_root_certificate" {
        for_each = radius.value.client_root_certificates
        content {
          name       = client_root_certificate.value.name
          thumbprint = client_root_certificate.value.thumbprint
        }
      }

      dynamic "server_root_certificate" {
        for_each = radius.value.server_root_certificates
        content {
          name             = server_root_certificate.value.name
          public_cert_data = server_root_certificate.value.public_cert_data
        }
      }
    }
  }

  # Optional+Computed on the provider: emit null when the spec leaves
  # it empty so ARM's default selection applies and reads don't drift.
  vpn_protocols = length(var.spec.vpn_protocols) > 0 ? var.spec.vpn_protocols : null

  tags = local.final_tags
}

# The composed policy groups: standalone ARM children of the
# configuration, one per spec entry, keyed by name (named
# member-matching rules a point-to-site gateway maps to address
# pools). The policy_group_ids output publishes each group's ARM id.
resource "azurerm_vpn_server_configuration_policy_group" "policy_groups" {
  for_each = { for policy_group in var.spec.policy_groups : policy_group.name => policy_group }

  name                        = each.value.name
  vpn_server_configuration_id = azurerm_vpn_server_configuration.main.id

  # Both ForceNew on the group: is_default marks the catch-all group.
  is_default = each.value.is_default
  priority   = each.value.priority

  dynamic "policy" {
    for_each = each.value.policies
    content {
      name  = policy.value.name
      type  = policy.value.type
      value = policy.value.value
    }
  }
}

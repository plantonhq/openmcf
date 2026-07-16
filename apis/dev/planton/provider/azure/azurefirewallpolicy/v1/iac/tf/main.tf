# Create the firewall policy -- the reusable rule-and-inspection document
# Azure Firewall instances enforce. The policy carries WHAT the firewall
# does (threat intelligence, DNS proxying, TLS inspection, IDPS, SNAT
# posture); rules live in separate rule collection group resources; the
# firewall instance carries WHERE enforcement runs.
#
# The sku and threat-intelligence mode are always sent explicitly
# (Standard/Alert when unspecified) -- Azure's own defaults, made
# deterministic so both engines produce identical payloads. The sku is
# ForceNew: changing the tier replaces the policy, which is why the tier
# must be chosen deliberately up front (Premium features -- IDPS, TLS
# inspection -- are gated to PREMIUM in the spec validation).
resource "azurerm_firewall_policy" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  sku                      = local.sku
  threat_intelligence_mode = local.threat_intelligence_mode

  # Inheritance: the base policy's rules and settings apply beneath this
  # policy's own. ARM validates region/tier pairing at apply time.
  base_policy_id = (
    var.spec.base_policy_id != null && var.spec.base_policy_id != ""
  ) ? var.spec.base_policy_id : null

  sql_redirect_allowed = var.spec.sql_redirect_allowed

  # SNAT ranges: sent only when the user overrides Azure's IANA-private
  # default.
  private_ip_ranges = (
    length(var.spec.private_ip_ranges) > 0
  ) ? var.spec.private_ip_ranges : null

  # Azure only records "Enabled" for auto-learn -- disabling is done by
  # omission on the wire, so the flag is sent only when true (an explicit
  # false would still read back as absent and churn state).
  auto_learn_private_ranges_enabled = (
    var.spec.auto_learn_private_ranges_enabled
  ) ? true : null

  dynamic "threat_intelligence_allowlist" {
    for_each = var.spec.threat_intelligence_allowlist != null ? [var.spec.threat_intelligence_allowlist] : []
    content {
      ip_addresses = threat_intelligence_allowlist.value.ip_addresses
      fqdns        = threat_intelligence_allowlist.value.fqdns
    }
  }

  dynamic "dns" {
    for_each = var.spec.dns != null ? [var.spec.dns] : []
    content {
      servers       = dns.value.servers
      proxy_enabled = dns.value.proxy_enabled
    }
  }

  # IDPS is Premium-only (the spec validation front-loads the gate). The
  # engine mode is sent only when specified -- Azure defaults an
  # unspecified mode to Off.
  dynamic "intrusion_detection" {
    for_each = var.spec.intrusion_detection != null ? [var.spec.intrusion_detection] : []
    content {
      mode = (
        intrusion_detection.value.mode != null
      ) ? lookup(local.idps_state_wire, intrusion_detection.value.mode, null) : null

      private_ranges = intrusion_detection.value.private_ranges

      dynamic "signature_overrides" {
        for_each = intrusion_detection.value.signature_overrides
        content {
          id = signature_overrides.value.id
          state = (
            signature_overrides.value.state != null
          ) ? lookup(local.idps_state_wire, signature_overrides.value.state, null) : null
        }
      }

      dynamic "traffic_bypass" {
        for_each = intrusion_detection.value.traffic_bypass
        content {
          name        = traffic_bypass.value.name
          description = traffic_bypass.value.description
          # The proto protocol names ARE the wire vocabulary
          # (ANY/ICMP/TCP/UDP); the provider validates case-insensitively.
          protocol              = traffic_bypass.value.protocol
          source_addresses      = traffic_bypass.value.source_addresses
          source_ip_groups      = traffic_bypass.value.source_ip_groups
          destination_addresses = traffic_bypass.value.destination_addresses
          destination_ip_groups = traffic_bypass.value.destination_ip_groups
          destination_ports     = traffic_bypass.value.destination_ports
        }
      }
    }
  }

  # TLS inspection reads the CA certificate from Key Vault through the
  # policy's user-assigned identity -- the identity block travels with
  # the certificate in practice.
  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = lookup(local.identity_type_wire, identity.value.type, identity.value.type)
      identity_ids = identity.value.user_assigned_identity_ids
    }
  }

  dynamic "tls_certificate" {
    for_each = var.spec.tls_certificate != null ? [var.spec.tls_certificate] : []
    content {
      key_vault_secret_id = tls_certificate.value.key_vault_secret_id
      name                = tls_certificate.value.name
    }
  }

  dynamic "insights" {
    for_each = var.spec.insights != null ? [var.spec.insights] : []
    content {
      enabled                            = insights.value.enabled
      default_log_analytics_workspace_id = insights.value.default_log_analytics_workspace_id
      retention_in_days                  = insights.value.retention_in_days

      dynamic "log_analytics_workspace" {
        for_each = insights.value.log_analytics_workspaces
        content {
          id                = log_analytics_workspace.value.workspace_id
          firewall_location = log_analytics_workspace.value.firewall_location
        }
      }
    }
  }

  dynamic "explicit_proxy" {
    for_each = var.spec.explicit_proxy != null ? [var.spec.explicit_proxy] : []
    content {
      enabled         = explicit_proxy.value.enabled
      http_port       = explicit_proxy.value.http_port
      https_port      = explicit_proxy.value.https_port
      enable_pac_file = explicit_proxy.value.enable_pac_file
      pac_file_port   = explicit_proxy.value.pac_file_port
      pac_file        = explicit_proxy.value.pac_file
    }
  }

  tags = local.final_tags
}

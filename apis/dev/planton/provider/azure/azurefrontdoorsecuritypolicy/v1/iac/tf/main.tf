# The security policy -- the association that attaches a Front Door WAF
# policy to the hostnames a profile serves. The WAF enforces nothing
# until this exists; this resource IS the enforcement seam.
#
# The provider's security_policies -> firewall -> association nesting is
# a one-choice ARM union (WebApplicationFirewall is the only
# security-policy type); the spec models it flat and this block is where
# the wrapper shape gets rebuilt.
resource "azurerm_cdn_frontdoor_security_policy" "main" {
  name                     = var.spec.security_policy_name
  cdn_frontdoor_profile_id = var.spec.profile_id

  security_policies {
    firewall {
      cdn_frontdoor_firewall_policy_id = var.spec.firewall_policy_id

      association {
        # Endpoint ids protect the generated *.azurefd.net hostname;
        # custom-domain ids protect that custom hostname. Azure accepts
        # both types interchangeably.
        dynamic "domain" {
          for_each = local.domain_ids
          content {
            cdn_frontdoor_domain_id = domain.value
          }
        }

        # The service accepts exactly one pattern ("/*") -- a constant,
        # not configuration; scope enforcement by choosing WHICH domains
        # to associate. NOTE the engine dialect: the pulumi bridge
        # flattens this one-item list to a single string ("/*") -- same
        # ARM payload.
        patterns_to_match = ["/*"]
      }
    }
  }
}

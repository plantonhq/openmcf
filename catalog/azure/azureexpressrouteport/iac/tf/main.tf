# Create the ExpressRoute Port -- the physical port pair on a Microsoft
# edge router at the peering location (ExpressRoute Direct). The port
# bills its full monthly rate FROM CREATION, whether or not any
# cross-connect exists, and some subscriptions need Microsoft enrollment
# for ExpressRoute Direct before ARM accepts this create.
resource "azurerm_express_route_port" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # The physical facts of the port pair -- all fixed at creation.
  peering_location  = var.spec.peering_location
  bandwidth_in_gbps = var.spec.bandwidth_in_gbps
  encapsulation     = lookup(local.encapsulation_wire, var.spec.encapsulation, var.spec.encapsulation)

  billing_type = local.billing_type

  # The port's managed identity -- what MACsec uses to read the CAK/CKN
  # secrets from Key Vault (the spec's contracts guarantee a
  # user-assigned identity is present whenever MACsec keys are set).
  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type
      identity_ids = identity.value.identity_ids
    }
  }

  # ARM always creates the link pair with the port; these blocks only
  # MANIPULATE the existing links (admin state, MACsec). The provider
  # applies link configuration in a second call after the port exists,
  # so a fresh port briefly reports links disabled even when
  # admin_enabled is true. Unset MACsec cipher applies ARM's default
  # (GcmAes128); the CKN/CAK keys travel together (spec-validated).
  dynamic "link1" {
    for_each = var.spec.link1 != null ? [var.spec.link1] : []
    content {
      admin_enabled = link1.value.admin_enabled
      macsec_cipher = (
        link1.value.macsec_cipher == null
        ? "GcmAes128"
        : lookup(local.macsec_cipher_wire, link1.value.macsec_cipher, link1.value.macsec_cipher)
      )
      macsec_ckn_keyvault_secret_id = link1.value.macsec_ckn_keyvault_secret_id != "" ? link1.value.macsec_ckn_keyvault_secret_id : null
      macsec_cak_keyvault_secret_id = link1.value.macsec_cak_keyvault_secret_id != "" ? link1.value.macsec_cak_keyvault_secret_id : null
      macsec_sci_enabled            = link1.value.macsec_sci_enabled
    }
  }

  dynamic "link2" {
    for_each = var.spec.link2 != null ? [var.spec.link2] : []
    content {
      admin_enabled = link2.value.admin_enabled
      macsec_cipher = (
        link2.value.macsec_cipher == null
        ? "GcmAes128"
        : lookup(local.macsec_cipher_wire, link2.value.macsec_cipher, link2.value.macsec_cipher)
      )
      macsec_ckn_keyvault_secret_id = link2.value.macsec_ckn_keyvault_secret_id != "" ? link2.value.macsec_ckn_keyvault_secret_id : null
      macsec_cak_keyvault_secret_id = link2.value.macsec_cak_keyvault_secret_id != "" ? link2.value.macsec_cak_keyvault_secret_id : null
      macsec_sci_enabled            = link2.value.macsec_sci_enabled
    }
  }

  tags = local.final_tags
}

# The composed authorizations: standalone ARM children of the port, one
# per spec entry, keyed by name (the provider serializes them against
# the port with its own lock -- ARM allows one port mutation at a time).
# ARM GENERATES each key; the name-keyed authorization_keys output
# surfaces them (sensitive) so a circuit in another subscription can be
# built on this port's capacity.
resource "azurerm_express_route_port_authorization" "authorizations" {
  for_each = { for authorization in var.spec.authorizations : authorization.name => authorization }

  name                    = each.value.name
  express_route_port_name = azurerm_express_route_port.main.name
  resource_group_name     = var.spec.resource_group
}

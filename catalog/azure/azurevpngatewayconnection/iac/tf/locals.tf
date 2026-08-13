# The connection carries no tags and no resource group of its own: ARM
# addresses it as a child of the VPN gateway, and the provider's schema
# has no tags argument -- so this module needs none of the family's
# usual tag-merging locals.
locals {
  # The spec's enum NAMES mapped onto ARM's vocabulary. Unset applies
  # ARM's defaults ("IKEv2"/"Default") explicitly so the rendered plan
  # shows the real values -- mirroring the Pulumi module's nil handling.
  protocol_wire = {
    "IKE_V1" = "IKEv1"
    "IKE_V2" = "IKEv2"
  }

  connection_mode_wire = {
    "DEFAULT"        = "Default"
    "INITIATOR_ONLY" = "InitiatorOnly"
    "RESPONDER_ONLY" = "ResponderOnly"
  }
}

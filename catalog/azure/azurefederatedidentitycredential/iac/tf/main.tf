# Create the federated identity credential -- a keyless OIDC trust rule on
# the referenced user-assigned managed identity.
#
# Lifecycle notes worth knowing before operating this resource:
# - issuer, subject, and audience update IN PLACE; name and the parent
#   identity are the credential's ARM identity, so changing either replaces
#   it (delete + create).
# - The provider serializes writes per parent identity (ARM rejects
#   concurrent credential writes on one identity), so several credentials on
#   the same identity apply sequentially -- expected, not a hang.
# - The resource group is derived from the parent identity's ARM ID (azurerm
#   does the same internally), so the module never asks the user to restate
#   derivable state that could then disagree with the parent.
resource "azurerm_federated_identity_credential" "main" {
  name = var.spec.name

  # The v4-canonical parent argument (the deprecated parent_id alias is
  # intentionally not used). The identity's ARM ID both parents the
  # credential and locates its resource group.
  user_assigned_identity_id = var.spec.user_assigned_identity

  # The three-way trust match Azure AD evaluates on every token exchange:
  # the incoming token's iss, sub, and aud claims must equal these values
  # exactly.
  issuer  = var.spec.issuer
  subject = var.spec.subject

  # ARM models the audience as a single-element list (exactly one audience
  # per credential today); the schema caps the list at one, so the module
  # passes the one configured (or defaulted) value.
  audience = [local.audience]
}

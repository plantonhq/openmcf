# Create the Machine Learning compute instance -- a single always-on VM
# serving as one data scientist's cloud workstation, as an ARM child of
# its workspace (.../workspaces/{ws}/computes/{name}).
#
# The provider has NO update path for this resource: EVERY argument is
# ForceNew, tags included -- any change replaces the instance (its OS
# disk and local files go with it). The instance always runs in its
# workspace's region (the service's own rule; there is no location
# argument), and its name is reserved region-wide per subscription.
#
# One contract lives at apply time, not manifest time: when
# node_public_ip_enabled is false, the provider requires subnet_resource_id
# UNLESS the workspace runs a managed network -- it depends on the live
# workspace's isolation mode (recorded on the spec fields).
resource "azurerm_machine_learning_compute_instance" "main" {
  name                          = var.spec.name
  machine_learning_workspace_id = var.spec.workspace_id
  virtual_machine_size          = var.spec.virtual_machine_size

  # "personal" is the only value the provider accepts; unset omits the
  # property and leaves the service default.
  authorization_type = var.spec.authorization_type != "" ? var.spec.authorization_type : null

  # The admin-provisions-for-the-team pattern: assign the instance to a
  # user other than the deploying principal.
  dynamic "assign_to_user" {
    for_each = var.spec.assign_to_user != null ? [var.spec.assign_to_user] : []
    content {
      tenant_id = assign_to_user.value.tenant_id != "" ? assign_to_user.value.tenant_id : null
      object_id = assign_to_user.value.object_id != "" ? assign_to_user.value.object_id : null
    }
  }

  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_wire[identity.value.type]
      identity_ids = length(identity.value.identity_ids) > 0 ? identity.value.identity_ids : null
    }
  }

  # Optional-with-default-true on the provider: emit null when the spec
  # leaves them unset so the provider default applies.
  local_auth_enabled     = var.spec.local_auth_enabled
  node_public_ip_enabled = var.spec.node_public_ip_enabled

  # Absent block means the SSH port is DISABLED (the provider's own
  # contract); the service assigns the username and port, surfaced as
  # outputs.
  dynamic "ssh" {
    for_each = var.spec.ssh != null ? [var.spec.ssh] : []
    content {
      public_key = ssh.value.public_key
    }
  }

  # Only legal when the workspace does NOT use a managed network.
  subnet_resource_id = var.spec.subnet_id != "" ? var.spec.subnet_id : null

  description = var.spec.description != "" ? var.spec.description : null

  tags = local.final_tags
}

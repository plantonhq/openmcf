# Create the Machine Learning compute cluster -- the auto-scaling pool
# of VMs that training jobs and pipelines run on, as an ARM child of
# its workspace (.../workspaces/{ws}/computes/{name}).
#
# Only identity, scale_settings and tags update in place -- every other
# argument is ForceNew (the provider's own contract). Uniquely in the
# ML family, `location` here is the NODES' region and may differ from
# the workspace's; the provider writes the cluster envelope at the
# WORKSPACE's region, so ARM reads the envelope back there (recorded on
# the spec's region field).
resource "azurerm_machine_learning_compute_cluster" "main" {
  name                          = var.spec.name
  machine_learning_workspace_id = var.spec.workspace_id
  location                      = var.spec.region
  vm_size                       = var.spec.vm_size

  # Enum name -> wire value; the spec enum is required, so there is no
  # unspecified fallback.
  vm_priority = local.vm_priority_wire[var.spec.vm_priority]

  # The one substantive in-place-updatable setting.
  scale_settings {
    min_node_count                       = var.spec.scale_settings.min_node_count
    max_node_count                       = var.spec.scale_settings.max_node_count
    scale_down_nodes_after_idle_duration = var.spec.scale_settings.scale_down_nodes_after_idle_duration
  }

  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_wire[identity.value.type]
      identity_ids = length(identity.value.identity_ids) > 0 ? identity.value.identity_ids : null
    }
  }

  # The admin account created on every node. At least one credential is
  # set (spec CEL mirrors the provider's AtLeastOneOf). The password is
  # sensitive -- resolved from secret references, masked in plan output.
  dynamic "ssh" {
    for_each = var.spec.ssh != null ? [var.spec.ssh] : []
    content {
      admin_username = ssh.value.admin_username
      admin_password = ssh.value.admin_password != "" ? ssh.value.admin_password : null
      key_value      = ssh.value.key_value != "" ? ssh.value.key_value : null
    }
  }

  # Plain bool: false is the provider's own default, so passing the
  # zero value through is exact.
  ssh_public_access_enabled = var.spec.ssh_public_access_enabled

  # Optional-with-default-true on the provider: emit null when the spec
  # leaves them unset so the provider default applies.
  local_auth_enabled     = var.spec.local_auth_enabled
  node_public_ip_enabled = var.spec.node_public_ip_enabled

  # Optional+Computed on the provider: unset lets Azure network the
  # nodes (a workspace managed network assigns one, read back).
  subnet_resource_id = var.spec.subnet_id != "" ? var.spec.subnet_id : null

  description = var.spec.description != "" ? var.spec.description : null

  tags = local.final_tags
}

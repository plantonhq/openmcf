# The Dapr component: a pluggable backend (state store, pub/sub broker,
# secret store, binding) registered once on the environment and consumed
# by Dapr-enabled apps whose dapr.app_id appears in scopes. Name, type,
# and environment are ForceNew.
resource "azurerm_container_app_environment_dapr_component" "main" {
  name                         = var.spec.component_name
  container_app_environment_id = var.spec.container_app_environment_id
  component_type               = var.spec.component_type
  version                      = var.spec.version

  # Documented default resolved in variables.tf; left false so a broken
  # component fails loudly at sidecar startup instead of surfacing as
  # runtime errors on first use.
  init_timeout  = var.spec.init_timeout
  ignore_errors = var.spec.ignore_errors

  # Connection strings and keys travel as component secrets referenced
  # from metadata by secret_name -- never as plain metadata values.
  dynamic "secret" {
    for_each = var.spec.secrets
    content {
      name  = secret.value.name
      value = secret.value.value
    }
  }

  # The component's configuration entries; keys depend on the component
  # type. The spec's CEL guarantees value XOR secret_name per entry.
  dynamic "metadata" {
    for_each = var.spec.metadata
    content {
      name        = metadata.value.name
      value       = metadata.value.value
      secret_name = metadata.value.secret_name
    }
  }

  # An empty scopes list exposes the component to every Dapr-enabled app
  # in the environment -- scope production components deliberately.
  scopes = length(var.spec.scopes) > 0 ? var.spec.scopes : null
}

locals {
  # Optional MySQL auth plugin: null when unset so the provider's own
  # default (caching_sha2_password) applies without a diff -- the provider
  # suppresses the default-vs-unset diff itself.
  mysql_auth_plugin = try(var.spec.mysql_auth_plugin, "") != "" ? var.spec.mysql_auth_plugin : null
}

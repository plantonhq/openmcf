# Local values and computed configuration

locals {
  # Build combined labels
  standard_labels = {
    "managed-by"    = "planton"
    "resource"      = var.metadata.name
    "resource-kind" = "KubernetesSecret"
  }

  labels = merge(local.standard_labels, var.spec.labels)

  # The service-account-token type mandates an annotation that binds the token
  # to its ServiceAccount; the cluster's token controller acts on it.
  service_account_annotation = var.spec.service_account_token != null ? {
    "kubernetes.io/service-account.name" = var.spec.service_account_token.service_account_name
  } : {}

  # Build annotations
  annotations = merge(var.spec.annotations, local.service_account_annotation)

  # Determine secret type from which variant is set
  secret_type = (
    var.spec.opaque != null ? "Opaque" :
    var.spec.tls != null ? "kubernetes.io/tls" :
    var.spec.docker_config_json != null ? "kubernetes.io/dockerconfigjson" :
    var.spec.basic_auth != null ? "kubernetes.io/basic-auth" :
    var.spec.ssh_auth != null ? "kubernetes.io/ssh-auth" :
    var.spec.service_account_token != null ? "kubernetes.io/service-account-token" :
    "Opaque"
  )

  # Build the docker config JSON structure when docker_config_json variant is set
  docker_config_json = var.spec.docker_config_json != null ? jsonencode({
    auths = {
      (var.spec.docker_config_json.registry_server) = {
        username = var.spec.docker_config_json.username
        password = var.spec.docker_config_json.password
        email    = var.spec.docker_config_json.email
        auth     = base64encode("${var.spec.docker_config_json.username}:${var.spec.docker_config_json.password}")
      }
    }
  }) : null

  # Compute secret data map based on which variant is set.
  # The data attribute takes plain strings (the provider base64-encodes them).
  # service-account-token intentionally falls through to {}: the cluster's token
  # controller populates token/ca.crt/namespace after creation.
  secret_data = (
    var.spec.opaque != null ? var.spec.opaque.data :
    var.spec.tls != null ? {
      "tls.crt" = var.spec.tls.tls_crt
      "tls.key" = var.spec.tls.tls_key
    } :
    var.spec.docker_config_json != null ? {
      ".dockerconfigjson" = local.docker_config_json
    } :
    var.spec.basic_auth != null ? {
      "username" = var.spec.basic_auth.username
      "password" = var.spec.basic_auth.password
    } :
    var.spec.ssh_auth != null ? {
      "ssh-privatekey" = var.spec.ssh_auth.ssh_private_key
    } :
    {}
  )

  # Binary entries arrive base64-encoded and flow to the binary_data attribute,
  # which expects base64 (unlike data, which takes plain strings).
  # try() is used for safe nested access because HCL does not short-circuit &&.
  secret_binary_data = try(var.spec.opaque.binary_data, {})
}

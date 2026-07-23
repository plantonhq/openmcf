# Workload-scoped satellite Secrets.
#
# env secret: literal secret env values collected across the app container,
# every sidecar, and every init container, materialized into ONE Secret
# (locals.env_secret_data). Secrets referenced from existing Kubernetes
# Secrets are wired directly as env references and never pass through here.
resource "kubernetes_secret_v1" "env_secrets" {
  count = length(local.env_secret_data) > 0 ? 1 : 0

  metadata {
    name      = local.env_secret_name
    namespace = local.namespace
    labels    = local.final_labels
  }

  type = "Opaque"
  # The provider's `data` argument takes plaintext and base64-encodes it on the
  # wire (there is no string_data argument on kubernetes_secret_v1).
  data = local.env_secret_data

  depends_on = [kubernetes_namespace.this]
}

# Docker-registry credential for private image pulls, materialized when the
# platform injects a dockerconfigjson at deploy time. Pods reference it through
# the pod-level image pull secret list (see locals.image_pull_secret_names).
resource "kubernetes_secret_v1" "image_pull" {
  count = local.create_image_pull_secret ? 1 : 0

  metadata {
    name      = local.image_pull_secret_name
    namespace = local.namespace
    labels    = local.final_labels
  }

  type = "kubernetes.io/dockerconfigjson"
  data = {
    ".dockerconfigjson" = var.docker_config_json
  }

  depends_on = [kubernetes_namespace.this]
}

# Module-owned Secret for declared pre-shared API keys:
# `<metadata.name>-authn-keys`, data key `keys` = the comma-joined key
# list (the chart's keysSecret contract — the Deployment reads
# OPENFGA_AUTHN_PRESHARED_KEYS from exactly that key). This Secret is
# the ONLY place the key material lands: chart values carry just its
# NAME (authn.preshared.keysSecret); the chart's plain-values key list,
# which would render every key into the Deployment manifest, is
# deliberately never used.
#
# Not created on the existing-Secret arm (the user owns that one) or
# when authn is unset/oidc. Created in the release namespace before the
# release — a secretKeyRef can only read Secrets in the workload's own
# namespace. Pulumi twin: authn_secret.go.

resource "kubernetes_secret_v1" "authn_keys" {
  count = local.materialize_authn_secret ? 1 : 0

  metadata {
    name      = local.authn_keys_secret_name
    namespace = local.namespace
    labels    = local.labels
  }

  type = "Opaque"

  # The provider's `data` takes plaintext and base64-encodes it on write
  # (the resource has no string_data argument; Pulumi twin: StringData).
  data = {
    keys = join(",", var.spec.authn.preshared.keys)
  }

  depends_on = [
    kubernetes_namespace_v1.openfga,
  ]
}

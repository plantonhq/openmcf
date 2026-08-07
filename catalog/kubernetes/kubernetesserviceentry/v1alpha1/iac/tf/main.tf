# Applies the ServiceEntry custom resource through kubectl_manifest (alekc/kubectl):
# no plan-time cluster dependency (plannable before the CRDs exist), applied
# server-side. No wait, deliberately: the CR is configuration its controller
# consumes; applying it server-side-validated is the whole contract. Pulumi
# equivalent: the typed CR without await annotations.
resource "kubectl_manifest" "service_entry" {
  yaml_body = yamlencode({
    apiVersion = "networking.istio.io/v1"
    kind       = "ServiceEntry"
    metadata = {
      name      = var.metadata.name
      namespace = var.spec.namespace
      labels    = local.labels
    }
    spec = local.manifest_spec
  })

  server_side_apply = true
}

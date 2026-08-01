# The k8s.keycloak.org/v2beta1 Keycloak CR — the single declaration the
# Keycloak Operator (a KubernetesKeycloakOperator install, the
# prerequisite) reconciles into the StatefulSet (named exactly after this
# resource), `<name>-service`, `<name>-discovery`, the NetworkPolicy and
# — unless the spec brings its own bootstrap-admin Secret — the
# create-once `<name>-initial-admin` credential Secret.
#
# The CR applies through kubectl_manifest (alekc/kubectl): unlike the
# hashicorp provider's kubernetes_manifest resource it needs no cluster
# connection at plan time — a Keycloak can be PLANNED before the
# operator's CRDs exist, which is what lets an infra chart deploy the
# operator and its servers in one run (and lets offline plan proofs work).
#
# No wait_for block, deliberately: server readiness depends on the
# operator (image pulls, database schema migrations, cluster formation)
# that is not part of applying the resource — the verifier owns
# readiness, the same never-block-on-a-controller posture as the sibling
# operator-CR modules.
resource "kubectl_manifest" "keycloak" {
  yaml_body = yamlencode({
    apiVersion = local.api_version
    kind       = "Keycloak"
    metadata = {
      name      = local.keycloak_name
      namespace = local.namespace
      labels    = local.labels
    }
    spec = local.keycloak_spec
  })

  server_side_apply = true
  force_conflicts   = true

  # BACKGROUND deletion, explicitly: the OPERATOR owns the Keycloak CR's
  # cascade — its finalizer tears down the StatefulSet, Services and
  # NetworkPolicy. Foreground propagation would block the delete on
  # children the operator keeps reconciling. Pulumi twin: the
  # pulumi.com/deletionPropagationPolicy annotation.
  delete_cascade = "Background"

  lifecycle {
    # FAIL LOUDLY on names past the operator's naming budget: every child
    # derives from this name by suffixing (`-network-policy` is the
    # longest at 15 characters) and StatefulSet pod hostnames must stay
    # DNS-legal — past 48 the derived names silently break the contract
    # the exported outputs are built on. Twin: the Pulumi module's
    # Resources() guard.
    precondition {
      condition     = length(var.metadata.name) <= 48
      error_message = "The Keycloak operator derives child names from the resource name (suffixes up to 15 characters) and pod hostnames must stay DNS-legal — use a name of at most 48 characters."
    }
  }

  depends_on = [
    kubernetes_namespace_v1.namespace,
  ]
}

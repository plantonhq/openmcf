# KubernetesNeo4j Terraform module.
#
# Installs Neo4j from the official Helm chart as a real Helm release. The
# typed spec renders into chart values (locals.helm_values); a declared
# admin password materializes as the "<name>-auth" Kubernetes Secret the
# chart consumes via neo4j.passwordFromSecret; the helm_values escape hatch
# is passed as a SECOND values document, which the provider merges over the
# first with Helm -f semantics — the exact semantic twin of the Pulumi
# module's buildHelmValues + mergeMaps.
#
# ORDERING IS LOAD-BEARING: the chart looks the passwordFromSecret Secret up
# AT TEMPLATE TIME and fails the install when it is missing, so the auth
# Secret is an explicit dependency of the release — it exists first, always.

# The optional installation namespace. Created before the release; deleted
# with the resource.
resource "kubernetes_namespace_v1" "neo4j" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# The declared admin password, materialized as an Opaque Secret with the
# chart's contract: ONE key, NEO4J_AUTH, whose value is "neo4j/<password>".
# Because neo4j.passwordFromSecret points here, the chart renders no auth
# Secret of its own — this Secret is the only place the credential lands,
# and it never transits chart values. Created only for the password arm
# (the existing_secret arm references a Secret the user owns).
resource "kubernetes_secret_v1" "auth" {
  count = local.create_auth_secret ? 1 : 0

  metadata {
    name      = local.auth_secret_name
    namespace = local.namespace
    labels    = local.labels
  }

  type = "Opaque"

  data = {
    NEO4J_AUTH = sensitive("neo4j/${local.auth_password}")
  }

  depends_on = [kubernetes_namespace_v1.neo4j]
}

resource "helm_release" "neo4j" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the server to become Ready — a database that never starts (bad
  # image, unschedulable pod, unbindable volume, a JVM that OOMs on boot)
  # should fail THIS apply, not the first driver connection. Neo4j
  # recovers/upgrades store files on startup, so the budget is generous.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  # Two documents, merged in order by the provider (helm -f semantics):
  # the typed rendering first, the user's escape hatch last.
  values = [
    yamlencode(local.helm_values),
    try(var.spec.helm_values, ""),
  ]

  depends_on = [
    kubernetes_namespace_v1.neo4j,
    kubernetes_secret_v1.auth,
  ]
}

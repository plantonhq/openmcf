# KubernetesKafkaUi Terraform module.
#
# Installs the kafbat UI console from the served kafka-ui Helm chart as a
# real Helm release. The typed spec renders into chart values
# (locals.helm_values): the cluster wiring becomes yamlApplicationConfig
# with ${ENV_VAR} placeholders in every password position, each placeholder
# wired to its source Secret through envs.secretMappings; the declared
# console login password materializes as the "<name>-secrets" Secret; the
# helm_values escape hatch is passed as a SECOND values document, which the
# provider merges over the first with Helm -f semantics — the exact
# semantic twin of the Pulumi module's buildHelmValues + mergeMaps.
#
# The release is named after metadata.name (NOT a fixed chart name):
# several consoles coexist in one cluster, each rendering its own
# Deployment and Service under the chart's fullname for that release.

# The optional installation namespace. Created before the release; deleted
# with the resource.
resource "kubernetes_namespace_v1" "kafka_ui" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# The LITERAL credentials the spec declares, materialized as an Opaque
# Secret — today that is exactly the console login password, under the
# key the KAFKA_UI_USER_PASSWORD secret mapping points at. Only the FIRST
# declared user's password is stored: LOGIN_FORM supports a single account
# (Spring Boot's default security user — see locals.tf), so keys for
# further users would be dead data no mapping references. Referenced
# credentials (sasl / schema registry / Connect password_secret entries)
# are NOT copied here — their mappings point at the source Secrets
# directly. This Secret is the only place the declared literal lands; it
# never transits chart values. Pulumi twin: secret.go's consoleSecret with
# the same name, key, and contents.
resource "kubernetes_secret_v1" "console" {
  count = local.auth_enabled ? 1 : 0

  metadata {
    name      = local.console_secret_name
    namespace = local.namespace
    labels    = local.labels
  }

  type = "Opaque"

  data = {
    (local.console_password_secret_key) = var.spec.auth.user.password
  }

  depends_on = [kubernetes_namespace_v1.kafka_ui]
}

resource "helm_release" "kafka_ui" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the console to become Ready — a console that never starts
  # (bad image, unschedulable pod, unresolvable cluster config) should
  # fail THIS apply, not the first browser hit.
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
    kubernetes_namespace_v1.kafka_ui,
    kubernetes_secret_v1.console,
  ]
}

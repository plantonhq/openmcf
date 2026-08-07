# KubernetesExternalDns Terraform module.
#
# Installs ExternalDNS from the official Helm chart as a real Helm release.
# The typed spec renders into chart values (locals.helm_values); declared
# provider credentials materialize as Kubernetes Secrets that the chart's
# env/volume wiring consumes; the helm_values escape hatch is passed as a
# SECOND values document, which the provider merges over the first with
# Helm -f semantics — the exact semantic twin of the Pulumi module's
# buildHelmValues + mergeMaps.
#
# The release is named after metadata.name (NOT a fixed chart name):
# multiple ExternalDNS instances per cluster — one per DNS provider / zone
# set, separated by TXT owner IDs — are a first-class upstream pattern.

# The optional installation namespace. Created before the release; deleted
# with the resource.
resource "kubernetes_namespace_v1" "external_dns" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# ---- credential secrets ----------------------------------------------------
# The selected provider's DECLARED credentials materialize as
# deterministically-named Secrets, so the credential itself never appears in
# chart values or pod specs — the chart wires env/volume references to them
# (locals.tf). Providers running keyless (or the webhook/in-memory arms)
# materialize nothing.

# Cloudflare: token consumed via CF_API_TOKEN (locals.tf env wiring).
resource "kubernetes_secret_v1" "cloudflare_credentials" {
  count = try(var.spec.cloudflare, null) != null ? 1 : 0

  metadata {
    name      = local.cloudflare_secret_name
    namespace = local.namespace
    labels    = local.labels
  }

  data = {
    "api-token" = var.spec.cloudflare.api_token
  }

  depends_on = [kubernetes_namespace_v1.external_dns]
}

# AWS static keys: consumed via AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY.
# Keyless installs (workload identity / node role) leave both empty and
# materialize nothing.
resource "kubernetes_secret_v1" "aws_credentials" {
  count = try(var.spec.aws_route53.access_key_id, "") != "" ? 1 : 0

  metadata {
    name      = local.aws_secret_name
    namespace = local.namespace
    labels    = local.labels
  }

  data = {
    "access-key-id"     = var.spec.aws_route53.access_key_id
    "secret-access-key" = var.spec.aws_route53.secret_access_key
  }

  depends_on = [kubernetes_namespace_v1.external_dns]
}

# GCP service-account key: mounted as a file with
# GOOGLE_APPLICATION_CREDENTIALS pointing at it (ADC's file path form).
resource "kubernetes_secret_v1" "gcp_credentials" {
  count = try(var.spec.google_cloud_dns.service_account_key_json, "") != "" ? 1 : 0

  metadata {
    name      = local.gcp_secret_name
    namespace = local.namespace
    labels    = local.labels
  }

  data = {
    "credentials.json" = var.spec.google_cloud_dns.service_account_key_json
  }

  depends_on = [kubernetes_namespace_v1.external_dns]
}

# Azure: the controller reads EVERYTHING (including identity mode) from a
# mounted azure.json — materialized from the typed fields whenever the
# azure_dns arm is set (even keyless modes need the file). Mounted at the
# controller's default config path, so no --azure-config-file override is
# needed.
resource "kubernetes_secret_v1" "azure_config" {
  count = try(var.spec.azure_dns, null) != null ? 1 : 0

  metadata {
    name      = local.azure_secret_name
    namespace = local.namespace
    labels    = local.labels
  }

  data = {
    "azure.json" = local.azure_config_json
  }

  depends_on = [kubernetes_namespace_v1.external_dns]
}

resource "helm_release" "external_dns" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the controller Deployment to become Available — a controller
  # that never starts (bad image, unschedulable pod) should fail THIS
  # apply, not the first record sync. Note the controller validates
  # provider CREDENTIALS at first sync, not at startup, so a live install
  # with wrong credentials still installs green and surfaces in controller
  # logs/records — by design (matching upstream behavior).
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 300

  # Two documents, merged in order by the provider (helm -f semantics):
  # the typed rendering first, the user's escape hatch last.
  values = [
    yamlencode(local.helm_values),
    try(var.spec.helm_values, ""),
  ]

  depends_on = [
    kubernetes_namespace_v1.external_dns,
    kubernetes_secret_v1.cloudflare_credentials,
    kubernetes_secret_v1.aws_credentials,
    kubernetes_secret_v1.gcp_credentials,
    kubernetes_secret_v1.azure_config,
  ]
}

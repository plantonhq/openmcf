# KubernetesKyverno Terraform module.
#
# Installs the Kyverno policy engine from the official chart as a real
# Helm release. The typed spec renders into chart values
# (locals.typed_helm_values); the helm_values escape hatch is passed as a
# SECOND values document, which the provider merges over the first with
# Helm -f semantics — the exact semantic twin of the Pulumi module's
# buildHelmValues + mergeMaps.
#
# WEBHOOK LIFECYCLE: the chart templates NO webhook configurations — the
# admission controller REGISTERS them at runtime and the chart's
# pre-delete cleanup hook is supposed to remove them at uninstall
# (webhooksCleanup, rendered explicitly). At the pinned chart the hook's
# delete-webhooks helper deletes ValidatingAdmissionPolicies instead of
# ValidatingWebhookConfigurations (upstream kyverno/kyverno#16492), so
# validating configs survive helm uninstall. The module-owned ConfigMap
# below runs the label-selected kubectl delete AFTER the release is
# gone — belt-and-suspenders until upstream ships a fixed
# readiness-checker. The spec's top-level comment still carries the
# manual unstick command for force-deleted releases.
#
# CRD LIFECYCLE: the policy CRDs are chart-TEMPLATED via the crds
# subchart — installed and DELETED with the release unless
# crds.keep_on_uninstall injects the resource-policy annotation.
# Destroying the engine cascade-deletes every policy on the cluster in
# the default posture (the spec's CRD field comments carry the warning).

# The optional installation namespace. Created before the release; deleted
# with the resource.
resource "kubernetes_namespace_v1" "kyverno" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# Destroy-ordered sentinel: helm_release depends on this ConfigMap, so
# destroy tears the release down FIRST and then runs the provisioner —
# after the admission pods are gone, so they cannot re-register the
# webhooks the chart's broken delete-webhooks helper left behind.
resource "kubernetes_config_map_v1" "webhook_gc" {
  metadata {
    name      = local.webhook_gc_config_map_name
    namespace = local.namespace
    labels    = local.labels
  }

  data = {
    purpose = "kyverno-webhook-gc"
  }

  provisioner "local-exec" {
    when = destroy
    # Destroy-time provisioners cannot interpolate getenv()/other
    # resources — only self/count/each — so the kubeconfig is resolved
    # by the shell from the tofu process environment (KUBECONFIG, or
    # KUBE_CONFIG_PATH which the Terraform Kubernetes provider uses).
    # $$ escapes so Terraform leaves the shell expansions alone.
    command = "KUBECONFIG=\"$${KUBECONFIG:-$${KUBE_CONFIG_PATH}}\" kubectl delete validatingwebhookconfiguration,mutatingwebhookconfiguration -l webhook.kyverno.io/managed-by=kyverno --ignore-not-found=true"
  }

  depends_on = [kubernetes_namespace_v1.kyverno]
}

resource "helm_release" "kyverno" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the controllers to become Ready — an admission engine that
  # never starts should fail THIS apply, not the first policy apply. The
  # timeout covers the four controller rollouts plus the runtime webhook
  # registration on a cold image pull.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the user's escape hatch second — and
  # fullnameOverride re-pinned LAST, the one deliberate exception to the
  # escape hatch's last-word contract (twin of the Pulumi module). Both
  # exported name outputs derive from the fullname; letting an override
  # move it would break them.
  values = concat(
    [yamlencode(local.typed_helm_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : [],
    [yamlencode({ fullnameOverride = local.release_name })]
  )

  depends_on = [
    kubernetes_namespace_v1.kyverno,
    kubernetes_config_map_v1.webhook_gc,
  ]
}

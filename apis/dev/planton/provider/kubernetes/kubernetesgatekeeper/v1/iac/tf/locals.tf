# Computed values for the KubernetesGatekeeper module.
# Every resolution here has an exact twin in the Pulumi module — keep
# them in lockstep: same rendered chart values, same outputs.
#
# HCL DISCIPLINE: conditional keys are contributed via merge() of
# `cond ? { key = value } : {}` singleton maps — a ternary whose branches
# are differently-shaped objects fails plan-time type unification.
# Optional nested blocks are read with try(): HCL's && does NOT
# short-circuit.
#
# NAMING: gatekeeper's chart HARDCODES its resource names
# (gatekeeper-webhook-service, gatekeeper-webhook-server-cert, the
# webhook configuration names) — there is no fullname derivation to pin
# and the engine is a per-cluster singleton by construction.

locals {
  # Pinned chart identity; chart_version resolves to the pinned default
  # when unset so both engines install the same chart whether or not the
  # platform's defaulting middleware ran. Chart and app versions move in
  # lockstep (chart 3.23.0 = Gatekeeper v3.23.0).
  helm_chart_name       = "gatekeeper"
  helm_chart_repo       = "https://open-policy-agent.github.io/gatekeeper/charts"
  default_chart_version = "3.23.0"
  chart_version         = try(var.spec.chart_version, "") != "" ? var.spec.chart_version : local.default_chart_version

  namespace    = var.spec.namespace
  release_name = var.metadata.name

  # Chart-FIXED names (hardcoded in the templates — no fullname
  # derivation). Exported as outputs.
  webhook_service_name     = "gatekeeper-webhook-service"
  webhook_cert_secret_name = "gatekeeper-webhook-server-cert"

  # Planton governance labels for the module-created namespace (never
  # injected into the chart's own resources; Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesGatekeeper"
    },
    try(var.metadata.id, "") != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    try(var.metadata.org, "") != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    try(var.metadata.env, "") != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # The self-management exemption label on the module-created namespace
  # must be DECLARED here, not left to the chart's post-install
  # label_namespace hook: the hook stamps it onto a namespace this module
  # owns, so an undeclared label is permanent config↔live drift — every
  # later apply would STRIP the exemption and let Gatekeeper police its
  # own namespace (fail-closed via the label-guard check webhook).
  # Rendered whenever the hook is enabled (empty = the chart default,
  # true), matching what the hook itself would stamp.
  label_namespace_enabled = try(var.spec.hooks.label_namespace, null) == null ? true : var.spec.hooks.label_namespace
  namespace_labels = merge(
    local.labels,
    local.label_namespace_enabled ? { "admission.gatekeeper.sh/ignore" = "no-self-managing" } : {}
  )

  # ---- shared renderers -----------------------------------------------------------
  # Single null-pruned comprehensions — never `cond ? {} : {for…}` (the
  # HCL type-unification class; a for-comprehension is a map, {} is an
  # object, and the ternary fails plan-time type-checking).
  resources_block = {
    for k, v in {
      requests = try(var.spec.resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = var.spec.resources.requests.cpu
          memory = var.spec.resources.requests.memory
        } : rk => rv if rv != null && rv != ""
      }
      limits = try(var.spec.resources.limits, null) == null ? null : {
        for rk, rv in {
          cpu    = var.spec.resources.limits.cpu
          memory = var.spec.resources.limits.memory
        } : rk => rv if rv != null && rv != ""
      }
    } : k => v if v != null
  }

  audit_resources_block = {
    for k, v in {
      requests = try(var.spec.audit.resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = var.spec.audit.resources.requests.cpu
          memory = var.spec.audit.resources.requests.memory
        } : rk => rv if rv != null && rv != ""
      }
      limits = try(var.spec.audit.resources.limits, null) == null ? null : {
        for rk, rv in {
          cpu    = var.spec.audit.resources.limits.cpu
          memory = var.spec.audit.resources.limits.memory
        } : rk => rv if rv != null && rv != ""
      }
    } : k => v if v != null
  }

  scheduling_tolerations = [
    for t in try(var.spec.scheduling.tolerations, []) : { for k, v in {
      key               = t.key != "" ? t.key : null
      operator          = t.operator != "" ? t.operator : null
      value             = t.value != "" ? t.value : null
      effect            = t.effect != "" ? t.effect : null
      tolerationSeconds = t.toleration_seconds
    } : k => v if v != null }
  ]

  # ---- controller-manager + audit blocks ----------------------------------------------
  # The scheduling entries apply to BOTH deployments — placing the engine
  # on dedicated nodes must move the audit loop with it, or audit
  # findings and enforcement diverge (the Pulumi twin fans out
  # identically).
  #
  # EXTERNAL CERT — chart-truth asymmetry at the pin: with
  # externalCertInjection enabled the AUDIT deployment auto-disables its
  # cert rotation but the CONTROLLER-MANAGER reads only its own flag;
  # without disableCertRotation=true the embedded rotator keeps
  # overwriting the injected Secret. The module sets it explicitly.
  controller_manager_block = merge(
    length(try(var.spec.exempt_namespaces, [])) > 0 ? { exemptNamespaces = var.spec.exempt_namespaces } : {},
    length(try(var.spec.exempt_namespace_prefixes, [])) > 0 ? { exemptNamespacePrefixes = var.spec.exempt_namespace_prefixes } : {},
    length(local.resources_block) > 0 ? { resources = local.resources_block } : {},
    length(try(var.spec.scheduling.node_selector, {})) > 0 ? { nodeSelector = var.spec.scheduling.node_selector } : {},
    length(local.scheduling_tolerations) > 0 ? { tolerations = local.scheduling_tolerations } : {},
    try(var.spec.external_cert, null) != null ? { disableCertRotation = true } : {}
  )

  audit_block = merge(
    length(local.audit_resources_block) > 0 ? { resources = local.audit_resources_block } : {},
    length(try(var.spec.scheduling.node_selector, {})) > 0 ? { nodeSelector = var.spec.scheduling.node_selector } : {},
    length(local.scheduling_tolerations) > 0 ? { tolerations = local.scheduling_tolerations } : {}
  )

  # ---- lifecycle hook blocks ------------------------------------------------------------
  post_install_block = {
    for k, v in {
      labelNamespace = try(var.spec.hooks.label_namespace, null) == null ? null : { enabled = var.spec.hooks.label_namespace }
      probeWebhook   = try(var.spec.hooks.probe_webhook, null) == null ? null : { enabled = var.spec.hooks.probe_webhook }
    } : k => v if v != null
  }

  # ---- image override (air-gap) ------------------------------------------------------------
  # The chart's image keys are repository + RELEASE (the tag key is named
  # "release" — chart-truth) with a separate crdRepository for the hook
  # containers; crdRepository and the curl hook image ride helm_values
  # for full air-gap installs (the spec comment teaches it). Each key
  # renders independently: pull_secret_name alone is a legal shape
  # (authenticated pulls of the DEFAULT image, e.g. Docker Hub
  # rate-limit credentials) — the Pulumi twin renders it the same way.
  image_block = {
    for k, v in {
      repository  = try(var.spec.image.repo, "") != "" ? var.spec.image.repo : null
      release     = try(var.spec.image.tag, "") != "" ? var.spec.image.tag : null
      pullSecrets = try(var.spec.image.pull_secret_name, "") != "" ? [{ name = var.spec.image.pull_secret_name }] : null
    } : k => v if v != null
  }

  # ---- typed chart values (Pulumi twin: buildHelmValues) --------------------------------------
  typed_helm_values = merge(
    try(var.spec.replicas, null) != null ? { replicas = var.spec.replicas } : {},

    # validating webhook
    try(var.spec.validating_webhook.enabled, null) != null ? { disableValidatingWebhook = !var.spec.validating_webhook.enabled } : {},
    try(var.spec.validating_webhook.failure_policy, "") != "" ? { validatingWebhookFailurePolicy = var.spec.validating_webhook.failure_policy } : {},
    try(var.spec.validating_webhook.timeout_seconds, null) != null ? { validatingWebhookTimeoutSeconds = var.spec.validating_webhook.timeout_seconds } : {},
    try(var.spec.validating_webhook.enable_delete_operations, false) ? { enableDeleteOperations = true } : {},
    try(var.spec.validating_webhook.check_ignore_failure_policy, "") != "" ? { validatingWebhookCheckIgnoreFailurePolicy = var.spec.validating_webhook.check_ignore_failure_policy } : {},

    # mutating webhook
    try(var.spec.mutating_webhook.enabled, null) != null ? { disableMutation = !var.spec.mutating_webhook.enabled } : {},
    try(var.spec.mutating_webhook.failure_policy, "") != "" ? { mutatingWebhookFailurePolicy = var.spec.mutating_webhook.failure_policy } : {},
    try(var.spec.mutating_webhook.timeout_seconds, null) != null ? { mutatingWebhookTimeoutSeconds = var.spec.mutating_webhook.timeout_seconds } : {},
    try(var.spec.mutating_webhook.mutation_annotations, false) ? { mutationAnnotations = true } : {},

    # audit
    try(var.spec.audit.interval_seconds, null) != null ? { auditInterval = var.spec.audit.interval_seconds } : {},
    try(var.spec.audit.constraint_violations_limit, null) != null ? { constraintViolationsLimit = var.spec.audit.constraint_violations_limit } : {},
    try(var.spec.audit.from_cache, false) ? { auditFromCache = true } : {},
    try(var.spec.audit.match_kind_only, false) ? { auditMatchKindOnly = true } : {},
    try(var.spec.audit.chunk_size, null) != null ? { auditChunkSize = var.spec.audit.chunk_size } : {},

    # engine capabilities
    try(var.spec.engine.enable_external_data, null) != null ? { enableExternalData = var.spec.engine.enable_external_data } : {},
    try(var.spec.engine.enable_k8s_native_validation, null) != null ? { enableK8sNativeValidation = var.spec.engine.enable_k8s_native_validation } : {},
    try(var.spec.engine.enable_generator_resource_expansion, null) != null ? { enableGeneratorResourceExpansion = var.spec.engine.enable_generator_resource_expansion } : {},
    length(try(var.spec.engine.disabled_builtins, [])) > 0 ? { disabledBuiltins = var.spec.engine.disabled_builtins } : {},
    try(var.spec.engine.log_denies, false) ? { logDenies = true } : {},
    try(var.spec.engine.log_level, "") != "" ? { logLevel = var.spec.engine.log_level } : {},

    # external webhook certificate (the cert-manager arm)
    try(var.spec.external_cert, null) != null ? {
      externalCertInjection = {
        enabled    = true
        secretName = var.spec.external_cert.secret_name
      }
    } : {},

    # lifecycle hooks
    length(local.post_install_block) > 0 ? { postInstall = local.post_install_block } : {},
    try(var.spec.hooks.upgrade_crds, null) != null ? { upgradeCRDs = { enabled = var.spec.hooks.upgrade_crds } } : {},
    # The extra `name` key works around a chart bug at the 3.23.0 pin:
    # the hook's ClusterRoleBinding subject renders
    # .Values.preUninstall.deleteWebhookConfigurations.name — a key that
    # does not exist in the chart's own values (the SA name lives under
    # .serviceAccount.name) — so enabling the arm without it fails every
    # uninstall on the CRB's empty subject. Value = the chart's SA-name
    # default. Pulumi twin renders it identically.
    try(var.spec.hooks.delete_webhook_configurations_on_uninstall, false) ? {
      preUninstall = { deleteWebhookConfigurations = { enabled = true, name = "gatekeeper-delete-webhook-configs" } }
    } : {},

    # image override
    length(local.image_block) > 0 ? { image = local.image_block } : {},

    # deployments
    length(local.controller_manager_block) > 0 ? { controllerManager = local.controller_manager_block } : {},
    length(local.audit_block) > 0 ? { audit = local.audit_block } : {}
  )
}

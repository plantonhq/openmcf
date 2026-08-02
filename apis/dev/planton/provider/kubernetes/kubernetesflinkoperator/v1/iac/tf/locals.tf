# Computed values for the KubernetesFlinkOperator module. Every resolution
# here has an exact twin in the Pulumi module's locals.go / values.go —
# keep them in lockstep.
#
# HCL DISCIPLINE (applies to every conditional object in this file):
# conditional entries are written as `key = cond ? value : null` inside ONE
# object literal, pruned with `{ for k, v in {...} : k => v if v != null }`.
# The two tempting alternatives are both broken: `cond ? {...} : {}`
# ternaries fail plan-time type unification when branches carry different
# attributes, and merge() over conditional lists silently UNIFIES
# primitive-only sibling objects into map(string). The null-prune form
# preserves every value's type.
#
# Optional nested blocks are read with try(): HCL's && does NOT
# short-circuit, so chained null checks still dereference the null.

locals {
  # Chart identity — must stay byte-identical with the Pulumi module's
  # vars: cross-engine chart drift deploys two different products from one
  # manifest.
  helm_chart_name = "flink-kubernetes-operator"

  # Release name = metadata.name.
  release_name = var.metadata.name

  # Chart version resolved to the pinned default when unset, so both
  # engines install the same chart whether or not the platform's
  # defaulting middleware ran — mirror of the Pulumi module's
  # DefaultChartVersion. Chart version = operator version = the image
  # tag the module pins. The flink.apache.org CRDs ship from the
  # chart's crds/ directory: Helm installs them once and NEVER
  # upgrades them — bumping this version does not touch the CRDs
  # (apply the new release's CRD files manually when a bump changes
  # them).
  chart_version = coalesce(try(var.spec.chart_version, null), "1.15.0")

  # The chart is served PER VERSION from a versioned Apache downloads
  # directory — the version is part of the repository URL itself, not
  # just of the chart pin. Twin of the Pulumi module's
  # HelmChartRepoTemplate.
  helm_chart_repo = "https://downloads.apache.org/flink/flink-kubernetes-operator-${local.chart_version}/"

  namespace = var.spec.namespace

  # Resource-identity labels stamped on the module-created satellites
  # (the namespace and the keystore-password Secret — never injected
  # into the chart's own resources; Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesFlinkOperator"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- webhook (chart: webhook.*) ---------------------------------------------
  # Unset-or-true = enabled (the upstream default this spec keeps). With
  # the webhook enabled the chart renders cert-manager Issuer/Certificate
  # UNCONDITIONALLY (cert-manager is this kind's registry prerequisite)
  # and both webhook configurations are failurePolicy Fail.
  webhook_enabled = coalesce(try(var.spec.webhook_enabled, null), true)

  # The module-owned keystore-password Secret. The chart's own default
  # is a HARDCODED PUBLIC PASSWORD ("password1234", base64 in
  # templates/webhook/secret.yaml behind keystore.useDefaultPassword=
  # true) — it must never ship, so the module generates a random
  # password per install and points the chart at it.
  webhook_keystore_secret_name = "${var.metadata.name}-webhook-keystore"

  # The webhook Service name is CHART-FIXED (templates/webhook/
  # service.yaml hardcodes it) — not fullname-derived. Empty when the
  # webhook is disabled; matches stack_outputs.proto.
  webhook_service = local.webhook_enabled ? "flink-operator-webhook-service" : ""

  # ---- operator configuration (chart: defaultConfiguration) --------------------
  # The operator is Flink-config-file configured. Leader election is
  # module-owned: any replica count beyond 1 REQUIRES it (the chart's
  # own contract — it refuses multi-replica installs without it), so
  # the two leader-election keys render exactly when replicas > 1 —
  # never a spec knob that could drift from the replica count. Key
  # spelling verified against the operator docs
  # (kubernetes.operator.leader-election.enabled / .lease-name).
  replicas               = try(var.spec.replicas, null)
  leader_election_needed = coalesce(local.replicas, 1) > 1
  operator_config = merge(
    try(var.spec.operator_config, {}),
    local.leader_election_needed ? {
      "kubernetes.operator.leader-election.enabled"    = "true"
      "kubernetes.operator.leader-election.lease-name" = "${var.metadata.name}-lease"
    } : {}
  )
  # NOTE the Flink conf format is "key: value" (colon-space), NOT
  # "key=value" — the file is YAML-flavored Flink configuration.
  flink_conf_file = join("\n", [
    for k in sort(keys(local.operator_config)) : "${k}: ${local.operator_config[k]}"
  ])

  # ---- job service account (chart: jobServiceAccount) --------------------------
  # "flink" is the chart default — the name every FlinkDeployment
  # references by default. The chart marks the service account
  # `helm.sh/resource-policy: keep`: it survives uninstall so running
  # jobs never lose their identity.
  job_service_account = coalesce(try(var.spec.job_service_account, null), "flink")

  # ---- operator container resources (shared ContainerResources) ----------------
  # Twin of the Pulumi module's resourcesMap. The chart ships NO
  # default requests/limits for the operator container — the resources
  # key renders only when the spec sets them. Helm deep-merges per
  # key: a partial block overrides only the halves it carries.
  operator_resources = try(var.spec.resources, null) == null ? null : {
    for k, v in {
      limits = try(var.spec.resources.limits, null) == null ? null : {
        for lk, lv in {
          cpu    = try(var.spec.resources.limits.cpu, "") != "" ? var.spec.resources.limits.cpu : null
          memory = try(var.spec.resources.limits.memory, "") != "" ? var.spec.resources.limits.memory : null
        } : lk => lv if lv != null
      }
      requests = try(var.spec.resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = try(var.spec.resources.requests.cpu, "") != "" ? var.spec.resources.requests.cpu : null
          memory = try(var.spec.resources.requests.memory, "") != "" ? var.spec.resources.requests.memory : null
        } : rk => rv if rv != null
      }
    } : k => v if v != null && length(v) > 0
  }

  # ---- operator pod block (resources + scheduling) ------------------------------
  # priorityClassName lives under operatorPod in this chart's values
  # (template-verified) — alongside nodeSelector/tolerations.
  operator_pod = {
    for k, v in {
      resources    = local.operator_resources != null && length(local.operator_resources) > 0 ? local.operator_resources : null
      nodeSelector = length(try(var.spec.scheduling.node_selector, {})) > 0 ? var.spec.scheduling.node_selector : null
      tolerations = length(try(var.spec.scheduling.tolerations, [])) > 0 ? [
        for t in var.spec.scheduling.tolerations : {
          for tk, tv in {
            key               = try(t.key, "") != "" ? t.key : null
            operator          = try(t.operator, "") != "" ? t.operator : null
            value             = try(t.value, "") != "" ? t.value : null
            effect            = try(t.effect, "") != "" ? t.effect : null
            tolerationSeconds = try(t.toleration_seconds, null)
          } : tk => tv if tv != null
        }
      ] : null
      priorityClassName = try(var.spec.scheduling.priority_class_name, "") != "" ? var.spec.scheduling.priority_class_name : null
    } : k => v if v != null
  }

  # ---- operator image -----------------------------------------------------------
  # tag is ALWAYS pinned to the chart version: the chart's own default
  # tag is the unpinned "latest" — the one values.yaml default that
  # must never stand. image_registry replaces ONLY the registry part
  # (chart default `ghcr.io/apache/flink-kubernetes-operator`); this
  # never rewrites the Flink images deployments run — those ride each
  # KubernetesFlinkDeployment's own image field. Twin of the Pulumi
  # module.
  operator_image = {
    for k, v in {
      repository = try(var.spec.image_registry, "") != "" ? "${var.spec.image_registry}/apache/flink-kubernetes-operator" : null
      tag        = local.chart_version
    } : k => v if v != null
  }

  # ---- typed chart values (twin of the Pulumi module's buildHelmValues) --------
  # Chart-default-matching values render only on divergence — with the
  # deliberate always-rendered exceptions called out inline
  # (nameOverride, image.tag, and the keystore wiring whenever the
  # webhook is on).
  typed_values = {
    for k, v in {
      # nameOverride is THIS chart's identity pin: the operator
      # Deployment (and its selector labels) render from the
      # `flink-operator.name` helper (default .Chart.Name |
      # nameOverride) — fullnameOverride alone is a no-op for the
      # Deployment and leaves it at the chart-constant
      # `flink-kubernetes-operator` (verified live: the pinned name was
      # NotFound while the chart-named Deployment served). Both keys
      # are set so any fullname-derived fallback also hangs off the
      # resource name.
      nameOverride     = local.release_name
      fullnameOverride = local.release_name

      image = local.operator_image

      # Scopes RBAC AND the admission webhook's namespaceSelector to
      # exactly these namespaces (template-verified: the chart wires
      # kubernetes.operator.watched.namespaces from the same list).
      watchNamespaces = length(try(var.spec.watch_namespaces, [])) > 0 ? var.spec.watch_namespaces : null

      # Enabled: the module-owned keystore password replaces the
      # chart's hardcoded-public default. Disabled: webhook.create=
      # false removes the webhook, the certificate machinery, and the
      # cert-manager dependency. Null-pruned, never a two-branch
      # ternary — the branches carry different attributes and HCL
      # fails plan-time type unification on them (the discipline
      # header's first class).
      webhook = {
        for wk, wv in {
          create = local.webhook_enabled ? null : false
          keystore = !local.webhook_enabled ? null : {
            useDefaultPassword = false
            passwordSecretRef = {
              name = local.webhook_keystore_secret_name
              key  = "password"
            }
          }
        } : wk => wv if wv != null
      }

      # Rendered on presence — an explicit 1 re-states the chart
      # default harmlessly; >1 pairs with the leader-election keys
      # rendered into operator_config above (the chart REFUSES
      # multi-replica installs without leader election, by design).
      replicas = local.replicas

      # create:true and append:true ARE the chart defaults — rendered
      # explicitly alongside the file for self-documentation (append
      # keeps the chart's built-in conf underneath; ours layers over).
      defaultConfiguration = length(local.operator_config) > 0 ? {
        create            = true
        append            = true
        "flink-conf.yaml" = local.flink_conf_file
      } : null

      # Rendered only on divergence from the chart default "flink".
      # The chart keeps the service account past uninstall
      # (helm.sh/resource-policy: keep) — running jobs never lose
      # their identity.
      jobServiceAccount = local.job_service_account != "flink" ? { name = local.job_service_account } : null

      operatorPod = length(local.operator_pod) > 0 ? local.operator_pod : null

      imagePullSecrets = length(try(var.spec.image_pull_secrets, [])) > 0 ? [for s in var.spec.image_pull_secrets : { name = s }] : null
    } : k => v if v != null && (!can(length(v)) || length(v) > 0)
  }

  # ---- keystore re-pin (third values document, merged AFTER the escape hatch) --
  # useDefaultPassword=false is re-pinned whenever the webhook is
  # enabled: the chart's hardcoded-password default must not resurface
  # through helm_values. Twin of the Pulumi module's post-merge re-pin.
  keystore_repin_values = {
    webhook = {
      keystore = {
        useDefaultPassword = false
      }
    }
  }
}

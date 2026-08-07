# Computed values for the KubernetesPerconaMysqlOperator module. Every
# resolution here has an exact twin in the Pulumi module's locals.go /
# values.go — keep them in lockstep.
#
# HCL DISCIPLINE (applies to every conditional object in this file):
# conditional entries are written as `key = cond ? value : null` inside ONE
# object literal, pruned with `{ for k, v in {...} : k => v if v != null }`.
# The two tempting alternatives are both broken: `cond ? {...} : {}`
# ternaries fail plan-time type unification when branches carry different
# attributes, and `merge(concat(cond ? [{...}] : [], ...)...)` silently
# UNIFIES primitive-only sibling objects into map(string) — numbers and
# booleans arrive in the chart values as strings. The null-prune form
# preserves every value's type.
#
# Optional nested blocks are read with try(): HCL's && does NOT
# short-circuit, so chained null checks still dereference the null — and
# var.spec is typed 'any', so an absent attribute is an error, not a null.

locals {
  # Chart identity — must stay byte-identical with the Pulumi module's
  # vars: cross-engine chart drift deploys two different products from one
  # manifest. Percona publishes the operator charts for every database
  # product from this one repository.
  helm_chart_repo = "https://percona.github.io/percona-helm-charts"
  helm_chart_name = "pxc-operator"

  # Release name = metadata.name. The chart derives every resource name
  # from the release (Deployment, ServiceAccount, RBAC through its
  # fullname helper), so distinct release names keep multiple
  # namespace-scoped installations from colliding.
  release_name = var.metadata.name

  # Chart version resolved to the pinned default when unset, so both
  # engines install the same chart whether or not the platform's
  # defaulting middleware ran — mirror of the Pulumi module's
  # DefaultChartVersion. Chart and operator versions move TOGETHER for
  # this chart (chart 1.20.0 ships operator 1.20.0); the chart pin
  # governs.
  chart_version = coalesce(try(var.spec.chart_version, null), "1.20.0")

  # The chart's own operatorImageRepository default — needed for the
  # tag-without-repository image override: the chart's repository value
  # is always suffixed with the chart's app version, so pinning a
  # DIFFERENT tag requires the chart's full-image override
  # ("<repository>:<tag>"), and the repository half must come from
  # somewhere when the spec leaves it empty.
  default_image_repository = "percona/percona-xtradb-cluster-operator"

  namespace = var.spec.namespace

  # Widened watch (cluster-wide or a namespace fence) grants the operator
  # ClusterRole permissions, under which it registers ONE shared
  # cluster-scoped ValidatingWebhookConfiguration at startup. The module
  # owns that object's lifecycle in the widened arms (see main.tf) — this
  # is the render gate.
  watch_widened = try(var.spec.watch.cluster_wide, false) || length(try(var.spec.watch.namespaces, [])) > 0

  # Resource-identity labels stamped on the namespace this module creates
  # (never injected into the chart's own resources — Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesPerconaMysqlOperator"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- operator container resources (shared ContainerResources shape) -------
  # Twin of the Pulumi module's resourcesMap. The chart SHIPS default
  # requests/limits (requests 100m/20Mi, limits 200m/500Mi) — the
  # resources key renders only when the spec sets them, so the chart
  # defaults survive an empty spec. Helm deep-merges per key, so a
  # partial spec block overrides only the halves it carries.
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

  # ---- image ------------------------------------------------------------------
  # The image override is TWO chart values with a precedence rule: the
  # chart's image helper uses `image` verbatim when non-empty, else
  # "<operatorImageRepository>:<chart app version>". A repository alone
  # therefore maps to operatorImageRepository (tag stays the chart
  # version); any custom TAG requires the full `image` override, with the
  # repository half falling back to the chart's own default repository
  # when the spec leaves it empty.
  image_repository = try(var.spec.image.repository, "")
  image_tag        = try(var.spec.image.tag, "")

  # ---- typed chart values (twin of the Pulumi module's buildHelmValues) ------
  # Chart-default-matching values render only on divergence
  # (watchAllNamespaces, logStructured, disableTelemetry, and the
  # xtrabackupSidecar gate all default off upstream), so the rendered
  # values stay minimal on both engines.
  typed_values = {
    for k, v in {
      replicaCount = try(var.spec.replicas, null)
      resources    = local.operator_resources != null && length(local.operator_resources) > 0 ? local.operator_resources : null

      # Watch scope maps to two independent chart values: cluster-wide
      # RBAC (watchAllNamespaces) or a comma-joined namespace fence
      # (watchNamespace). Spec CEL rules make the two arms mutually
      # exclusive, so at most one renders; both unset means the operator
      # watches its own namespace (the upstream default). The chart's own
      # createNamespace value (create the WATCHED namespaces) is never
      # rendered — the module owns only the installation namespace, and
      # watched namespaces must already exist.
      watchAllNamespaces = try(var.spec.watch.cluster_wide, false) ? true : null
      watchNamespace     = length(try(var.spec.watch.namespaces, [])) > 0 ? join(",", var.spec.watch.namespaces) : null

      # Both integers in this chart's values file (unlike the sibling
      # MongoDB operator chart, which declares its default as a string).
      maxConcurrentReconciles = try(var.spec.max_concurrent_reconciles, null)
      s3WorkersLimit          = try(var.spec.s3_workers_limit, null)

      # logStructured matches the chart default (false) — rendered only
      # when on. logLevel renders whenever the spec carries a value (the
      # chart default is "INFO"; re-stating it is harmless and keeps
      # rendering presence-driven, not value-driven).
      logStructured = try(var.spec.log.structured, false) ? true : null
      logLevel      = try(var.spec.log.level, "") != "" ? var.spec.log.level : null

      # Chart default false (telemetry on) — rendered only on explicit
      # opt-out.
      disableTelemetry = try(var.spec.disable_telemetry, false) ? true : null

      # Leader election flattens to four top-level chart values, each
      # rendered on presence: the enabled flag whenever the spec carries
      # it (the chart default is true, so an explicit true is a harmless
      # re-statement), the three timing knobs whenever non-empty.
      leaderElectionEnabled = try(var.spec.leader_election.enabled, null)
      leaseDuration         = try(var.spec.leader_election.lease_duration, "") != "" ? var.spec.leader_election.lease_duration : null
      renewDeadline         = try(var.spec.leader_election.renew_deadline, "") != "" ? var.spec.leader_election.renew_deadline : null
      retryPeriod           = try(var.spec.leader_election.retry_period, "") != "" ? var.spec.leader_election.retry_period : null

      # The chart folds every gate into one PXCO_FEATURE_GATES environment
      # variable; xtrabackupSidecar is the only gate it declares. Rendered
      # only when on (chart default false) — Helm deep-merges the map, so
      # the chart's own entry is replaced, not duplicated.
      featureGates = try(var.spec.xtrabackup_sidecar, false) ? { xtrabackupSidecar = true } : null

      nodeSelector = length(try(var.spec.node_selector, {})) > 0 ? var.spec.node_selector : null
      tolerations = length(try(var.spec.tolerations, [])) > 0 ? [
        for t in var.spec.tolerations : {
          for tk, tv in {
            key               = try(t.key, "") != "" ? t.key : null
            operator          = try(t.operator, "") != "" ? t.operator : null
            value             = try(t.value, "") != "" ? t.value : null
            effect            = try(t.effect, "") != "" ? t.effect : null
            tolerationSeconds = try(t.toleration_seconds, null)
          } : tk => tv if tv != null
        }
      ] : null

      # Pull secrets are the raw Kubernetes object list ([{name: ...}]) —
      # the chart pipes imagePullSecrets straight into the pod spec with
      # toYaml.
      imagePullSecrets = length(try(var.spec.image_pull_secrets, [])) > 0 ? [
        for s in var.spec.image_pull_secrets : { name = s }
      ] : null

      operatorImageRepository = local.image_repository != "" ? local.image_repository : null
      image                   = local.image_tag != "" ? "${local.image_repository != "" ? local.image_repository : local.default_image_repository}:${local.image_tag}" : null
    } : k => v if v != null
  }
}

# Computed values for the KubernetesSolrOperator module. Every resolution
# here has an exact twin in the Pulumi module's locals.go / values.go —
# keep them in lockstep.
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
  # manifest. Chart versions carry NO `v` prefix (the operator image/CRD
  # artifacts DO — chart 0.9.1 ships operator v0.9.1); the SERVED index
  # (https://solr.apache.org/charts) governs the version.
  helm_chart_repo = "https://solr.apache.org/charts"
  helm_chart_name = "solr-operator"

  # Release name = metadata.name.
  release_name = var.metadata.name

  # Chart version resolved to the pinned default when unset, so both
  # engines install the same chart whether or not the platform's
  # defaulting middleware ran — mirror of the Pulumi module's
  # DefaultChartVersion.
  chart_version = coalesce(try(var.spec.chart_version, null), "0.9.1")

  namespace = var.spec.namespace

  # Resource-identity labels stamped on the namespace this module creates
  # (never injected into the chart's own resources — Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesSolrOperator"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- module-owned CRDs ------------------------------------------------------
  # The solr-operator chart ships NO CRDs (they are separate release
  # artifacts) — the module owns the four files staged at ../crds: the
  # three solr.apache.org CRDs plus the ZookeeperCluster CRD of the
  # bundled zookeeper-operator dependency.
  #
  # Each file is split into YAML documents on SEPARATOR LINES ("---" at
  # column 0 — "\n---" never matches the "---" substrings the CRD
  # schemas embed in description text, which all sit mid-line or
  # indented). The upstream files open with license-comment headers and
  # the ZookeeperCluster file carries a doubled leading "---"; the
  # can(yamldecode) filter drops the comment-only and empty fragments.
  # Keyed by each CRD's OWN metadata.name (never a file name or split
  # index) so state addresses stay stable.
  crd_documents = {
    for doc in flatten([
      for f in fileset("${path.module}/../crds", "*.yaml") :
      split("\n---", file("${path.module}/../crds/${f}"))
    ]) :
    yamldecode(doc).metadata.name => doc
    if trimspace(doc) != "" && can(yamldecode(doc).metadata.name)
  }

  # ---- operator container resources (shared ContainerResources shape) -------
  # Twin of the Pulumi module's resourcesMap. The chart ships NO default
  # resources (`resources: {}` — the operator is lightweight); the key
  # renders only when the spec sets it.
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

  # ---- bundled zookeeper-operator values -------------------------------------
  # NOTE the dash in the chart's values key: "zookeeper-operator" is the
  # SUBCHART name, quoted wherever it appears as a map key. crd.create
  # is pinned false UNCONDITIONALLY — the module owns the
  # ZookeeperCluster CRD (applied apply_only below), so the subchart
  # must never install its own copy and put it under Helm's
  # delete-on-uninstall lifecycle. install renders on presence (chart
  # default true); use renders only on divergence (chart default false,
  # and it is ignored whenever install is true).
  zookeeper_operator_values = {
    for k, v in {
      crd     = { create = false }
      install = try(var.spec.zookeeper_operator.install, null)
      use     = try(var.spec.zookeeper_operator.use_existing, false) ? true : null
    } : k => v if v != null
  }

  # ---- operator -> Solr mutual TLS --------------------------------------------
  # Secret names arrive as resolved literals (the value-or-ref foreign
  # keys resolve before Terraform runs). Scalars with chart defaults
  # (ca_cert_secret_key, insecure_skip_verify, watch_for_updates) render
  # on presence.
  mtls_values = try(var.spec.mtls, null) == null ? null : {
    for k, v in {
      clientCertSecret   = try(var.spec.mtls.client_cert_secret, "") != "" ? var.spec.mtls.client_cert_secret : null
      caCertSecret       = try(var.spec.mtls.ca_cert_secret, "") != "" ? var.spec.mtls.ca_cert_secret : null
      caCertSecretKey    = try(var.spec.mtls.ca_cert_secret_key, null) != null ? var.spec.mtls.ca_cert_secret_key : null
      insecureSkipVerify = try(var.spec.mtls.insecure_skip_verify, null)
      watchForUpdates    = try(var.spec.mtls.watch_for_updates, null)
    } : k => v if v != null
  }

  # ---- image source ------------------------------------------------------------
  # image.repository / image.tag are the air-gap override; the pull
  # secret is the chart's image.imagePullSecret — a SINGULAR string (the
  # chart accepts exactly one), not the usual imagePullSecrets object
  # list.
  image_values = {
    for k, v in {
      repository      = try(var.spec.image.repository, "") != "" ? var.spec.image.repository : null
      tag             = try(var.spec.image.tag, "") != "" ? var.spec.image.tag : null
      imagePullSecret = try(var.spec.image_pull_secret, "") != "" ? var.spec.image_pull_secret : null
    } : k => v if v != null
  }

  # ---- typed chart values (twin of the Pulumi module's buildHelmValues) ------
  # Chart-default-matching values render only on divergence, so the
  # rendered values stay minimal on both engines. The ONE always-rendered
  # entry is "zookeeper-operator" (it always carries crd.create = false).
  typed_values = {
    for k, v in {
      replicaCount = try(var.spec.replicas, null)
      resources    = local.operator_resources != null && length(local.operator_resources) > 0 ? local.operator_resources : null

      # The chart takes a COMMA-JOINED string (templates/_helpers.tpl
      # splits it back apart), not a YAML list. Empty = the operator
      # watches ALL namespaces (the chart default; also the chart's
      # ClusterRole-vs-Role switch).
      watchNamespaces = length(try(var.spec.watch_namespaces, [])) > 0 ? join(",", var.spec.watch_namespaces) : null

      "zookeeper-operator" = local.zookeeper_operator_values

      # Both nest a single enable flag and default true in the chart.
      # Rendered on presence — an explicit true re-states the chart
      # default harmlessly, an explicit false is the actual opt-out.
      leaderElection = try(var.spec.leader_election_enabled, null) != null ? { enable = var.spec.leader_election_enabled } : null
      metrics        = try(var.spec.metrics_enabled, null) != null ? { enable = var.spec.metrics_enabled } : null

      mTLS = local.mtls_values != null && length(local.mtls_values) > 0 ? local.mtls_values : null

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

      image = length(local.image_values) > 0 ? local.image_values : null
    } : k => v if v != null
  }

  # ---- deployment name (the chart's fullname template, replayed) -------------
  # templates/_helpers.tpl "solr-operator.fullname" with no
  # nameOverride/fullnameOverride (the typed spec sets neither): the
  # release name itself when it already contains "solr-operator",
  # otherwise "<release>-solr-operator" — truncated to 63 chars with one
  # trailing "-" trimmed. A helm_values override of either template knob
  # would change the real name without being reflected here — the
  # escape-hatch caveat both engines document.
  deployment_name = trimsuffix(
    substr(
      strcontains(local.release_name, local.helm_chart_name) ? local.release_name : "${local.release_name}-${local.helm_chart_name}",
      0, 63
    ),
  "-")
}

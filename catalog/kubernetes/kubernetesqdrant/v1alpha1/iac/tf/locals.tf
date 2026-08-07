# Computed values for the KubernetesQdrant module. Every resolution here
# has an exact twin in the Pulumi module's locals.go / values.go — keep
# them in lockstep.
#
# HCL DISCIPLINE (applies to every conditional object in this file):
# conditional entries are written as `key = cond ? value : null` inside ONE
# object literal, pruned with `{ for k, v in {...} : k => v if v != null }`.
# The two tempting alternatives are both broken: `cond ? {...} : {}`
# ternaries fail plan-time type unification when branches carry different
# attributes, and merge() over primitive-only sibling objects silently
# UNIFIES them into map(string) — numbers and booleans arrive in the chart
# values as strings. The null-prune form preserves every value's type.
#
# Optional nested blocks are read with try(): HCL's && does NOT
# short-circuit, so chained null checks still dereference the null. Optional
# SCALARS inside optional blocks are read with try(coalesce(x), null) — the
# null-safe read (coalesce rejects a lone null; try turns that rejection
# into the fallback).

locals {
  # Chart identity — must stay byte-identical with the Pulumi module's
  # vars: cross-engine chart-name drift deploys two different products
  # from one manifest.
  helm_chart_name = "qdrant"
  helm_chart_repo = "https://qdrant.github.io/qdrant-helm"

  # Release name — metadata.name, NOT a fixed chart name: several Qdrant
  # clusters coexist in one Kubernetes cluster. fullnameOverride below
  # pins every chart child name to this.
  release_name = var.metadata.name

  # Chart version resolved to the pinned default when unset, so both
  # engines install the same chart whether or not the platform's
  # defaulting middleware ran — mirror of the Pulumi module's
  # DefaultChartVersion. The chart version tracks the Qdrant release it
  # ships.
  chart_version = coalesce(var.spec.chart_version, "1.18.2")

  namespace = var.spec.namespace

  # Resource-identity labels stamped on the module-created satellites
  # (the namespace — never injected into the chart's own resources; Helm
  # owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesQdrant"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # The chart's main ClusterIP Service is qdrant.fullname — pinned to the
  # resource name via fullnameOverride. Feeds the endpoint outputs.
  service_name = local.release_name

  # ---- storage ---------------------------------------------------------
  storage_size       = try(coalesce(var.spec.storage.size), null) != null ? var.spec.storage.size : "10Gi"
  storage_class      = try(var.spec.storage.storage_class, "")
  snapshots_declared = try(var.spec.snapshots, null) != null
  snapshots_size     = try(coalesce(var.spec.snapshots.size), null) != null ? var.spec.snapshots.size : "10Gi"
  snapshots_class    = try(var.spec.snapshots.storage_class, "")

  # ---- api keys ----------------------------------------------------------
  # The chart owns key materialization either way: the generate arm
  # renders `apiKey: true` (the chart creates a random key ONCE and keeps
  # it stable across upgrades via its lookup); the existing arm renders
  # the chart's valueFrom shape — the chart reads the referenced Secret
  # AT TEMPLATE TIME (it must exist before the install) and copies the
  # key into its own `<name>-apikey` Secret. Key MATERIAL never lands in
  # chart values.
  # The two arms carry DIFFERENT TYPES (bool true for generate, an object
  # for the existing-secret ref), which HCL's ternary cannot unify —
  # jsonencode makes every branch a string and jsondecode restores the
  # dynamic value at the use site.
  api_key_existing = try(var.spec.api_key.existing_secret, null)
  api_key_generate = try(var.spec.api_key.generate, false)
  api_key_value_json = local.api_key_existing != null ? jsonencode({
    valueFrom = {
      secretKeyRef = {
        name = local.api_key_existing.name
        key  = local.api_key_existing.key
      }
    }
  }) : (local.api_key_generate ? jsonencode(true) : null)

  read_only_api_key_existing = try(var.spec.read_only_api_key.existing_secret, null)
  read_only_api_key_generate = try(var.spec.read_only_api_key.generate, false)
  read_only_api_key_value_json = local.read_only_api_key_existing != null ? jsonencode({
    valueFrom = {
      secretKeyRef = {
        name = local.read_only_api_key_existing.name
        key  = local.read_only_api_key_existing.key
      }
    }
  }) : (local.read_only_api_key_generate ? jsonencode(true) : null)

  # Name of the chart-owned Secret carrying the API key material
  # (`<name>-apikey`, keys api-key / read-only-api-key) — exported when
  # any key is declared.
  api_key_secret_name           = try(var.spec.api_key, null) != null ? "${local.release_name}-apikey" : ""
  read_only_api_key_secret_name = try(var.spec.read_only_api_key, null) != null ? "${local.release_name}-apikey" : ""

  # ---- tls -----------------------------------------------------------------
  # The typed tls block turns on the engine's service.enable_tls and
  # points config.tls at the mounted certificate Secret; the chart's own
  # probes switch to HTTPS off the same flag. The mount path must stay
  # identical to the Pulumi module's TlsMountPath.
  tls_secret     = try(var.spec.tls.secret, "")
  tls_enabled    = local.tls_secret != ""
  tls_mount_path = "/qdrant/tls"
  http_scheme    = local.tls_enabled ? "https" : "http"

  # ---- typed chart values (twin of the Pulumi module's buildHelmValues) --
  helm_values = {
    for k, v in {
      # fullnameOverride pins qdrant.fullname to the resource name: child
      # names (the Service, `-headless`, the `-apikey` Secret, the
      # ConfigMap) derive deterministically and the longest suffix stays
      # far from the 63-char ceiling.
      fullnameOverride = local.release_name

      # Distributed mode is the chart default (config.cluster.enabled:
      # true); pod 0 bootstraps consensus and later pods join over p2p.
      replicaCount = try(coalesce(var.spec.replicas), null) != null ? var.spec.replicas : 1

      resources = try(var.spec.resources, null) == null ? null : {
        for rk, rv in {
          requests = try(var.spec.resources.requests, null) == null ? null : {
            for qk, qv in {
              cpu    = var.spec.resources.requests.cpu != "" ? var.spec.resources.requests.cpu : null
              memory = var.spec.resources.requests.memory != "" ? var.spec.resources.requests.memory : null
            } : qk => qv if qv != null
          }
          limits = try(var.spec.resources.limits, null) == null ? null : {
            for lk, lv in {
              cpu    = var.spec.resources.limits.cpu != "" ? var.spec.resources.limits.cpu : null
              memory = var.spec.resources.limits.memory != "" ? var.spec.resources.limits.memory : null
            } : lk => lv if lv != null
          }
        } : rk => rv if rv != null && rv != {}
      }

      # The chart's persistence block always renders a PVC; size resolved
      # to the spec default, storageClassName only when declared (absent
      # = the cluster's default class).
      persistence = {
        for pk, pv in {
          size             = local.storage_size
          storageClassName = local.storage_class != "" ? local.storage_class : null
        } : pk => pv if pv != null
      }

      # Rendered only when declared: a separate volume for snapshots and
      # snapshot shard transfers (upstream sizing guidance: like the main
      # volume).
      snapshotPersistence = local.snapshots_declared ? {
        for sk, sv in {
          enabled          = true
          size             = local.snapshots_size
          storageClassName = local.snapshots_class != "" ? local.snapshots_class : null
        } : sk => sv if sv != null
      } : null

      apiKey         = local.api_key_value_json != null ? jsondecode(local.api_key_value_json) : null
      readOnlyApiKey = local.read_only_api_key_value_json != null ? jsondecode(local.read_only_api_key_value_json) : null

      config = local.tls_enabled ? {
        service = { enable_tls = true }
        tls = {
          cert = "${local.tls_mount_path}/tls.crt"
          key  = "${local.tls_mount_path}/tls.key"
        }
      } : null

      additionalVolumes = local.tls_enabled ? [
        {
          name   = "qdrant-tls"
          secret = { secretName = local.tls_secret }
        }
      ] : null

      additionalVolumeMounts = local.tls_enabled ? [
        {
          name      = "qdrant-tls"
          mountPath = local.tls_mount_path
          readOnly  = true
        }
      ] : null

      nodeSelector = length(try(var.spec.scheduling.node_selector, {})) > 0 ? var.spec.scheduling.node_selector : null

      tolerations = length(try(var.spec.scheduling.tolerations, [])) > 0 ? [
        for t in var.spec.scheduling.tolerations : {
          for tk, tv in {
            key               = t.key != "" ? t.key : null
            operator          = t.operator != "" ? t.operator : null
            value             = t.value != "" ? t.value : null
            effect            = t.effect != "" ? t.effect : null
            tolerationSeconds = try(t.toleration_seconds, null)
          } : tk => tv if tv != null
        }
      ] : null

      # The chart's affinity value is a raw PodSpec affinity object; the
      # typed flag renders the chart's own documented anti-affinity
      # recipe (spread members across nodes so one node loss takes one
      # member, not the quorum).
      affinity = try(var.spec.scheduling.pod_anti_affinity, false) ? {
        podAntiAffinity = {
          requiredDuringSchedulingIgnoredDuringExecution = [
            {
              labelSelector = {
                matchLabels = {
                  "app.kubernetes.io/name"     = local.helm_chart_name
                  "app.kubernetes.io/instance" = local.release_name
                }
              }
              topologyKey = "kubernetes.io/hostname"
            }
          ]
        }
      } : null

      priorityClassName = try(var.spec.scheduling.priority_class_name, "") != "" ? var.spec.scheduling.priority_class_name : null

      metrics = var.spec.service_monitor_enabled ? {
        serviceMonitor = { enabled = true }
      } : null

      # The chart's image.repository carries the registry
      # (docker.io/qdrant/qdrant); useUnprivilegedImage switches to the
      # `qdrant-unprivileged` variant for restricted PSS environments.
      image = (try(var.spec.image, null) != null && (try(var.spec.image.repository, "") != "" || try(var.spec.image.tag, "") != "" || try(var.spec.image.use_unprivileged, false))) ? {
        for ik, iv in {
          repository           = var.spec.image.repository != "" ? var.spec.image.repository : null
          tag                  = var.spec.image.tag != "" ? var.spec.image.tag : null
          useUnprivilegedImage = var.spec.image.use_unprivileged ? true : null
        } : ik => iv if iv != null
      } : null
    } : k => v if v != null && v != {}
  }
}

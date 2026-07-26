# Computed values for the KubernetesSeaweedFs module. Every resolution here
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
  helm_chart_name = "seaweedfs"
  helm_chart_repo = "https://seaweedfs.github.io/seaweedfs/helm"

  # Release name — metadata.name, NOT a fixed chart name: several
  # SeaweedFS stores coexist in one cluster. fullnameOverride below pins
  # every componentName child (`<name>-master`, `-filer`, `-s3`,
  # `-admin`) to this.
  release_name = var.metadata.name

  # Chart version resolved to the pinned default when unset, so both
  # engines install the same chart whether or not the platform's
  # defaulting middleware ran — mirror of the Pulumi module's
  # DefaultChartVersion. Chart 4.40.0 ships appVersion "4.40".
  chart_version = coalesce(var.spec.chart_version, "4.40.0")

  namespace = var.spec.namespace

  # Resource-identity labels stamped on the module-created satellites
  # (namespace, the admin-auth Secret — never injected into the chart's
  # own resources; Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesSeaweedFs"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- s3 posture --------------------------------------------------------
  # The optional bools default to TRUE (this kind IS the catalog's S3
  # store) — only an explicit false turns them off. Dedicated means the
  # gateway runs as its own Deployment; embedded runs it on the filer.
  s3_enabled       = try(coalesce(var.spec.s3.enabled), true)
  s3_auth          = try(coalesce(var.spec.s3.enable_auth), true)
  s3_dedicated     = try(var.spec.s3.dedicated, null) != null
  s3_config_secret = try(var.spec.s3.existing_config_secret, "")
  s3_domain_name   = try(var.spec.s3.domain_name, "")
  s3_buckets       = try(var.spec.s3.buckets, [])

  # Chart-derived service names (componentName = "<fullname>-<suffix>";
  # fullnameOverride pins fullname to the resource name). The `-s3`
  # Service exists for the embedded and dedicated shapes alike.
  master_service_name = "${local.release_name}-master"
  filer_service_name  = "${local.release_name}-filer"
  s3_service_name     = "${local.release_name}-s3"
  admin_service_name  = "${local.release_name}-admin"

  # Name of the Secret carrying the S3 credentials: the chart-generated
  # `<fullname>-s3-secret` (admin + read-only pairs, stable across
  # upgrades, kept on uninstall), the referenced existing config secret,
  # or "" when auth is off.
  s3_credentials_secret_name = (local.s3_enabled && local.s3_auth) ? (
    local.s3_config_secret != "" ? local.s3_config_secret : "${local.release_name}-s3-secret"
  ) : ""

  # ---- admin console --------------------------------------------------------
  admin_enabled       = try(var.spec.admin.enabled, false)
  admin_existing      = try(var.spec.admin.existing_auth_secret, "")
  create_admin_secret = local.admin_enabled && local.admin_existing == ""
  admin_auth_secret_name = local.admin_enabled ? (
    local.create_admin_secret ? "${local.release_name}-admin-auth" : local.admin_existing
  ) : ""

  # ---- per-tier data volumes ---------------------------------------------------
  # STORAGE POSTURE: the chart's out-of-the-box storage is hostPath under
  # /ssd and /storage (bare-metal grain). This module deliberately maps
  # every data volume to a PersistentVolumeClaim and every logs volume to
  # emptyDir — portable across every managed cloud and kind cluster; the
  # escape hatch can restore hostPath for bare-metal fleets. storageClass
  # renders only when declared (absent = the cluster's default class —
  # the chart renders storageClassName empty, which Kubernetes treats as
  # nil).
  master_data = {
    for k, v in {
      type         = "persistentVolumeClaim"
      size         = try(coalesce(var.spec.master.data_volume.size), null) != null && try(var.spec.master.data_volume.size, "") != "" ? var.spec.master.data_volume.size : "5Gi"
      storageClass = try(var.spec.master.data_volume.storage_class, "") != "" ? var.spec.master.data_volume.storage_class : null
    } : k => v if v != null
  }

  volume_data = {
    for k, v in {
      name         = "data"
      type         = "persistentVolumeClaim"
      size         = try(var.spec.volume.data_volume.size, "") != "" ? var.spec.volume.data_volume.size : "30Gi"
      storageClass = try(var.spec.volume.data_volume.storage_class, "") != "" ? var.spec.volume.data_volume.storage_class : null
      maxVolumes   = try(var.spec.volume.max_volumes, 0)
    } : k => v if v != null
  }

  filer_data = {
    for k, v in {
      type         = "persistentVolumeClaim"
      size         = try(var.spec.filer.data_volume.size, "") != "" ? var.spec.filer.data_volume.size : "10Gi"
      storageClass = try(var.spec.filer.data_volume.storage_class, "") != "" ? var.spec.filer.data_volume.storage_class : null
    } : k => v if v != null
  }

  admin_data = try(var.spec.admin.data_volume, null) == null ? null : {
    for k, v in {
      type         = "persistentVolumeClaim"
      size         = var.spec.admin.data_volume.size != "" ? var.spec.admin.data_volume.size : "10Gi"
      storageClass = var.spec.admin.data_volume.storage_class != "" ? var.spec.admin.data_volume.storage_class : null
    } : k => v if v != null
  }

  # ---- container resources (shared shape renderer, per tier) ---------------------
  master_resources = try(var.spec.master.resources, null) == null ? null : {
    for k, v in {
      requests = try(var.spec.master.resources.requests, null) == null ? null : {
        for qk, qv in {
          cpu    = var.spec.master.resources.requests.cpu != "" ? var.spec.master.resources.requests.cpu : null
          memory = var.spec.master.resources.requests.memory != "" ? var.spec.master.resources.requests.memory : null
        } : qk => qv if qv != null
      }
      limits = try(var.spec.master.resources.limits, null) == null ? null : {
        for lk, lv in {
          cpu    = var.spec.master.resources.limits.cpu != "" ? var.spec.master.resources.limits.cpu : null
          memory = var.spec.master.resources.limits.memory != "" ? var.spec.master.resources.limits.memory : null
        } : lk => lv if lv != null
      }
    } : k => v if v != null && v != {}
  }

  volume_resources = try(var.spec.volume.resources, null) == null ? null : {
    for k, v in {
      requests = try(var.spec.volume.resources.requests, null) == null ? null : {
        for qk, qv in {
          cpu    = var.spec.volume.resources.requests.cpu != "" ? var.spec.volume.resources.requests.cpu : null
          memory = var.spec.volume.resources.requests.memory != "" ? var.spec.volume.resources.requests.memory : null
        } : qk => qv if qv != null
      }
      limits = try(var.spec.volume.resources.limits, null) == null ? null : {
        for lk, lv in {
          cpu    = var.spec.volume.resources.limits.cpu != "" ? var.spec.volume.resources.limits.cpu : null
          memory = var.spec.volume.resources.limits.memory != "" ? var.spec.volume.resources.limits.memory : null
        } : lk => lv if lv != null
      }
    } : k => v if v != null && v != {}
  }

  filer_resources = try(var.spec.filer.resources, null) == null ? null : {
    for k, v in {
      requests = try(var.spec.filer.resources.requests, null) == null ? null : {
        for qk, qv in {
          cpu    = var.spec.filer.resources.requests.cpu != "" ? var.spec.filer.resources.requests.cpu : null
          memory = var.spec.filer.resources.requests.memory != "" ? var.spec.filer.resources.requests.memory : null
        } : qk => qv if qv != null
      }
      limits = try(var.spec.filer.resources.limits, null) == null ? null : {
        for lk, lv in {
          cpu    = var.spec.filer.resources.limits.cpu != "" ? var.spec.filer.resources.limits.cpu : null
          memory = var.spec.filer.resources.limits.memory != "" ? var.spec.filer.resources.limits.memory : null
        } : lk => lv if lv != null
      }
    } : k => v if v != null && v != {}
  }

  s3_dedicated_resources = try(var.spec.s3.dedicated.resources, null) == null ? null : {
    for k, v in {
      requests = try(var.spec.s3.dedicated.resources.requests, null) == null ? null : {
        for qk, qv in {
          cpu    = var.spec.s3.dedicated.resources.requests.cpu != "" ? var.spec.s3.dedicated.resources.requests.cpu : null
          memory = var.spec.s3.dedicated.resources.requests.memory != "" ? var.spec.s3.dedicated.resources.requests.memory : null
        } : qk => qv if qv != null
      }
      limits = try(var.spec.s3.dedicated.resources.limits, null) == null ? null : {
        for lk, lv in {
          cpu    = var.spec.s3.dedicated.resources.limits.cpu != "" ? var.spec.s3.dedicated.resources.limits.cpu : null
          memory = var.spec.s3.dedicated.resources.limits.memory != "" ? var.spec.s3.dedicated.resources.limits.memory : null
        } : lk => lv if lv != null
      }
    } : k => v if v != null && v != {}
  }

  admin_resources = try(var.spec.admin.resources, null) == null ? null : {
    for k, v in {
      requests = try(var.spec.admin.resources.requests, null) == null ? null : {
        for qk, qv in {
          cpu    = var.spec.admin.resources.requests.cpu != "" ? var.spec.admin.resources.requests.cpu : null
          memory = var.spec.admin.resources.requests.memory != "" ? var.spec.admin.resources.requests.memory : null
        } : qk => qv if qv != null
      }
      limits = try(var.spec.admin.resources.limits, null) == null ? null : {
        for lk, lv in {
          cpu    = var.spec.admin.resources.limits.cpu != "" ? var.spec.admin.resources.limits.cpu : null
          memory = var.spec.admin.resources.limits.memory != "" ? var.spec.admin.resources.limits.memory : null
        } : lk => lv if lv != null
      }
    } : k => v if v != null && v != {}
  }

  # ---- s3 buckets (the chart's createBuckets shape, hook-consumed) ------------------
  s3_create_buckets = [
    for b in local.s3_buckets : {
      for k, v in {
        name          = b.name
        anonymousRead = b.anonymous_read ? true : null
        ttl           = b.ttl != "" ? b.ttl : null
        objectLock    = b.object_lock ? true : null
        versioning    = b.versioning ? "Enabled" : null
      } : k => v if v != null
    }
  ]

  # The chart's s3-secret and bucket hook read auth and
  # existingConfigSecret from BOTH the s3.* and filer.s3.* paths, so the
  # module renders them on both — only the enabled flags differ between
  # the embedded and dedicated shapes.
  s3_common = {
    for k, v in {
      enableAuth           = local.s3_auth
      existingConfigSecret = local.s3_config_secret != "" ? local.s3_config_secret : null
      domainName           = local.s3_domain_name != "" ? local.s3_domain_name : null
    } : k => v if v != null
  }

  # ---- typed chart values (twin of the Pulumi module's buildHelmValues) --
  helm_values = {
    for k, v in {
      # fullnameOverride pins seaweedfs.fullname to the resource name:
      # the componentName children derive deterministically and the
      # longest suffix stays far from the 63-char ceiling (the chart's
      # componentName helper truncates internally as a second guard).
      fullnameOverride = local.release_name

      # Cross-tier replication: one typed placement code flips the
      # chart's enableReplication and overrides master + filer placement
      # together. ServiceMonitors are gated per component on the ONE
      # monitoring flag.
      global = (var.spec.replication != "" || var.spec.service_monitor_enabled) ? {
        seaweedfs = {
          for gk, gv in {
            enableReplication    = var.spec.replication != "" ? true : null
            replicationPlacement = var.spec.replication != "" ? var.spec.replication : null
            monitoring           = var.spec.service_monitor_enabled ? { enabled = true } : null
          } : gk => gv if gv != null
        }
      } : null

      # The chart's top-level image block wins over the global defaults:
      # image.repository REPLACES the whole image name; image.tag
      # defaults to the chart's appVersion.
      image = (try(var.spec.image, null) != null && (try(var.spec.image.registry, "") != "" || try(var.spec.image.repository, "") != "" || try(var.spec.image.tag, "") != "")) ? {
        for ik, iv in {
          registry   = var.spec.image.registry != "" ? var.spec.image.registry : null
          repository = var.spec.image.repository != "" ? var.spec.image.repository : null
          tag        = var.spec.image.tag != "" ? var.spec.image.tag : null
        } : ik => iv if iv != null
      } : null

      master = {
        for mk, mv in {
          replicas          = try(coalesce(var.spec.master.replicas), null)
          volumeSizeLimitMB = try(coalesce(var.spec.master.volume_size_limit_mb), null)
          data              = local.master_data
          logs              = { type = "emptyDir" }
          resources         = local.master_resources
        } : mk => mv if mv != null
      }

      # dataDirs is a LIST in the chart; the typed surface models the
      # canonical single-PVC entry (named "data"), sized per pod. More
      # exotic layouts (multiple dirs, hostPath fleets) ride helm_values
      # — note lists REPLACE on merge, so an override provides the whole
      # list.
      volume = {
        for vk, vv in {
          replicas            = try(coalesce(var.spec.volume.replicas), null)
          dataDirs            = [local.volume_data]
          logs                = { type = "emptyDir" }
          index               = try(var.spec.volume.index_mode, "") != "" ? var.spec.volume.index_mode : null
          minFreeSpacePercent = try(coalesce(var.spec.volume.min_free_space_percent), null)
          resources           = local.volume_resources
        } : vk => vv if vv != null
      }

      # The filer's embedded leveldb metadata store lives on the data PVC
      # (the chart's WEED_LEVELDB2_ENABLED default) — external shared
      # stores (Postgres/MySQL) ride extra env vars + helm_values.
      # EMBEDDED s3 (default): the gateway serves from the filer pods.
      filer = {
        for fk, fv in {
          replicas             = try(coalesce(var.spec.filer.replicas), null)
          data                 = local.filer_data
          logs                 = { type = "emptyDir" }
          encryptVolumeData    = try(var.spec.filer.encrypt_volume_data, false) ? true : null
          extraEnvironmentVars = length(try(var.spec.filer.extra_environment_vars, {})) > 0 ? var.spec.filer.extra_environment_vars : null
          resources            = local.filer_resources
          s3 = merge(
            { enabled = local.s3_enabled && !local.s3_dedicated },
            local.s3_common,
          )
        } : fk => fv if fv != null
      }

      # DEDICATED s3: its own Deployment that scales independently of
      # metadata. Buckets always render under s3.createBuckets (the
      # hook's preferred path for either shape).
      s3 = merge(
        { enabled = local.s3_enabled && local.s3_dedicated },
        local.s3_common,
        {
          for sk, sv in {
            createBuckets = length(local.s3_create_buckets) > 0 ? local.s3_create_buckets : null
            replicas      = try(coalesce(var.spec.s3.dedicated.replicas), null)
            resources     = local.s3_dedicated_resources
            logs          = local.s3_dedicated ? { type = "emptyDir" } : null
          } : sk => sv if sv != null
        },
      )

      # The console is never installed open: the chart requires
      # userKey/pwKey alongside existingSecret, and the module always
      # points at a credentials Secret — the module-materialized
      # `<name>-admin-auth` or the referenced existing one.
      admin = local.admin_enabled ? {
        for ak, av in {
          enabled = true
          secret = {
            existingSecret = local.admin_auth_secret_name
            userKey        = "user"
            pwKey          = "password"
          }
          data      = local.admin_data
          resources = local.admin_resources
        } : ak => av if av != null
      } : null
    } : k => v if v != null && v != {}
  }
}

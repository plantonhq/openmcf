# Computed values for the KubernetesMysql module. Every resolution here
# has an exact twin in the Pulumi module's locals.go / cluster.go /
# secrets.go — keep them in lockstep: same resource names, same rendered
# CR body, same Secret names/keys/types.
#
# HCL DISCIPLINE (applies to every conditional object in this file):
# conditional entries are written as `key = cond ? value : null` inside ONE
# object literal, pruned with `{ for k, v in {...} : k => v if v != null }`.
# The two tempting alternatives are both broken: `cond ? {...} : {}`
# ternaries fail plan-time type unification when branches carry different
# attributes, and `merge(concat(cond ? [{...}] : [], ...)...)` silently
# UNIFIES primitive-only sibling objects into map(string) — numbers and
# booleans arrive at the API as strings and server-side validation rejects
# the object. The null-prune form preserves every value's type: sizes,
# replica counts, PDB bounds, tolerationSeconds, PITR cadence render as
# YAML numbers; the presence-gated booleans as booleans.
#
# Optional nested blocks are read with try(): HCL's && does NOT
# short-circuit, so chained null checks still dereference the null.
#
# PRESENCE-SENSITIVE KEYS (rendered only when they carry signal, exactly
# like the Pulumi module):
#   - tls: OMITTED when the spec block is absent (operator default =
#     enabled with self-generated certificates).
#   - logcollector: OMITTED when the spec block is absent — the module
#     still renders enabled=true + the pinned image (upstream default).
#   - unsafeFlags: only flags that are true render; the whole block is
#     omitted when none are.
#   - affinity: antiAffinityTopologyKey renders only when the spec sets
#     one (upstream default kubernetes.io/hostname); the literal "none"
#     passes through verbatim — the operator's own OFF switch.

locals {
  cluster_name = var.metadata.name
  namespace    = var.spec.namespace

  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesMysql"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  cr_version               = "1.20.0"
  pxc_default_image        = "percona/percona-xtradb-cluster:8.4.8-8.1"
  haproxy_image            = "percona/haproxy:2.8.18-1"
  proxysql_image           = "percona/proxysql2:2.7.3-1.3"
  logcollector_image       = "percona/fluentbit:5.0.6-1"
  backup_image             = "percona/percona-xtrabackup:8.4.0-5.1"
  version_service_endpoint = "https://check.percona.com"
  upgrade_apply            = "disabled"
  default_update_strategy  = "SmartUpdate"
  default_instances        = 3
  default_proxy_replicas   = 3
  default_pitr_upload_secs = 60

  finalizers            = ["percona.com/delete-pxc-pods-in-order"]
  users_secret_name     = "${local.cluster_name}-secrets"
  is_haproxy            = try(var.spec.proxy.proxysql, null) == null
  primary_service_name  = local.is_haproxy ? "${local.cluster_name}-haproxy" : "${local.cluster_name}-proxysql"
  replicas_service_name = (
    local.is_haproxy &&
    !(try(var.spec.proxy.haproxy.expose_replicas.enabled, null) != null && try(var.spec.proxy.haproxy.expose_replicas.enabled, true) == false)
  ) ? "${local.cluster_name}-haproxy-replicas" : ""
  kube_endpoint         = "${local.primary_service_name}.${local.namespace}.svc.cluster.local:3306"

  backup = try(var.spec.backup, null)

  user_password_secrets = {
    for u in var.spec.users : "${local.cluster_name}-user-${u.name}" => u.password
    if try(u.password, "") != ""
  }

  backup_credential_secrets = local.backup == null ? {} : {
    for s in local.backup.storages : "${local.cluster_name}-backup-${s.name}" => (
      try(s.s3.access_keys, null) != null ? {
        AWS_ACCESS_KEY_ID     = s.s3.access_keys.access_key_id
        AWS_SECRET_ACCESS_KEY = s.s3.access_keys.secret_access_key
        } : try(s.azure, null) != null ? {
        AZURE_STORAGE_ACCOUNT_NAME = s.azure.storage_account
        AZURE_STORAGE_ACCOUNT_KEY  = s.azure.access_key
      } : null
    )
    if try(s.s3.access_keys, null) != null || try(s.azure, null) != null
  }

  unsafe_flags = {
    for k, v in {
      pxcSize           = try(var.spec.unsafe.cluster_size, false) ? true : null
      tls               = try(var.spec.unsafe.tls, false) ? true : null
      proxySize         = try(var.spec.unsafe.proxy_size, false) ? true : null
      backupIfUnhealthy = try(var.spec.unsafe.backup_if_unhealthy, false) ? true : null
    } : k => v if v != null
  }

  tls_body = try(var.spec.tls, null) == null ? null : {
    for k, v in {
      enabled = try(var.spec.tls.enabled, null) == null ? true : var.spec.tls.enabled
      SANs    = length(try(var.spec.tls.sans, [])) > 0 ? var.spec.tls.sans : null
      issuerConf = try(var.spec.tls.issuer, "") != "" ? {
        name  = var.spec.tls.issuer
        kind  = coalesce(try(var.spec.tls.issuer_kind, null), "ClusterIssuer")
        group = "cert-manager.io"
      } : null
    } : k => v if v != null
  }

  pxc_body = {
    for k, v in {
      size  = coalesce(try(var.spec.instances, null), local.default_instances)
      image = try(var.spec.image_name, "") != "" ? var.spec.image_name : local.pxc_default_image

      configuration = try(var.spec.mysql_config, "") != "" ? var.spec.mysql_config : null

      autoRecovery = try(var.spec.auto_recovery, null) == null ? null : (
        var.spec.auto_recovery ? null : false
      )

      resources = try(var.spec.resources, null) == null ? null : {
        for rk, rv in {
          limits = try(var.spec.resources.limits, null) == null ? null : {
            for lk, lv in {
              cpu    = try(var.spec.resources.limits.cpu, "") != "" ? var.spec.resources.limits.cpu : null
              memory = try(var.spec.resources.limits.memory, "") != "" ? var.spec.resources.limits.memory : null
            } : lk => lv if lv != null
          }
          requests = try(var.spec.resources.requests, null) == null ? null : {
            for rk2, rv2 in {
              cpu    = try(var.spec.resources.requests.cpu, "") != "" ? var.spec.resources.requests.cpu : null
              memory = try(var.spec.resources.requests.memory, "") != "" ? var.spec.resources.requests.memory : null
            } : rk2 => rv2 if rv2 != null
          }
        } : rk => rv if rv != null
      }

      volumeSpec = {
        persistentVolumeClaim = {
          for pk, pv in {
            storageClassName = try(var.spec.storage.storage_class, "") != "" ? var.spec.storage.storage_class : null
            resources        = { requests = { storage = var.spec.storage.size } }
          } : pk => pv if pv != null
        }
      }

      affinity = try(var.spec.scheduling.anti_affinity_topology_key, "") != "" ? {
        antiAffinityTopologyKey = var.spec.scheduling.anti_affinity_topology_key
      } : null

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

      podDisruptionBudget = (
        try(var.spec.pod_disruption_budget.max_unavailable, 0) > 0 || try(var.spec.pod_disruption_budget.min_available, 0) > 0
        ) ? {
        for pk, pv in {
          maxUnavailable = try(var.spec.pod_disruption_budget.max_unavailable, 0) > 0 ? var.spec.pod_disruption_budget.max_unavailable : null
          minAvailable   = try(var.spec.pod_disruption_budget.min_available, 0) > 0 ? var.spec.pod_disruption_budget.min_available : null
        } : pk => pv if pv != null
      } : null

      imagePullSecrets = length(var.spec.image_pull_secrets) > 0 ? [
        for s in var.spec.image_pull_secrets : { name = s }
      ] : null
    } : k => v if v != null
  }

  # The disabled arm renders as just { enabled = false } — every other key
  # is nulled (and pruned) when ProxySQL owns the proxy role. The enabled
  # flag itself is the ONE key both arms carry, which keeps this a single
  # null-pruned object instead of an inconsistently-typed ternary.
  haproxy_body = {
    for k, v in {
      enabled = local.is_haproxy
      size    = local.is_haproxy ? coalesce(try(var.spec.proxy.haproxy.replicas, null), local.default_proxy_replicas) : null
      image   = local.is_haproxy ? local.haproxy_image : null

      configuration = local.is_haproxy && try(var.spec.proxy.haproxy.config, "") != "" ? var.spec.proxy.haproxy.config : null

      resources = !local.is_haproxy || try(var.spec.proxy.haproxy.resources, null) == null ? null : {
        for rk, rv in {
          limits = try(var.spec.proxy.haproxy.resources.limits, null) == null ? null : {
            for lk, lv in {
              cpu    = try(var.spec.proxy.haproxy.resources.limits.cpu, "") != "" ? var.spec.proxy.haproxy.resources.limits.cpu : null
              memory = try(var.spec.proxy.haproxy.resources.limits.memory, "") != "" ? var.spec.proxy.haproxy.resources.limits.memory : null
            } : lk => lv if lv != null
          }
          requests = try(var.spec.proxy.haproxy.resources.requests, null) == null ? null : {
            for rk2, rv2 in {
              cpu    = try(var.spec.proxy.haproxy.resources.requests.cpu, "") != "" ? var.spec.proxy.haproxy.resources.requests.cpu : null
              memory = try(var.spec.proxy.haproxy.resources.requests.memory, "") != "" ? var.spec.proxy.haproxy.resources.requests.memory : null
            } : rk2 => rv2 if rv2 != null
          }
        } : rk => rv if rv != null
      }

      exposePrimary = local.is_haproxy && try(var.spec.proxy.haproxy.expose_primary, null) != null ? {
        for ek, ev in {
          type        = coalesce(try(var.spec.proxy.haproxy.expose_primary.type, null), "ClusterIP")
          annotations = length(try(var.spec.proxy.haproxy.expose_primary.annotations, {})) > 0 ? var.spec.proxy.haproxy.expose_primary.annotations : null
        } : ek => ev if ev != null
      } : null

      exposeReplicas = local.is_haproxy && try(var.spec.proxy.haproxy.expose_replicas, null) != null ? {
        for ek, ev in {
          enabled     = try(var.spec.proxy.haproxy.expose_replicas.enabled, null) == null ? true : var.spec.proxy.haproxy.expose_replicas.enabled
          onlyReaders = try(var.spec.proxy.haproxy.expose_replicas.only_readers, false) ? true : null
          type        = coalesce(try(var.spec.proxy.haproxy.expose_replicas.type, null), "ClusterIP")
          annotations = length(try(var.spec.proxy.haproxy.expose_replicas.annotations, {})) > 0 ? var.spec.proxy.haproxy.expose_replicas.annotations : null
        } : ek => ev if ev != null
      } : null

      imagePullSecrets = local.is_haproxy && length(var.spec.image_pull_secrets) > 0 ? [
        for s in var.spec.image_pull_secrets : { name = s }
      ] : null
    } : k => v if v != null
  }

  # Same single-object discipline as haproxy_body: enabled carries the
  # arm choice; every ProxySQL-only key nulls (and prunes) away when
  # HAProxy owns the proxy role.
  proxysql_body = {
    for k, v in {
      enabled = !local.is_haproxy
      size    = !local.is_haproxy ? coalesce(try(var.spec.proxy.proxysql.replicas, null), local.default_proxy_replicas) : null
      image   = !local.is_haproxy ? local.proxysql_image : null

      configuration = !local.is_haproxy && try(var.spec.proxy.proxysql.config, "") != "" ? var.spec.proxy.proxysql.config : null

      resources = local.is_haproxy || try(var.spec.proxy.proxysql.resources, null) == null ? null : {
        for rk, rv in {
          limits = try(var.spec.proxy.proxysql.resources.limits, null) == null ? null : {
            for lk, lv in {
              cpu    = try(var.spec.proxy.proxysql.resources.limits.cpu, "") != "" ? var.spec.proxy.proxysql.resources.limits.cpu : null
              memory = try(var.spec.proxy.proxysql.resources.limits.memory, "") != "" ? var.spec.proxy.proxysql.resources.limits.memory : null
            } : lk => lv if lv != null
          }
          requests = try(var.spec.proxy.proxysql.resources.requests, null) == null ? null : {
            for rk2, rv2 in {
              cpu    = try(var.spec.proxy.proxysql.resources.requests.cpu, "") != "" ? var.spec.proxy.proxysql.resources.requests.cpu : null
              memory = try(var.spec.proxy.proxysql.resources.requests.memory, "") != "" ? var.spec.proxy.proxysql.resources.requests.memory : null
            } : rk2 => rv2 if rv2 != null
          }
        } : rk => rv if rv != null
      }

      volumeSpec = local.is_haproxy ? null : {
        persistentVolumeClaim = {
          for pk, pv in {
            storageClassName = try(var.spec.proxy.proxysql.storage.storage_class, "") != "" ? var.spec.proxy.proxysql.storage.storage_class : null
            resources        = { requests = { storage = var.spec.proxy.proxysql.storage.size } }
          } : pk => pv if pv != null
        }
      }

      expose = !local.is_haproxy && try(var.spec.proxy.proxysql.expose_primary, null) != null ? {
        for ek, ev in {
          type        = coalesce(try(var.spec.proxy.proxysql.expose_primary.type, null), "ClusterIP")
          annotations = length(try(var.spec.proxy.proxysql.expose_primary.annotations, {})) > 0 ? var.spec.proxy.proxysql.expose_primary.annotations : null
        } : ek => ev if ev != null
      } : null

      imagePullSecrets = !local.is_haproxy && length(var.spec.image_pull_secrets) > 0 ? [
        for s in var.spec.image_pull_secrets : { name = s }
      ] : null
    } : k => v if v != null
  }

  users = [
    for u in var.spec.users : {
      for k, v in {
        name  = u.name
        dbs   = length(try(u.dbs, [])) > 0 ? u.dbs : null
        hosts = length(try(u.hosts, [])) > 0 ? u.hosts : null
        grants = length(try(u.grants, [])) > 0 ? u.grants : null
        withGrantOption = try(u.with_grant_option, false) ? true : null
        passwordSecretRef = try(u.password, "") != "" ? {
          name = "${local.cluster_name}-user-${u.name}"
          key  = "password"
        } : null
      } : k => v if v != null
    }
  ]

  backup_storages = local.backup == null ? null : {
    for s in local.backup.storages : s.name => {
      for k, v in {
        type = try(s.s3, null) != null ? "s3" : try(s.azure, null) != null ? "azure" : "filesystem"
        verifyTLS = try(s.verify_tls, null) == null ? null : s.verify_tls

        s3 = try(s.s3, null) == null ? null : {
          for sk, sv in {
            bucket            = s.s3.bucket
            region            = try(s.s3.region, "") != "" ? s.s3.region : null
            prefix            = try(s.s3.prefix, "") != "" ? s.s3.prefix : null
            endpointUrl       = try(s.s3.endpoint_url, "") != "" ? s.s3.endpoint_url : null
            forcePathStyle    = try(s.s3.force_path_style, false) ? true : null
            credentialsSecret = "${local.cluster_name}-backup-${s.name}"
          } : sk => sv if sv != null
        }

        azure = try(s.azure, null) == null ? null : {
          for ak, av in {
            container         = s.azure.container
            prefix            = try(s.azure.prefix, "") != "" ? s.azure.prefix : null
            endpointUrl       = try(s.azure.endpoint_url, "") != "" ? s.azure.endpoint_url : null
            credentialsSecret = "${local.cluster_name}-backup-${s.name}"
          } : ak => av if av != null
        }

        volume = try(s.pvc, null) == null ? null : {
          persistentVolumeClaim = {
            for pk, pv in {
              storageClassName = try(s.pvc.volume.storage_class, "") != "" ? s.pvc.volume.storage_class : null
              accessModes      = ["ReadWriteOnce"]
              resources        = { requests = { storage = s.pvc.volume.size } }
            } : pk => pv if pv != null
          }
        }
      } : k => v if v != null
    }
  }

  backup_schedules = local.backup == null ? null : [
    for sch in local.backup.schedules : {
      for k, v in {
        name        = sch.name
        schedule    = sch.schedule
        storageName = sch.storage_name
        retention = try(sch.keep, null) == null ? null : {
          type              = "count"
          count             = sch.keep
          deleteFromStorage = try(sch.delete_from_storage, null) == null ? true : sch.delete_from_storage
        }
      } : k => v if v != null
    }
  ]

  pitr_body = try(local.backup.pitr, null) == null || !try(local.backup.pitr.enabled, false) ? null : {
    for k, v in {
      enabled            = true
      storageName        = local.backup.pitr.storage_name
      timeBetweenUploads = coalesce(try(local.backup.pitr.time_between_uploads, null), local.default_pitr_upload_secs)
    } : k => v if v != null
  }

  backup_body = local.backup == null ? null : {
    for k, v in {
      image    = local.backup_image
      storages = local.backup_storages
      schedule = length(local.backup_schedules) > 0 ? local.backup_schedules : null
      pitr     = local.pitr_body
      imagePullSecrets = length(var.spec.image_pull_secrets) > 0 ? [
        for s in var.spec.image_pull_secrets : { name = s }
      ] : null
    } : k => v if v != null
  }

  # An absent spec block still renders enabled=true + the pinned image
  # (the upstream default posture) — try() resolves the absent block to
  # the same nulls an empty block carries, so ONE null-pruned object
  # serves both cases.
  logcollector_body = {
    for k, v in {
      enabled = try(var.spec.log_collector.enabled, null) == null ? true : var.spec.log_collector.enabled
      image   = local.logcollector_image
      resources = try(var.spec.log_collector.resources, null) == null ? null : {
        for rk, rv in {
          limits = try(var.spec.log_collector.resources.limits, null) == null ? null : {
            for lk, lv in {
              cpu    = try(var.spec.log_collector.resources.limits.cpu, "") != "" ? var.spec.log_collector.resources.limits.cpu : null
              memory = try(var.spec.log_collector.resources.limits.memory, "") != "" ? var.spec.log_collector.resources.limits.memory : null
            } : lk => lv if lv != null
          }
          requests = try(var.spec.log_collector.resources.requests, null) == null ? null : {
            for rk2, rv2 in {
              cpu    = try(var.spec.log_collector.resources.requests.cpu, "") != "" ? var.spec.log_collector.resources.requests.cpu : null
              memory = try(var.spec.log_collector.resources.requests.memory, "") != "" ? var.spec.log_collector.resources.requests.memory : null
            } : rk2 => rv2 if rv2 != null
          }
        } : rk => rv if rv != null
      }
    } : k => v if v != null
  }

  mysql_manifest = {
    apiVersion = "pxc.percona.com/v1"
    kind       = "PerconaXtraDBCluster"
    metadata = {
      name       = local.cluster_name
      namespace  = local.namespace
      labels     = local.labels
      finalizers = local.finalizers
    }
    spec = {
      for k, v in {
        crVersion    = local.cr_version
        secretsName  = local.users_secret_name
        updateStrategy = coalesce(try(var.spec.update_strategy, null), local.default_update_strategy)
        upgradeOptions = {
          versionServiceEndpoint = local.version_service_endpoint
          apply                  = local.upgrade_apply
        }
        pause = try(var.spec.pause, false) ? true : null

        pxc          = local.pxc_body
        haproxy      = local.haproxy_body
        proxysql     = local.proxysql_body
        tls          = local.tls_body
        users        = length(local.users) > 0 ? local.users : null
        backup       = local.backup_body
        logcollector = local.logcollector_body
        unsafeFlags  = length(local.unsafe_flags) > 0 ? local.unsafe_flags : null
      } : k => v if v != null
    }
  }
}

# Computed values for the KubernetesOpenSearch module. Every resolution
# here has an exact twin in the Pulumi module's locals.go / cluster.go —
# keep them in lockstep: same rendered CR body, same outputs.
#
# HCL DISCIPLINE (applies to every conditional object in this file):
# conditional entries are written as `key = cond ? value : null` inside ONE
# object literal, pruned with `{ for k, v in {...} : k => v if v != null }`.
# The two tempting alternatives are both broken: `cond ? {...} : {}`
# ternaries fail plan-time type unification when branches carry different
# attributes, and `merge(concat(cond ? [{...}] : [], ...)...)` silently
# UNIFIES primitive-only sibling objects into map(string) — numbers and
# booleans arrive at the API as strings and server-side validation rejects
# the object. The null-prune form preserves every value's type: replicas
# and httpPort render as YAML numbers, the presence-gated booleans as
# booleans, PDB bounds as int-or-string per intstr semantics.
#
# Optional nested blocks are read with try(): HCL's && does NOT
# short-circuit, so chained null checks still dereference the null.
#
# PRESENCE-SENSITIVE FIELDS (rendered only on divergence from the CRD/
# operator default, exactly like the Pulumi module):
#   - general.setVMMaxMapCount / tls generate / perNode: CRD default FALSE,
#     spec default TRUE — true renders explicitly, false is omitted.
#   - general.drainDataNodes, nodePools[].pdb.enable, dashboards tls
#     enable/generate, additionalVolumes[].restartPods: plain bools —
#     only true renders.

locals {
  # ClusterName is metadata.name — the naming root the operator derives
  # every object from: StatefulSets `<name>-<pool>`, the main Service
  # `<name>` (this module pins general.serviceName to it), the discovery
  # Service `<name>-discovery`, the admin credentials Secret
  # `<name>-admin-password`, the Dashboards deployment/Service
  # `<name>-dashboards`.
  cluster_name = var.metadata.name

  # image.pull_secret_name joins image_pull_secrets, deduplicated — see
  # the imagePullSecrets rendering below.
  image_pull_secret_names = distinct(concat(
    try(var.spec.image_pull_secrets, []),
    try(coalesce(var.spec.image.pull_secret_name), "") != "" ? [var.spec.image.pull_secret_name] : []
  ))
  namespace    = var.spec.namespace

  # Resource-identity labels stamped on the module-created objects
  # (namespace, OpenSearchCluster). The operator derives ITS objects'
  # identity from the cluster name; these labels tie the family back to
  # the Planton resource.
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesOpenSearch"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  http_port = try(coalesce(var.spec.http_port), 9200)

  # The HTTP endpoint scheme is https REGARDLESS of the security block:
  # the operator itself always talks https to the cluster (upstream
  # pkg/builders/cluster.go URLForCluster returns
  # "https://<svc>.svc.<dns-base>:<port>" unconditionally, and the node
  # readiness probe curls https://localhost:<port>). With spec.security
  # absent the TLS reconciler generates nothing (pkg/reconcilers/tls.go
  # Reconcile returns early) and the opensearchproject image's demo
  # security configuration serves the HTTP layer over TLS instead
  # (pkg/reconcilers/securityconfig.go: "Cluster is running with demo
  # certificates").
  http_endpoint = "https://${local.cluster_name}.${local.namespace}.svc.cluster.local:${local.http_port}"

  # The operator creates `<name>-admin-password` unconditionally
  # (pkg/reconcilers/cluster.go → builders.PasswordSecret), but with a
  # custom security config the operator's bootstrapped credentials do not
  # exist — the user's admin_credentials_secret is authoritative — so the
  # handle is exported empty.
  admin_credentials_secret_name = try(var.spec.security.config, null) != null ? "" : "${local.cluster_name}-admin-password"

  dashboards_enabled = try(var.spec.dashboards.enabled, false)

  # The operator names the Dashboards Service
  # `<general.serviceName>-dashboards` (pkg/builders/dashboards.go) on the
  # fixed port 5601; this module pins serviceName to the cluster name.
  dashboards_service_name = local.dashboards_enabled ? "${local.cluster_name}-dashboards" : ""
  dashboards_endpoint = local.dashboards_enabled ? format(
    "%s://%s.%s.svc.cluster.local:5601",
    try(var.spec.dashboards.tls.enable, false) ? "https" : "http",
    local.dashboards_service_name,
    local.namespace
  ) : ""

  port_forward_command = "kubectl port-forward svc/${local.cluster_name} -n ${local.namespace} ${local.http_port}:${local.http_port}"

  # ---- spec.general -----------------------------------------------------
  # Module-owned constants: serviceName is ALWAYS the cluster name (the
  # operator names the main Service after it — the exported service_name
  # contract depends on this) and vendor is always "opensearch". The CRD's
  # ImageSpec takes ONE image string (not repo/tag fields) — the shared
  # ContainerImage folds into `repo:tag`.
  general_body = {
    for k, v in {
      httpPort    = local.http_port
      serviceName = local.cluster_name
      vendor      = "opensearch"
      version     = var.spec.version

      setVMMaxMapCount = try(coalesce(var.spec.set_vm_max_map_count), true) ? true : null
      drainDataNodes   = try(var.spec.drain_data_nodes, false) ? true : null

      additionalConfig = length(try(var.spec.additional_config, {})) > 0 ? var.spec.additional_config : null
      annotations      = length(try(var.spec.service_annotations, {})) > 0 ? var.spec.service_annotations : null
      pluginsList      = length(try(var.spec.plugins_list, [])) > 0 ? var.spec.plugins_list : null

      image = try(var.spec.image.repo, "") != "" ? (
        try(var.spec.image.tag, "") != "" ? "${var.spec.image.repo}:${var.spec.image.tag}" : var.spec.image.repo
      ) : null
      # image.pull_secret_name joins image_pull_secrets (deduplicated) —
      # a private image override naturally travels with its own pull
      # secret (twin of the Pulumi module).
      imagePullSecrets = length(local.image_pull_secret_names) > 0 ? [
        for s in local.image_pull_secret_names : { name = s }
      ] : null

      keystore = length(try(var.spec.keystore, [])) > 0 ? [
        for entry in var.spec.keystore : {
          for ek, ev in {
            secret      = { name = entry.secret }
            keyMappings = length(try(entry.key_mappings, {})) > 0 ? entry.key_mappings : null
          } : ek => ev if ev != null
        }
      ] : null

      snapshotRepositories = length(try(var.spec.snapshot_repositories, [])) > 0 ? [
        for repo in var.spec.snapshot_repositories : {
          for rk, rv in {
            name     = repo.name
            type     = repo.type
            settings = length(try(repo.settings, {})) > 0 ? repo.settings : null
          } : rk => rv if rv != null
        }
      ] : null

      monitoring = try(var.spec.monitoring.enabled, false) ? {
        for mk, mv in {
          enable               = true
          scrapeInterval       = try(var.spec.monitoring.scrape_interval, "") != "" ? var.spec.monitoring.scrape_interval : null
          monitoringUserSecret = try(var.spec.monitoring.monitoring_user_secret, "") != "" ? var.spec.monitoring.monitoring_user_secret : null
          pluginUrl            = try(var.spec.monitoring.plugin_url, "") != "" ? var.spec.monitoring.plugin_url : null
        } : mk => mv if mv != null
      } : null

      additionalVolumes = length(try(var.spec.additional_volumes, [])) > 0 ? [
        for volume in var.spec.additional_volumes : {
          for vk, vv in {
            name        = volume.name
            path        = volume.path
            subPath     = try(volume.sub_path, "") != "" ? volume.sub_path : null
            secret      = try(volume.secret_name, "") != "" ? { secretName = volume.secret_name } : null
            configMap   = try(volume.config_map_name, "") != "" ? { name = volume.config_map_name } : null
            restartPods = try(volume.restart_pods, false) ? true : null
          } : vk => vv if vv != null
        }
      ] : null
    } : k => v if v != null
  }

  # ---- spec.bootstrap ----------------------------------------------------
  bootstrap_body = try(var.spec.bootstrap, null) == null ? null : {
    for k, v in {
      resources        = local.bootstrap_resources
      jvm              = try(var.spec.bootstrap.jvm, "") != "" ? var.spec.bootstrap.jvm : null
      additionalConfig = length(try(var.spec.bootstrap.additional_config, {})) > 0 ? var.spec.bootstrap.additional_config : null
    } : k => v if v != null
  }

  bootstrap_resources = try(var.spec.bootstrap.resources, null) == null ? null : {
    for k, v in {
      limits = try(var.spec.bootstrap.resources.limits, null) == null ? null : {
        cpu    = var.spec.bootstrap.resources.limits.cpu
        memory = var.spec.bootstrap.resources.limits.memory
      }
      requests = try(var.spec.bootstrap.resources.requests, null) == null ? null : {
        cpu    = var.spec.bootstrap.resources.requests.cpu
        memory = var.spec.bootstrap.resources.requests.memory
      }
    } : k => v if v != null
  }

  # ---- spec.nodePools ----------------------------------------------------
  # One entry per spec.node_pools pool — component is the pool name
  # (StatefulSets become `<cluster>-<pool>`). PDB bounds follow intstr
  # semantics: a string that parses as an integer renders as a YAML number
  # ("2" → 2), anything else (percentages like "25%") as a string — the
  # Pulumi twin applies strconv.Atoi for the identical result.
  node_pools_body = [
    for pool in var.spec.node_pools : {
      for k, v in {
        component = pool.name
        replicas  = pool.replicas
        roles     = pool.roles

        resources = try(pool.resources, null) == null ? null : {
          for rk, rv in {
            limits = try(pool.resources.limits, null) == null ? null : {
              cpu    = pool.resources.limits.cpu
              memory = pool.resources.limits.memory
            }
            requests = try(pool.resources.requests, null) == null ? null : {
              cpu    = pool.resources.requests.cpu
              memory = pool.resources.requests.memory
            }
          } : rk => rv if rv != null
        }

        jvm      = try(pool.jvm, "") != "" ? pool.jvm : null
        diskSize = try(pool.disk_size, "") != "" ? pool.disk_size : null

        # The CRD's PVCSource key is `storageClass` (not storageClassName);
        # accessModes is required by the operator's PVC template — pinned
        # ReadWriteOnce. An emptyDir arm without a size limit renders the
        # empty object (`emptyDir: {}`), same as the typed SDK.
        persistence = try(pool.persistence, null) == null ? null : {
          for pk, pv in {
            pvc = try(pool.persistence.pvc, null) == null ? null : {
              for ck, cv in {
                accessModes  = ["ReadWriteOnce"]
                storageClass = try(pool.persistence.pvc.storage_class, "") != "" ? pool.persistence.pvc.storage_class : null
              } : ck => cv if cv != null
            }
            emptyDir = try(pool.persistence.empty_dir, null) == null ? null : {
              for ck, cv in {
                sizeLimit = try(pool.persistence.empty_dir.size_limit, "") != "" ? pool.persistence.empty_dir.size_limit : null
              } : ck => cv if cv != null
            }
          } : pk => pv if pv != null
        }

        additionalConfig = length(try(pool.additional_config, {})) > 0 ? pool.additional_config : null
        nodeSelector     = length(try(pool.node_selector, {})) > 0 ? pool.node_selector : null

        tolerations = length(try(pool.tolerations, [])) > 0 ? [
          for t in pool.tolerations : {
            for tk, tv in {
              key               = try(t.key, "") != "" ? t.key : null
              operator          = try(t.operator, "") != "" ? t.operator : null
              value             = try(t.value, "") != "" ? t.value : null
              effect            = try(t.effect, "") != "" ? t.effect : null
              tolerationSeconds = try(t.toleration_seconds, null)
            } : tk => tv if tv != null
          }
        ] : null

        # try(tonumber(v), v) — NOT the `can() ? tonumber() : v` ternary:
        # a conditional UNIFIES its number/string branches to string and
        # "2" would render quoted; try() keeps the dynamic type.
        pdb = try(pool.pdb, null) == null ? null : {
          for pk, pv in {
            enable         = try(pool.pdb.enable, false) ? true : null
            minAvailable   = try(pool.pdb.min_available, "") != "" ? try(tonumber(pool.pdb.min_available), pool.pdb.min_available) : null
            maxUnavailable = try(pool.pdb.max_unavailable, "") != "" ? try(tonumber(pool.pdb.max_unavailable), pool.pdb.max_unavailable) : null
          } : pk => pv if pv != null
        }
      } : k => v if v != null
    }
  ]

  # ---- spec.security -----------------------------------------------------
  # Only blocks the spec declares render (null-prune). generate/perNode:
  # CRD default FALSE (existing certificates expected), spec default TRUE
  # — true renders explicitly, false is omitted.
  security_transport = try(var.spec.security.transport_tls, null) == null ? null : {
    for k, v in {
      generate = try(coalesce(var.spec.security.transport_tls.generate), true) ? true : null
      perNode  = try(coalesce(var.spec.security.transport_tls.per_node), true) ? true : null
      secret   = try(var.spec.security.transport_tls.secret, "") != "" ? { name = var.spec.security.transport_tls.secret } : null
      caSecret = try(var.spec.security.transport_tls.ca_secret, "") != "" ? { name = var.spec.security.transport_tls.ca_secret } : null
      nodesDn  = length(try(var.spec.security.transport_tls.nodes_dn, [])) > 0 ? var.spec.security.transport_tls.nodes_dn : null
      adminDn  = length(try(var.spec.security.transport_tls.admin_dn, [])) > 0 ? var.spec.security.transport_tls.admin_dn : null
    } : k => v if v != null
  }

  security_http = try(var.spec.security.http_tls, null) == null ? null : {
    for k, v in {
      generate = try(coalesce(var.spec.security.http_tls.generate), true) ? true : null
      secret   = try(var.spec.security.http_tls.secret, "") != "" ? { name = var.spec.security.http_tls.secret } : null
    } : k => v if v != null
  }

  security_tls = local.security_transport == null && local.security_http == null ? null : {
    for k, v in {
      transport = local.security_transport
      http      = local.security_http
    } : k => v if v != null
  }

  security_config = try(var.spec.security.config, null) == null ? null : {
    for k, v in {
      securityConfigSecret   = try(var.spec.security.config.security_config_secret, "") != "" ? { name = var.spec.security.config.security_config_secret } : null
      adminSecret            = try(var.spec.security.config.admin_secret, "") != "" ? { name = var.spec.security.config.admin_secret } : null
      adminCredentialsSecret = try(var.spec.security.config.admin_credentials_secret, "") != "" ? { name = var.spec.security.config.admin_credentials_secret } : null
    } : k => v if v != null
  }

  security_body = {
    for k, v in {
      tls    = local.security_tls
      config = local.security_config == null ? null : (length(local.security_config) > 0 ? local.security_config : null)
    } : k => v if v != null
  }

  # ---- spec.dashboards ---------------------------------------------------
  # Rendered only when enabled. Version defaults to the CLUSTER version
  # (module-owned: Dashboards refuses mismatched clusters, and the CRD's
  # version field is required). replicas is required by the CRD — always
  # rendered with the resolved default (1).
  dashboards_body = !local.dashboards_enabled ? null : {
    for k, v in {
      enable   = true
      replicas = try(coalesce(var.spec.dashboards.replicas), 1)
      version  = try(var.spec.dashboards.version, "") != "" ? var.spec.dashboards.version : var.spec.version

      resources = try(var.spec.dashboards.resources, null) == null ? null : {
        for rk, rv in {
          limits = try(var.spec.dashboards.resources.limits, null) == null ? null : {
            cpu    = var.spec.dashboards.resources.limits.cpu
            memory = var.spec.dashboards.resources.limits.memory
          }
          requests = try(var.spec.dashboards.resources.requests, null) == null ? null : {
            cpu    = var.spec.dashboards.resources.requests.cpu
            memory = var.spec.dashboards.resources.requests.memory
          }
        } : rk => rv if rv != null
      }

      tls = try(var.spec.dashboards.tls.enable, false) ? {
        for tk, tv in {
          enable   = true
          generate = try(coalesce(var.spec.dashboards.tls.generate), true) ? true : null
          secret   = try(var.spec.dashboards.tls.secret, "") != "" ? { name = var.spec.dashboards.tls.secret } : null
        } : tk => tv if tv != null
      } : null

      basePath                    = try(var.spec.dashboards.base_path, "") != "" ? var.spec.dashboards.base_path : null
      additionalConfig            = length(try(var.spec.dashboards.additional_config, {})) > 0 ? var.spec.dashboards.additional_config : null
      opensearchCredentialsSecret = try(var.spec.dashboards.opensearch_credentials_secret, "") != "" ? { name = var.spec.dashboards.opensearch_credentials_secret } : null

      service = try(var.spec.dashboards.service, null) == null ? null : {
        for sk, sv in {
          type                     = try(coalesce(var.spec.dashboards.service.type), "ClusterIP")
          loadBalancerSourceRanges = length(try(var.spec.dashboards.service.load_balancer_source_ranges, [])) > 0 ? var.spec.dashboards.service.load_balancer_source_ranges : null
        } : sk => sv if sv != null
      }

      pluginsList = length(try(var.spec.dashboards.plugins_list, [])) > 0 ? var.spec.dashboards.plugins_list : null
    } : k => v if v != null
  }

  # ---- the OpenSearchCluster CR -------------------------------------------
  cluster_manifest = {
    apiVersion = "opensearch.opster.io/v1"
    kind       = "OpenSearchCluster"
    metadata = {
      name      = local.cluster_name
      namespace = local.namespace
      labels    = local.labels
    }
    spec = {
      for k, v in {
        general    = local.general_body
        bootstrap  = local.bootstrap_body
        nodePools  = local.node_pools_body
        security   = length(local.security_body) > 0 ? local.security_body : null
        dashboards = local.dashboards_body
      } : k => v if v != null
    }
  }
}

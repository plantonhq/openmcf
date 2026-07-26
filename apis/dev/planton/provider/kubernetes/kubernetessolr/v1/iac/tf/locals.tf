# Computed values for the KubernetesSolr module. Every resolution here has
# an exact twin in the Pulumi module's locals.go / solr_cloud.go — keep them
# in lockstep: same derived names, same rendered CR body.
#
# HCL DISCIPLINE (applies to every conditional object in this file):
# conditional entries are written as `key = cond ? value : null` inside ONE
# object literal, pruned with `{ for k, v in {...} : k => v if v != null }`.
# The two tempting alternatives are both broken: `cond ? {...} : {}`
# ternaries fail plan-time type unification when branches carry different
# attributes, and merging primitive-only sibling objects silently UNIFIES
# them into map(string) — numbers and booleans arrive at the API as strings
# and server-side validation rejects the object. The null-prune form
# preserves every value's type: replicas, podPort, the managed update
# budgets, tolerationSeconds render as YAML numbers; the presence-gated
# booleans as booleans.
#
# Optional nested blocks are read with try(): HCL's && does NOT
# short-circuit, so chained null checks still dereference the null.
#
# INT-OR-STRING: the managed update budgets keep the CRD's int-or-string
# semantics with try(tonumber(x), x) — a ternary would UNIFY the number
# branch into a string ("2" must reach the API as the number 2 while "25%"
# stays a string).
#
# PRESENCE-SENSITIVE FIELDS (rendered only when explicitly set, exactly
# like the Pulumi module):
#   - availability.pdb_enabled: absent already means enabled upstream.
#   - scaling.vacate_pods_on_scale_down / populate_pods_on_scale_up: both
#     default true upstream — absence already means "move replicas".
#   - the plain booleans (probes_require_auth, verify_client_hostname,
#     use_external_address, hide_common, hide_nodes): only true renders.

locals {
  # ClusterName is metadata.name — the naming root the operator derives
  # every object from: StatefulSet `<name>-solrcloud`, common Service
  # `<name>-solrcloud-common`, basic-auth Secret
  # `<name>-solrcloud-basic-auth`, provided ZooKeeper client service
  # `<name>-solrcloud-zookeeper-client`.
  cluster_name = var.metadata.name
  namespace    = var.spec.namespace

  # Resource-identity labels stamped on the module-created objects
  # (namespace, SolrCloud). The operator derives ITS objects' identity
  # from the SolrCloud name; these labels tie the family back to the
  # Planton resource.
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesSolr"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # The operator's naming contract (exported as outputs).
  common_service_name = "${local.cluster_name}-solrcloud-common"

  # TLS drives the endpoint scheme/port: the common service listens on 80
  # without TLS and 443 with TLS (80/443 are the scheme defaults so the
  # endpoint carries no port suffix).
  tls_enabled         = try(var.spec.tls, null) != null
  endpoint_scheme     = local.tls_enabled ? "https" : "http"
  common_service_port = local.tls_enabled ? 443 : 80
  internal_endpoint   = "${local.endpoint_scheme}://${local.common_service_name}.${local.namespace}.svc.cluster.local"

  # Basic auth is the only operator-managed authentication type; the
  # operator generates `<name>-solrcloud-basic-auth` only when the user
  # did not bring their own credential Secret.
  security_basic = try(var.spec.security.authentication_type, "") == "basic"
  basic_auth_secret_name_output = (
    local.security_basic && try(var.spec.security.basic_auth_secret, "") == ""
  ) ? "${local.cluster_name}-solrcloud-basic-auth" : ""

  # ---- ZooKeeper wiring -----------------------------------------------------
  zk_provided = try(var.spec.zookeeper.provided, null)
  zk_external = try(var.spec.zookeeper.external, null)

  # "/" is the shared chroot default of both arms (and of the operator).
  zk_chroot = (
    local.zk_external != null ? try(coalesce(local.zk_external.chroot, "/"), "/") :
    local.zk_provided != null ? try(coalesce(local.zk_provided.chroot, "/"), "/") :
    "/"
  )

  # The connection string the cluster uses. The provided arm (and the
  # empty default, which the operator treats as provided-with-defaults)
  # lands on the operator-named client service; the chroot is appended
  # only when it diverges from "/".
  zk_base_connection = (
    local.zk_external != null
    ? local.zk_external.connection_string
    : "${local.cluster_name}-solrcloud-zookeeper-client:2181"
  )
  zookeeper_connection_string = (
    local.zk_chroot != "" && local.zk_chroot != "/"
    ? "${local.zk_base_connection}${local.zk_chroot}"
    : local.zk_base_connection
  )

  zk_provided_body = local.zk_provided == null ? null : {
    for k, v in {
      replicas = local.zk_provided.replicas
      chroot   = try(coalesce(local.zk_provided.chroot, "/"), "/")

      persistence = try(local.zk_provided.persistence, null) == null ? null : {
        spec = {
          for sk, sv in {
            resources        = try(local.zk_provided.persistence.size, "") != "" ? { requests = { storage = local.zk_provided.persistence.size } } : null
            storageClassName = try(local.zk_provided.persistence.storage_class, "") != "" ? local.zk_provided.persistence.storage_class : null
          } : sk => sv if sv != null
        }
      }

      zookeeperPodPolicy = try(local.zk_provided.resources, null) == null ? null : {
        resources = {
          for rk, rv in {
            limits = try(local.zk_provided.resources.limits, null) == null ? null : {
              cpu    = local.zk_provided.resources.limits.cpu
              memory = local.zk_provided.resources.limits.memory
            }
            requests = try(local.zk_provided.resources.requests, null) == null ? null : {
              cpu    = local.zk_provided.resources.requests.cpu
              memory = local.zk_provided.resources.requests.memory
            }
          } : rk => rv if rv != null
        }
      }
    } : k => v if v != null
  }

  # An empty zookeeper block renders NOTHING — the operator then defaults
  # to a provided 3-node ensemble on its own.
  zookeeper_ref_body = {
    for k, v in {
      provided = local.zk_provided_body
      connectionInfo = local.zk_external == null ? null : {
        internalConnectionString = local.zk_external.connection_string
        chroot                   = try(coalesce(local.zk_external.chroot, "/"), "/")
      }
    } : k => v if v != null
  }

  # ---- data storage -----------------------------------------------------------
  storage_persistent = try(var.spec.storage.persistent, null)
  storage_ephemeral  = try(var.spec.storage.ephemeral, null)

  data_storage_body = {
    for k, v in {
      persistent = local.storage_persistent == null ? null : {
        reclaimPolicy = try(coalesce(local.storage_persistent.reclaim_policy, "Retain"), "Retain")
        pvcTemplate = {
          spec = {
            for sk, sv in {
              resources        = { requests = { storage = local.storage_persistent.size } }
              storageClassName = try(local.storage_persistent.storage_class, "") != "" ? local.storage_persistent.storage_class : null
            } : sk => sv if sv != null
          }
        }
      }
      # emptyDir prunes to {} when no size limit is declared — the
      # rendered CR still carries `emptyDir: {}` (the arm marker).
      ephemeral = local.storage_ephemeral == null ? null : {
        emptyDir = {
          for ek, ev in {
            sizeLimit = try(local.storage_ephemeral.size_limit, "") != "" ? local.storage_ephemeral.size_limit : null
          } : ek => ev if ev != null
        }
      }
    } : k => v if v != null
  }

  # ---- pod options ---------------------------------------------------------------
  # Solr node scheduling/resource knobs, gathered into the operator's
  # customSolrKubeOptions.podOptions — rendered only when any is set.
  pod_options_body = {
    for k, v in {
      resources = try(var.spec.resources, null) == null ? null : {
        for rk, rv in {
          limits = try(var.spec.resources.limits, null) == null ? null : {
            cpu    = var.spec.resources.limits.cpu
            memory = var.spec.resources.limits.memory
          }
          requests = try(var.spec.resources.requests, null) == null ? null : {
            cpu    = var.spec.resources.requests.cpu
            memory = var.spec.resources.requests.memory
          }
        } : rk => rv if rv != null
      }
      nodeSelector = length(var.spec.node_selector) > 0 ? var.spec.node_selector : null
      tolerations = length(var.spec.tolerations) > 0 ? [
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
    } : k => v if v != null
  }

  # ---- addressability --------------------------------------------------------------
  # Always rendered: podPort carries the proto default (8983); the
  # external block models the operator's own Ingress / ExternalDNS
  # exposure when declared (additionalDomainNames deliberately unmodeled).
  external_body = try(var.spec.external, null) == null ? null : {
    for k, v in {
      method             = var.spec.external.method
      domainName         = var.spec.external.domain_name
      useExternalAddress = try(var.spec.external.use_external_address, false) ? true : null
      hideCommon         = try(var.spec.external.hide_common, false) ? true : null
      hideNodes          = try(var.spec.external.hide_nodes, false) ? true : null
    } : k => v if v != null
  }

  solr_addressability_body = {
    for k, v in {
      podPort  = var.spec.pod_port
      external = local.external_body
    } : k => v if v != null
  }

  # ---- update strategy ---------------------------------------------------------------
  # The managed budgets keep the CRD's int-or-string semantics:
  # try(tonumber(x), x) renders "2" as the number 2 and "25%" as a string.
  update_strategy_managed_body = {
    for k, v in {
      maxPodsUnavailable          = try(var.spec.update_strategy.max_pods_unavailable, "") != "" ? try(tonumber(var.spec.update_strategy.max_pods_unavailable), var.spec.update_strategy.max_pods_unavailable) : null
      maxShardReplicasUnavailable = try(var.spec.update_strategy.max_shard_replicas_unavailable, "") != "" ? try(tonumber(var.spec.update_strategy.max_shard_replicas_unavailable), var.spec.update_strategy.max_shard_replicas_unavailable) : null
    } : k => v if v != null
  }

  update_strategy_body = try(var.spec.update_strategy, null) == null ? null : {
    for k, v in {
      method          = try(coalesce(var.spec.update_strategy.method, "Managed"), "Managed")
      managed         = length(local.update_strategy_managed_body) > 0 ? local.update_strategy_managed_body : null
      restartSchedule = try(var.spec.update_strategy.restart_schedule, "") != "" ? var.spec.update_strategy.restart_schedule : null
    } : k => v if v != null
  }

  # ---- availability / scaling -----------------------------------------------------------
  # pdb_enabled is optional-with-default: absence already means enabled to
  # the operator, so only an EXPLICIT value renders.
  availability_body = try(var.spec.availability.pdb_enabled, null) == null ? null : {
    podDisruptionBudget = { enabled = var.spec.availability.pdb_enabled }
  }

  scaling_body = {
    for k, v in {
      vacatePodsOnScaleDown = try(var.spec.scaling.vacate_pods_on_scale_down, null)
      populatePodsOnScaleUp = try(var.spec.scaling.populate_pods_on_scale_up, null)
    } : k => v if v != null
  }

  # ---- TLS / security ---------------------------------------------------------------------
  tls_body = !local.tls_enabled ? null : {
    for k, v in {
      pkcs12Secret = {
        name = var.spec.tls.pkcs12_secret.name
        key  = var.spec.tls.pkcs12_secret.key
      }
      keyStorePasswordSecret = {
        name = var.spec.tls.keystore_password_secret.name
        key  = var.spec.tls.keystore_password_secret.key
      }
      trustStoreSecret = try(var.spec.tls.truststore_secret, null) == null ? null : {
        name = var.spec.tls.truststore_secret.name
        key  = var.spec.tls.truststore_secret.key
      }
      trustStorePasswordSecret = try(var.spec.tls.truststore_password_secret, null) == null ? null : {
        name = var.spec.tls.truststore_password_secret.name
        key  = var.spec.tls.truststore_password_secret.key
      }
      clientAuth           = try(coalesce(var.spec.tls.client_auth, "None"), "None")
      verifyClientHostname = try(var.spec.tls.verify_client_hostname, false) ? true : null
    } : k => v if v != null
  }

  # Rendered only when basic auth is enabled — the CRD's
  # authenticationType value is capitalized ("Basic") while the spec's
  # enum is lowercase. A declared-but-empty security block means security
  # stays disabled and nothing renders.
  security_body = !local.security_basic ? null : {
    for k, v in {
      authenticationType = "Basic"
      basicAuthSecret    = try(var.spec.security.basic_auth_secret, "") != "" ? var.spec.security.basic_auth_secret : null
      probesRequireAuth  = try(var.spec.security.probes_require_auth, false) ? true : null
      bootstrapSecurityJson = try(var.spec.security.bootstrap_security_json, null) == null ? null : {
        name = var.spec.security.bootstrap_security_json.name
        key  = var.spec.security.bootstrap_security_json.key
      }
    } : k => v if v != null
  }

  # ---- backup repositories ----------------------------------------------------------------
  # One entry per declared repository; exactly one backend arm each (the
  # spec's CEL guarantees it). Declared S3 credentials ride secretKeyRefs;
  # an empty credentials block means the nodes' ambient identity (IRSA) —
  # the keyless path — and renders nothing.
  backup_repositories_body = [
    for r in var.spec.backup_repositories : {
      for k, v in {
        name = r.name

        s3 = try(r.s3, null) == null ? null : {
          for sk, sv in {
            region       = r.s3.region
            bucket       = r.s3.bucket
            baseLocation = try(r.s3.base_location, "") != "" ? r.s3.base_location : null
            endpoint     = try(r.s3.endpoint, "") != "" ? r.s3.endpoint : null
            credentials = try(r.s3.credentials, null) == null ? null : {
              for ck, cv in {
                accessKeyIdSecret = try(r.s3.credentials.access_key_id_secret, null) == null ? null : {
                  name = r.s3.credentials.access_key_id_secret.name
                  key  = r.s3.credentials.access_key_id_secret.key
                }
                secretAccessKeySecret = try(r.s3.credentials.secret_access_key_secret, null) == null ? null : {
                  name = r.s3.credentials.secret_access_key_secret.name
                  key  = r.s3.credentials.secret_access_key_secret.key
                }
              } : ck => cv if cv != null
            }
          } : sk => sv if sv != null
        }

        gcs = try(r.gcs, null) == null ? null : {
          for gk, gv in {
            bucket = r.gcs.bucket
            gcsCredentialSecret = try(r.gcs.gcs_credential_secret, null) == null ? null : {
              name = r.gcs.gcs_credential_secret.name
              key  = r.gcs.gcs_credential_secret.key
            }
            baseLocation = try(r.gcs.base_location, "") != "" ? r.gcs.base_location : null
          } : gk => gv if gv != null
        }

        volume = try(r.volume, null) == null ? null : {
          for vk, vv in {
            source = {
              persistentVolumeClaim = { claimName = r.volume.pvc_claim_name }
            }
            directory = try(r.volume.directory, "") != "" ? r.volume.directory : null
          } : vk => vv if vv != null
        }
      } : k => v if v != null
    }
  ]

  # ---- the SolrCloud CR ------------------------------------------------------------------------
  solrcloud_manifest = {
    apiVersion = "solr.apache.org/v1beta1"
    kind       = "SolrCloud"
    metadata = {
      name      = local.cluster_name
      namespace = local.namespace
      labels    = local.labels
    }
    spec = {
      for k, v in {
        replicas = var.spec.replicas
        solrImage = {
          # The tag is the spec's required version; the repository falls
          # back to the official "solr" image. imagePullPolicy is
          # deliberately omitted (operator default).
          repository = try(var.spec.image_repository, "") != "" ? var.spec.image_repository : "solr"
          tag        = var.spec.version
        }
        solrJavaMem  = try(var.spec.java_mem, "") != "" ? var.spec.java_mem : null
        solrOpts     = try(var.spec.solr_opts, "") != "" ? var.spec.solr_opts : null
        solrLogLevel = try(coalesce(var.spec.log_level, "INFO"), "INFO")
        solrGCTune   = try(var.spec.gc_tune, "") != "" ? var.spec.gc_tune : null

        zookeeperRef = length(local.zookeeper_ref_body) > 0 ? local.zookeeper_ref_body : null
        dataStorage  = length(local.data_storage_body) > 0 ? local.data_storage_body : null

        customSolrKubeOptions = length(local.pod_options_body) > 0 ? { podOptions = local.pod_options_body } : null

        solrAddressability = local.solr_addressability_body

        updateStrategy = local.update_strategy_body
        availability   = local.availability_body
        scaling        = length(local.scaling_body) > 0 ? local.scaling_body : null

        solrTLS      = local.tls_body
        solrSecurity = local.security_body

        backupRepositories = length(local.backup_repositories_body) > 0 ? local.backup_repositories_body : null

        solrModules    = length(var.spec.solr_modules) > 0 ? var.spec.solr_modules : null
        additionalLibs = length(var.spec.additional_libs) > 0 ? var.spec.additional_libs : null
      } : k => v if v != null
    }
  }
}

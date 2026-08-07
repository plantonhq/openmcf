# Computed values for the KubernetesKafka module. Every resolution here
# has an exact twin in the Pulumi module's locals.go / kafka.go /
# nodepools.go — keep them in lockstep (same keys rendered and omitted,
# numbers as numbers, booleans as booleans).
#
# HCL DISCIPLINE (applies to every conditional object in this file):
# conditional entries are written as `key = cond ? value : null` inside ONE
# object literal, pruned with `{ for k, v in {...} : k => v if v != null }`.
# The two tempting alternatives are both broken: `cond ? {...} : {}`
# ternaries fail plan-time type unification when branches carry different
# attributes, and `merge(concat(cond ? [{...}] : [], ...)...)` silently
# UNIFIES primitive-only sibling objects into map(string) — numbers and
# booleans arrive in the rendered CR as strings, which server-side apply
# rejects. The null-prune form preserves every value's type.
#
# Optional nested blocks are read with try(): HCL's && does NOT
# short-circuit, so chained null checks still dereference the null. Optional
# SCALARS inside present blocks are read with try(coalesce(x), "") — a
# present-but-null optional passes a bare try() and then poisons string
# templates.

locals {
  namespace = var.spec.namespace

  # Cluster name = metadata.name — the Kafka CR name and the value of the
  # strimzi.io/cluster label on every dependent resource
  # (KafkaNodePool here; KafkaTopic/KafkaUser in their own kinds).
  cluster_name = var.metadata.name

  # Resource-identity labels stamped on every module-created object.
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesKafka"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # Strimzi naming contracts (twins of the Pulumi module's locals).
  bootstrap_service_name      = "${local.cluster_name}-kafka-bootstrap"
  cluster_ca_cert_secret_name = "${local.cluster_name}-cluster-ca-cert"
  metrics_config_map_name     = "${local.cluster_name}-kafka-metrics"

  # The first internal listener (explicit "internal" or unset type, which
  # defaults to internal) supplies the in-cluster bootstrap endpoint.
  # Clusters exposing ONLY external listeners export an empty endpoint —
  # an honest signal that in-cluster clients have no plain path.
  internal_listener_ports = [
    for l in var.spec.listeners : l.port
    if coalesce(try(l.type, null), "internal") == "internal"
  ]
  internal_bootstrap_endpoint = length(local.internal_listener_ports) > 0 ? "${local.bootstrap_service_name}.${local.namespace}.svc.cluster.local:${local.internal_listener_ports[0]}" : ""

  # ---- listener rendering ------------------------------------------------
  # Every listener carries the CRD-required quartet (name/port/type/tls);
  # authentication and configuration render only when declared.
  listeners = [
    for l in var.spec.listeners : {
      for k, v in {
        name = l.name
        port = l.port
        type = coalesce(try(l.type, null), "internal")
        tls  = try(l.tls, false)

        authentication = try(l.authentication, null) == null ? null : {
          for ak, av in {
            type = l.authentication.type
            # custom-type knobs render only on the custom arm.
            sasl           = l.authentication.type == "custom" ? try(l.authentication.sasl, false) : null
            listenerConfig = l.authentication.type == "custom" && length(try(l.authentication.listener_config, {})) > 0 ? l.authentication.listener_config : null
          } : ak => av if av != null
        }

        configuration = local.listener_configurations[l.name]
      } : k => v if v != null
    }
  ]

  # ---- per-listener configuration (staged: build once, prune once) ----------
  # Stage 1: the bootstrap sub-body per listener (null-pruned; {} when the
  # listener declares no bootstrap block or every field is empty). The
  # for-expression evaluates UNCONDITIONALLY — every read is try()-guarded,
  # so an absent bootstrap block yields all-null entries that prune to {}.
  # A `cond ? {} : {for ...}` ternary here would be the
  # inconsistent-conditional-types plan failure.
  listener_bootstrap_bodies = {
    for l in var.spec.listeners : l.name => {
      for bk, bv in {
        host             = try(coalesce(l.configuration.bootstrap.host), "") != "" ? l.configuration.bootstrap.host : null
        annotations      = length(try(l.configuration.bootstrap.annotations, {})) > 0 ? l.configuration.bootstrap.annotations : null
        labels           = length(try(l.configuration.bootstrap.labels, {})) > 0 ? l.configuration.bootstrap.labels : null
        loadBalancerIP   = try(coalesce(l.configuration.bootstrap.load_balancer_ip), "") != "" ? l.configuration.bootstrap.load_balancer_ip : null
        nodePort         = try(l.configuration.bootstrap.node_port, null)
        alternativeNames = length(try(l.configuration.bootstrap.alternative_names, [])) > 0 ? l.configuration.bootstrap.alternative_names : null
      } : bk => bv if bv != null
    }
  }

  # Stage 2: the full configuration body per listener (null-pruned; {}
  # when the listener declares no configuration or every field is empty).
  # Unconditional for the same reason as stage 1.
  listener_configuration_bodies = {
    for l in var.spec.listeners : l.name => {
      for k, v in {
        class                 = try(coalesce(l.configuration.class), "") != "" ? l.configuration.class : null
        externalTrafficPolicy = try(coalesce(l.configuration.external_traffic_policy), "") != "" ? l.configuration.external_traffic_policy : null

        loadBalancerSourceRanges      = length(try(l.configuration.load_balancer_source_ranges, [])) > 0 ? l.configuration.load_balancer_source_ranges : null
        allocateLoadBalancerNodePorts = try(l.configuration.allocate_load_balancer_node_ports, null)
        createBootstrapService        = try(l.configuration.create_bootstrap_service, null)
        useServiceDnsDomain           = try(l.configuration.use_service_dns_domain, false) ? true : null
        maxConnections                = try(l.configuration.max_connections, null)
        maxConnectionCreationRate     = try(l.configuration.max_connection_creation_rate, null)
        preferredNodePortAddressType  = try(coalesce(l.configuration.preferred_node_port_address_type), "") != "" ? l.configuration.preferred_node_port_address_type : null
        publishNotReadyAddresses      = try(l.configuration.publish_not_ready_addresses, false) ? true : null

        # cert-manager writes tls.crt/tls.key; the spec defaults mirror
        # that (twin of the Pulumi module's fallbacks).
        brokerCertChainAndKey = try(l.configuration.broker_cert_chain_and_key, null) == null ? null : {
          secretName  = l.configuration.broker_cert_chain_and_key.secret_name
          certificate = coalesce(try(l.configuration.broker_cert_chain_and_key.certificate, null), "tls.crt")
          key         = coalesce(try(l.configuration.broker_cert_chain_and_key.key, null), "tls.key")
        }

        bootstrap = length(local.listener_bootstrap_bodies[l.name]) > 0 ? local.listener_bootstrap_bodies[l.name] : null

        brokers = length(try(l.configuration.brokers, [])) > 0 ? [
          for b in l.configuration.brokers : {
            for wk, wv in {
              broker         = try(b.broker, 0)
              host           = try(coalesce(b.host), "") != "" ? b.host : null
              advertisedHost = try(coalesce(b.advertised_host), "") != "" ? b.advertised_host : null
              advertisedPort = try(b.advertised_port, null)
              annotations    = length(try(b.annotations, {})) > 0 ? b.annotations : null
              labels         = length(try(b.labels, {})) > 0 ? b.labels : null
              loadBalancerIP = try(coalesce(b.load_balancer_ip), "") != "" ? b.load_balancer_ip : null
              nodePort       = try(b.node_port, null)
            } : wk => wv if wv != null
          }
        ] : null
      } : k => v if v != null
    }
  }

  # Stage 3: null when empty, so the listener body prunes the key — the
  # exact twin of the Pulumi module's nil-when-empty return.
  listener_configurations = {
    for name, body in local.listener_configuration_bodies : name => length(body) > 0 ? body : null
  }

  # ---- kafka block ---------------------------------------------------------
  kafka_body = {
    for k, v in {
      version         = try(coalesce(var.spec.kafka_version), "") != "" ? var.spec.kafka_version : null
      metadataVersion = try(coalesce(var.spec.metadata_version), "") != "" ? var.spec.metadata_version : null
      listeners       = local.listeners
      config          = length(try(var.spec.config, {})) > 0 ? var.spec.config : null

      authorization = try(var.spec.authorization, null) == null ? null : {
        for ak, av in {
          type       = var.spec.authorization.type
          superUsers = length(try(var.spec.authorization.super_users, [])) > 0 ? var.spec.authorization.super_users : null
          # custom-type knobs render only on the custom arm.
          authorizerClass  = var.spec.authorization.type == "custom" ? var.spec.authorization.authorizer_class : null
          supportsAdminApi = var.spec.authorization.type == "custom" && try(var.spec.authorization.supports_admin_api, false) ? true : null
        } : ak => av if av != null
      }

      rack = try(var.spec.rack, null) == null ? null : {
        topologyKey = var.spec.rack.topology_key
      }

      jvmOptions = try(var.spec.jvm, null) == null ? null : (
        length({
          for jk, jv in {
            "-Xms" = try(coalesce(var.spec.jvm.xms), "") != "" ? var.spec.jvm.xms : null
            "-Xmx" = try(coalesce(var.spec.jvm.xmx), "") != "" ? var.spec.jvm.xmx : null
          } : jk => jv if jv != null
        }) > 0 ? {
          for jk, jv in {
            "-Xms" = try(coalesce(var.spec.jvm.xms), "") != "" ? var.spec.jvm.xms : null
            "-Xmx" = try(coalesce(var.spec.jvm.xmx), "") != "" ? var.spec.jvm.xmx : null
          } : jk => jv if jv != null
        } : null
      )

      # The module owns the rules ConfigMap (main.tf); the CR only points
      # at it.
      metricsConfig = try(var.spec.metrics.enabled, false) ? {
        type = "jmxPrometheusExporter"
        valueFrom = {
          configMapKeyRef = {
            name = local.metrics_config_map_name
            key  = "kafka-metrics-config.yml"
          }
        }
      } : null
    } : k => v if v != null
  }

  # ---- entity operator -------------------------------------------------------
  # Each sub-operator renders when enabled (the spec defaults both true).
  # When BOTH are disabled the block is omitted entirely — Strimzi deploys
  # no entity operator pod, and KafkaTopic/KafkaUser declarations for this
  # cluster become inert (the spec comments warn about this).
  topic_operator_enabled = coalesce(try(var.spec.entity_operator.topic_operator_enabled, null), true)
  user_operator_enabled  = coalesce(try(var.spec.entity_operator.user_operator_enabled, null), true)
  entity_operator_body = (local.topic_operator_enabled || local.user_operator_enabled) ? {
    for k, v in {
      topicOperator = local.topic_operator_enabled ? {} : null
      userOperator  = local.user_operator_enabled ? {} : null
    } : k => v if v != null
  } : null

  # ---- resources helper (shared ContainerResources shape) ---------------------
  cruise_control_resources = try(var.spec.cruise_control.resources, null) == null ? null : {
    for k, v in {
      limits = try(var.spec.cruise_control.resources.limits, null) == null ? null : {
        for lk, lv in {
          cpu    = try(var.spec.cruise_control.resources.limits.cpu, "") != "" ? var.spec.cruise_control.resources.limits.cpu : null
          memory = try(var.spec.cruise_control.resources.limits.memory, "") != "" ? var.spec.cruise_control.resources.limits.memory : null
        } : lk => lv if lv != null
      }
      requests = try(var.spec.cruise_control.resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = try(var.spec.cruise_control.resources.requests.cpu, "") != "" ? var.spec.cruise_control.resources.requests.cpu : null
          memory = try(var.spec.cruise_control.resources.requests.memory, "") != "" ? var.spec.cruise_control.resources.requests.memory : null
        } : rk => rv if rv != null
      }
    } : k => v if v != null && length(v) > 0
  }

  kafka_exporter_resources = try(var.spec.kafka_exporter.resources, null) == null ? null : {
    for k, v in {
      limits = try(var.spec.kafka_exporter.resources.limits, null) == null ? null : {
        for lk, lv in {
          cpu    = try(var.spec.kafka_exporter.resources.limits.cpu, "") != "" ? var.spec.kafka_exporter.resources.limits.cpu : null
          memory = try(var.spec.kafka_exporter.resources.limits.memory, "") != "" ? var.spec.kafka_exporter.resources.limits.memory : null
        } : lk => lv if lv != null
      }
      requests = try(var.spec.kafka_exporter.resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = try(var.spec.kafka_exporter.resources.requests.cpu, "") != "" ? var.spec.kafka_exporter.resources.requests.cpu : null
          memory = try(var.spec.kafka_exporter.resources.requests.memory, "") != "" ? var.spec.kafka_exporter.resources.requests.memory : null
        } : rk => rv if rv != null
      }
    } : k => v if v != null && length(v) > 0
  }

  # ---- optional cluster companions ---------------------------------------------
  cruise_control_body = try(var.spec.cruise_control.enabled, false) ? {
    for k, v in {
      config    = length(try(var.spec.cruise_control.config, {})) > 0 ? var.spec.cruise_control.config : null
      resources = local.cruise_control_resources != null && length(local.cruise_control_resources) > 0 ? local.cruise_control_resources : null
      autoRebalance = length(try(var.spec.cruise_control.auto_rebalance_modes, [])) > 0 ? [
        for m in var.spec.cruise_control.auto_rebalance_modes : { mode = m }
      ] : null
    } : k => v if v != null
  } : null

  kafka_exporter_body = try(var.spec.kafka_exporter.enabled, false) ? {
    for k, v in {
      groupRegex = try(coalesce(var.spec.kafka_exporter.group_regex), "") != "" ? var.spec.kafka_exporter.group_regex : null
      topicRegex = try(coalesce(var.spec.kafka_exporter.topic_regex), "") != "" ? var.spec.kafka_exporter.topic_regex : null
      resources  = local.kafka_exporter_resources != null && length(local.kafka_exporter_resources) > 0 ? local.kafka_exporter_resources : null
    } : k => v if v != null
  } : null

  cluster_ca_body = try(var.spec.cluster_ca, null) == null ? null : (
    length({
      for k, v in {
        validityDays = try(var.spec.cluster_ca.validity_days, null)
        renewalDays  = try(var.spec.cluster_ca.renewal_days, null)
      } : k => v if v != null
    }) > 0 ? {
      for k, v in {
        validityDays = try(var.spec.cluster_ca.validity_days, null)
        renewalDays  = try(var.spec.cluster_ca.renewal_days, null)
      } : k => v if v != null
    } : null
  )

  clients_ca_body = try(var.spec.clients_ca, null) == null ? null : (
    length({
      for k, v in {
        validityDays = try(var.spec.clients_ca.validity_days, null)
        renewalDays  = try(var.spec.clients_ca.renewal_days, null)
      } : k => v if v != null
    }) > 0 ? {
      for k, v in {
        validityDays = try(var.spec.clients_ca.validity_days, null)
        renewalDays  = try(var.spec.clients_ca.renewal_days, null)
      } : k => v if v != null
    } : null
  )

  # ---- the Kafka CR manifest -----------------------------------------------------
  kafka_manifest = {
    apiVersion = "kafka.strimzi.io/v1"
    kind       = "Kafka"
    metadata = {
      name      = local.cluster_name
      namespace = local.namespace
      labels    = local.labels
    }
    spec = {
      for k, v in {
        kafka                  = local.kafka_body
        entityOperator         = local.entity_operator_body
        cruiseControl          = local.cruise_control_body
        kafkaExporter          = local.kafka_exporter_body
        clusterCa              = local.cluster_ca_body
        clientsCa              = local.clients_ca_body
        maintenanceTimeWindows = length(try(var.spec.maintenance_time_windows, [])) > 0 ? var.spec.maintenance_time_windows : null
      } : k => v if v != null
    }
  }

  # ---- node pool manifests (keyed by pool name — the import-ID basis) -------------
  # The Strimzi pod template carries affinity and tolerations but NO
  # nodeSelector — a node_selector map therefore translates to a
  # requiredDuringSchedulingIgnoredDuringExecution nodeAffinity with one
  # matchExpressions entry per label, sorted by key (semantically identical
  # for exact-match selection; the Pulumi module renders the same
  # translation).
  node_pool_manifests = {
    for pool in var.spec.node_pools : pool.name => {
      apiVersion = "kafka.strimzi.io/v1"
      kind       = "KafkaNodePool"
      metadata = {
        name      = pool.name
        namespace = local.namespace
        # The binding label rides ON TOP of the identity labels.
        labels = merge(local.labels, { "strimzi.io/cluster" = local.cluster_name })
      }
      spec = {
        for k, v in {
          replicas = pool.replicas
          roles    = pool.roles

          # One literal with per-arm gated keys (a chained per-arm ternary
          # selecting differently shaped objects is the inconsistent-
          # conditional-types plan failure wearing a oneof costume). The
          # resolved type gates which keys survive the prune:
          # ephemeral = type only; persistent-claim = size/class/
          # deleteClaim; jbod = volumes (each volume its own null-pruned
          # literal, at most one carrying the KRaft metadata marker).
          storage = {
            for sk, sv in {
              type = coalesce(try(pool.storage.type, null), "persistent-claim")
              size = coalesce(try(pool.storage.type, null), "persistent-claim") == "persistent-claim" ? pool.storage.size : null
              class = (
                coalesce(try(pool.storage.type, null), "persistent-claim") == "persistent-claim" &&
                try(coalesce(pool.storage.storage_class), "") != ""
              ) ? pool.storage.storage_class : null
              deleteClaim = (
                coalesce(try(pool.storage.type, null), "persistent-claim") == "persistent-claim" &&
                try(pool.storage.delete_claim, false)
              ) ? true : null
              volumes = coalesce(try(pool.storage.type, null), "persistent-claim") == "jbod" ? [
                for volume in pool.storage.volumes : {
                  for vk, vv in {
                    id            = try(volume.id, 0)
                    type          = "persistent-claim"
                    size          = volume.size
                    class         = try(coalesce(volume.storage_class), "") != "" ? volume.storage_class : null
                    deleteClaim   = try(volume.delete_claim, false) ? true : null
                    kraftMetadata = try(volume.kraft_metadata, false) ? "shared" : null
                  } : vk => vv if vv != null
                }
              ] : null
            } : sk => sv if sv != null
          }

          resources = try(pool.resources, null) == null ? null : (
            length({
              for k2, v2 in {
                limits = try(pool.resources.limits, null) == null ? null : {
                  for lk, lv in {
                    cpu    = try(pool.resources.limits.cpu, "") != "" ? pool.resources.limits.cpu : null
                    memory = try(pool.resources.limits.memory, "") != "" ? pool.resources.limits.memory : null
                  } : lk => lv if lv != null
                }
                requests = try(pool.resources.requests, null) == null ? null : {
                  for rk, rv in {
                    cpu    = try(pool.resources.requests.cpu, "") != "" ? pool.resources.requests.cpu : null
                    memory = try(pool.resources.requests.memory, "") != "" ? pool.resources.requests.memory : null
                  } : rk => rv if rv != null
                }
              } : k2 => v2 if v2 != null && length(v2) > 0
            }) > 0 ? {
              for k2, v2 in {
                limits = try(pool.resources.limits, null) == null ? null : {
                  for lk, lv in {
                    cpu    = try(pool.resources.limits.cpu, "") != "" ? pool.resources.limits.cpu : null
                    memory = try(pool.resources.limits.memory, "") != "" ? pool.resources.limits.memory : null
                  } : lk => lv if lv != null
                }
                requests = try(pool.resources.requests, null) == null ? null : {
                  for rk, rv in {
                    cpu    = try(pool.resources.requests.cpu, "") != "" ? pool.resources.requests.cpu : null
                    memory = try(pool.resources.requests.memory, "") != "" ? pool.resources.requests.memory : null
                  } : rk => rv if rv != null
                }
              } : k2 => v2 if v2 != null && length(v2) > 0
            } : null
          )

          template = (length(try(pool.node_selector, {})) > 0 || length(try(pool.tolerations, [])) > 0) ? {
            pod = {
              for pk, pv in {
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
                affinity = length(try(pool.node_selector, {})) > 0 ? {
                  nodeAffinity = {
                    requiredDuringSchedulingIgnoredDuringExecution = {
                      nodeSelectorTerms = [{
                        matchExpressions = [
                          for key in sort(keys(pool.node_selector)) : {
                            key      = key
                            operator = "In"
                            values   = [pool.node_selector[key]]
                          }
                        ]
                      }]
                    }
                  }
                } : null
              } : pk => pv if pv != null
            }
          } : null
        } : k => v if v != null
      }
    }
  }
}

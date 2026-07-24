# Computed values for the KubernetesKafkaMirrorMaker2 module. Every
# resolution here has an exact twin in the Pulumi module's locals.go /
# mirrormaker2.go — keep them in lockstep (same keys rendered and omitted,
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

  # Engine name = metadata.name — the KafkaMirrorMaker2 CR name and the
  # stem of every operator naming contract below.
  mirrormaker_name = var.metadata.name

  # Resource-identity labels stamped on every module-created object.
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesKafkaMirrorMaker2"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # Group-identity fallbacks (twins of the Pulumi module's locals): the
  # spec defaults the alias to "target" and derives the group ID and the
  # three storage topics from metadata.name — one engine identity per
  # resource by construction. These MUST be unique among Connect-protocol
  # workloads sharing the target cluster.
  target_alias         = coalesce(try(var.spec.target.alias, null), "target")
  group_id             = coalesce(try(var.spec.target.group_id, null), local.mirrormaker_name)
  config_storage_topic = coalesce(try(var.spec.target.config_storage_topic, null), "${local.mirrormaker_name}-mirrormaker2-configs")
  status_storage_topic = coalesce(try(var.spec.target.status_storage_topic, null), "${local.mirrormaker_name}-mirrormaker2-status")
  offset_storage_topic = coalesce(try(var.spec.target.offset_storage_topic, null), "${local.mirrormaker_name}-mirrormaker2-offsets")

  # Strimzi naming contracts (twins of the Pulumi module's locals).
  metrics_config_map_name = "${local.mirrormaker_name}-mm2-metrics"
  rest_api_endpoint       = "http://${local.mirrormaker_name}-mirrormaker2-api.${local.namespace}.svc.cluster.local:8083"

  # ---- shared client-connection rendering ----------------------------------
  # The target and every mirror source carry the SAME shared client
  # messages (StrimziKafkaClientTls / StrimziKafkaClientAuthentication);
  # the staged bodies below render them once per cluster — mirror-source
  # bodies keyed by alias (the spec enforces alias uniqueness).
  #
  # TLS: the certificates the CLIENT trusts when verifying brokers, each
  # naming a Secret with either one file (certificate) or a glob (pattern)
  # — the spec enforces exactly one of the two.
  target_tls_body = try(var.spec.target.tls, null) == null ? null : {
    trustedCertificates = [
      for c in var.spec.target.tls.trusted_certificates : {
        for ck, cv in {
          secretName  = c.secret_name
          certificate = try(coalesce(c.certificate), "") != "" ? c.certificate : null
          pattern     = try(coalesce(c.pattern), "") != "" ? c.pattern : null
        } : ck => cv if cv != null
      }
    ]
  }

  mirror_source_tls_bodies = {
    for m in var.spec.mirrors : m.source.alias => try(m.source.tls, null) == null ? null : {
      trustedCertificates = [
        for c in m.source.tls.trusted_certificates : {
          for ck, cv in {
            secretName  = c.secret_name
            certificate = try(coalesce(c.certificate), "") != "" ? c.certificate : null
            pattern     = try(coalesce(c.pattern), "") != "" ? c.pattern : null
          } : ck => cv if cv != null
        }
      ]
    }
  }

  # Authentication: each type carries only its own credential shape (the
  # spec's CEL rules guarantee the referenced fields are present) — tls =
  # the client certificate the workload presents (KubernetesKafkaUser
  # credential Secrets carry user.crt/user.key; the fallbacks mirror the
  # spec defaults), the SASL trio = username + password Secret reference
  # (password key defaults to "password"), custom = bring-your-own
  # mechanism via sasl + config.
  target_authentication_body = try(var.spec.target.authentication, null) == null ? null : {
    for ak, av in {
      type = var.spec.target.authentication.type

      certificateAndKey = var.spec.target.authentication.type == "tls" ? {
        secretName  = var.spec.target.authentication.certificate_and_key.secret_name
        certificate = coalesce(try(var.spec.target.authentication.certificate_and_key.certificate, null), "user.crt")
        key         = coalesce(try(var.spec.target.authentication.certificate_and_key.key, null), "user.key")
      } : null

      username = contains(["scram-sha-512", "scram-sha-256", "plain"], var.spec.target.authentication.type) ? var.spec.target.authentication.username : null
      passwordSecret = contains(["scram-sha-512", "scram-sha-256", "plain"], var.spec.target.authentication.type) ? {
        secretName = var.spec.target.authentication.password_secret.secret_name
        password   = coalesce(try(var.spec.target.authentication.password_secret.password, null), "password")
      } : null

      sasl   = var.spec.target.authentication.type == "custom" ? try(var.spec.target.authentication.sasl, false) : null
      config = var.spec.target.authentication.type == "custom" && length(try(var.spec.target.authentication.config, {})) > 0 ? var.spec.target.authentication.config : null
    } : ak => av if av != null
  }

  mirror_source_authentication_bodies = {
    for m in var.spec.mirrors : m.source.alias => try(m.source.authentication, null) == null ? null : {
      for ak, av in {
        type = m.source.authentication.type

        certificateAndKey = m.source.authentication.type == "tls" ? {
          secretName  = m.source.authentication.certificate_and_key.secret_name
          certificate = coalesce(try(m.source.authentication.certificate_and_key.certificate, null), "user.crt")
          key         = coalesce(try(m.source.authentication.certificate_and_key.key, null), "user.key")
        } : null

        username = contains(["scram-sha-512", "scram-sha-256", "plain"], m.source.authentication.type) ? m.source.authentication.username : null
        passwordSecret = contains(["scram-sha-512", "scram-sha-256", "plain"], m.source.authentication.type) ? {
          secretName = m.source.authentication.password_secret.secret_name
          password   = coalesce(try(m.source.authentication.password_secret.password, null), "password")
        } : null

        sasl   = m.source.authentication.type == "custom" ? try(m.source.authentication.sasl, false) : null
        config = m.source.authentication.type == "custom" && length(try(m.source.authentication.config, {})) > 0 ? m.source.authentication.config : null
      } : ak => av if av != null
    }
  }

  # ---- per-mirror connector tuning -----------------------------------------
  # Each mirror runs a MirrorSourceConnector (records + topic
  # configuration) and a MirrorCheckpointConnector (consumer-group offset
  # translation). replication.policy.class in a connector's config selects
  # topic naming: the default prefixes mirrored topics with the source
  # alias; IdentityReplicationPolicy keeps original names (set the SAME
  # class on both connectors — the spec comments carry the contract).
  mirror_source_connector_bodies = {
    for m in var.spec.mirrors : m.source.alias => try(m.source_connector, null) == null ? null : {
      for ck, cv in {
        tasksMax = try(m.source_connector.tasks_max, null)
        config   = length(try(m.source_connector.config, {})) > 0 ? m.source_connector.config : null
        autoRestart = try(m.source_connector.auto_restart, null) == null ? null : {
          for ak, av in {
            enabled     = try(m.source_connector.auto_restart.enabled, false)
            maxRestarts = try(m.source_connector.auto_restart.max_restarts, null)
          } : ak => av if av != null
        }
      } : ck => cv if cv != null
    }
  }

  mirror_checkpoint_connector_bodies = {
    for m in var.spec.mirrors : m.source.alias => try(m.checkpoint_connector, null) == null ? null : {
      for ck, cv in {
        tasksMax = try(m.checkpoint_connector.tasks_max, null)
        config   = length(try(m.checkpoint_connector.config, {})) > 0 ? m.checkpoint_connector.config : null
        autoRestart = try(m.checkpoint_connector.auto_restart, null) == null ? null : {
          for ak, av in {
            enabled     = try(m.checkpoint_connector.auto_restart.enabled, false)
            maxRestarts = try(m.checkpoint_connector.auto_restart.max_restarts, null)
          } : ak => av if av != null
        }
      } : ck => cv if cv != null
    }
  }

  # ---- target block ---------------------------------------------------------
  target_body = {
    for k, v in {
      alias              = local.target_alias
      bootstrapServers   = var.spec.target.bootstrap_servers
      groupId            = local.group_id
      configStorageTopic = local.config_storage_topic
      statusStorageTopic = local.status_storage_topic
      offsetStorageTopic = local.offset_storage_topic
      tls                = local.target_tls_body
      authentication     = local.target_authentication_body
      config             = length(try(var.spec.target.config, {})) > 0 ? var.spec.target.config : null
    } : k => v if v != null
  }

  # ---- mirrors --------------------------------------------------------------
  mirrors = [
    for m in var.spec.mirrors : {
      for k, v in {
        source = {
          for sk, sv in {
            alias            = m.source.alias
            bootstrapServers = m.source.bootstrap_servers
            tls              = local.mirror_source_tls_bodies[m.source.alias]
            authentication   = local.mirror_source_authentication_bodies[m.source.alias]
            config           = length(try(m.source.config, {})) > 0 ? m.source.config : null
          } : sk => sv if sv != null
        }
        topicsPattern        = try(coalesce(m.topics_pattern), "") != "" ? m.topics_pattern : null
        topicsExcludePattern = try(coalesce(m.topics_exclude_pattern), "") != "" ? m.topics_exclude_pattern : null
        groupsPattern        = try(coalesce(m.groups_pattern), "") != "" ? m.groups_pattern : null
        groupsExcludePattern = try(coalesce(m.groups_exclude_pattern), "") != "" ? m.groups_exclude_pattern : null
        sourceConnector      = local.mirror_source_connector_bodies[m.source.alias]
        checkpointConnector  = local.mirror_checkpoint_connector_bodies[m.source.alias]
      } : k => v if v != null
    }
  ]

  # ---- worker resources / JVM ------------------------------------------------
  resources_body = try(var.spec.resources, null) == null ? null : (
    length({
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
    }) > 0 ? {
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
    } : null
  )

  jvm_options_body = try(var.spec.jvm, null) == null ? null : (
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

  # ---- worker scheduling -------------------------------------------------------
  # The Strimzi pod template carries affinity and tolerations but NO
  # nodeSelector — a node_selector map therefore translates to a
  # requiredDuringSchedulingIgnoredDuringExecution nodeAffinity with one
  # matchExpressions entry per label, sorted by key (semantically identical
  # for exact-match selection; the Pulumi module renders the same
  # translation).
  pod_template_body = (length(try(var.spec.node_selector, {})) > 0 || length(try(var.spec.tolerations, [])) > 0) ? {
    pod = {
      for pk, pv in {
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
        affinity = length(try(var.spec.node_selector, {})) > 0 ? {
          nodeAffinity = {
            requiredDuringSchedulingIgnoredDuringExecution = {
              nodeSelectorTerms = [{
                matchExpressions = [
                  for key in sort(keys(var.spec.node_selector)) : {
                    key      = key
                    operator = "In"
                    values   = [var.spec.node_selector[key]]
                  }
                ]
              }]
            }
          }
        } : null
      } : pk => pv if pv != null
    }
  } : null

  # ---- the KafkaMirrorMaker2 CR manifest ------------------------------------------
  mirrormaker2_manifest = {
    apiVersion = "kafka.strimzi.io/v1"
    kind       = "KafkaMirrorMaker2"
    metadata = {
      name      = local.mirrormaker_name
      namespace = local.namespace
      labels    = local.labels
    }
    spec = {
      for k, v in {
        version = try(coalesce(var.spec.version), "") != "" ? var.spec.version : null

        # Twin of the Pulumi module's fallback: the CRD requires a worker
        # count and the spec defaults to one.
        replicas = coalesce(try(var.spec.replicas, null), 1)

        target  = local.target_body
        mirrors = local.mirrors

        resources  = local.resources_body
        jvmOptions = local.jvm_options_body

        rack = try(var.spec.rack, null) == null ? null : {
          topologyKey = var.spec.rack.topology_key
        }

        # The module owns the rules ConfigMap (main.tf); the CR only
        # points at it.
        metricsConfig = try(var.spec.metrics.enabled, false) ? {
          type = "jmxPrometheusExporter"
          valueFrom = {
            configMapKeyRef = {
              name = local.metrics_config_map_name
              key  = "metrics-config.yml"
            }
          }
        } : null

        template = local.pod_template_body
      } : k => v if v != null
    }
  }
}

# The canonical Strimzi JMX Prometheus Exporter rule set for Kafka
# Connect-protocol workers (MirrorMaker 2 runs the Connect engine), copied
# VERBATIM from the pinned upstream example
# (examples/metrics/kafka-connect-metrics.yaml, ConfigMap key
# metrics-config.yml). The Pulumi module embeds the SAME content
# (metrics_rules.go) — the two must stay byte-identical so both engines
# render one ConfigMap shape.
locals {
  mirrormaker2_metrics_rules = <<-EOT
    # Inspired by kafka-connect rules
    # https://github.com/prometheus/jmx_exporter/blob/master/example_configs/kafka-connect.yml
    # See https://github.com/prometheus/jmx_exporter for more info about JMX Prometheus Exporter metrics
    lowercaseOutputName: true
    lowercaseOutputLabelNames: true
    rules:
    #kafka.connect:type=app-info,client-id="{clientid}"
    #kafka.consumer:type=app-info,client-id="{clientid}"
    #kafka.producer:type=app-info,client-id="{clientid}"
    - pattern: 'kafka.(.+)<type=app-info, client-id=(.+)><>start-time-ms'
      name: kafka_$1_start_time_seconds
      labels:
        clientId: "$2"
      help: "Kafka $1 JMX metric start time seconds"
      type: GAUGE
      valueFactor: 0.001
    - pattern: 'kafka.(.+)<type=app-info, client-id=(.+)><>(commit-id|version): (.+)'
      name: kafka_$1_$3_info
      value: 1
      labels:
        clientId: "$2"
        $3: "$4"
      help: "Kafka $1 JMX metric info version and commit-id"
      type: UNTYPED

    #kafka.consumer:type=consumer-fetch-manager-metrics,client-id="{clientid}",topic="{topic}"", partition="{partition}"
    - pattern: kafka.consumer<type=consumer-fetch-manager-metrics, client-id=(.+), topic=(.+), partition=(.+)><>(.+-total)
      name: kafka_consumer_fetch_manager_$4
      labels:
        clientId: "$1"
        topic: "$2"
        partition: "$3"
      help: "Kafka Consumer JMX metric type consumer-fetch-manager-metrics"
      type: COUNTER
    - pattern: kafka.consumer<type=consumer-fetch-manager-metrics, client-id=(.+), topic=(.+), partition=(.+)><>(compression-rate|.+-avg|.+-replica|.+-lag|.+-lead)
      name: kafka_consumer_fetch_manager_$4
      labels:
        clientId: "$1"
        topic: "$2"
        partition: "$3"
      help: "Kafka Consumer JMX metric type consumer-fetch-manager-metrics"
      type: GAUGE

    #kafka.producer:type=producer-topic-metrics,client-id="{clientid}",topic="{topic}"
    - pattern: kafka.producer<type=producer-topic-metrics, client-id=(.+), topic=(.+)><>(.+-total)
      name: kafka_producer_topic_$3
      labels:
        clientId: "$1"
        topic: "$2"
      help: "Kafka Producer JMX metric type producer-topic-metrics"
      type: COUNTER
    - pattern: kafka.producer<type=producer-topic-metrics, client-id=(.+), topic=(.+)><>(compression-rate|.+-avg|.+rate)
      name: kafka_producer_topic_$3
      labels:
        clientId: "$1"
        topic: "$2"
      help: "Kafka Producer JMX metric type producer-topic-metrics"
      type: GAUGE

    #kafka.connect:type=connect-node-metrics,client-id="{clientid}",node-id="{nodeid}"
    #kafka.consumer:type=consumer-node-metrics,client-id=consumer-1,node-id="{nodeid}"
    - pattern: kafka.(.+)<type=(.+)-metrics, client-id=(.+), node-id=(.+)><>(.+-total)
      name: kafka_$2_$5
      labels:
        clientId: "$3"
        nodeId: "$4"
      help: "Kafka $1 JMX metric type $2"
      type: COUNTER
    - pattern: kafka.(.+)<type=(.+)-metrics, client-id=(.+), node-id=(.+)><>(.+-avg|.+-rate)
      name: kafka_$2_$5
      labels:
        clientId: "$3"
        nodeId: "$4"
      help: "Kafka $1 JMX metric type $2"
      type: GAUGE

    #kafka.connect:type=kafka-metrics-count,client-id="{clientid}"
    #kafka.consumer:type=consumer-fetch-manager-metrics,client-id="{clientid}"
    #kafka.consumer:type=consumer-coordinator-metrics,client-id="{clientid}"
    #kafka.consumer:type=consumer-metrics,client-id="{clientid}"
    - pattern: kafka.(.+)<type=(.+)-metrics, client-id=(.*)><>(.+-total)
      name: kafka_$2_$4
      labels:
        clientId: "$3"
      help: "Kafka $1 JMX metric type $2"
      type: COUNTER
    - pattern: kafka.(.+)<type=(.+)-metrics, client-id=(.*)><>(.+-avg|.+-bytes|.+-count|.+-ratio|.+-age|.+-flight|.+-threads|.+-connectors|.+-tasks|.+-ago)
      name: kafka_$2_$4
      labels:
        clientId: "$3"
      help: "Kafka $1 JMX metric type $2"
      type: GAUGE

    #kafka.connect:type=connector-metrics,connector="{connector}"
    - pattern: 'kafka.connect<type=connector-metrics, connector=(.+)><>(connector-class|connector-type|connector-version|status): (.+)'
      name: kafka_connect_connector_$2
      value: 1
      labels:
        connector: "$1"
        $2: "$3"
      help: "Kafka Connect $2 JMX metric type connector"
      type: GAUGE

    #kafka.connect:type=connector-task-metrics,connector="{connector}",task="{task}<> status"
    - pattern: 'kafka.connect<type=connector-task-metrics, connector=(.+), task=(.+)><>status: ([a-z-]+)'
      name: kafka_connect_connector_task_status
      value: 1
      labels:
        connector: "$1"
        task: "$2"
        status: "$3"
      help: "Kafka Connect JMX Connector task status"
      type: GAUGE

    #kafka.connect:type=task-error-metrics,connector="{connector}",task="{task}"
    #kafka.connect:type=source-task-metrics,connector="{connector}",task="{task}"
    #kafka.connect:type=sink-task-metrics,connector="{connector}",task="{task}"
    #kafka.connect:type=connector-task-metrics,connector="{connector}",task="{task}"
    - pattern: kafka.connect<type=(.+)-metrics, connector=(.+), task=(.+)><>(.+-total)
      name: kafka_connect_$1_$4
      labels:
        connector: "$2"
        task: "$3"
      help: "Kafka Connect JMX metric type $1"
      type: COUNTER
    - pattern: kafka.connect<type=(.+)-metrics, connector=(.+), task=(.+)><>(.+-count|.+-ms|.+-ratio|.+-seq-no|.+-rate|.+-max|.+-avg|.+-failures|.+-requests|.+-timestamp|.+-logged|.+-errors|.+-retries|.+-skipped)
      name: kafka_connect_$1_$4
      labels:
        connector: "$2"
        task: "$3"
      help: "Kafka Connect JMX metric type $1"
      type: GAUGE

    #kafka.connect:type=connect-worker-metrics,connector="{connector}"
    - pattern: kafka.connect<type=connect-worker-metrics, connector=(.+)><>([a-z-]+)
      name: kafka_connect_worker_$2
      labels:
        connector: "$1"
      help: "Kafka Connect JMX metric $1"
      type: GAUGE

    #kafka.connect:type=connect-worker-metrics
    - pattern: kafka.connect<type=connect-worker-metrics><>([a-z-]+-total)
      name: kafka_connect_worker_$1
      help: "Kafka Connect JMX metric worker"
      type: COUNTER
    - pattern: kafka.connect<type=connect-worker-metrics><>([a-z-]+)
      name: kafka_connect_worker_$1
      help: "Kafka Connect JMX metric worker"
      type: GAUGE

    #kafka.connect:type=connect-worker-rebalance-metrics,leader-name|connect-protocol
    - pattern: 'kafka.connect<type=connect-worker-rebalance-metrics><>(leader-name|connect-protocol): (.+)'
      name: kafka_connect_worker_rebalance_$1
      value: 1
      labels:
          $1: "$2"
      help: "Kafka Connect $2 JMX metric type worker rebalance"
      type: UNTYPED

    #kafka.connect:type=connect-worker-rebalance-metrics
    - pattern: kafka.connect<type=connect-worker-rebalance-metrics><>([a-z-]+-total)
      name: kafka_connect_worker_rebalance_$1
      help: "Kafka Connect JMX metric rebalance information"
      type: COUNTER
    - pattern: kafka.connect<type=connect-worker-rebalance-metrics><>([a-z-]+)
      name: kafka_connect_worker_rebalance_$1
      help: "Kafka Connect JMX metric rebalance information"
      type: GAUGE

    #kafka.connect:type=connect-coordinator-metrics
    - pattern: kafka.connect<type=connect-coordinator-metrics><>(assigned-connectors|assigned-tasks)
      name: kafka_connect_coordinator_$1
      help: "Kafka Connect JMX metric assignment information"
      type: GAUGE
  EOT
}

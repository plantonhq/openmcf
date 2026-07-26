# Computed values for the KubernetesKafkaUi module. Every resolution here
# has an exact twin in the Pulumi module's locals.go / helm_release.go /
# secret.go — keep them in lockstep.
#
# HCL DISCIPLINE (applies to every conditional object in this file):
# conditional entries are written as `key = cond ? value : null` inside ONE
# object literal, pruned with `{ for k, v in {...} : k => v if v != null }`.
# The two tempting alternatives are both broken: `cond ? {...} : {}`
# ternaries fail plan-time type unification when branches carry different
# attributes, and `merge(concat(cond ? [{...}] : [], ...)...)` silently
# UNIFIES primitive-only sibling objects into map(string) — numbers and
# booleans arrive in the chart values as strings. The null-prune form
# preserves every value's type. Lists of same-shaped elements assemble with
# concat(); merge() is used ONLY where every value is a string (the Kafka
# properties maps) or every value is the same object shape (the secret
# mappings), where unification is exact.
#
# Optional nested blocks are read with try(): HCL's && does NOT
# short-circuit, so chained null checks still dereference the null.
#
# THE PLACEHOLDER / secretMappings MECHANISM (why no secret ever lands in
# rendered values): the chart writes yamlApplicationConfig verbatim into a
# ConfigMap — anything rendered there is world-readable, lands in Helm
# release history and in Terraform state. So every password position in
# the application config carries a literal ${ENV_VAR} placeholder (written
# $${...} in HCL), and envs.secretMappings wires each env var to a
# Kubernetes Secret key (the chart renders valueFrom.secretKeyRef
# entries). The kafbat UI is a Spring Boot app: Spring resolves the
# ${ENV_VAR} placeholders in the mounted config.yml against the container
# environment at startup — credentials exist only inside the running
# container. Referenced credentials (sasl / schema registry / Connect
# password_secret) map straight to their source Secrets; the one LITERAL
# credential (the console login password) materializes into the
# module-owned "<name>-secrets" Secret (main.tf) and maps from there.
#
# Env var naming is deterministic and index-based so both engines emit
# identical placeholders:
#   KAFKA_CLUSTER_<i>_PASSWORD                  — cluster i sasl
#   KAFKA_CLUSTER_<i>_SCHEMA_REGISTRY_PASSWORD  — cluster i registry
#   KAFKA_CLUSTER_<i>_CONNECT_<j>_PASSWORD      — cluster i, connect j
#   KAFKA_UI_USER_PASSWORD                      — console login

locals {
  # Chart identity — must stay byte-identical with the Pulumi module's
  # vars: cross-engine chart-name drift deploys two different products
  # from one manifest.
  helm_chart_name = "kafka-ui"
  helm_chart_repo = "https://ui.charts.kafbat.io"

  # Chart version resolved to the pinned default when unset, so both
  # engines install the same chart whether or not the platform's
  # defaulting middleware ran — mirror of the Pulumi module's
  # DefaultChartVersion. Verified against the SERVED repository index at
  # https://ui.charts.kafbat.io/index.yaml: kafka-ui 1.6.4 (appVersion
  # v1.5.0). Chart and app versions move separately; the chart pin
  # governs.
  chart_version = coalesce(var.spec.chart_version, "1.6.4")

  # Release name — metadata.name, NOT a fixed chart name: several
  # consoles coexist in one cluster, so each manifest gets its own
  # release.
  release_name = var.metadata.name

  namespace = var.spec.namespace

  # Resource-identity labels stamped on the module-created satellites
  # (namespace, the console Secret — never injected into the chart's own
  # resources; Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesKafkaUi"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # Whether login_form authentication was declared. Gates the console
  # Secret, the LOGIN_FORM values, and the KAFKA_UI_USER_PASSWORD secret
  # mapping.
  auth_enabled = try(var.spec.auth, null) != null

  # Deterministic name of the module-materialized Secret holding every
  # LITERAL credential the spec declares (today: the console login
  # password). Referenced credentials stay in their source Secrets.
  console_secret_name = "${var.metadata.name}-secrets"

  # Key inside the console Secret holding the login password — must match
  # the kubernetes_secret_v1.console data key in main.tf.
  console_password_secret_key = "console-user-password"

  # Console replicas / service exposure, resolved to the spec defaults
  # (1 / ClusterIP / 80) when unset — always rendered so both engines
  # emit the identical service block and the endpoint output is
  # deterministic. Optional scalars arrive as null: try(coalesce(x),
  # null) is the null-safe read (coalesce rejects a lone null; try turns
  # that rejection into the fallback).
  replicas     = coalesce(var.spec.replicas, 1)
  service_type = try(coalesce(var.spec.service_type), null) != null ? var.spec.service_type : "ClusterIP"
  service_port = try(coalesce(var.spec.service_port), null) != null ? var.spec.service_port : 80

  # FULLNAME PINNED TO metadata.name (the catalog's Helm-kind
  # convention): the typed values set fullnameOverride, so the chart
  # names its Deployment/Service exactly after the resource — outputs
  # stay deterministic, multiple consoles per cluster never collide on
  # derived names, and verifiers/exposure kinds address the Service by
  # the resource name alone.
  service_name = local.release_name

  # ---- per-cluster Kafka client properties --------------------------------
  # The user's properties map merged with the module-owned security
  # properties derived from the typed tls/sasl blocks (module-owned
  # entries win — the spec forbids credentials in properties). Every
  # value is a string, so merge() unification is exact here.
  #
  # THE PEM TRUSTSTORE TRICK: Kafka clients since KIP-651 accept
  # ssl.truststore.type=PEM with ssl.truststore.location pointing at a
  # plain PEM certificate file — no JKS/PKCS12 conversion, no truststore
  # password. The module mounts the CA Secret as-is (tls_volumes below)
  # and points the truststore at the mounted key, so a Strimzi
  # cluster-CA Secret works directly.
  #
  # The JAAS line: ScramLoginModule for SCRAM-* mechanisms,
  # PlainLoginModule for PLAIN (the spec CEL rule admits nothing else).
  # Username inline (not sensitive); password is the $${...} placeholder.
  cluster_properties = [
    for i, c in var.spec.clusters : merge(
      c.properties,
      c.tls != null ? {
        "security.protocol"       = c.sasl != null ? "SASL_SSL" : "SSL"
        "ssl.truststore.type"     = "PEM"
        "ssl.truststore.location" = "/etc/kafkaui/cluster-${i}-ca/${coalesce(try(c.tls.ca_certificate, null), "ca.crt")}"
      } : {},
      c.sasl != null ? {
        "security.protocol" = c.tls != null ? "SASL_SSL" : "SASL_PLAINTEXT"
        "sasl.mechanism"    = c.sasl.mechanism
        "sasl.jaas.config"  = "${startswith(c.sasl.mechanism, "SCRAM") ? "org.apache.kafka.common.security.scram.ScramLoginModule" : "org.apache.kafka.common.security.plain.PlainLoginModule"} required username=\"${c.sasl.username}\" password=\"$${KAFKA_CLUSTER_${i}_PASSWORD}\";"
      } : {}
    )
  ]

  # ---- kafka.clusters entries (ClustersProperties shape) ------------------
  # readOnly hides every mutating console action for this cluster (topic
  # create/delete, message produce, config edits) — an app-side switch,
  # not a Kafka ACL: the right posture for production clusters on a
  # shared console. Rendered only when true (the app default is false).
  clusters_rendered = [
    for i, c in var.spec.clusters : {
      for k, v in {
        name             = c.name
        bootstrapServers = c.bootstrap_servers
        readOnly         = c.read_only ? true : null
        properties       = length(local.cluster_properties[i]) > 0 ? local.cluster_properties[i] : null
        schemaRegistry   = try(c.schema_registry.url, null)
        schemaRegistryAuth = c.schema_registry != null && (try(c.schema_registry.username, "") != "" || try(c.schema_registry.password_secret, null) != null) ? {
          for ak, av in {
            username = try(c.schema_registry.username, "") != "" ? c.schema_registry.username : null
            password = try(c.schema_registry.password_secret, null) != null ? "$${KAFKA_CLUSTER_${i}_SCHEMA_REGISTRY_PASSWORD}" : null
          } : ak => av if av != null
        } : null
        kafkaConnect = length(c.kafka_connect) > 0 ? [
          for j, kc in c.kafka_connect : {
            for kk, kv in {
              name     = kc.name
              address  = kc.address
              username = kc.username != "" ? kc.username : null
              password = kc.password_secret != null ? "$${KAFKA_CLUSTER_${i}_CONNECT_${j}_PASSWORD}" : null
            } : kk => kv if kv != null
          }
        ] : null
      } : k => v if v != null
    }
  ]

  # ---- application config (mounted as /kafka-ui/config.yml) ---------------
  # LOGIN_FORM rides Spring Boot's DEFAULT security user — the app
  # (io.kafbat.ui.config.auth.BasicAuthSecurityConfig) registers no user
  # store of its own, so exactly ONE account exists:
  # spring.security.user.name/password — which is why the spec models a
  # single `user` (multi-user login needs LDAP/OAuth2 through
  # helm_values).
  yaml_application_config = {
    for k, v in {
      kafka = { clusters = local.clusters_rendered }
      auth  = { type = local.auth_enabled ? "LOGIN_FORM" : "DISABLED" }
      spring = local.auth_enabled ? {
        security = {
          user = {
            name     = var.spec.auth.user.username
            password = "$${KAFKA_UI_USER_PASSWORD}"
          }
        }
      } : null
    } : k => v if v != null
  }

  # ---- secret mappings (one entry per rendered placeholder) ---------------
  # Every value is the same {name, keyName} shape, so merge() unification
  # is exact. coalesce("", "password") also resolves an explicit empty
  # key to the spec default — same as the Pulumi twin's key == "" check.
  sasl_password_mappings = {
    for i, c in var.spec.clusters :
    "KAFKA_CLUSTER_${i}_PASSWORD" => {
      name    = c.sasl.password_secret.secret_name
      keyName = coalesce(c.sasl.password_secret.key, "password")
    } if c.sasl != null
  }

  schema_registry_password_mappings = {
    for i, c in var.spec.clusters :
    "KAFKA_CLUSTER_${i}_SCHEMA_REGISTRY_PASSWORD" => {
      name    = c.schema_registry.password_secret.secret_name
      keyName = coalesce(c.schema_registry.password_secret.key, "password")
    } if try(c.schema_registry.password_secret, null) != null
  }

  connect_password_mappings = merge({}, [
    for i, c in var.spec.clusters : {
      for j, kc in c.kafka_connect :
      "KAFKA_CLUSTER_${i}_CONNECT_${j}_PASSWORD" => {
        name    = kc.password_secret.secret_name
        keyName = coalesce(kc.password_secret.key, "password")
      } if kc.password_secret != null
    }
  ]...)

  console_password_mappings = local.auth_enabled ? {
    KAFKA_UI_USER_PASSWORD = {
      name    = local.console_secret_name
      keyName = local.console_password_secret_key
    }
  } : {}

  secret_mappings = merge(
    local.sasl_password_mappings,
    local.schema_registry_password_mappings,
    local.connect_password_mappings,
    local.console_password_mappings,
  )

  # ---- CA volumes ----------------------------------------------------------
  # One secret volume per TLS cluster ("cluster-<i>-ca", index-named so
  # entries stay stable and unique across clusters), mounted where the
  # rendered ssl.truststore.location points. The chart passes volumes /
  # volumeMounts through to the Deployment verbatim. The source index i
  # is preserved through the filter, so names line up with the rendered
  # truststore paths even when only some clusters use TLS.
  tls_volumes = [
    for i, c in var.spec.clusters : {
      name   = "cluster-${i}-ca"
      secret = { secretName = c.tls.ca_secret_name }
    } if c.tls != null
  ]

  tls_volume_mounts = [
    for i, c in var.spec.clusters : {
      name      = "cluster-${i}-ca"
      mountPath = "/etc/kafkaui/cluster-${i}-ca"
      readOnly  = true
    } if c.tls != null
  ]

  # ---- container resources (shared ContainerResources shape) --------------
  kafka_ui_resources = try(var.spec.resources, null) == null ? null : {
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

  # ---- typed chart values (twin of the Pulumi module's buildHelmValues) --
  helm_values = {
    for k, v in {
      # Always rendered with resolved defaults so both engines emit the
      # identical documents.
      # Pins the chart's fullname to the resource name (see service_name
      # above).
      fullnameOverride = local.release_name
      replicaCount     = local.replicas
      service = {
        type = local.service_type
        port = local.service_port
      }
      yamlApplicationConfig = local.yaml_application_config

      resources    = local.kafka_ui_resources
      nodeSelector = length(var.spec.node_selector) > 0 ? var.spec.node_selector : null
      tolerations = length(var.spec.tolerations) > 0 ? [
        for t in var.spec.tolerations : {
          for tk, tv in {
            key               = t.key != "" ? t.key : null
            operator          = t.operator != "" ? t.operator : null
            value             = t.value != "" ? t.value : null
            effect            = t.effect != "" ? t.effect : null
            tolerationSeconds = try(t.toleration_seconds, null)
          } : tk => tv if tv != null
        }
      ] : null

      # Chart default is ghcr.io/kafbat/kafka-ui at the chart appVersion —
      # only the registry seam is typed (air-gapped mirrors).
      image = var.spec.image_registry != "" ? { registry = var.spec.image_registry } : null

      envs = length(local.secret_mappings) > 0 ? { secretMappings = local.secret_mappings } : null

      volumes      = length(local.tls_volumes) > 0 ? local.tls_volumes : null
      volumeMounts = length(local.tls_volume_mounts) > 0 ? local.tls_volume_mounts : null
    } : k => v if v != null
  }
}

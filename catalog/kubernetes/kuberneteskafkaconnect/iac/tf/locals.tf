# Computed values for the KubernetesKafkaConnect module. Every resolution
# here has an exact twin in the Pulumi module's locals.go / connect.go —
# keep them in lockstep (same keys rendered and omitted, numbers as
# numbers, booleans as booleans).
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

  # Connect cluster name = metadata.name — the KafkaConnect CR name and
  # the value KubernetesKafkaConnector resources bind to (rendered as
  # their strimzi.io/cluster label).
  connect_name = var.metadata.name

  # Resource-identity labels stamped on every module-created object.
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesKafkaConnect"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # Strimzi naming contracts (twins of the Pulumi module's locals). The
  # REST API endpoint is read-only inspection — connector management is
  # declarative through KubernetesKafkaConnector (the operator reverts
  # REST-API-made changes on annotated clusters).
  rest_api_service_name   = "${local.connect_name}-connect-api"
  rest_api_endpoint       = "http://${local.rest_api_service_name}.${local.namespace}.svc.cluster.local:8083"
  metrics_config_map_name = "${local.connect_name}-connect-metrics"

  # The GROUP IDENTITY quartet defaults from metadata.name — these values
  # MUST be unique per Connect cluster sharing a Kafka cluster (two
  # clusters sharing a group.id or a storage topic corrupt each other's
  # state), and deriving the defaults from the resource name keeps
  # distinct clusters distinct without the author inventing four names.
  group_id             = try(coalesce(var.spec.group_id), "") != "" ? var.spec.group_id : local.connect_name
  config_storage_topic = try(coalesce(var.spec.config_storage_topic), "") != "" ? var.spec.config_storage_topic : "${local.connect_name}-connect-configs"
  status_storage_topic = try(coalesce(var.spec.status_storage_topic), "") != "" ? var.spec.status_storage_topic : "${local.connect_name}-connect-status"
  offset_storage_topic = try(coalesce(var.spec.offset_storage_topic), "") != "" ? var.spec.offset_storage_topic : "${local.connect_name}-connect-offsets"

  # ---- client TLS trust ----------------------------------------------------
  # Presence of spec.tls enables TLS on the Kafka connection. Each trusted
  # certificate names a Secret plus EXACTLY ONE of certificate (a single
  # file) or pattern (a glob over the Secret's files) — the proto's XOR
  # validation guarantees one and only one is set, so the body renders
  # whichever arm is present.
  tls_body = try(var.spec.tls, null) == null ? null : {
    trustedCertificates = [
      for cert in var.spec.tls.trusted_certificates : {
        for k, v in {
          secretName  = cert.secret_name
          certificate = try(coalesce(cert.certificate), "") != "" ? cert.certificate : null
          pattern     = try(coalesce(cert.pattern), "") != "" ? cert.pattern : null
        } : k => v if v != null
      }
    ]
  }

  # ---- client authentication -------------------------------------------------
  # The type renders verbatim; each arm renders ONLY its own credential
  # fields (the proto's CEL rules guarantee the arm's fields are present):
  # tls -> certificateAndKey; scram-sha-512/-256/plain -> username +
  # passwordSecret; custom -> sasl + config. KubernetesKafkaUser credential
  # Secrets carry user.crt/user.key/password — the coalesce fallbacks
  # mirror the spec defaults (twin of the Pulumi module's fallbacks).
  sasl_auth_types = ["scram-sha-512", "scram-sha-256", "plain"]

  authentication_body = try(var.spec.authentication, null) == null ? null : {
    for k, v in {
      type = var.spec.authentication.type

      certificateAndKey = var.spec.authentication.type == "tls" ? {
        secretName  = var.spec.authentication.certificate_and_key.secret_name
        certificate = coalesce(try(var.spec.authentication.certificate_and_key.certificate, null), "user.crt")
        key         = coalesce(try(var.spec.authentication.certificate_and_key.key, null), "user.key")
      } : null

      username = contains(local.sasl_auth_types, var.spec.authentication.type) ? var.spec.authentication.username : null
      passwordSecret = contains(local.sasl_auth_types, var.spec.authentication.type) ? {
        secretName = var.spec.authentication.password_secret.secret_name
        password   = coalesce(try(var.spec.authentication.password_secret.password, null), "password")
      } : null

      # custom-type knobs render only on the custom arm; sasl renders even
      # when false — false is the declared value, not an absent one.
      sasl   = var.spec.authentication.type == "custom" ? try(var.spec.authentication.sasl, false) : null
      config = var.spec.authentication.type == "custom" && length(try(var.spec.authentication.config, {})) > 0 ? var.spec.authentication.config : null
    } : k => v if v != null
  }

  # ---- OCI image-volume plugins ------------------------------------------------
  # Plugins mounted from OCI artifacts as Kubernetes image volumes. The
  # artifact type is ALWAYS the literal "image" — it is the only artifact
  # type this arm supports, so the module owns it rather than asking the
  # author to repeat it.
  oci_plugins = length(try(var.spec.plugins, [])) > 0 ? [
    for plugin in var.spec.plugins : {
      name = plugin.name
      artifacts = [
        for artifact in plugin.artifacts : {
          for k, v in {
            type       = "image"
            reference  = artifact.reference
            pullPolicy = try(coalesce(artifact.pull_policy), "") != "" ? artifact.pull_policy : null
          } : k => v if v != null
        }
      ]
    }
  ] : null

  # ---- operator-driven image build -----------------------------------------------
  # output names the registry destination (type defaults to docker — the
  # Kubernetes path; imagestream is OpenShift-only); plugins declare the
  # artifacts baked into the image. Each artifact renders exactly its
  # type's fields — a url-family artifact (jar/tgz/zip/other) never
  # carries Maven coordinates and a maven artifact never carries a url
  # (the proto's CEL rules enforce the same partition), so the operator's
  # schema validation sees only the keys its per-type sub-schema allows.
  url_artifact_types = ["jar", "tgz", "zip", "other"]

  build_body = try(var.spec.build, null) == null ? null : {
    output = {
      for k, v in {
        type                   = coalesce(try(var.spec.build.output.type, null), "docker")
        image                  = var.spec.build.output.image
        pushSecret             = try(coalesce(var.spec.build.output.push_secret), "") != "" ? var.spec.build.output.push_secret : null
        additionalBuildOptions = length(try(var.spec.build.output.additional_build_options, [])) > 0 ? var.spec.build.output.additional_build_options : null
        additionalPushOptions  = length(try(var.spec.build.output.additional_push_options, [])) > 0 ? var.spec.build.output.additional_push_options : null
      } : k => v if v != null
    }
    plugins = [
      for plugin in var.spec.build.plugins : {
        name = plugin.name
        artifacts = [
          for artifact in plugin.artifacts : {
            for k, v in {
              type       = artifact.type
              url        = contains(local.url_artifact_types, artifact.type) ? artifact.url : null
              sha512sum  = contains(local.url_artifact_types, artifact.type) && try(coalesce(artifact.sha512sum), "") != "" ? artifact.sha512sum : null
              insecure   = contains(local.url_artifact_types, artifact.type) && try(artifact.insecure, false) ? true : null
              fileName   = artifact.type == "other" && try(coalesce(artifact.file_name), "") != "" ? artifact.file_name : null
              repository = artifact.type == "maven" && try(coalesce(artifact.repository), "") != "" ? artifact.repository : null
              group      = artifact.type == "maven" ? artifact.group : null
              artifact   = artifact.type == "maven" ? artifact.artifact : null
              version    = artifact.type == "maven" ? artifact.version : null
            } : k => v if v != null
          }
        ]
      }
    ]
  }

  # ---- worker resources (shared ContainerResources shape) -----------------------
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

  # ---- JVM heap --------------------------------------------------------------
  jvm_options = try(var.spec.jvm, null) == null ? null : (
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

  # ---- pod template ------------------------------------------------------------
  # The Strimzi pod template carries affinity and tolerations but NO
  # nodeSelector — a node_selector map therefore translates to a
  # requiredDuringSchedulingIgnoredDuringExecution nodeAffinity with one
  # matchExpressions entry per label, sorted by key (semantically identical
  # for exact-match selection; the Pulumi module renders the same
  # translation).
  template_body = (length(try(var.spec.node_selector, {})) > 0 || length(try(var.spec.tolerations, [])) > 0) ? {
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

  # ---- the KafkaConnect CR manifest ------------------------------------------------
  connect_manifest = {
    apiVersion = "kafka.strimzi.io/v1"
    kind       = "KafkaConnect"
    metadata = {
      name      = local.connect_name
      namespace = local.namespace
      labels    = local.labels
      # Module-owned and unconditional: this annotation is what makes
      # KubernetesKafkaConnector declarations work — the operator
      # reconciles KafkaConnector resources against this cluster and
      # reverts REST-API-made changes it does not own.
      annotations = {
        "strimzi.io/use-connector-resources" = "true"
      }
    }
    spec = {
      for k, v in {
        version  = try(coalesce(var.spec.version), "") != "" ? var.spec.version : null
        replicas = var.spec.replicas
        # When build is also configured the operator runs the image it
        # builds and ignores this — the spec documents set-one-or-the-
        # other; the module renders what was declared.
        image              = try(coalesce(var.spec.image), "") != "" ? var.spec.image : null
        bootstrapServers   = var.spec.bootstrap_servers
        groupId            = local.group_id
        configStorageTopic = local.config_storage_topic
        statusStorageTopic = local.status_storage_topic
        offsetStorageTopic = local.offset_storage_topic
        tls                = local.tls_body
        authentication     = local.authentication_body
        config             = length(try(var.spec.config, {})) > 0 ? var.spec.config : null
        plugins            = local.oci_plugins
        build              = local.build_body
        resources          = local.resources_body
        jvmOptions         = local.jvm_options

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

        template = local.template_body
      } : k => v if v != null
    }
  }
}

# KubernetesJupyterHub locals — every resolution here has an exact twin in
# the Pulumi module's locals.go/values.go; keep them in lockstep.

locals {
  # Chart identity — must stay byte-identical with the Pulumi module's
  # vars (cross-engine chart-name drift deploys two different products
  # from one manifest).
  helm_chart_name = "jupyterhub"
  helm_chart_repo = "https://hub.jupyter.org/helm-chart/"

  release_name  = var.metadata.name
  namespace     = var.spec.namespace
  chart_version = try(var.spec.chart_version, "") != "" ? var.spec.chart_version : "4.4.0"

  # Resource-identity labels stamped on the module-created satellites
  # (namespace + module-owned Secrets — never injected into the chart's
  # own resources; Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesJupyterHub"
    },
    try(var.metadata.id, "") != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    try(var.metadata.org, "") != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    try(var.metadata.env, "") != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ------------------------------ database ------------------------------
  # The chart's own hub.db vocabulary: sqlite-pvc (default) / postgres /
  # mysql. External arms carry a CREDENTIAL-FREE url — the hub exports
  # PGPASSWORD/MYSQL_PWD from the `hub.db.password` key of the mounted
  # existing-secret at startup (chart truth: jupyterhub_config.py).
  hub_postgres = try(var.spec.hub.database.postgres, null)
  hub_mysql    = try(var.spec.hub.database.mysql, null)
  hub_sqlite   = try(var.spec.hub.database.sqlite_pvc, null)

  db_type = local.hub_postgres != null ? "postgres" : (local.hub_mysql != null ? "mysql" : "sqlite-pvc")

  db_port_postgres = try(local.hub_postgres.port, null) != null ? local.hub_postgres.port : 5432
  db_port_mysql    = try(local.hub_mysql.port, null) != null ? local.hub_mysql.port : 3306

  db_url = local.db_type == "postgres" ? format(
    "postgresql+psycopg2://%s@%s:%d/%s",
    try(local.hub_postgres.username, "") != "" ? local.hub_postgres.username : "jupyterhub",
    local.hub_postgres.host,
    local.db_port_postgres,
    try(local.hub_postgres.database_name, "") != "" ? local.hub_postgres.database_name : "jupyterhub"
    ) : (local.db_type == "mysql" ? format(
      "mysql+pymysql://%s@%s:%d/%s",
      try(local.hub_mysql.username, "") != "" ? local.hub_mysql.username : "jupyterhub",
      local.hub_mysql.host,
      local.db_port_mysql,
      try(local.hub_mysql.database_name, "") != "" ? local.hub_mysql.database_name : "jupyterhub"
  ) : "")

  db_password_secret = local.db_type == "postgres" ? local.hub_postgres.password_secret.secret_name : (
    local.db_type == "mysql" ? local.hub_mysql.password_secret.secret_name : ""
  )
  db_password_secret_key_raw = local.db_type == "postgres" ? try(local.hub_postgres.password_secret.secret_key, "") : (
    local.db_type == "mysql" ? try(local.hub_mysql.password_secret.secret_key, "") : ""
  )
  db_password_secret_key = local.db_password_secret_key_raw != "" && local.db_password_secret_key_raw != null ? local.db_password_secret_key_raw : "password"

  # Module-owned hub Secret (`<name>-hub-secret`) — the chart's
  # hub.existingSecret seam; exists only on the external database arms.
  hub_secret_enabled = local.db_type != "sqlite-pvc"
  hub_secret_name    = "${local.release_name}-hub-secret"

  # ---------------------------- authentication --------------------------
  # The secured default: an absent authentication block is the
  # shared-password arm with a module-generated password — the chart's
  # own default (any username, NO password) never ships.
  auth        = try(var.spec.authentication, null)
  auth_method = try(local.auth.native, null) != null ? "native" : (
    try(local.auth.github, null) != null ? "github" : (
      try(local.auth.google, null) != null ? "google" : (
        try(local.auth.oidc, null) != null ? "oidc" : "shared_password"
      )
    )
  )

  shared_password_byo = try(local.auth.shared_password.password_secret, null) != null
  shared_password_secret_name = local.auth_method == "shared_password" ? (
    local.shared_password_byo ? local.auth.shared_password.password_secret.secret_name : "${local.release_name}-auth"
  ) : ""
  shared_password_secret_key_raw = local.shared_password_byo ? try(local.auth.shared_password.password_secret.secret_key, "") : "password"
  shared_password_secret_key    = local.shared_password_secret_key_raw != "" && local.shared_password_secret_key_raw != null ? local.shared_password_secret_key_raw : "password"
  shared_password_module_owned  = local.auth_method == "shared_password" && !local.shared_password_byo

  oauth_client_secret_ref = local.auth_method == "github" ? local.auth.github.client_secret_secret : (
    local.auth_method == "google" ? local.auth.google.client_secret_secret : (
      local.auth_method == "oidc" ? local.auth.oidc.client_secret_secret : null
    )
  )
  oauth_client_secret_key = local.oauth_client_secret_ref != null ? (
    try(local.oauth_client_secret_ref.secret_key, "") != "" && try(local.oauth_client_secret_ref.secret_key, "") != null ? local.oauth_client_secret_ref.secret_key : "password"
  ) : "password"

  # Environment variable names the extraConfig snippets read — the
  # leak-free path for secret material (identical constants in the
  # Pulumi module's vars).
  shared_password_env_var    = "PLANTON_SHARED_PASSWORD"
  oauth_client_secret_env_var = "PLANTON_OAUTH_CLIENT_SECRET"

  # Roster: JupyterHub 5 denies sign-in unless an allow rule matches —
  # an empty roster means "any authenticated identity", declared
  # EXPLICITLY via allow_all.
  allowed_users = try(local.auth.allowed_users, [])
  admin_users   = try(local.auth.admin_users, [])
  # merge-of-conditionals (each branch against {}) — HCL rejects a
  # single ternary whose branches carry different object attributes.
  authenticator_config = merge(
    length(local.allowed_users) > 0 ? { allowed_users = local.allowed_users } : {},
    length(local.allowed_users) == 0 ? { allow_all = true } : {},
    length(local.admin_users) > 0 ? { admin_users = local.admin_users } : {}
  )

  # hub.config per arm — PUBLIC identity settings only; secrets ride env
  # indirection. Full class paths: entry-point shortnames are
  # registration details, the class path is unambiguous at any version.
  oidc_scopes         = length(try(local.auth.oidc.scopes, [])) > 0 ? local.auth.oidc.scopes : ["openid", "email", "profile"]
  oidc_username_claim = try(local.auth.oidc.username_claim, "") != "" && try(local.auth.oidc.username_claim, "") != null ? local.auth.oidc.username_claim : "preferred_username"
  oidc_login_service  = try(local.auth.oidc.login_service, "") != "" && try(local.auth.oidc.login_service, "") != null ? local.auth.oidc.login_service : "OIDC"

  # Every ternary below carries ONE attribute — a branch bundling
  # differently-typed attributes against {} cannot unify (the HCL
  # type-unification class).
  authenticator_class_by_method = {
    shared_password = "jupyterhub.auth.DummyAuthenticator"
    native          = "nativeauthenticator.NativeAuthenticator"
    github          = "oauthenticator.github.GitHubOAuthenticator"
    google          = "oauthenticator.google.GoogleOAuthenticator"
    oidc            = "oauthenticator.generic.GenericOAuthenticator"
  }

  hub_config = merge(
    {
      Authenticator = local.authenticator_config
      JupyterHub    = { authenticator_class = local.authenticator_class_by_method[local.auth_method] }
    },
    local.auth_method == "native" && length(local.native_config) > 0 ? { NativeAuthenticator = local.native_config } : {},
    local.auth_method == "github" ? {
      GitHubOAuthenticator = merge(
        {
          client_id          = local.auth.github.client_id
          oauth_callback_url = local.auth.github.oauth_callback_url
        },
        # Org/team membership checks need the read:org scope; the two
        # attributes ride separate ternaries (list vs list unify, but
        # keeping the discipline uniform costs nothing).
        length(try(local.auth.github.allowed_organizations, [])) > 0 ? { allowed_organizations = local.auth.github.allowed_organizations } : {},
        length(try(local.auth.github.allowed_organizations, [])) > 0 ? { scope = ["read:org"] } : {}
      )
    } : {},
    local.auth_method == "google" ? {
      GoogleOAuthenticator = merge(
        {
          client_id          = local.auth.google.client_id
          oauth_callback_url = local.auth.google.oauth_callback_url
        },
        length(try(local.auth.google.hosted_domains, [])) > 0 ? { hosted_domain = local.auth.google.hosted_domains } : {}
      )
    } : {},
    local.auth_method == "oidc" ? {
      GenericOAuthenticator = {
        client_id          = local.auth.oidc.client_id
        oauth_callback_url = local.auth.oidc.oauth_callback_url
        authorize_url      = local.auth.oidc.authorize_url
        token_url          = local.auth.oidc.token_url
        userdata_url       = local.auth.oidc.userdata_url
        scope              = local.oidc_scopes
        username_claim     = local.oidc_username_claim
        login_service      = local.oidc_login_service
      }
    } : {}
  )

  native_config = merge(
    try(local.auth.native.open_signup, false) ? { open_signup = true } : {},
    try(local.auth.native.minimum_password_length, null) != null ? { minimum_password_length = local.auth.native.minimum_password_length } : {}
  )

  # The python snippets consuming the env vars (twin of the Pulumi
  # module's authBlocks).
  hub_extra_config = merge(
    local.auth_method == "shared_password" ? {
      plantonSharedPassword = "import os\nc.DummyAuthenticator.password = os.environ[\"${local.shared_password_env_var}\"]\n"
    } : {},
    local.auth_method == "github" ? {
      plantonOauthClientSecret = "import os\nc.GitHubOAuthenticator.client_secret = os.environ[\"${local.oauth_client_secret_env_var}\"]\n"
    } : {},
    local.auth_method == "google" ? {
      plantonOauthClientSecret = "import os\nc.GoogleOAuthenticator.client_secret = os.environ[\"${local.oauth_client_secret_env_var}\"]\n"
    } : {},
    local.auth_method == "oidc" ? {
      plantonOauthClientSecret = "import os\nc.GenericOAuthenticator.client_secret = os.environ[\"${local.oauth_client_secret_env_var}\"]\n"
    } : {}
  )

  hub_extra_env = merge(
    local.auth_method == "shared_password" ? {
      (local.shared_password_env_var) = {
        valueFrom = {
          secretKeyRef = {
            name = local.shared_password_secret_name
            key  = local.shared_password_secret_key
          }
        }
      }
    } : {},
    local.oauth_client_secret_ref != null ? {
      (local.oauth_client_secret_env_var) = {
        valueFrom = {
          secretKeyRef = {
            name = local.oauth_client_secret_ref.secret_name
            key  = local.oauth_client_secret_key
          }
        }
      }
    } : {}
  )

  # ---- resources helper renderings (the null-prune idiom) ----------------
  component_resources = {
    hub   = try(var.spec.hub.resources, null)
    proxy = try(var.spec.proxy.resources, null)
  }

  rendered_resources = {
    for name, r in local.component_resources : name => r == null ? null : {
      for rk, rv in {
        requests = try(r.requests, null) == null ? null : {
          for qk, qv in {
            cpu    = r.requests.cpu != "" ? r.requests.cpu : null
            memory = r.requests.memory != "" ? r.requests.memory : null
          } : qk => qv if qv != null
        }
        limits = try(r.limits, null) == null ? null : {
          for lk, lv in {
            cpu    = r.limits.cpu != "" ? r.limits.cpu : null
            memory = r.limits.memory != "" ? r.limits.memory : null
          } : lk => lv if lv != null
        }
      } : rk => rv if rv != null && rv != {}
    }
  }

  # ------------------------------ hub block ------------------------------
  hub_spec = try(var.spec.hub, null)

  hub_db_pvc_block = merge(
    try(local.hub_sqlite.storage_size, "") != "" && try(local.hub_sqlite.storage_size, "") != null ? { storage = local.hub_sqlite.storage_size } : {},
    try(local.hub_sqlite.storage_class, "") != "" ? { storageClassName = local.hub_sqlite.storage_class } : {}
  )

  hub_db_block = merge(
    { type = local.db_type },
    local.db_type == "sqlite-pvc" && length(local.hub_db_pvc_block) > 0 ? { pvc = local.hub_db_pvc_block } : {},
    local.db_type != "sqlite-pvc" ? { url = local.db_url } : {}
  )

  hub_resources = local.rendered_resources.hub

  hub_block = merge(
    { db = local.hub_db_block },
    local.hub_secret_enabled ? { existingSecret = local.hub_secret_name } : {},
    try(local.hub_spec.concurrent_spawn_limit, null) != null ? { concurrentSpawnLimit = local.hub_spec.concurrent_spawn_limit } : {},
    try(local.hub_spec.active_server_limit, null) != null ? { activeServerLimit = local.hub_spec.active_server_limit } : {},
    try(local.hub_spec.allow_named_servers, false) ? { allowNamedServers = true } : {},
    try(local.hub_spec.allow_named_servers, false) && try(local.hub_spec.named_server_limit_per_user, null) != null ? { namedServerLimitPerUser = local.hub_spec.named_server_limit_per_user } : {},
    try(local.hub_spec.shutdown_on_logout, false) ? { shutdownOnLogout = true } : {},
    local.hub_resources != null && local.hub_resources != {} ? { resources = local.hub_resources } : {},
    { config = local.hub_config },
    length(local.hub_extra_config) > 0 ? { extraConfig = local.hub_extra_config } : {},
    length(local.hub_extra_env) > 0 ? { extraEnv = local.hub_extra_env } : {},
    length(try(var.spec.scheduling.core_node_selector, {})) > 0 ? { nodeSelector = var.spec.scheduling.core_node_selector } : {},
    local.network_policy_disabled ? { networkPolicy = { enabled = false } } : {}
  )

  # ------------------------------ proxy ----------------------------------
  # DELIBERATE chart-default override: the chart ships
  # proxy.service.type LoadBalancer; this kind composes exposure from
  # first-class kinds, so the front door stays ClusterIP unless the
  # spec says otherwise.
  proxy_service_type = try(var.spec.proxy.service_type, "") != "" && try(var.spec.proxy.service_type, "") != null ? var.spec.proxy.service_type : "ClusterIP"
  proxy_resources    = local.rendered_resources.proxy

  proxy_chp_block = merge(
    local.proxy_resources != null && local.proxy_resources != {} ? { resources = local.proxy_resources } : {},
    length(try(var.spec.scheduling.core_node_selector, {})) > 0 ? { nodeSelector = var.spec.scheduling.core_node_selector } : {},
    local.network_policy_disabled ? { networkPolicy = { enabled = false } } : {}
  )

  proxy_block = merge(
    {
      service = merge(
        { type = local.proxy_service_type },
        length(try(var.spec.proxy.service_annotations, {})) > 0 ? { annotations = var.spec.proxy.service_annotations } : {}
      )
    },
    length(local.proxy_chp_block) > 0 ? { chp = local.proxy_chp_block } : {}
  )

  # ---------------------------- singleuser --------------------------------
  single_user = try(var.spec.single_user, null)

  single_user_storage = try(local.single_user.storage, null)
  # merge-of-conditionals: the three storage arms carry different object
  # attributes, so they merge against {} instead of chaining ternaries.
  single_user_storage_block = merge(
    try(local.single_user_storage.dynamic, null) != null ? { type = "dynamic" } : {},
    try(local.single_user_storage.dynamic.capacity, "") != "" && try(local.single_user_storage.dynamic.capacity, "") != null ? { capacity = local.single_user_storage.dynamic.capacity } : {},
    try(local.single_user_storage.dynamic.storage_class, "") != "" ? { dynamic = { storageClass = local.single_user_storage.dynamic.storage_class } } : {},
    # The static arm splits across two single-attribute ternaries — a
    # `{type=string, static=object}` branch against {} cannot unify
    # (mixed attribute types defeat HCL's map conversion).
    try(local.single_user_storage.static, null) != null ? { type = "static" } : {},
    try(local.single_user_storage.static, null) != null ? {
      static = {
        pvcName = local.single_user_storage.static.pvc_name
        subPath = try(local.single_user_storage.static.sub_path, "") != "" && try(local.single_user_storage.static.sub_path, "") != null ? local.single_user_storage.static.sub_path : "{username}"
      }
    } : {},
    try(local.single_user_storage.none, null) != null ? { type = "none" } : {}
  )

  # The chart takes CPU as a NUMBER (values schema) — tonumber() fails
  # loudly on a non-numeric string, the engine-level twin of the Pulumi
  # module's ParseFloat error.
  single_user_profiles = [
    for profile in try(local.single_user.profiles, []) : merge(
      { display_name = profile.display_name },
      try(profile.description, "") != "" ? { description = profile.description } : {},
      try(profile.default, false) ? { default = true } : {},
      length(merge(
        try(profile.image, null) != null ? { image = "${profile.image.repository}:${profile.image.tag}" } : {},
        try(profile.memory_guarantee, "") != "" ? { mem_guarantee = profile.memory_guarantee } : {},
        try(profile.memory_limit, "") != "" ? { mem_limit = profile.memory_limit } : {},
        try(profile.cpu_guarantee, "") != "" ? { cpu_guarantee = tonumber(profile.cpu_guarantee) } : {},
        try(profile.cpu_limit, "") != "" ? { cpu_limit = tonumber(profile.cpu_limit) } : {}
        )) > 0 ? {
        kubespawner_override = merge(
          try(profile.image, null) != null ? { image = "${profile.image.repository}:${profile.image.tag}" } : {},
          try(profile.memory_guarantee, "") != "" ? { mem_guarantee = profile.memory_guarantee } : {},
          try(profile.memory_limit, "") != "" ? { mem_limit = profile.memory_limit } : {},
          try(profile.cpu_guarantee, "") != "" ? { cpu_guarantee = tonumber(profile.cpu_guarantee) } : {},
          try(profile.cpu_limit, "") != "" ? { cpu_limit = tonumber(profile.cpu_limit) } : {}
        )
      } : {}
    )
  ]

  # A single merge covers both the absent and present single_user cases:
  # every part is try()-guarded, so an absent block contributes nothing.
  single_user_block = merge(
    try(local.single_user.image, null) != null ? {
      image = {
        name = local.single_user.image.repository
        tag  = local.single_user.image.tag
      }
    } : {},
    length(merge(
      try(local.single_user.memory_guarantee, "") != "" && try(local.single_user.memory_guarantee, "") != null ? { guarantee = local.single_user.memory_guarantee } : {},
      try(local.single_user.memory_limit, "") != "" ? { limit = local.single_user.memory_limit } : {}
      )) > 0 ? {
      memory = merge(
        try(local.single_user.memory_guarantee, "") != "" && try(local.single_user.memory_guarantee, "") != null ? { guarantee = local.single_user.memory_guarantee } : {},
        try(local.single_user.memory_limit, "") != "" ? { limit = local.single_user.memory_limit } : {}
      )
    } : {},
    length(merge(
      try(local.single_user.cpu_guarantee, "") != "" ? { guarantee = tonumber(local.single_user.cpu_guarantee) } : {},
      try(local.single_user.cpu_limit, "") != "" ? { limit = tonumber(local.single_user.cpu_limit) } : {}
      )) > 0 ? {
      cpu = merge(
        try(local.single_user.cpu_guarantee, "") != "" ? { guarantee = tonumber(local.single_user.cpu_guarantee) } : {},
        try(local.single_user.cpu_limit, "") != "" ? { limit = tonumber(local.single_user.cpu_limit) } : {}
      )
    } : {},
    length(local.single_user_storage_block) > 0 ? { storage = local.single_user_storage_block } : {},
    try(local.single_user.default_url, "") != "" ? { defaultUrl = local.single_user.default_url } : {},
    try(local.single_user.start_timeout_seconds, null) != null ? { startTimeout = local.single_user.start_timeout_seconds } : {},
    length(try(local.single_user.extra_env, {})) > 0 ? { extraEnv = local.single_user.extra_env } : {},
    length(local.single_user_profiles) > 0 ? { profileList = local.single_user_profiles } : {},
    length(try(var.spec.scheduling.user_node_selector, {})) > 0 ? { nodeSelector = var.spec.scheduling.user_node_selector } : {},
    local.network_policy_disabled ? { networkPolicy = { enabled = false } } : {}
  )

  # ---------------------------- scheduling --------------------------------
  scheduling = try(var.spec.scheduling, null)

  placeholder_replicas = try(local.scheduling.user_placeholder_replicas, null)
  scheduling_block = merge(
    try(local.scheduling.user_scheduler_enabled, null) != null ? {
      userScheduler = { enabled = local.scheduling.user_scheduler_enabled }
    } : {},
    local.placeholder_replicas != null ? {
      userPlaceholder = {
        enabled  = local.placeholder_replicas > 0
        replicas = local.placeholder_replicas
      }
    } : {},
    # Placeholders only pre-warm capacity when real users can EVICT
    # them — the pod-priority machinery (chart default off).
    local.placeholder_replicas != null && try(local.placeholder_replicas, 0) > 0 ? { podPriority = { enabled = true } } : {}
  )

  # ------------------------------ culling ---------------------------------
  culling = try(var.spec.culling, null)
  # try()-guarded parts make the outer null check unnecessary (and a
  # `null ? {} : merge(...)` ternary would trip HCL type unification).
  cull_block = merge(
    try(local.culling.enabled, null) != null ? { enabled = local.culling.enabled } : {},
    try(local.culling.timeout_seconds, null) != null ? { timeout = local.culling.timeout_seconds } : {},
    try(local.culling.every_seconds, null) != null ? { every = local.culling.every_seconds } : {},
    try(local.culling.max_age_seconds, null) != null ? { maxAge = local.culling.max_age_seconds } : {},
    try(local.culling.cull_users, false) ? { users = true } : {}
  )

  # ----------------------------- pre-puller -------------------------------
  pre_puller = try(var.spec.pre_puller, null)
  pre_puller_block = merge(
    try(local.pre_puller.hook_enabled, null) != null ? { hook = { enabled = local.pre_puller.hook_enabled } } : {},
    try(local.pre_puller.continuous_enabled, null) != null ? { continuous = { enabled = local.pre_puller.continuous_enabled } } : {}
  )

  # One spec toggle drives the chart's three per-component NetworkPolicy
  # switches identically (they default true; only an explicit false
  # needs rendering).
  network_policy_disabled = try(var.spec.network_policy_enabled, null) != null ? !var.spec.network_policy_enabled : false

  # ------------------------- assembled values -----------------------------
  helm_values = merge(
    {
      hub   = local.hub_block
      proxy = local.proxy_block
    },
    length(local.single_user_block) > 0 ? { singleuser = local.single_user_block } : {},
    length(local.scheduling_block) > 0 ? { scheduling = local.scheduling_block } : {},
    length(local.cull_block) > 0 ? { cull = local.cull_block } : {},
    length(local.pre_puller_block) > 0 ? { prePuller = local.pre_puller_block } : {}
  )

  # ------------------------------- outputs --------------------------------
  # CHART-FIXED resource names: at fullnameOverride "" every resource
  # renders a bare name — a per-NAMESPACE singleton.
  proxy_public_service_name = "proxy-public"
  hub_service_name          = "hub"
  proxy_public_endpoint     = "http://${local.proxy_public_service_name}.${local.namespace}.svc.cluster.local:80"
  port_forward_command      = "kubectl port-forward svc/${local.proxy_public_service_name} -n ${local.namespace} 8080:80"

  shared_password_output_name = local.auth_method == "shared_password" ? local.shared_password_secret_name : ""
  shared_password_output_key  = local.auth_method == "shared_password" ? local.shared_password_secret_key : ""
}

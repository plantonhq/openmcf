# Computed values for the KubernetesSecretStore module.
#
# The CR spec rendered here is the Terraform twin of the shared Go builder
# (pkg/iac/pulumi/pulumimodule/provider/kubernetes/externalsecretsstore)
# that both the KubernetesClusterSecretStore and KubernetesSecretStore
# Pulumi modules use — and the KubernetesClusterSecretStore Terraform
# module carries the same rendering with a cluster-scoped kind. Keep all of
# them in lockstep: same CRD field names, same credential Secret name, same
# secret data keys.
#
# CREDENTIAL MODEL: wherever the CRD expects a secretRef, the spec carries
# the credential VALUE (sensitive). local.credential_secrets collects the
# one deterministic "<resource-name>-credentials" Secret to materialize;
# the CR's secretRefs point at it. A namespaced store's secret references
# resolve in its OWN namespace by default, so refs carry no explicit
# namespace (the cluster twin stamps one on every ref).
#
# Optional nested blocks are read with try(): HCL's && does NOT
# short-circuit, so chained null checks still dereference the null.

locals {
  # SecretStore metadata.name — the name ExternalSecrets in the same
  # namespace reference (kind SecretStore, the upstream default).
  store_name = var.metadata.name

  # Namespace the SecretStore and its credential Secrets live in (resolved
  # literal from the spec's value-or-ref).
  namespace = var.spec.namespace

  labels = merge(concat(
    [{
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesSecretStore"
    }],
    (var.metadata.id != null && var.metadata.id != "") ? [{ "planton.ai/resource-id" = var.metadata.id }] : [],
    (var.metadata.org != null && var.metadata.org != "") ? [{ "planton.ai/organization" = var.metadata.org }] : [],
    (var.metadata.env != null && var.metadata.env != "") ? [{ "planton.ai/environment" = var.metadata.env }] : []
  )...)

  aws   = try(var.spec.config.aws, null)
  gcp   = try(var.spec.config.gcp_secret_manager, null)
  azure = try(var.spec.config.azure_key_vault, null)
  vault = try(var.spec.config.vault, null)
  k8s   = try(var.spec.config.kubernetes, null)
  fake  = try(var.spec.config.fake, null)

  # One deterministic credential Secret per store — the exact twin of the
  # Go builder's "<resource-name>-credentials" derivation.
  credential_secret_name = "${local.store_name}-credentials"

  # Twin of the Go builder's secretRef closure: every credential reference
  # addresses the one credential Secret by fixed key. Namespaced stores
  # omit the namespace — refs default to the store's own.
  secret_ref = {
    for key in [
      "access-key-id",
      "secret-access-key",
      "credentials.json",
      "client-id",
      "client-secret",
      "vault-token",
      "vault-approle-secret-id",
      "bearer-token",
    ] : key => { name = local.credential_secret_name, key = key }
  }

  # ---- backend provider blocks ------------------------------------------
  # One rendered block per backend arm, null when the arm is not selected
  # (the proto oneof guarantees exactly one is). Null-prune idiom
  # throughout this file: one object literal whose conditional entries are
  # null when absent, filtered by the for-expression — heterogeneous
  # conditional merges (`cond ? {...} : {}` / concat-spread) fail HCL type
  # unification when sibling entries infer as different object types.

  provider_aws = local.aws == null ? null : {
    for k, v in {
      # service renders unconditionally, matching the Go builder — the
      # platform applies the "SecretsManager" default before input reaches
      # the module.
      service = try(local.aws.service, null) != null ? local.aws.service : ""
      region  = local.aws.region
      role    = try(local.aws.role, "") != "" ? local.aws.role : null
      prefix  = try(local.aws.prefix, "") != "" ? local.aws.prefix : null
      # Static keys and the ServiceAccount reference are mutually
      # exclusive by proto validation; no auth block at all = the
      # operator's ambient identity (its controller ServiceAccount / node
      # role) — upstream's documented fallback.
      auth = (try(local.aws.access_key_id, "") != "" || try(local.aws.service_account_name, "") != "") ? {
        for ak, av in {
          secretRef = try(local.aws.access_key_id, "") != "" ? {
            accessKeyIDSecretRef     = local.secret_ref["access-key-id"]
            secretAccessKeySecretRef = local.secret_ref["secret-access-key"]
          } : null
          jwt = (try(local.aws.access_key_id, "") == "" && try(local.aws.service_account_name, "") != "") ? {
            # On a SecretStore the serviceAccountRef namespace defaults to
            # the store's own namespace; an explicit override is unusual,
            # so it is carried only when the spec sets one.
            serviceAccountRef = {
              for rk, rv in {
                name      = local.aws.service_account_name
                namespace = try(local.aws.service_account_namespace, "") != "" ? local.aws.service_account_namespace : null
              } : rk => rv if rv != null
            }
          } : null
        } : ak => av if av != null
      } : null
    } : k => v if v != null
  }

  provider_gcp = local.gcp == null ? null : {
    for k, v in {
      projectID = local.gcp.project_id
      location  = try(local.gcp.location, "") != "" ? local.gcp.location : null
      auth = (try(local.gcp.service_account_key_json, "") != "" || try(local.gcp.service_account_name, "") != "") ? {
        for ak, av in {
          secretRef = try(local.gcp.service_account_key_json, "") != "" ? {
            secretAccessKeySecretRef = local.secret_ref["credentials.json"]
          } : null
          workloadIdentity = (try(local.gcp.service_account_key_json, "") == "" && try(local.gcp.service_account_name, "") != "") ? {
            serviceAccountRef = {
              for rk, rv in {
                name      = local.gcp.service_account_name
                namespace = try(local.gcp.service_account_namespace, "") != "" ? local.gcp.service_account_namespace : null
              } : rk => rv if rv != null
            }
          } : null
        } : ak => av if av != null
      } : null
    } : k => v if v != null
  }

  provider_azure = local.azure == null ? null : {
    for k, v in {
      vaultUrl   = local.azure.vault_url
      authType   = try(local.azure.auth_type, null)
      tenantId   = try(local.azure.tenant_id, "") != "" ? local.azure.tenant_id : null
      identityId = try(local.azure.identity_id, "") != "" ? local.azure.identity_id : null
      serviceAccountRef = try(local.azure.service_account_name, "") != "" ? {
        for rk, rv in {
          name      = local.azure.service_account_name
          namespace = try(local.azure.service_account_namespace, "") != "" ? local.azure.service_account_namespace : null
        } : rk => rv if rv != null
      } : null
      authSecretRef = try(local.azure.client_id, "") != "" ? {
        clientId     = local.secret_ref["client-id"]
        clientSecret = local.secret_ref["client-secret"]
      } : null
    } : k => v if v != null
  }

  provider_vault = local.vault == null ? null : {
    for k, v in {
      server    = local.vault.server
      path      = try(local.vault.path, "") != "" ? local.vault.path : null
      version   = try(local.vault.version, null)
      namespace = try(local.vault.namespace, "") != "" ? local.vault.namespace : null
      # The CRD carries caBundle as base64 []byte JSON; both engines pass
      # the SAME base64-encoded PEM string through unchanged.
      caBundle = try(local.vault.ca_bundle, "") != "" ? local.vault.ca_bundle : null
      # Exactly one auth arm per proto validation (the Go builder errors
      # when none is set; here an unset arm simply renders no entry).
      auth = {
        for ak, av in {
          tokenSecretRef = try(local.vault.token, null) != null ? local.secret_ref["vault-token"] : null
          appRole = try(local.vault.app_role, null) == null ? null : {
            for pk, pv in {
              roleId    = local.vault.app_role.role_id
              secretRef = local.secret_ref["vault-approle-secret-id"]
              path      = try(local.vault.app_role.path, "") != "" ? local.vault.app_role.path : null
            } : pk => pv if pv != null
          }
          kubernetes = try(local.vault.kubernetes, null) == null ? null : {
            for kk, kv in {
              role      = local.vault.kubernetes.role
              mountPath = try(local.vault.kubernetes.mount_path, "") != "" ? local.vault.kubernetes.mount_path : null
              # No per-backend namespace exists here and a namespaced
              # store's refs default to its own namespace — name only.
              serviceAccountRef = try(local.vault.kubernetes.service_account_name, "") != "" ? {
                name = local.vault.kubernetes.service_account_name
              } : null
            } : kk => kv if kv != null
          }
        } : ak => av if av != null
      }
    } : k => v if v != null
  }

  provider_k8s = local.k8s == null ? null : {
    for k, v in {
      server = (try(local.k8s.server_url, "") != "" || try(local.k8s.ca_bundle, "") != "") ? {
        for sk, sv in {
          url      = try(local.k8s.server_url, "") != "" ? local.k8s.server_url : null
          caBundle = try(local.k8s.ca_bundle, "") != "" ? local.k8s.ca_bundle : null
        } : sk => sv if sv != null
      } : null
      remoteNamespace = try(local.k8s.remote_namespace, "") != "" ? local.k8s.remote_namespace : null
      # Bearer token and the ServiceAccount reference are mutually
      # exclusive by proto validation.
      auth = (try(local.k8s.token, "") != "" || try(local.k8s.service_account_name, "") != "") ? {
        for ak, av in {
          token = try(local.k8s.token, "") != "" ? {
            bearerToken = local.secret_ref["bearer-token"]
          } : null
          serviceAccount = (try(local.k8s.token, "") == "" && try(local.k8s.service_account_name, "") != "") ? {
            name = local.k8s.service_account_name
          } : null
        } : ak => av if av != null
      } : null
    } : k => v if v != null
  }

  provider_fake = local.fake == null ? null : {
    data = [
      for entry in local.fake.data : {
        for ek, ev in {
          key     = entry.key
          value   = entry.value
          version = try(entry.version, "") != "" ? entry.version : null
        } : ek => ev if ev != null
      }
    ]
  }

  store_provider = {
    for k, v in {
      aws        = local.provider_aws
      gcpsm      = local.provider_gcp
      azurekv    = local.provider_azure
      vault      = local.provider_vault
      kubernetes = local.provider_k8s
      fake       = local.provider_fake
    } : k => v if v != null
  }

  # ---- credential Secret to materialize ----------------------------------
  # stringData for the one "<resource-name>-credentials" Secret. Keys are
  # the exact twins of the Go builder's credentialData assignments; empty
  # map = the selected backend declares no static credentials, and no
  # Secret is created.
  credential_data = {
    for k, v in {
      "access-key-id"           = try(local.aws.access_key_id, "") != "" ? local.aws.access_key_id : null
      "secret-access-key"       = try(local.aws.access_key_id, "") != "" ? local.aws.secret_access_key : null
      "credentials.json"        = try(local.gcp.service_account_key_json, "") != "" ? local.gcp.service_account_key_json : null
      "client-id"               = try(local.azure.client_id, "") != "" ? local.azure.client_id : null
      "client-secret"           = try(local.azure.client_id, "") != "" ? local.azure.client_secret : null
      "vault-token"             = try(local.vault.token, null) != null ? local.vault.token.token : null
      "vault-approle-secret-id" = try(local.vault.app_role, null) != null ? local.vault.app_role.secret_id : null
      "bearer-token"            = try(local.k8s.token, "") != "" ? local.k8s.token : null
    } : k => v if v != null
  }

  credential_secrets = {
    for name, data in { (local.credential_secret_name) = local.credential_data } :
    name => data if length(data) > 0
  }

  # ---- the CR spec --------------------------------------------------------
  store_spec = {
    for k, v in {
      provider        = local.store_provider
      controller      = try(var.spec.config.controller_class, "") != "" ? var.spec.config.controller_class : null
      refreshInterval = try(var.spec.config.refresh_interval, "") != "" ? var.spec.config.refresh_interval : null
      retrySettings = try(var.spec.config.retry, null) == null ? null : (
        (try(var.spec.config.retry.max_retries, null) != null || try(var.spec.config.retry.retry_interval, "") != "") ? {
          for rk, rv in {
            # maxRetries stays a number all the way into yamlencode — the
            # null-prune object literal preserves attribute types.
            maxRetries    = try(var.spec.config.retry.max_retries, null)
            retryInterval = try(var.spec.config.retry.retry_interval, "") != "" ? var.spec.config.retry.retry_interval : null
          } : rk => rv if rv != null
        } : null
      )
    } : k => v if v != null
  }
}

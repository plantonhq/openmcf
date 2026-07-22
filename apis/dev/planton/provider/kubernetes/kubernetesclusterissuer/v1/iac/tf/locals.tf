# Computed values for the KubernetesClusterIssuer module.
#
# The CR spec rendered here is the Terraform twin of the shared Go builder
# (pkg/iac/pulumi/pulumimodule/provider/kubernetes/certmanagerissuer) that
# both the KubernetesClusterIssuer and KubernetesIssuer Pulumi modules use —
# and the KubernetesIssuer Terraform module carries the same rendering with
# a namespace-scoped kind. Keep all of them in lockstep: same CRD field
# names, same credential Secret names, same secret data keys.
#
# CREDENTIAL MODEL: wherever the CRD expects a secretRef, the spec carries
# the credential VALUE (sensitive). local.credential_secrets collects the
# Kubernetes Secrets to materialize (deterministic names derived from the
# resource name); the CR's secretRefs point at them.
#
# Optional nested blocks are read with try(): HCL's && does NOT
# short-circuit, so chained null checks still dereference the null.

locals {
  cluster_issuer_name = var.metadata.name

  # Namespace credential Secrets are created in: cert-manager's
  # cluster-resource namespace — the only namespace cert-manager reads
  # Secrets from for cluster-scoped resources.
  secrets_namespace = var.spec.cert_manager_namespace

  labels = merge(concat(
    [{
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesClusterIssuer"
    }],
    (var.metadata.id != null && var.metadata.id != "") ? [{ "planton.ai/resource-id" = var.metadata.id }] : [],
    (var.metadata.org != null && var.metadata.org != "") ? [{ "planton.ai/organization" = var.metadata.org }] : [],
    (var.metadata.env != null && var.metadata.env != "") ? [{ "planton.ai/environment" = var.metadata.env }] : []
  )...)

  acme  = try(var.spec.config.acme, null)
  ca    = try(var.spec.config.ca, null)
  self  = try(var.spec.config.self_signed, null)
  vault = try(var.spec.config.vault, null)

  # ---- ACME ------------------------------------------------------------
  acme_account_key_secret_name = local.acme != null ? "${local.cluster_issuer_name}-acme-account-key" : ""

  acme_eab_secret_name = local.acme != null && try(local.acme.external_account_binding, null) != null ? "${local.cluster_issuer_name}-acme-eab" : ""

  # One rendered solver per spec solver, index-stable so credential Secret
  # names stay deterministic. Null-prune idiom throughout this file: one
  # object literal whose conditional entries are null when absent, filtered
  # by the for-expression — heterogeneous conditional merges (`cond ? {...}
  # : {}` / concat-spread) fail HCL type unification when sibling entries
  # infer as different object types.
  acme_solvers = [
    for i, s in try(local.acme.solvers, []) : {
      for k, v in {
        selector = try(s.selector, null) == null ? null : (
          length(try(s.selector.dns_zones, [])) > 0 || length(try(s.selector.dns_names, [])) > 0 || length(try(s.selector.match_labels, {})) > 0
          ) ? {
          for sk, sv in {
            dnsZones    = length(try(s.selector.dns_zones, [])) > 0 ? s.selector.dns_zones : null
            dnsNames    = length(try(s.selector.dns_names, [])) > 0 ? s.selector.dns_names : null
            matchLabels = length(try(s.selector.match_labels, {})) > 0 ? s.selector.match_labels : null
          } : sk => sv if sv != null
        } : null
        http01 = try(s.http01, null) != null ? local.http01_rendered[i] : null
        dns01  = try(s.dns01, null) != null ? local.dns01_pruned[i] : null
      } : k => v if v != null
    }
  ]

  http01_rendered = {
    for i, s in try(local.acme.solvers, []) : i => {
      for k, v in {
        ingress = try(s.http01.ingress, null) == null ? null : {
          for ik, iv in {
            ingressClassName = try(s.http01.ingress.ingress_class_name, "") != "" ? s.http01.ingress.ingress_class_name : null
            name             = try(s.http01.ingress.name, "") != "" ? s.http01.ingress.name : null
            serviceType      = try(s.http01.ingress.service_type, null)
          } : ik => iv if iv != null
        }
        gatewayHTTPRoute = try(s.http01.gateway_http_route, null) == null ? null : {
          for gk, gv in {
            parentRefs = [for ref in try(s.http01.gateway_http_route.parent_refs, []) : {
              for rk, rv in {
                name        = ref.name
                namespace   = try(ref.namespace, "") != "" ? ref.namespace : null
                sectionName = try(ref.section_name, "") != "" ? ref.section_name : null
              } : rk => rv if rv != null
            }]
            labels      = length(try(s.http01.gateway_http_route.labels, {})) > 0 ? s.http01.gateway_http_route.labels : null
            serviceType = try(s.http01.gateway_http_route.service_type, null)
          } : gk => gv if gv != null
        }
      } : k => v if v != null
    } if try(s.http01, null) != null
  }

  # DNS-01 provider rendering. CRD JSON quirks (verified against the pinned
  # CRD schema): "cloudDNS", "azureDNS", "acmeDNS" capitalization,
  # "clientSecretSecretRef", Route53 nested auth.kubernetes.serviceAccountRef.
  #
  # Built with the null-prune idiom — one object literal whose conditional
  # entries are null when absent, filtered by the for-expression — because
  # concat-spread unification panics when sibling pieces infer as maps of
  # different element types (a heterogeneity this shape guarantees).
  dns01_pruned = {
    for i, s in try(local.acme.solvers, []) : i => {
      for k, v in {
        cnameStrategy = try(s.dns01.cname_strategy, null) != null ? (s.dns01.cname_strategy == "follow" ? "Follow" : "None") : null

        cloudflare = try(s.dns01.cloudflare, null) == null ? null : {
          for ck, cv in {
            email             = try(s.dns01.cloudflare.api_key.email, null)
            apiTokenSecretRef = try(s.dns01.cloudflare.api_token, null) != null ? { name = "${local.cluster_issuer_name}-solver${i}-cloudflare", key = "api-token" } : null
            apiKeySecretRef   = try(s.dns01.cloudflare.api_key, null) != null ? { name = "${local.cluster_issuer_name}-solver${i}-cloudflare", key = "api-key" } : null
          } : ck => cv if cv != null
        }

        route53 = try(s.dns01.route53, null) == null ? null : {
          for rk, rv in {
            region                   = s.dns01.route53.region
            hostedZoneID             = try(s.dns01.route53.hosted_zone_id, "") != "" ? s.dns01.route53.hosted_zone_id : null
            role                     = try(s.dns01.route53.assume_role_arn, "") != "" ? s.dns01.route53.assume_role_arn : null
            accessKeyID              = try(s.dns01.route53.static_credentials, null) != null ? s.dns01.route53.static_credentials.access_key_id : null
            secretAccessKeySecretRef = try(s.dns01.route53.static_credentials, null) != null ? { name = "${local.cluster_issuer_name}-solver${i}-route53", key = "secret-access-key" } : null
            auth = try(s.dns01.route53.service_account, null) == null ? null : {
              kubernetes = {
                serviceAccountRef = {
                  for sk, sv in {
                    name      = s.dns01.route53.service_account.service_account_name
                    audiences = length(try(s.dns01.route53.service_account.audiences, [])) > 0 ? s.dns01.route53.service_account.audiences : null
                  } : sk => sv if sv != null
                }
              }
            }
          } : rk => rv if rv != null
        }

        azureDNS = try(s.dns01.azure_dns, null) == null ? null : {
          for ak, av in {
            subscriptionID        = s.dns01.azure_dns.subscription_id
            resourceGroupName     = s.dns01.azure_dns.resource_group_name
            hostedZoneName        = try(s.dns01.azure_dns.hosted_zone_name, "") != "" ? s.dns01.azure_dns.hosted_zone_name : null
            zoneType              = try(s.dns01.azure_dns.zone_type, null) != null ? (s.dns01.azure_dns.zone_type == "private" ? "AzurePrivateZone" : "AzurePublicZone") : null
            environment           = try(s.dns01.azure_dns.environment, null)
            clientID              = try(s.dns01.azure_dns.client_secret, "") != "" ? s.dns01.azure_dns.client_id : null
            tenantID              = try(s.dns01.azure_dns.client_secret, "") != "" ? s.dns01.azure_dns.tenant_id : null
            clientSecretSecretRef = try(s.dns01.azure_dns.client_secret, "") != "" ? { name = "${local.cluster_issuer_name}-solver${i}-azure-dns", key = "client-secret" } : null
            managedIdentity = try(s.dns01.azure_dns.managed_identity, null) == null ? null : {
              for mk, mv in {
                clientID   = try(s.dns01.azure_dns.managed_identity.client_id, "") != "" ? s.dns01.azure_dns.managed_identity.client_id : null
                resourceID = try(s.dns01.azure_dns.managed_identity.resource_id, "") != "" ? s.dns01.azure_dns.managed_identity.resource_id : null
              } : mk => mv if mv != null
            }
          } : ak => av if av != null
        }

        cloudDNS = try(s.dns01.gcp_cloud_dns, null) == null ? null : {
          for gk, gv in {
            project                 = s.dns01.gcp_cloud_dns.project_id
            hostedZoneName          = try(s.dns01.gcp_cloud_dns.hosted_zone_name, "") != "" ? s.dns01.gcp_cloud_dns.hosted_zone_name : null
            serviceAccountSecretRef = try(s.dns01.gcp_cloud_dns.service_account_key_json, "") != "" ? { name = "${local.cluster_issuer_name}-solver${i}-clouddns", key = "key.json" } : null
          } : gk => gv if gv != null
        }

        digitalocean = try(s.dns01.digitalocean, null) == null ? null : {
          tokenSecretRef = { name = "${local.cluster_issuer_name}-solver${i}-digitalocean", key = "access-token" }
        }

        rfc2136 = try(s.dns01.rfc2136, null) == null ? null : {
          for fk, fv in {
            nameserver          = s.dns01.rfc2136.nameserver
            tsigKeyName         = try(s.dns01.rfc2136.tsig_key_name, "") != "" ? s.dns01.rfc2136.tsig_key_name : null
            tsigSecretSecretRef = try(s.dns01.rfc2136.tsig_key_name, "") != "" ? { name = "${local.cluster_issuer_name}-solver${i}-rfc2136", key = "tsig-secret" } : null
            tsigAlgorithm       = try(s.dns01.rfc2136.tsig_algorithm, "") != "" ? s.dns01.rfc2136.tsig_algorithm : null
          } : fk => fv if fv != null
        }

        acmeDNS = try(s.dns01.acme_dns, null) == null ? null : {
          host             = s.dns01.acme_dns.host
          accountSecretRef = { name = "${local.cluster_issuer_name}-solver${i}-acme-dns", key = "acmedns.json" }
        }

        akamai = try(s.dns01.akamai, null) == null ? null : {
          serviceConsumerDomain = s.dns01.akamai.service_consumer_domain
          clientTokenSecretRef  = { name = "${local.cluster_issuer_name}-solver${i}-akamai", key = "client-token" }
          clientSecretSecretRef = { name = "${local.cluster_issuer_name}-solver${i}-akamai", key = "client-secret" }
          accessTokenSecretRef  = { name = "${local.cluster_issuer_name}-solver${i}-akamai", key = "access-token" }
        }

        webhook = try(s.dns01.webhook, null) == null ? null : {
          for wk, wv in {
            groupName  = s.dns01.webhook.group_name
            solverName = s.dns01.webhook.solver_name
            config     = try(s.dns01.webhook.config_yaml, "") != "" ? yamldecode(s.dns01.webhook.config_yaml) : null
          } : wk => wv if wv != null
        }
      } : k => v if v != null
    } if try(s.dns01, null) != null
  }

  # ---- credential Secrets to materialize --------------------------------
  # Map: secret name -> stringData. Names and keys are the exact twins of
  # the shared Go builder's CredentialSecret derivations.
  solver_credential_secrets = merge(concat([{}], [
    for i, s in try(local.acme.solvers, []) : merge(
      try(s.dns01.cloudflare.api_token, null) != null ? {
        "${local.cluster_issuer_name}-solver${i}-cloudflare" = { "api-token" = s.dns01.cloudflare.api_token.token }
      } : {},
      try(s.dns01.cloudflare.api_key, null) != null ? {
        "${local.cluster_issuer_name}-solver${i}-cloudflare" = { "api-key" = s.dns01.cloudflare.api_key.key }
      } : {},
      try(s.dns01.route53.static_credentials, null) != null ? {
        "${local.cluster_issuer_name}-solver${i}-route53" = { "secret-access-key" = s.dns01.route53.static_credentials.secret_access_key }
      } : {},
      try(s.dns01.azure_dns.client_secret, "") != "" ? {
        "${local.cluster_issuer_name}-solver${i}-azure-dns" = { "client-secret" = s.dns01.azure_dns.client_secret }
      } : {},
      try(s.dns01.gcp_cloud_dns.service_account_key_json, "") != "" ? {
        "${local.cluster_issuer_name}-solver${i}-clouddns" = { "key.json" = s.dns01.gcp_cloud_dns.service_account_key_json }
      } : {},
      try(s.dns01.digitalocean, null) != null ? {
        "${local.cluster_issuer_name}-solver${i}-digitalocean" = { "access-token" = s.dns01.digitalocean.token }
      } : {},
      try(s.dns01.rfc2136.tsig_key_name, "") != "" ? {
        "${local.cluster_issuer_name}-solver${i}-rfc2136" = { "tsig-secret" = s.dns01.rfc2136.tsig_secret }
      } : {},
      try(s.dns01.acme_dns, null) != null ? {
        "${local.cluster_issuer_name}-solver${i}-acme-dns" = { "acmedns.json" = s.dns01.acme_dns.account_json }
      } : {},
      try(s.dns01.akamai, null) != null ? {
        "${local.cluster_issuer_name}-solver${i}-akamai" = {
          "client-token"  = s.dns01.akamai.client_token
          "client-secret" = s.dns01.akamai.client_secret
          "access-token"  = s.dns01.akamai.access_token
        }
      } : {}
    )
  ])...)

  credential_secrets = merge(concat(
    [local.solver_credential_secrets],
    (local.acme_eab_secret_name != "") ? [{
      (local.acme_eab_secret_name) = { "key" = local.acme.external_account_binding.hmac_key }
    }] : [],
    (local.vault != null && try(local.vault.token_auth, null) != null) ? [{
      "${local.cluster_issuer_name}-vault-token" = { "token" = local.vault.token_auth.token }
    }] : [],
    (local.vault != null && try(local.vault.app_role_auth, null) != null) ? [{
      "${local.cluster_issuer_name}-vault-approle" = { "secret-id" = local.vault.app_role_auth.secret_id }
    }] : []
  )...)

  # ---- the CR spec -------------------------------------------------------
  issuer_spec = {
    for k, v in {
      acme = local.acme == null ? null : {
        for ak, av in {
          email               = local.acme.email
          server              = local.acme.server
          privateKeySecretRef = { name = local.acme_account_key_secret_name }
          solvers             = local.acme_solvers
          profile             = try(local.acme.profile, "") != "" ? local.acme.profile : null
          preferredChain      = try(local.acme.preferred_chain, "") != "" ? local.acme.preferred_chain : null
          caBundle            = try(local.acme.ca_bundle, "") != "" ? local.acme.ca_bundle : null
          # skipTLSVerify and the two behavior flags are emitted only when
          # true — false is the API default, and omitting it keeps the
          # applied object byte-identical with the Pulumi engine's.
          skipTLSVerify               = try(local.acme.skip_tls_verify, false) ? true : null
          disableAccountKeyGeneration = try(local.acme.disable_account_key_generation, false) ? true : null
          enableDurationFeature       = try(local.acme.enable_duration_feature, false) ? true : null
          externalAccountBinding = local.acme_eab_secret_name == "" ? null : {
            keyID        = local.acme.external_account_binding.key_id
            keySecretRef = { name = local.acme_eab_secret_name, key = "key" }
          }
        } : ak => av if av != null
      }

      ca = local.ca == null ? null : {
        for ck, cv in {
          secretName             = local.ca.ca_secret_name
          crlDistributionPoints  = length(try(local.ca.crl_distribution_points, [])) > 0 ? local.ca.crl_distribution_points : null
          ocspServers            = length(try(local.ca.ocsp_servers, [])) > 0 ? local.ca.ocsp_servers : null
          issuingCertificateURLs = length(try(local.ca.issuing_certificate_urls, [])) > 0 ? local.ca.issuing_certificate_urls : null
        } : ck => cv if cv != null
      }

      selfSigned = local.self == null ? null : {
        for sk, sv in {
          crlDistributionPoints = length(try(local.self.crl_distribution_points, [])) > 0 ? local.self.crl_distribution_points : null
        } : sk => sv if sv != null
      }

      vault = local.vault == null ? null : {
        for vk, vv in {
          server     = local.vault.server
          path       = local.vault.path
          namespace  = try(local.vault.vault_namespace, "") != "" ? local.vault.vault_namespace : null
          caBundle   = try(local.vault.ca_bundle, "") != "" ? local.vault.ca_bundle : null
          serverName = try(local.vault.server_name, "") != "" ? local.vault.server_name : null
          auth = {
            for auth_k, auth_v in {
              tokenSecretRef = try(local.vault.token_auth, null) == null ? null : {
                name = "${local.cluster_issuer_name}-vault-token", key = "token"
              }
              appRole = try(local.vault.app_role_auth, null) == null ? null : {
                path      = local.vault.app_role_auth.path
                roleId    = local.vault.app_role_auth.role_id
                secretRef = { name = "${local.cluster_issuer_name}-vault-approle", key = "secret-id" }
              }
              kubernetes = try(local.vault.kubernetes_auth, null) == null ? null : {
                role      = local.vault.kubernetes_auth.role
                mountPath = "/v1/auth/${coalesce(try(local.vault.kubernetes_auth.mount_path, null), "kubernetes")}"
                serviceAccountRef = {
                  for rk, rv in {
                    name      = local.vault.kubernetes_auth.service_account_name
                    audiences = length(try(local.vault.kubernetes_auth.audiences, [])) > 0 ? local.vault.kubernetes_auth.audiences : null
                  } : rk => rv if rv != null
                }
              }
            } : auth_k => auth_v if auth_v != null
          }
        } : vk => vv if vv != null
      }
    } : k => v if v != null
  }
}

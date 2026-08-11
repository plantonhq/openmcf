variable "metadata" {
  description = "Metadata for the resource, including name and labels"
  type = object({
    name    = string,
    id      = optional(string),
    org     = optional(string),
    env     = optional(string),
    labels  = optional(map(string)),
    tags    = optional(list(string)),
    version = optional(object({ id = string, message = string }))
  })
}

variable "spec" {
  description = "Specification for the GCP Workload Identity Pool Provider"
  type = object({
    # The pool this provider belongs to — the bare pool ID. The CLI's tfvars
    # converter resolves StringValueOrRef fields to their literal string
    # before the module runs, so this arrives as a plain string.
    # Immutable in GCP (ForceNew).
    workload_identity_pool_id = string

    # The provider ID — the final component of the provider's resource name.
    # 4-32 chars of lowercase letters, digits, hyphens; the "gcp-" prefix is
    # reserved by Google. Immutable in GCP (ForceNew).
    workload_identity_pool_provider_id = string

    # The GCP project that owns the pool. Arrives as a plain string (resolved
    # reference). If empty, the provider's default project is used (locals.tf).
    project_id = optional(string, "")

    # Human-readable name shown in the GCP console (max 32 chars). Mutable.
    display_name = optional(string)

    # Which issuer this provider trusts and any operational notes
    # (max 256 chars). Mutable.
    description = optional(string)

    # Kill switch: a disabled provider rejects new token exchanges
    # (already-issued credentials remain valid until expiry). Mutable.
    disabled = optional(bool, false)

    # Claim -> Google-attribute CEL mappings ("google.subject" is required for
    # OIDC; AWS/SAML/X.509 have server-side defaults when omitted).
    attribute_mapping = optional(map(string), {})

    # CEL expression gating which otherwise valid credentials are accepted
    # (max 4096 chars). For multi-tenant issuers like GitHub Actions, always
    # set one.
    attribute_condition = optional(string)

    # The issuer — exactly one of aws / oidc / saml / x509.
    aws = optional(object({
      # The 12-digit AWS account ID whose workloads may federate.
      account_id = string
    }))

    oidc = optional(object({
      # The OIDC issuer URL; must match the `iss` claim of incoming tokens.
      issuer_uri = string
      # Acceptable `aud` values (max 10). Empty means the audience must equal
      # the provider's full canonical resource name — the safest default.
      allowed_audiences = optional(list(string), [])
      # JWKS JSON for issuers unreachable via .well-known discovery
      # (public keys, not a secret). Empty means use discovery.
      jwks_json = optional(string)
    }))

    saml = optional(object({
      # The IdP's SAML configuration metadata XML document.
      idp_metadata_xml = string
    }))

    x509 = optional(object({
      trust_store = object({
        # Incoming end-entity certificates must chain to one of these.
        trust_anchors = list(object({
          pem_certificate = string
        }))
        # Intermediates available for building the chain to a trust anchor.
        intermediate_cas = optional(list(object({
          pem_certificate = string
        })), [])
      })
    }))

    # DELETE (default) / PREVENT / ABANDON; empty falls through to the
    # provider default (DELETE).
    deletion_policy = optional(string, "")
  })

  validation {
    condition     = can(regex("^[a-z0-9-]{4,32}$", var.spec.workload_identity_pool_provider_id))
    error_message = "workload_identity_pool_provider_id must be 4-32 characters of lowercase letters, digits, or hyphens."
  }

  validation {
    condition     = !startswith(var.spec.workload_identity_pool_provider_id, "gcp-")
    error_message = "the prefix 'gcp-' is reserved by Google — choose a provider ID that does not start with it."
  }

  validation {
    condition = length([
      for issuer in [var.spec.aws, var.spec.oidc, var.spec.saml, var.spec.x509] : issuer
      if issuer != null
    ]) == 1
    error_message = "exactly one issuer (aws, oidc, saml, or x509) must be configured."
  }

  validation {
    condition     = var.spec.oidc == null || contains(keys(var.spec.attribute_mapping), "google.subject")
    error_message = "OIDC providers must map google.subject in attribute_mapping (e.g. {\"google.subject\" = \"assertion.sub\"}) — GCP rejects the provider without it."
  }
}

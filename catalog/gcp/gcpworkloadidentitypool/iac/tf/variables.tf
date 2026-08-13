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
  description = "Specification for the GCP Workload Identity Pool"
  type = object({
    # The pool ID — the final component of the pool's resource name.
    # 4-32 chars of lowercase letters, digits, hyphens; the "gcp-" prefix is
    # reserved by Google. Immutable in GCP (ForceNew).
    workload_identity_pool_id = string

    # The GCP project that owns this pool. The CLI's tfvars converter
    # resolves StringValueOrRef fields to their literal string before the
    # module runs, so this arrives as a plain string.
    # If empty, the provider's default project is used (see locals.tf).
    project_id = optional(string, "")

    # Human-readable name shown in the GCP console (max 32 chars). Mutable.
    display_name = optional(string)

    # What this pool federates and who owns it (max 256 chars). Mutable.
    description = optional(string)

    # Kill switch: a disabled pool rejects all token exchanges and existing
    # tokens stop granting access; re-enabling restores them. Mutable.
    disabled = optional(bool, false)

    # Pool operating mode: FEDERATION_ONLY (default), TRUST_DOMAIN, or
    # SYSTEM_TRUST_DOMAIN. Immutable in the API — the update fails server-side.
    mode = optional(string)

    # mTLS workload-certificate issuance for the pool's trust domain
    # (TRUST_DOMAIN pools).
    inline_certificate_issuance_config = optional(object({
      # Region -> CA Service pool resource path used for issuance in that
      # region. Exactly one of ca_pools or use_default_shared_ca supplies
      # the signing authority (spec-enforced).
      ca_pools = optional(map(string), {})
      # Issue from the GCP-provisioned default shared CA instead of your
      # own CA pools.
      use_default_shared_ca = optional(bool, false)
      # Certificate key algorithm; server-defaults to ECDSA_P256.
      key_algorithm = optional(string)
      # Certificate lifetime like "86400s" (24h) .. "2592000s" (30d);
      # server-defaults to 86400s.
      lifetime = optional(string)
      # Percent of remaining lifetime at which rotation begins (50-80);
      # server-defaults to 50.
      rotation_window_percentage = optional(number)
    }))

    # Additional trust domains whose certificates this pool's trust domain
    # accepts (a domain always trusts itself).
    inline_trust_config = optional(object({
      additional_trust_bundles = list(object({
        trust_domain = string
        trust_anchors = list(object({
          pem_certificate = string
        }))
        # Additionally trust the GCP-managed regional root certificates.
        trust_default_shared_ca = optional(bool, false)
      }))
    }))

    # Which workloads may receive a managed identity from this pool
    # (max 50 rules; applied by GCP in a second API call after create).
    attestation_rules = optional(list(object({
      google_cloud_resource = string
    })), [])

    # What happens to the pool in GCP on destroy:
    # DELETE (provider default; ~30-day soft delete), PREVENT, or ABANDON.
    deletion_policy = optional(string, "")
  })

  validation {
    condition     = can(regex("^[a-z0-9-]{4,32}$", var.spec.workload_identity_pool_id))
    error_message = "workload_identity_pool_id must be 4-32 characters of lowercase letters, digits, or hyphens."
  }

  validation {
    condition     = !startswith(var.spec.workload_identity_pool_id, "gcp-")
    error_message = "the prefix 'gcp-' is reserved by Google — choose a pool ID that does not start with it."
  }

  validation {
    condition     = var.spec.mode == null || contains(["FEDERATION_ONLY", "TRUST_DOMAIN"], coalesce(var.spec.mode, "FEDERATION_ONLY"))
    error_message = "mode must be FEDERATION_ONLY or TRUST_DOMAIN (SYSTEM_TRUST_DOMAIN pools are Google-managed and cannot be created)."
  }

  validation {
    # try() guards the null case: HCL's || evaluates BOTH operands (it does
    # not short-circuit), so dereferencing the object on the right-hand side
    # would fail whenever the block is omitted.
    condition     = var.spec.inline_certificate_issuance_config == null || (length(try(var.spec.inline_certificate_issuance_config.ca_pools, {})) > 0) != try(var.spec.inline_certificate_issuance_config.use_default_shared_ca, false)
    error_message = "choose exactly one certificate authority source: ca_pools or use_default_shared_ca."
  }
}

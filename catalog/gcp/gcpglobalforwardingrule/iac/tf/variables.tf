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
  description = "Specification for the GCP Compute Engine global forwarding rule"
  type = object({
    # The GCP project that owns the rule. The CLI's tfvars converter
    # resolves StringValueOrRef fields to their literal string before the
    # module runs, so this arrives as a plain string.
    # If empty, the provider's default project is used (see locals.tf).
    project_id = optional(string, "")

    # Name of the rule in GCP (RFC1035; PSC-for-Google-APIs names are
    # 1-20 letters/digits). Empty defaults to metadata.name (see
    # locals.tf). Immutable (ForceNew).
    forwarding_rule_name = optional(string, "")

    description = optional(string, "")

    # The target receiving matched traffic — required; arrives as a plain
    # string: a proxy self-link, a PSC bundle name (all-apis / vpc-sc), or
    # a service attachment URI. The only mutable wiring on the resource
    # (setTarget swaps it in place).
    target = string

    # The VIP: a literal IP, an address resource URL, or empty for an
    # ephemeral Google-assigned IP. Arrives as a plain string (a
    # GcpGlobalAddress ref resolves to its literal IP). Immutable.
    ip_address = optional(string, "")

    # IP protocol (middleware defaults it to TCP — GCP's own default).
    # Immutable.
    ip_protocol = optional(string, "")

    # IP version for an auto-assigned ephemeral address (IPV4/IPV6); only
    # meaningful when ip_address is empty. Immutable.
    ip_version = optional(string, "")

    # Load balancer family. Middleware defaults it to EXTERNAL; the spec's
    # NONE sentinel selects the Private Service Connect form, which the
    # API expects as an EMPTY scheme (see locals.tf). Immutable except the
    # EXTERNAL → EXTERNAL_MANAGED canary migration.
    load_balancing_scheme = optional(string, "")

    # Port or contiguous range ("443", "8080-8090"); how the port-80
    # redirect rule and the port-443 serving rule share one VIP. Not used
    # by PSC rules. Immutable.
    port_range = optional(string, "")

    # VPC network (internal schemes + PSC only; arrives as a plain
    # self-link string). Immutable.
    network = optional(string, "")

    # Subnetwork for internal load balancing (plain self-link string).
    # Immutable.
    subnetwork = optional(string, "")

    # Networking tier; global rules are PREMIUM-only (spec CEL enforces
    # it). Empty keeps GCP's computed default. Immutable.
    network_tier = optional(string, "")

    # Traffic Director xDS scoping filters (INTERNAL_SELF_MANAGED only).
    # Immutable.
    metadata_filters = optional(list(object({
      filter_match_criteria = string
      filter_labels = list(object({
        name  = string
        value = string
      }))
    })), [])

    # Service Directory registration for PSC-for-Google-APIs frontends.
    # Immutable.
    service_directory_registration = optional(object({
      namespace                = optional(string, "")
      service_directory_region = optional(string, "")
    }))

    # Skip the auto-created PSC DNS zone (PSC only). Immutable.
    no_automate_dns_zone = optional(bool, false)

    # Labels to organize and bill this rule. Mutable.
    labels = optional(map(string), {})

    # EXTERNAL → EXTERNAL_MANAGED backend-bucket canary migration state
    # (PREPARE / TEST_BY_PERCENTAGE / TEST_ALL_TRAFFIC) and the traffic
    # fraction for TEST_BY_PERCENTAGE. Mutable.
    external_managed_backend_bucket_migration_state              = optional(string, "")
    external_managed_backend_bucket_migration_testing_percentage = optional(number, 0)

    # DELETE (default), PREVENT, or ABANDON — what destroy does to the
    # forwarding rule.
    deletion_policy = optional(string, "")
  })
}

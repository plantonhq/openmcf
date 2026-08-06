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
  description = "Specification for the GCP Compute Engine global target HTTPS proxy"
  type = object({
    # The GCP project that owns the proxy. The CLI's tfvars converter
    # resolves StringValueOrRef fields to their literal string before the
    # module runs, so this arrives as a plain string.
    # If empty, the provider's default project is used (see locals.tf).
    project_id = optional(string, "")

    # Name of the proxy in GCP (RFC1035). Empty defaults to metadata.name
    # (see locals.tf). Immutable (ForceNew).
    proxy_name = optional(string, "")

    description = optional(string, "")

    # The URL map the proxy routes through — required; arrives as a plain
    # string (self-link or name). Mutable in place (setUrlMap).
    url_map = string

    # Compute Engine SSL certificates (1-15); each StringValueOrRef arrives
    # as a plain self-link string. Exactly one certificate mechanism may be
    # used (enforced by the spec's CEL upstream). Mutable in place
    # (setSslCertificates) — rotation is attach-new-then-detach-old.
    ssl_certificates = optional(list(string), [])

    # Certificate Manager certificates (cross-region internal ALB only);
    # arrive as plain resource-name strings. Mutable.
    certificate_manager_certificates = optional(list(string), [])

    # Certificate Manager certificate map URI (external ALBs only; the
    # SNI-scale mechanism). Mutable.
    certificate_map = optional(string, "")

    # SSL policy self-link constraining TLS versions/ciphers; empty keeps
    # GCP's permissive default policy. Mutable in place (setSslPolicy).
    ssl_policy = optional(string, "")

    # Network security ServerTlsPolicy resource name — the mTLS lever, and
    # Traffic Director's only TLS mechanism. Mutable and clearable.
    server_tls_policy = optional(string, "")

    # QUIC (HTTP/3) negotiation: NONE (GCP decides), ENABLE, or DISABLE.
    # Middleware defaults it to NONE, which matches GCP's own default.
    quic_override = optional(string, "")

    # TLS 1.3 0-RTT early data: STRICT, PERMISSIVE, UNRESTRICTED, or
    # DISABLED; empty keeps GCP's default (DISABLED). Immutable.
    tls_early_data = optional(string, "")

    # Idle client keep-alive in seconds (5-1200); 0 lets GCP apply its
    # default. Only honored by EXTERNAL_MANAGED load balancers. Immutable.
    http_keep_alive_timeout_sec = optional(number, 0)

    # Traffic Director binding; false for internet-facing frontends.
    # Immutable.
    proxy_bind = optional(bool, false)
  })
}

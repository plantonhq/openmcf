locals {
  # Reference fields (domain, value) arrive flattened as plain strings: the
  # Planton orchestrator resolves valueFrom references before Terraform runs.
  domain = var.spec.domain
  value  = var.spec.value

  # The spec's enum value names ARE the DigitalOcean record types (A, AAAA,
  # CNAME, MX, TXT, SRV, NS, CAA, SOA), so the type wires through directly.
  type = var.spec.type
}

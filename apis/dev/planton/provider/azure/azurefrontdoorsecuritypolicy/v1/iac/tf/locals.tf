locals {
  # No Azure tags: ARM does not support tags on Front Door security
  # policies, so the platform's identity tags live on the profile.

  # The spec's domain list, defensive-coalesced for the provider's
  # for_each (the spec requires at least one entry).
  domain_ids = coalesce(var.spec.domain_ids, [])
}

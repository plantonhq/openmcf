locals {
  # No Azure tags: ARM does not support tags on Front Door secrets, so
  # the platform's identity tags live on the profile.

  # No derived values: the secret is a thin immutable wrapper -- profile
  # parent, name, and the wrapped Key Vault certificate id, all passed
  # through as-is.
}

locals {
  # The CMK configuration carries no Azure tags: it is a property of the
  # namespace, not an ARM object of its own, so the platform's identity
  # tags live on the namespace (and its cluster).

  # The optional unwrapping identity: the wire may carry an empty string
  # for an unset optional reference, and sending "" would be rejected --
  # null omits the argument so Azure falls back to the namespace's
  # system-assigned identity.
  user_assigned_identity_id = (
    var.spec.user_assigned_identity_id != null && var.spec.user_assigned_identity_id != ""
    ? var.spec.user_assigned_identity_id
    : null
  )
}

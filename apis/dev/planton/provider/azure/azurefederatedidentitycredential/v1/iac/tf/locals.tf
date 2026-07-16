# Federated identity credentials carry no ARM tags (they are untagged child
# resources of the identity), so the usual metadata-derived tag locals are
# intentionally absent from this module.
locals {
  # Normalize an empty audience to the token-exchange default, mirroring the
  # spec's declared default so both engines deploy identical trust rules for
  # an identical spec.
  audience = (
    var.spec.audience == null || var.spec.audience == ""
    ? "api://AzureADTokenExchange"
    : var.spec.audience
  )
}

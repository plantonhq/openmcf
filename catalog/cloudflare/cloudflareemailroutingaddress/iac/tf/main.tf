# An account-scoped Email Routing destination address. Creating it sends a
# verification email to the mailbox; it is usable as a forwarding target only
# after the owner clicks the verification link.
resource "cloudflare_email_routing_address" "main" {
  account_id = var.spec.account_id
  email      = var.spec.email
  # Explicit verification-state override; empty leaves the state to the normal
  # emailed-link flow. Cloudflare permits non-admin callers only to flip a
  # verified address back to "unverified".
  # PARITY-EXCEPTION: the Pulumi module cannot send this field -- SDK v6.17.0
  # lacks it. Both engines converge when the SDK catches up.
  status = try(var.spec.status, "") != "" ? var.spec.status : null
}

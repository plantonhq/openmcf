resource "aws_route53_zone" "this" {
  # metadata.name IS the domain (ForceNew — a zone cannot be renamed).
  name          = local.zone_name
  comment       = local.comment
  force_destroy = var.spec.force_destroy

  # Reusable delegation sets are public-zone-only and conflict with the vpc
  # block — both couplings are CEL-enforced in the spec.
  delegation_set_id = local.delegation_set_id

  # Public zones only. Tri-state pass-through: AWS keeps the feature's
  # current state when the argument is absent and requires an EXPLICIT false
  # to switch it back off, so the spec's presence travels as-is (null when
  # unset, true/false when the manifest says so).
  enable_accelerated_recovery = var.spec.enable_accelerated_recovery

  # A private zone is defined by its VPC set: AWS creates the zone attached to
  # the first VPC and associates the rest; the provider manages the whole set
  # declaratively (adding/removing entries associates/disassociates).
  dynamic "vpc" {
    for_each = local.vpc_associations
    content {
      vpc_id     = vpc.value.vpc_id
      vpc_region = vpc.value.vpc_region
    }
  }

  tags = local.aws_tags
}

# DNSSEC key-signing key — the asymmetric KMS key behind the zone's signatures.
# The key must live in us-east-1 with key spec ECC_NIST_P256 and a key policy
# allowing dnssec-route53.amazonaws.com (documented on the spec).
resource "aws_route53_key_signing_key" "this" {
  count = local.dnssec_enabled ? 1 : 0

  hosted_zone_id             = aws_route53_zone.this.zone_id
  key_management_service_arn = var.spec.dnssec.kms_key_arn
  name                       = local.ksk_name
  # ACTIVE unless the spec deactivates the key (the diagnostics lever);
  # spec default and provider default agree, so this is always sent.
  status = local.ksk_status
}

# Signing only flips on after the KSK exists — an explicit dependency, not
# just ordering.
resource "aws_route53_hosted_zone_dnssec" "this" {
  count = local.dnssec_enabled ? 1 : 0

  hosted_zone_id = aws_route53_zone.this.zone_id
  signing_status = "SIGNING"

  depends_on = [aws_route53_key_signing_key.this]
}

# Query logging — public zones only (CEL-enforced). The destination log group
# must live in us-east-1 and carry an account-level CloudWatch Logs resource
# policy allowing route53.amazonaws.com; that policy is account-scoped (max 10
# per region) and deliberately NOT created per zone.
resource "aws_route53_query_log" "this" {
  count = local.query_logging_enabled ? 1 : 0

  zone_id                  = aws_route53_zone.this.zone_id
  cloudwatch_log_group_arn = var.spec.query_logging.cloudwatch_log_group_arn
}

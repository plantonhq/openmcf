# One AWS Private CA with composed activation, issued certificates,
# the ACM renewal permission, and the resource policy managed in-line.
#
# Lifecycle facts the render below depends on:
#   - a fresh CA sits in PENDING_CERTIFICATE until a certificate is
#     INSTALLED on it; the CA resource's own certificate attribute is
#     read at create (still empty), so the ca_certificate output
#     comes from the ACTIVATION path, never the CA attribute;
#   - a ROOT self-signs: its own CSR issued against ITSELF with the
#     RootCACertificate template, then installed (three provider
#     calls the raw provider makes users wire by hand);
#   - a SUBORDINATE's certificate is issued by the PARENT CA with a
#     path-length template; the issue call's signing algorithm is
#     this spec's signing_algorithm, so same-family hierarchies are
#     assumed (the norm - a cross-family hierarchy needs out-of-band
#     activation);
#   - issued certificates require the CA ACTIVE - they depend on the
#     activation resource; delete REVOKES them (the provider's
#     documented delete semantic);
#   - deleting the CA parks it restorable for
#     permanent_deletion_time_in_days; billing stops at delete.

# Template ARNs are partition-scoped
# (arn:{partition}:acm-pca:::template/...).
data "aws_partition" "this" {}

resource "aws_acmpca_certificate_authority" "this" {
  type       = var.spec.type
  usage_mode = var.spec.usage_mode != "" ? var.spec.usage_mode : null

  certificate_authority_configuration {
    key_algorithm     = var.spec.key_algorithm
    signing_algorithm = var.spec.signing_algorithm

    subject {
      common_name         = var.spec.subject.common_name != "" ? var.spec.subject.common_name : null
      organization        = var.spec.subject.organization != "" ? var.spec.subject.organization : null
      organizational_unit = var.spec.subject.organizational_unit != "" ? var.spec.subject.organizational_unit : null
      country             = var.spec.subject.country != "" ? var.spec.subject.country : null
      state               = var.spec.subject.state != "" ? var.spec.subject.state : null
      locality            = var.spec.subject.locality != "" ? var.spec.subject.locality : null
    }
  }

  key_storage_security_standard   = var.spec.key_storage_security_standard != "" ? var.spec.key_storage_security_standard : null
  permanent_deletion_time_in_days = var.spec.permanent_deletion_time_in_days > 0 ? var.spec.permanent_deletion_time_in_days : null
  enabled                         = var.spec.enabled != null ? var.spec.enabled : true

  dynamic "revocation_configuration" {
    for_each = var.spec.revocation != null ? [var.spec.revocation] : []
    content {
      dynamic "crl_configuration" {
        for_each = revocation_configuration.value.crl != null && try(revocation_configuration.value.crl.enabled, false) ? [revocation_configuration.value.crl] : []
        content {
          enabled            = true
          expiration_in_days = crl_configuration.value.expiration_in_days
          s3_bucket_name     = crl_configuration.value.s3_bucket_name
          s3_object_acl      = crl_configuration.value.s3_object_acl != "" ? crl_configuration.value.s3_object_acl : null
          custom_cname       = crl_configuration.value.custom_cname != "" ? crl_configuration.value.custom_cname : null
          custom_path        = crl_configuration.value.custom_path != "" ? crl_configuration.value.custom_path : null
        }
      }

      dynamic "ocsp_configuration" {
        for_each = revocation_configuration.value.ocsp != null && try(revocation_configuration.value.ocsp.enabled, false) ? [revocation_configuration.value.ocsp] : []
        content {
          enabled           = true
          ocsp_custom_cname = ocsp_configuration.value.custom_cname != "" ? ocsp_configuration.value.custom_cname : null
        }
      }
    }
  }

  tags = local.aws_tags
}

# ROOT self-activation: issue the CA's own CSR against itself with the
# root template, then install it.
resource "aws_acmpca_certificate" "root_ca" {
  count = local.is_root ? 1 : 0

  certificate_authority_arn   = aws_acmpca_certificate_authority.this.arn
  certificate_signing_request = aws_acmpca_certificate_authority.this.certificate_signing_request
  signing_algorithm           = var.spec.signing_algorithm
  template_arn                = "arn:${data.aws_partition.this.partition}:acm-pca:::template/RootCACertificate/V1"

  validity {
    type  = var.spec.root_ca_validity != null ? var.spec.root_ca_validity.type : "YEARS"
    value = var.spec.root_ca_validity != null ? var.spec.root_ca_validity.value : "10"
  }
}

# SUBORDINATE activation: the parent CA signs this CA's CSR with a
# path-length template.
resource "aws_acmpca_certificate" "subordinate_ca" {
  count = local.activates_subordinate ? 1 : 0

  certificate_authority_arn   = var.spec.subordinate_activation.parent_ca_arn
  certificate_signing_request = aws_acmpca_certificate_authority.this.certificate_signing_request
  signing_algorithm           = var.spec.signing_algorithm
  template_arn                = "arn:${data.aws_partition.this.partition}:acm-pca:::template/SubordinateCACertificate_PathLen${var.spec.subordinate_activation.path_length}/V1"

  validity {
    type  = var.spec.subordinate_activation.validity.type
    value = var.spec.subordinate_activation.validity.value
  }
}

# Install whichever activation certificate was issued.
resource "aws_acmpca_certificate_authority_certificate" "this" {
  count = local.has_composed_activation ? 1 : 0

  certificate_authority_arn = aws_acmpca_certificate_authority.this.arn
  certificate               = local.is_root ? aws_acmpca_certificate.root_ca[0].certificate : aws_acmpca_certificate.subordinate_ca[0].certificate
  certificate_chain         = local.is_root ? null : aws_acmpca_certificate.subordinate_ca[0].certificate_chain
}

# Certificates issued from this CA, keyed by name. Issuing needs the
# CA ACTIVE, which the activation above makes it.
resource "aws_acmpca_certificate" "issued" {
  for_each = { for certificate in var.spec.issued_certificates : certificate.name => certificate }

  certificate_authority_arn   = aws_acmpca_certificate_authority.this.arn
  certificate_signing_request = each.value.csr
  signing_algorithm           = each.value.signing_algorithm
  template_arn                = each.value.template_arn != "" ? each.value.template_arn : null
  api_passthrough             = each.value.api_passthrough != "" ? each.value.api_passthrough : null

  validity {
    type  = each.value.validity.type
    value = each.value.validity.value
  }

  depends_on = [aws_acmpca_certificate_authority_certificate.this]
}

# The ACM auto-renewal grant - all three actions, per AWS's documented
# requirement (a partial grant fails silently at renewal time).
resource "aws_acmpca_permission" "acm" {
  count = var.spec.acm_renewal_permission ? 1 : 0

  certificate_authority_arn = aws_acmpca_certificate_authority.this.arn
  principal                 = "acm.amazonaws.com"
  actions                   = ["IssueCertificate", "GetCertificate", "ListPermissions"]
}

resource "aws_acmpca_policy" "this" {
  count = var.spec.policy != "" ? 1 : 0

  resource_arn = aws_acmpca_certificate_authority.this.arn
  policy       = var.spec.policy
}

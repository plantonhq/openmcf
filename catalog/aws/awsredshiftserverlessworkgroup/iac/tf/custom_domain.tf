# Fronts the workgroup's endpoint with a custom DNS name and an ACM
# certificate (AWS keeps one custom domain per workgroup). The CNAME
# record pointing the domain at the workgroup endpoint stays yours to
# manage; certificate renewals through ACM update the association's
# expiry in place.
resource "aws_redshiftserverless_custom_domain_association" "this" {
  count = var.spec.custom_domain != null ? 1 : 0

  workgroup_name                = aws_redshiftserverless_workgroup.this.workgroup_name
  custom_domain_name            = var.spec.custom_domain.domain_name
  custom_domain_certificate_arn = var.spec.custom_domain.certificate_arn
}

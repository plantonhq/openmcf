output "certificate_authority_arn" {
  description = "The CA's ARN - the join key every consumer uses"
  value       = aws_acmpca_certificate_authority.this.arn
}

output "certificate_authority_id" {
  description = "The CA's id (the UUID tail of the ARN)"
  value       = aws_acmpca_certificate_authority.this.id
}

output "ca_certificate" {
  description = "The CA's certificate, PEM - what relying parties trust (from the composed activation; empty until an out-of-band activation installs one)"
  value       = local.is_root ? aws_acmpca_certificate.root_ca[0].certificate : (local.activates_subordinate ? aws_acmpca_certificate.subordinate_ca[0].certificate : "")
}

output "ca_certificate_chain" {
  description = "The CA's certificate chain, PEM (empty for a ROOT - it IS the anchor)"
  value       = local.activates_subordinate ? coalesce(aws_acmpca_certificate.subordinate_ca[0].certificate_chain, "") : ""
}

output "ca_csr" {
  description = "The CA's certificate signing request, PEM - what an EXTERNAL parent signs when activating a subordinate out of band"
  value       = aws_acmpca_certificate_authority.this.certificate_signing_request
}

output "issued_certificate_arns" {
  description = "Issued certificates' ARNs keyed by each entry's name - each value is that certificate's import ID"
  value       = { for name, certificate in aws_acmpca_certificate.issued : name => certificate.arn }
}

output "activation_certificate_arn" {
  description = "The ARN of the CA's own activation certificate (root self-signed XOR parent-issued) - empty until activated"
  value       = local.is_root ? aws_acmpca_certificate.root_ca[0].arn : (local.activates_subordinate ? aws_acmpca_certificate.subordinate_ca[0].arn : "")
}

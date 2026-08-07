output "service_arn" {
  description = "The full ARN of the App Runner service -- the handle IAM policies, VPC Ingress Connections, and deployment triggers reference."
  value       = aws_apprunner_service.this.arn
}

output "service_id" {
  description = "The AWS-assigned service identifier, unique within the account and region."
  value       = aws_apprunner_service.this.service_id
}

output "service_url" {
  description = "The default HTTPS endpoint of the service (scheme-less). For private services, the domain a VPC Ingress Connection resolves to."
  value       = aws_apprunner_service.this.service_url
}

output "service_name" {
  description = "The service name (metadata.name) the service was created under."
  value       = aws_apprunner_service.this.service_name
}

output "service_status" {
  description = "The service's lifecycle status at the end of the deployment (RUNNING when serving traffic)."
  value       = aws_apprunner_service.this.status
}

output "custom_domains" {
  description = "Per-domain DNS material for the associated custom domains: the dns_target to point each domain at, plus the certificate-validation CNAMEs to create (compose them into AwsRoute53DnsRecord resources). Empty when no custom domains are associated."
  value = [
    for domain_name, assoc in aws_apprunner_custom_domain_association.this : {
      domain_name = domain_name
      dns_target  = assoc.dns_target
      status      = assoc.status
      certificate_validation_records = [
        for r in assoc.certificate_validation_records : {
          record_name  = r.name
          record_type  = r.type
          record_value = r.value
        }
      ]
    }
  ]
}

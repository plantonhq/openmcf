output "peering_connection_id" {
  description = "The peering connection's id (pcx-...) - what route tables and accept-arm instances reference, and the provider's import ID"
  value = coalesce(
    one(aws_vpc_peering_connection.this[*].id),
    one(aws_vpc_peering_connection_accepter.this[*].id),
  )
}

output "accept_status" {
  description = "The connection's acceptance status after this side's deploy"
  value = coalesce(
    one(aws_vpc_peering_connection.this[*].accept_status),
    one(aws_vpc_peering_connection_accepter.this[*].accept_status),
  )
}

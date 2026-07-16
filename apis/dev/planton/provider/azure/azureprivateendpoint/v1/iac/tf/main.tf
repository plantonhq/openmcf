# Create the private endpoint -- a network interface that gives a Private
# Link-enabled service a private IP inside the subnet. The private service
# connection and (optional) DNS zone group are inline blocks on the endpoint;
# the DNS zone group is what makes the service FQDN resolve to the private IP
# instead of the public one, so it is created whenever any zone is referenced.
resource "azurerm_private_endpoint" "endpoint" {
  name                          = var.spec.name
  location                      = var.spec.region
  resource_group_name           = var.spec.resource_group
  subnet_id                     = var.spec.subnet_id
  custom_network_interface_name = var.spec.custom_network_interface_name != "" ? var.spec.custom_network_interface_name : null
  tags                          = local.final_tags

  private_service_connection {
    name                 = local.connection_name
    is_manual_connection = var.spec.private_service_connection.is_manual_connection

    # Exactly one of resource id / alias is set (spec-guaranteed); send null
    # for the unset one so both engines deploy the identical connection.
    private_connection_resource_id    = var.spec.private_service_connection.private_connection_resource_id != "" ? var.spec.private_service_connection.private_connection_resource_id : null
    private_connection_resource_alias = var.spec.private_service_connection.connection_alias != "" ? var.spec.private_service_connection.connection_alias : null
    subresource_names                 = length(var.spec.private_service_connection.subresource_names) > 0 ? var.spec.private_service_connection.subresource_names : null

    # request_message only accompanies a manual connection (spec-guaranteed
    # pairing); null when empty.
    request_message = var.spec.private_service_connection.request_message != "" ? var.spec.private_service_connection.request_message : null
  }

  dynamic "private_dns_zone_group" {
    for_each = length(var.spec.private_dns_zone_ids) > 0 ? [1] : []
    content {
      name                 = local.dns_zone_group_name
      private_dns_zone_ids = var.spec.private_dns_zone_ids
    }
  }

  # Static IP assignments; when empty the endpoint uses dynamic allocation.
  dynamic "ip_configuration" {
    for_each = var.spec.ip_configurations
    content {
      name               = ip_configuration.value.name
      private_ip_address = ip_configuration.value.private_ip_address
      subresource_name   = ip_configuration.value.subresource_name != "" ? ip_configuration.value.subresource_name : null
      member_name        = ip_configuration.value.member_name != "" ? ip_configuration.value.member_name : null
    }
  }
}

# Application security group membership is expressed member-side in Azure's
# model, as its own association resource (one per group).
resource "azurerm_private_endpoint_application_security_group_association" "asg" {
  count = length(var.spec.application_security_group_ids)

  private_endpoint_id           = azurerm_private_endpoint.endpoint.id
  application_security_group_id = var.spec.application_security_group_ids[count.index]
}

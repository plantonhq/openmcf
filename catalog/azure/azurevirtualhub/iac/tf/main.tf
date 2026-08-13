# Create the Virtual Hub -- the managed regional router of a Virtual
# WAN. A Standard hub bills (~$0.25/hr class) from creation, and ARM
# takes 15-30 minutes to bring the hub's router to a Provisioned
# routing state; deleting a hub requires its connections and gateways
# to be gone first.
resource "azurerm_virtual_hub" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # Required by the spec although the provider marks it optional: the
  # provider's schema also serves the standalone Route Server
  # construction (not modeled by this kind), and ARM rejects a WAN hub
  # without an address prefix.
  virtual_wan_id = var.spec.virtual_wan_id
  address_prefix = var.spec.address_prefix

  sku                    = local.sku
  hub_routing_preference = local.hub_routing_preference

  # Off is ARM's default; the WAN's own allow_branch_to_branch_traffic
  # must ALSO be on for branch-to-branch transit to actually flow.
  branch_to_branch_traffic_enabled = var.spec.branch_to_branch_traffic_enabled

  # ARM's floor and default is 2 routing infrastructure units; rendered
  # explicitly so the plan shows the real capacity.
  virtual_router_auto_scale_min_capacity = (
    var.spec.virtual_router_auto_scale_min_capacity == null
    ? 2
    : var.spec.virtual_router_auto_scale_min_capacity
  )

  # The hub resource's classic inline routes (applied to the default
  # route table). The modern per-table form lives in route_tables below.
  dynamic "route" {
    for_each = var.spec.routes
    content {
      address_prefixes    = route.value.address_prefixes
      next_hop_ip_address = route.value.next_hop_ip_address
    }
  }

  tags = local.final_tags
}

# The composed custom route tables: standalone ARM children of the hub,
# one per spec entry, keyed by name (how spoke isolation and
# shared-services topologies are built). Routes are managed INLINE on
# each table -- never mix in the provider's standalone route resource,
# which fights inline routes over the same ARM collection.
resource "azurerm_virtual_hub_route_table" "route_tables" {
  for_each = { for route_table in var.spec.route_tables : route_table.name => route_table }

  name           = each.value.name
  virtual_hub_id = azurerm_virtual_hub.main.id
  labels         = each.value.labels

  dynamic "route" {
    for_each = each.value.routes
    content {
      name              = route.value.name
      destinations_type = lookup(local.destinations_type_wire, route.value.destinations_type, route.value.destinations_type)
      destinations      = route.value.destinations
      next_hop          = route.value.next_hop
      # next_hop_type is left to the provider's default -- ARM's only
      # value is "ResourceId"; there is nothing to configure.
    }
  }
}

# The composed route maps: BGP route match/transform policies that
# connections reference as inbound_route_map_id / outbound_route_map_id.
resource "azurerm_route_map" "route_maps" {
  for_each = { for route_map in var.spec.route_maps : route_map.name => route_map }

  name           = each.value.name
  virtual_hub_id = azurerm_virtual_hub.main.id

  dynamic "rule" {
    for_each = each.value.rules
    content {
      name = rule.value.name

      # Unset leaves the provider's default ("Unknown": evaluation
      # stops after the match).
      next_step_if_matched = (
        rule.value.next_step_if_matched == null
        ? null
        : lookup(local.next_step_wire, rule.value.next_step_if_matched, rule.value.next_step_if_matched)
      )

      dynamic "match_criterion" {
        for_each = rule.value.match_criteria
        content {
          match_condition = lookup(local.match_condition_wire, match_criterion.value.match_condition, match_criterion.value.match_condition)
          as_path         = match_criterion.value.as_path
          community       = match_criterion.value.community
          route_prefix    = match_criterion.value.route_prefix
        }
      }

      dynamic "action" {
        for_each = rule.value.actions
        content {
          type = lookup(local.action_type_wire, action.value.type, action.value.type)

          # The spec guarantees non-DROP actions carry at least one
          # parameter (mirroring the provider's own create-time rule).
          dynamic "parameter" {
            for_each = action.value.parameters
            content {
              as_path      = parameter.value.as_path
              community    = parameter.value.community
              route_prefix = parameter.value.route_prefix
            }
          }
        }
      }
    }
  }
}

# The composed BGP peerings between the hub's router and NVAs in
# connected spokes. All fields are ForceNew on ARM's side; routes only
# flow once the peer is reachable through a hub connection.
resource "azurerm_virtual_hub_bgp_connection" "bgp_connections" {
  for_each = { for bgp_connection in var.spec.bgp_connections : bgp_connection.name => bgp_connection }

  name           = each.value.name
  virtual_hub_id = azurerm_virtual_hub.main.id
  peer_asn       = each.value.peer_asn
  peer_ip        = each.value.peer_ip

  virtual_network_connection_id = (
    each.value.virtual_network_connection_id != ""
    ? each.value.virtual_network_connection_id
    : null
  )
}

# The hub's routing intent (at most one): steers Internet/private
# traffic through a security appliance in THIS hub. Setting it takes
# over the hub's routing policy -- per-connection route-table
# customization and routing intent are mutually exclusive on ARM's side.
# Keyed by the intent's name so the resource address matches the ARM
# child segment (what the blind import round-trip derives its ID from).
resource "azurerm_virtual_hub_routing_intent" "routing_intent" {
  for_each = var.spec.routing_intent != null ? { (var.spec.routing_intent.name) = var.spec.routing_intent } : {}

  name           = each.value.name
  virtual_hub_id = azurerm_virtual_hub.main.id

  dynamic "routing_policy" {
    for_each = each.value.routing_policies
    content {
      name         = routing_policy.value.name
      destinations = [for destination in routing_policy.value.destinations : lookup(local.routing_policy_destination_wire, destination, destination)]
      next_hop     = routing_policy.value.next_hop
    }
  }
}

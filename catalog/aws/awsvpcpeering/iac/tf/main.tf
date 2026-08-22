# One side of a VPC peering connection, as a request-XOR-accept mode
# union (the spec's CEL guarantees exactly one arm).
#
# Lifecycle facts the render below depends on:
#   - vpc_id / peer_vpc_id / peer_owner_id / peer_region are fixed for
#     life (replace-on-change);
#   - auto_accept works same-account, same-region only (the provider
#     hard-errors on peer_region + auto_accept; the spec's CELs
#     front-load both walls);
#   - DNS-resolution options need an ACTIVE connection - on a pending
#     cross-account request AWS rejects the modification until
#     accepted (the accept-arm instance sets the accepter side);
#   - the ACCEPT arm's destroy is a NO-OP at AWS: it abandons
#     management without deleting the peering (only the requester side
#     deletes);
#   - one VPC pair supports at most one peering - AWS returns the
#     EXISTING connection id for a duplicate request (never declare
#     the same pair twice).

# --- request arm -----------------------------------------------------------

resource "aws_vpc_peering_connection" "this" {
  count = var.spec.request != null ? 1 : 0

  vpc_id        = var.spec.request.vpc_id
  peer_vpc_id   = var.spec.request.peer_vpc_id
  peer_owner_id = var.spec.request.peer_owner_id != "" ? var.spec.request.peer_owner_id : null
  peer_region   = var.spec.request.peer_region != "" ? var.spec.request.peer_region : null
  auto_accept   = var.spec.request.auto_accept

  # Options are managed in-line as the single owner - the standalone
  # aws_vpc_peering_connection_options resource fights this form (the
  # provider's own docs warn the two overwrite each other).
  dynamic "requester" {
    for_each = var.spec.request.requester_allow_remote_vpc_dns_resolution ? [1] : []
    content {
      allow_remote_vpc_dns_resolution = true
    }
  }

  dynamic "accepter" {
    for_each = var.spec.request.accepter_allow_remote_vpc_dns_resolution ? [1] : []
    content {
      allow_remote_vpc_dns_resolution = true
    }
  }

  tags = local.aws_tags
}

# --- accept arm ------------------------------------------------------------

resource "aws_vpc_peering_connection_accepter" "this" {
  count = var.spec.accept != null ? 1 : 0

  vpc_peering_connection_id = var.spec.accept.vpc_peering_connection_id
  auto_accept               = var.spec.accept.auto_accept

  dynamic "accepter" {
    for_each = var.spec.accept.accepter_allow_remote_vpc_dns_resolution ? [1] : []
    content {
      allow_remote_vpc_dns_resolution = true
    }
  }

  tags = local.aws_tags
}

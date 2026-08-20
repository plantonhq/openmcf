# One RDS Proxy with its default-target-group pool tuning, additional
# endpoints, and database target managed in-line.
#
# Lifecycle facts the render below depends on:
#   - engine_family, vpc_subnet_ids, and both network-type dials are
#     ForceNew on the proxy; everything else updates in place;
#   - the default target group is a PATCH satellite: it always exists
#     on the proxy, its provider delete is a no-op, and managing it
#     here just tunes the pool (name is always "default");
#   - the target registration waits out a database still in CREATING
#     (the provider retries for 5 minutes) - it is still ordered after
#     the target group so plans read cleanly;
#   - a proxy fronts exactly ONE database (instance XOR cluster - the
#     spec's CEL wall mirrors AWS's contract).

resource "aws_db_proxy" "this" {
  name          = var.metadata.name
  engine_family = var.spec.engine_family
  role_arn      = var.spec.role_arn
  vpc_subnet_ids = var.spec.vpc_subnet_ids
  vpc_security_group_ids = length(var.spec.vpc_security_group_ids) > 0 ? var.spec.vpc_security_group_ids : null

  dynamic "auth" {
    for_each = var.spec.auth
    content {
      # SECRETS is the only auth scheme AWS supports - pinned here,
      # never spec surface.
      auth_scheme                = "SECRETS"
      secret_arn                 = auth.value.secret_arn
      description                = auth.value.description != "" ? auth.value.description : null
      iam_auth                   = auth.value.iam_auth != "" ? auth.value.iam_auth : null
      client_password_auth_type  = auth.value.client_password_auth_type != "" ? auth.value.client_password_auth_type : null
      username                   = auth.value.username != "" ? auth.value.username : null
    }
  }

  require_tls                    = var.spec.require_tls ? true : null
  idle_client_timeout            = var.spec.idle_client_timeout > 0 ? var.spec.idle_client_timeout : null
  debug_logging                  = var.spec.debug_logging ? true : null
  default_auth_scheme            = var.spec.default_auth_scheme != "" ? var.spec.default_auth_scheme : null
  endpoint_network_type          = var.spec.endpoint_network_type != "" ? var.spec.endpoint_network_type : null
  target_connection_network_type = var.spec.target_connection_network_type != "" ? var.spec.target_connection_network_type : null

  tags = local.aws_tags
}

# Always managed: with no pool tuning configured AWS keeps its
# defaults, and managing the group gives the target a concrete
# ordering point plus the ARN output either way.
resource "aws_db_proxy_default_target_group" "this" {
  db_proxy_name = aws_db_proxy.this.name

  dynamic "connection_pool_config" {
    for_each = var.spec.connection_pool != null ? [var.spec.connection_pool] : []
    content {
      connection_borrow_timeout    = connection_pool_config.value.connection_borrow_timeout
      init_query                   = connection_pool_config.value.init_query != "" ? connection_pool_config.value.init_query : null
      max_connections_percent      = connection_pool_config.value.max_connections_percent
      max_idle_connections_percent = connection_pool_config.value.max_idle_connections_percent
      session_pinning_filters      = length(connection_pool_config.value.session_pinning_filters) > 0 ? connection_pool_config.value.session_pinning_filters : null
    }
  }
}

# Additional endpoints, keyed by name.
resource "aws_db_proxy_endpoint" "this" {
  for_each = { for endpoint in var.spec.endpoints : endpoint.name => endpoint }

  db_proxy_name          = aws_db_proxy.this.name
  db_proxy_endpoint_name = each.value.name
  vpc_subnet_ids         = each.value.vpc_subnet_ids
  vpc_security_group_ids = length(each.value.vpc_security_group_ids) > 0 ? each.value.vpc_security_group_ids : null
  target_role            = each.value.target_role != "" ? each.value.target_role : null

  tags = local.aws_tags
}

# The one database this proxy fronts (instance XOR cluster).
resource "aws_db_proxy_target" "this" {
  count = var.spec.target != null ? 1 : 0

  db_proxy_name     = aws_db_proxy.this.name
  target_group_name = aws_db_proxy_default_target_group.this.name

  db_instance_identifier = var.spec.target.db_instance_identifier != "" ? var.spec.target.db_instance_identifier : null
  db_cluster_identifier  = var.spec.target.db_cluster_identifier != "" ? var.spec.target.db_cluster_identifier : null
}

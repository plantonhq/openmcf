locals {
  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsRestApiGateway"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # The resolved stage name ("prod" when the stage is omitted or
  # unnamed).
  stage_name = var.spec.stage != null && try(var.spec.stage.name, "") != "" ? var.spec.stage.name : "prod"

  # Routes keyed by "METHOD /path" (the for_each keys both engines
  # share).
  routes = { for r in var.spec.routes : "${r.method} ${r.path}" => r }

  # Every tree node the route paths imply ("/users/{id}" implies
  # "/users" and "/users/{id}"). The root path "/" needs no resource.
  resource_paths = toset(flatten([
    for r in var.spec.routes : [
      for i in range(1, length(split("/", trimprefix(r.path, "/"))) + 1) :
      "/${join("/", slice(split("/", trimprefix(r.path, "/")), 0, i))}"
    ] if r.path != "/"
  ]))

  # Nodes bucketed by depth: the tree renders level-by-level (a
  # for_each resource cannot reference itself), and the spec caps
  # paths at five segments to keep the levels bounded.
  paths_level1 = toset([for p in local.resource_paths : p if length(split("/", trimprefix(p, "/"))) == 1])
  paths_level2 = toset([for p in local.resource_paths : p if length(split("/", trimprefix(p, "/"))) == 2])
  paths_level3 = toset([for p in local.resource_paths : p if length(split("/", trimprefix(p, "/"))) == 3])
  paths_level4 = toset([for p in local.resource_paths : p if length(split("/", trimprefix(p, "/"))) == 4])
  paths_level5 = toset([for p in local.resource_paths : p if length(split("/", trimprefix(p, "/"))) == 5])

  # Resource IDs keyed by path, all levels merged (methods and outputs
  # read this).
  resource_id_by_path = merge(
    { for p, r in aws_api_gateway_resource.level1 : p => r.id },
    { for p, r in aws_api_gateway_resource.level2 : p => r.id },
    { for p, r in aws_api_gateway_resource.level3 : p => r.id },
    { for p, r in aws_api_gateway_resource.level4 : p => r.id },
    { for p, r in aws_api_gateway_resource.level5 : p => r.id },
  )

  # Named satellites keyed by name (route references resolve through
  # these).
  models             = { for m in var.spec.models : m.name => m }
  request_validators = { for v in var.spec.request_validators : v.name => v }
  authorizers        = { for a in var.spec.authorizers : a.name => a }
  gateway_responses  = { for g in var.spec.gateway_responses : g.response_type => g }

  # Route responses flattened to "METHOD /path|status" (one method
  # response + its integration response mapping per entry).
  route_responses = merge([
    for key, r in local.routes : {
      for resp in r.responses : "${key}|${resp.status_code}" => {
        route_key = key
        route     = r
        response  = resp
      }
    }
  ]...)

  # Import-derivation echo maps: the route and route-response for_each
  # keys are COMPOSITE ("GET /orders/{id}", "GET /orders/{id}|200"), so
  # the blind import path needs each key's resource id, method, and
  # status code as individually addressable map entries keyed exactly
  # like the instances (see iac/import-map.yaml).
  route_resource_ids = {
    for key, r in local.routes : key => r.path == "/" ? aws_api_gateway_rest_api.this.root_resource_id : local.resource_id_by_path[r.path]
  }
  route_methods = { for key, r in local.routes : key => r.method }
  response_resource_ids = {
    for key, rr in local.route_responses : key => rr.route.path == "/" ? aws_api_gateway_rest_api.this.root_resource_id : local.resource_id_by_path[rr.route.path]
  }
  response_methods      = { for key, rr in local.route_responses : key => rr.route.method }
  response_status_codes = { for key, rr in local.route_responses : key => rr.response.status_code }

  # Documentation parts keyed by declaration position (locations are
  # composite; position is the stable cross-engine key).
  documentation_parts = var.spec.documentation != null ? { for i, p in var.spec.documentation.parts : tostring(i) => p } : {}

  # Method settings keyed by method path.
  method_settings = var.spec.stage != null ? { for m in var.spec.stage.method_settings : m.method_path => m } : {}

  # The deployment trigger fingerprints the API DEFINITION (everything
  # except the stage and documentation): any change redeploys - the
  # declarative behavior REST APIs' explicit-snapshot model would
  # otherwise lose. Each engine hashes its own canonical rendering;
  # redeploy-on-change behavior is what parity requires.
  definition_hash = sha1(jsonencode([
    var.spec.routes,
    var.spec.openapi,
    var.spec.models,
    var.spec.request_validators,
    var.spec.authorizers,
    var.spec.gateway_responses,
    var.spec.policy,
    var.spec.binary_media_types,
    var.spec.minimum_compression_size,
    var.spec.api_key_source,
    var.spec.endpoint_configuration,
    var.spec.endpoint_access_mode,
    var.spec.security_policy,
    var.spec.disable_execute_api_endpoint,
  ]))
}

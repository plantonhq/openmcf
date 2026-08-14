# AwsRestApiGateway — Terraform/OpenTofu module

Deploys an Amazon API Gateway REST API (`aws_api_gateway_rest_api`)
plus the derived resource/method/integration tree, an explicit
deployment, one stage, and the API-scoped satellites.

Module facts worth knowing before editing:

- **Typed `routes` XOR `openapi`.** The spec enforces exactly one; the
  module renders the matching arm.
- **The resource tree is derived from paths** (max five segments).
  `for_each` keys are `"${method} ${path}"`. Parent resources are
  created even when no route lands on the intermediate path.
- **Deployment trigger is `local.definition_hash`.** REST APIs do not
  auto-deploy; hashing the definition is how every spec change
  redeploys. The deployment uses `create_before_destroy`.
- **A standalone `aws_api_gateway_rest_api_policy`** owns the resource
  policy so it can be added or removed without replacing the API.
- **One stage**, default name `prod`. Canary traffic shifting is not
  modeled.

Outputs mirror the Pulumi module key-for-key: `rest_api_id`,
`rest_api_arn`, `execution_arn`, `root_resource_id`, `stage_name`,
`stage_arn`, `invoke_url`, `deployment_id`, optional client-certificate
fields, and the name-keyed satellite ID maps.

# AwsRestApiGateway — Pulumi module (Go)

Deploys an Amazon API Gateway REST API (`apigateway.RestApi`) plus the
derived resource/method/integration tree, an explicit deployment, one
stage, and the API-scoped satellites.

Module facts worth knowing before editing:

- **Typed `routes` XOR `openapi`.** The spec enforces exactly one; the
  module renders the matching arm.
- **`tree.go` derives the resource tree from paths** (max five
  segments) and wires methods, integrations, and method responses onto
  the leaves. Keep the Terraform `locals.tf` derivation in lockstep.
- **Deployment trigger hashes the full API definition** so every spec
  change redeploys. REST APIs do not auto-deploy.
- **Satellites live in `satellites.go`** (authorizers, models,
  validators, gateway responses, documentation, client certificate)
  and are referenced by name from routes.
- **One stage**, default name `prod`. Canary traffic shifting is not
  modeled.

Outputs mirror the Terraform module key-for-key: `rest_api_id`,
`rest_api_arn`, `execution_arn`, `root_resource_id`, `stage_name`,
`stage_arn`, `invoke_url`, `deployment_id`, optional client-certificate
fields, and the name-keyed satellite ID maps.

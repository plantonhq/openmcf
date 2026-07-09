# GcpVertexAiDeployedIndex Terraform Module

This directory contains the Terraform module for deploying a GCP Vertex AI Vector Search index onto an index endpoint.

## Usage

```hcl
module "vertex_ai_deployed_index" {
  source = "./path/to/module"

  metadata = {
    name = "my-deployed-index"
  }

  spec = {
    location          = "us-central1"
    deployed_index_id = "products_v1"
    index             = "projects/my-project/locations/us-central1/indexes/1234567890"
    index_endpoint    = "projects/my-project/locations/us-central1/indexEndpoints/9876543210"
  }
}
```

## Inputs

| Name | Type | Required | Description |
|------|------|----------|-------------|
| metadata | object | yes | Planton resource metadata |
| spec | object | yes | GcpVertexAiDeployedIndex specification |
| provider_config | object | no | GCP provider configuration |

## Outputs

| Name | Description |
|------|-------------|
| name | Name of the DeployedIndex resource |
| deployed_index_id | The user-chosen deployment handle |
| create_time | RFC3339 creation timestamp |
| index_sync_time | Timestamp up to which the deployment reflects the source index's updates |
| match_grpc_address | Private gRPC query address (peered endpoints only) |
| service_attachment | PSC service attachment (PSC endpoints only) |
| index_endpoint | Full resource path of the endpoint this deployment lives on |

## Mutability

Only the replica bounds inside the sizing arm update in place (the provider
PATCHes them via `mutateDeployedIndex`). Everything else — including,
unusually, `display_name` — replaces the deployment (undeploy + redeploy).

## No Labels, No Project

The GCP API gives this resource class no labels and no project field (the
deployment lives inside the endpoint resource). The module therefore applies
no platform attribution labels — none can exist here.

## Provider Requirements

- `hashicorp/google` ~> 6.0

Note: deploy timeouts are 45 minutes (create/update) and 20 (delete) —
deploys genuinely take tens of minutes.

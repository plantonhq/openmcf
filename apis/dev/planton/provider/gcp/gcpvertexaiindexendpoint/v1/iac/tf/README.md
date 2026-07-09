# GcpVertexAiIndexEndpoint Terraform Module

This directory contains the Terraform module for provisioning a GCP Vertex AI Vector Search index endpoint.

## Usage

```hcl
module "vertex_ai_index_endpoint" {
  source = "./path/to/module"

  metadata = {
    name = "my-index-endpoint"
  }

  spec = {
    location                = "us-central1"
    display_name            = "My Index Endpoint"
    public_endpoint_enabled = true
  }
}
```

`spec.project_id` is optional: when empty, the endpoint lands in the provider's
default project.

## Inputs

| Name | Type | Required | Description |
|------|------|----------|-------------|
| metadata | object | yes | Planton resource metadata |
| spec | object | yes | GcpVertexAiIndexEndpoint specification |
| provider_config | object | no | GCP provider configuration |

## Outputs

| Name | Description |
|------|-------------|
| index_endpoint_id | Fully qualified endpoint resource path (the deployed index's composition key) |
| index_endpoint_name | The GCP-assigned numeric endpoint ID |
| public_endpoint_domain_name | Public query domain (public arm only) |
| create_time | RFC3339 creation timestamp |
| update_time | RFC3339 last-update timestamp |

## Connectivity

Three mutually exclusive, immutable modes:

1. **Public** -- Set `spec.public_endpoint_enabled = true`
2. **VPC-peered** -- Set `spec.network` (self-links are normalized to the
   API's relative `projects/{project}/global/networks/{name}` form)
3. **Private Service Connect** -- Set `spec.private_service_connect_config`

## Provider Requirements

- `hashicorp/google` ~> 6.0

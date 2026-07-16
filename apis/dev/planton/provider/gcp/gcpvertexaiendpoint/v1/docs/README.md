# GcpVertexAiEndpoint: Design Notes

## Service Overview

Vertex AI Endpoints are the serving layer in Google Cloud's AI Platform. An endpoint is a stable URL that acts as a gateway for model prediction requests. It decouples the infrastructure (the endpoint itself, its networking, and encryption) from the operational concern of deploying and routing traffic to specific models.

### What an Endpoint Does

- Provides a stable prediction URL that persists across model deployments
- Manages traffic splitting between multiple deployed model versions
- Handles auto-scaling of prediction serving infrastructure
- Supports private networking (VPC peering and PSC) for secure inference
- Enables monitoring and logging of prediction requests

### What an Endpoint Does NOT Do

- It does not train models
- It does not manage model artifacts (that's Model Registry)
- It does not run batch predictions (that's BatchPredictionJob)
- Model deployment to the endpoint is a separate operational step

## Deployment Landscape

### Terraform: `google_vertex_ai_endpoint`

```hcl
resource "google_vertex_ai_endpoint" "endpoint" {
  name         = "1234567890"
  display_name = "My Endpoint"
  location     = "us-central1"
  network      = "projects/12345/global/networks/my-vpc"

  encryption_spec {
    kms_key_name = "projects/.../cryptoKeys/my-key"
  }
}
```

Key characteristics:
- `name` is Required, numeric-only, max 10 digits (no leading zeros) -- the API never generates it
- `network` and `private_service_connect_config` are mutually exclusive
- `dedicated_endpoint_enabled` conflicts with PSC
- Labels are supported (non-authoritative)
- `encryption_spec`, `location`, `network`, and `name` are ForceNew (immutable after creation)

### Pulumi: `vertex.AiEndpoint`

```go
endpoint, _ := vertex.NewAiEndpoint(ctx, "endpoint", &vertex.AiEndpointArgs{
    Name:        pulumi.StringPtr("1234567890"),
    DisplayName: pulumi.String("My Endpoint"),
    Location:    pulumi.String("us-central1"),
})
```

Same schema and mutual-exclusion rules as Terraform. `Name` must always be sent explicitly: leaving it unset triggers engine auto-naming, which produces a non-numeric name the API rejects.

## Feature Coverage

| Feature | Coverage |
|---------|----------|
| Display name + description | Full |
| Location (region) | Full |
| VPC peering via `network` | Full -- references GcpVpcNetwork |
| Private Service Connect (allowlist + secure/IAM-authorized mode) | Full |
| CMEK encryption | Full -- references GcpKmsKey |
| Dedicated endpoint DNS | Full |
| Request/response logging to BigQuery | Full |
| User labels | Full -- merged beneath platform labels |
| Endpoint name (optional, identity-derived when omitted) | Full |

### Deliberate Exclusions (with reasons)

| Provider surface | Reason |
|------------------|--------|
| `traffic_split` | A JSON map keyed by deployed-model IDs that only exist after models are deployed -- an operational step outside this component. Managing it from infrastructure state would clobber operational traffic changes and perma-diff against reality. Revisit if model deployment itself becomes a first-class kind. |
| `region` | The provider carries both `location` (required) and `region` (optional) for the same axis. `region` is not vestigial -- the provider resolves the regional API host (`https://{region}-aiplatform.googleapis.com`) from it -- so both modules pin `region` to the spec's `location` internally. One honest spec field wins; the plumbing detail stays in the modules. |
| `private_service_connect_config.psc_automation_configs` | Not in the released 6.x provider line (newer-line surface only). |
| `model_deployment_monitoring_job` output | Only populated after model deployment; always empty from IaC. |

## Endpoint Name Derivation

The API requires a numeric endpoint ID and never generates one. When the spec omits `endpoint_name`, both IaC modules derive the same stable ID from the resource's identity:

1. Build the identity string `"{org}/{env}/{name}"` from resource metadata.
2. Take the first 12 hex characters (48 bits) of its SHA-256 digest.
3. Map into `[1000000000, 9999999999]` -- always 10 digits, never a leading zero.

The derivation is implemented identically in the Terraform module (`locals.tf`) and the Pulumi module, so the same manifest yields the same endpoint ID on either engine, and re-applies never regenerate it. The resolved value is exported as the `endpoint_name` stack output.

**Reservation caveat (live-verified):** GCP reserves a deleted endpoint's numeric ID -- creating an endpoint with a previously deleted ID fails with `409 ALREADY_EXISTS` even after the endpoint is fully gone (the GET returns 404). Deterministic derivation therefore means destroy-then-recreate of the *same resource identity* collides with its own ghost. This is the honest trade-off: stable IDs across engines and re-applies, at the cost of not being able to immediately recreate a destroyed endpoint under the same identity. The escape hatches are an explicit `endpoint_name` or a changed resource name.

## Design Decisions

1. **Flattened `encryption_spec`** to `kms_key_name` -- a single field doesn't need a wrapper message
2. **PSC as sub-message** (not flat bool) -- `project_allowlist` and the secure-PSC flag are essential access-control levers
3. **`endpoint_name` optional with identity-based derivation** -- the API demands a numeric ID; deriving it from resource identity keeps manifests clean while staying deterministic across engines
4. **Request/response logging destination as a plain string** -- the `bq://` URI scheme has no matching stack output on the BigQuery kinds, so a reference shape would never resolve; the accepted URI forms are documented on the field

## Networking Deep Dive

### Public Endpoints (Default)

Requests go to `{region}-aiplatform.googleapis.com`. Simple but shared infrastructure.

### VPC-Peered Endpoints

The endpoint gets a private IP within the peered VPC. Requires:
- Private Services Access configured on the VPC
- `network` field set to the fully qualified network path

### Private Service Connect

Modern approach using PSC service attachments. Benefits:
- No VPC peering required
- Fine-grained access control via `project_allowlist`
- Optional IAM authorization on connections (`enable_secure_private_service_connect`)
- Strongest network isolation
- Cannot combine with `dedicated_endpoint_enabled`

## References

- [Vertex AI Endpoints Documentation](https://cloud.google.com/vertex-ai/docs/predictions/overview)
- [Private Endpoints](https://cloud.google.com/vertex-ai/docs/predictions/using-private-service-connect)
- [Terraform Resource](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/vertex_ai_endpoint)
- [Pulumi Resource](https://www.pulumi.com/registry/packages/gcp/api-docs/vertex/aiendpoint/)

# Regional Mapped Domain

This preset fronts a REST API stage at `api.example.com` with a
REGIONAL custom domain and a root base-path mapping.

## When to Use

- Production hostnames for a REST API
- The first custom domain in an environment (REGIONAL is the default
  for new work)

## What You Get

- A REGIONAL domain bound to the named ACM certificate
- A root mapping onto the named AwsRestApiGateway's `prod` stage
- `regional_domain_name` / `regional_zone_id` outputs for the Route 53
  alias

## Customize

- Replace `api.example.com` with your hostname (it must match the
  certificate)
- Point `certificateArn` and `restApiId` at your resources
- Add more `basePathMappings` entries with a `basePath` (for example
  `v1`) when one hostname fronts several APIs

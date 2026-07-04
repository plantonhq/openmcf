# Container-Based Lambda Function

This preset creates a Lambda function from a container image in ECR. The function name comes from `metadata.name`. Runtime and entrypoint are defined by the image — leave `runtime` and `handler` empty.

## When to Use

- Dependencies exceeding the 250 MB zip limit (ML models, large SDKs, native binaries)
- Custom OS packages or runtimes not available as managed Lambda runtimes
- Teams with existing Docker build pipelines who want the same image locally and in Lambda
- Images up to 10 GB (vs 250 MB for zip deployments)

## Key Configuration Choices

- **Container image** (`image_uri`) — ECR URI; runtime and handler come from the image CMD/ENTRYPOINT (override with `image_config` if needed)
- **512 MB memory** (`memory_size_mb: 512`) — container functions often need more headroom than zip deployments
- **60-second timeout** (`timeout_seconds: 60`) — longer default for heavier container workloads
- **ARM64** (`architecture: arm64`) — Graviton is typically ~20% cheaper; switch to `x86_64` if required
- **Composed execution role** (`role_arn.valueFrom`) — references an `AwsIamRole` with `AWSLambdaBasicExecutionRole`
- **No runtime/handler** — required empty for container deployments (CEL-enforced)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aws-region>` | AWS region for the function and ECR repository | Your deployment region |
| `<execution-role-name>` | Name of the `AwsIamRole` resource for execution | `AwsIamRole` metadata.name |
| `<ecr-image-uri>` | ECR image URI (e.g. `123456789012.dkr.ecr.us-west-2.amazonaws.com/my-function:latest`) | ECR console or your image pipeline |

## Related Presets

- **01-zip-basic** — use instead for lightweight zip-based functions with managed runtimes

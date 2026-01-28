# Overview

The **AwsCloudFront** API resource simplifies the deployment and management of AWS CloudFront distributions, providing a globally-distributed CDN for static and dynamic content delivery. By defining an opinionated, minimal configuration that covers 80% of production use cases, **AwsCloudFront** eliminates the complexity of managing multi-region certificates, origin access controls, and cache behaviors. It is part of OpenMCF's broader multi-cloud deployment framework, supporting both Pulumi and Terraform under the hood to fit seamlessly into your existing infrastructure-as-code workflows.

With a single YAML manifest conforming to the familiar Kubernetes-like structure (`apiVersion`, `kind`, `metadata`, `spec`, `status`), teams can validate, provision, and maintain CloudFront distributions consistently across multiple environments. Whether you're serving a static website from S3, routing traffic to an Application Load Balancer, or deploying a multi-origin CDN with custom domains, **AwsCloudFront** handles the essential configuration for you—including HTTPS enforcement, cost-optimized edge locations, and validation of the critical us-east-1 certificate requirement.

---

## Key Features

- **Opinionated 80/20 Configuration**  
  Exposes the essential fields needed for 80% of production CloudFront deployments: origins, custom domains, ACM certificates, price classes, and default root objects. Advanced features remain accessible via direct Pulumi or Terraform modules.

- **Multi-Origin Support with Default Selection**  
  Configure multiple origins (S3 buckets, ALBs, custom endpoints) with one marked as the default target. The API validates that exactly one origin is marked as default, preventing misconfigurations.

- **Automatic us-east-1 Certificate Validation**  
  The protobuf API enforces that ACM certificate ARNs must be in the `us-east-1` region—a common CloudFront gotcha. This validation happens before deployment, not at apply-time.

- **Cost-Optimized Defaults**  
  Defaults to `PRICE_CLASS_100` (North America and Europe only), minimizing costs for applications serving primarily NA/EU users. Users must explicitly opt into global price classes.

- **HTTPS-Only and Secure by Default**  
  All cache behaviors enforce `redirect-to-https` for viewer protocol policy, and TLSv1.2 is the minimum protocol version. Custom origins use `https-only` origin protocol policy.

- **Consistent Resource Model**  
  Uses OpenMCF's standard resource layout with comprehensive protobuf validations, including unique alias validation, certificate requirement when aliases are set, and origin path pattern enforcement.

- **Pulumi & Terraform Integration**  
  Provision the same CloudFront distribution specification using either Pulumi or Terraform, unified by the OpenMCF CLI's straightforward commands and orchestration.

---

## Module Architecture

The Pulumi implementation for **AwsCloudFront** follows a clean, modular structure:

### File Organization

```
iac/pulumi/
├── main.go              # Entrypoint that invokes the module
├── Pulumi.yaml          # Pulumi project metadata
├── Makefile             # Build and test automation
├── debug.sh             # Delve debugger integration script
├── README.md            # CLI usage and debugging guide
├── overview.md          # This file
└── module/
    ├── main.go          # Main Resources() function, provider setup, output exports
    ├── locals.go        # Locals struct for convenient access to stack input
    ├── outputs.go       # Output key constants matching AwsCloudFrontStackOutputs proto
    └── distribution.go  # Core logic for creating the CloudFront distribution
```

### Key Components

#### 1. `main.go` (Entrypoint)
The top-level `main.go` initializes the Pulumi program, deserializes the stack input from the OpenMCF CLI, and invokes the `module.Resources()` function.

#### 2. `module/main.go` (Orchestration)
The `Resources()` function orchestrates the entire deployment:
- Initializes `Locals` struct from the stack input for convenient access to spec fields
- Creates the AWS provider using credentials from `ProviderConfig` or defaults to the environment
- Invokes `createDistribution()` to provision the CloudFront distribution
- Exports outputs (`distribution_id`, `domain_name`, `hosted_zone_id`) mapped to the protobuf `AwsCloudFrontStackOutputs` schema

#### 3. `module/distribution.go` (Core Logic)
This is the heart of the implementation. It:
- Iterates over the `origins` array from the spec, constructing Pulumi `DistributionOriginArgs` for each
- Identifies the origin marked as `is_default` and assigns it as the `targetOriginId` for the default cache behavior
- Applies secure defaults: `https-only` origin protocol policy, `redirect-to-https` viewer protocol policy, TLSv1.2 minimum version
- Configures the viewer certificate: if `certificate_arn` is provided, uses SNI-only ACM certificate; otherwise, uses the CloudFront default certificate
- Maps the protobuf `PriceClass` enum to CloudFront's `PriceClass_100`, `PriceClass_200`, or `PriceClass_All` values
- Creates the `cloudfront.Distribution` resource with all configured parameters

#### 4. `module/locals.go` (Convenience)
Defines the `Locals` struct that provides shorthand access to `Target` (the full `AwsCloudFront` resource) and `Spec` (the `AwsCloudFrontSpec` configuration). This eliminates repetitive `stackInput.Target.Spec` references throughout the code.

#### 5. `module/outputs.go` (Output Constants)
Defines string constants for Pulumi output keys that match the protobuf `AwsCloudFrontStackOutputs` schema:
- `OpDistributionId`: CloudFront distribution ID (e.g., `E2QWRUHAPOMQZL`)
- `OpDomainName`: CloudFront domain name (e.g., `d123.cloudfront.net`)
- `OpHostedZoneId`: Route53 hosted zone ID for CloudFront distributions (always `Z2FDTNDATAQYW2`)

---

## CloudFront-Specific Design Decisions

### Why Custom Origin Config Instead of S3 Origin Config?
The Pulumi implementation uses `CustomOriginConfig` for all origins rather than detecting S3 buckets and applying `S3OriginConfig`. This is intentional:
- **Simplicity**: The 80/20 API doesn't expose Origin Access Control (OAC) or Origin Access Identity (OAI) configuration. Users must configure S3 bucket policies separately if needed.
- **Consistency**: All origins are treated uniformly, whether they're S3, ALB, API Gateway, or custom domains.
- **Escape Hatch**: For users requiring OAC/OAI automation, they can use the Pulumi or Terraform modules directly or extend the protobuf spec in a future API version.

### Why Only Default Cache Behavior?
The implementation provisions a single default cache behavior that applies to all requests. Path-based routing via `ordered_cache_behavior` is not exposed in the API. This aligns with the 80/20 principle:
- **Default Use Case**: Most CloudFront deployments serve a single origin (S3 static site or ALB) without path-based routing.
- **Clarity**: A single cache behavior simplifies configuration and eliminates ambiguity about which behavior applies to a given path.
- **Future Extension**: Path-based routing can be added to the protobuf spec in a future API version if demand warrants it.

### Why Hardcode the CloudFront Hosted Zone ID?
The output `hosted_zone_id` is hardcoded to `Z2FDTNDATAQYW2`, which is the Route53 alias hosted zone ID for all CloudFront distributions. This is an AWS constant and never changes. Hardcoding it:
- Eliminates an unnecessary API call to retrieve it
- Ensures the output is available immediately without waiting for CloudFront distribution creation
- Matches AWS documentation and best practices

### Why Default TTL of 3600 Seconds?
The default cache behavior uses `default_ttl = 3600` (1 hour). This is a balanced default:
- Long enough to achieve meaningful cache hit ratios for static content
- Short enough to avoid stale content for semi-dynamic sites
- Can be overridden by `Cache-Control` headers from the origin if desired

---

## Integration with OpenMCF Stack Input Pattern

The Pulumi module expects an `AwsCloudFrontStackInput` protobuf message as input, which includes:
- `target`: The full `AwsCloudFront` resource (metadata + spec + status)
- `provider_config`: AWS credentials (access key, secret key, region, optional session token)

The OpenMCF CLI serializes this stack input and passes it to the Pulumi program via standard input. The Pulumi entrypoint deserializes it and invokes the module. This pattern:
- Decouples the CLI from the IaC backend (Pulumi or Terraform)
- Enables consistent validation and secret injection at the CLI level
- Allows the same manifest to be deployed with either `openmcf pulumi up` or `openmcf terraform apply`

---

## Common CloudFront Gotchas (Handled by This Implementation)

### 1. The us-east-1 Certificate Requirement
**The Problem**: CloudFront requires all ACM certificates to be in the `us-east-1` region, regardless of your origin's region. Developers often provision certificates in their application's region (e.g., `eu-west-1`) and encounter cryptic errors at deploy-time.

**How This Module Helps**: The protobuf validation enforces the `us-east-1` region in the certificate ARN pattern at validation time, not deploy-time. Invalid ARNs are rejected before Pulumi is even invoked.

### 2. The Custom Origin HTTPS Requirement
**The Problem**: When using a custom origin (not S3), CloudFront requires TLS configuration. Developers often forget to specify the origin protocol policy, resulting in default HTTP behavior.

**How This Module Helps**: The implementation hardcodes `OriginProtocolPolicy: pulumi.String("https-only")` for all custom origins, enforcing secure communication from CloudFront to the origin.

### 3. The Viewer Protocol Redirect
**The Problem**: Allowing `allow-all` viewer protocol policy permits insecure HTTP connections. This is a security and SEO issue (Google penalizes sites without HTTPS).

**How This Module Helps**: The default cache behavior is hardcoded to `ViewerProtocolPolicy: pulumi.String("redirect-to-https")`, ensuring all HTTP requests are automatically upgraded to HTTPS.

### 4. The TLS Version Security Hole
**The Problem**: Old TLS versions (TLSv1.0, TLSv1.1) are vulnerable to attacks. CloudFront's default minimum protocol version varies by configuration.

**How This Module Helps**: The viewer certificate is configured with `MinimumProtocolVersion: pulumi.String("TLSv1.2_2021")`, enforcing modern TLS standards.

---

## Debugging the Pulumi Module

For local development and debugging, the Pulumi module includes a `debug.sh` script that runs the compiled Go binary under Delve (the Go debugger). To enable it:

1. Uncomment the `binary` option in `Pulumi.yaml`:
   ```yaml
   options:
     binary: ./debug.sh
   ```

2. Run the Pulumi command as usual:
   ```bash
   openmcf pulumi preview \
     --manifest ../hack/manifest.yaml \
     --stack org/project/stack \
     --module-dir .
   ```

3. Pulumi will invoke `debug.sh`, which compiles the program and starts it under Delve on port 2345. You can attach your IDE's debugger to this port.

For more details, refer to the OpenMCF documentation: `docs/pages/docs/guide/debug-pulumi-modules.mdx`

---

## Testing Strategy

The protobuf spec includes comprehensive unit tests in `spec_test.go` that validate:
- Field-level validations (aliases uniqueness, certificate ARN format, origin domain pattern)
- CEL expressions (aliases require certificate, exactly one default origin)
- Enum enforcement (price class defined values)

These tests run via Ginkgo and use `protovalidate` for validation. All tests must pass before merging changes.

---

## Next Steps

- Refer to the [README.md](./README.md) for detailed CLI usage, deployment commands, and debugging instructions.
- Review the [examples.md](./examples.md) to explore common CloudFront deployment patterns, from minimal static sites to multi-origin distributions with custom domains.
- Check out the wider OpenMCF documentation for deeper insights into multi-cloud deployments, advanced features, and CLI usage patterns.
- For the research behind the 80/20 API design decisions, see `docs/README.md` at the component root.


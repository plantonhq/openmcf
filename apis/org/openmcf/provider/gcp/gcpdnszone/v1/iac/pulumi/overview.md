### Overview

The provided Pulumi module is designed to manage Google Cloud Platform (GCP) DNS Zones through a Kubernetes-like API resource model. The module leverages the `GcpDnsZone` API resource, which includes essential fields such as `apiVersion`, `kind`, `metadata`, `spec`, and `status`. The `GcpDnsZoneSpec` defines the configuration for the DNS zone, including the GCP project ID (using `StringValueOrRef` for flexible referencing), IAM service accounts for managing DNS records, and a list of DNS records with specific types, names, values, and TTL settings. By utilizing the Pulumi Google provider, the module programmatically creates a Managed Zone in GCP, ensuring that the DNS name adheres to proper formatting conventions. It also handles the replacement of dots with hyphens in the zone name to comply with GCP's naming requirements.

### StringValueOrRef Pattern

The `project_id` field uses the `StringValueOrRef` pattern, which allows users to either:
- Provide a direct literal value: `{value: "my-project-123"}`
- Reference another resource's output: `{valueFrom: {kind: GcpProject, name: "my-project", fieldPath: "status.outputs.project_id"}}`

**Current limitation**: Reference resolution is not yet implemented. Only literal `value` is used. References will fail silently (empty string).

**Future work**: Implement reference resolution in a shared library that all modules can use.

### IAM and DNS Record Management

In addition to setting up the Managed Zone, the module configures IAM bindings to grant specified service accounts the necessary permissions to manage DNS records within the zone. Although the current implementation grants broader `dns.admin` roles at the project level due to limitations in the GCP provider, it sets the foundation for more granular permission controls in the future. The module iterates through the provided DNS records, creating each one within the Managed Zone with the appropriate type, name, values, and TTL. Finally, it exports critical attributes such as the Managed Zone name, nameservers, and project ID to the stack outputs, facilitating seamless integration and status tracking within the broader infrastructure deployment workflow.

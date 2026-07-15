// Package pulumiannotationkeys defines the manifest annotation keys that carry
// Pulumi stack-location configuration.
//
// Platform-behavior signals live in metadata.annotations, never metadata.labels:
// labels are derived into cloud-provider tags by planton IaC modules, so a platform
// key there would leak internal detail onto the user's real cloud resources.
package pulumiannotationkeys

const (
	// StackFqdnAnnotationKey is the primary annotation that takes precedence over individual components
	// Format: "organization/project/stack"
	StackFqdnAnnotationKey = "pulumi.planton.dev/stack.fqdn"

	// OrganizationAnnotationKey is used when stack.fqdn is not present
	OrganizationAnnotationKey = "pulumi.planton.dev/organization"

	// ProjectAnnotationKey is used when stack.fqdn is not present
	ProjectAnnotationKey = "pulumi.planton.dev/project"

	// StackNameAnnotationKey is used when stack.fqdn is not present
	StackNameAnnotationKey = "pulumi.planton.dev/stack.name"
)

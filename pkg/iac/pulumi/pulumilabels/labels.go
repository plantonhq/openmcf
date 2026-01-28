package pulumilabels

const (
	// StackFqdnLabelKey is the primary label that takes precedence over individual components
	// Format: "organization/project/stack"
	StackFqdnLabelKey = "pulumi.openmcf.org/stack.fqdn"

	// OrganizationLabelKey is used when stack.fqdn is not present
	OrganizationLabelKey = "pulumi.openmcf.org/organization"

	// ProjectLabelKey is used when stack.fqdn is not present
	ProjectLabelKey = "pulumi.openmcf.org/project"

	// StackNameLabelKey is used when stack.fqdn is not present
	StackNameLabelKey = "pulumi.openmcf.org/stack.name"
)

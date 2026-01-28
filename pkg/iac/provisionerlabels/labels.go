package provisionerlabels

const (
	// ProvisionerLabelKey specifies which IaC provisioner to use
	// Supported values: "pulumi", "tofu", "terraform" (case-insensitive)
	ProvisionerLabelKey = "openmcf.org/provisioner"
)

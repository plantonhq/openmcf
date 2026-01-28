package kuberneteslabelkeys

import (
	"github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/labels/labelkeys"
)

var (
	Resource     = labelkeys.WithDomainPrefix("resource")
	Organization = labelkeys.WithDomainPrefix("organization")
	Environment  = labelkeys.WithDomainPrefix("environment")
	ResourceKind = labelkeys.WithDomainPrefix("kind")
	ResourceId   = labelkeys.WithDomainPrefix("id")
	ResourceName = labelkeys.WithDomainPrefix("name")
)

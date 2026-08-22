package module

import (
	digitaloceandatabasekafkatopicv1alpha1 "github.com/plantonhq/planton/catalog/digitalocean/digitaloceandatabasekafkatopic/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references for the rest of the module. The Kafka
// topic resource has no tag surface, so no Planton label set applies.
type Locals struct {
	DigitalOceanDatabaseKafkaTopic *digitaloceandatabasekafkatopicv1alpha1.DigitalOceanDatabaseKafkaTopic
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *digitaloceandatabasekafkatopicv1alpha1.DigitalOceanDatabaseKafkaTopicStackInput) *Locals {
	return &Locals{
		DigitalOceanDatabaseKafkaTopic: stackInput.Target,
	}
}

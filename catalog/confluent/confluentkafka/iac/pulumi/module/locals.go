package module

import (
	confluentkafkav1alpha1 "github.com/plantonhq/planton/catalog/confluent/confluentkafka/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	ConfluentKafka *confluentkafkav1alpha1.ConfluentKafka
}

func initializeLocals(ctx *pulumi.Context, stackInput *confluentkafkav1alpha1.ConfluentKafkaStackInput) *Locals {
	locals := &Locals{}

	//assign value for the locals variable to make it available across the project
	locals.ConfluentKafka = stackInput.Target

	return locals
}

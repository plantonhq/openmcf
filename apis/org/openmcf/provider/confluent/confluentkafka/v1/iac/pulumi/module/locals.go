package module

import (
	confluentkafkav1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/confluent/confluentkafka/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	ConfluentKafka *confluentkafkav1.ConfluentKafka
}

func initializeLocals(ctx *pulumi.Context, stackInput *confluentkafkav1.ConfluentKafkaStackInput) *Locals {
	locals := &Locals{}

	//assign value for the locals variable to make it available across the project
	locals.ConfluentKafka = stackInput.Target

	return locals
}

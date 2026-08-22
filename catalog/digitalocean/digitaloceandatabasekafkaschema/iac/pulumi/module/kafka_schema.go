package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// kafkaSchema registers the schema subject and exports its outputs. EVERY
// argument is create-only upstream (the provider has no update path): any
// change replaces the subject and drops all previously registered versions
// -- including a whitespace-only reformat of the definition, which is
// compared verbatim.
func kafkaSchema(
	ctx *pulumi.Context,
	locals *Locals,
	digitalOceanProvider *digitalocean.Provider,
) (*digitalocean.DatabaseKafkaSchemaRegistry, error) {
	spec := locals.DigitalOceanDatabaseKafkaSchema.Spec

	createdSchema, err := digitalocean.NewDatabaseKafkaSchemaRegistry(
		ctx,
		"schema",
		&digitalocean.DatabaseKafkaSchemaRegistryArgs{
			// References are resolved to the literal cluster UUID before
			// the module runs.
			ClusterId:   pulumi.String(spec.Cluster.GetValue()),
			SubjectName: pulumi.String(spec.SubjectName),
			SchemaType:  pulumi.String(spec.SchemaType),
			Schema:      pulumi.String(spec.Schema),
		},
		pulumi.Provider(digitalOceanProvider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to register digitalocean kafka schema subject")
	}

	ctx.Export(OpClusterId, createdSchema.ClusterId)
	ctx.Export(OpSubjectName, createdSchema.SubjectName)

	return createdSchema, nil
}

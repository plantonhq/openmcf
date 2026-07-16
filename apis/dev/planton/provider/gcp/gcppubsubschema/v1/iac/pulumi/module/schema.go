package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/pubsub"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// schema provisions the Pub/Sub schema. A shareable resource: one schema
// can validate messages on many topics (each attaches it by reference),
// so the event contract is evolved in one place. Definition changes
// commit a new revision in place (up to 20 revisions per schema) rather
// than replacing the resource; only renaming replaces it.
func schema(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpPubSubSchema.Spec

	// Enable the Pub/Sub API — the control plane that owns the schema.
	// disable_on_destroy stays false: tearing down one schema must never
	// disable the API for everything else in the project.
	pubsubApiArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("pubsub.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		pubsubApiArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdPubsubApi, err := projects.NewService(ctx,
		"gcppsch-pubsub.googleapis.com", pubsubApiArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable pubsub.googleapis.com api")
	}

	args := &pubsub.SchemaArgs{
		Name: pulumi.String(spec.SchemaName),
		// Type and definition travel together: the definition text is
		// parsed as the declared language (AVRO JSON or a protobuf
		// message), and later revisions must keep the same type.
		Type:       pulumi.StringPtr(spec.Type),
		Definition: pulumi.StringPtr(spec.Definition),
	}

	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project.
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
	}

	createdSchema, err := pubsub.NewSchema(ctx, "schema", args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdPubsubApi}))
	if err != nil {
		return errors.Wrap(err, "failed to create pub/sub schema")
	}

	// The resource ID is the fully qualified schema path
	// (projects/{project}/schemas/{name}) — the exact string a topic's
	// schema_settings.schema reference consumes.
	ctx.Export(OpSchemaId, createdSchema.ID())
	ctx.Export(OpSchemaName, createdSchema.Name)

	return nil
}

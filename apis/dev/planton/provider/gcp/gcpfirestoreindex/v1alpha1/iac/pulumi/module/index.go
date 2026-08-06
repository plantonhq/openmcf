package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/firestore"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// firestoreIndex provisions a Firestore composite index — the prerequisite
// for any query that filters or orders on multiple fields. Every property is
// immutable: changing anything replaces the index (Firestore rebuilds it in
// the background; the old index serves queries until the new one is ready).
// Firestore appends __name__ automatically.
func firestoreIndex(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpFirestoreIndex.Spec

	// Enable the Firestore API — composite indexes are managed through the
	// Firestore Admin API. disable_on_destroy stays false: tearing down one
	// index must never disable the API for everything else in the project.
	firestoreApiArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("firestore.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		firestoreApiArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdFirestoreApi, err := projects.NewService(ctx,
		"fsidx-firestore.googleapis.com", firestoreApiArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable firestore.googleapis.com api")
	}

	fields := firestore.IndexFieldArray{}
	for _, field := range spec.Fields {
		fieldArgs := &firestore.IndexFieldArgs{
			FieldPath: pulumi.String(field.FieldPath),
		}
		if field.Order != "" {
			fieldArgs.Order = pulumi.StringPtr(field.Order)
		}
		if field.ArrayConfig != "" {
			fieldArgs.ArrayConfig = pulumi.StringPtr(field.ArrayConfig)
		}
		if field.VectorConfig != nil {
			fieldArgs.VectorConfig = &firestore.IndexFieldVectorConfigArgs{
				Dimension: pulumi.IntPtr(int(field.VectorConfig.Dimension)),
				Flat:      &firestore.IndexFieldVectorConfigFlatArgs{},
			}
		}
		fields = append(fields, fieldArgs)
	}

	args := &firestore.IndexArgs{
		Collection: pulumi.String(spec.Collection),
		Fields:     fields,
	}

	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project (omit the arg entirely).
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
	}
	if spec.Database.GetValue() != "" {
		args.Database = pulumi.StringPtr(spec.Database.GetValue())
	}
	if spec.GetQueryScope() != "" {
		args.QueryScope = pulumi.StringPtr(spec.GetQueryScope())
	}
	if spec.GetApiScope() != "" {
		args.ApiScope = pulumi.StringPtr(spec.GetApiScope())
	}
	if spec.Density != "" {
		args.Density = pulumi.StringPtr(spec.Density)
	}

	createdIndex, err := firestore.NewIndex(ctx, "firestore-index", args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdFirestoreApi}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create firestore index")
	}

	ctx.Export(OpIndexId, createdIndex.Name)
	ctx.Export(OpCollection, pulumi.String(spec.Collection))

	return nil
}

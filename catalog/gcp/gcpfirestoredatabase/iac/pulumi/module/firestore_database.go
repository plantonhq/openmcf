package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/firestore"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// firestoreDatabase provisions the Firestore database: the top-level
// container for collections, documents, and indexes. A project supports
// multiple named databases beside the special "(default)" one.
//
// locationId, name, databaseEdition, the CMEK key, and
// appEngineIntegrationMode are immutable (ForceNew in the provider).
// type is mutable but switching between Native and Datastore Mode is a
// significant operational change. deleteProtectionState is GCP's own
// guard — deletion by any client fails while ENABLED.
func firestoreDatabase(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpFirestoreDatabase.Spec

	// Enable the Firestore API — the control plane the database is
	// managed through. disable_on_destroy stays false: tearing down one
	// database must never disable the API for everything else in the
	// project.
	firestoreApiArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("firestore.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		firestoreApiArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdFirestoreApi, err := projects.NewService(ctx,
		"fst-firestore.googleapis.com", firestoreApiArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable firestore.googleapis.com api")
	}

	args := &firestore.DatabaseArgs{
		LocationId: pulumi.String(spec.LocationId),
		Type:       pulumi.String(spec.Type),
		Name:       pulumi.StringPtr(spec.DatabaseName),
	}

	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project (omit the arg entirely; an empty string
	// would be sent verbatim and rejected).
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
	}

	if spec.ConcurrencyMode != "" {
		args.ConcurrencyMode = pulumi.StringPtr(spec.ConcurrencyMode)
	}
	if spec.PointInTimeRecoveryEnablement != "" {
		args.PointInTimeRecoveryEnablement = pulumi.StringPtr(spec.PointInTimeRecoveryEnablement)
	}
	if spec.GetDeleteProtectionState() != "" {
		args.DeleteProtectionState = pulumi.StringPtr(spec.GetDeleteProtectionState())
	}
	if spec.DatabaseEdition != "" {
		args.DatabaseEdition = pulumi.StringPtr(spec.DatabaseEdition)
	}
	if spec.AppEngineIntegrationMode != "" {
		args.AppEngineIntegrationMode = pulumi.StringPtr(spec.AppEngineIntegrationMode)
	}

	// ENTERPRISE-only data-access switches (spec CELs enforce the edition
	// pairing pre-deploy, matching the API).
	if spec.FirestoreDataAccessMode != "" {
		args.FirestoreDataAccessMode = pulumi.StringPtr(spec.FirestoreDataAccessMode)
	}
	if spec.MongodbCompatibleDataAccessMode != "" {
		args.MongodbCompatibleDataAccessMode = pulumi.StringPtr(spec.MongodbCompatibleDataAccessMode)
	}
	if spec.RealtimeUpdatesMode != "" {
		args.RealtimeUpdatesMode = pulumi.StringPtr(spec.RealtimeUpdatesMode)
	}

	// Create-time resource-manager tags (org policy / IAM conditions).
	if len(spec.ResourceManagerTags) > 0 {
		args.Tags = pulumi.ToStringMap(spec.ResourceManagerTags)
	}

	// CMEK encryption.
	if spec.KmsKeyName.GetValue() != "" {
		args.CmekConfig = &firestore.DatabaseCmekConfigArgs{
			KmsKeyName: pulumi.String(spec.KmsKeyName.GetValue()),
		}
	}

	// Defaults to DELETE so the IaC tool manages the full lifecycle (the
	// provider's own default ABANDON would leave the database behind on
	// destroy). PREVENT and ABANDON are deliberate choices — identical to
	// the Terraform module.
	deletionPolicy := spec.GetDeletionPolicy()
	if deletionPolicy == "" {
		deletionPolicy = "DELETE"
	}
	args.DeletionPolicy = pulumi.StringPtr(deletionPolicy)

	createdDatabase, err := firestore.NewDatabase(ctx, "firestore-database", args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdFirestoreApi}))
	if err != nil {
		return errors.Wrap(err, "failed to create firestore database")
	}

	// The resource ID already carries the projects/{project}/databases/{name}
	// form, with the project resolved even when the spec rode the
	// provider default — identical to the Terraform module's id output.
	ctx.Export(OpDatabaseId, createdDatabase.ID())
	ctx.Export(OpDatabaseName, createdDatabase.Name)
	ctx.Export(OpUid, createdDatabase.Uid)
	ctx.Export(OpCreateTime, createdDatabase.CreateTime)
	ctx.Export(OpEarliestVersionTime, createdDatabase.EarliestVersionTime)
	ctx.Export(OpVersionRetentionPeriod, createdDatabase.VersionRetentionPeriod)
	ctx.Export(OpKeyPrefix, createdDatabase.KeyPrefix)
	ctx.Export(OpUpdateTime, createdDatabase.UpdateTime)

	return nil
}

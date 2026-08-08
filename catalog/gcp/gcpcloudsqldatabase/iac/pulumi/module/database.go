package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/sql"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// database provisions a logical database inside a Cloud SQL instance.
// Databases carry their own lifecycle: create and drop application databases
// freely without touching the instance node they live on.
//
// No API enablement here: the instance this database lives on cannot exist
// without sqladmin.googleapis.com already enabled (its own module enables
// it), so a database module enabling the API again would only add churn.
//
// Charset/collation semantics are engine-specific — MySQL accepts any
// supported pair, PostgreSQL requires UTF8 at creation with an OS-locale
// collation, SQL Server ignores charset entirely. The API validates the
// combination at deploy time.
func database(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpCloudSqlDatabase.Spec

	args := &sql.DatabaseArgs{
		Name:     pulumi.String(spec.DatabaseName),
		Instance: pulumi.String(spec.Instance.GetValue()),
	}

	// An empty project falls back to the provider's default project — the
	// ambient-project contract every GCP kind honors.
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	// Empty means the engine default; sending "" would be rejected.
	if spec.Charset != "" {
		args.Charset = pulumi.StringPtr(spec.Charset)
	}
	if spec.Collation != "" {
		args.Collation = pulumi.StringPtr(spec.Collation)
	}

	// DELETE (default) drops the database; ABANDON removes it from IaC
	// management — the documented workaround when live connections block a
	// PostgreSQL drop; PREVENT fails destroying previews.
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
	}

	createdDatabase, err := sql.NewDatabase(ctx, "database", args, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to create database")
	}

	ctx.Export(OpDatabaseName, createdDatabase.Name)
	ctx.Export(OpSelfLink, createdDatabase.SelfLink)

	return nil
}

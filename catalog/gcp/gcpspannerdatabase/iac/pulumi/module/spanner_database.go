package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/spanner"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// spannerDatabase provisions the Spanner database — schemas, data,
// encryption, and retention on a Spanner instance. name, dialect,
// encryption, and instance are immutable; retention, drop protection, and
// appended DDL update in place.
func spannerDatabase(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpSpannerDatabase.Spec

	// Enable the Spanner API so a fresh project can host databases.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("spanner.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"spanner-spanner.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable spanner.googleapis.com api")
	}

	// IaC-side deletion guard (spec default TRUE): a destroy fails while
	// set, before touching GCP. Always set explicitly so both engines share
	// identical destroy semantics regardless of provider defaults.
	deletionProtection := true
	if spec.DeletionProtection != nil {
		deletionProtection = spec.GetDeletionProtection()
	}

	args := &spanner.DatabaseArgs{
		Instance:           pulumi.String(spec.Instance.GetValue()),
		Name:               pulumi.StringPtr(locals.DatabaseName),
		DeletionProtection: pulumi.BoolPtr(deletionProtection),
	}

	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project (omit the arg entirely).
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
	}

	// Dialect (immutable, permanent choice).
	if spec.DatabaseDialect != "" {
		args.DatabaseDialect = pulumi.StringPtr(spec.DatabaseDialect)
	}

	// Version retention period for point-in-time recovery.
	if spec.VersionRetentionPeriod != "" {
		args.VersionRetentionPeriod = pulumi.StringPtr(spec.VersionRetentionPeriod)
	}

	// DDL is append-only after creation: new statements apply via
	// UpdateDDL; editing or removing an existing entry forces database
	// recreation. The provider quotes identifiers per dialect.
	if len(spec.Ddl) > 0 {
		args.Ddls = pulumi.ToStringArray(spec.Ddl)
	}

	// GCP API-side lock: while true, no interface (console, gcloud, IaC)
	// can delete the database, and the PARENT INSTANCE cannot be deleted
	// either.
	if spec.EnableDropProtection {
		args.EnableDropProtection = pulumi.BoolPtr(true)
	}

	// CMEK: immutable. kms_key_name for regional instance configs,
	// kms_key_names (one key per region) for multi-region configs — the
	// spec enforces exactly one shape pre-deploy.
	if spec.EncryptionConfig != nil {
		encryptionArgs := &spanner.DatabaseEncryptionConfigArgs{}
		if spec.EncryptionConfig.KmsKeyName.GetValue() != "" {
			encryptionArgs.KmsKeyName = pulumi.StringPtr(spec.EncryptionConfig.KmsKeyName.GetValue())
		}
		if len(spec.EncryptionConfig.KmsKeyNames) > 0 {
			keyNames := make([]string, 0, len(spec.EncryptionConfig.KmsKeyNames))
			for _, keyName := range spec.EncryptionConfig.KmsKeyNames {
				keyNames = append(keyNames, keyName.GetValue())
			}
			encryptionArgs.KmsKeyNames = pulumi.ToStringArray(keyNames)
		}
		args.EncryptionConfig = encryptionArgs
	}

	// Default time zone (affects time-zone-dependent SQL functions).
	if spec.DefaultTimeZone != "" {
		args.DefaultTimeZone = pulumi.StringPtr(spec.DefaultTimeZone)
	}

	createdDatabase, err := spanner.NewDatabase(ctx, "spanner-database", args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdProjectService}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create spanner database")
	}

	// database_id is built from the created resource's resolved attributes
	// so the output is correct under the ambient-project fallback (the spec
	// project may be empty).
	ctx.Export(OpDatabaseId, pulumi.Sprintf(
		"projects/%s/instances/%s/databases/%s",
		createdDatabase.Project,
		createdDatabase.Instance,
		createdDatabase.Name,
	))
	ctx.Export(OpDatabaseName, createdDatabase.Name)
	ctx.Export(OpState, createdDatabase.State)

	return nil
}

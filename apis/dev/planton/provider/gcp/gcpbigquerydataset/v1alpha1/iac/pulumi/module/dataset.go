package module

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/bigquery"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// dataset provisions the BigQuery dataset — the container that pins data
// location and owns the defaults every contained table inherits
// (expiration, CMEK, collation) plus the dataset-level ACL. dataset_id,
// location, and is_case_insensitive are immutable; everything else updates
// in place.
func dataset(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpBigQueryDataset.Spec

	// Enable the BigQuery API so a fresh project can host datasets.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("bigquery.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"bigquery-bigquery.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable bigquery.googleapis.com api")
	}

	args := &bigquery.DatasetArgs{
		DatasetId: pulumi.String(spec.DatasetId),
		Location:  pulumi.StringPtr(spec.Location),
		Labels:    pulumi.ToStringMap(locals.GcpLabels),
	}

	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project (omit the arg entirely).
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
	}

	// Empty-/zero-means-unset scalars are omitted so the API applies its
	// own server-side defaults (max_time_travel_hours 168, LOGICAL billing).
	if spec.FriendlyName != "" {
		args.FriendlyName = pulumi.StringPtr(spec.FriendlyName)
	}
	if spec.Description != "" {
		args.Description = pulumi.StringPtr(spec.Description)
	}
	if spec.DefaultCollation != "" {
		args.DefaultCollation = pulumi.StringPtr(spec.DefaultCollation)
	}
	if spec.StorageBillingModel != "" {
		args.StorageBillingModel = pulumi.StringPtr(spec.StorageBillingModel)
	}
	if spec.DefaultTableExpirationMs > 0 {
		args.DefaultTableExpirationMs = pulumi.IntPtr(int(spec.DefaultTableExpirationMs))
	}
	if spec.DefaultPartitionExpirationMs > 0 {
		args.DefaultPartitionExpirationMs = pulumi.IntPtr(int(spec.DefaultPartitionExpirationMs))
	}
	// The provider types max_time_travel_hours as a string.
	if spec.MaxTimeTravelHours > 0 {
		args.MaxTimeTravelHours = pulumi.StringPtr(strconv.Itoa(int(spec.MaxTimeTravelHours)))
	}
	if spec.IsCaseInsensitive {
		args.IsCaseInsensitive = pulumi.BoolPtr(true)
	}

	// When false (the safe default), destroy fails while the dataset
	// contains tables — the guard against deleting data with its container.
	if spec.DeleteContentsOnDestroy {
		args.DeleteContentsOnDestroy = pulumi.BoolPtr(true)
	}

	if len(spec.ResourceTags) > 0 {
		args.ResourceTags = pulumi.ToStringMap(spec.ResourceTags)
	}

	// CMEK default for all new tables. The BigQuery service agent must hold
	// cryptoKeyEncrypterDecrypter on the key before the first table write.
	if spec.KmsKeyName.GetValue() != "" {
		args.DefaultEncryptionConfiguration = &bigquery.DatasetDefaultEncryptionConfigurationArgs{
			KmsKeyName: pulumi.String(spec.KmsKeyName.GetValue()),
		}
	}

	// The spec's access list is AUTHORITATIVE: these entries become the
	// dataset's complete ACL. An entry is either a principal grant (role +
	// one identity) or a resource authorization (view/routine/dataset with
	// implicit read access and no role) — the spec's CEL rules guarantee
	// the shape before this module ever runs.
	if len(spec.Access) > 0 {
		accessArray := bigquery.DatasetAccessTypeArray{}
		for _, entry := range spec.Access {
			accessArgs := &bigquery.DatasetAccessTypeArgs{}

			if entry.Role != "" {
				accessArgs.Role = pulumi.StringPtr(entry.Role)
			}
			if entry.UserByEmail != "" {
				accessArgs.UserByEmail = pulumi.StringPtr(entry.UserByEmail)
			}
			if entry.GroupByEmail != "" {
				accessArgs.GroupByEmail = pulumi.StringPtr(entry.GroupByEmail)
			}
			if entry.Domain != "" {
				accessArgs.Domain = pulumi.StringPtr(entry.Domain)
			}
			if entry.SpecialGroup != "" {
				accessArgs.SpecialGroup = pulumi.StringPtr(entry.SpecialGroup)
			}
			if entry.IamMember != "" {
				accessArgs.IamMember = pulumi.StringPtr(entry.IamMember)
			}
			if entry.View != nil {
				accessArgs.View = &bigquery.DatasetAccessViewArgs{
					ProjectId: pulumi.String(entry.View.ProjectId),
					DatasetId: pulumi.String(entry.View.DatasetId),
					TableId:   pulumi.String(entry.View.TableId),
				}
			}
			if entry.Routine != nil {
				accessArgs.Routine = &bigquery.DatasetAccessRoutineArgs{
					ProjectId: pulumi.String(entry.Routine.ProjectId),
					DatasetId: pulumi.String(entry.Routine.DatasetId),
					RoutineId: pulumi.String(entry.Routine.RoutineId),
				}
			}
			// The provider nests the grantee dataset reference one level
			// deeper (dataset { dataset { ... } target_types }); the spec
			// flattens that single-purpose wrapper.
			if entry.Dataset != nil {
				accessArgs.Dataset = &bigquery.DatasetAccessDatasetArgs{
					Dataset: &bigquery.DatasetAccessDatasetDatasetArgs{
						ProjectId: pulumi.String(entry.Dataset.ProjectId),
						DatasetId: pulumi.String(entry.Dataset.DatasetId),
					},
					TargetTypes: pulumi.ToStringArray(entry.Dataset.TargetTypes),
				}
			}
			if entry.Condition != nil {
				conditionArgs := &bigquery.DatasetAccessConditionArgs{
					Expression: pulumi.String(entry.Condition.Expression),
				}
				if entry.Condition.Title != "" {
					conditionArgs.Title = pulumi.StringPtr(entry.Condition.Title)
				}
				if entry.Condition.Description != "" {
					conditionArgs.Description = pulumi.StringPtr(entry.Condition.Description)
				}
				if entry.Condition.Location != "" {
					conditionArgs.Location = pulumi.StringPtr(entry.Condition.Location)
				}
				accessArgs.Condition = conditionArgs
			}

			accessArray = append(accessArray, accessArgs)
		}
		args.Accesses = accessArray
	}

	// Immutable: converts the dataset into a read-only projection of an
	// external source (e.g. AWS Glue) through a BigQuery Omni connection.
	if spec.ExternalDatasetReference != nil {
		args.ExternalDatasetReference = &bigquery.DatasetExternalDatasetReferenceArgs{
			ExternalSource: pulumi.String(spec.ExternalDatasetReference.ExternalSource),
			Connection:     pulumi.String(spec.ExternalDatasetReference.Connection),
		}
	}

	// Hive Metastore compatibility metadata for open-source engines.
	if spec.ExternalCatalogOptions != nil {
		catalogArgs := &bigquery.DatasetExternalCatalogDatasetOptionsArgs{}
		if spec.ExternalCatalogOptions.DefaultStorageLocationUri != "" {
			catalogArgs.DefaultStorageLocationUri = pulumi.StringPtr(spec.ExternalCatalogOptions.DefaultStorageLocationUri)
		}
		if len(spec.ExternalCatalogOptions.Parameters) > 0 {
			catalogArgs.Parameters = pulumi.ToStringMap(spec.ExternalCatalogOptions.Parameters)
		}
		args.ExternalCatalogDatasetOptions = catalogArgs
	}

	createdDataset, err := bigquery.NewDataset(ctx, "bigquery-dataset", args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdProjectService}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create bigquery dataset")
	}

	ctx.Export(OpDatasetId, createdDataset.DatasetId)
	ctx.Export(OpSelfLink, createdDataset.SelfLink)
	// Read from the created resource so the output is correct under the
	// ambient-project fallback (the spec project may be empty).
	ctx.Export(OpProject, createdDataset.Project)
	ctx.Export(OpCreationTime, createdDataset.CreationTime)
	ctx.Export(OpLocation, createdDataset.Location)
	ctx.Export(OpEtag, createdDataset.Etag)

	return nil
}

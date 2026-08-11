package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/bigtable"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// bigtableInstance provisions the Bigtable instance: the logical
// container, served by one or more clusters (physical replicas, each in
// its own zone). Multi-cluster instances replicate automatically; the
// client library routes and fails over transparently.
//
// Per-cluster immutability is enforced server-side: zone, storageType,
// kmsKeyName, and nodeScalingFactor cannot change on an existing
// clusterId. numNodes and autoscaling bounds resize in place.
func bigtableInstance(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpBigtableInstance.Spec

	// Enable the Bigtable Admin API — the control plane instance and
	// table management run through. disable_on_destroy stays false:
	// tearing down one instance must never disable the API for
	// everything else in the project.
	adminApiArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("bigtableadmin.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		adminApiArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdAdminApi, err := projects.NewService(ctx,
		"bt-bigtableadmin.googleapis.com", adminApiArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable bigtableadmin.googleapis.com api")
	}

	// Build cluster configurations from the spec.
	clusters := bigtable.InstanceClusterArray{}
	for _, c := range spec.Clusters {
		clusterArgs := bigtable.InstanceClusterArgs{
			ClusterId: pulumi.String(c.ClusterId),
			Zone:      pulumi.String(c.Zone),
		}

		// Storage type (optional, middleware applies default "SSD").
		if c.GetStorageType() != "" {
			clusterArgs.StorageType = pulumi.StringPtr(c.GetStorageType())
		}

		// Node scaling factor (optional, GCP defaults to 1X).
		if c.NodeScalingFactor != "" {
			clusterArgs.NodeScalingFactor = pulumi.StringPtr(c.NodeScalingFactor)
		}

		// CMEK encryption (optional).
		if c.KmsKeyName != nil && c.KmsKeyName.GetValue() != "" {
			clusterArgs.KmsKeyName = pulumi.StringPtr(c.KmsKeyName.GetValue())
		}

		// Scaling: either fixed num_nodes or autoscaling (mutually exclusive,
		// validated by proto CEL). If neither is set, Bigtable auto-allocates.
		if c.NumNodes > 0 {
			clusterArgs.NumNodes = pulumi.IntPtr(int(c.NumNodes))
		}
		if c.AutoscalingConfig != nil {
			autoscalingArgs := &bigtable.InstanceClusterAutoscalingConfigArgs{
				MinNodes:  pulumi.Int(int(c.AutoscalingConfig.MinNodes)),
				MaxNodes:  pulumi.Int(int(c.AutoscalingConfig.MaxNodes)),
				CpuTarget: pulumi.Int(int(c.AutoscalingConfig.CpuTarget)),
			}
			if c.AutoscalingConfig.StorageTarget > 0 {
				autoscalingArgs.StorageTarget = pulumi.IntPtr(int(c.AutoscalingConfig.StorageTarget))
			}
			clusterArgs.AutoscalingConfig = autoscalingArgs
		}

		clusters = append(clusters, clusterArgs)
	}

	// Terraform-side deletion guard (spec default TRUE): always sent
	// explicitly so destroy behavior is identical on both engines — a
	// manifest that never mentions deletion protection must behave the
	// same everywhere.
	deletionProtection := true
	if spec.DeletionProtection != nil {
		deletionProtection = spec.GetDeletionProtection()
	}

	args := &bigtable.InstanceArgs{
		Name:               pulumi.String(spec.InstanceName),
		Labels:             pulumi.ToStringMap(locals.GcpLabels),
		Clusters:           clusters,
		DeletionProtection: pulumi.BoolPtr(deletionProtection),
		ForceDestroy:       pulumi.BoolPtr(spec.ForceDestroy),
	}

	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project (omit the arg entirely; an empty string
	// would be sent verbatim and rejected).
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
	}

	// Display name (optional, GCP defaults to instance name).
	if spec.DisplayName != "" {
		args.DisplayName = pulumi.StringPtr(spec.DisplayName)
	}

	// Edition gates feature availability (ENTERPRISE_PLUS unlocks
	// multi-location automated-backup placement on tables). Unset lets
	// the provider apply its ENTERPRISE default; upgrades apply in
	// place, there is no downgrade path.
	if spec.Edition != "" {
		args.Edition = pulumi.StringPtr(spec.Edition)
	}

	// Resource Manager tags for org-policy and IAM conditions.
	// Create-time only (ForceNew): a tag change replaces the instance.
	if len(spec.ResourceManagerTags) > 0 {
		args.Tags = pulumi.ToStringMap(spec.ResourceManagerTags)
	}

	// What a PERMITTED destroy does once deletion_protection allows one:
	// DELETE (default), PREVENT (destroy fails), or ABANDON (drop from
	// state, keep the instance running in GCP).
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
	}

	createdInstance, err := bigtable.NewInstance(ctx, "bigtable-instance", args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdAdminApi}))
	if err != nil {
		return errors.Wrap(err, "failed to create bigtable instance")
	}

	// Export outputs.
	ctx.Export(OpInstanceId, createdInstance.ID())
	ctx.Export(OpInstanceName, pulumi.String(spec.InstanceName))

	return nil
}

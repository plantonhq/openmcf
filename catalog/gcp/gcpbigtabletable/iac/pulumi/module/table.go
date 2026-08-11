package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/bigtable"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// bigtableTable provisions the table (column families, pre-splits, change
// streams, automated backups, deletion protection) plus one GC-policy
// resource per column family that declares one — the API's own
// granularity, folded into this kind because a GC policy has no
// independent life apart from its family.
//
// Table name, instance, and splitKeys are immutable (ForceNew in the
// provider); changing splitKeys REPLACES the table and its data. Column
// families are mutable. deletionProtection is the API-side guard (spec
// default PROTECTED): deletion by ANY client fails until it is set
// UNPROTECTED — stronger than an IaC-side guard for a data-bearing
// resource.
func bigtableTable(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpBigtableTable.Spec

	// Enable the Bigtable Admin API — table and GC-policy management run
	// through it. disable_on_destroy stays false: tearing down one table
	// must never disable the API for everything else in the project.
	adminApiArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("bigtableadmin.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		adminApiArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdAdminApi, err := projects.NewService(ctx,
		"bttbl-bigtableadmin.googleapis.com", adminApiArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable bigtableadmin.googleapis.com api")
	}

	// Column families are created with no GC policy on the table object;
	// per-family retention lives in the GCPolicy resources below, so a
	// policy change never touches the table itself.
	columnFamilies := bigtable.TableColumnFamilyArray{}
	for _, cf := range spec.ColumnFamilies {
		cfArgs := &bigtable.TableColumnFamilyArgs{
			Family: pulumi.String(cf.Family),
		}
		if cf.Type != "" {
			cfArgs.Type = pulumi.StringPtr(cf.Type)
		}
		columnFamilies = append(columnFamilies, cfArgs)
	}

	// The spec default PROTECTED arrives materialized; send it explicitly
	// so destroy behavior is identical on both engines (the API default
	// is UNPROTECTED — weaker than a data-bearing resource deserves).
	deletionProtection := "PROTECTED"
	if spec.DeletionProtection != nil && spec.GetDeletionProtection() != "" {
		deletionProtection = spec.GetDeletionProtection()
	}

	args := &bigtable.TableArgs{
		Name:               pulumi.StringPtr(locals.TableName),
		InstanceName:       pulumi.String(spec.Instance.GetValue()),
		DeletionProtection: pulumi.StringPtr(deletionProtection),
		ColumnFamilies:     columnFamilies,
	}

	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project (omit the arg entirely).
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
	}
	if len(spec.SplitKeys) > 0 {
		splitKeys := make(pulumi.StringArray, 0, len(spec.SplitKeys))
		for _, k := range spec.SplitKeys {
			splitKeys = append(splitKeys, pulumi.String(k))
		}
		args.SplitKeys = splitKeys
	}
	if spec.ChangeStreamRetention != "" {
		args.ChangeStreamRetention = pulumi.StringPtr(spec.ChangeStreamRetention)
	}
	if spec.RowKeySchema != "" {
		args.RowKeySchema = pulumi.StringPtr(spec.RowKeySchema)
	}
	if spec.AutomatedBackupPolicy != nil {
		backupPolicyArgs := &bigtable.TableAutomatedBackupPolicyArgs{
			RetentionPeriod: pulumi.StringPtr(spec.AutomatedBackupPolicy.RetentionPeriod),
			Frequency:       pulumi.StringPtr(spec.AutomatedBackupPolicy.Frequency),
		}
		// Zones backups may be created in (empty = all zones of the
		// instance); ENTERPRISE_PLUS instances only. Optional+Computed in
		// the provider — sent only when set so unset never fights the
		// server-populated read-back.
		if len(spec.AutomatedBackupPolicy.Locations) > 0 {
			backupPolicyArgs.Locations = pulumi.ToStringArray(spec.AutomatedBackupPolicy.Locations)
		}
		args.AutomatedBackupPolicy = backupPolicyArgs
	}

	// One spec field drives the destroy behavior of BOTH objects this
	// kind manages (the table and its per-family GC policies): DELETE
	// (default), PREVENT (destroy fails), or ABANDON (drop from state,
	// keep the table).
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
	}

	createdTable, err := bigtable.NewTable(ctx, "table", args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdAdminApi}))
	if err != nil {
		return errors.Wrap(err, "failed to create bigtable table")
	}

	// One GC policy per column family that declares one. Bigtable never
	// deletes old cell versions without a GC policy, so an unbounded
	// family accumulates every write forever. Policies are mutable in
	// place; deleting one resets the family to "no GC" rather than
	// deleting data.
	for _, cf := range spec.ColumnFamilies {
		if cf.GcPolicy == nil {
			continue
		}

		gcArgs := &bigtable.GCPolicyArgs{
			InstanceName: pulumi.String(spec.Instance.GetValue()),
			Table:        createdTable.Name,
			ColumnFamily: pulumi.String(cf.Family),
		}
		if spec.ProjectId.GetValue() != "" {
			gcArgs.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
		}

		// The raw JSON tree and the typed fields are mutually exclusive
		// (enforced by spec CEL) — exactly one shape reaches the provider.
		if cf.GcPolicy.GcRules != "" {
			gcArgs.GcRules = pulumi.StringPtr(cf.GcPolicy.GcRules)
		}
		if cf.GcPolicy.Mode != "" {
			gcArgs.Mode = pulumi.StringPtr(cf.GcPolicy.Mode)
		}
		if cf.GcPolicy.MaxAge != "" {
			gcArgs.MaxAge = &bigtable.GCPolicyMaxAgeArgs{
				Duration: pulumi.StringPtr(cf.GcPolicy.MaxAge),
			}
		}
		if cf.GcPolicy.MaxVersions > 0 {
			gcArgs.MaxVersions = bigtable.GCPolicyMaxVersionArray{
				&bigtable.GCPolicyMaxVersionArgs{
					Number: pulumi.Int(int(cf.GcPolicy.MaxVersions)),
				},
			}
		}
		// Allows EXPANDING what is eligible for collection on a
		// replicated instance — Bigtable otherwise rejects the change as
		// a data-loss safety measure.
		if cf.GcPolicy.IgnoreWarnings {
			gcArgs.IgnoreWarnings = pulumi.BoolPtr(true)
		}
		// Mirrors the table's policy (one spec field, both objects).
		// ABANDON is also the escape hatch when Bigtable rejects a
		// GC-policy delete on a replicated instance.
		if spec.DeletionPolicy != "" {
			gcArgs.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
		}

		if _, err := bigtable.NewGCPolicy(ctx, "gc-policy-"+cf.Family, gcArgs,
			pulumi.Provider(gcpProvider),
			pulumi.DependsOn([]pulumi.Resource{createdTable})); err != nil {
			return errors.Wrapf(err, "failed to create gc policy for column family %s", cf.Family)
		}
	}

	ctx.Export(OpTableId, createdTable.ID())
	ctx.Export(OpTableName, createdTable.Name)
	ctx.Export(OpInstanceName, pulumi.String(spec.Instance.GetValue()))

	return nil
}

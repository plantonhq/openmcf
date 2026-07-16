package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/alloydb"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// user provisions a database user on an AlloyDB cluster. Users are
// first-class nodes: one per application with its own credential instead of
// sharing the cluster's initial superuser.
func user(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpAlloydbUser.Spec

	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("alloydb.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"alloydb-alloydb.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable alloydb.googleapis.com api")
	}

	args := &alloydb.UserArgs{
		Cluster:  pulumi.String(spec.Cluster.GetValue()),
		UserId:   pulumi.String(spec.UserId),
		UserType: pulumi.String(spec.GetUserType()),
	}

	if len(spec.DatabaseRoles) > 0 {
		args.DatabaseRoles = pulumi.ToStringArray(spec.DatabaseRoles)
	}

	if spec.Password != "" {
		args.Password = pulumi.ToSecret(pulumi.String(spec.Password)).(pulumi.StringOutput)
	}

	createdUser, err := alloydb.NewUser(ctx, "user", args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdProjectService}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create alloydb user")
	}

	ctx.Export(OpName, createdUser.Name)
	ctx.Export(OpUserId, createdUser.UserId)
	ctx.Export(OpClusterId, createdUser.Cluster)

	return nil
}

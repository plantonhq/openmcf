package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/monitoring"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// dashboard provisions the Cloud Monitoring dashboard from the spec's one
// JSON document. The provider validates the document is JSON at plan time
// and suppresses diffs on server-added keys (etag, name), so a dashboard
// exported from the GCP console round-trips cleanly.
func dashboard(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpMonitoringDashboard.Spec

	// Enable the Cloud Monitoring API so a fresh project can host the
	// dashboard. disable_on_destroy stays false (the provider default):
	// tearing down one dashboard must never disable monitoring for
	// everything else in the project. Matches the Terraform module.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("monitoring.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"dashboard-monitoring.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable monitoring.googleapis.com api")
	}

	args := &monitoring.DashboardArgs{
		DashboardJson: pulumi.String(spec.DashboardJson),
	}

	// Unset defers to the provider default (DELETE).
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.String(spec.DeletionPolicy)
	}

	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project (empty string would be sent verbatim and
	// rejected).
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	createdDashboard, err := monitoring.NewDashboard(ctx, "dashboard", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{createdProjectService}))
	if err != nil {
		return errors.Wrap(err, "failed to create dashboard")
	}

	// The provider computes the dashboard's resource name into `id`
	// (projects/{project}/dashboards/{dashboard_id}).
	ctx.Export(OpDashboardName, createdDashboard.ID())

	return nil
}

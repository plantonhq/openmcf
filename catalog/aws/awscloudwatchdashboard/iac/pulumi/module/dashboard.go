package module

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cloudwatch"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// dashboard creates the CloudWatch dashboard and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - the dashboard's name is spec.dashboard_name (an explicit field -
//     dashboard names carry uppercase metadata.name cannot), and
//     changing it replaces the dashboard;
//   - create and update are the same AWS call (PutDashboard is a pure
//     upsert), so every body change applies in place;
//   - dashboards are untaggable at AWS - no tags argument exists on
//     the resource.
func dashboard(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	// The body arrives as a Struct; the provider wants the document as
	// a JSON string. AWS normalizes it server-side and the provider
	// diffs it semantically, so key order never causes drift.
	bodyJson, err := json.Marshal(spec.DashboardBody.AsMap())
	if err != nil {
		return errors.Wrap(err, "marshal dashboard body")
	}

	createdDashboard, err := cloudwatch.NewDashboard(ctx, "dashboard", &cloudwatch.DashboardArgs{
		DashboardName: pulumi.String(spec.DashboardName),
		DashboardBody: pulumi.String(string(bodyJson)),
	}, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create dashboard")
	}

	ctx.Export(OpDashboardName, createdDashboard.DashboardName)
	ctx.Export(OpDashboardArn, createdDashboard.DashboardArn)
	return nil
}

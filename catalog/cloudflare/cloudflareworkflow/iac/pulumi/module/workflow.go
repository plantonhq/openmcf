package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// workflow registers the workflow: the binding of a class exported by a
// deployed Worker script to a named workflow in the account. Cloudflare's
// create IS a PUT (name-as-upsert): registering an existing name adopts and
// overwrites it rather than failing -- names must be chosen deliberately.
// account_id and workflow_name force replacement; class_name and
// script_name update in place (a full-body PUT to the same endpoint).
func workflow(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareWorkflow.Spec

	args := &cloudflare.WorkflowArgs{
		AccountId:    pulumi.String(spec.AccountId),
		WorkflowName: pulumi.String(spec.WorkflowName),
		ClassName:    pulumi.String(spec.ClassName),
		ScriptName:   pulumi.String(spec.ScriptName.GetValue()),
	}

	// Retention values are dynamic at the API (integer milliseconds or a
	// duration expression); the spec carries both forms as strings and they
	// pass through verbatim. Only set what the manifest states -- an absent
	// tree keeps Cloudflare's defaults.
	if spec.DefaultRetention != nil {
		retention := cloudflare.WorkflowDefaultRetentionArgs{}
		if spec.DefaultRetention.ErrorRetention != "" {
			retention.ErrorRetention = pulumi.String(spec.DefaultRetention.ErrorRetention)
		}
		if spec.DefaultRetention.SuccessRetention != "" {
			retention.SuccessRetention = pulumi.String(spec.DefaultRetention.SuccessRetention)
		}
		args.DefaultRetention = retention
	}

	if spec.Limits != nil && spec.Limits.Steps != nil {
		args.Limits = cloudflare.WorkflowLimitsArgs{
			Steps: pulumi.IntPtr(int(spec.Limits.GetSteps())),
		}
	}

	if len(spec.Schedules) > 0 {
		schedules := cloudflare.WorkflowScheduleArray{}
		for _, schedule := range spec.Schedules {
			schedules = append(schedules, cloudflare.WorkflowScheduleArgs{
				Cron: pulumi.String(schedule.Cron),
			})
		}
		args.Schedules = schedules
	}

	createdWorkflow, err := cloudflare.NewWorkflow(
		ctx,
		"workflow",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create workflow")
	}

	ctx.Export(OpWorkflowName, createdWorkflow.WorkflowName)
	ctx.Export(OpVersionId, createdWorkflow.VersionId)

	return nil
}

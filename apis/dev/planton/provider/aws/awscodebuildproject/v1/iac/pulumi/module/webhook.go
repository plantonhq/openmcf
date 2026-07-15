package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/codebuild"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// webhook creates the folded 1:1 webhook satellite for source-triggered
// builds. With manual_creation, AWS mints the payload URL and HMAC secret
// WITHOUT registering anything at the provider -- the operator wires the
// repository webhook by hand from this module's outputs. Only called when
// spec.webhook is configured.
func webhook(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
	proj *codebuild.Project,
) (*codebuild.Webhook, error) {
	spec := locals.AwsCodeBuildProject.Spec.Webhook

	args := &codebuild.WebhookArgs{
		ProjectName: proj.Name,
	}

	if spec.BuildType != "" {
		args.BuildType = pulumi.StringPtr(spec.BuildType)
	}
	if spec.ManualCreation {
		args.ManualCreation = pulumi.BoolPtr(true)
	}

	if len(spec.FilterGroups) > 0 {
		var filterGroups codebuild.WebhookFilterGroupArray
		for _, fg := range spec.FilterGroups {
			var filters codebuild.WebhookFilterGroupFilterArray
			for _, f := range fg.Filters {
				filterArgs := &codebuild.WebhookFilterGroupFilterArgs{
					Type:    pulumi.String(f.Type),
					Pattern: pulumi.String(f.Pattern),
				}
				if f.ExcludeMatchedPattern {
					filterArgs.ExcludeMatchedPattern = pulumi.BoolPtr(true)
				}
				filters = append(filters, filterArgs)
			}
			filterGroups = append(filterGroups, &codebuild.WebhookFilterGroupArgs{
				Filters: filters,
			})
		}
		args.FilterGroups = filterGroups
	}

	// Organization/group-scoped webhooks fire for every repository in scope
	// (runner projects, org-wide CI) instead of a single repository.
	if spec.ScopeConfiguration != nil {
		scArgs := &codebuild.WebhookScopeConfigurationArgs{
			Name:  pulumi.String(spec.ScopeConfiguration.Name),
			Scope: pulumi.String(spec.ScopeConfiguration.Scope),
		}
		if spec.ScopeConfiguration.Domain != "" {
			scArgs.Domain = pulumi.StringPtr(spec.ScopeConfiguration.Domain)
		}
		args.ScopeConfiguration = scArgs
	}

	// Comment-approval gate for PR-triggered builds -- protects CI secrets
	// from untrusted fork code.
	if spec.PullRequestBuildPolicy != nil {
		prArgs := &codebuild.WebhookPullRequestBuildPolicyArgs{
			RequiresCommentApproval: pulumi.String(spec.PullRequestBuildPolicy.RequiresCommentApproval),
		}
		if len(spec.PullRequestBuildPolicy.ApproverRoles) > 0 {
			prArgs.ApproverRoles = pulumi.ToStringArray(spec.PullRequestBuildPolicy.ApproverRoles)
		}
		args.PullRequestBuildPolicy = prArgs
	}

	created, err := codebuild.NewWebhook(ctx, "codebuild-webhook", args,
		pulumi.Provider(provider),
		pulumi.DependsOn([]pulumi.Resource{proj}),
	)
	if err != nil {
		return nil, errors.Wrap(err, "create codebuild webhook")
	}

	return created, nil
}

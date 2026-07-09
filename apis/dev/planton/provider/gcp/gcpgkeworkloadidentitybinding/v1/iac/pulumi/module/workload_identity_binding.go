package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// workloadIdentityBinding creates one ADDITIVE IAM grant on the Google
// Service Account: the workload-identity principal receives
// roles/iam.workloadIdentityUser — the GCP half of GKE Workload Identity.
// The Kubernetes half (the iam.gke.io/gcp-service-account annotation on the
// KSA) belongs to the workload's own deployment.
func workloadIdentityBinding(ctx *pulumi.Context,
	locals *Locals,
	gcpProvider *gcp.Provider) error {

	spec := locals.GcpGkeWorkloadIdentityBinding.Spec

	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project. The member string needs a concrete
	// project, so the fallback is made concrete by reading the provider's
	// resolved project from the client config (the Pulumi counterpart of
	// the Terraform module's google_client_config data source).
	poolProject := locals.PoolProject
	if poolProject == "" {
		clientConfig, err := organizations.GetClientConfig(ctx, pulumi.Provider(gcpProvider))
		if err != nil {
			return errors.Wrap(err, "failed to read provider client config for the default project")
		}
		if clientConfig.Project == "" {
			return errors.New("project_id is empty and the provider has no default project configured")
		}
		poolProject = clientConfig.Project
	}

	// The workload-identity principal, constructed from its parts so a
	// typo'd member string is impossible by construction — identical
	// construction to the Terraform module.
	workloadIdentityMember := fmt.Sprintf(
		"serviceAccount:%s.svc.id.goog[%s/%s]",
		poolProject,
		spec.KsaNamespace,
		spec.KsaName,
	)

	args := &serviceaccount.IAMMemberArgs{
		ServiceAccountId: pulumi.String(locals.ServiceAccountId),
		Role:             pulumi.String("roles/iam.workloadIdentityUser"),
		Member:           pulumi.String(workloadIdentityMember),
	}

	// An IAM Condition is part of the grant's identity: the same grant with
	// and without a condition are two independent bindings in the policy.
	if spec.Condition != nil {
		conditionArgs := &serviceaccount.IAMMemberConditionArgs{
			Title:      pulumi.String(spec.Condition.Title),
			Expression: pulumi.String(spec.Condition.Expression),
		}
		if spec.Condition.Description != "" {
			conditionArgs.Description = pulumi.StringPtr(spec.Condition.Description)
		}
		args.Condition = conditionArgs
	}

	createdIamMember, err := serviceaccount.NewIAMMember(
		ctx,
		"workload-identity-binding",
		args,
		pulumi.Provider(gcpProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create IAM member for workload identity binding")
	}

	ctx.Export(OpMember, createdIamMember.Member)
	ctx.Export(OpServiceAccountEmail,
		pulumi.String(spec.ServiceAccountEmail.GetValue()))

	return nil
}

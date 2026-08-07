package module

import (
	"github.com/pkg/errors"
	gcpprojectv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpproject/v1alpha1"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/organizations"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// project provisions the Google Cloud project — the Layer-0 container
// every other GCP resource lives in. IAM grants are deliberately NOT
// bundled here; model them as first-class GcpProjectIamMember resources.
func project(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) (*organizations.Project, error) {
	spec := locals.GcpProject.Spec

	projectArgs := &organizations.ProjectArgs{
		Name:      pulumi.String(locals.DisplayName),
		ProjectId: pulumi.String(spec.ProjectId),
		Labels:    pulumi.ToStringMap(locals.GcpLabels),
		// False by default: deleting the auto-created "default" network is
		// a standard hardening step, and explicit GcpVpcNetwork resources
		// are the composable path.
		AutoCreateNetwork: pulumi.Bool(spec.GetAutoCreateNetwork()),
		DeletionPolicy:    pulumi.String(locals.DeletionPolicy),
	}

	if spec.BillingAccountId != "" {
		projectArgs.BillingAccount = pulumi.StringPtr(spec.BillingAccountId)
	}

	// Resource Manager tags bind at create time only; changing them
	// afterwards forces recreation (bind tag values out-of-band instead).
	if len(spec.Tags) > 0 {
		projectArgs.Tags = pulumi.ToStringMap(spec.Tags)
	}

	// Exactly one of org_id / folder_id is sent, selected by parent_type.
	if spec.ParentType == gcpprojectv1alpha1.GcpProjectParentType_organization {
		projectArgs.OrgId = pulumi.String(spec.ParentId)
	}
	if spec.ParentType == gcpprojectv1alpha1.GcpProjectParentType_folder {
		projectArgs.FolderId = pulumi.String(spec.ParentId)
	}

	createdProject, err := organizations.NewProject(ctx, spec.ProjectId, projectArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create GCP project")
	}

	return createdProject, nil
}

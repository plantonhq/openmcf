package module

import (
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/sagemaker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	OpDomainId                                 = "domain_id"
	OpDomainArn                                = "domain_arn"
	OpDomainUrl                                = "domain_url"
	OpHomeEfsFileSystemId                      = "home_efs_file_system_id"
	OpSecurityGroupIdForDomainBoundary         = "security_group_id_for_domain_boundary"
	OpSingleSignOnApplicationArn               = "single_sign_on_application_arn"
	OpSingleSignOnManagedApplicationInstanceId = "single_sign_on_managed_application_instance_id"
	OpUserProfileArns                          = "user_profile_arns"
	OpSpaceArns                                = "space_arns"
	OpSpaceUrls                                = "space_urls"
)

func outputs(ctx *pulumi.Context, createdDomain *sagemaker.Domain, profileArns, spaceArns, spaceUrls pulumi.StringMap) {
	ctx.Export(OpDomainId, createdDomain.ID())
	ctx.Export(OpDomainArn, createdDomain.Arn)
	ctx.Export(OpDomainUrl, createdDomain.Url)
	ctx.Export(OpHomeEfsFileSystemId, createdDomain.HomeEfsFileSystemId)
	ctx.Export(OpSecurityGroupIdForDomainBoundary, createdDomain.SecurityGroupIdForDomainBoundary)
	// The two SSO outputs are only populated when auth_mode is SSO; empty
	// strings under IAM auth.
	ctx.Export(OpSingleSignOnApplicationArn, createdDomain.SingleSignOnApplicationArn)
	ctx.Export(OpSingleSignOnManagedApplicationInstanceId, createdDomain.SingleSignOnManagedApplicationInstanceId)
	// Folded satellite maps, keyed by the spec names (the import feed).
	ctx.Export(OpUserProfileArns, profileArns)
	ctx.Export(OpSpaceArns, spaceArns)
	ctx.Export(OpSpaceUrls, spaceUrls)
}

package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/redshiftserverless"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// customDomain fronts the workgroup's endpoint with a custom DNS name
// and an ACM certificate (AWS keeps one custom domain per workgroup).
// The CNAME record pointing the domain at the workgroup endpoint stays
// yours to manage; certificate renewals through ACM update the
// association's expiry in place. Returns the association (nil without a
// custom domain) so the expiry output keeps a stable shape.
func customDomain(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
	createdWorkgroup *redshiftserverless.Workgroup,
) (*redshiftserverless.CustomDomainAssociation, error) {
	spec := locals.AwsRedshiftServerlessWorkgroup.Spec
	if spec.CustomDomain == nil {
		return nil, nil
	}

	createdAssociation, err := redshiftserverless.NewCustomDomainAssociation(ctx, "custom-domain",
		&redshiftserverless.CustomDomainAssociationArgs{
			WorkgroupName:              createdWorkgroup.WorkgroupName,
			CustomDomainName:           pulumi.String(spec.CustomDomain.DomainName),
			CustomDomainCertificateArn: pulumi.String(spec.CustomDomain.CertificateArn.GetValue()),
		},
		pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{createdWorkgroup}))
	if err != nil {
		return nil, errors.Wrap(err, "failed to associate custom domain")
	}
	return createdAssociation, nil
}

package module

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cognito"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// domain provisions the pool's hosted-UI domain (one per pool -- honestly
// folded into the pool resource) and exports the domain join keys.
//
// Two shapes exist, distinguished by the presence of a dot in the domain:
//   - Prefix domain ("myapp-auth"): AWS serves the UI at
//     {prefix}.auth.{region}.amazoncognito.com; ready in about a minute.
//   - Custom domain ("auth.example.com"): AWS fronts the UI with a managed
//     CloudFront distribution; the deployer points DNS at the exported
//     distribution domain (creation can take tens of minutes on AWS's side).
func domain(ctx *pulumi.Context, locals *Locals, createdPool *cognito.UserPool, provider *aws.Provider) error {
	if locals.Spec.Domain == nil || locals.Spec.Domain.Domain == "" {
		// No domain configured -- export empty values so the output contract
		// stays shape-stable for downstream references.
		ctx.Export(OpUserPoolDomain, pulumi.String(""))
		ctx.Export(OpHostedUiUrl, pulumi.String(""))
		ctx.Export(OpCloudfrontDistribution, pulumi.String(""))
		ctx.Export(OpCloudfrontDistributionArn, pulumi.String(""))
		ctx.Export(OpCloudfrontHostedZoneId, pulumi.String(""))
		return nil
	}

	domainSpec := locals.Spec.Domain
	resourceName := locals.Target.Metadata.Name + "-domain"

	args := &cognito.UserPoolDomainArgs{
		Domain:     pulumi.String(domainSpec.Domain),
		UserPoolId: createdPool.ID(),
	}

	// A certificate is what makes a domain "custom" to AWS: its presence
	// switches the domain onto a CloudFront distribution. The spec's CEL
	// already requires it for dotted domains.
	isCustomDomain := strings.Contains(domainSpec.Domain, ".")
	if isCustomDomain && domainSpec.CertificateArn.GetValue() != "" {
		args.CertificateArn = pulumi.StringPtr(domainSpec.CertificateArn.GetValue())
	}

	// Only forward an explicit choice -- omitted means AWS's default (managed
	// login for new domains).
	if domainSpec.ManagedLoginVersion != nil {
		args.ManagedLoginVersion = pulumi.IntPtr(int(*domainSpec.ManagedLoginVersion))
	}

	created, err := cognito.NewUserPoolDomain(ctx, resourceName, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "failed to create Cognito user pool domain")
	}

	// The RAW domain string is the join key (ALB authenticate-cognito actions
	// take it as user_pool_domain); the full URL is a separate convenience
	// output for application configs.
	ctx.Export(OpUserPoolDomain, pulumi.String(domainSpec.Domain))
	if isCustomDomain {
		ctx.Export(OpHostedUiUrl, pulumi.Sprintf("https://%s", domainSpec.Domain))
	} else {
		ctx.Export(OpHostedUiUrl, pulumi.String(
			fmt.Sprintf("https://%s.auth.%s.amazoncognito.com", domainSpec.Domain, locals.Spec.Region)))
	}

	// CloudFront alias targets, meaningful only for custom domains (AWS
	// reports the distribution the domain rides on). For prefix domains the
	// provider returns the shared Cognito CloudFront value, which is not a
	// DNS alias target -- export empty strings to keep the outputs honest.
	if isCustomDomain {
		ctx.Export(OpCloudfrontDistribution, created.CloudfrontDistribution)
		ctx.Export(OpCloudfrontDistributionArn, created.CloudfrontDistributionArn)
		ctx.Export(OpCloudfrontHostedZoneId, created.CloudfrontDistributionZoneId)
	} else {
		ctx.Export(OpCloudfrontDistribution, pulumi.String(""))
		ctx.Export(OpCloudfrontDistributionArn, pulumi.String(""))
		ctx.Export(OpCloudfrontHostedZoneId, pulumi.String(""))
	}

	return nil
}

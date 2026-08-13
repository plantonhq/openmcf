package module

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/apigateway"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// domain creates the custom domain, its base-path mappings, and any
// private access associations, and exports outputs.
//
// Lifecycle facts the renders below depend on:
//   - the certificate reference fans in by endpoint type: EDGE domains
//     take certificate_arn (the us-east-1 cert), REGIONAL and PRIVATE
//     domains take regional_certificate_arn -- one spec field, wired to
//     the right provider argument here (the same fan-in in both
//     engines);
//   - uploaded certificate material follows the same fan-in
//     (certificate_name for EDGE, regional_certificate_name for
//     REGIONAL);
//   - domain create/update waits on DomainNameStatus AVAILABLE (up to
//     60 minutes upstream -- enhanced security policies trigger a
//     post-create update);
//   - base-path mapping creation retries briefly upstream to absorb
//     stage-propagation lag on freshly deployed stages;
//   - an access association is a Plugin Framework resource with ARN
//     identity and NO update -- every change replaces it, which is free.
func domain(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &apigateway.DomainNameArgs{
		DomainName: pulumi.String(spec.DomainName),
		Tags:       pulumi.ToStringMap(locals.AwsTags),
	}

	// Endpoint configuration. The bridge flattens the provider's max-1
	// types list to a singular string.
	endpointConfig := &apigateway.DomainNameEndpointConfigurationArgs{
		Types: pulumi.String(locals.EndpointType),
	}
	if spec.EndpointConfiguration != nil && spec.EndpointConfiguration.IpAddressType != "" {
		endpointConfig.IpAddressType = pulumi.String(spec.EndpointConfiguration.IpAddressType)
	}
	args.EndpointConfiguration = endpointConfig

	// Certificate fan-in by endpoint type (see the header comment).
	if spec.CertificateArn.GetValue() != "" {
		if locals.EndpointType == "EDGE" {
			args.CertificateArn = pulumi.String(spec.CertificateArn.GetValue())
		} else {
			args.RegionalCertificateArn = pulumi.String(spec.CertificateArn.GetValue())
		}
	}
	if spec.UploadedCertificate != nil {
		if locals.EndpointType == "EDGE" {
			args.CertificateName = pulumi.String(spec.UploadedCertificate.Name)
		} else {
			args.RegionalCertificateName = pulumi.String(spec.UploadedCertificate.Name)
		}
		args.CertificateBody = pulumi.String(spec.UploadedCertificate.Body)
		if spec.UploadedCertificate.Chain != "" {
			args.CertificateChain = pulumi.String(spec.UploadedCertificate.Chain)
		}
		args.CertificatePrivateKey = pulumi.String(spec.UploadedCertificate.PrivateKey)
	}

	if spec.EndpointAccessMode != "" {
		args.EndpointAccessMode = pulumi.String(spec.EndpointAccessMode)
	}
	if spec.SecurityPolicy != "" {
		args.SecurityPolicy = pulumi.String(spec.SecurityPolicy)
	}
	if spec.MutualTls != nil {
		mtls := &apigateway.DomainNameMutualTlsAuthenticationArgs{
			TruststoreUri: pulumi.String(spec.MutualTls.TruststoreUri),
		}
		if spec.MutualTls.TruststoreVersion != "" {
			mtls.TruststoreVersion = pulumi.String(spec.MutualTls.TruststoreVersion)
		}
		args.MutualTlsAuthentication = mtls
	}
	if spec.OwnershipVerificationCertificateArn.GetValue() != "" {
		args.OwnershipVerificationCertificateArn = pulumi.String(spec.OwnershipVerificationCertificateArn.GetValue())
	}
	if spec.RoutingMode != "" {
		args.RoutingMode = pulumi.String(spec.RoutingMode)
	}
	if spec.Policy != nil {
		policyJson, err := json.Marshal(spec.Policy.AsMap())
		if err != nil {
			return errors.Wrap(err, "marshal domain policy")
		}
		args.Policy = pulumi.String(string(policyJson))
	}

	created, err := apigateway.NewDomainName(ctx, "domain", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create domain name")
	}

	// Base-path mappings. Keyed by base path ("(root)" for the empty
	// path) -- the same keys as the Terraform for_each and the output
	// map.
	mappingIds := pulumi.StringMap{}
	for _, m := range spec.BasePathMappings {
		key := m.BasePath
		if key == "" {
			key = "(root)"
		}
		mappingArgs := &apigateway.BasePathMappingArgs{
			DomainName: created.DomainName,
			RestApi:    pulumi.String(m.RestApiId.GetValue()),
		}
		if m.BasePath != "" {
			mappingArgs.BasePath = pulumi.String(m.BasePath)
		}
		if m.StageName.GetValue() != "" {
			mappingArgs.StageName = pulumi.String(m.StageName.GetValue())
		}
		mapping, err := apigateway.NewBasePathMapping(ctx, "mapping-"+key, mappingArgs, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create base path mapping %q", key)
		}
		mappingIds[key] = mapping.ID().ToStringOutput()
	}
	ctx.Export(OpBasePathMappingIds, mappingIds)

	// Private access associations, keyed by the VPC endpoint they grant.
	associationArns := pulumi.StringMap{}
	for _, a := range spec.AccessAssociations {
		vpce := a.VpcEndpointId.GetValue()
		association, err := apigateway.NewDomainNameAccessAssociation(ctx, "access-"+vpce, &apigateway.DomainNameAccessAssociationArgs{
			AccessAssociationSource:     pulumi.String(vpce),
			AccessAssociationSourceType: pulumi.String("VPCE"),
			DomainNameArn:               created.Arn,
			Tags:                        pulumi.ToStringMap(locals.AwsTags),
		}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create access association %q", vpce)
		}
		associationArns[vpce] = association.Arn
	}
	ctx.Export(OpAccessAssociationArns, associationArns)

	ctx.Export(OpDomainName, created.DomainName)
	ctx.Export(OpDomainNameArn, created.Arn)
	ctx.Export(OpDomainNameId, created.DomainNameId)
	ctx.Export(OpRegionalDomainName, created.RegionalDomainName)
	ctx.Export(OpRegionalZoneId, created.RegionalZoneId)
	ctx.Export(OpCloudfrontDomainName, created.CloudfrontDomainName)
	ctx.Export(OpCloudfrontZoneId, created.CloudfrontZoneId)
	return nil
}

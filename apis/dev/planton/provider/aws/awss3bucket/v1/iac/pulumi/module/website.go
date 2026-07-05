package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// website creates the static-website-hosting satellite when configured and
// returns it (nil otherwise) so the website outputs can be exported from the
// resource that actually owns them.
func website(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	createdBucket *s3.BucketV2) (*s3.BucketWebsiteConfigurationV2, error) {
	spec := locals.Spec
	if spec.Website == nil {
		return nil, nil
	}

	args := &s3.BucketWebsiteConfigurationV2Args{
		Bucket: createdBucket.ID(),
	}

	// CEL guarantees index-mode XOR redirect-all-mode, so exactly one of
	// these two branches configures the resource.
	if spec.Website.IndexDocumentSuffix != "" {
		args.IndexDocument = &s3.BucketWebsiteConfigurationV2IndexDocumentArgs{
			Suffix: pulumi.String(spec.Website.IndexDocumentSuffix),
		}
		if spec.Website.ErrorDocumentKey != "" {
			args.ErrorDocument = &s3.BucketWebsiteConfigurationV2ErrorDocumentArgs{
				Key: pulumi.String(spec.Website.ErrorDocumentKey),
			}
		}
		if len(spec.Website.RoutingRules) > 0 {
			routingRules := s3.BucketWebsiteConfigurationV2RoutingRuleArray{}
			for _, r := range spec.Website.RoutingRules {
				rule := &s3.BucketWebsiteConfigurationV2RoutingRuleArgs{}
				if r.Condition != nil {
					condition := &s3.BucketWebsiteConfigurationV2RoutingRuleConditionArgs{}
					if r.Condition.HttpErrorCodeReturnedEquals != "" {
						condition.HttpErrorCodeReturnedEquals = pulumi.StringPtr(r.Condition.HttpErrorCodeReturnedEquals)
					}
					if r.Condition.KeyPrefixEquals != "" {
						condition.KeyPrefixEquals = pulumi.StringPtr(r.Condition.KeyPrefixEquals)
					}
					rule.Condition = condition
				}
				redirect := &s3.BucketWebsiteConfigurationV2RoutingRuleRedirectArgs{}
				if r.Redirect.HostName != "" {
					redirect.HostName = pulumi.StringPtr(r.Redirect.HostName)
				}
				if r.Redirect.HttpRedirectCode != "" {
					redirect.HttpRedirectCode = pulumi.StringPtr(r.Redirect.HttpRedirectCode)
				}
				if r.Redirect.Protocol != "" {
					redirect.Protocol = pulumi.StringPtr(r.Redirect.Protocol)
				}
				if r.Redirect.ReplaceKeyPrefixWith != "" {
					redirect.ReplaceKeyPrefixWith = pulumi.StringPtr(r.Redirect.ReplaceKeyPrefixWith)
				}
				if r.Redirect.ReplaceKeyWith != "" {
					redirect.ReplaceKeyWith = pulumi.StringPtr(r.Redirect.ReplaceKeyWith)
				}
				rule.Redirect = redirect
				routingRules = append(routingRules, rule)
			}
			args.RoutingRules = routingRules
		}
	}
	if spec.Website.RedirectAllRequestsTo != nil {
		redirectAll := &s3.BucketWebsiteConfigurationV2RedirectAllRequestsToArgs{
			HostName: pulumi.String(spec.Website.RedirectAllRequestsTo.HostName),
		}
		if spec.Website.RedirectAllRequestsTo.Protocol != "" {
			redirectAll.Protocol = pulumi.StringPtr(spec.Website.RedirectAllRequestsTo.Protocol)
		}
		args.RedirectAllRequestsTo = redirectAll
	}

	createdWebsite, err := s3.NewBucketWebsiteConfigurationV2(ctx, "website", args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to configure website hosting")
	}
	return createdWebsite, nil
}

// logging creates the server-access-logging satellite when configured.
func logging(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	createdBucket *s3.BucketV2) error {
	spec := locals.Spec
	if spec.Logging == nil {
		return nil
	}

	args := &s3.BucketLoggingV2Args{
		Bucket:       createdBucket.ID(),
		TargetBucket: pulumi.String(spec.Logging.TargetBucket.GetValue()),
		TargetPrefix: pulumi.String(spec.Logging.TargetPrefix),
	}

	// Partitioned key format makes access logs directly queryable via Athena
	// date partitions; the flat SimplePrefix format is the AWS default.
	if spec.Logging.PartitionedPrefixDateSource != "" {
		args.TargetObjectKeyFormat = &s3.BucketLoggingV2TargetObjectKeyFormatArgs{
			PartitionedPrefix: &s3.BucketLoggingV2TargetObjectKeyFormatPartitionedPrefixArgs{
				PartitionDateSource: pulumi.String(spec.Logging.PartitionedPrefixDateSource),
			},
		}
	}

	if _, err := s3.NewBucketLoggingV2(ctx, "logging", args, pulumi.Provider(provider)); err != nil {
		return errors.Wrap(err, "failed to configure access logging")
	}
	return nil
}

// cors creates the CORS satellite when the spec defines rules.
func cors(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	createdBucket *s3.BucketV2) error {
	spec := locals.Spec
	if len(spec.CorsRules) == 0 {
		return nil
	}

	corsRules := s3.BucketCorsConfigurationV2CorsRuleArray{}
	for _, r := range spec.CorsRules {
		rule := &s3.BucketCorsConfigurationV2CorsRuleArgs{
			AllowedMethods: pulumi.ToStringArray(r.AllowedMethods),
			AllowedOrigins: pulumi.ToStringArray(r.AllowedOrigins),
		}
		if r.Id != "" {
			rule.Id = pulumi.StringPtr(r.Id)
		}
		if len(r.AllowedHeaders) > 0 {
			rule.AllowedHeaders = pulumi.ToStringArray(r.AllowedHeaders)
		}
		if len(r.ExposeHeaders) > 0 {
			rule.ExposeHeaders = pulumi.ToStringArray(r.ExposeHeaders)
		}
		if r.MaxAgeSeconds > 0 {
			rule.MaxAgeSeconds = pulumi.IntPtr(int(r.MaxAgeSeconds))
		}
		corsRules = append(corsRules, rule)
	}

	if _, err := s3.NewBucketCorsConfigurationV2(ctx, "cors", &s3.BucketCorsConfigurationV2Args{
		Bucket:    createdBucket.ID(),
		CorsRules: corsRules,
	}, pulumi.Provider(provider)); err != nil {
		return errors.Wrap(err, "failed to configure CORS")
	}
	return nil
}

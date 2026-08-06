package module

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/sesv2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// EmailIdentityResult holds the created email identity outputs.
type EmailIdentityResult struct {
	Arn                pulumi.StringOutput
	EmailIdentity      pulumi.StringOutput
	IdentityType       pulumi.StringOutput
	VerificationStatus pulumi.StringOutput
	DkimTokens         pulumi.StringArrayOutput
}

// emailIdentity creates the SESv2 email identity and its mail-from, feedback,
// and policy satellites — mirroring the Terraform module's resource split.
func emailIdentity(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
) (*EmailIdentityResult, error) {
	spec := locals.AwsSesEmailIdentity.Spec

	args := &sesv2.EmailIdentityArgs{
		// The identity string IS the AWS identifier -- deliberately from the spec,
		// not metadata.name, because the identity must be the exact DNS name mail
		// is sent from.
		EmailIdentity: pulumi.String(spec.EmailIdentity),
		Tags:          pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.ConfigurationSet != nil && spec.ConfigurationSet.GetValue() != "" {
		args.ConfigurationSetName = pulumi.StringPtr(spec.ConfigurationSet.GetValue())
	}

	// DKIM signing: Easy DKIM and BYODKIM are mutually exclusive arms.
	if locals.Dkim != nil {
		dkimArgs := &sesv2.EmailIdentityDkimSigningAttributesArgs{}
		if locals.Byodkim {
			dkimArgs.DomainSigningPrivateKey = pulumi.StringPtr(locals.Dkim.GetDomainSigningPrivateKey())
			dkimArgs.DomainSigningSelector = pulumi.StringPtr(locals.Dkim.GetDomainSigningSelector())
		} else if locals.Dkim.GetNextSigningKeyLength() != "" {
			dkimArgs.NextSigningKeyLength = pulumi.StringPtr(locals.Dkim.GetNextSigningKeyLength())
		}
		args.DkimSigningAttributes = dkimArgs
	}

	createdIdentity, err := sesv2.NewEmailIdentity(
		ctx,
		spec.EmailIdentity,
		args,
		pulumi.Provider(provider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create ses email identity")
	}

	// Custom MAIL FROM domain -- created only when the manifest configures it.
	if spec.MailFrom != nil {
		mailFromArgs := &sesv2.EmailIdentityMailFromAttributesArgs{
			EmailIdentity:  createdIdentity.EmailIdentity,
			MailFromDomain: pulumi.StringPtr(spec.MailFrom.MailFromDomain),
		}
		behavior := spec.MailFrom.GetBehaviorOnMxFailure()
		if behavior != "" {
			mailFromArgs.BehaviorOnMxFailure = pulumi.StringPtr(behavior)
		} else {
			mailFromArgs.BehaviorOnMxFailure = pulumi.StringPtr("USE_DEFAULT_VALUE")
		}
		if _, err := sesv2.NewEmailIdentityMailFromAttributes(
			ctx,
			"mail-from",
			mailFromArgs,
			pulumi.Provider(provider),
			pulumi.DependsOn([]pulumi.Resource{createdIdentity}),
		); err != nil {
			return nil, errors.Wrap(err, "failed to create mail-from attributes")
		}
	}

	// Bounce/complaint email forwarding -- materialized only when the manifest
	// takes an explicit position.
	if spec.EmailForwardingEnabled != nil {
		if _, err := sesv2.NewEmailIdentityFeedbackAttributes(
			ctx,
			"feedback",
			&sesv2.EmailIdentityFeedbackAttributesArgs{
				EmailIdentity:          createdIdentity.EmailIdentity,
				EmailForwardingEnabled: pulumi.Bool(spec.GetEmailForwardingEnabled()),
			},
			pulumi.Provider(provider),
			pulumi.DependsOn([]pulumi.Resource{createdIdentity}),
		); err != nil {
			return nil, errors.Wrap(err, "failed to create feedback attributes")
		}
	}

	// Authorization policies -- one AWS sub-resource per named entry.
	for name, policy := range locals.Policies {
		policyJSON, err := json.Marshal(policy.Policy.AsMap())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to marshal policy %s", name)
		}
		if _, err := sesv2.NewEmailIdentityPolicy(
			ctx,
			fmt.Sprintf("policy-%s", name),
			&sesv2.EmailIdentityPolicyArgs{
				EmailIdentity: createdIdentity.EmailIdentity,
				PolicyName:    pulumi.String(name),
				Policy:        pulumi.String(policyJSON),
			},
			pulumi.Provider(provider),
			pulumi.DependsOn([]pulumi.Resource{createdIdentity}),
		); err != nil {
			return nil, errors.Wrapf(err, "failed to create identity policy %s", name)
		}
	}

	return &EmailIdentityResult{
		Arn:                createdIdentity.Arn,
		EmailIdentity:      createdIdentity.EmailIdentity,
		IdentityType:       createdIdentity.IdentityType,
		VerificationStatus: createdIdentity.VerificationStatus,
		DkimTokens:         createdIdentity.DkimSigningAttributes.Tokens(),
	}, nil
}

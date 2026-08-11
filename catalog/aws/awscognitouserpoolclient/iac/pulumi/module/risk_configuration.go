package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cognito"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// riskConfiguration creates the CLIENT-SCOPED aws_cognito_risk_configuration
// for this app client, overriding the pool-wide configuration set on the
// AwsCognitoUserPool spec. Requires the pool's threat protection to be
// active (advanced_security_mode AUDIT or ENFORCED) -- a cross-resource
// requirement AWS enforces at apply time.
func riskConfiguration(ctx *pulumi.Context, locals *Locals, createdClient *cognito.UserPoolClient, provider *aws.Provider) error {
	rc := locals.Spec.RiskConfiguration
	if rc == nil {
		return nil
	}

	args := &cognito.RiskConfigurationArgs{
		UserPoolId: pulumi.String(locals.Spec.UserPoolId.GetValue()),
		ClientId:   createdClient.ID(),
	}

	if at := rc.AccountTakeover; at != nil {
		// The provider requires the actions block; the spec's CEL requires at
		// least one action inside it.
		actions := &cognito.RiskConfigurationAccountTakeoverRiskConfigurationActionsArgs{}
		if at.LowAction != nil {
			actions.LowAction = &cognito.RiskConfigurationAccountTakeoverRiskConfigurationActionsLowActionArgs{
				EventAction: pulumi.String(at.LowAction.EventAction),
				Notify:      pulumi.Bool(at.LowAction.Notify),
			}
		}
		if at.MediumAction != nil {
			actions.MediumAction = &cognito.RiskConfigurationAccountTakeoverRiskConfigurationActionsMediumActionArgs{
				EventAction: pulumi.String(at.MediumAction.EventAction),
				Notify:      pulumi.Bool(at.MediumAction.Notify),
			}
		}
		if at.HighAction != nil {
			actions.HighAction = &cognito.RiskConfigurationAccountTakeoverRiskConfigurationActionsHighActionArgs{
				EventAction: pulumi.String(at.HighAction.EventAction),
				Notify:      pulumi.Bool(at.HighAction.Notify),
			}
		}

		accountTakeover := &cognito.RiskConfigurationAccountTakeoverRiskConfigurationArgs{
			Actions: actions,
		}

		if nc := at.NotifyConfiguration; nc != nil {
			notify := &cognito.RiskConfigurationAccountTakeoverRiskConfigurationNotifyConfigurationArgs{
				SourceArn: pulumi.String(nc.SourceArn.GetValue()),
			}
			if nc.From != "" {
				notify.From = pulumi.StringPtr(nc.From)
			}
			if nc.ReplyTo != "" {
				notify.ReplyTo = pulumi.StringPtr(nc.ReplyTo)
			}
			if nc.BlockEmail != nil {
				notify.BlockEmail = &cognito.RiskConfigurationAccountTakeoverRiskConfigurationNotifyConfigurationBlockEmailArgs{
					Subject:  pulumi.String(nc.BlockEmail.Subject),
					HtmlBody: pulumi.String(nc.BlockEmail.HtmlBody),
					TextBody: pulumi.String(nc.BlockEmail.TextBody),
				}
			}
			if nc.MfaEmail != nil {
				notify.MfaEmail = &cognito.RiskConfigurationAccountTakeoverRiskConfigurationNotifyConfigurationMfaEmailArgs{
					Subject:  pulumi.String(nc.MfaEmail.Subject),
					HtmlBody: pulumi.String(nc.MfaEmail.HtmlBody),
					TextBody: pulumi.String(nc.MfaEmail.TextBody),
				}
			}
			if nc.NoActionEmail != nil {
				notify.NoActionEmail = &cognito.RiskConfigurationAccountTakeoverRiskConfigurationNotifyConfigurationNoActionEmailArgs{
					Subject:  pulumi.String(nc.NoActionEmail.Subject),
					HtmlBody: pulumi.String(nc.NoActionEmail.HtmlBody),
					TextBody: pulumi.String(nc.NoActionEmail.TextBody),
				}
			}
			accountTakeover.NotifyConfiguration = notify
		}

		args.AccountTakeoverRiskConfiguration = accountTakeover
	}

	if cc := rc.CompromisedCredentials; cc != nil {
		compromised := &cognito.RiskConfigurationCompromisedCredentialsRiskConfigurationArgs{
			Actions: &cognito.RiskConfigurationCompromisedCredentialsRiskConfigurationActionsArgs{
				EventAction: pulumi.String(cc.EventAction),
			},
		}
		// Empty means AWS's default (all supported events) -- send absence.
		if len(cc.EventFilter) > 0 {
			compromised.EventFilters = pulumi.ToStringArray(cc.EventFilter)
		}
		args.CompromisedCredentialsRiskConfiguration = compromised
	}

	if re := rc.RiskException; re != nil {
		exception := &cognito.RiskConfigurationRiskExceptionConfigurationArgs{}
		if len(re.BlockedIpRanges) > 0 {
			exception.BlockedIpRangeLists = pulumi.ToStringArray(re.BlockedIpRanges)
		}
		if len(re.SkippedIpRanges) > 0 {
			exception.SkippedIpRangeLists = pulumi.ToStringArray(re.SkippedIpRanges)
		}
		args.RiskExceptionConfiguration = exception
	}

	if _, err := cognito.NewRiskConfiguration(ctx,
		locals.Target.Metadata.Name,
		args,
		pulumi.Provider(provider),
		pulumi.Parent(createdClient)); err != nil {
		return errors.Wrap(err, "failed to create risk configuration")
	}

	return nil
}

package verify

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apprunner"
	apprunnertypes "github.com/aws/aws-sdk-go-v2/service/apprunner/types"
	pkgerrors "github.com/pkg/errors"
)

// App Runner resources have lifecycle-state semantics rather than hard
// deletes: a deleted service stays describable with status DELETED, and the
// versioned companions (auto scaling configuration, VPC connector,
// observability configuration) flip to INACTIVE on deletion while their
// describe calls keep returning them. Every verifier here therefore treats
// the terminal lifecycle states as "absent" -- the NAT-gateway/ECS class of
// state-aware verification.
//
// Status casing is compared case-insensitively throughout: the App Runner
// API is inconsistent across its own resources (auto scaling configurations
// return lowercase "active" while the SDK enum constant is "ACTIVE"; the
// other resources return uppercase), so a strict == against the enum
// constant silently fails on real responses.

func apprunnerClient(cfg aws.Config, region string) *apprunner.Client {
	if region != "" {
		cfg.Region = region
	}
	return apprunner.NewFromConfig(cfg)
}

func isAppRunnerNotFound(err error) bool {
	var notFound *apprunnertypes.ResourceNotFoundException
	return pkgerrors.As(err, &notFound)
}

// appRunnerServiceVerifier verifies an AwsAppRunnerService via
// DescribeService, keyed on the service_arn output.
type appRunnerServiceVerifier struct{}

func (*appRunnerServiceVerifier) IDOutputKey() string { return "service_arn" }

func (*appRunnerServiceVerifier) VerifyExists(ctx context.Context, cfg aws.Config, arn, region string) error {
	status, found, err := appRunnerServiceStatus(ctx, cfg, arn, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsapprunnerservice verify-exists failed for %q", arn)
	}
	if !found || status == apprunnertypes.ServiceStatusDeleted {
		return pkgerrors.Errorf("awsapprunnerservice %q not found (or deleted) after deploy", arn)
	}
	return nil
}

func (*appRunnerServiceVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, arn, region string) error {
	status, found, err := appRunnerServiceStatus(ctx, cfg, arn, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsapprunnerservice verify-absent failed for %q", arn)
	}
	if found && status != apprunnertypes.ServiceStatusDeleted {
		return pkgerrors.Errorf("awsapprunnerservice %q still exists after destroy (status %q)", arn, status)
	}
	return nil
}

func appRunnerServiceStatus(ctx context.Context, cfg aws.Config, arn, region string) (apprunnertypes.ServiceStatus, bool, error) {
	out, err := apprunnerClient(cfg, region).DescribeService(ctx, &apprunner.DescribeServiceInput{
		ServiceArn: aws.String(arn),
	})
	if err != nil {
		if isAppRunnerNotFound(err) {
			return "", false, nil
		}
		return "", false, err
	}
	return out.Service.Status, true, nil
}

// appRunnerAutoScalingConfigurationVerifier verifies an
// AwsAppRunnerAutoScalingConfiguration via DescribeAutoScalingConfiguration,
// keyed on the configuration_arn output. Deletion flips the revision to
// INACTIVE rather than removing it.
type appRunnerAutoScalingConfigurationVerifier struct{}

func (*appRunnerAutoScalingConfigurationVerifier) IDOutputKey() string { return "configuration_arn" }

func (*appRunnerAutoScalingConfigurationVerifier) VerifyExists(ctx context.Context, cfg aws.Config, arn, region string) error {
	active, found, err := appRunnerAutoScalingConfigurationActive(ctx, cfg, arn, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsapprunnerautoscalingconfiguration verify-exists failed for %q", arn)
	}
	if !found || !active {
		return pkgerrors.Errorf("awsapprunnerautoscalingconfiguration %q not found (or inactive) after deploy", arn)
	}
	return nil
}

func (*appRunnerAutoScalingConfigurationVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, arn, region string) error {
	active, found, err := appRunnerAutoScalingConfigurationActive(ctx, cfg, arn, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsapprunnerautoscalingconfiguration verify-absent failed for %q", arn)
	}
	if found && active {
		return pkgerrors.Errorf("awsapprunnerautoscalingconfiguration %q still active after destroy", arn)
	}
	return nil
}

func appRunnerAutoScalingConfigurationActive(ctx context.Context, cfg aws.Config, arn, region string) (bool, bool, error) {
	out, err := apprunnerClient(cfg, region).DescribeAutoScalingConfiguration(ctx, &apprunner.DescribeAutoScalingConfigurationInput{
		AutoScalingConfigurationArn: aws.String(arn),
	})
	if err != nil {
		if isAppRunnerNotFound(err) {
			return false, false, nil
		}
		return false, false, err
	}
	return strings.EqualFold(string(out.AutoScalingConfiguration.Status), string(apprunnertypes.AutoScalingConfigurationStatusActive)), true, nil
}

// appRunnerVpcConnectorVerifier verifies an AwsAppRunnerVpcConnector via
// DescribeVpcConnector, keyed on the vpc_connector_arn output. Deletion
// flips the connector to INACTIVE rather than removing it.
type appRunnerVpcConnectorVerifier struct{}

func (*appRunnerVpcConnectorVerifier) IDOutputKey() string { return "vpc_connector_arn" }

func (*appRunnerVpcConnectorVerifier) VerifyExists(ctx context.Context, cfg aws.Config, arn, region string) error {
	active, found, err := appRunnerVpcConnectorActive(ctx, cfg, arn, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsapprunnervpcconnector verify-exists failed for %q", arn)
	}
	if !found || !active {
		return pkgerrors.Errorf("awsapprunnervpcconnector %q not found (or inactive) after deploy", arn)
	}
	return nil
}

func (*appRunnerVpcConnectorVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, arn, region string) error {
	active, found, err := appRunnerVpcConnectorActive(ctx, cfg, arn, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsapprunnervpcconnector verify-absent failed for %q", arn)
	}
	if found && active {
		return pkgerrors.Errorf("awsapprunnervpcconnector %q still active after destroy", arn)
	}
	return nil
}

func appRunnerVpcConnectorActive(ctx context.Context, cfg aws.Config, arn, region string) (bool, bool, error) {
	out, err := apprunnerClient(cfg, region).DescribeVpcConnector(ctx, &apprunner.DescribeVpcConnectorInput{
		VpcConnectorArn: aws.String(arn),
	})
	if err != nil {
		if isAppRunnerNotFound(err) {
			return false, false, nil
		}
		return false, false, err
	}
	return strings.EqualFold(string(out.VpcConnector.Status), string(apprunnertypes.VpcConnectorStatusActive)), true, nil
}

// appRunnerObservabilityConfigurationVerifier verifies an
// AwsAppRunnerObservabilityConfiguration via
// DescribeObservabilityConfiguration, keyed on the configuration_arn output.
// Deletion flips the revision to INACTIVE rather than removing it.
type appRunnerObservabilityConfigurationVerifier struct{}

func (*appRunnerObservabilityConfigurationVerifier) IDOutputKey() string {
	return "configuration_arn"
}

func (*appRunnerObservabilityConfigurationVerifier) VerifyExists(ctx context.Context, cfg aws.Config, arn, region string) error {
	active, found, err := appRunnerObservabilityConfigurationActive(ctx, cfg, arn, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsapprunnerobservabilityconfiguration verify-exists failed for %q", arn)
	}
	if !found || !active {
		return pkgerrors.Errorf("awsapprunnerobservabilityconfiguration %q not found (or inactive) after deploy", arn)
	}
	return nil
}

func (*appRunnerObservabilityConfigurationVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, arn, region string) error {
	active, found, err := appRunnerObservabilityConfigurationActive(ctx, cfg, arn, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsapprunnerobservabilityconfiguration verify-absent failed for %q", arn)
	}
	if found && active {
		return pkgerrors.Errorf("awsapprunnerobservabilityconfiguration %q still active after destroy", arn)
	}
	return nil
}

func appRunnerObservabilityConfigurationActive(ctx context.Context, cfg aws.Config, arn, region string) (bool, bool, error) {
	out, err := apprunnerClient(cfg, region).DescribeObservabilityConfiguration(ctx, &apprunner.DescribeObservabilityConfigurationInput{
		ObservabilityConfigurationArn: aws.String(arn),
	})
	if err != nil {
		if isAppRunnerNotFound(err) {
			return false, false, nil
		}
		return false, false, err
	}
	return strings.EqualFold(string(out.ObservabilityConfiguration.Status), string(apprunnertypes.ObservabilityConfigurationStatusActive)), true, nil
}

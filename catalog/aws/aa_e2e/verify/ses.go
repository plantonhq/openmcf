package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	pkgerrors "github.com/pkg/errors"
)

// sesConfigurationSetVerifier verifies an AwsSesConfigurationSet via
// GetConfigurationSet, keyed on configuration_set_name.
type sesConfigurationSetVerifier struct{}

func (*sesConfigurationSetVerifier) IDOutputKey() string { return "configuration_set_name" }

func (*sesConfigurationSetVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := sesConfigurationSetExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awssesconfigurationset verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awssesconfigurationset %q not found after deploy", id)
	}
	return nil
}

func (*sesConfigurationSetVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := sesConfigurationSetExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awssesconfigurationset verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awssesconfigurationset %q still exists after destroy", id)
	}
	return nil
}

func sesConfigurationSetExists(ctx context.Context, cfg aws.Config, name, region string) (bool, error) {
	client := sesv2.NewFromConfig(cfg, func(o *sesv2.Options) {
		if region != "" {
			o.Region = region
		}
	})
	_, err := client.GetConfigurationSet(ctx, &sesv2.GetConfigurationSetInput{
		ConfigurationSetName: aws.String(name),
	})
	if err == nil {
		return true, nil
	}
	var notFound *types.NotFoundException
	if errors.As(err, &notFound) {
		return false, nil
	}
	return false, err
}

// sesEmailIdentityVerifier verifies an AwsSesEmailIdentity via GetEmailIdentity,
// keyed on email_identity.
type sesEmailIdentityVerifier struct{}

func (*sesEmailIdentityVerifier) IDOutputKey() string { return "email_identity" }

func (*sesEmailIdentityVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := sesEmailIdentityExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awssesemailidentity verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awssesemailidentity %q not found after deploy", id)
	}
	return nil
}

func (*sesEmailIdentityVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := sesEmailIdentityExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awssesemailidentity verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awssesemailidentity %q still exists after destroy", id)
	}
	return nil
}

func sesEmailIdentityExists(ctx context.Context, cfg aws.Config, identity, region string) (bool, error) {
	client := sesv2.NewFromConfig(cfg, func(o *sesv2.Options) {
		if region != "" {
			o.Region = region
		}
	})
	_, err := client.GetEmailIdentity(ctx, &sesv2.GetEmailIdentityInput{
		EmailIdentity: aws.String(identity),
	})
	if err == nil {
		return true, nil
	}
	var notFound *types.NotFoundException
	if errors.As(err, &notFound) {
		return false, nil
	}
	return false, err
}

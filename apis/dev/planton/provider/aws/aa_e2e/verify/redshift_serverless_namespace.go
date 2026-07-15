package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshiftserverless"
	"github.com/aws/aws-sdk-go-v2/service/redshiftserverless/types"
	pkgerrors "github.com/pkg/errors"
)

// redshiftServerlessNamespaceVerifier verifies an
// AwsRedshiftServerlessNamespace via GetNamespace, keyed on the
// namespace_name output. A namespace mid-deletion stays describable with
// a DELETING status before the service starts returning the typed
// ResourceNotFoundException -- the same lifecycle class as the RDS-shaped
// kinds -- so existence is "described AND not deleting", and absence
// accepts either signal.
type redshiftServerlessNamespaceVerifier struct{}

func (*redshiftServerlessNamespaceVerifier) IDOutputKey() string { return "namespace_name" }

func (*redshiftServerlessNamespaceVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := redshiftServerlessNamespaceExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsredshiftserverlessnamespace verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsredshiftserverlessnamespace %q not found after deploy", id)
	}
	return nil
}

func (*redshiftServerlessNamespaceVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := redshiftServerlessNamespaceExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsredshiftserverlessnamespace verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsredshiftserverlessnamespace %q still exists after destroy", id)
	}
	return nil
}

// redshiftServerlessNamespaceExists reports whether the namespace is
// present and not already on its way out. A ResourceNotFoundException is
// treated as absent.
func redshiftServerlessNamespaceExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := redshiftserverless.NewFromConfig(cfg, func(o *redshiftserverless.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.GetNamespace(ctx, &redshiftserverless.GetNamespaceInput{NamespaceName: &id})
	if err != nil {
		var notFound *types.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	if out.Namespace == nil {
		return false, nil
	}
	if out.Namespace.Status == types.NamespaceStatusDeleting {
		return false, nil
	}
	return true, nil
}

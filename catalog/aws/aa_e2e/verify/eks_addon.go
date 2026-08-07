package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	pkgerrors "github.com/pkg/errors"
)

// eksAddonVerifier verifies an AwsEksAddon via DescribeAddon.
// DescribeAddon needs BOTH the cluster and add-on names, and the harness
// carries exactly one identifier per component -- so this verifier keys
// on the add-on ARN, which encodes both:
// arn:aws:eks:<region>:<account>:addon/<cluster>/<addon-name>/<uuid>.
// The typed ResourceNotFoundException is the absence signal (raised for a
// missing add-on AND for a missing cluster, which covers teardown order
// either way).
type eksAddonVerifier struct{}

func (*eksAddonVerifier) IDOutputKey() string { return "addon_arn" }

func (*eksAddonVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := eksAddonExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awseksaddon verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awseksaddon %q not found after deploy", id)
	}
	return nil
}

func (*eksAddonVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := eksAddonExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awseksaddon verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awseksaddon %q still exists after destroy", id)
	}
	return nil
}

func eksAddonExists(ctx context.Context, cfg aws.Config, arn, region string) (bool, error) {
	clusterName, addonName, err := parseAddonArn(arn)
	if err != nil {
		return false, err
	}
	client := eks.NewFromConfig(cfg, func(o *eks.Options) {
		if region != "" {
			o.Region = region
		}
	})
	_, err = client.DescribeAddon(ctx, &eks.DescribeAddonInput{
		ClusterName: &clusterName,
		AddonName:   &addonName,
	})
	if err != nil {
		var notFound *types.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// parseAddonArn extracts the cluster and add-on names from an add-on ARN's
// resource part: addon/<cluster>/<addon-name>/<uuid>.
func parseAddonArn(arn string) (clusterName string, addonName string, err error) {
	const marker = ":addon/"
	idx := strings.Index(arn, marker)
	if idx < 0 {
		return "", "", pkgerrors.Errorf("not an EKS add-on ARN: %q", arn)
	}
	parts := strings.Split(arn[idx+len(marker):], "/")
	if len(parts) < 2 {
		return "", "", pkgerrors.Errorf("malformed EKS add-on ARN: %q", arn)
	}
	return parts[0], parts[1], nil
}

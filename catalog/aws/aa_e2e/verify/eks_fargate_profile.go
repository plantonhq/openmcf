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

// eksFargateProfileVerifier verifies an AwsEksFargateProfile via
// DescribeFargateProfile. The API needs BOTH the cluster and profile
// names, and the harness carries exactly one identifier per component --
// so this verifier keys on the profile ARN, which encodes both:
// arn:aws:eks:<region>:<account>:fargateprofile/<cluster>/<profile>/<uuid>.
// The typed ResourceNotFoundException is the absence signal (raised for a
// missing profile AND for a missing cluster, which covers teardown order
// either way).
type eksFargateProfileVerifier struct{}

func (*eksFargateProfileVerifier) IDOutputKey() string { return "fargate_profile_arn" }

func (*eksFargateProfileVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := eksFargateProfileExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awseksfargateprofile verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awseksfargateprofile %q not found after deploy", id)
	}
	return nil
}

func (*eksFargateProfileVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := eksFargateProfileExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awseksfargateprofile verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awseksfargateprofile %q still exists after destroy", id)
	}
	return nil
}

func eksFargateProfileExists(ctx context.Context, cfg aws.Config, arn, region string) (bool, error) {
	clusterName, profileName, err := parseFargateProfileArn(arn)
	if err != nil {
		return false, err
	}
	client := eks.NewFromConfig(cfg, func(o *eks.Options) {
		if region != "" {
			o.Region = region
		}
	})
	_, err = client.DescribeFargateProfile(ctx, &eks.DescribeFargateProfileInput{
		ClusterName:        &clusterName,
		FargateProfileName: &profileName,
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

// parseFargateProfileArn extracts the cluster and profile names from a
// profile ARN's resource part: fargateprofile/<cluster>/<profile>/<uuid>.
func parseFargateProfileArn(arn string) (clusterName string, profileName string, err error) {
	const marker = ":fargateprofile/"
	idx := strings.Index(arn, marker)
	if idx < 0 {
		return "", "", pkgerrors.Errorf("not an EKS Fargate profile ARN: %q", arn)
	}
	parts := strings.Split(arn[idx+len(marker):], "/")
	if len(parts) < 2 {
		return "", "", pkgerrors.Errorf("malformed EKS Fargate profile ARN: %q", arn)
	}
	return parts[0], parts[1], nil
}

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

// eksAccessEntryVerifier verifies an AwsEksAccessEntry via
// DescribeAccessEntry. The API needs the cluster name and the PRINCIPAL
// ARN, and the harness carries exactly one identifier per component -- so
// this verifier keys on the entry ARN, which encodes both:
// arn:<partition>:eks:<region>:<account>:access-entry/<cluster>/<principal
// type>/<principal account>/<principal path+name...>/<uuid>. The principal
// ARN is reconstructed from those segments (joining the middle segments
// preserves IAM paths). The typed ResourceNotFoundException is the absence
// signal (raised for a missing entry AND for a missing cluster, which
// covers teardown order either way).
type eksAccessEntryVerifier struct{}

func (*eksAccessEntryVerifier) IDOutputKey() string { return "access_entry_arn" }

func (*eksAccessEntryVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := eksAccessEntryExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awseksaccessentry verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awseksaccessentry %q not found after deploy", id)
	}
	return nil
}

func (*eksAccessEntryVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := eksAccessEntryExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awseksaccessentry verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awseksaccessentry %q still exists after destroy", id)
	}
	return nil
}

func eksAccessEntryExists(ctx context.Context, cfg aws.Config, arn, region string) (bool, error) {
	clusterName, principalArn, err := parseAccessEntryArn(arn)
	if err != nil {
		return false, err
	}
	client := eks.NewFromConfig(cfg, func(o *eks.Options) {
		if region != "" {
			o.Region = region
		}
	})
	_, err = client.DescribeAccessEntry(ctx, &eks.DescribeAccessEntryInput{
		ClusterName:  &clusterName,
		PrincipalArn: &principalArn,
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

// parseAccessEntryArn extracts the cluster name and reconstructs the
// principal ARN from an access-entry ARN. Resource part:
// access-entry/<cluster>/<type>/<account>/<path+name...>/<uuid>, where
// <type> is "role" or "user". The partition is lifted from the entry ARN
// itself so non-standard partitions (GovCloud, China) round-trip.
func parseAccessEntryArn(arn string) (clusterName string, principalArn string, err error) {
	const marker = ":access-entry/"
	idx := strings.Index(arn, marker)
	if idx < 0 {
		return "", "", pkgerrors.Errorf("not an EKS access entry ARN: %q", arn)
	}
	arnParts := strings.SplitN(arn, ":", 3)
	if len(arnParts) < 3 || arnParts[0] != "arn" {
		return "", "", pkgerrors.Errorf("malformed EKS access entry ARN: %q", arn)
	}
	partition := arnParts[1]

	parts := strings.Split(arn[idx+len(marker):], "/")
	// cluster + type + account + at least one name segment + uuid.
	if len(parts) < 5 {
		return "", "", pkgerrors.Errorf("malformed EKS access entry ARN: %q", arn)
	}
	clusterName = parts[0]
	principalType := parts[1]
	account := parts[2]
	pathAndName := strings.Join(parts[3:len(parts)-1], "/")

	principalArn = "arn:" + partition + ":iam::" + account + ":" + principalType + "/" + pathAndName
	return clusterName, principalArn, nil
}

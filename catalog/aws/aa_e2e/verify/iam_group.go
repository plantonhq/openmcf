package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	pkgerrors "github.com/pkg/errors"
)

// iamGroupVerifier verifies an AwsIamGroup via GetGroup, keyed on the
// group name. IAM is a global service, so the region parameter is
// ignored. A deleted group returns the typed NoSuchEntity error, which
// is the "absent" signal; any other error is a genuine failure and
// must surface. GetGroup also returns the group's members, so a green
// destroy proves the declarative membership was unwound with the group
// (IAM refuses to delete a group that still has members or policies).
type iamGroupVerifier struct{}

func (*iamGroupVerifier) IDOutputKey() string { return "group_name" }

func (*iamGroupVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, _ string) error {
	exists, err := iamGroupExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsiamgroup verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsiamgroup %q not found after deploy", id)
	}
	return nil
}

func (*iamGroupVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, _ string) error {
	exists, err := iamGroupExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsiamgroup verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsiamgroup %q still exists after destroy", id)
	}
	return nil
}

func iamGroupExists(ctx context.Context, cfg aws.Config, groupName string) (bool, error) {
	client := iam.NewFromConfig(cfg)
	_, err := client.GetGroup(ctx, &iam.GetGroupInput{GroupName: aws.String(groupName)})
	if err != nil {
		if isIamNoSuchEntity(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

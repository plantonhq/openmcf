package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	pkgerrors "github.com/pkg/errors"
)

// ecrRepositoryVerifier verifies an AwsEcrRepo via DescribeRepositories,
// keyed on repository_name.
type ecrRepositoryVerifier struct{}

func (*ecrRepositoryVerifier) IDOutputKey() string { return "repository_name" }

func (*ecrRepositoryVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := ecrRepositoryExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsecrrepo verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsecrrepo %q not found after deploy", id)
	}
	return nil
}

func (*ecrRepositoryVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := ecrRepositoryExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsecrrepo verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsecrrepo %q still exists after destroy", id)
	}
	return nil
}

func ecrRepositoryExists(ctx context.Context, cfg aws.Config, repositoryName, region string) (bool, error) {
	client := ecr.NewFromConfig(cfg, func(o *ecr.Options) {
		if region != "" {
			o.Region = region
		}
	})
	_, err := client.DescribeRepositories(ctx, &ecr.DescribeRepositoriesInput{
		RepositoryNames: []string{repositoryName},
	})
	if err == nil {
		return true, nil
	}
	var notFound *types.RepositoryNotFoundException
	if errors.As(err, &notFound) {
		return false, nil
	}
	return false, err
}

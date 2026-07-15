package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	ecrtypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"
	pkgerrors "github.com/pkg/errors"
)

// ecrRepoVerifier verifies an AwsEcrRepo via DescribeRepositories, keyed on
// the repository_name output. ECR deletion is synchronous — a deleted
// repository immediately returns the typed RepositoryNotFoundException, so
// existence is a plain describe.
type ecrRepoVerifier struct{}

func (*ecrRepoVerifier) IDOutputKey() string { return "repository_name" }

func (*ecrRepoVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := ecrRepoExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsecrrepo verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsecrrepo %q not found after deploy", id)
	}
	return nil
}

func (*ecrRepoVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := ecrRepoExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsecrrepo verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsecrrepo %q still exists after destroy", id)
	}
	return nil
}

func ecrRepoExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := ecr.NewFromConfig(cfg, func(o *ecr.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeRepositories(ctx, &ecr.DescribeRepositoriesInput{
		RepositoryNames: []string{id},
	})
	if err != nil {
		var notFound *ecrtypes.RepositoryNotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	return len(out.Repositories) > 0, nil
}

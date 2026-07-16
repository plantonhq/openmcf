package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codebuild"
	pkgerrors "github.com/pkg/errors"
)

// codeBuildProjectVerifier verifies an AwsCodeBuildProject via
// BatchGetProjects.
//
// CodeBuild has no typed NotFound error for projects: BatchGetProjects
// succeeds and reports unresolved names in projectsNotFound instead, so
// existence is decided by whether the name resolved. Deletion is synchronous
// (DeleteProject is a single control-plane call), so no transitional
// deleting state needs handling.
type codeBuildProjectVerifier struct{}

func (*codeBuildProjectVerifier) IDOutputKey() string { return "project_name" }

func (*codeBuildProjectVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := codeBuildProjectExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscodebuildproject verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awscodebuildproject %q not found after deploy", id)
	}
	return nil
}

func (*codeBuildProjectVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := codeBuildProjectExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscodebuildproject verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awscodebuildproject %q still exists after destroy", id)
	}
	return nil
}

func codeBuildProjectExists(ctx context.Context, cfg aws.Config, name, region string) (bool, error) {
	client := codebuild.NewFromConfig(cfg, func(o *codebuild.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.BatchGetProjects(ctx, &codebuild.BatchGetProjectsInput{Names: []string{name}})
	if err != nil {
		return false, err
	}
	return len(out.Projects) > 0, nil
}

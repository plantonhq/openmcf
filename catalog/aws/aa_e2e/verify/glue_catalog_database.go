package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	gluetypes "github.com/aws/aws-sdk-go-v2/service/glue/types"
	pkgerrors "github.com/pkg/errors"
)

// glueCatalogDatabaseVerifier verifies an AwsGlueCatalogDatabase via
// GetDatabase in the deploying account's own catalog (the catalog ID defaults
// server-side to the caller's account). Glue reports a missing database with
// the typed EntityNotFoundException; deletion is synchronous.
type glueCatalogDatabaseVerifier struct{}

func (*glueCatalogDatabaseVerifier) IDOutputKey() string { return "database_name" }

func (*glueCatalogDatabaseVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := glueCatalogDatabaseExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsgluecatalogdatabase verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsgluecatalogdatabase %q not found after deploy", id)
	}
	return nil
}

func (*glueCatalogDatabaseVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := glueCatalogDatabaseExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsgluecatalogdatabase verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsgluecatalogdatabase %q still exists after destroy", id)
	}
	return nil
}

func glueCatalogDatabaseExists(ctx context.Context, cfg aws.Config, name, region string) (bool, error) {
	client := glue.NewFromConfig(cfg, func(o *glue.Options) {
		if region != "" {
			o.Region = region
		}
	})
	_, err := client.GetDatabase(ctx, &glue.GetDatabaseInput{Name: &name})
	if err != nil {
		var notFound *gluetypes.EntityNotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

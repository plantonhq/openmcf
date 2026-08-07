package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// cloudSqlUserVerifier probes a database user on a Cloud SQL instance via the
// sqladmin API.
type cloudSqlUserVerifier struct{}

func (v *cloudSqlUserVerifier) IDOutputKey() string { return "user_name" }

func (v *cloudSqlUserVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	instance := outputs["instance_name"]
	if instance == "" {
		return errors.New("instance_name output missing after deploy — required to verify the user")
	}
	name := outputs["user_name"]
	if name == "" {
		return errors.New("user_name output missing after deploy")
	}

	if _, err := svc.SqlAdmin.Users.Get(svc.Project, instance, name).Context(ctx).Do(); err != nil {
		return errors.Wrapf(err, "cloud sql user %s on instance %s not found after deploy", name, instance)
	}
	return nil
}

func (v *cloudSqlUserVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	instance := outputs["instance_name"]
	name := outputs["user_name"]
	if instance == "" || name == "" {
		return nil
	}

	_, err := svc.SqlAdmin.Users.Get(svc.Project, instance, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing cloud sql user %s on instance %s after destroy", name, instance)
	}
	return errors.Errorf("cloud sql user %s on instance %s still exists after destroy", name, instance)
}

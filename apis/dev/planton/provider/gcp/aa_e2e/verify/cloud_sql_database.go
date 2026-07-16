package verify

import (
	"context"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

var databaseSelfLinkRE = regexp.MustCompile(`/instances/([^/]+)/databases/`)

// cloudSqlDatabaseVerifier probes a logical database inside a Cloud SQL
// instance via the sqladmin API.
type cloudSqlDatabaseVerifier struct{}

func (v *cloudSqlDatabaseVerifier) IDOutputKey() string { return "database_name" }

func (v *cloudSqlDatabaseVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["database_name"]
	if name == "" {
		return errors.New("database_name output missing after deploy")
	}

	instance := outputs["instance_name"]
	if instance == "" {
		instance = instanceFromDatabaseSelfLink(outputs["self_link"])
	}
	if instance == "" {
		return errors.New("could not resolve instance name from outputs — need instance_name or self_link")
	}

	if _, err := svc.SqlAdmin.Databases.Get(svc.Project, instance, name).Context(ctx).Do(); err != nil {
		return errors.Wrapf(err, "cloud sql database %s on instance %s not found after deploy", name, instance)
	}
	return nil
}

func (v *cloudSqlDatabaseVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["database_name"]
	if name == "" {
		return nil
	}

	instance := outputs["instance_name"]
	if instance == "" {
		instance = instanceFromDatabaseSelfLink(outputs["self_link"])
	}
	if instance == "" {
		return nil
	}

	_, err := svc.SqlAdmin.Databases.Get(svc.Project, instance, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing cloud sql database %s on instance %s after destroy", name, instance)
	}
	return errors.Errorf("cloud sql database %s on instance %s still exists after destroy", name, instance)
}

func instanceFromDatabaseSelfLink(selfLink string) string {
	if selfLink == "" {
		return ""
	}
	if m := databaseSelfLinkRE.FindStringSubmatch(selfLink); len(m) == 2 {
		return m[1]
	}
	if idx := strings.Index(selfLink, "/instances/"); idx >= 0 {
		rest := selfLink[idx+len("/instances/"):]
		if end := strings.Index(rest, "/"); end > 0 {
			return rest[:end]
		}
	}
	return ""
}

package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// composerUserWorkloadsConfigMapVerifier probes a Composer user
// workloads ConfigMap via the Composer API. The name output is the
// fully qualified resource name
// (projects/{p}/locations/{r}/environments/{e}/userWorkloadsConfigMaps/{n}).
// Posture assertion: the ConfigMap carries at least one data entry —
// proof the payload landed, not just the shell.
type composerUserWorkloadsConfigMapVerifier struct{}

func (v *composerUserWorkloadsConfigMapVerifier) IDOutputKey() string { return "name" }

func (v *composerUserWorkloadsConfigMapVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["name"]
	if name == "" {
		return errors.New("name output missing after deploy")
	}

	configMap, err := svc.Composer.Projects.Locations.Environments.UserWorkloadsConfigMaps.Get(name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "composer user workloads config map %s not found after deploy", name)
	}
	if len(configMap.Data) == 0 {
		return errors.Errorf("composer user workloads config map %s has no data entries after deploy", name)
	}
	return nil
}

func (v *composerUserWorkloadsConfigMapVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["name"]
	if name == "" {
		return nil
	}

	_, err := svc.Composer.Projects.Locations.Environments.UserWorkloadsConfigMaps.Get(name).Context(ctx).Do()
	if err == nil {
		return errors.Errorf("composer user workloads config map %s still exists after destroy", name)
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) && apiErr.Code == 404 {
		return nil
	}
	// When the whole chain (environment included) was destroyed, the API
	// answers for the missing parent instead of the config map.
	if apiErr != nil && apiErr.Code == 400 {
		return nil
	}
	return errors.Wrapf(err, "unexpected error probing composer user workloads config map %s after destroy", name)
}

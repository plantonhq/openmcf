package verify

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// redisInstanceVerifier probes a Memorystore for Redis instance via the redis
// locations API. Posture assertions confirm the instance is READY, that the
// endpoint output matches the live host, and that a live AUTH-enabled
// instance carried its AUTH string into the outputs.
type redisInstanceVerifier struct{}

func (v *redisInstanceVerifier) IDOutputKey() string { return "instance_name" }

// instancePath builds the projects/{p}/locations/{region}/instances/{name}
// resource path the redis API addresses instances by.
func (v *redisInstanceVerifier) instancePath(svc *Services, outputs map[string]string) (string, error) {
	name := outputs["instance_name"]
	region := outputs["region"]
	if name == "" || region == "" {
		return "", errors.New("instance_name or region output missing")
	}
	return fmt.Sprintf("projects/%s/locations/%s/instances/%s", svc.Project, region, name), nil
}

func (v *redisInstanceVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	path, err := v.instancePath(svc, outputs)
	if err != nil {
		return errors.Wrap(err, "after deploy")
	}

	inst, err := svc.Redis.Projects.Locations.Instances.Get(path).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "redis instance %s not found after deploy", path)
	}
	if inst.State != "READY" {
		return errors.Errorf("redis instance %s state is %q, want READY", path, inst.State)
	}
	if host := outputs["host"]; host != "" && inst.Host != host {
		return errors.Errorf("redis instance %s host mismatch: output %q, live %q", path, host, inst.Host)
	}
	// Assert live-state → output, never the reverse: Pulumi masks secret
	// outputs as "[secret]" when read without --show-secrets, so an
	// "empty" secret output is indistinguishable from a populated one and
	// can never anchor an assertion.
	if inst.AuthEnabled && outputs["auth_string"] == "" {
		return errors.Errorf("redis instance %s: AUTH is enabled live but the auth_string output is empty", path)
	}
	return nil
}

func (v *redisInstanceVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	path, err := v.instancePath(svc, outputs)
	if err != nil {
		return nil
	}

	_, err = svc.Redis.Projects.Locations.Instances.Get(path).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing redis instance %s after destroy", path)
	}
	return errors.Errorf("redis instance %s still exists after destroy", path)
}

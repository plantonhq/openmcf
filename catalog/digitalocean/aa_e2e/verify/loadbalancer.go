package verify

import (
	"context"

	"github.com/digitalocean/godo"
	pkgerrors "github.com/pkg/errors"
)

// loadBalancerVerifier verifies a DigitalOceanLoadBalancer via
// GET /v2/load_balancers/{id}.
type loadBalancerVerifier struct{}

func (*loadBalancerVerifier) IDOutputKey() string { return "load_balancer_id" }

func (*loadBalancerVerifier) VerifyExists(ctx context.Context, client *godo.Client, id string) error {
	exists, err := loadBalancerExists(ctx, client, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "digitaloceanloadbalancer verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("digitaloceanloadbalancer %q not found after deploy", id)
	}
	return nil
}

func (*loadBalancerVerifier) VerifyAbsent(ctx context.Context, client *godo.Client, id string) error {
	exists, err := loadBalancerExists(ctx, client, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "digitaloceanloadbalancer verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("digitaloceanloadbalancer %q still exists after destroy", id)
	}
	return nil
}

func loadBalancerExists(ctx context.Context, client *godo.Client, id string) (bool, error) {
	_, _, err := client.LoadBalancers.Get(ctx, id)
	if err != nil {
		if isNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

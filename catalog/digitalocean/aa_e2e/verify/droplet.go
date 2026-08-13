package verify

import (
	"context"
	"strconv"

	"github.com/digitalocean/godo"
	pkgerrors "github.com/pkg/errors"
)

// dropletVerifier verifies a DigitalOceanDroplet via GET /v2/droplets/{id}.
// Droplet ids are integers in the API; the stack output carries the decimal
// string form.
type dropletVerifier struct{}

func (*dropletVerifier) IDOutputKey() string { return "droplet_id" }

func (*dropletVerifier) VerifyExists(ctx context.Context, client *godo.Client, id string) error {
	exists, err := dropletExists(ctx, client, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "digitaloceandroplet verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("digitaloceandroplet %q not found after deploy", id)
	}
	return nil
}

func (*dropletVerifier) VerifyAbsent(ctx context.Context, client *godo.Client, id string) error {
	exists, err := dropletExists(ctx, client, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "digitaloceandroplet verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("digitaloceandroplet %q still exists after destroy", id)
	}
	return nil
}

func dropletExists(ctx context.Context, client *godo.Client, id string) (bool, error) {
	numericID, err := strconv.Atoi(id)
	if err != nil {
		return false, pkgerrors.Wrapf(err, "droplet id %q is not the API's integer id", id)
	}
	_, _, err = client.Droplets.Get(ctx, numericID)
	if err != nil {
		if isNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

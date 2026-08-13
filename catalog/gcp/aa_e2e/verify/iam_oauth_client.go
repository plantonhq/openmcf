package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// iamOauthClientVerifier probes a workforce OAuth client through the IAM
// API (projects/{project}/locations/{location}/oauthClients/{id}).
//
// Deleted clients are SOFT-deleted for ~30 days: the API may keep serving
// the resource with state DELETED instead of returning 404. VERIFY-CLN
// accepts either shape as gone.
type iamOauthClientVerifier struct{}

func (v *iamOauthClientVerifier) IDOutputKey() string { return "client_name" }

func (v *iamOauthClientVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["client_name"]
	client, err := svc.Iam.Projects.Locations.OauthClients.Get(name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "oauth client %s not found after deploy", name)
	}
	if client.State == "DELETED" {
		return errors.Errorf("oauth client %s is in DELETED state after deploy", name)
	}
	return nil
}

func (v *iamOauthClientVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["client_name"]
	client, err := svc.Iam.Projects.Locations.OauthClients.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing oauth client %s after destroy", name)
	}
	if client.State == "DELETED" {
		// The 30-day soft-delete window: the resource answers but is
		// deleted — the expected post-destroy state.
		return nil
	}
	return errors.Errorf("oauth client %s still active (state %s) after destroy", name, client.State)
}

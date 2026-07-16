// GKE bearer-token minting: an OAuth2 access token from a Google service-account key.
package kubetoken

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
)

// cloudPlatformScope is the OAuth2 scope GKE API servers accept for IAM-authenticated
// requests (the same scope gcloud requests for cluster access).
const cloudPlatformScope = "https://www.googleapis.com/auth/cloud-platform"

// GkeTokenOptions carries the Google service-account key that signs the token request.
type GkeTokenOptions struct {
	// ServiceAccountKeyJSON is the raw key-file JSON downloaded from the GCP console
	// or minted via `gcloud iam service-accounts keys create`.
	ServiceAccountKeyJSON string

	// TokenURL overrides Google's token endpoint. Empty in production; tests point it
	// at a local fake so token minting is verifiable fully offline.
	TokenURL string
}

// MintGkeToken exchanges the service-account key for a short-lived OAuth2 access token
// (Google issues ~1h tokens; the expiry comes from the token response, never assumed).
func MintGkeToken(ctx context.Context, opts GkeTokenOptions) (Token, error) {
	if opts.ServiceAccountKeyJSON == "" {
		return Token{}, errors.New("a service-account key is required to mint a GKE token")
	}

	cfg, err := google.JWTConfigFromJSON([]byte(opts.ServiceAccountKeyJSON), cloudPlatformScope)
	if err != nil {
		return Token{}, errors.Wrap(err, "parsing Google service-account key")
	}
	if opts.TokenURL != "" {
		cfg.TokenURL = opts.TokenURL
	}

	token, err := cfg.TokenSource(ctx).Token()
	if err != nil {
		return Token{}, errors.Wrap(err, "exchanging service-account JWT for an access token")
	}

	return Token{Value: token.AccessToken, ExpiresAt: token.Expiry}, nil
}

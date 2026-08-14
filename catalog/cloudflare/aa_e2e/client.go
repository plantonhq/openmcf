package aa_e2e

// Cloudflare v4 REST client for E2E verification: authenticated GETs with a
// typed exists/absent answer and 429-aware Retry-After backoff.

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	apiBaseURL = "https://api.cloudflare.com/client/v4"

	// maxAttempts bounds the 429 retry loop. Cloudflare rate-limits at 1,200
	// requests per 5 minutes per token; verification traffic is far below
	// that, so more than a handful of retries means something else is wrong.
	maxAttempts = 4

	// defaultRetryDelay applies when a 429 response carries no Retry-After
	// header. Doubles per attempt.
	defaultRetryDelay = 2 * time.Second
)

// Client is a lightweight HTTP client for the Cloudflare v4 REST API, scoped
// to what E2E verification needs: authenticated GETs with a typed
// exists/absent answer. It retries 429 responses honoring Retry-After --
// Cloudflare enforces API rate limits, so backoff is part of the client
// contract here, not an optimization.
type Client struct {
	apiToken   string
	accountID  string
	httpClient *http.Client
}

// NewClient returns a client for the given API token and account.
func NewClient(apiToken, accountID string) *Client {
	return &Client{
		apiToken:   apiToken,
		accountID:  accountID,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// AccountID returns the Cloudflare account the harness is scoped to. Account-
// scoped verifiers prepend it to their API paths.
func (c *Client) AccountID() string {
	return c.accountID
}

// cloudflareEnvelope is the standard v4 response wrapper. Only the pieces
// verification reads are modeled.
type cloudflareEnvelope struct {
	Success bool `json:"success"`
	Errors  []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

// ResourceExists reports whether a GET on the given path (relative to
// /client/v4/, e.g. "zones/abc123") finds a resource. 200 means present;
// 404 means absent; a 400 carrying Cloudflare error code 7003 ("could not
// route to that endpoint") also means absent -- Cloudflare answers a GET on
// a deleted or unknown object identifier that way on several endpoints
// instead of a clean 404. Anything else is a real error.
func (c *Client) ResourceExists(ctx context.Context, path string) (bool, error) {
	resp, body, err := c.get(ctx, path)
	if err != nil {
		return false, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	case http.StatusBadRequest:
		var envelope cloudflareEnvelope
		if json.Unmarshal(body, &envelope) == nil {
			for _, e := range envelope.Errors {
				if e.Code == 7003 {
					return false, nil
				}
			}
		}
		return false, errors.Errorf("GET %s returned 400: %s", path, body)
	default:
		return false, errors.Errorf("GET %s returned %d: %s", path, resp.StatusCode, body)
	}
}

// ResourceActive reports whether a GET on the given path finds a resource
// that is not soft-deleted. Some Cloudflare families (cloudflared tunnels,
// teamnet routes) never 404 a destroyed object: the GET keeps answering 200
// with a non-null `result.deleted_at` (the provider's own destroy check
// asserts deleted_at != nil rather than expecting an error). This probe is
// OPT-IN per verifier: it parses the standard v4 envelope, so it must never
// replace ResourceExists on endpoints that return raw bodies (e.g. the KV
// value endpoint).
func (c *Client) ResourceActive(ctx context.Context, path string) (bool, error) {
	resp, body, err := c.get(ctx, path)
	if err != nil {
		return false, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var envelope struct {
			Result struct {
				DeletedAt *string `json:"deleted_at"`
			} `json:"result"`
		}
		if err := json.Unmarshal(body, &envelope); err != nil {
			return false, errors.Errorf("GET %s returned 200 with an unparseable body: %s", path, body)
		}
		if envelope.Result.DeletedAt != nil && *envelope.Result.DeletedAt != "" {
			return false, nil
		}
		return true, nil
	case http.StatusNotFound:
		return false, nil
	case http.StatusBadRequest:
		var envelope cloudflareEnvelope
		if json.Unmarshal(body, &envelope) == nil {
			for _, e := range envelope.Errors {
				if e.Code == 7003 {
					return false, nil
				}
			}
		}
		return false, errors.Errorf("GET %s returned 400: %s", path, body)
	default:
		return false, errors.Errorf("GET %s returned %d: %s", path, resp.StatusCode, body)
	}
}

// VerifyConnectivity proves the token is valid and active using Cloudflare's
// purpose-built, side-effect-free endpoint, then proves the token can see the
// configured account. Failing fast here keeps credential problems from
// surfacing later as confusing mid-lane verification errors.
func (c *Client) VerifyConnectivity(ctx context.Context) error {
	resp, body, err := c.get(ctx, "user/tokens/verify")
	if err != nil {
		return errors.Wrap(err, "token verification request failed")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("token verification returned %d: %s", resp.StatusCode, body)
	}

	resp, body, err = c.get(ctx, "accounts/"+c.accountID)
	if err != nil {
		return errors.Wrap(err, "account lookup request failed")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("account %s lookup returned %d (is the token scoped to this account?): %s",
			c.accountID, resp.StatusCode, body)
	}
	return nil
}

// get performs one authenticated GET with 429-aware retries and returns the
// final response alongside its fully-read body (the response body is already
// closed).
func (c *Client) get(ctx context.Context, path string) (*http.Response, []byte, error) {
	url := fmt.Sprintf("%s/%s", apiBaseURL, path)

	for attempt := 1; ; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to build request for %s", path)
		}
		req.Header.Set("Authorization", "Bearer "+c.apiToken)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "GET %s failed", path)
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != http.StatusTooManyRequests {
			return resp, body, nil
		}
		if attempt >= maxAttempts {
			return nil, nil, errors.Errorf("GET %s rate-limited after %d attempts", path, attempt)
		}

		delay := defaultRetryDelay * time.Duration(1<<(attempt-1))
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			if seconds, err := strconv.Atoi(ra); err == nil && seconds > 0 {
				delay = time.Duration(seconds) * time.Second
			}
		}
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		}
	}
}

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
	"github.com/plantonhq/planton/catalog/cloudflare/aa_e2e/verify"
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

// ResourcePresent reports whether a GET finds a resource that is still
// present after applying opts. 404 and Cloudflare's 400/7003 unknown-object
// answer are always absent. A 200 is present unless SoftDeleted sees a
// non-null result.deleted_at or result.status matches AbsentStatuses.
// Parses the v4 envelope, so it must never replace ResourceExists on
// raw-body endpoints (e.g. the KV value endpoint).
func (c *Client) ResourcePresent(ctx context.Context, path string, opts verify.EnvelopePresence) (bool, error) {
	resp, body, err := c.get(ctx, path)
	if err != nil {
		return false, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		if !opts.SoftDeleted && len(opts.AbsentStatuses) == 0 {
			return true, nil
		}
		var envelope struct {
			Result struct {
				DeletedAt *string `json:"deleted_at"`
				Status    string  `json:"status"`
			} `json:"result"`
		}
		if err := json.Unmarshal(body, &envelope); err != nil {
			return false, errors.Errorf("GET %s returned 200 with an unparseable body: %s", path, body)
		}
		if opts.SoftDeleted && envelope.Result.DeletedAt != nil && *envelope.Result.DeletedAt != "" {
			return false, nil
		}
		for _, status := range opts.AbsentStatuses {
			if envelope.Result.Status == status {
				return false, nil
			}
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

// ResourceActive is the deleted_at-aware probe: a 200 whose envelope carries
// a non-null result.deleted_at counts as ABSENT. Kept as a one-line wrapper
// so existing tunnel/route registrations and the API interface do not move.
func (c *Client) ResourceActive(ctx context.Context, path string) (bool, error) {
	return c.ResourcePresent(ctx, path, verify.EnvelopePresence{SoftDeleted: true})
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

package helmcrds

import (
	"context"
	"io"
	"net/http"
	"time"
)

// bundleFetchTimeout bounds one bundle download. Upstream CRD bundles run to a
// few megabytes (Solr's is ~2 MB); a minute is generous without letting a
// hung host stall a plan indefinitely.
const bundleFetchTimeout = 60 * time.Second

// fetchBundle downloads the pinned upstream CRD bundle and splits it into
// documents. The bundle branch exists for charts that ship no CRDs at all and
// publish them beside the chart instead.
func fetchBundle(ctx context.Context, src Source) ([]string, error) {
	url := src.ResolvedBundleURL()
	ctx, cancel := context.WithTimeout(ctx, bundleFetchTimeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, bundleFetchFailure(src, 0, err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, bundleFetchFailure(src, 0, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, bundleFetchFailure(src, response.StatusCode, nil)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, bundleFetchFailure(src, 0, err)
	}
	return splitDocuments(string(body)), nil
}

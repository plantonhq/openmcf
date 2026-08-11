package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// certificateMapVerifier probes a Certificate Manager certificate map
// and confirms it exists. The map_id output is the full resource name
// (projects/{p}/locations/global/certificateMaps/{name}) — exactly what
// the Certificate Manager GET takes. Entries are deployed and destroyed
// with the map; entry posture belongs to the proof lane's live API
// reads.
type certificateMapVerifier struct{}

func (v *certificateMapVerifier) IDOutputKey() string { return "map_id" }

func (v *certificateMapVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["map_id"]
	if _, err := svc.CertificateManager.Projects.Locations.CertificateMaps.Get(name).Context(ctx).Do(); err != nil {
		return errors.Wrapf(err, "certificate map %s not found after deploy", name)
	}
	return nil
}

func (v *certificateMapVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["map_id"]
	_, err := svc.CertificateManager.Projects.Locations.CertificateMaps.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing certificate map %s after destroy", name)
	}
	return errors.Errorf("certificate map %s still exists after destroy", name)
}

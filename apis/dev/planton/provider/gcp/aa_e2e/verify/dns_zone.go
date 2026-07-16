package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// dnsZoneVerifier probes a Cloud DNS managed zone via the dns API.
// Posture assertions confirm nameservers were assigned and the zone's
// visibility matches its configuration (public zones have no private
// visibility config; private zones do).
type dnsZoneVerifier struct{}

func (v *dnsZoneVerifier) IDOutputKey() string { return "zone_name" }

func (v *dnsZoneVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	zoneName := outputs["zone_name"]
	if zoneName == "" {
		return errors.New("zone_name output missing after deploy")
	}

	zone, err := svc.DNS.ManagedZones.Get(svc.Project, zoneName).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "dns managed zone %s not found after deploy", zoneName)
	}
	if len(zone.NameServers) == 0 {
		return errors.Errorf("dns managed zone %s has no nameservers after deploy", zoneName)
	}

	switch zone.Visibility {
	case "public":
		if zone.PrivateVisibilityConfig != nil {
			return errors.Errorf("dns managed zone %s visibility is public but private_visibility_config is set", zoneName)
		}
	case "private":
		if zone.PrivateVisibilityConfig == nil {
			return errors.Errorf("dns managed zone %s visibility is private but private_visibility_config is missing", zoneName)
		}
	default:
		return errors.Errorf("dns managed zone %s visibility is %q, want public or private", zoneName, zone.Visibility)
	}

	if wantVis := outputs["visibility"]; wantVis != "" && zone.Visibility != wantVis {
		return errors.Errorf("dns managed zone %s visibility mismatch: output %q, live %q", zoneName, wantVis, zone.Visibility)
	}
	return nil
}

func (v *dnsZoneVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	zoneName := outputs["zone_name"]
	if zoneName == "" {
		return nil
	}

	_, err := svc.DNS.ManagedZones.Get(svc.Project, zoneName).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing dns managed zone %s after destroy", zoneName)
	}
	return errors.Errorf("dns managed zone %s still exists after destroy", zoneName)
}

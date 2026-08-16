// The zone-settings singleton verifiers. These kinds manage zone-scoped
// SETTINGS surfaces that always exist while the zone exists: the settings
// object predates the deploy and survives the destroy (the provider's
// destroy is a NO-OP for the class -- measured at v5.23.0: zone_setting,
// waiting_room_settings, the cache toggles, universal_ssl_setting,
// total_tls, zone_auto_origin_tls_kex, certificate_authorities_hostname_
// associations all have literally empty Delete implementations).
//
// The verification contract that follows:
//
//   - VerifyExists asserts the kind's settings surface ANSWERS for the
//     deployed zone -- the deploy targeted a real, readable settings object.
//     Value-level assertion (did always_use_https actually become "on"?)
//     needs per-scenario expectations the harness does not carry; the proof
//     lane owns that design call and records it on the queue entry.
//   - VerifyAbsent asserts the surface STILL ANSWERS after destroy -- the
//     honest no-op-destroy posture. Asserting disappearance would fail every
//     honest run. Zero-orphan cleanup comes from the run-scoped zone fixture:
//     the abandoned settings die with the throwaway zone, so nothing outlives
//     the lane.
//
// One shared struct, not three bespoke ones: unlike the per-kind destroy
// contrasts seen elsewhere in the catalog, these three kinds carry exactly
// the same contract and differ only in which settings endpoint they probe.
package verify

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

// settingsSingletonVerifier verifies a zone-scoped settings surface whose
// object always exists while the zone does (no-op destroy class).
type settingsSingletonVerifier struct {
	component string
	// pathFormat is the settings surface's GET path with one %s for zone_id.
	pathFormat string
}

// IDOutputKey: settings singletons have no resource id -- the zone IS the
// identity (the fallback-origin precedent).
func (v *settingsSingletonVerifier) IDOutputKey() string { return "zone_id" }

func (v *settingsSingletonVerifier) VerifyExists(ctx context.Context, api API, outputs map[string]string) error {
	return v.assertSurfaceAnswers(ctx, api, outputs, "after deploy")
}

// VerifyAbsent asserts the settings surface still answers -- destroy is a
// no-op for this class, so "gone" would be a false expectation. Do not "fix"
// this into an absence check: it would fail every honest run.
func (v *settingsSingletonVerifier) VerifyAbsent(ctx context.Context, api API, outputs map[string]string) error {
	return v.assertSurfaceAnswers(ctx, api, outputs, "after destroy (settings surfaces persist -- destroy is a no-op for this class)")
}

func (v *settingsSingletonVerifier) assertSurfaceAnswers(ctx context.Context, api API, outputs map[string]string, when string) error {
	zoneID := outputs["zone_id"]
	if zoneID == "" {
		return errors.Errorf("%s outputs carry no zone_id -- cannot verify", v.component)
	}
	path := fmt.Sprintf(v.pathFormat, zoneID)
	exists, err := api.ResourceExists(ctx, path)
	if err != nil {
		return errors.Wrapf(err, "%s settings-surface probe failed %s", v.component, when)
	}
	if !exists {
		return errors.Errorf("%s settings surface for zone %s does not answer %s (GET %s)",
			v.component, zoneID, when, path)
	}
	return nil
}

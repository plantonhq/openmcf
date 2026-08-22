//go:build !codegen
// +build !codegen

// The E2E-proof join for the public report: read the provider's component
// E2E profiles and reduce them to the renderer's plain proof inputs. Lives
// beside the renderer (not inside it) so RenderPublicReport stays pure and
// the CLI and the drift gate share exactly one join -- two copies of this
// logic would eventually disagree about what "proven" means.

package providerparity

import (
	"os"
	"sort"
	"strings"

	"github.com/plantonhq/planton/pkg/crkreflect"
	"github.com/plantonhq/planton/pkg/e2e/profile"
	componentv1 "github.com/plantonhq/planton/qa/componente2eprofile/v1"
	sharedpb "github.com/plantonhq/planton/shared"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
)

// E2EProof is one kind's live-proof status, reduced from its E2E profile.
type E2EProof struct {
	// Green is the profile's recorded status: the kind's live lanes passed
	// as of the last proof run.
	Green bool
	// Engines are the validated provisioners, sorted (e.g. pulumi,
	// terraform).
	Engines []string
}

// Proven is THE definition of proven for every report surface: green live
// runs with both IaC engines validated. A green profile with one engine is
// progress, not proof.
func (p E2EProof) Proven() bool {
	return p.Green && len(p.Engines) >= 2
}

// BuildE2EProofs reads every component E2E profile of one provider and maps
// kind name -> proof status. A provider without an E2E harness (no
// aa_e2e/profile.yaml yet) yields nil proofs, not an error: the report then
// simply shows nothing proven, which is the truth.
func BuildE2EProofs(repoRoot string, provider cloudresourcekind.CloudResourceProvider) (map[string]E2EProof, error) {
	providerName := crkreflect.ProviderDirName(provider)
	if _, err := os.Stat(profile.ProviderProfilePath(repoRoot, providerName)); os.IsNotExist(err) {
		return nil, nil
	}
	result, err := profile.Discover(repoRoot, providerName, profile.FilterOpts{})
	if err != nil {
		return nil, err
	}
	proofs := map[string]E2EProof{}
	for _, ce := range result.Components {
		spec := ce.Profile.Spec
		if spec == nil {
			continue
		}
		kind := crkreflect.KindFromString(ce.Name)
		if kind == cloudresourcekind.CloudResourceKind_unspecified {
			continue
		}
		engines := make([]string, 0, len(spec.ValidatedProvisioners))
		for _, vp := range spec.ValidatedProvisioners {
			engines = append(engines, strings.ToLower(sharedpb.IacProvisioner_name[int32(vp)]))
		}
		sort.Strings(engines)
		proofs[kind.String()] = E2EProof{
			Green:   spec.Status == componentv1.ComponentE2EProfileSpec_green,
			Engines: engines,
		}
	}
	return proofs, nil
}

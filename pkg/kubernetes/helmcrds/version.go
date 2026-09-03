package helmcrds

import (
	"github.com/Masterminds/semver/v3"
)

// ExistingCRD is what a cluster already carries for one CRD this source owns:
// the name and the source version read from its stamp. A CRD with no version
// stamp (applied by something else, or by a release of this module older than
// the stamps) reports an empty Version and is never treated as a downgrade.
type ExistingCRD struct {
	Name    string
	Version string
}

// CheckNoDowngrade refuses when any existing CRD carries a strictly higher
// source version than the version about to be applied. Helm's own semver
// (Masterminds/semver) orders the versions, so prerelease and build metadata
// behave exactly as Helm's repository index would order them.
func CheckNoDowngrade(existing []ExistingCRD, requestedVersion string) error {
	requested, err := semver.NewVersion(requestedVersion)
	if err != nil {
		return UnparseableVersionFailure("requested", requestedVersion)
	}
	for _, crd := range existing {
		if crd.Version == "" {
			continue
		}
		current, err := semver.NewVersion(crd.Version)
		if err != nil {
			return UnparseableVersionFailure("cluster's "+crd.Name, crd.Version)
		}
		if requested.LessThan(current) {
			return DowngradeFailure(crd.Name, crd.Version, requestedVersion)
		}
	}
	return nil
}

package failure

import (
	"fmt"
	"strings"
)

// The Helm failures every chart install can meet, spoken once here so a Pulumi
// module (which meets them in-process through the Helm SDK) and a Terraform
// module (which meets them as the helm provider's raw text, explained by the
// layer that runs the engine) end in the same words. Helm's own text stays
// inside "observed" so nothing is hidden from a reader who knows Helm; the
// meaning and next step are what Helm never says.
//
// Each constructor's signature phrase is a fragment of its own observation;
// Explain uses it to stay silent when a module has already spoken.
const (
	signatureHelmVersionNotPublished    = "is not in the index of"
	signatureHelmOCIVersionNotPublished = "registry has no chart"
	signatureHelmRepositoryUnreachable  = "could not be reached from where this plan runs"
	signatureHelmStaleRepositoryCache   = "could not read a repository index on this machine"
)

// HelmVersionNotPublished: the pinned version is absent from an https
// repository's index.
func HelmVersionNotPublished(repository, chart, version, raw string) *Failure {
	return &Failure{
		Observed: fmt.Sprintf("chart %s %s %s %s (%s)", chart, version, signatureHelmVersionNotPublished, repository, raw),
		Meaning:  "the pinned chart_version has not been published to that repository, or the version string is misspelled",
		NextStep: fmt.Sprintf("set spec.chart_version to a version listed at %s/index.yaml, then re-run", strings.TrimSuffix(repository, "/")),
	}
}

// HelmOCIVersionNotPublished: the pinned version is absent from an OCI
// registry path. Registries have no index to read, so this is the only text.
func HelmOCIVersionNotPublished(repository, chart, version, raw string) *Failure {
	return &Failure{
		Observed: fmt.Sprintf("the OCI %s %s %s at %s (%s)", signatureHelmOCIVersionNotPublished, chart, version, repository, raw),
		Meaning:  "the pinned chart_version was never pushed to that registry path, or the version string is misspelled",
		NextStep: fmt.Sprintf("set spec.chart_version to a version the registry serves (`helm show chart %s/%s --version <v>` lists what exists), then re-run", strings.TrimSuffix(repository, "/"), chart),
	}
}

// HelmRepositoryUnreachable: the repository host does not resolve or egress to
// it is blocked from where the plan runs. The install needs the same path and
// would fail the same way; the render just fails first.
func HelmRepositoryUnreachable(repository, raw string) *Failure {
	return &Failure{
		Observed: fmt.Sprintf("the chart repository %s %s (%s)", repository, signatureHelmRepositoryUnreachable, raw),
		Meaning:  "the repository host does not resolve or egress to it is blocked; the release install needs the same path and would fail the same way, this check just fails first",
		NextStep: fmt.Sprintf("confirm DNS and egress from the runner to %s (`curl -I %s/index.yaml`), then re-run", repository, strings.TrimSuffix(repository, "/")),
	}
}

// HelmStaleRepositoryCache: an operator machine's Helm repository list names
// a repository whose index cache is missing. Helm consults that list even for
// a URL-addressed chart, so every render and install on the machine fails the
// same way; the runner starts from an empty list and never meets this.
func HelmStaleRepositoryCache(chart, raw string) *Failure {
	return &Failure{
		Observed: fmt.Sprintf("Helm %s while locating chart %s (%s)", signatureHelmStaleRepositoryCache, chart, raw),
		Meaning:  "the local Helm repository list has an entry whose index cache is missing; every render and every install on this machine fails the same way until it is repaired",
		NextStep: "run `helm repo update`, or remove the stale entry with `helm repo remove <name>`; the runner starts from an empty repository list and does not have this failure",
	}
}

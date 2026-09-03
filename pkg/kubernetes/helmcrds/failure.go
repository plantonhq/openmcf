package helmcrds

import (
	"fmt"
	"strings"
)

// Failure is the one shape every CRD-lifecycle error takes: what the platform
// observed (with the value), what it most likely means (one root cause), and
// the exact next step. The three lines are stable so an agent can match them;
// the values inside them vary.
type Failure struct {
	Observed string
	Meaning  string
	NextStep string
}

func (f *Failure) Error() string {
	return "observed: " + f.Observed + "\nmeaning: " + f.Meaning + "\nnext step: " + f.NextStep
}

// Raw substrings Helm produces for the failures the primitive anticipates. The
// provider's helm_template data source and the SDK are the same code, so one
// vocabulary classifies both; the Terraform block carries the same texts in
// its preconditions where it can see the cause.
const (
	helmVersionNotFound  = "not found in"
	helmRepoUnreachable  = "is not a valid chart repository or cannot be reached"
	helmOCINotFound      = "FetchReference"
	helmStaleRepoCache   = "no cached repo found"
	helmChartNotInRepo   = "chart \"" // "chart \"x\" version \"y\" not found in ..."
	helmRepoIndexMissing = "index.yaml"
)

// classifyLocateError turns Helm's chart-location error into a Failure. The
// raw Helm text is kept inside "observed" so nothing is hidden from a reader
// who knows Helm; the meaning and next step are what Helm never says.
func classifyLocateError(src Source, err error) error {
	raw := err.Error()
	chart := src.Chart + " " + src.Version
	switch {
	case strings.Contains(raw, helmStaleRepoCache):
		return &Failure{
			Observed: fmt.Sprintf("Helm could not read a repository index on this machine while locating chart %s (%s)", chart, raw),
			Meaning:  "the local Helm repository list has an entry whose index cache is missing; every render and every install on this machine fails the same way until it is repaired",
			NextStep: "run `helm repo update`, or remove the stale entry with `helm repo remove <name>`; the runner starts from an empty repository list and does not have this failure",
		}
	case strings.HasPrefix(src.Repository, "oci://") && strings.Contains(raw, helmOCINotFound):
		return &Failure{
			Observed: fmt.Sprintf("the OCI registry has no chart %s at %s (%s)", chart, src.Repository, raw),
			Meaning:  "the pinned chart_version was never pushed to that registry path, or the version string is misspelled",
			NextStep: fmt.Sprintf("set spec.chart_version to a version the registry serves (`helm show chart %s/%s --version <v>` lists what exists), then re-run", strings.TrimSuffix(src.Repository, "/"), src.Chart),
		}
	case strings.Contains(raw, helmVersionNotFound) && strings.Contains(raw, helmChartNotInRepo):
		return &Failure{
			Observed: fmt.Sprintf("chart %s is not in the index of %s (%s)", chart, src.Repository, raw),
			Meaning:  "the pinned chart_version has not been published to that repository, or the version string is misspelled",
			NextStep: fmt.Sprintf("set spec.chart_version to a version listed at %s/index.yaml, then re-run", strings.TrimSuffix(src.Repository, "/")),
		}
	case strings.Contains(raw, helmRepoUnreachable) || strings.Contains(raw, helmRepoIndexMissing):
		return &Failure{
			Observed: fmt.Sprintf("the chart repository %s could not be reached from where this plan runs (%s)", src.Repository, raw),
			Meaning:  "the repository host does not resolve or egress to it is blocked; the release install needs the same path and would fail the same way, this check just fails first",
			NextStep: fmt.Sprintf("confirm DNS and egress from the runner to %s (`curl -I %s/index.yaml`), then re-run", src.Repository, strings.TrimSuffix(src.Repository, "/")),
		}
	}
	return &Failure{
		Observed: fmt.Sprintf("Helm could not locate chart %s in %s (%s)", chart, src.Repository, raw),
		Meaning:  "the chart coordinates the module carries (repository, chart name, version) do not resolve from here",
		NextStep: "check spec.chart_version against the repository index and the module's chart identity, then re-run",
	}
}

func renderFailure(src Source, err error) error {
	return &Failure{
		Observed: fmt.Sprintf("rendering chart %s %s from %s failed (%s)", src.Chart, src.Version, src.Repository, err.Error()),
		Meaning:  "the chart's templates rejected the values the release would install with; the release install would fail identically",
		NextStep: "fix the value the error names (spec fields or helm_values), or pin a chart_version whose templates accept it",
	}
}

func zeroCRDsFailure(src Source) error {
	if src.IsBundle() {
		return &Failure{
			Observed: fmt.Sprintf("the CRD bundle at %s contained no CustomResourceDefinition documents", src.ResolvedBundleURL()),
			Meaning:  "upstream changed what that URL serves, or the version has no CRD bundle",
			NextStep: "open the URL, confirm it is the CRD bundle for this version, and correct the module's bundle URL template if upstream moved it",
		}
	}
	return &Failure{
		Observed: fmt.Sprintf("rendering chart %s %s with its CRD switch on (%s) produced no CustomResourceDefinition documents", src.Chart, src.Version, strings.TrimSpace(src.CRDOverride)),
		Meaning:  "the chart's CRD switch is not the value the module sets at this version (upstream renamed it), or the chart stopped shipping CRDs",
		NextStep: fmt.Sprintf("read the chart's values at this version (`helm show values %s --repo %s --version %s`) and update the module's CRD override to the switch the chart uses", src.Chart, src.Repository, src.Version),
	}
}

func bundleFetchFailure(src Source, status int, err error) error {
	url := src.ResolvedBundleURL()
	if err != nil {
		return &Failure{
			Observed: fmt.Sprintf("the CRD bundle at %s could not be fetched (%s)", url, err.Error()),
			Meaning:  "the bundle host does not resolve or egress to it is blocked from where this plan runs",
			NextStep: fmt.Sprintf("confirm DNS and egress to the bundle host (`curl -I %s`), then re-run", url),
		}
	}
	if status == 404 {
		return &Failure{
			Observed: fmt.Sprintf("the CRD bundle URL %s answered HTTP 404", url),
			Meaning:  "no bundle is published at this version, or upstream changed the URL pattern",
			NextStep: fmt.Sprintf("confirm version %s exists at the upstream download page; set spec.chart_version to a published version, or correct the module's bundle URL template", src.Version),
		}
	}
	return &Failure{
		Observed: fmt.Sprintf("the CRD bundle URL %s answered HTTP %d", url, status),
		Meaning:  "the upstream host is refusing or failing the request",
		NextStep: fmt.Sprintf("retry once the host answers 200 (`curl -I %s`); if it persists, the upstream download location has changed", url),
	}
}

// HelmManagedCRDsFailure is the refusal when a chart templates CRDs as ordinary
// release resources without Helm's keep annotation and the kind's policy has
// not accepted that. Helm would install and upgrade them with the release and
// delete them with it, taking every custom resource built on them; nothing the
// module applies can protect a resource Helm owns. The remedies are ordered
// from most to least protective.
func HelmManagedCRDsFailure(src Source, names []string) error {
	return &Failure{
		Observed: fmt.Sprintf("chart %s %s templates %d CustomResourceDefinition(s) as release resources without the %s: keep annotation (%s)", src.Chart, src.Version, len(names), HelmKeepAnnotation, strings.Join(names, ", ")),
		Meaning:  "Helm owns these CRDs: it installs and upgrades them with the release and deletes them when the release is uninstalled, along with every custom resource built on them; keep_on_uninstall cannot reach resources Helm owns",
		NextStep: "use the catalog's typed kind for this chart if one exists; otherwise turn on the chart's own keep switch in its values (cert-manager-style `crds.keep: true`) so the chart protects them itself; or set spec.crds.allow_helm_managed to true to accept that these CRDs live and die with the release",
	}
}

// DowngradeFailure is the refusal when the cluster already carries a CRD at a
// higher source version than the manifest asks for. Lowering a schema can
// strip fields from existing custom resources on their next write, so the
// primitive stops before touching anything.
func DowngradeFailure(crdName, existingVersion, requestedVersion string) error {
	return &Failure{
		Observed: fmt.Sprintf("the cluster carries CRD %s derived from chart version %s; the manifest asks for chart version %s", crdName, existingVersion, requestedVersion),
		Meaning:  "applying an older CRD schema over a newer one can strip fields from existing custom resources the next time anything writes them",
		NextStep: fmt.Sprintf("set spec.chart_version to %s or higher; if you accept the loss and want the older schema, delete the CRD first (`kubectl delete crd %s`) and re-run", existingVersion, crdName),
	}
}

// UnparseableVersionFailure covers a stamp or a manifest version that is not
// semver. Refusing is safer than guessing an order.
func UnparseableVersionFailure(which, value string) error {
	return &Failure{
		Observed: fmt.Sprintf("the %s chart version %q is not a semantic version", which, value),
		Meaning:  "the never-downgrade check cannot order versions it cannot parse, so it refuses rather than guess",
		NextStep: "use the exact version string the chart repository index publishes (for example 0.120.0)",
	}
}

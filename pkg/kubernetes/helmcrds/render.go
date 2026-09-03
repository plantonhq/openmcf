package helmcrds

import (
	"strings"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/registry"
)

// render templates the pinned chart client-side with CRDs included and
// returns every rendered document: the templated body AND the chart's crds/
// directory (Helm exposes the two on different surfaces; a kind may ship its
// CRDs on either). Nothing touches a cluster; the render only needs the chart
// repository, which the release install needs anyway.
//
// The Helm environment (repository config and cache) is deliberately the
// process's own, not an isolated one: the release install runs in the same
// environment, so a stale local repository entry must fail the render exactly
// as it would fail the install, and the Failure text carries the remedy.
func render(src Source, releaseName, namespace string) ([]string, error) {
	settings := cli.New()
	cfg := new(action.Configuration)
	registryClient, err := registry.NewClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create the Helm registry client")
	}
	cfg.RegistryClient = registryClient

	install := action.NewInstall(cfg)
	install.DryRun = true
	install.ClientOnly = true
	install.IncludeCRDs = true
	install.ReleaseName = releaseName
	install.Namespace = namespace
	install.ChartPathOptions.Version = src.Version
	if len(src.APIVersions) > 0 {
		install.APIVersions = chartutil.VersionSet(src.APIVersions)
	}
	if src.KubeVersion != "" {
		kubeVersion, err := chartutil.ParseKubeVersion(src.KubeVersion)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid KubeVersion %q", src.KubeVersion)
		}
		install.KubeVersion = kubeVersion
	}

	// OCI charts are addressed by a single reference; HTTP charts by
	// repository URL plus name. Helm's own CLI makes the same split.
	chartRef := src.Chart
	if strings.HasPrefix(src.Repository, "oci://") {
		chartRef = strings.TrimSuffix(src.Repository, "/") + "/" + src.Chart
	} else {
		install.ChartPathOptions.RepoURL = src.Repository
	}

	chartPath, err := install.ChartPathOptions.LocateChart(chartRef, settings)
	if err != nil {
		return nil, classifyLocateError(src, err)
	}
	chrt, err := loader.Load(chartPath)
	if err != nil {
		return nil, renderFailure(src, err)
	}

	values, err := mergeValueDocuments(append(append([]string{}, src.Values...), src.CRDOverride))
	if err != nil {
		return nil, renderFailure(src, err)
	}

	release, err := install.Run(chrt, values)
	if err != nil {
		return nil, renderFailure(src, err)
	}

	// With IncludeCRDs the SDK prepends the chart's crds/ directory to the
	// rendered manifest (verified against v3.20.2 with a crds/-directory
	// chart), so one manifest carries both surfaces; templated CRDs are in
	// it by construction.
	return splitDocuments(release.Manifest), nil
}

// mergeValueDocuments applies helm -f semantics: later documents override
// earlier ones, maps merge recursively, everything else replaces. Empty
// documents are skipped.
func mergeValueDocuments(documents []string) (map[string]interface{}, error) {
	merged := map[string]interface{}{}
	for _, doc := range documents {
		if strings.TrimSpace(doc) == "" {
			continue
		}
		values, err := chartutil.ReadValues([]byte(doc))
		if err != nil {
			return nil, errors.Wrap(err, "values document is not valid YAML")
		}
		merged = mergeMaps(merged, values.AsMap())
	}
	return merged, nil
}

func mergeMaps(base, overlay map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(base))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range overlay {
		if overlayMap, ok := v.(map[string]interface{}); ok {
			if baseMap, ok := out[k].(map[string]interface{}); ok {
				out[k] = mergeMaps(baseMap, overlayMap)
				continue
			}
		}
		out[k] = v
	}
	return out
}

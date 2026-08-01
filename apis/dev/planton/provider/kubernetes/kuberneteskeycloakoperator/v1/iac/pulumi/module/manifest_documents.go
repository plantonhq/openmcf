package module

import (
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// documentSeparator splits a multi-document YAML stream on `---` LINES —
// never the bare substring (the CRD schemas embed "---" inside
// description text) — the same boundary the Terraform twin's locals use,
// so both engines see an identical document set.
var documentSeparator = regexp.MustCompile(`(?m)^---\s*$`)

// manifestPartitions carries the release manifests split into the
// ordered apply groups. The partition exists because the groups need
// STRUCTURAL create/destroy ordering (see Resources): the bundle ships
// no Namespace document (the module authors one when create_namespace
// is set), and the CRDs are published as four separate files beside the
// bundle — applied LAST at create so destroy (the reverse) drains them
// FIRST while the operator still runs.
type manifestPartitions struct {
	// WorkloadsYaml carries the bundle's 16 documents: ServiceAccount,
	// ClusterRoles, bindings, Role, the metrics/health Service and the
	// operator Deployment (the watch-scope variant chosen by
	// spec.cluster_wide).
	WorkloadsYaml string
	// CrdsYaml carries the four k8s.keycloak.org CRDs — created last
	// (the JOSDK operator crash-loops until they exist), deleted FIRST
	// while the operator still runs.
	CrdsYaml string
}

// fetchManifestPartitions downloads the release manifests for the pinned
// tag and splits them into the ordered groups. Fetching in module code
// keeps the same posture as the Terraform twin's data.http: the
// manifests are read from the pinned, immutable tag at plan/apply time.
func fetchManifestPartitions(clusterWide bool) (*manifestPartitions, error) {
	bundleURL := BundleURL(clusterWide)
	raw, err := httpGet(bundleURL)
	if err != nil {
		return nil, err
	}

	var workloadDocs []string
	for _, document := range documentSeparator.Split(raw, -1) {
		if strings.TrimSpace(document) == "" {
			continue
		}
		kind, err := documentKind(document)
		if err != nil {
			return nil, err
		}
		switch kind {
		case "":
			// Comment-only fragments carry no kind and deploy nothing.
			continue
		case "Namespace", "CustomResourceDefinition":
			// The bundle ships neither: the module authors the
			// Namespace itself, and the CRDs live in four separate
			// files. Either kind appearing here means the pinned
			// asset's shape changed.
			return nil, errors.Errorf(
				"the operator bundle at %s unexpectedly carries a %s document — the pinned asset's shape changed",
				bundleURL, kind)
		default:
			workloadDocs = append(workloadDocs, document)
		}
	}
	if len(workloadDocs) == 0 {
		return nil, errors.Errorf(
			"the operator bundle at %s yielded no documents — the pinned asset's shape changed", bundleURL)
	}

	// The four CRD files are shared by both watch-scope variants and
	// carry exactly one CustomResourceDefinition each (plus a generator
	// comment header the splitter drops).
	var crdDocs []string
	for _, file := range vars.CrdFiles {
		crdURL := CrdURL(file)
		raw, err := httpGet(crdURL)
		if err != nil {
			return nil, err
		}
		found := 0
		for _, document := range documentSeparator.Split(raw, -1) {
			if strings.TrimSpace(document) == "" {
				continue
			}
			kind, err := documentKind(document)
			if err != nil {
				return nil, err
			}
			switch kind {
			case "":
				continue
			case "CustomResourceDefinition":
				crdDocs = append(crdDocs, document)
				found++
			default:
				return nil, errors.Errorf(
					"the CRD file at %s unexpectedly carries a %s document — the pinned asset's shape changed",
					crdURL, kind)
			}
		}
		if found != 1 {
			return nil, errors.Errorf(
				"the CRD file at %s yielded %d CustomResourceDefinition documents (expected exactly 1) — the pinned asset's shape changed",
				crdURL, found)
		}
	}

	return &manifestPartitions{
		WorkloadsYaml: strings.Join(workloadDocs, "\n---\n"),
		CrdsYaml:      strings.Join(crdDocs, "\n---\n"),
	}, nil
}

// namespaceDocumentYaml renders the module-authored Namespace document
// (the bundle ships none) carrying the standard Planton governance
// labels — the same document the Terraform twin's namespace_documents
// local builds, up to each engine's fleet-conventional label KEYS
// (Pulumi planton.ai/name|kind|id vs TF planton.ai/resource-*; the
// program-wide divergence, not this module's).
func namespaceDocumentYaml(locals *Locals) (string, error) {
	labels := map[string]interface{}{}
	for key, value := range locals.Labels {
		labels[key] = value
	}
	rendered, err := yaml.Marshal(map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Namespace",
		"metadata": map[string]interface{}{
			"name":   locals.Namespace,
			"labels": labels,
		},
	})
	if err != nil {
		return "", errors.Wrap(err, "rendering the namespace document")
	}
	return string(rendered), nil
}

func documentKind(document string) (string, error) {
	var meta struct {
		Kind string `yaml:"kind"`
	}
	if err := yaml.Unmarshal([]byte(document), &meta); err != nil {
		return "", errors.Wrap(err, "reading a manifest document's kind")
	}
	return meta.Kind, nil
}

func httpGet(url string) (string, error) {
	client := &http.Client{Timeout: 2 * time.Minute}
	response, err := client.Get(url)
	if err != nil {
		return "", errors.Wrapf(err, "fetching the release manifest from %s", url)
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != http.StatusOK {
		return "", errors.Errorf("fetching the release manifest from %s: HTTP %d", url, response.StatusCode)
	}
	raw, err := io.ReadAll(response.Body)
	if err != nil {
		return "", errors.Wrapf(err, "reading the release manifest body from %s", url)
	}
	return string(raw), nil
}

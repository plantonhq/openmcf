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

// documentSeparator splits a multi-document YAML stream on `---` lines —
// the same boundary the Terraform twin's locals use, so both engines see
// an identical document set.
var documentSeparator = regexp.MustCompile(`(?m)^---\s*$`)

// manifestPartitions carries the release manifest split into the three
// ordered apply groups. The partition exists because the groups need
// STRUCTURAL create/destroy ordering (see Resources) — a flat apply
// races the Namespace at create and, far worse, deletes the CRDs in the
// same pass as the operator that must stay alive to process the runtime
// InstallerSets' finalizers during the CRD drain (verified live: the
// tektoninstallersets CRD wedges in Terminating for the full timeout).
type manifestPartitions struct {
	// NamespaceYaml carries the fixed `tekton-operator` Namespace —
	// created first, deleted last.
	NamespaceYaml string
	// WorkloadsYaml carries everything that is neither Namespace nor
	// CRD: the two Deployments, RBAC, ConfigMaps, Services, the webhook
	// Secret.
	WorkloadsYaml string
	// CrdsYaml carries the 14 operator.tekton.dev CRDs — created last
	// (the operator crash-retries until they exist), deleted FIRST while
	// the operator still runs.
	CrdsYaml string
}

// fetchManifestPartitions downloads the released manifest and splits it
// into the three ordered groups. Fetching in module code keeps the same
// posture as the Terraform twin's data.http: the manifest is read from
// the pinned, immutable release asset at plan/apply time.
func fetchManifestPartitions(url string) (*manifestPartitions, error) {
	client := &http.Client{Timeout: 2 * time.Minute}
	response, err := client.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching the release manifest from %s", url)
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != http.StatusOK {
		return nil, errors.Errorf("fetching the release manifest from %s: HTTP %d", url, response.StatusCode)
	}
	raw, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading the release manifest body")
	}

	partitions := &manifestPartitions{}
	var namespaceDocs, workloadDocs, crdDocs []string
	for _, document := range documentSeparator.Split(string(raw), -1) {
		if strings.TrimSpace(document) == "" {
			continue
		}
		var meta struct {
			Kind string `yaml:"kind"`
		}
		if err := yaml.Unmarshal([]byte(document), &meta); err != nil {
			return nil, errors.Wrap(err, "reading a manifest document's kind")
		}
		switch meta.Kind {
		case "":
			// Comment-only fragments carry no kind and deploy nothing.
			continue
		case "Namespace":
			namespaceDocs = append(namespaceDocs, document)
		case "CustomResourceDefinition":
			crdDocs = append(crdDocs, document)
		default:
			workloadDocs = append(workloadDocs, document)
		}
	}
	if len(namespaceDocs) == 0 || len(crdDocs) == 0 || len(workloadDocs) == 0 {
		return nil, errors.Errorf(
			"the release manifest did not partition as expected (namespaces=%d, crds=%d, workloads=%d) — the pinned asset's shape changed",
			len(namespaceDocs), len(crdDocs), len(workloadDocs))
	}
	partitions.NamespaceYaml = strings.Join(namespaceDocs, "\n---\n")
	partitions.WorkloadsYaml = strings.Join(workloadDocs, "\n---\n")
	partitions.CrdsYaml = strings.Join(crdDocs, "\n---\n")
	return partitions, nil
}

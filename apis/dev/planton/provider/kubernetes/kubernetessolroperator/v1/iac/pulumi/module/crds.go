package module

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
	pulumiyaml "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"sigs.k8s.io/yaml"
)

// yamlDocumentSeparator matches a YAML document separator LINE ("---" at
// column 0). Splitting on the line (not the bare substring) matters: the
// CRD schemas embed "---" inside description text (`rw-rw----`, a
// "--- TODO" note in the ZookeeperCluster schema) and a substring split
// would corrupt those documents.
var yamlDocumentSeparator = regexp.MustCompile(`(?m)^---[ \t]*$`)

// crdManifestMeta is the minimal shape needed to identify a CRD document.
type crdManifestMeta struct {
	Kind     string `json:"kind"`
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
}

// applyCrds applies the module-owned CRD files staged at vars.CrdsDir —
// the solr-operator chart ships NO CRDs (they are separate release
// artifacts), so the module owns the three solr.apache.org CRDs plus the
// ZookeeperCluster CRD of the bundled zookeeper-operator dependency.
// One ConfigGroup per CRD, keyed by the CRD's OWN metadata.name, so
// resource addresses stay stable regardless of file naming or ordering.
//
// KEEP-ON-UNINSTALL: every applied CRD carries retainOnDelete, so
// destroying this stack removes the operator but NEVER the CRDs — and
// therefore never cascade-deletes SolrCloud/SolrBackup/ZookeeperCluster
// resources cluster-wide. This is the exact twin of the Terraform
// module's `apply_only = true` on kubectl_manifest (whose Delete is a
// no-op in the provider source).
//
// The retainOnDelete option must reach the ConfigGroup's CHILDREN (the
// actual CRD resources): the yaml SDK drops ordinary resource options
// when creating children (GetChildOptions forwards only
// version/pluginDownloadURL), so it is delivered through a resource
// TRANSFORMATION, which the SDK propagates from parent to children.
func applyCrds(ctx *pulumi.Context, kubernetesProvider pulumi.ProviderResource) ([]pulumi.Resource, error) {
	crdDocuments, err := loadCrdDocuments(vars.CrdsDir)
	if err != nil {
		return nil, err
	}

	retainOnDelete := pulumi.Transformations([]pulumi.ResourceTransformation{
		func(args *pulumi.ResourceTransformationArgs) *pulumi.ResourceTransformationResult {
			return &pulumi.ResourceTransformationResult{
				Props: args.Props,
				Opts:  append(args.Opts, pulumi.RetainOnDelete(true)),
			}
		},
	})

	// Deterministic creation order (map iteration is randomized).
	crdNames := make([]string, 0, len(crdDocuments))
	for name := range crdDocuments {
		crdNames = append(crdNames, name)
	}
	sort.Strings(crdNames)

	createdCrds := make([]pulumi.Resource, 0, len(crdNames))
	for _, crdName := range crdNames {
		createdCrd, err := pulumiyaml.NewConfigGroup(ctx, crdName,
			&pulumiyaml.ConfigGroupArgs{
				YAML: []string{crdDocuments[crdName]},
			},
			pulumi.Provider(kubernetesProvider),
			retainOnDelete,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to apply crd %s", crdName)
		}
		createdCrds = append(createdCrds, createdCrd)
	}

	return createdCrds, nil
}

// loadCrdDocuments reads every staged *.yaml file, splits it into YAML
// documents, and returns the CRD documents keyed by metadata.name.
// License-comment headers and empty documents (the ZookeeperCluster file
// opens with a comment block and a doubled "---") are skipped by the
// parse-and-check filter.
func loadCrdDocuments(dir string) (map[string]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read crds directory %s", dir)
	}

	crdDocuments := map[string]string{}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		content, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read crd file %s", entry.Name())
		}
		for _, document := range yamlDocumentSeparator.Split(string(content), -1) {
			if strings.TrimSpace(document) == "" {
				continue
			}
			meta := crdManifestMeta{}
			// Comment-only documents unmarshal to the zero value; anything
			// without a metadata.name is not a manifest to apply.
			if err := yaml.Unmarshal([]byte(document), &meta); err != nil || meta.Metadata.Name == "" {
				continue
			}
			crdDocuments[meta.Metadata.Name] = document
		}
	}

	if len(crdDocuments) == 0 {
		return nil, errors.Errorf("no CRD documents found in %s", dir)
	}
	return crdDocuments, nil
}

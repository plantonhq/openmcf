package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const twoDocumentStream = `apiVersion: gcp.planton.ai/v1alpha1
kind: CloudRunService
metadata:
  name: storefront
---
apiVersion: gcp.planton.ai/v1alpha1
kind: GcpMemorystore
metadata:
  name: storefront-cache
`

func TestRefuseMultiDocument_NamesEveryDocument(t *testing.T) {
	err := refuseMultiDocument([]byte(twoDocumentStream), "overlay-render.yaml")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "2 YAML documents")
	assert.Contains(t, err.Error(), "CloudRunService/storefront")
	assert.Contains(t, err.Error(), "GcpMemorystore/storefront-cache")
	assert.Contains(t, err.Error(), "overlay-render.yaml")
}

func TestRefuseMultiDocument_SingleDocumentPasses(t *testing.T) {
	single := "kind: CloudRunService\nmetadata:\n  name: storefront\n"
	assert.NoError(t, refuseMultiDocument([]byte(single), "manifest.yaml"))
}

func TestRefuseMultiDocument_SeparatorNoiseDoesNotCount(t *testing.T) {
	cases := map[string]string{
		"leading separator":  "---\nkind: CloudRunService\nmetadata:\n  name: storefront\n",
		"trailing separator": "kind: CloudRunService\nmetadata:\n  name: storefront\n---\n",
		"comment-only tail":  "kind: CloudRunService\nmetadata:\n  name: storefront\n---\n# just a comment\n",
	}
	for name, doc := range cases {
		assert.NoError(t, refuseMultiDocument([]byte(doc), "manifest.yaml"), name)
	}
}

func TestRefuseMultiDocument_MalformedYamlIsLeftToTheLoader(t *testing.T) {
	// The schema-aware loader owns malformed-YAML diagnosis; this gate must
	// not intercept it with a less helpful message.
	assert.NoError(t, refuseMultiDocument([]byte("kind: [unclosed"), "manifest.yaml"))
}

func TestLoadManifestBytes_RefusesMultiDocument(t *testing.T) {
	_, err := LoadManifestBytes([]byte(twoDocumentStream), "overlay-render.yaml")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exactly one manifest per run")
}

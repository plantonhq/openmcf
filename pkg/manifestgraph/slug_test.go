package manifestgraph

import (
	"testing"

	"github.com/plantonhq/planton/shared"
	"github.com/stretchr/testify/assert"
)

// TestGenerateSlug pins the port to the platform's slug generation rules —
// node identity on both lanes derives through these exact rules, so a
// divergence here is an identity divergence everywhere.
func TestGenerateSlug(t *testing.T) {
	cases := map[string]string{
		"storefront":            "storefront",
		"My Shared Producer":    "my-shared-producer",
		"My VPC (US-East)":      "my-vpc-us-east",
		"a__b":                  "a__b",
		"CAFÉ deploy":           "cafe-deploy",
		"  spaced  out  ":       "spaced-out",
		"trailing-hyphen-":      "trailing-hyphen",
		"--leading":             "leading",
		"dots.and.dashes-mixed": "dots-and-dashes-mixed",
		"":                      "",
		"   ":                   "",
	}
	for name, want := range cases {
		assert.Equal(t, want, GenerateSlug(name), "name %q", name)
	}
}

func TestResolveSlug_ExplicitSlugPassesThrough(t *testing.T) {
	meta := &shared.CloudResourceMetadata{Name: "My Shared Producer", Slug: "authored-slug"}
	assert.Equal(t, "authored-slug", ResolveSlug(meta))

	meta = &shared.CloudResourceMetadata{Name: "My Shared Producer"}
	assert.Equal(t, "my-shared-producer", ResolveSlug(meta))

	assert.Equal(t, "", ResolveSlug(nil))
}

package manifestgraph

import (
	"regexp"
	"strings"

	"github.com/plantonhq/planton/shared"
	"golang.org/x/text/unicode/norm"
)

// The slug grammar below is byte-for-byte the platform's resource-slug
// generation (lowercase, hyphen-separated, word characters only). Node and
// edge-target identity both derive through it, so a manifest named
// "My Shared VPC" and a reference naming "My Shared VPC" land on the same
// slug — the join that makes them one graph node.
var (
	nonWord     = regexp.MustCompile(`[^\w-]`)
	whitespace  = regexp.MustCompile(`\s+`)
	multiHyphen = regexp.MustCompile(`-{2,}`)
)

// GenerateSlug converts a human-readable resource name into its identity
// slug: unicode-decompose, replace non-word characters and whitespace with
// hyphens, lowercase, collapse hyphen runs, trim edge hyphens. An empty or
// blank name yields "".
func GenerateSlug(name string) string {
	if strings.TrimSpace(name) == "" {
		return ""
	}
	slug := norm.NFD.String(name)
	slug = nonWord.ReplaceAllString(slug, "-")
	slug = whitespace.ReplaceAllString(slug, "-")
	slug = strings.ToLower(slug)
	slug = multiHyphen.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

// ResolveSlug resolves a manifest's identity slug: an explicit metadata.slug
// passes through untouched (it IS the platform identity when present),
// otherwise the slug generates from metadata.name.
func ResolveSlug(meta *shared.CloudResourceMetadata) string {
	if meta == nil {
		return ""
	}
	if meta.GetSlug() != "" {
		return meta.GetSlug()
	}
	return GenerateSlug(meta.GetName())
}

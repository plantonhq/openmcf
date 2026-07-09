//go:build !codegen
// +build !codegen

package refcheck

import (
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
)

// ResolveRefPath reports whether a reference field path (e.g.
// "status.outputs.bucket_id", "spec.vpc_cidr", "metadata.name") resolves
// against the given kind. It returns an empty string when the path resolves,
// or a human-readable reason when it does not — the same contract the
// annotation analyzer uses internally, exported for consumers that validate
// concrete valueFrom references (such as infra-chart validation) against the
// same single source of truth.
func ResolveRefPath(kind cloudresourcekind.CloudResourceKind, refPath string) string {
	rootMd, rest, reason := targetRoot(kind, refPath)
	if reason != "" {
		return reason
	}
	return resolvePath(rootMd, rest)
}

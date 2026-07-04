// Package crkreflect provides runtime access to CloudResourceKind metadata
// encoded as protobuf enum value options.
package crkreflect

import (
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
)

// Prerequisites returns the direct prerequisites for the given kind.
// Returns nil if the kind has no prerequisites or if kind meta is unavailable.
func Prerequisites(kind cloudresourcekind.CloudResourceKind) []cloudresourcekind.CloudResourceKind {
	meta, err := KindMeta(kind)
	if err != nil {
		return nil
	}
	return meta.GetPrerequisites()
}

// TransitivePrerequisites returns all prerequisites in topological order
// (deploy first to last). Resolves the full transitive dependency graph:
// if A depends on B and B depends on C, returns [C, B].
//
// Returns an error if a cycle is detected (indicates a modeling mistake).
func TransitivePrerequisites(kind cloudresourcekind.CloudResourceKind) ([]cloudresourcekind.CloudResourceKind, error) {
	return TransitiveClosure(Prerequisites(kind))
}

// TransitiveClosure returns the given kinds plus all of their transitive
// prerequisites in topological order (deploy first to last), deduplicated.
// Every kind appears after everything it depends on, so callers can deploy
// the result front to back. Accepts multiple roots so a caller can expand a
// combined set (e.g. a kind's registry prerequisites merged with extras a
// test scenario composes) in one pass.
//
// Returns an error if a cycle is detected (indicates a modeling mistake).
func TransitiveClosure(kinds []cloudresourcekind.CloudResourceKind) ([]cloudresourcekind.CloudResourceKind, error) {
	var result []cloudresourcekind.CloudResourceKind
	visited := make(map[cloudresourcekind.CloudResourceKind]bool)
	inStack := make(map[cloudresourcekind.CloudResourceKind]bool)

	var root cloudresourcekind.CloudResourceKind
	var visit func(k cloudresourcekind.CloudResourceKind) error
	visit = func(k cloudresourcekind.CloudResourceKind) error {
		if inStack[k] {
			return cycleError(root, k)
		}
		if visited[k] {
			return nil
		}

		inStack[k] = true
		for _, prereq := range Prerequisites(k) {
			if err := visit(prereq); err != nil {
				return err
			}
		}
		inStack[k] = false
		visited[k] = true
		result = append(result, k)
		return nil
	}

	for _, k := range kinds {
		root = k
		if err := visit(k); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// HasPrerequisites returns true if the kind has any direct prerequisites.
func HasPrerequisites(kind cloudresourcekind.CloudResourceKind) bool {
	return len(Prerequisites(kind)) > 0
}

func cycleError(root, cycleAt cloudresourcekind.CloudResourceKind) error {
	return &CycleError{Root: root, CycleAt: cycleAt}
}

// CycleError indicates a circular dependency in the prerequisite graph.
type CycleError struct {
	Root    cloudresourcekind.CloudResourceKind
	CycleAt cloudresourcekind.CloudResourceKind
}

func (e *CycleError) Error() string {
	return "prerequisite cycle detected: " + e.Root.String() + " -> ... -> " + e.CycleAt.String()
}

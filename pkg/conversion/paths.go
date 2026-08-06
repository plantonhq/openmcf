package conversion

import "strings"

// Document paths are dot-separated JSON field names ("spec.displayName").
// Path segments never index into lists: an operation that reshapes list
// items does so with a CEL expression over the whole list value.

func splitPath(path string) []string {
	return strings.Split(path, ".")
}

// getPath returns the node at path and whether it exists.
func getPath(doc map[string]any, path string) (any, bool) {
	segments := splitPath(path)
	var current any = doc
	for _, seg := range segments {
		m, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}
		current, ok = m[seg]
		if !ok {
			return nil, false
		}
	}
	return current, true
}

// setPath writes value at path, creating intermediate maps as needed.
func setPath(doc map[string]any, path string, value any) {
	segments := splitPath(path)
	current := doc
	for _, seg := range segments[:len(segments)-1] {
		next, ok := current[seg].(map[string]any)
		if !ok {
			next = map[string]any{}
			current[seg] = next
		}
		current = next
	}
	current[segments[len(segments)-1]] = value
}

// deletePath removes the node at path; intermediate maps left empty are
// pruned so a fully-converted document has no hollow husks.
func deletePath(doc map[string]any, path string) {
	segments := splitPath(path)
	parents := make([]map[string]any, 0, len(segments))
	current := doc
	for _, seg := range segments[:len(segments)-1] {
		parents = append(parents, current)
		next, ok := current[seg].(map[string]any)
		if !ok {
			return
		}
		current = next
	}
	delete(current, segments[len(segments)-1])
	// Prune now-empty intermediates bottom-up.
	for i := len(parents) - 1; i >= 0; i-- {
		if len(current) == 0 {
			delete(parents[i], segments[i])
		}
		current = parents[i]
	}
}

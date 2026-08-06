package conversion

import (
	"fmt"
	"reflect"

	"github.com/google/cel-go/cel"
	"google.golang.org/protobuf/types/known/structpb"
)

// Direction selects which op list of a Spec to execute.
type Direction int

const (
	// Upgrade converts a document from Spec.From to Spec.To.
	Upgrade Direction = iota
	// Downgrade converts a document from Spec.To back to Spec.From.
	Downgrade
)

// DeclaredLoss records one value an executed drop op actually removed.
// Losses are part of the result, never a side channel: callers surface them
// to users (CLI output, UX warnings, migration reports).
type DeclaredLoss struct {
	Path   string
	Reason string
}

// Apply converts doc between the spec's versions in the given direction.
// The document is modified structurally (JSON field names) and its
// apiVersion is restamped to the destination version. The input map is not
// mutated; the converted document is returned.
func Apply(spec *Spec, direction Direction, doc map[string]any) (map[string]any, []DeclaredLoss, error) {
	ops := spec.Ops
	targetVersion := spec.To
	if direction == Downgrade {
		ops = spec.Reverse
		targetVersion = spec.From
	}

	converted := deepCopy(doc).(map[string]any)
	var losses []DeclaredLoss

	for i, op := range ops {
		switch {
		case op.Rename != nil:
			if v, ok := getPath(converted, op.Rename.FromPath); ok {
				deletePath(converted, op.Rename.FromPath)
				setPath(converted, op.Rename.ToPath, v)
			}

		case op.Map != nil:
			v, ok := getPath(converted, op.Map.Path)
			if !ok {
				continue
			}
			mapped, err := evalCel(op.Map.Expr, v)
			if err != nil {
				return nil, nil, fmt.Errorf("op %d (map %s): %w", i, op.Map.Path, err)
			}
			to := op.Map.To
			if to == "" {
				to = op.Map.Path
			}
			deletePath(converted, op.Map.Path)
			setPath(converted, to, mapped)

		case op.Default != nil:
			if _, ok := getPath(converted, op.Default.Path); !ok {
				setPath(converted, op.Default.Path, op.Default.Value)
			}

		case op.Drop != nil:
			if _, ok := getPath(converted, op.Drop.Path); ok {
				deletePath(converted, op.Drop.Path)
				losses = append(losses, DeclaredLoss{Path: op.Drop.Path, Reason: op.Drop.Lossy.Reason})
			}
		}
	}

	// Restamp the envelope: the document now speaks the destination version.
	if apiVersion, ok := getPath(converted, "apiVersion"); ok {
		if s, isString := apiVersion.(string); isString {
			if group, _, found := cutLast(s, "/"); found {
				setPath(converted, "apiVersion", group+"/"+targetVersion)
			}
		}
	}

	return converted, losses, nil
}

func cutLast(s, sep string) (before, after string, found bool) {
	idx := -1
	for i := len(s) - len(sep); i >= 0; i-- {
		if s[i:i+len(sep)] == sep {
			idx = i
			break
		}
	}
	if idx < 0 {
		return s, "", false
	}
	return s[:idx], s[idx+len(sep):], true
}

// evalCel evaluates a CEL expression with `value` bound to the source node.
func evalCel(expr string, value any) (any, error) {
	env, err := cel.NewEnv(cel.Variable("value", cel.DynType))
	if err != nil {
		return nil, fmt.Errorf("creating CEL environment: %w", err)
	}
	ast, issues := env.Compile(expr)
	if issues != nil && issues.Err() != nil {
		return nil, fmt.Errorf("the CEL expression %q does not compile: %w", expr, issues.Err())
	}
	program, err := env.Program(ast)
	if err != nil {
		return nil, fmt.Errorf("building CEL program for %q: %w", expr, err)
	}
	out, _, err := program.Eval(map[string]any{"value": value})
	if err != nil {
		return nil, fmt.Errorf("evaluating CEL expression %q: %w", expr, err)
	}
	// Conversion through structpb guarantees a JSON-shaped native value
	// (map[string]any / []any / string / float64 / bool / nil), so converted
	// documents serialize exactly like loaded ones.
	structValue, err := out.ConvertToNative(reflect.TypeOf(&structpb.Value{}))
	if err != nil {
		return nil, fmt.Errorf("converting CEL result of %q to a document value: %w", expr, err)
	}
	return structValue.(*structpb.Value).AsInterface(), nil
}

// deepCopy clones a JSON-shaped document tree.
func deepCopy(v any) any {
	switch t := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, val := range t {
			out[k] = deepCopy(val)
		}
		return out
	case []any:
		out := make([]any, len(t))
		for i, val := range t {
			out[i] = deepCopy(val)
		}
		return out
	default:
		return v
	}
}

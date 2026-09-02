package tekton

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"sigs.k8s.io/yaml"
)

// Images returns every container image the embedded content can run, sorted
// and deduplicated, exactly as written -- tag documents intent, digest is
// what the cluster pulls. This is the mirror/allowlist an air-gapped or
// egress-restricted build cluster needs, derived from the content instead of
// hand-copied beside it, so the two cannot drift.
//
// An image is either a literal on a step, or a step's "$(params.X)" resolved
// through the param that feeds it: the task's own default for X, or -- when
// a pipeline references the task -- the value the pipeline task passes for X
// (a literal, or "$(params.Y)" resolved through the pipeline's default for
// Y). A reference that resolves to nothing (the image being BUILT, supplied
// at dispatch) is not a pulled image and is skipped. Inline taskSpecs inside
// pipelines resolve the same way.
func Images() ([]string, error) {
	tasks, err := loadTaskDocs()
	if err != nil {
		return nil, err
	}

	set := map[string]struct{}{}
	for _, task := range tasks {
		for _, image := range task.resolveImages(nil) {
			set[image] = struct{}{}
		}
	}

	for _, track := range Tracks() {
		raw, _ := Track(track)
		var doc map[string]any
		if err := yaml.Unmarshal(raw, &doc); err != nil {
			return nil, fmt.Errorf("parsing pipeline %s: %w", track, err)
		}
		spec, _ := doc["spec"].(map[string]any)
		pipelineDefaults := paramDefaults(spec)
		for _, listKey := range []string{"tasks", "finally"} {
			entries, _ := spec[listKey].([]any)
			for _, entry := range entries {
				pipelineTask, _ := entry.(map[string]any)
				if pipelineTask == nil {
					continue
				}
				// The values a pipeline task passes down, with the pipeline's
				// own params already substituted.
				passed := passedParams(pipelineTask, pipelineDefaults)

				if inline, ok := pipelineTask["taskSpec"].(map[string]any); ok {
					for _, image := range (&taskDoc{spec: inline}).resolveImages(passed) {
						set[image] = struct{}{}
					}
					continue
				}
				ref, _ := pipelineTask["taskRef"].(map[string]any)
				refName, _ := ref["name"].(string)
				if task, ok := tasks[refName]; ok {
					for _, image := range task.resolveImages(passed) {
						set[image] = struct{}{}
					}
				}
			}
		}
	}

	images := make([]string, 0, len(set))
	for image := range set {
		images = append(images, image)
	}
	sort.Strings(images)
	return images, nil
}

// taskDoc is one parsed Task (or inline taskSpec): its params with defaults
// and its steps' image fields.
type taskDoc struct {
	spec map[string]any
}

// loadTaskDocs parses every embedded task, keyed by its Tekton metadata.name
// (the name a pipeline's taskRef uses).
func loadTaskDocs() (map[string]*taskDoc, error) {
	files, err := TaskFiles()
	if err != nil {
		return nil, err
	}
	out := make(map[string]*taskDoc, len(files))
	for stem, raw := range files {
		var doc map[string]any
		if err := yaml.Unmarshal(raw, &doc); err != nil {
			return nil, fmt.Errorf("parsing task %s: %w", stem, err)
		}
		metadata, _ := doc["metadata"].(map[string]any)
		name, _ := metadata["name"].(string)
		if name == "" {
			return nil, fmt.Errorf("task %s declares no metadata.name", stem)
		}
		spec, _ := doc["spec"].(map[string]any)
		out[name] = &taskDoc{spec: spec}
	}
	return out, nil
}

// resolveImages returns the literal images this task's steps run, given the
// params a caller passes (which win over the task's own defaults).
func (t *taskDoc) resolveImages(passed map[string]string) []string {
	defaults := paramDefaults(t.spec)
	var images []string
	steps, _ := t.spec["steps"].([]any)
	for _, entry := range steps {
		step, _ := entry.(map[string]any)
		image, _ := step["image"].(string)
		if image == "" {
			continue
		}
		if param, isRef := paramRef(image); isRef {
			if value, ok := passed[param]; ok {
				image = value
			} else if value, ok := defaults[param]; ok {
				image = value
			} else {
				continue // supplied at dispatch: not an image this content pulls
			}
		}
		images = append(images, image)
	}
	return images
}

// paramDefaults collects the string defaults a spec's params declare.
func paramDefaults(spec map[string]any) map[string]string {
	out := map[string]string{}
	params, _ := spec["params"].([]any)
	for _, entry := range params {
		param, _ := entry.(map[string]any)
		name, _ := param["name"].(string)
		value, ok := param["default"].(string)
		if name != "" && ok {
			out[name] = value
		}
	}
	return out
}

// passedParams collects the string values a pipeline task passes to its
// task, substituting the pipeline's own param defaults for "$(params.Y)"
// references. A value the pipeline cannot resolve is dropped, so the task's
// own default (or nothing) applies.
func passedParams(pipelineTask map[string]any, pipelineDefaults map[string]string) map[string]string {
	out := map[string]string{}
	params, _ := pipelineTask["params"].([]any)
	for _, entry := range params {
		param, _ := entry.(map[string]any)
		name, _ := param["name"].(string)
		value, ok := param["value"].(string)
		if name == "" || !ok {
			continue
		}
		if ref, isRef := paramRef(value); isRef {
			resolved, found := pipelineDefaults[ref]
			if !found {
				continue
			}
			value = resolved
		}
		out[name] = value
	}
	return out
}

// paramRefPattern matches a whole value that is exactly one Tekton param
// reference, "$(params.<name>)".
var paramRefPattern = regexp.MustCompile(`^\$\(params\.([A-Za-z0-9_.-]+)\)$`)

func paramRef(value string) (name string, ok bool) {
	m := paramRefPattern.FindStringSubmatch(strings.TrimSpace(value))
	if m == nil {
		return "", false
	}
	return m[1], true
}

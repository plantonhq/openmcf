package explain

import (
	"strings"
	"testing"
)

// markdownHeadingVocabulary is the published search grammar (see the
// RenderMarkdown doc comment). Tests pin both membership and order: agents
// grep these headings across every generated file, so a rename or reorder is
// a breaking change to every consumer at once.
var markdownHeadingVocabulary = []string{
	"## Example",
	"## Spec Fields",
	"## Field Details",
	"## Validation Rules",
	"## Outputs",
	"## References",
	"## See Also",
}

func renderKindMarkdown(t *testing.T, kind string, opts MarkdownOptions) string {
	t.Helper()
	engine := DefaultEngine()
	report, err := engine.Explain(mustResource(t, kind), nil)
	if err != nil {
		t.Fatal(err)
	}
	return RenderMarkdown(report, opts)
}

func TestRenderMarkdownRootView(t *testing.T) {
	text := renderKindMarkdown(t, "AwsVpc", MarkdownOptions{
		ExampleYAML: "apiVersion: aws.planton.dev/v1\nkind: AwsVpc\n",
		SeeAlso:     []MarkdownLink{{Title: "Overview", Path: "./README.md"}},
	})

	for _, want := range []string{
		"# AwsVpc\n",
		"do not\n> edit by hand", // the generated-file banner
		"**apiVersion**: `aws.planton.dev/v1`",
		"## Example", "```yaml", "kind: AwsVpc",
		"## Spec Fields",
		"| `spec.region` |", // table row, dotted path
		"## Field Details",
		"### spec.region",
		"## Outputs",
		"| `status.outputs.vpcId` |",
		"## See Also",
		"[Overview](./README.md)",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("markdown missing %q\n---\n%s", want, text)
		}
	}

	// region derives required from min_len; the table row must say so.
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "| `spec.region` |") && !strings.Contains(line, "| yes |") {
			t.Errorf("spec.region row should be marked required: %q", line)
		}
	}
}

func TestRenderMarkdownHeadingGrammar(t *testing.T) {
	text := renderKindMarkdown(t, "AwsRdsCluster", MarkdownOptions{ExampleYAML: "kind: AwsRdsCluster\n"})

	var seen []string
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "## ") {
			seen = append(seen, line)
		}
	}
	if len(seen) == 0 {
		t.Fatal("no H2 headings rendered")
	}
	// Every rendered heading must come from the vocabulary, in canonical order.
	cursor := 0
	for _, heading := range seen {
		found := false
		for cursor < len(markdownHeadingVocabulary) {
			if markdownHeadingVocabulary[cursor] == heading {
				found = true
				cursor++
				break
			}
			cursor++
		}
		if !found {
			t.Errorf("heading %q not in the vocabulary (or out of canonical order); rendered: %v", heading, seen)
		}
	}
}

func TestRenderMarkdownNestedListPaths(t *testing.T) {
	text := renderKindMarkdown(t, "AwsRdsCluster", MarkdownOptions{})
	// Children of a repeated message keep the `[]` marker so the path spells
	// the YAML nesting an author writes.
	if !strings.Contains(text, "### spec.instances[].instanceClass") {
		t.Errorf("nested list path missing:\n%s", text)
	}
}

func TestRenderMarkdownForeignKeys(t *testing.T) {
	text := renderKindMarkdown(t, "aws-security-group", MarkdownOptions{})

	// The pipe in the shared type label must not split the table row.
	if !strings.Contains(text, "`string \\| valueFrom`") {
		t.Errorf("foreign-key type label not pipe-escaped in table:\n%s", text)
	}
	if !strings.Contains(text, "## References") {
		t.Fatal("References section missing")
	}
	if !strings.Contains(text, "| `spec.vpcId` | AwsVpc | `status.outputs.vpc_id` |") {
		t.Errorf("outbound reference row missing:\n%s", text)
	}
	// The valueFrom authoring contract must appear at the field detail --
	// the only place an author learns a bare string does not parse.
	if !strings.Contains(text, "write as {value:") || !strings.Contains(text, "kind: AwsVpc") {
		t.Errorf("valueFrom contract missing from field details:\n%s", text)
	}
}

func TestRenderMarkdownEnumsFullyExpanded(t *testing.T) {
	text := renderKindMarkdown(t, "kubernetes-namespace", MarkdownOptions{})
	if !strings.Contains(text, "Allowed values (use exactly as shown):") {
		t.Fatalf("enum expansion missing:\n%s", text)
	}
	if !strings.Contains(text, "- `baseline` -- ") || !strings.Contains(text, "Minimally restrictive") {
		t.Errorf("enum value docs must render in full -- files have no terminal collapse:\n%s", text)
	}
}

func TestRenderMarkdownDeterminism(t *testing.T) {
	engine := DefaultEngine()
	report, err := engine.Explain(mustResource(t, "AwsRdsCluster"), nil)
	if err != nil {
		t.Fatal(err)
	}
	opts := MarkdownOptions{ExampleYAML: "kind: AwsRdsCluster\n"}
	if RenderMarkdown(report, opts) != RenderMarkdown(report, opts) {
		t.Fatal("render is not deterministic for identical input")
	}
}

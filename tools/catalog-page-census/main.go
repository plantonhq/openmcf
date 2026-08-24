// Command catalog-page-census reports on the validity of every COMPLETE
// manifest (a fenced yaml document declaring apiVersion + kind) embedded in
// the catalog pages (catalog/<provider>/<kind>/catalog.md), using the same
// load+validate machinery the GUIDE wisdom guardrails use. A catalog page
// whose copy-pasted manifest fails `planton validate` teaches every reader a
// broken shape — this census makes that class measurable and its burn-down
// verifiable. It is a REPORT, not a gate: run it any time
// (`go run ./tools/catalog-page-census/`), optionally with an argument to
// print full validation errors for pages whose path contains the argument.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/plantonhq/planton/internal/manifest"
)

var (
	yamlBlock  = regexp.MustCompile("(?s)```yaml\n(.*?)```")
	apiVersion = regexp.MustCompile(`(?m)^apiVersion:`)
	kindLine   = regexp.MustCompile(`(?m)^kind:`)
)

func main() {
	pages, err := filepath.Glob("catalog/*/*/catalog.md")
	if err != nil {
		panic(err)
	}
	sort.Strings(pages)

	var pagesTotal, pagesWithManifests, manifestsTotal, manifestsInvalid, pagesFailing int
	failuresByProvider := map[string]int{}

	for _, page := range pages {
		if strings.Contains(page, "/_test/") {
			continue
		}
		pagesTotal++
		raw, err := os.ReadFile(page)
		if err != nil {
			panic(err)
		}
		var docs []string
		for _, block := range yamlBlock.FindAllStringSubmatch(string(raw), -1) {
			for _, doc := range strings.Split(block[1], "\n---\n") {
				if apiVersion.MatchString(doc) && kindLine.MatchString(doc) {
					docs = append(docs, doc)
				}
			}
		}
		if len(docs) == 0 {
			continue
		}
		pagesWithManifests++
		pageFailed := false
		for i, doc := range docs {
			manifestsTotal++
			loaded, err := manifest.LoadManifestBytes([]byte(doc), page)
			if err == nil {
				err = manifest.ValidateLoaded(loaded)
			}
			if err != nil {
				manifestsInvalid++
				pageFailed = true
				if len(os.Args) > 1 && strings.Contains(page, os.Args[1]) {
					fmt.Printf("INVALID %s manifest#%d FULL ERROR:\n%v\n---\n", page, i+1, err)
				} else {
					fmt.Printf("INVALID %s manifest#%d: %v\n", page, i+1, firstLine(err.Error()))
				}
			}
		}
		if pageFailed {
			pagesFailing++
			failuresByProvider[strings.Split(page, "/")[1]]++
		}
	}

	fmt.Printf("\n=== catalog.md manifest-validity census ===\n")
	fmt.Printf("pages scanned: %d\n", pagesTotal)
	fmt.Printf("pages with complete manifests: %d\n", pagesWithManifests)
	fmt.Printf("manifests validated: %d\n", manifestsTotal)
	fmt.Printf("manifests INVALID: %d\n", manifestsInvalid)
	fmt.Printf("pages with at least one invalid manifest: %d\n", pagesFailing)
	providers := make([]string, 0, len(failuresByProvider))
	for p := range failuresByProvider {
		providers = append(providers, p)
	}
	sort.Strings(providers)
	for _, p := range providers {
		fmt.Printf("  %s: %d failing pages\n", p, failuresByProvider[p])
	}
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

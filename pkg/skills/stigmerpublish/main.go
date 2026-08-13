// Command stigmerpublish publishes packaged skill definitions to a Stigmer
// engine. It consumes a defspack output directory (the manifest plus the
// deterministic per-skill zips) and pushes those exact bytes, so the
// engine's content-addressed version identity equals the release manifest's
// checksum -- one content state, one checksum, on the CDN, in the engine,
// and in every consumer that verifies either.
//
// Pushing the packager's bytes (instead of re-zipping the skill directory
// here or in a CLI) is the point of this tool's existence: the engine keys
// a skill version on the SHA-256 of the uploaded archive, so only a
// deterministic archive makes "unchanged content is a no-op" true from CI
// and keeps the hosted engine's version identity equal to the release
// manifest. The Stigmer CLI re-zips per push today, which registers a new
// version on every CI run for unchanged content.
//
// SELF-RETIRING: this tool exists only until the official Stigmer CLI can
// push a pre-packaged archive (requested upstream:
// https://github.com/stigmer/stigmer/issues/671). When that ships, the
// publish workflows swap to the CLI and this tool is deleted.
//
// Run `go run ./pkg/skills/stigmerpublish -dir build/definitions` after
// `go run ./pkg/skills/defspack -version <tag> -out build/definitions`:
//
//	-dir          defspack output directory (manifest + skill zips)
//	-org          target engine org (default "planton")
//	-skill        limit to one skill slug (default: every skill in the manifest)
//	-verify-only  compare the engine's latest version hash against the
//	              manifest without pushing; non-zero exit on any mismatch
//	-server       override the engine address (tests and local engines)
//	-insecure     dial without TLS (loopback engines only)
//
// Auth is STIGMER_API_KEY from the environment, never a flag -- keys must
// not land in shell history or CI logs.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	stigmer "github.com/stigmer/stigmer/sdk/go/v3"
)

func main() {
	dir := flag.String("dir", "build/definitions", "defspack output directory (manifest + skill zips)")
	org := flag.String("org", "planton", "target engine org slug")
	skill := flag.String("skill", "", "limit to one skill slug (default: all skills in the manifest)")
	verifyOnly := flag.Bool("verify-only", false, "compare engine state against the manifest without pushing")
	server := flag.String("server", "", "override the engine address (default: Stigmer Cloud)")
	insecure := flag.Bool("insecure", false, "dial without TLS (loopback engines only)")
	flag.Parse()

	opts := []stigmer.ClientOption{}
	if key := os.Getenv("STIGMER_API_KEY"); key != "" {
		opts = append(opts, stigmer.WithAPIKey(key))
	} else if !*insecure {
		fatal(fmt.Errorf("STIGMER_API_KEY is not set"))
	}
	if *server != "" {
		opts = append(opts, stigmer.WithBaseURL(*server))
	}
	if *insecure {
		opts = append(opts, stigmer.WithInsecure())
	}
	client, err := stigmer.NewClient(opts...)
	if err != nil {
		fatal(err)
	}
	defer func() { _ = client.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if err := run(ctx, client.Skill, *dir, *org, *skill, *verifyOnly, os.Stdout); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "stigmerpublish:", err)
	os.Exit(1)
}

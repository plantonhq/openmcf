# Fix assetPrefix Static Export 404s for Unified Domain

**Date**: March 28, 2026
**Type**: Bug Fix
**Components**: Build System

## Summary

Added a post-build step to copy `out/_next` to `out/_site/_next` so that Next.js static export assets resolve correctly under the `assetPrefix: '/_site'` configuration. Without this, every CSS and JS file on planton.ai returned 404 from GitHub Pages, rendering the marketing site completely unstyled.

## Problem Statement

The planton.ai unified domain migration introduced `assetPrefix: '/_site'` in `next.config.ts` to prevent `_next` static asset collisions between the marketing site (GitHub Pages) and the console app (Kubernetes). Both apps serve assets from `/_next/...` by default. A Cloudflare Origin Rule routes `/_next/...` to K8s, so the website needs its assets under a different prefix (`/_site/_next/...`) that stays on the GitHub Pages origin.

### Pain Points

- `assetPrefix` changes the URL references in generated HTML (e.g. `<script src="/_site/_next/static/chunks/abc.js">`)
- But Next.js `output: 'export'` writes physical files to `out/_next/` regardless of assetPrefix — the output directory structure is unaffected
- GitHub Pages serves from the `out/` directory, so `/_site/_next/...` requests found no matching files
- Every CSS and JS asset returned 404, leaving the entire marketing site unstyled and non-functional

## Solution

A single post-build command that copies the `_next` directory into the `_site/_next` path that HTML expects:

```bash
mkdir -p out/_site && cp -r out/_next out/_site/_next
```

This was added to both the GitHub Actions workflow (production) and the Makefile (local preview).

### Why `cp -r` and not `mv`

The original `out/_next` is kept intact. It costs trivial disk space and ensures nothing breaks if any internal Next.js reference still uses the non-prefixed path.

### Why not a Cloudflare rewrite rule

A Cloudflare Transform Rule could rewrite `/_site/_next/*` to `/_next/*` at the edge, but that adds infrastructure to manage and reason about. The build-time copy is deterministic, has zero runtime cost, and keeps the complexity in the build pipeline where it belongs.

## Implementation Details

### GitHub Actions workflow

`.github/workflows/pages.yml` — added after `yarn build` and before the 404.html fallback:

```yaml
mkdir -p out/_site
cp -r out/_next out/_site/_next
```

### Makefile

`Makefile` — updated the `build` target so `make preview-site` also works correctly:

```makefile
build: install lint
	NODE_NO_WARNINGS=1 NEXT_TELEMETRY_DISABLED=1 $(YARN) build
	mkdir -p out/_site && cp -r out/_next out/_site/_next
```

## Benefits

- Marketing site fully styled and functional again on planton.ai
- Local `make preview-site` matches production behavior
- No additional infrastructure or Cloudflare rules needed
- Build-time fix — zero runtime cost, deterministic

## Impact

- **Users**: planton.ai marketing site renders correctly with all CSS and JS loaded
- **Developers**: `make preview-site` and `make build` produce the correct output for the assetPrefix configuration
- **Infrastructure**: No changes — fix is entirely in the build pipeline

## Related Work

- Unified domain migration project: `planton/_projects/20260314.01.unified-domain-migration/`
- `assetPrefix: '/_site'` was added in commit `c687eaf` as part of T03 (Website Asset Prefix)
- Cloudflare Origin Rule routes `/_site/*` to GitHub Pages and `/_next/*` to K8s

---

**Status**: Live
**Timeline**: 15 minutes (diagnosis + fix + deploy + verification)

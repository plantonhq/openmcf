---
title: "Planton Is Now on a Single Domain"
date: 2026-03-28
category: improvement
tags:
  - platform
  - console
excerpt: "The entire platform — marketing site, docs, and console — now lives at planton.ai. No more switching between subdomains."
author:
  - name: Swarup Donepudi
    title: Founder
---

Everything on Planton now lives under one domain. The marketing site, documentation, and console application are all served from `planton.ai`. The separate `console.planton.ai` subdomain is retired — if you have bookmarks or links pointing there, they automatically redirect to the equivalent page on `planton.ai`.

## What Changed

Previously, using Planton meant moving between two domains: `planton.ai` for the marketing site and docs, and `console.planton.ai` for the console application. This split created a disjointed experience — your browser treated them as separate sites, authentication cookies couldn't be shared, and the GitHub-style organization URLs (`planton.ai/{org}/...`) couldn't work across the boundary.

Now, a single edge routing layer decides where each request goes based on the URL path. Marketing pages (`/`, `/docs`, `/pricing`, `/blog`, `/changelog`, and others) are served from the static site. Everything else — `/dashboard`, `/login`, organization pages, and all console routes — is served from the console application. This happens transparently at the edge with no impact on page load time.

## What This Means for You

**GitHub-style URLs work end-to-end.** Organization URLs like `planton.ai/{org}/infra-hub` now resolve correctly without crossing domain boundaries.

**Logged-in users land on their dashboard.** When you visit `planton.ai` while logged in, the site detects your session and takes you to `/dashboard` instead of showing the marketing page.

**Old links still work.** Any URL on `console.planton.ai` permanently redirects (301) to the same path on `planton.ai`. Bookmarks, CI scripts, and shared links continue to work without changes.

**One re-login required.** Because the authentication cookie moved from `console.planton.ai` to `planton.ai`, existing sessions were invalidated. You'll need to sign in once after this change. All subsequent sessions use the new domain automatically.

## Why This Matters

A unified domain is not just a cosmetic change. Shared authentication cookies mean the marketing site and console can coordinate — like redirecting logged-in users to the dashboard. It also simplifies CORS, cookie security, and link sharing. Every URL you see, share, or bookmark is now on `planton.ai`.

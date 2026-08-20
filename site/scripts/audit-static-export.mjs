/**
 * Static-export audit crawler.
 *
 * Loads every HTML page of the built static export (out/) in a headless
 * browser and records what each page ACTUALLY fetches — the ground truth
 * that neither source-level nor built-output text analysis can give, because
 * a URL string sitting in a compiled chunk is not evidence that any page
 * renders it, and viewport-lazy media appears in no initial HTML at all.
 *
 * For every page it:
 *   1. navigates and waits for the network to settle,
 *   2. scrolls to the bottom of the page (triggering viewport-lazy loads),
 *   3. waits for the network to settle again,
 *   4. records every same-origin URL requested (the page's network manifest),
 *   5. captures a full-page screenshot.
 *
 * Outputs, under --out:
 *   manifests.json            route -> sorted same-origin request paths
 *   screenshots/<route>.png   full-page screenshot per route
 *
 * Two audit runs (e.g. before and after a cleanup) are comparable by
 * diffing manifests.json with /_next/ paths excluded (chunk hashes change
 * on every build) and by comparing screenshots.
 *
 * Usage:
 *   npx serve out -l 4175 &
 *   node scripts/audit-static-export.mjs --base http://localhost:4175 --out /tmp/site-audit
 */

import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { createRequire } from 'node:module';

const require = createRequire(import.meta.url);
const puppeteer = require('puppeteer');

const args = process.argv.slice(2);
const arg = (name, fallback) => {
  const i = args.indexOf(`--${name}`);
  return i >= 0 ? args[i + 1] : fallback;
};
const BASE = arg('base', 'http://localhost:4175');
const OUT_DIR = arg('out', '/tmp/site-audit');
const CONCURRENCY = Number(arg('concurrency', '4'));

const siteRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
const exportDir = path.join(siteRoot, 'out');

// Enumerate routes from the export's HTML files. Both out/foo.html and
// out/foo/index.html shapes occur; dedupe to one route each.
function enumerateRoutes() {
  const routes = new Set();
  const walk = (dir) => {
    for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
      const full = path.join(dir, entry.name);
      const rel = path.relative(exportDir, full);
      if (entry.isDirectory()) {
        if (['_next', '_site', '_pagefind', 'videos'].includes(rel.split(path.sep)[0])) continue;
        walk(full);
      } else if (entry.name.endsWith('.html')) {
        let route = '/' + rel.replace(/\.html$/, '').replace(/\/?index$/, '');
        if (route === '/404' || route === '/500') continue;
        routes.add(route === '' || route === '/' ? '/' : route.replace(/\/$/, '') || '/');
      }
    }
  };
  walk(exportDir);
  return [...routes].sort();
}

async function autoScroll(page) {
  await page.evaluate(async () => {
    await new Promise((resolve) => {
      let total = 0;
      const step = 800;
      const timer = setInterval(() => {
        window.scrollBy(0, step);
        total += step;
        if (total >= document.body.scrollHeight) {
          clearInterval(timer);
          resolve();
        }
      }, 60);
    });
  });
}

async function auditRoute(browser, route) {
  const page = await browser.newPage();
  await page.setViewport({ width: 1440, height: 900 });
  const requested = new Set();
  page.on('request', (req) => {
    const u = new URL(req.url());
    if (u.origin === new URL(BASE).origin) requested.add(u.pathname);
  });
  try {
    await page.goto(BASE + route, { waitUntil: 'networkidle2', timeout: 45000 });
    await autoScroll(page);
    await new Promise((r) => setTimeout(r, 800));
    const shotName = (route === '/' ? 'index' : route.slice(1).replace(/\//g, '__')) + '.png';
    await page.screenshot({ path: path.join(OUT_DIR, 'screenshots', shotName), fullPage: true });
    return { route, requests: [...requested].sort(), error: null };
  } catch (err) {
    return { route, requests: [...requested].sort(), error: String(err) };
  } finally {
    await page.close();
  }
}

const routes = enumerateRoutes();
fs.mkdirSync(path.join(OUT_DIR, 'screenshots'), { recursive: true });
console.log(`auditing ${routes.length} routes against ${BASE}`);

const browser = await puppeteer.launch({ headless: 'shell' });
const results = {};
let failed = 0;
const queue = [...routes];
await Promise.all(
  Array.from({ length: CONCURRENCY }, async () => {
    while (queue.length) {
      const route = queue.shift();
      const r = await auditRoute(browser, route);
      results[r.route] = { requests: r.requests, error: r.error };
      if (r.error) {
        failed++;
        console.error(`  FAIL ${r.route}: ${r.error}`);
      }
    }
  }),
);
await browser.close();

fs.writeFileSync(path.join(OUT_DIR, 'manifests.json'), JSON.stringify(results, null, 1));
const fetched = new Set(Object.values(results).flatMap((r) => r.requests));
console.log(`done: ${routes.length} routes, ${failed} failures, ${fetched.size} distinct same-origin paths fetched`);
console.log(`manifests: ${path.join(OUT_DIR, 'manifests.json')}`);

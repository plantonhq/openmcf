/**
 * THE single source of positioning truth for the website.
 *
 * The product positioning follows a three-level rule, and this file exists
 * so the rule cannot drift the way the last one did (the "Vercel for
 * backend" line was coined for Service Hub, escaped its scope, and spent
 * six months describing the whole product):
 *
 *   Level 1 — Planton, the umbrella. Never uses an analogy.
 *   Level 2 — each hub gets exactly one analogy, and the analogy never
 *             escapes its hub. Infra Hub is "Cursor for Cloud
 *             Infrastructure"; Service Hub is "Vercel for Backend, In
 *             Your Own Cloud".
 *
 * Any sentence on the site that states what Planton IS must read from
 * here — the same law src/data/pricing.ts applies to prices.
 *
 * Vocabulary: "Infra Chart" is the product noun; "template" is the plain
 * word that explains it. Pair them on first mention, and never capitalize
 * "Template" as if it were a product name.
 */

export const POSITIONING = {
  /** Level 1 — the umbrella. No analogy, ever. */
  umbrella: {
    tagline: 'The Self-Service Cloud Platform',
    sentence:
      'Planton turns your own cloud account into a self-service platform. AI designs the infrastructure, verifies the cost and permissions before anything is created, and publishes it as templates your whole team can deploy. Your services then ship onto that infrastructure straight from Git.',
  },

  /** Level 2 — one analogy per hub, scoped to that hub only. */
  infraHub: {
    name: 'Infra Hub',
    analogy: 'Cursor for Cloud Infrastructure',
    line: 'Describe what you need, watch it compose on a live canvas, see the cloud bill and the IAM policy before anything is created, deploy, and publish it as an Infra Chart — a template your team reuses.',
  },
  serviceHub: {
    name: 'Service Hub',
    analogy: 'Vercel for Backend, In Your Own Cloud',
    line: 'Connect a Git repository and every push becomes a running deployment — no pipeline YAML, no Dockerfile required, with results written back into GitHub checks and deployments.',
  },
} as const;

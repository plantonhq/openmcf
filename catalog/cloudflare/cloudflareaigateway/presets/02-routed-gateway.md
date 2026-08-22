# Routed gateway (cheap-first)

A gateway whose `cheap-first` route sends free-tier requests to a Workers AI model and everything else to GPT-4o, with the cheap arm falling back to the smart one on failure. Clients address the graph by route name; the graph decides where each request lands.

**When to use it:** tiered products where model cost should follow customer tier, or any cheap-model-first strategy with a smart fallback.

**What to change:** the `conditions` expression (a JSON document over request metadata), the two models and providers, and the graph shape itself. Remember: ANY edit to a route's `elements` replaces that route object -- the plan showing a replace is the designed behavior, not a mistake.

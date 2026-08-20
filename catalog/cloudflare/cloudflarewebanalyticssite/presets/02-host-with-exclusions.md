# Host with path exclusions

The measure-the-storefront-only shape: a hostname-identified site (works whether or not the site is behind Cloudflare) running the lite beacon, with admin and checkout paths excluded from measurement and everything else included. Order is positional -- the exclusion sits ahead of the catch-all. Remember the rules are write-only at the provider: once this manifest owns them, dashboard edits to the rule set will be overwritten on the next apply.

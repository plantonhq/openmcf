# Web-Tier Firewall

This preset creates a firewall for public web servers: HTTPS and HTTP open to the world, SSH restricted to a management network, all outbound traffic allowed, and Droplet membership driven by the `web` tag so every Droplet carrying that tag is protected automatically.

## When to Use

- Public-facing web servers, API backends, or reverse proxies
- Any tier where Droplets come and go and membership should follow a tag
- The internet-facing half of a classic web/database two-tier layout (pair with `02-database-tier`)

## Key Configuration Choices

- **Tag-based targeting** (`tags: [web]`) -- membership follows the tag automatically as Droplets are created and destroyed; DigitalOcean creates the tag implicitly on first use. Prefer tags over `dropletIds` for anything long-lived.
- **HTTPS and HTTP from everywhere** -- both address families (`0.0.0.0/0` and `::/0`).
- **SSH from a management CIDR only** -- never expose port 22 to the world in production.
- **Open egress** -- web tiers typically need unrestricted outbound (package mirrors, APIs, DNS). `portRange: all` is the provider's canonical spelling; writing `1-65535` reads back as `all` and creates a permanent diff.
- **ICMP outbound** -- icmp rules take no port range (the provider drops one if set).

## Placeholders to Replace

- `metadata.name` / `firewallName` -- your firewall's name.
- The SSH rule's `sourceAddresses` (`203.0.113.0/24` is a documentation example) -- your office/VPN CIDR.

## Related Presets

- **02-database-tier** -- the private half: inbound only from this tier's tag, restricted egress.

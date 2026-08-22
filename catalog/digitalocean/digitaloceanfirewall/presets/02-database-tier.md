# Database-Tier Firewall

This preset creates a locked-down firewall for self-managed database Droplets: PostgreSQL reachable only from Droplets tagged `web`, SSH only from a management network, and egress restricted to DNS and HTTPS (package updates). Everything else, in both directions, is denied by default.

## When to Use

- Self-managed databases (PostgreSQL, MySQL, Redis) running on Droplets
- Any internal tier that should never accept traffic from the internet
- The private half of a classic web/database two-tier layout (pair with `01-web-tier`)

## Key Configuration Choices

- **Tag-to-tag rules** (`sourceTags: [web]`) -- the database accepts connections from whatever Droplets carry the `web` tag, so scaling the web tier never touches firewall configuration.
- **Restricted egress** -- unlike the web tier, this tier's outbound is limited to DNS (udp/53) and HTTPS (tcp/443 for OS package mirrors). A database host has no business opening arbitrary outbound connections; this contains exfiltration paths.
- **Port 5432** -- PostgreSQL. Change to `3306` (MySQL), `6379` (Redis), or your engine's port.
- **SSH from a management CIDR only** -- same discipline as every tier.

## Placeholders to Replace

- `metadata.name` / `firewallName` -- your firewall's name.
- The SSH rule's `sourceAddresses` (`203.0.113.0/24` is a documentation example) -- your office/VPN CIDR.
- `portRange: "5432"` -- your database engine's port if not PostgreSQL.

## Related Presets

- **01-web-tier** -- the public half whose `web` tag this preset trusts.

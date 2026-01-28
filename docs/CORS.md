# CORS Configuration

OpenMCF handles CORS automatically. Most users don't need to configure anything.

## Quick Start

```bash
docker compose up
```

That's it! CORS is configured and working.

## When Do I Need to Configure CORS?

### ✅ You DON'T need configuration if:
- Running with `docker compose up`
- Accessing from `localhost:3000`
- Using default setup

### ⚙️ You MIGHT need configuration if:
- Using custom domains
- Behind a reverse proxy (Caddy/nginx)
- Corporate network with specific requirements

## Configuration Options

### Custom Domains

If accessing from custom domains:

```yaml
# docker-compose.yml
environment:
  - CORS_ALLOWED_ORIGINS=https://app.mycompany.com,https://admin.mycompany.com
```

### With Reverse Proxy

If using Caddy/nginx that handles CORS:

```yaml
# docker-compose.yml
environment:
  - ENABLE_CORS=false  # Let reverse proxy handle CORS
```

Then configure CORS in your reverse proxy.

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ENABLE_CORS` | `true` | Enable backend CORS |
| `CORS_ALLOWED_ORIGINS` | `localhost:3000,localhost:3001` | Allowed origins (comma-separated) |

## Common Issues

### "Access blocked by CORS policy"

**Solution:** Add your domain to allowed origins

```bash
CORS_ALLOWED_ORIGINS=https://yourdomain.com docker compose up
```

### "missing trailer" error

This is fixed in the latest version. Make sure you're using:
```bash
docker pull ghcr.io/plantonhq/openmcf:latest
```

### Duplicate CORS headers

If using a reverse proxy, disable backend CORS:
```bash
ENABLE_CORS=false docker compose up
```

## Technical Details

For maintainers and advanced users, see [CORS Architecture](_cursor/cors-architecture.md).

## Need Help?

- 📖 [Full Documentation](https://docs.openmcf.org)
- 💬 [GitHub Discussions](https://github.com/plantonhq/openmcf/discussions)
- 🐛 [Report an Issue](https://github.com/plantonhq/openmcf/issues)


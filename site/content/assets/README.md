# Content Assets

This folder contains assets for the planton.ai website that are synced to Cloudflare R2.

## Structure

```
assets/
├── _inbox/       # Drop zone for raw images (gitignored)
├── _rules/       # Cursor rules for asset management
└── images/       # Organized, processed images
```

## URL Mapping

Everything in this folder mirrors to `https://assets.planton.ai/site/`.

| Local Path                 | URL                                                       |
| -------------------------- | --------------------------------------------------------- |
| `images/service-hub/x.png` | `https://assets.planton.ai/site/images/service-hub/x.png` |

## Adding Images

See `_rules/README.md` for the complete workflow using the `@process-planton-ai-images` cursor rule.

Quick version:

1. Drop images in `_inbox/`
2. Invoke `@process-planton-ai-images` with context
3. Use the R2 URL in docs

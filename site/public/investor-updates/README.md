# Investor Updates Content

This directory contains investor updates for planton.ai/legal/investor-updates.

## File Naming Convention

Files must be named with a date prefix:

```
YYYY-MM-DD-slug.md
```

Examples:
- `2026-02-01-february-update.md`
- `2026-03-15-series-a-progress.md`
- `2026-04-01-quarterly-metrics-q1.md`

**Important**: Unlike the fastlane feature on donepudi.me, the date prefix is KEPT in the URL for investor updates. This makes chronological sorting obvious in shared links:

- URL: `/legal/investor-updates/2026-02-01-february-update`

## Frontmatter

Each update requires the following frontmatter:

```yaml
---
title: "February 2026 Update: Building Distribution"
date: "2026-02-01"
tags: [product, growth, fundraising]  # optional
excerpt: "Custom excerpt for preview"  # optional, auto-generated if not provided
---
```

### Required Fields

| Field | Description |
|-------|-------------|
| `title` | The update title (displayed in timeline and page header) |
| `date` | Date in YYYY-MM-DD format |

### Optional Fields

| Field | Description |
|-------|-------------|
| `tags` | Array of tags for categorization |
| `excerpt` | Custom excerpt for the timeline preview (auto-generated from content if not provided) |

## Content Guidelines

### Recommended Sections

1. **Summary/Highlights** - Key takeaways for busy readers
2. **Metrics** - MRR, customers, usage stats
3. **Wins** - What went well
4. **Challenges** - What's hard, what we're working through
5. **Use of Funds** - How we spent money this period
6. **Next Month** - What we're focused on
7. **Asks** - How investors can help (intros, advice, etc.)

### Tone

- **Honest**: Share both wins and challenges
- **Specific**: Use numbers when possible
- **Concise**: Respect readers' time
- **Forward-looking**: What's next, not just what happened

### Example Structure

```markdown
---
title: "February 2026 Update"
date: "2026-02-01"
tags: [product, metrics, growth]
---

## Highlights

- MRR: $X → $Y (+Z%)
- New customers: N
- Key win: [description]

## Metrics Dashboard

| Metric | Last Month | This Month | Change |
|--------|------------|------------|--------|
| MRR | $X | $Y | +Z% |
| Customers | A | B | +C |
| Churn | X% | Y% | -Z% |

## What Went Well

- [Win 1]
- [Win 2]

## Challenges

- [Challenge 1 and how we're addressing it]
- [Challenge 2]

## Use of Funds

- Engineering: $X (Y%)
- Marketing: $X (Y%)
- Operations: $X (Y%)

## Next Month Focus

1. [Priority 1]
2. [Priority 2]
3. [Priority 3]

## How You Can Help

- Intros to [type of person/company]
- Advice on [specific topic]
```

## Rendering

Updates are rendered with a simple Markdown processor that supports:

- Headers (H1, H2, H3)
- Bold and italic text
- Links
- Unordered lists

For more complex formatting, consider keeping content simple and clear.

## URL Structure

- Timeline: `planton.ai/legal/investor-updates`
- Individual update: `planton.ai/legal/investor-updates/YYYY-MM-DD-slug`

Each update can be:
- Expanded inline on the timeline page
- Opened in a new tab via direct URL
- Link copied for sharing

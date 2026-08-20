# Shared Platform Split

This preset creates a cost-center category where shared platform
costs (the platform team's tag plus AWS Support) are RE-ALLOCATED
across the product values proportionally to each product's own spend
— full chargeback, no orphan bucket.

## When to Use

- Chargeback where shared/platform costs must land on the products
  that benefit, not in a "shared" line nobody owns
- Allocating support plans, shared tooling, or central accounts

## What You Get

- Product values from their tags, a platform value composed with
  or-logic (tag OR the Support service), and unmatched spend in
  "shared"
- A PROPORTIONAL split-charge rule dissolving the platform value into
  the products

## Customize

- `method: EVEN` splits equally; `method: FIXED` takes a `parameters`
  entry with one percentage per target summing to 100
- A split source cannot be a target of another split rule (AWS
  rejects circular allocation)
- Order matters: put narrow product rules ABOVE broad shared rules —
  first match wins

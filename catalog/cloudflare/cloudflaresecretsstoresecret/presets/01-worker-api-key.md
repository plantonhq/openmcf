---
display_name: Worker API Key
---

# Shared provider API key

An AI provider's API key stored once and readable by both Workers and the AI Gateway -- the classic Bring-Your-Own-Keys shape. The value arrives as a managed-secret reference (never plaintext) and rotates in one place for every consumer.

**When to use it:** any credential more than one Worker (or the gateway) presents to an upstream service.

**What to change:** the secret's `name` (the contract consumers reference), the managed-secret reference in `value`, and the `scopes` -- keep them alphabetical, the spec enforces the order the API returns.

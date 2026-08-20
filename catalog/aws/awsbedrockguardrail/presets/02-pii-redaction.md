# PII Redaction

This preset creates a guardrail focused on sensitive-information handling:
contact details are masked (ANONYMIZE — the model still answers, with
`{NAME}`-style placeholders), while payment data, government identifiers,
and credentials are blocked outright. A custom regex masks internal ticket
identifiers as an example of organization-specific patterns.

## When to Use

- Assistants over support conversations, CRM data, or user-generated text
- Compliance postures that require PII never reach model logs or outputs

## Key Configuration Choices

- **ANONYMIZE for utility, BLOCK for liability.** Masked entities keep the
  conversation flowing; blocked entities (cards, SSNs, credentials) stop
  it — tune per entity as your data-protection rules require.
- **The regex arm is the extension point** for identifiers AWS's PII
  taxonomy doesn't know (employee IDs, order numbers, internal hostnames).

## After Deployment

Watch guardrail traces for match rates; flip noisy entities to
detect-only (`inputAction: NONE`) while tuning rather than removing them.

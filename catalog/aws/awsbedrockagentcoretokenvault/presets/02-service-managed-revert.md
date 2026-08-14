# Service-Managed Key (Revert)

This preset returns the region's token vault to AWS-owned encryption
— the default posture, and the ONLY way back from a customer-managed
key (destroying the component does not revert the setting).

## When to Use

- Reverting a customer-managed key before retiring/deleting that key
- Standing down a compliance posture that no longer applies

## What You Get

- The default vault back under AWS's service-managed key
- Safe to schedule the old CMK's deletion AFTER this applies

## Customize

- Nothing to customize beyond the region — this is the whole posture

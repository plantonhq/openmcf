---
display_name: MFA Hardening
---

# MFA hardening

The security-first shape: MFA required for every application by default (security keys or TOTP), unmatched requests denied, the dashboard locked read-only so IaC is the only write path, and the Access service key rotating every 30 days.

Two API-side rules to respect: `mfa_required_for_all_apps` needs at least one allowed authenticator AND an MFA session duration (both set here), and `ssh_piv_key` alone is rejected while any non-infrastructure application exists. The dashboard lock (`is_ui_read_only`) applies to everyone regardless of role -- the toggle reason tells them where changes live now.

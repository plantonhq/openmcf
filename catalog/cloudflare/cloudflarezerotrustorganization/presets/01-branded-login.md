# Branded login

The everyday shape: brand the Access login page and set the session default. The team domain (`auth_domain`) is set once and never touched casually -- changing it breaks every Access login and WARP enrollment at once. Everything unset (MFA policy, custom pages, key rotation) stays exactly as the dashboard has it: unset means unmanaged on this singleton.

Swap in your account id, team domain, and brand colors; apply mutates the organization the account already carries (Zero Trust onboarding created it), and destroy reverts nothing.

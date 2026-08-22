# Locked fleet baseline

The managed-fleet shape: users cannot turn WARP off (`switch_locked`), cannot unenroll (`allowed_to_leave: false`), and cannot switch modes; a disabled client self-reconnects after ten minutes, with the standard captive-portal grace for hotel and airport Wi-Fi. Roll this out only after split-tunnel routing is proven -- a bad route with a locked switch hits every device at once.

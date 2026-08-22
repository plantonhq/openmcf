# Developer LAN access

The engineering shape: a 30-minute LAN window on a /24, the lab network excluded from the tunnel, and `lab.internal` resolving against the lab resolver through this profile's own fallback list (full-replacement -- `localhost` and `home.arpa` are re-declared deliberately). Precedence 200 ranks after the contractor profile's 100 (lower wins), so an engineer flagged as both gets the tighter contractor profile.

# Cross-Account Accepter

The accepter side of a cross-account peering: the partner account's request arrives `pending-acceptance`, this instance accepts it by `pcx-` id (paste the id the requester shares, or wire another instance's output in same-account chains) and enables DNS resolution for its side. Destroying this instance abandons management — only the requester deletes the peering.

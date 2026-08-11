# Azure VPN Gateway Connection -- Operational Guide

Judgment that does not fit in field references.

## Provisioned is not connected -- believe it

A Succeeded deployment means ARM accepted the parameters, nothing
more. The tunnel reaches Connected only when the branch device
negotiates: matching pre-shared key, compatible proposals, reachable
endpoint. Debug a Connecting tunnel on the BRANCH side first (key,
ciphers, firewall) -- redeploying the Azure side rarely changes
anything.

## The link pin is permanent; the link set is not

Each tunnel's `vpnSiteLinkId` is fixed at creation -- repointing a
tunnel at a different site link means replacing that tunnel (cheap:
seconds of renegotiation, no gateway impact). ADDING a tunnel for a
branch's new second ISP is an in-place update; do that instead of
editing the existing one.

## Let Azure generate keys unless the branch cannot

An omitted `sharedKey` gets an Azure-generated value -- strong and
never in a manifest. Set the key only when the branch device demands a
pre-agreed one, and then REFERENCE a secret; the field is sensitive
and reference-capable for exactly that. One key per tunnel, both ends
identical.

## Pinned proposals fail closed

With `ipsecPolicies` set, the tunnel offers EXACTLY that proposal --
there is no fallback to Azure's defaults. A branch device that would
have negotiated fine against the default set can go dark after
pinning. Pin from the device vendor's documented matrix, not from a
generic hardening guide; with GCM IKE encryption, use the MATCHING GCM
integrity value (ARM's rule, enforced by the spec vocabulary).

## BGP is a three-object agreement

A BGP tunnel needs: the site link's `bgp` block (the branch speaker),
this link's `bgpEnabled` (fixed at creation), and agreement on ASNs
(the gateway is 65515; the branch must differ and avoid
65515-65520). Miss any leg and routes silently fall back to the
site's static `addressCidrs` -- traffic may still flow, masking the
misconfiguration until the static list drifts.

## Internet security changes the branch's default route

`internetSecurityEnabled: true` advertises 0.0.0.0/0 down the tunnel
-- the branch's internet traffic now hairpins through Azure. Enable it
deliberately, paired with a hub firewall (routing intent); enabling it
"because security" on a branch with local breakout expectations is a
support ticket, not a hardening win.

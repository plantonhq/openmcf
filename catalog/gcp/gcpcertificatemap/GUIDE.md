# GcpCertificateMap Guide

Operational judgment for running certificate maps as code — the things
the spec reference cannot tell you.

## Ship a PRIMARY entry or unmatched SNI fails cold

A map with only hostname entries hard-fails the TLS handshake for any
SNI it does not know — including health checkers and clients that send
no SNI at all. A `matcher: PRIMARY` entry (usually the wildcard or the
apex certificate) is the safety net; treat a PRIMARY-less production
map as a misconfiguration.

## Rotation happens in the certificate LIST, never the entry

The entry's certificates list is its only mutable surface. Rotate by
ADDING the replacement certificate to the list, waiting for it to
serve, then removing the old one — the entry keeps serving throughout.
Changing hostname/matcher/name REPLACES the entry, which is a window
where that hostname has no binding at the proxy.

## PROVISIONING certificates can be attached — serving is another matter

Certificate Manager accepts map entries that reference managed
certificates still in PROVISIONING; the handshake serves them only once
they turn ACTIVE. Wire the map and proxy first, complete DNS
authorization after — the order works, but watch the certificate's
managed_state before declaring the domain live.

## Detach from proxies BEFORE destroying

A GcpTargetHttpsProxy holding this map's URI fails every handshake the
moment the map (or its matching entry) is gone — GCP does not block the
delete. Teardown order: point the proxy at a replacement map (or direct
certificates), then destroy the map. `PREVENT` is the honest posture
while any proxy references it.

## Maps and certificates are global; regional certs do not belong

Certificate maps are a GLOBAL resource for global external ALBs and
classic HTTPS proxies. Regional Certificate Manager certificates
(regional ALBs) cannot ride a map — attach them directly on the
regional proxy instead.

## Fifteen certificates per entry is an API wall, not a suggestion

The per-entry cap is enforced here at manifest time (and by the API).
Needing more than a handful on ONE hostname usually signals migration
debt (RSA+ECDSA pairs plus staged rotations); finish the rotations
rather than raising the count.

## Teardown discipline

`DELETE` removes entries and map (detach from proxies first — above).
`ABANDON` keeps the map serving unmanaged — the escape for handing TLS
routing to another team. `PREVENT` protects live TLS routing from
accidental teardown; entries inherit the same policy.
